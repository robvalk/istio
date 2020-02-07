// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inject

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ghodss/yaml"
	"github.com/howeyc/fsnotify"
	"k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	clientset "k8s.io/client-go/kubernetes"

	meshconfig "istio.io/api/mesh/v1alpha1"
	"istio.io/istio/pilot/cmd"
	"istio.io/istio/pilot/cmd/pilot-agent/status"
	"istio.io/istio/pkg/log"
)

const proxyUIDAnnotation = "sidecar.istio.io/proxyUID"

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = v1beta1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1beta1.AddToScheme(runtimeScheme)
}

const (
	watchDebounceDelay = 100 * time.Millisecond
)

// Webhook implements a mutating webhook for automatic proxy injection.
type Webhook struct {
	mu                     sync.RWMutex
	sidecarConfig          *Config
	sidecarTemplateVersion string
	meshConfig             *meshconfig.MeshConfig

	healthCheckInterval time.Duration
	healthCheckFile     string

	server               *http.Server
	meshFile             string
	configFile           string
	keyCertWatcher       *fsnotify.Watcher
	configWatcher        *fsnotify.Watcher
	caFile               string
	certFile             string
	keyFile              string
	webhookConfigFile    string
	cert                 *tls.Certificate
	namespace            string
	deploymentName       string
	webhookConfigName    string
	clientset            clientset.Interface
	ownerRefs            []metav1.OwnerReference
	webhookConfiguration *admissionregistrationv1beta1.MutatingWebhookConfiguration
	manageWebhookConfig  bool
}

func loadConfig(injectFile, meshFile string) (*Config, *meshconfig.MeshConfig, error) {
	data, err := ioutil.ReadFile(injectFile)
	if err != nil {
		return nil, nil, err
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		log.Warnf("Failed to parse injectFile %s", string(data))
		return nil, nil, err
	}
	meshConfig, err := cmd.ReadMeshConfig(meshFile)
	if err != nil {
		return nil, nil, err
	}

	log.Infof("New configuration: sha256sum %x", sha256.Sum256(data))
	log.Infof("Policy: %v", c.Policy)
	log.Infof("AlwaysInjectSelector: %v", c.AlwaysInjectSelector)
	log.Infof("NeverInjectSelector: %v", c.NeverInjectSelector)
	log.Infof("Template: |\n  %v", strings.Replace(c.Template, "\n", "\n  ", -1))

	return &c, meshConfig, nil
}

// WebhookParameters configures parameters for the sidecar injection
// webhook.
type WebhookParameters struct {
	// ConfigFile is the path to the sidecar injection configuration file.
	ConfigFile string

	// MeshFile is the path to the mesh configuration file.
	MeshFile string

	// CACertFile is the path to the x509 CA bundle file.
	CACertFile string

	// CertFile is the path to the x509 certificate for https.
	CertFile string

	// KeyFile is the path to the x509 private key matching `CertFile`.
	KeyFile string

	// Port is the webhook port, e.g. typically 443 for https.
	Port int

	// WebhookConfigFile is the path to the mutatingwebhookconfiguration
	WebhookConfigFile string

	// Namespace is the namespace in which the deployment and service resides.
	Namespace string

	// Name of the webhook
	WebhookConfigName string

	// The webhook deployment name
	DeploymentName string

	// HealthCheckInterval configures how frequently the health check
	// file is updated. Value of zero disables the health check
	// update.
	HealthCheckInterval time.Duration

	// HealthCheckFile specifies the path to the health check file
	// that is periodically updated.
	HealthCheckFile string

	// ManageWebhookConfig determines whether the MutatingWebhookConfiguration
	// should be watched and updated by the webhook itself. This can be disabled,
	// in which case we simply assume the configuration to be correct. This is
	// helpful if you want to remove ClusterRole permissions from the injector.
	// NOTE: setting this to false requires the WebhookConfiguration to be created
	// and updated by something that is external to the injector, e.g. an operator
	ManageWebhookConfig bool

	Clientset clientset.Interface
}

