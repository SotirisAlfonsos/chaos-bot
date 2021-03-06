// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: manager.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StatusResponse_Status int32

const (
	StatusResponse_SUCCESS StatusResponse_Status = 0
	StatusResponse_FAIL    StatusResponse_Status = 1
)

// Enum value maps for StatusResponse_Status.
var (
	StatusResponse_Status_name = map[int32]string{
		0: "SUCCESS",
		1: "FAIL",
	}
	StatusResponse_Status_value = map[string]int32{
		"SUCCESS": 0,
		"FAIL":    1,
	}
)

func (x StatusResponse_Status) Enum() *StatusResponse_Status {
	p := new(StatusResponse_Status)
	*p = x
	return p
}

func (x StatusResponse_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StatusResponse_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_manager_proto_enumTypes[0].Descriptor()
}

func (StatusResponse_Status) Type() protoreflect.EnumType {
	return &file_manager_proto_enumTypes[0]
}

func (x StatusResponse_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StatusResponse_Status.Descriptor instead.
func (StatusResponse_Status) EnumDescriptor() ([]byte, []int) {
	return file_manager_proto_rawDescGZIP(), []int{5, 0}
}

type ServiceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *ServiceRequest) Reset() {
	*x = ServiceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manager_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceRequest) ProtoMessage() {}

