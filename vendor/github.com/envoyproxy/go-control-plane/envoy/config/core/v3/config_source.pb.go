// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/config/core/v3/config_source.proto

package envoy_config_core_v3

import (
	fmt "fmt"
	_ "github.com/cncf/udpa/go/udpa/annotations"
	_ "github.com/envoyproxy/go-control-plane/envoy/annotations"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ApiVersion int32

const (
	ApiVersion_AUTO ApiVersion = 0
	ApiVersion_V2   ApiVersion = 1
	ApiVersion_V3   ApiVersion = 2
)

var ApiVersion_name = map[int32]string{
	0: "AUTO",
	1: "V2",
	2: "V3",
}

var ApiVersion_value = map[string]int32{
	"AUTO": 0,
	"V2":   1,
	"V3":   2,
}

func (x ApiVersion) String() string {
	return proto.EnumName(ApiVersion_name, int32(x))
}

func (ApiVersion) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a3116f5d2bd13e64, []int{0}
}

type ApiConfigSource_ApiType int32

const (
	ApiConfigSource_hidden_envoy_deprecated_UNSUPPORTED_REST_LEGACY ApiConfigSource_ApiType = 0 // Deprecated: Do not use.
	ApiConfigSource_REST                                            ApiConfigSource_ApiType = 1
	ApiConfigSource_GRPC                                            ApiConfigSource_ApiType = 2
	ApiConfigSource_DELTA_GRPC                                      ApiConfigSource_ApiType = 3
)

var ApiConfigSource_ApiType_name = map[int32]string{
	0: "hidden_envoy_deprecated_UNSUPPORTED_REST_LEGACY",
	1: "REST",
	2: "GRPC",
	3: "DELTA_GRPC",
}

var ApiConfigSource_ApiType_value = map[string]int32{
	"hidden_envoy_deprecated_UNSUPPORTED_REST_LEGACY": 0,
	"REST":       1,
	"GRPC":       2,
	"DELTA_GRPC": 3,
}

func (x ApiConfigSource_ApiType) String() string {
	return proto.EnumName(ApiConfigSource_ApiType_name, int32(x))
}

func (ApiConfigSource_ApiType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_a3116f5d2bd13e64, []int{0, 0}
}

type ApiConfigSource struct {
	ApiType                   ApiConfigSource_ApiType `protobuf:"varint,1,opt,name=api_type,json=apiType,proto3,enum=envoy.config.core.v3.ApiConfigSource_ApiType" json:"api_type,omitempty"`
	TransportApiVersion       ApiVersion              `protobuf:"varint,8,opt,name=transport_api_version,json=transportApiVersion,proto3,enum=envoy.config.core.v3.ApiVersion" json:"transport_api_version,omitempty"`
	ClusterNames              []string                `protobuf:"bytes,2,rep,name=cluster_names,json=clusterNames,proto3" json:"cluster_names,omitempty"`
	GrpcServices              []*GrpcService          `protobuf:"bytes,4,rep,name=grpc_services,json=grpcServices,proto3" json:"grpc_services,omitempty"`
	RefreshDelay              *duration.Duration      `protobuf:"bytes,3,opt,name=refresh_delay,json=refreshDelay,proto3" json:"refresh_delay,omitempty"`
	RequestTimeout            *duration.Duration      `protobuf:"bytes,5,opt,name=request_timeout,json=requestTimeout,proto3" json:"request_timeout,omitempty"`
	RateLimitSettings         *RateLimitSettings      `protobuf:"bytes,6,opt,name=rate_limit_settings,json=rateLimitSettings,proto3" json:"rate_limit_settings,omitempty"`
	SetNodeOnFirstMessageOnly bool                    `protobuf:"varint,7,opt,name=set_node_on_first_message_only,json=setNodeOnFirstMessageOnly,proto3" json:"set_node_on_first_message_only,omitempty"`
	XXX_NoUnkeyedLiteral      struct{}                `json:"-"`
	XXX_unrecognized          []byte                  `json:"-"`
	XXX_sizecache             int32                   `json:"-"`
}