// NewWebhook creates a new instance of a mutating webhook for automatic sidecar injection.
func NewWebhook(p WebhookParameters) (*Webhook, error) {
	sidecarConfig, meshConfig, err := loadConfig(p.ConfigFile, p.MeshFile)
	if err != nil {
		return nil, err
	}
	pair, err := tls.LoadX509KeyPair(p.CertFile, p.KeyFile)
	if err != nil {
		return nil, err
	}

	certKeyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	// watch the parent directory of the target files so we can catch
	// symlink updates of k8s ConfigMaps volumes.
	for _, file := range []string{p.ConfigFile, p.MeshFile, p.CertFile, p.KeyFile} {
		watchDir, _ := filepath.Split(file)
		if err := certKeyWatcher.Watch(watchDir); err != nil {
			return nil, fmt.Errorf("could not watch %v: %v", file, err)
		}
	}

	// configuration must be updated whenever the caBundle changes.
	// NOTE: Use a separate watcher to differentiate config/ca from cert/key updates. This is
	// useful to avoid unnecessary updates and, more importantly, makes its easier to more
	// accurately capture logs/metrics when files change.
	configWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if p.ManageWebhookConfig {
		// only actually watch the files if we manage the webhook config
		for _, file := range []string{p.CACertFile, p.WebhookConfigFile} {
			watchDir, _ := filepath.Split(file)
			if err := configWatcher.Watch(watchDir); err != nil {
				return nil, fmt.Errorf("could not watch %v: %v", file, err)
			}
		}
	}

	wh := &Webhook{
		server: &http.Server{
			Addr: fmt.Sprintf(":%v", p.Port),
		},
		sidecarConfig:          sidecarConfig,
		sidecarTemplateVersion: sidecarTemplateVersionHash(sidecarConfig.Template),
		meshConfig:             meshConfig,
		configFile:             p.ConfigFile,
		meshFile:               p.MeshFile,
		keyCertWatcher:         certKeyWatcher,
		configWatcher:          configWatcher,
		healthCheckInterval:    p.HealthCheckInterval,
		healthCheckFile:        p.HealthCheckFile,
		certFile:               p.CertFile,
		keyFile:                p.KeyFile,
		cert:                   &pair,
		caFile:                 p.CACertFile,
		webhookConfigFile:      p.WebhookConfigFile,
		deploymentName:         p.DeploymentName,
		webhookConfigName:      p.WebhookConfigName,
		manageWebhookConfig:    p.ManageWebhookConfig,
		namespace:              p.Namespace,
		clientset:              p.Clientset,
	}
	if wh.manageWebhookConfig {
		if registryPullerDeployment, err := wh.clientset.ExtensionsV1beta1().Deployments(wh.namespace).Get(wh.deploymentName, metav1.GetOptions{}); err != nil {
			log.Errorf("Could not find %s/%s deployment to set ownerRef. The mutatingwebhookconfiguration must be deleted manually",
				wh.namespace, wh.deploymentName)
		} else {
			wh.ownerRefs = []metav1.OwnerReference{
				*metav1.NewControllerRef(
					registryPullerDeployment,
					extensionsv1beta1.SchemeGroupVersion.WithKind("Deployment"),
				),
			}
		}
	}

	// mtls disabled because apiserver webhook cert usage is still TBD.
	wh.server.TLSConfig = &tls.Config{GetCertificate: wh.getCert}
	h := http.NewServeMux()
	h.HandleFunc("/inject", wh.serveInject)
	wh.server.Handler = h

	return wh, nil
}

func (wh *Webhook) stop() {
	wh.keyCertWatcher.Close() // nolint: errcheck
	wh.configWatcher.Close()  // nolint: errcheck
	wh.server.Close()         // nolint: errcheck
}