func (x *ServiceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_manager_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceRequest.ProtoReflect.Descriptor instead.
func (*ServiceRequest) Descriptor() ([]byte, []int) {
	return file_manager_proto_rawDescGZIP(), []int{0}
}

func (x *ServiceRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type DockerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *DockerRequest) Reset() {
	*x = DockerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manager_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DockerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DockerRequest) ProtoMessage() {}

func (x *DockerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_manager_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DockerRequest.ProtoReflect.Descriptor instead.
func (*DockerRequest) Descriptor() ([]byte, []int) {
	return file_manager_proto_rawDescGZIP(), []int{1}
}

func (x *DockerRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type CPURequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Percentage int32 `protobuf:"varint,1,opt,name=percentage,proto3" json:"percentage,omitempty"`
}

func (x *CPURequest) Reset() {
	*x = CPURequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manager_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CPURequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CPURequest) ProtoMessage() {}

func (x *CPURequest) ProtoReflect() protoreflect.Message {
	mi := &file_manager_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CPURequest.ProtoReflect.Descriptor instead.
func (*CPURequest) Descriptor() ([]byte, []int) {
	return file_manager_proto_rawDescGZIP(), []int{2}
}

func (x *CPURequest) GetPercentage() int32 {
	if x != nil {
		return x.Percentage
	}
	return 0
}

type ServerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ServerRequest) Reset() {
	*x = ServerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manager_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerRequest) ProtoMessage() {}

func (x *ServerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_manager_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerRequest.ProtoReflect.Descriptor instead.
func (*ServerRequest) Descriptor() ([]byte, []int) {
	return file_manager_proto_rawDescGZIP(), []int{3}
}

type NetworkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Device        string  `protobuf:"bytes,1,opt,name=device,proto3" json:"device,omitempty"`
	Latency       uint32  `protobuf:"varint,2,opt,name=latency,proto3" json:"latency,omitempty"`
	DelayCorr     float32 `protobuf:"fixed32,3,opt,name=delayCorr,proto3" json:"delayCorr,omitempty"`
	Limit         uint32  `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
	Loss          float32 `protobuf:"fixed32,5,opt,name=loss,proto3" json:"loss,omitempty"`
	LossCorr      float32 `protobuf:"fixed32,6,opt,name=lossCorr,proto3" json:"lossCorr,omitempty"`
	Gap           uint32  `protobuf:"varint,7,opt,name=gap,proto3" json:"gap,omitempty"`
	Duplicate     float32 `protobuf:"fixed32,8,opt,name=duplicate,proto3" json:"duplicate,omitempty"`
	DuplicateCorr float32 `protobuf:"fixed32,9,opt,name=duplicateCorr,proto3" json:"duplicateCorr,omitempty"`
	Jitter        uint32  `protobuf:"varint,10,opt,name=jitter,proto3" json:"jitter,omitempty"`
	ReorderProb   float32 `protobuf:"fixed32,11,opt,name=reorderProb,proto3" json:"reorderProb,omitempty"`
	ReorderCorr   float32 `protobuf:"fixed32,12,opt,name=reorderCorr,proto3" json:"reorderCorr,omitempty"`
	CorruptProb   float32 `protobuf:"fixed32,13,opt,name=corruptProb,proto3" json:"corruptProb,omitempty"`
	CorruptCorr   float32 `protobuf:"fixed32,14,opt,name=corruptCorr,proto3" json:"corruptCorr,omitempty"`
}

func (x *NetworkRequest) Reset() {
	*x = NetworkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manager_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkRequest) ProtoMessage() {}

func (x *NetworkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_manager_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkRequest.ProtoReflect.Descriptor instead.
func (*NetworkRequest) Descriptor() ([]byte, []int) {
	return file_manager_proto_rawDescGZIP(), []int{4}
}

func (x *NetworkRequest) GetDevice() string {
	if x != nil {
		return x.Device
	}
	return ""
}

func (x *NetworkRequest) GetLatency() uint32 {
	if x != nil {
		return x.Latency
	}
	return 0
}

func (x *NetworkRequest) GetDelayCorr() float32 {
	if x != nil {
		return x.DelayCorr
	}
	return 0
}

func (x *NetworkRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *NetworkRequest) GetLoss() float32 {
	if x != nil {
		return x.Loss
	}
	return 0
}

func (x *NetworkRequest) GetLossCorr() float32 {
	if x != nil {
		return x.LossCorr
	}
	return 0
}

func (x *NetworkRequest) GetGap() uint32 {
	if x != nil {
		return x.Gap
	}
	return 0
}

func (x *NetworkRequest) GetDuplicate() float32 {
	if x != nil {
		return x.Duplicate
	}
	return 0
}

func (x *NetworkRequest) GetDuplicateCorr() float32 {
	if x != nil {
		return x.DuplicateCorr
	}
	return 0
}

func (x *NetworkRequest) GetJitter() uint32 {
	if x != nil {
		return x.Jitter
	}
	return 0
}

func (x *NetworkRequest) GetReorderProb() float32 {
	if x != nil {
		return x.ReorderProb
	}
	return 0
}

func (x *NetworkRequest) GetReorderCorr() float32 {
	if x != nil {
		return x.ReorderCorr
	}
	return 0
}

func (x *NetworkRequest) GetCorruptProb() float32 {
	if x != nil {
		return x.CorruptProb
	}
	return 0
}

func (x *NetworkRequest) GetCorruptCorr() float32 {
	if x != nil {
		return x.CorruptCorr
	}
	return 0
}

type StatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  StatusResponse_Status `protobuf:"varint,1,opt,name=status,proto3,enum=v1.StatusResponse_Status" json:"status,omitempty"`
	Message string                `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *StatusResponse) Reset() {
	*x = StatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_manager_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatusResponse) ProtoMessage() {}

func (x *StatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_manager_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatusResponse.ProtoReflect.Descriptor instead.
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return file_manager_proto_rawDescGZIP(), []int{5}
}

func (x *StatusResponse) GetStatus() StatusResponse_Status {
	if x != nil {
		return x.Status
	}
	return StatusResponse_SUCCESS
}

func (x *StatusResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_manager_proto protoreflect.FileDescriptor

var file_manager_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x76, 0x31, 0x22, 0x24, 0x0a, 0x0e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x23, 0x0a, 0x0d, 0x44, 0x6f, 0x63,
	0x6b, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2c,
	0x0a, 0x0a, 0x43, 0x50, 0x55, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a,
	0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0a, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x61, 0x67, 0x65, 0x22, 0x0f, 0x0a, 0x0d,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x9c, 0x03,
	0x0a, 0x0e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6c, 0x61, 0x74, 0x65,
	0x6e, 0x63, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x6c, 0x61, 0x74, 0x65, 0x6e,
	0x63, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x43, 0x6f, 0x72, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x43, 0x6f, 0x72, 0x72,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x6f, 0x73, 0x73, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x04, 0x6c, 0x6f, 0x73, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f,
	0x73, 0x73, 0x43, 0x6f, 0x72, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02, 0x52, 0x08, 0x6c, 0x6f,
	0x73, 0x73, 0x43, 0x6f, 0x72, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x67, 0x61, 0x70, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x03, 0x67, 0x61, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x75, 0x70, 0x6c,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x02, 0x52, 0x09, 0x64, 0x75, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x64, 0x75, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x74, 0x65, 0x43, 0x6f, 0x72, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0d, 0x64,
	0x75, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x72, 0x72, 0x12, 0x16, 0x0a, 0x06,
	0x6a, 0x69, 0x74, 0x74, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6a, 0x69,
	0x74, 0x74, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x50,
	0x72, 0x6f, 0x62, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0b, 0x72, 0x65, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x50, 0x72, 0x6f, 0x62, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x43, 0x6f, 0x72, 0x72, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0b, 0x72, 0x65, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x72, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x72, 0x72,
	0x75, 0x70, 0x74, 0x50, 0x72, 0x6f, 0x62, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0b, 0x63,
	0x6f, 0x72, 0x72, 0x75, 0x70, 0x74, 0x50, 0x72, 0x6f, 0x62, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f,
	0x72, 0x72, 0x75, 0x70, 0x74, 0x43, 0x6f, 0x72, 0x72, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x02, 0x52,
	0x0b, 0x63, 0x6f, 0x72, 0x72, 0x75, 0x70, 0x74, 0x43, 0x6f, 0x72, 0x72, 0x22, 0x7e, 0x0a, 0x0e,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x1f, 0x0a, 0x06, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53,
	0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x41, 0x49, 0x4c, 0x10, 0x01, 0x32, 0x70, 0x0a, 0x07,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x4b, 0x69, 0x6c, 0x6c, 0x12,
	0x12, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x33, 0x0a, 0x07, 0x52, 0x65, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x6d,
	0x0a, 0x06, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x04, 0x4b, 0x69, 0x6c, 0x6c,
	0x12, 0x11, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x32, 0x0a, 0x07, 0x52, 0x65, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x12, 0x11, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x6f, 0x63, 0x6b, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x65, 0x0a,
	0x03, 0x43, 0x50, 0x55, 0x12, 0x2d, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x0e, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x50, 0x55, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x2f, 0x0a, 0x07, 0x52, 0x65, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x12, 0x0e,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x50, 0x55, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x32, 0x39, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x2f,
	0x0a, 0x04, 0x4b, 0x69, 0x6c, 0x6c, 0x12, 0x11, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32,
	0x71, 0x0a, 0x07, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12, 0x31, 0x0a, 0x05, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x12, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x33, 0x0a,
	0x07, 0x52, 0x65, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x12, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x09, 0x5a, 0x07, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_manager_proto_rawDescOnce sync.Once
	file_manager_proto_rawDescData = file_manager_proto_rawDesc
)

func file_manager_proto_rawDescGZIP() []byte {
	file_manager_proto_rawDescOnce.Do(func() {
		file_manager_proto_rawDescData = protoimpl.X.CompressGZIP(file_manager_proto_rawDescData)
	})
	return file_manager_proto_rawDescData
}

var file_manager_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_manager_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_manager_proto_goTypes = []interface{}{
	(StatusResponse_Status)(0), // 0: v1.StatusResponse.Status
	(*ServiceRequest)(nil),     // 1: v1.ServiceRequest
	(*DockerRequest)(nil),      // 2: v1.DockerRequest
	(*CPURequest)(nil),         // 3: v1.CPURequest
	(*ServerRequest)(nil),      // 4: v1.ServerRequest
	(*NetworkRequest)(nil),     // 5: v1.NetworkRequest
	(*StatusResponse)(nil),     // 6: v1.StatusResponse
}
var file_manager_proto_depIdxs = []int32{
	0,  // 0: v1.StatusResponse.status:type_name -> v1.StatusResponse.Status
	1,  // 1: v1.Service.Kill:input_type -> v1.ServiceRequest
	1,  // 2: v1.Service.Recover:input_type -> v1.ServiceRequest
	2,  // 3: v1.Docker.Kill:input_type -> v1.DockerRequest
	2,  // 4: v1.Docker.Recover:input_type -> v1.DockerRequest
	3,  // 5: v1.CPU.Start:input_type -> v1.CPURequest
	3,  // 6: v1.CPU.Recover:input_type -> v1.CPURequest
	4,  // 7: v1.Server.Kill:input_type -> v1.ServerRequest
	5,  // 8: v1.Network.Start:input_type -> v1.NetworkRequest
	5,  // 9: v1.Network.Recover:input_type -> v1.NetworkRequest
	6,  // 10: v1.Service.Kill:output_type -> v1.StatusResponse
	6,  // 11: v1.Service.Recover:output_type -> v1.StatusResponse
	6,  // 12: v1.Docker.Kill:output_type -> v1.StatusResponse
	6,  // 13: v1.Docker.Recover:output_type -> v1.StatusResponse
	6,  // 14: v1.CPU.Start:output_type -> v1.StatusResponse
	6,  // 15: v1.CPU.Recover:output_type -> v1.StatusResponse
	6,  // 16: v1.Server.Kill:output_type -> v1.StatusResponse
	6,  // 17: v1.Network.Start:output_type -> v1.StatusResponse
	6,  // 18: v1.Network.Recover:output_type -> v1.StatusResponse
	10, // [10:19] is the sub-list for method output_type
	1,  // [1:10] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_manager_proto_init() }
func file_manager_proto_init() {
	if File_manager_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_manager_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_manager_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DockerRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_manager_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CPURequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_manager_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_manager_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_manager_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatusResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_manager_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   5,
		},
		GoTypes:           file_manager_proto_goTypes,
		DependencyIndexes: file_manager_proto_depIdxs,
		EnumInfos:         file_manager_proto_enumTypes,
		MessageInfos:      file_manager_proto_msgTypes,
	}.Build()
	File_manager_proto = out.File
	file_manager_proto_rawDesc = nil
	file_manager_proto_goTypes = nil
	file_manager_proto_depIdxs = nil
}