func (m *ApiConfigSource) Reset()         { *m = ApiConfigSource{} }
func (m *ApiConfigSource) String() string { return proto.CompactTextString(m) }
func (*ApiConfigSource) ProtoMessage()    {}
func (*ApiConfigSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_a3116f5d2bd13e64, []int{0}
}

func (m *ApiConfigSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ApiConfigSource.Unmarshal(m, b)
}
func (m *ApiConfigSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ApiConfigSource.Marshal(b, m, deterministic)
}
func (m *ApiConfigSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ApiConfigSource.Merge(m, src)
}
func (m *ApiConfigSource) XXX_Size() int {
	return xxx_messageInfo_ApiConfigSource.Size(m)
}
func (m *ApiConfigSource) XXX_DiscardUnknown() {
	xxx_messageInfo_ApiConfigSource.DiscardUnknown(m)
}

var xxx_messageInfo_ApiConfigSource proto.InternalMessageInfo

func (m *ApiConfigSource) GetApiType() ApiConfigSource_ApiType {
	if m != nil {
		return m.ApiType
	}
	return ApiConfigSource_hidden_envoy_deprecated_UNSUPPORTED_REST_LEGACY
}

func (m *ApiConfigSource) GetTransportApiVersion() ApiVersion {
	if m != nil {
		return m.TransportApiVersion
	}
	return ApiVersion_AUTO
}

func (m *ApiConfigSource) GetClusterNames() []string {
	if m != nil {
		return m.ClusterNames
	}
	return nil
}

func (m *ApiConfigSource) GetGrpcServices() []*GrpcService {
	if m != nil {
		return m.GrpcServices
	}
	return nil
}

func (m *ApiConfigSource) GetRefreshDelay() *duration.Duration {
	if m != nil {
		return m.RefreshDelay
	}
	return nil
}

func (m *ApiConfigSource) GetRequestTimeout() *duration.Duration {
	if m != nil {
		return m.RequestTimeout
	}
	return nil
}

func (m *ApiConfigSource) GetRateLimitSettings() *RateLimitSettings {
	if m != nil {
		return m.RateLimitSettings
	}
	return nil
}

func (m *ApiConfigSource) GetSetNodeOnFirstMessageOnly() bool {
	if m != nil {
		return m.SetNodeOnFirstMessageOnly
	}
	return false
}

type AggregatedConfigSource struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AggregatedConfigSource) Reset()         { *m = AggregatedConfigSource{} }
func (m *AggregatedConfigSource) String() string { return proto.CompactTextString(m) }
func (*AggregatedConfigSource) ProtoMessage()    {}
func (*AggregatedConfigSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_a3116f5d2bd13e64, []int{1}
}

func (m *AggregatedConfigSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AggregatedConfigSource.Unmarshal(m, b)
}
func (m *AggregatedConfigSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AggregatedConfigSource.Marshal(b, m, deterministic)
}
func (m *AggregatedConfigSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AggregatedConfigSource.Merge(m, src)
}
func (m *AggregatedConfigSource) XXX_Size() int {
	return xxx_messageInfo_AggregatedConfigSource.Size(m)
}
func (m *AggregatedConfigSource) XXX_DiscardUnknown() {
	xxx_messageInfo_AggregatedConfigSource.DiscardUnknown(m)
}

var xxx_messageInfo_AggregatedConfigSource proto.InternalMessageInfo

type SelfConfigSource struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SelfConfigSource) Reset()         { *m = SelfConfigSource{} }
func (m *SelfConfigSource) String() string { return proto.CompactTextString(m) }
func (*SelfConfigSource) ProtoMessage()    {}
func (*SelfConfigSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_a3116f5d2bd13e64, []int{2}
}

func (m *SelfConfigSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SelfConfigSource.Unmarshal(m, b)
}
func (m *SelfConfigSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SelfConfigSource.Marshal(b, m, deterministic)
}
func (m *SelfConfigSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SelfConfigSource.Merge(m, src)
}
func (m *SelfConfigSource) XXX_Size() int {
	return xxx_messageInfo_SelfConfigSource.Size(m)
}
func (m *SelfConfigSource) XXX_DiscardUnknown() {
	xxx_messageInfo_SelfConfigSource.DiscardUnknown(m)
}