// Run implements the webhook server
func (wh *Webhook) Run(stop <-chan struct{}) {
	go func() {
		if err := wh.server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("admission webhook ListenAndServeTLS failed: %v", err)
		}
	}()
	defer wh.stop()

	webhookChangedCh := make(chan struct{})
	if wh.manageWebhookConfig {
		// Try to create the initial webhook configuration (if it doesn't
		// already exist). Setup a persistent monitor to reconcile the
		// configuration if the observed configuration doesn't match
		// the desired configuration.
		if err := wh.rebuildWebhookConfig(); err == nil {
			wh.createOrUpdateWebhookConfig()
		}
		webhookChangedCh = wh.monitorWebhookChanges(stop)
	}

	var healthC <-chan time.Time
	if wh.healthCheckInterval != 0 && wh.healthCheckFile != "" {
		t := time.NewTicker(wh.healthCheckInterval)
		healthC = t.C
		defer t.Stop()
	}
	var keyCertTimerC <-chan time.Time
	var configTimerC <-chan time.Time

	for {
		select {
		case <-keyCertTimerC:
			keyCertTimerC = nil
			sidecarConfig, meshConfig, err := loadConfig(wh.configFile, wh.meshFile)
			if err != nil {
				log.Errorf("update error: %v", err)
				break
			}

			version := sidecarTemplateVersionHash(sidecarConfig.Template)
			pair, err := tls.LoadX509KeyPair(wh.certFile, wh.keyFile)
			if err != nil {
				log.Errorf("reload cert error: %v", err)
				break
			}
			wh.mu.Lock()
			wh.sidecarConfig = sidecarConfig
			wh.sidecarTemplateVersion = version
			wh.meshConfig = meshConfig
			wh.cert = &pair
			wh.mu.Unlock()
		case <-configTimerC:
			configTimerC = nil

			// rebuild the desired configuration and reconcile with the
			// existing configuration.
			if err := wh.rebuildWebhookConfig(); err == nil {
				wh.createOrUpdateWebhookConfig()
			}
		case <-webhookChangedCh:
			// reconcile the desired configuration
			wh.createOrUpdateWebhookConfig()
		case event := <-wh.keyCertWatcher.Event:
			// use a timer to debounce configuration updates
			if (event.IsModify() || event.IsCreate()) && keyCertTimerC == nil {
				keyCertTimerC = time.After(watchDebounceDelay)
			}
		case event := <-wh.configWatcher.Event:
			// use a timer to debounce configuration updates
			if (event.IsModify() || event.IsCreate()) && configTimerC == nil {
				configTimerC = time.After(watchDebounceDelay)
			}
		case err := <-wh.keyCertWatcher.Error:
			log.Errorf("keyCertWatcher error: %v", err)
		case err := <-wh.configWatcher.Error:
			log.Errorf("configWatcher error: %v", err)
		case <-healthC:
			content := []byte(`ok`)
			if err := ioutil.WriteFile(wh.healthCheckFile, content, 0644); err != nil {
				log.Errorf("Health check update of %q failed: %v", wh.healthCheckFile, err)
			}
		case <-stop:
			return
		}
	}
}

func (wh *Webhook) getCert(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	wh.mu.Lock()
	defer wh.mu.Unlock()
	return wh.cert, nil
}

// It would be great to use https://github.com/mattbaird/jsonpatch to
// generate RFC6902 JSON patches. Unfortunately, it doesn't produce
// correct patches for object removal. Fortunately, our patching needs
// are fairly simple so generating them manually isn't horrible (yet).
type rfc6902PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// JSONPatch `remove` is applied sequentially. Remove items in reverse
// order to avoid renumbering indices.
func removeContainers(containers []corev1.Container, removed []string, path string) (patch []rfc6902PatchOperation) {
	names := map[string]bool{}
	for _, name := range removed {
		names[name] = true
	}
	for i := len(containers) - 1; i >= 0; i-- {
		if _, ok := names[containers[i].Name]; ok {
			patch = append(patch, rfc6902PatchOperation{
				Op:   "remove",
				Path: fmt.Sprintf("%v/%v", path, i),
			})
		}
	}
	return patch
}

func removeVolumes(volumes []corev1.Volume, removed []string, path string) (patch []rfc6902PatchOperation) {
	names := map[string]bool{}
	for _, name := range removed {
		names[name] = true
	}
	for i := len(volumes) - 1; i >= 0; i-- {
		if _, ok := names[volumes[i].Name]; ok {
			patch = append(patch, rfc6902PatchOperation{
				Op:   "remove",
				Path: fmt.Sprintf("%v/%v", path, i),
			})
		}
	}
	return patch
}

func removeImagePullSecrets(imagePullSecrets []corev1.LocalObjectReference, removed []string, path string) (patch []rfc6902PatchOperation) {
	names := map[string]bool{}
	for _, name := range removed {
		names[name] = true
	}
	for i := len(imagePullSecrets) - 1; i >= 0; i-- {
		if _, ok := names[imagePullSecrets[i].Name]; ok {
			patch = append(patch, rfc6902PatchOperation{
				Op:   "remove",
				Path: fmt.Sprintf("%v/%v", path, i),
			})
		}
	}
	return patch
}

func addContainer(target, added []corev1.Container, basePath string) (patch []rfc6902PatchOperation) {
	saJwtSecretMountName := ""
	var saJwtSecretMount corev1.VolumeMount
	// find service account secret volume mount(/var/run/secrets/kubernetes.io/serviceaccount,
	// https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/#service-account-automation) from app container
	for _, add := range target {
		for _, vmount := range add.VolumeMounts {
			if vmount.MountPath == "/var/run/secrets/kubernetes.io/serviceaccount" {
				saJwtSecretMountName = vmount.Name
				saJwtSecretMount = vmount
			}
		}
	}
	first := len(target) == 0
	var value interface{}
	for _, add := range added {
		if add.Name == sidecarContainerName && saJwtSecretMountName != "" {
			// add service account secret volume mount(/var/run/secrets/kubernetes.io/serviceaccount,
			// https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/#service-account-automation) to istio-proxy container,
			// so that envoy could fetch/pass k8s sa jwt and pass to sds server, which will be used to request workload identity for the pod.
			add.VolumeMounts = append(add.VolumeMounts, saJwtSecretMount)
		}
		value = add
		path := basePath
		if first {
			first = false
			value = []corev1.Container{add}
		} else {
			path = path + "/-"
		}
		patch = append(patch, rfc6902PatchOperation{
			Op:    "add",
			Path:  path,
			Value: value,
		})
	}
	return patch
}

func addSecurityContext(target *corev1.PodSecurityContext, basePath string) (patch []rfc6902PatchOperation) {
	patch = append(patch, rfc6902PatchOperation{
		Op:    "add",
		Path:  basePath,
		Value: target,
	})
	return patch
}

func addVolume(target, added []corev1.Volume, basePath string) (patch []rfc6902PatchOperation) {
	first := len(target) == 0
	var value interface{}
	for _, add := range added {
		value = add
		path := basePath
		if first {
			first = false
			value = []corev1.Volume{add}
		} else {
			path = path + "/-"
		}
		patch = append(patch, rfc6902PatchOperation{
			Op:    "add",
			Path:  path,
			Value: value,
		})
	}
	return patch
}

func addImagePullSecrets(target, added []corev1.LocalObjectReference, basePath string) (patch []rfc6902PatchOperation) {
	first := len(target) == 0
	var value interface{}
	for _, add := range added {
		value = add
		path := basePath
		if first {
			first = false
			value = []corev1.LocalObjectReference{add}
		} else {
			path = path + "/-"
		}
		patch = append(patch, rfc6902PatchOperation{
			Op:    "add",
			Path:  path,
			Value: value,
		})
	}
	return patch
}

func addPodDNSConfig(target *corev1.PodDNSConfig, basePath string) (patch []rfc6902PatchOperation) {
	patch = append(patch, rfc6902PatchOperation{
		Op:    "add",
		Path:  basePath,
		Value: target,
	})
	return patch
}

// escape JSON Pointer value per https://tools.ietf.org/html/rfc6901
func escapeJSONPointerValue(in string) string {
	step := strings.Replace(in, "~", "~0", -1)
	return strings.Replace(step, "/", "~1", -1)
}

func updateAnnotation(target map[string]string, added map[string]string) (patch []rfc6902PatchOperation) {
	if target == nil {
		target = map[string]string{}
		patch = append(patch, rfc6902PatchOperation{
			Op:    "add",
			Path:  "/metadata/annotations",
			Value: added,
		})
	} else {
		for key, value := range added {
			op := "add"
			if target[key] != "" {
				op = "replace"
			}
			patch = append(patch, rfc6902PatchOperation{
				Op:    op,
				Path:  "/metadata/annotations/" + escapeJSONPointerValue(key),
				Value: value,
			})
		}
	}
	return patch
}