var xxx_messageInfo_SelfConfigSource proto.InternalMessageInfo

type RateLimitSettings struct {
	MaxTokens            *wrappers.UInt32Value `protobuf:"bytes,1,opt,name=max_tokens,json=maxTokens,proto3" json:"max_tokens,omitempty"`
	FillRate             *wrappers.DoubleValue `protobuf:"bytes,2,opt,name=fill_rate,json=fillRate,proto3" json:"fill_rate,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *RateLimitSettings) Reset()         { *m = RateLimitSettings{} }
func (m *RateLimitSettings) String() string { return proto.CompactTextString(m) }
func (*RateLimitSettings) ProtoMessage()    {}
func (*RateLimitSettings) Descriptor() ([]byte, []int) {
	return fileDescriptor_a3116f5d2bd13e64, []int{3}
}

func (m *RateLimitSettings) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RateLimitSettings.Unmarshal(m, b)
}
func (m *RateLimitSettings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RateLimitSettings.Marshal(b, m, deterministic)
}
func (m *RateLimitSettings) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RateLimitSettings.Merge(m, src)
}
func (m *RateLimitSettings) XXX_Size() int {
	return xxx_messageInfo_RateLimitSettings.Size(m)
}
func (m *RateLimitSettings) XXX_DiscardUnknown() {
	xxx_messageInfo_RateLimitSettings.DiscardUnknown(m)
}

var xxx_messageInfo_RateLimitSettings proto.InternalMessageInfo

func (m *RateLimitSettings) GetMaxTokens() *wrappers.UInt32Value {
	if m != nil {
		return m.MaxTokens
	}
	return nil
}

func (m *RateLimitSettings) GetFillRate() *wrappers.DoubleValue {
	if m != nil {
		return m.FillRate
	}
	return nil
}

type ConfigSource struct {
	// Types that are valid to be assigned to ConfigSourceSpecifier:
	//	*ConfigSource_Path
	//	*ConfigSource_ApiConfigSource
	//	*ConfigSource_Ads
	//	*ConfigSource_Self
	ConfigSourceSpecifier isConfigSource_ConfigSourceSpecifier `protobuf_oneof:"config_source_specifier"`
	InitialFetchTimeout   *duration.Duration                   `protobuf:"bytes,4,opt,name=initial_fetch_timeout,json=initialFetchTimeout,proto3" json:"initial_fetch_timeout,omitempty"`
	ResourceApiVersion    ApiVersion                           `protobuf:"varint,6,opt,name=resource_api_version,json=resourceApiVersion,proto3,enum=envoy.config.core.v3.ApiVersion" json:"resource_api_version,omitempty"`
	XXX_NoUnkeyedLiteral  struct{}                             `json:"-"`
	XXX_unrecognized      []byte                               `json:"-"`
	XXX_sizecache         int32                                `json:"-"`
}

func (m *ConfigSource) Reset()         { *m = ConfigSource{} }
func (m *ConfigSource) String() string { return proto.CompactTextString(m) }
func (*ConfigSource) ProtoMessage()    {}
func (*ConfigSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_a3116f5d2bd13e64, []int{4}
}

func (m *ConfigSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigSource.Unmarshal(m, b)
}
func (m *ConfigSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigSource.Marshal(b, m, deterministic)
}
func (m *ConfigSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigSource.Merge(m, src)
}
func (m *ConfigSource) XXX_Size() int {
	return xxx_messageInfo_ConfigSource.Size(m)
}
func (m *ConfigSource) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigSource.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigSource proto.InternalMessageInfo

type isConfigSource_ConfigSourceSpecifier interface {
	isConfigSource_ConfigSourceSpecifier()
}

type ConfigSource_Path struct {
	Path string `protobuf:"bytes,1,opt,name=path,proto3,oneof"`
}

type ConfigSource_ApiConfigSource struct {
	ApiConfigSource *ApiConfigSource `protobuf:"bytes,2,opt,name=api_config_source,json=apiConfigSource,proto3,oneof"`
}

type ConfigSource_Ads struct {
	Ads *AggregatedConfigSource `protobuf:"bytes,3,opt,name=ads,proto3,oneof"`
}

type ConfigSource_Self struct {
	Self *SelfConfigSource `protobuf:"bytes,5,opt,name=self,proto3,oneof"`
}

func (*ConfigSource_Path) isConfigSource_ConfigSourceSpecifier() {}

func (*ConfigSource_ApiConfigSource) isConfigSource_ConfigSourceSpecifier() {}

func (*ConfigSource_Ads) isConfigSource_ConfigSourceSpecifier() {}

func (*ConfigSource_Self) isConfigSource_ConfigSourceSpecifier() {}

func (m *ConfigSource) GetConfigSourceSpecifier() isConfigSource_ConfigSourceSpecifier {
	if m != nil {
		return m.ConfigSourceSpecifier
	}
	return nil
}

func (m *ConfigSource) GetPath() string {
	if x, ok := m.GetConfigSourceSpecifier().(*ConfigSource_Path); ok {
		return x.Path
	}
	return ""
}

func (m *ConfigSource) GetApiConfigSource() *ApiConfigSource {
	if x, ok := m.GetConfigSourceSpecifier().(*ConfigSource_ApiConfigSource); ok {
		return x.ApiConfigSource
	}
	return nil
}

func (m *ConfigSource) GetAds() *AggregatedConfigSource {
	if x, ok := m.GetConfigSourceSpecifier().(*ConfigSource_Ads); ok {
		return x.Ads
	}
	return nil
}

func (m *ConfigSource) GetSelf() *SelfConfigSource {
	if x, ok := m.GetConfigSourceSpecifier().(*ConfigSource_Self); ok {
		return x.Self
	}
	return nil
}

func (m *ConfigSource) GetInitialFetchTimeout() *duration.Duration {
	if m != nil {
		return m.InitialFetchTimeout
	}
	return nil
}

func (m *ConfigSource) GetResourceApiVersion() ApiVersion {
	if m != nil {
		return m.ResourceApiVersion
	}
	return ApiVersion_AUTO
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*ConfigSource) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*ConfigSource_Path)(nil),
		(*ConfigSource_ApiConfigSource)(nil),
		(*ConfigSource_Ads)(nil),
		(*ConfigSource_Self)(nil),
	}
}

func init() {
	proto.RegisterEnum("envoy.config.core.v3.ApiVersion", ApiVersion_name, ApiVersion_value)
	proto.RegisterEnum("envoy.config.core.v3.ApiConfigSource_ApiType", ApiConfigSource_ApiType_name, ApiConfigSource_ApiType_value)
	proto.RegisterType((*ApiConfigSource)(nil), "envoy.config.core.v3.ApiConfigSource")
	proto.RegisterType((*AggregatedConfigSource)(nil), "envoy.config.core.v3.AggregatedConfigSource")
	proto.RegisterType((*SelfConfigSource)(nil), "envoy.config.core.v3.SelfConfigSource")
	proto.RegisterType((*RateLimitSettings)(nil), "envoy.config.core.v3.RateLimitSettings")
	proto.RegisterType((*ConfigSource)(nil), "envoy.config.core.v3.ConfigSource")
}

func init() {
	proto.RegisterFile("envoy/config/core/v3/config_source.proto", fileDescriptor_a3116f5d2bd13e64)
}

var fileDescriptor_a3116f5d2bd13e64 = []byte{
	// 943 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x4b, 0x6f, 0x23, 0x45,
	0x10, 0xf6, 0x8c, 0xbd, 0x89, 0xd3, 0x79, 0x39, 0x9d, 0x2c, 0xeb, 0x8d, 0xc0, 0x72, 0x1c, 0xb2,
	0x6b, 0x22, 0x18, 0x4b, 0xf6, 0x2d, 0xb0, 0x08, 0x3b, 0xce, 0x4b, 0xca, 0x26, 0xd6, 0xd8, 0x09,
	0x42, 0x42, 0xb4, 0x3a, 0x9e, 0xf2, 0xa4, 0xc5, 0x78, 0x7a, 0xe8, 0xee, 0x31, 0xf1, 0x0d, 0x71,
	0x42, 0xe2, 0x1f, 0xf0, 0x0b, 0x10, 0x67, 0x4e, 0xdc, 0x11, 0x5c, 0xf9, 0x21, 0xdc, 0x90, 0x40,
	0x7b, 0x42, 0x3d, 0x33, 0xde, 0xf8, 0x15, 0x2d, 0xf8, 0xe2, 0xa9, 0xae, 0xef, 0xfb, 0xba, 0xaa,
	0xeb, 0x81, 0xca, 0xe0, 0x0f, 0xf8, 0xb0, 0xd2, 0xe5, 0x7e, 0x8f, 0xb9, 0x95, 0x2e, 0x17, 0x50,
	0x19, 0xd4, 0x12, 0x93, 0x48, 0x1e, 0x8a, 0x2e, 0x58, 0x81, 0xe0, 0x8a, 0xe3, 0xad, 0x08, 0x69,
	0xc5, 0x2e, 0x4b, 0x23, 0xad, 0x41, 0x6d, 0xfb, 0xf9, 0x5c, 0xbe, 0x2b, 0x82, 0x2e, 0x91, 0x20,
	0x06, 0x6c, 0x44, 0xdf, 0x2e, 0xb8, 0x9c, 0xbb, 0x1e, 0x54, 0x22, 0xeb, 0x26, 0xec, 0x55, 0x9c,
	0x50, 0x50, 0xc5, 0xb8, 0xff, 0x90, 0xff, 0x6b, 0x41, 0x83, 0x00, 0x84, 0x4c, 0xfc, 0xbb, 0xf1,
	0x45, 0xd4, 0xf7, 0xb9, 0x8a, 0x78, 0xb2, 0xe2, 0x40, 0x20, 0xa0, 0x3b, 0x2e, 0xf2, 0x4e, 0xe8,
	0x04, 0x74, 0x02, 0x23, 0x15, 0x55, 0xe1, 0x48, 0x63, 0x67, 0xc6, 0x3d, 0x00, 0x21, 0x19, 0xf7,
	0x99, 0xef, 0x26, 0x90, 0x27, 0x03, 0xea, 0x31, 0x87, 0x2a, 0xa8, 0x8c, 0x3e, 0x62, 0x47, 0xe9,
	0xaf, 0x47, 0x68, 0xbd, 0x1e, 0xb0, 0xc3, 0x28, 0xd1, 0x76, 0xf4, 0x30, 0xd8, 0x46, 0x59, 0x1a,
	0x30, 0xa2, 0x86, 0x01, 0xe4, 0x8d, 0xa2, 0x51, 0x5e, 0xab, 0x7e, 0x60, 0xcd, 0x7b, 0x25, 0x6b,
	0x8a, 0xa8, 0xed, 0xce, 0x30, 0x80, 0x46, 0xf6, 0x55, 0xe3, 0xd1, 0xb7, 0x86, 0x99, 0x33, 0xec,
	0x45, 0x1a, 0x1f, 0xe1, 0x2f, 0xd0, 0x63, 0x25, 0xa8, 0x2f, 0x03, 0x2e, 0x14, 0xd1, 0xea, 0x49,
	0x88, 0xf9, 0x6c, 0x74, 0x41, 0xf1, 0xc1, 0x0b, 0xae, 0x63, 0xdc, 0x98, 0xe6, 0xe6, 0x6b, 0xa1,
	0x7b, 0x37, 0xde, 0x45, 0xab, 0x5d, 0x2f, 0x94, 0x0a, 0x04, 0xf1, 0x69, 0x1f, 0x64, 0xde, 0x2c,
	0xa6, 0xcb, 0x4b, 0xf6, 0x4a, 0x72, 0x78, 0xa1, 0xcf, 0xf0, 0x31, 0x5a, 0x1d, 0x2f, 0xa1, 0xcc,
	0x67, 0x8a, 0xe9, 0xf2, 0x72, 0x75, 0x67, 0xfe, 0xe5, 0x27, 0x22, 0xe8, 0xb6, 0x63, 0xa4, 0xbd,
	0xe2, 0xde, 0x1b, 0x12, 0x7f, 0x8c, 0x56, 0x05, 0xf4, 0x04, 0xc8, 0x5b, 0xe2, 0x80, 0x47, 0x87,
	0xf9, 0x74, 0xd1, 0x28, 0x2f, 0x57, 0x9f, 0x5a, 0x71, 0xb1, 0xad, 0x51, 0xb1, 0xad, 0x66, 0xd2,
	0x0c, 0xf6, 0x4a, 0x82, 0x6f, 0x6a, 0x38, 0x3e, 0x47, 0xeb, 0x02, 0xbe, 0x0a, 0x41, 0x2a, 0xa2,
	0x58, 0x1f, 0x78, 0xa8, 0xf2, 0x8f, 0xde, 0xa0, 0x10, 0xe5, 0xff, 0x93, 0x61, 0xee, 0xa7, 0xec,
	0xb5, 0x84, 0xdb, 0x89, 0xa9, 0xf8, 0x53, 0xb4, 0x29, 0xa8, 0x02, 0xe2, 0xb1, 0x3e, 0x53, 0x44,
	0x82, 0x52, 0xcc, 0x77, 0x65, 0x7e, 0x21, 0x52, 0x7c, 0x3e, 0x3f, 0x37, 0x9b, 0x2a, 0x38, 0xd7,
	0xf8, 0x76, 0x02, 0xb7, 0x37, 0xc4, 0xf4, 0x11, 0xae, 0xa3, 0x82, 0x04, 0x45, 0x7c, 0xee, 0x00,
	0xe1, 0x3e, 0xe9, 0x31, 0x21, 0x15, 0xe9, 0x83, 0x94, 0xd4, 0xd5, 0x07, 0xde, 0x30, 0xbf, 0x58,
	0x34, 0xca, 0x59, 0xfb, 0xa9, 0x04, 0x75, 0xc1, 0x1d, 0xb8, 0xf4, 0x8f, 0x35, 0xe4, 0x65, 0x8c,
	0xb8, 0xf4, 0xbd, 0x61, 0xc9, 0x43, 0x8b, 0x49, 0x53, 0xe0, 0x17, 0xa8, 0x72, 0xcb, 0x1c, 0x07,
	0x7c, 0x12, 0x45, 0x44, 0x46, 0x6d, 0x0e, 0x0e, 0xb9, 0xba, 0x68, 0x5f, 0xb5, 0x5a, 0x97, 0x76,
	0xe7, 0xa8, 0x49, 0xec, 0xa3, 0x76, 0x87, 0x9c, 0x1f, 0x9d, 0xd4, 0x0f, 0x3f, 0xcb, 0xa5, 0xb6,
	0xb3, 0x3f, 0xfe, 0xfd, 0xf3, 0xf7, 0xa6, 0x91, 0x35, 0x70, 0x16, 0x65, 0xb4, 0x2b, 0x17, 0x7d,
	0x9d, 0xd8, 0xad, 0xc3, 0x9c, 0x89, 0xd7, 0x10, 0x6a, 0x1e, 0x9d, 0x77, 0xea, 0x24, 0xb2, 0xd3,
	0x07, 0xe5, 0x1f, 0x7e, 0xfd, 0xae, 0xb0, 0x8b, 0x76, 0xe2, 0x94, 0x69, 0xc0, 0xac, 0x41, 0x35,
	0x4e, 0x79, 0xaa, 0x53, 0x4b, 0x67, 0xe8, 0xad, 0xba, 0xeb, 0x0a, 0x70, 0xf5, 0xfd, 0xe3, 0x9e,
	0x83, 0x8a, 0xd6, 0xd8, 0x4f, 0x16, 0xc8, 0xa4, 0xc6, 0x5c, 0x42, 0xe9, 0x05, 0xca, 0xb5, 0xc1,
	0xeb, 0x4d, 0x88, 0xbc, 0xa7, 0x45, 0xde, 0x45, 0xa5, 0x59, 0x91, 0x69, 0x68, 0xe9, 0x37, 0x03,
	0x6d, 0xcc, 0x54, 0x03, 0x7f, 0x88, 0x50, 0x9f, 0xde, 0x11, 0xc5, 0xbf, 0x04, 0x5f, 0x46, 0x43,
	0xb8, 0x5c, 0x7d, 0x7b, 0xa6, 0x39, 0xae, 0xce, 0x7c, 0x55, 0xab, 0x5e, 0x53, 0x2f, 0x04, 0x7b,
	0xa9, 0x4f, 0xef, 0x3a, 0x11, 0x1c, 0x9f, 0xa1, 0xa5, 0x1e, 0xf3, 0x3c, 0xa2, 0x2b, 0x9a, 0x37,
	0x1f, 0xe0, 0x36, 0x79, 0x78, 0xe3, 0x41, 0xc4, 0x6d, 0xac, 0xbd, 0x6a, 0x2c, 0xe3, 0xa5, 0x9d,
	0x54, 0xf2, 0xb3, 0xb3, 0x9a, 0xae, 0x83, 0x3a, 0xd8, 0xd7, 0x89, 0xec, 0xa1, 0xdd, 0xd9, 0x44,
	0x66, 0x62, 0x2e, 0xfd, 0x99, 0x46, 0x2b, 0x13, 0x7b, 0x64, 0x0b, 0x65, 0x02, 0xaa, 0x6e, 0xa3,
	0xf0, 0x97, 0x4e, 0x53, 0x76, 0x64, 0xe1, 0x36, 0xda, 0xd0, 0xf3, 0x3f, 0xb1, 0x8b, 0x93, 0x28,
	0xf7, 0xfe, 0xd3, 0x9a, 0x39, 0x4d, 0xd9, 0xeb, 0x74, 0x6a, 0x65, 0x7d, 0x82, 0xd2, 0xd4, 0x91,
	0xc9, 0x1c, 0xbe, 0xff, 0x80, 0xcc, 0xdc, 0xfa, 0x9d, 0xa6, 0x6c, 0x4d, 0xc5, 0x1f, 0xa1, 0x8c,
	0x04, 0xaf, 0x97, 0x0c, 0xe2, 0xb3, 0xf9, 0x12, 0xd3, 0xd5, 0xd3, 0x49, 0x69, 0x16, 0x7e, 0x89,
	0x1e, 0x33, 0x9f, 0x29, 0x46, 0x3d, 0xd2, 0x03, 0xd5, 0xbd, 0x7d, 0x3d, 0xd7, 0x99, 0x37, 0x6d,
	0x86, 0xcd, 0x84, 0x77, 0xac, 0x69, 0xa3, 0x91, 0xfe, 0x1c, 0x6d, 0x09, 0x88, 0x9f, 0x66, 0x62,
	0x59, 0x2e, 0xfc, 0xef, 0x65, 0x89, 0x47, 0x3a, 0xf7, 0xde, 0x83, 0x3d, 0x5d, 0xd4, 0x22, 0x2a,
	0xcc, 0x16, 0x75, 0x3c, 0xb7, 0x46, 0x01, 0x3d, 0x99, 0x28, 0x12, 0x91, 0x01, 0x74, 0x59, 0x8f,
	0x81, 0xc0, 0xe9, 0x7f, 0x1a, 0xc6, 0xfe, 0x33, 0x84, 0xc6, 0x16, 0x70, 0x16, 0x65, 0xea, 0x57,
	0x9d, 0xcb, 0x5c, 0x0a, 0x2f, 0x20, 0xf3, 0xba, 0x9a, 0x33, 0xa2, 0xff, 0x5a, 0xce, 0x6c, 0xd4,
	0x7f, 0xf9, 0xe6, 0xf7, 0x3f, 0x16, 0xcc, 0x9c, 0x89, 0x4a, 0x8c, 0xc7, 0xa1, 0x07, 0x82, 0xdf,
	0x0d, 0xe7, 0x66, 0xd1, 0xd8, 0x18, 0x8f, 0xa1, 0xa5, 0x9f, 0xab, 0x65, 0xdc, 0x2c, 0x44, 0xef,
	0x56, 0xfb, 0x37, 0x00, 0x00, 0xff, 0xff, 0x11, 0x95, 0x60, 0xb9, 0xd7, 0x07, 0x00, 0x00,
}