func createPatch(pod *corev1.Pod, prevStatus *SidecarInjectionStatus, annotations map[string]string, sic *SidecarInjectionSpec) ([]byte, error) {
	var patch []rfc6902PatchOperation

	// Remove any containers previously injected by kube-inject using
	// container and volume name as unique key for removal.
	patch = append(patch, removeContainers(pod.Spec.InitContainers, prevStatus.InitContainers, "/spec/initContainers")...)
	patch = append(patch, removeContainers(pod.Spec.Containers, prevStatus.Containers, "/spec/containers")...)
	patch = append(patch, removeVolumes(pod.Spec.Volumes, prevStatus.Volumes, "/spec/volumes")...)
	patch = append(patch, removeImagePullSecrets(pod.Spec.ImagePullSecrets, prevStatus.ImagePullSecrets, "/spec/imagePullSecrets")...)

	rewrite := ShouldRewriteAppProbers(sic)
	addAppProberCmd := func() {
		if !rewrite {
			return
		}
		sidecar := FindSidecar(sic.Containers)
		if sidecar == nil {
			log.Errorf("sidecar not found in the template, skip addAppProberCmd")
			return
		}
		// We don't have to escape json encoding here when using golang libraries.
		if prober := DumpAppProbers(&pod.Spec); prober != "" {
			sidecar.Env = append(sidecar.Env, corev1.EnvVar{Name: status.KubeAppProberEnvName, Value: prober})
		}
	}
	addAppProberCmd()

	patch = append(patch, addContainer(pod.Spec.InitContainers, sic.InitContainers, "/spec/initContainers")...)
	patch = append(patch, addContainer(pod.Spec.Containers, sic.Containers, "/spec/containers")...)
	patch = append(patch, addVolume(pod.Spec.Volumes, sic.Volumes, "/spec/volumes")...)
	patch = append(patch, addImagePullSecrets(pod.Spec.ImagePullSecrets, sic.ImagePullSecrets, "/spec/imagePullSecrets")...)

	if sic.DNSConfig != nil {
		patch = append(patch, addPodDNSConfig(sic.DNSConfig, "/spec/dnsConfig")...)
	}

	if pod.Spec.SecurityContext != nil {
		patch = append(patch, addSecurityContext(pod.Spec.SecurityContext, "/spec/securityContext")...)
	}

	patch = append(patch, updateAnnotation(pod.Annotations, annotations)...)

	if rewrite {
		patch = append(patch, createProbeRewritePatch(&pod.Spec, sic)...)
	}

	return json.Marshal(patch)
}

// Retain deprecated hardcoded container and volumes names to aid in
// backwards compatible migration to the new SidecarInjectionStatus.
var (
	initContainerName    = "istio-init"
	sidecarContainerName = "istio-proxy"

	legacyInitContainerNames = []string{initContainerName, "enable-core-dump"}
	legacyContainerNames     = []string{sidecarContainerName}
	legacyVolumeNames        = []string{"istio-certs", "istio-envoy"}
)

func injectionStatus(pod *corev1.Pod) *SidecarInjectionStatus {
	var statusBytes []byte
	if pod.ObjectMeta.Annotations != nil {
		if value, ok := pod.ObjectMeta.Annotations[annotationStatus.name]; ok {
			statusBytes = []byte(value)
		}
	}

	// default case when injected pod has explicit status
	var status SidecarInjectionStatus
	if err := json.Unmarshal(statusBytes, &status); err == nil {
		// heuristic assumes status is valid if any of the resource
		// lists is non-empty.
		if len(status.InitContainers) != 0 ||
			len(status.Containers) != 0 ||
			len(status.Volumes) != 0 ||
			len(status.ImagePullSecrets) != 0 {
			return &status
		}
	}

	// backwards compatibility case when injected pod has legacy
	// status. Infer status from the list of legacy hardcoded
	// container and volume names.
	return &SidecarInjectionStatus{
		InitContainers: legacyInitContainerNames,
		Containers:     legacyContainerNames,
		Volumes:        legacyVolumeNames,
	}
}

func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{Result: &metav1.Status{Message: err.Error()}}
}

func (wh *Webhook) inject(ar *v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	req := ar.Request
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		log.Errorf("Could not unmarshal raw object: %v %s", err,
			string(req.Object.Raw))
		return toAdmissionResponse(err)
	}

	// Deal with potential empty fields, e.g., when the pod is created by a deployment
	podName := potentialPodName(&pod.ObjectMeta)
	if pod.ObjectMeta.Namespace == "" {
		pod.ObjectMeta.Namespace = req.Namespace
	}

	log.Infof("AdmissionReview for Kind=%v Namespace=%v Name=%v (%v) UID=%v Rfc6902PatchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, podName, req.UID, req.Operation, req.UserInfo)
	log.Debugf("Object: %v", string(req.Object.Raw))
	log.Debugf("OldObject: %v", string(req.OldObject.Raw))

	partialInjection := false
	if !injectRequired(ignoredNamespaces, wh.sidecarConfig, &pod.Spec, &pod.ObjectMeta) {
		if wasInjectedThroughIstioctl(&pod) {
			log.Infof("Performing partial injection into pre-injected pod %s/%s (injecting Multus annotation and runAsUser id)", pod.ObjectMeta.Namespace, podName)
			partialInjection = true
		} else {
			log.Infof("Skipping %s/%s due to policy check", pod.ObjectMeta.Namespace, podName)
			return &v1beta1.AdmissionResponse{
				Allowed: true,
			}
		}
	}

	// due to bug https://github.com/kubernetes/kubernetes/issues/57923,
	// k8s sa jwt token volume mount file is only accessible to root user, not istio-proxy(the user that istio proxy runs as).
	// workaround by https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#set-the-security-context-for-a-pod
	if wh.meshConfig.EnableSdsTokenMount && wh.meshConfig.SdsUdsPath != "" {
		var grp = int64(1337)
		pod.Spec.SecurityContext = &corev1.PodSecurityContext{
			FSGroup: &grp,
		}
	}

	proxyUID, err := getProxyUID(pod)
	if err != nil {
		log.Infof("Could not get proxyUID from annotation: %v", err)
	}
	if proxyUID == nil {
		if pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.RunAsUser != nil {
			uid := uint64(*pod.Spec.SecurityContext.RunAsUser) + 1
			proxyUID = &uid
		}
		for _, c := range pod.Spec.Containers {
			if c.SecurityContext != nil && c.SecurityContext.RunAsUser != nil {
				uid := uint64(*c.SecurityContext.RunAsUser) + 1
				if proxyUID == nil || uid > *proxyUID {
					proxyUID = &uid
				}
			}
		}
	}
	if proxyUID == nil {
		uid := DefaultSidecarProxyUID
		proxyUID = &uid
	}

	spec, status, err := injectionData(wh.sidecarConfig.Template, wh.sidecarTemplateVersion, &pod.ObjectMeta, &pod.Spec, &pod.ObjectMeta, wh.meshConfig.DefaultConfig, wh.meshConfig, *proxyUID) // nolint: lll
	if err != nil {
		log.Infof("Injection data: err=%v spec=%v\n", err, status)
		return toAdmissionResponse(err)
	}

	var patchBytes []byte
	if partialInjection {
		patchBytes, err = createPartialPatch(&pod, spec.Annotations, *proxyUID)
	} else {
		replaceProxyRunAsUserID(spec, *proxyUID)

		annotations := map[string]string{annotationStatus.name: status}
		for k, v := range spec.Annotations {
			annotations[k] = v
		}
		patchBytes, err = createPatch(&pod, injectionStatus(&pod), annotations, spec)
	}

	if err != nil {
		log.Infof("AdmissionResponse: err=%v spec=%v\n", err, spec)

		return toAdmissionResponse(err)
	}

	log.Infof("AdmissionResponse: patch=%v\n", string(patchBytes))

	reviewResponse := v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}
	return &reviewResponse
}

func wasInjectedThroughIstioctl(pod *corev1.Pod) bool {
	_, found := pod.Annotations[annotationStatus.name]
	return found
}

func replaceProxyRunAsUserID(spec *SidecarInjectionSpec, proxyUID uint64) {
	for i, c := range spec.InitContainers {
		if c.Name == initContainerName {
			for j, arg := range c.Args {
				if arg == "-u" {
					spec.InitContainers[i].Args[j+1] = strconv.FormatUint(proxyUID, 10)
					break
				}
			}
			break
		}
	}
	for i, c := range spec.Containers {
		if c.Name == sidecarContainerName {
			if c.SecurityContext == nil {
				securityContext := corev1.SecurityContext{}
				spec.Containers[i].SecurityContext = &securityContext
			}
			proxyUIDasInt64 := int64(proxyUID)
			spec.Containers[i].SecurityContext.RunAsUser = &proxyUIDasInt64
			break
		}
	}
}

func createPartialPatch(pod *corev1.Pod, annotations map[string]string, proxyUID uint64) ([]byte, error) {
	var patch []rfc6902PatchOperation
	patch = append(patch, patchProxyRunAsUserID(pod, proxyUID)...)
	patch = append(patch, updateAnnotation(pod.Annotations, annotations)...)
	return json.Marshal(patch)
}

func patchProxyRunAsUserID(pod *corev1.Pod, proxyUID uint64) (patch []rfc6902PatchOperation) {
	for i, c := range pod.Spec.InitContainers {
		if c.Name == initContainerName {
			for j, arg := range c.Args {
				if arg == "-u" {
					patch = append(patch, rfc6902PatchOperation{
						Op:    "replace",
						Path:  fmt.Sprintf("/spec/initContainers/%d/args/%d", i, j+1), // j+1 because the uid is the next argument (after -u)
						Value: strconv.FormatUint(proxyUID, 10),
					})
					break
				}
			}
			break
		}
	}

	for i, c := range pod.Spec.Containers {
		if c.Name == sidecarContainerName {
			if c.SecurityContext == nil {
				proxyUIDasInt64 := int64(proxyUID)
				securityContext := corev1.SecurityContext{
					RunAsUser: &proxyUIDasInt64,
				}
				patch = append(patch, rfc6902PatchOperation{
					Op:    "add",
					Path:  fmt.Sprintf("/spec/containers/%d/securityContext", i),
					Value: securityContext,
				})
			} else if c.SecurityContext.RunAsUser == nil {
				patch = append(patch, rfc6902PatchOperation{
					Op:    "add",
					Path:  fmt.Sprintf("/spec/containers/%d/securityContext/runAsUser", i),
					Value: proxyUID,
				})
			} else {
				patch = append(patch, rfc6902PatchOperation{
					Op:    "replace",
					Path:  fmt.Sprintf("/spec/containers/%d/securityContext/runAsUser", i),
					Value: proxyUID,
				})
			}
			break
		}
	}

	return patch
}

func getProxyUID(pod corev1.Pod) (*uint64, error) {
	if pod.Annotations != nil {
		if annotationValue, found := pod.Annotations[proxyUIDAnnotation]; found {
			proxyUID, err := strconv.ParseUint(annotationValue, 10, 64)
			if err != nil {
				return nil, err
			}
			return &proxyUID, nil
		}
	}
	return nil, nil
}

func (wh *Webhook) serveInject(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		log.Errorf("no body found")
		http.Error(w, "no body found", http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Errorf("contentType=%s, expect application/json", contentType)
		http.Error(w, "invalid Content-Type, want `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	var reviewResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		log.Errorf("Could not decode body: %v", err)
		reviewResponse = toAdmissionResponse(err)
	} else {
		reviewResponse = wh.inject(&ar)
	}

	response := v1beta1.AdmissionReview{}
	if reviewResponse != nil {
		response.Response = reviewResponse
		if ar.Request != nil {
			response.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(response)
	if err != nil {
		log.Errorf("Could not encode response: %v", err)
		http.Error(w, fmt.Sprintf("could encode response: %v", err), http.StatusInternalServerError)
	}
	if _, err := w.Write(resp); err != nil {
		log.Errorf("Could not write response: %v", err)
		http.Error(w, fmt.Sprintf("could write response: %v", err), http.StatusInternalServerError)
	}
}