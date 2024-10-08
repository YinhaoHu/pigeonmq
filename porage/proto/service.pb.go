// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.12.4
// source: service.proto

package service

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// CreateLedgerRequest is the request message for the CreateLedger RPC.
type CreateLedgerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LedgerId uint64 `protobuf:"varint,1,opt,name=ledger_id,json=ledgerId,proto3" json:"ledger_id,omitempty"`
}

func (x *CreateLedgerRequest) Reset() {
	*x = CreateLedgerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateLedgerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateLedgerRequest) ProtoMessage() {}

func (x *CreateLedgerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateLedgerRequest.ProtoReflect.Descriptor instead.
func (*CreateLedgerRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreateLedgerRequest) GetLedgerId() uint64 {
	if x != nil {
		return x.LedgerId
	}
	return 0
}

// AppendEntryOnLedgerRequest is the request message for the PutEntryOnLedger RPC.
type AppendEntryOnLedgerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LedgerId uint64 `protobuf:"varint,1,opt,name=ledger_id,json=ledgerId,proto3" json:"ledger_id,omitempty"`
	Payload  []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *AppendEntryOnLedgerRequest) Reset() {
	*x = AppendEntryOnLedgerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntryOnLedgerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntryOnLedgerRequest) ProtoMessage() {}

func (x *AppendEntryOnLedgerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntryOnLedgerRequest.ProtoReflect.Descriptor instead.
func (*AppendEntryOnLedgerRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1}
}

func (x *AppendEntryOnLedgerRequest) GetLedgerId() uint64 {
	if x != nil {
		return x.LedgerId
	}
	return 0
}

func (x *AppendEntryOnLedgerRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

// AppendEntryOnLedgerResponse is the response message for the PutEntryOnLedger RPC.
type AppendEntryOnLedgerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EntryId int64 `protobuf:"varint,1,opt,name=entry_id,json=entryId,proto3" json:"entry_id,omitempty"`
}

func (x *AppendEntryOnLedgerResponse) Reset() {
	*x = AppendEntryOnLedgerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppendEntryOnLedgerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntryOnLedgerResponse) ProtoMessage() {}

func (x *AppendEntryOnLedgerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntryOnLedgerResponse.ProtoReflect.Descriptor instead.
func (*AppendEntryOnLedgerResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{2}
}

func (x *AppendEntryOnLedgerResponse) GetEntryId() int64 {
	if x != nil {
		return x.EntryId
	}
	return 0
}

// GetEntryFromLedgerRequest is the request message for the GetEntryFromLedger RPC.
type GetEntryFromLedgerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LedgerId uint64 `protobuf:"varint,1,opt,name=ledger_id,json=ledgerId,proto3" json:"ledger_id,omitempty"`
	EntryId  int64  `protobuf:"varint,2,opt,name=entry_id,json=entryId,proto3" json:"entry_id,omitempty"`
}

func (x *GetEntryFromLedgerRequest) Reset() {
	*x = GetEntryFromLedgerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEntryFromLedgerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEntryFromLedgerRequest) ProtoMessage() {}

func (x *GetEntryFromLedgerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEntryFromLedgerRequest.ProtoReflect.Descriptor instead.
func (*GetEntryFromLedgerRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetEntryFromLedgerRequest) GetLedgerId() uint64 {
	if x != nil {
		return x.LedgerId
	}
	return 0
}

func (x *GetEntryFromLedgerRequest) GetEntryId() int64 {
	if x != nil {
		return x.EntryId
	}
	return 0
}

// GetEntryFromLedgerResponse is the response message for the GetEntryFromLedger RPC.
type GetEntryFromLedgerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *GetEntryFromLedgerResponse) Reset() {
	*x = GetEntryFromLedgerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetEntryFromLedgerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEntryFromLedgerResponse) ProtoMessage() {}

func (x *GetEntryFromLedgerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetEntryFromLedgerResponse.ProtoReflect.Descriptor instead.
func (*GetEntryFromLedgerResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{4}
}

func (x *GetEntryFromLedgerResponse) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

// DeleteLedgerRequest is the request message for the DeleteLedger RPC.
type DeleteLedgerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LedgerId uint64 `protobuf:"varint,1,opt,name=ledger_id,json=ledgerId,proto3" json:"ledger_id,omitempty"`
}

func (x *DeleteLedgerRequest) Reset() {
	*x = DeleteLedgerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteLedgerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteLedgerRequest) ProtoMessage() {}

func (x *DeleteLedgerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteLedgerRequest.ProtoReflect.Descriptor instead.
func (*DeleteLedgerRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteLedgerRequest) GetLedgerId() uint64 {
	if x != nil {
		return x.LedgerId
	}
	return 0
}

// ListLedgersResponse is the response message for the ListLedgers RPC.
type ListLedgersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LedgerIds []uint64 `protobuf:"varint,1,rep,packed,name=ledger_ids,json=ledgerIds,proto3" json:"ledger_ids,omitempty"`
}

func (x *ListLedgersResponse) Reset() {
	*x = ListLedgersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListLedgersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListLedgersResponse) ProtoMessage() {}

func (x *ListLedgersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListLedgersResponse.ProtoReflect.Descriptor instead.
func (*ListLedgersResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{6}
}

func (x *ListLedgersResponse) GetLedgerIds() []uint64 {
	if x != nil {
		return x.LedgerIds
	}
	return nil
}

// ListWorkersResponse is the response message for the ListWorkers RPC. A map from worker name(string) to
// description(WorkerDescription) is returned.
type ListWorkersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Workers map[string]*WorkerDescription `protobuf:"bytes,1,rep,name=workers,proto3" json:"workers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ListWorkersResponse) Reset() {
	*x = ListWorkersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListWorkersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWorkersResponse) ProtoMessage() {}

func (x *ListWorkersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWorkersResponse.ProtoReflect.Descriptor instead.
func (*ListWorkersResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{7}
}

func (x *ListWorkersResponse) GetWorkers() map[string]*WorkerDescription {
	if x != nil {
		return x.Workers
	}
	return nil
}

// WorkerDescription is the description of a worker.
type WorkerDescription struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Description string `protobuf:"bytes,1,opt,name=description,proto3" json:"description,omitempty"`
}

func (x *WorkerDescription) Reset() {
	*x = WorkerDescription{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerDescription) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerDescription) ProtoMessage() {}

func (x *WorkerDescription) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerDescription.ProtoReflect.Descriptor instead.
func (*WorkerDescription) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{8}
}

func (x *WorkerDescription) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// LedgerLengthRequest is the request message for the LedgerLength RPC.
type LedgerLengthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LedgerId uint64 `protobuf:"varint,1,opt,name=ledger_id,json=ledgerId,proto3" json:"ledger_id,omitempty"`
}

func (x *LedgerLengthRequest) Reset() {
	*x = LedgerLengthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LedgerLengthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LedgerLengthRequest) ProtoMessage() {}

func (x *LedgerLengthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LedgerLengthRequest.ProtoReflect.Descriptor instead.
func (*LedgerLengthRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{9}
}

func (x *LedgerLengthRequest) GetLedgerId() uint64 {
	if x != nil {
		return x.LedgerId
	}
	return 0
}

// LedgerLengthResponse is the response message for the LedgerLength RPC.
type LedgerLengthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Length int64 `protobuf:"varint,1,opt,name=length,proto3" json:"length,omitempty"`
}

func (x *LedgerLengthResponse) Reset() {
	*x = LedgerLengthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LedgerLengthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LedgerLengthResponse) ProtoMessage() {}

func (x *LedgerLengthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LedgerLengthResponse.ProtoReflect.Descriptor instead.
func (*LedgerLengthResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{10}
}

func (x *LedgerLengthResponse) GetLength() int64 {
	if x != nil {
		return x.Length
	}
	return 0
}

var File_service_proto protoreflect.FileDescriptor

var file_service_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0d, 0x70, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x32, 0x0a, 0x13, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x49, 0x64, 0x22,
	0x53, 0x0a, 0x1a, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x4f, 0x6e,
	0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x08, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x22, 0x38, 0x0a, 0x1b, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x4f, 0x6e, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x49, 0x64, 0x22, 0x53,
	0x0a, 0x19, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x46, 0x72, 0x6f, 0x6d, 0x4c, 0x65,
	0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6c,
	0x65, 0x64, 0x67, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08,
	0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x74, 0x72,
	0x79, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72,
	0x79, 0x49, 0x64, 0x22, 0x36, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x46,
	0x72, 0x6f, 0x6d, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x32, 0x0a, 0x13, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x49, 0x64, 0x22,
	0x34, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x04, 0x52, 0x09, 0x6c, 0x65, 0x64, 0x67,
	0x65, 0x72, 0x49, 0x64, 0x73, 0x22, 0xbe, 0x01, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x49, 0x0a,
	0x07, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f,
	0x2e, 0x70, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x07, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x73, 0x1a, 0x5c, 0x0a, 0x0c, 0x57, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x36, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x35, 0x0a, 0x11, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x32, 0x0a,
	0x13, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x2e, 0x0a, 0x14, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x4c, 0x65, 0x6e, 0x67, 0x74,
	0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x6e,
	0x67, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74,
	0x68, 0x32, 0xfd, 0x04, 0x0a, 0x0d, 0x50, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x4c, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x65, 0x64,
	0x67, 0x65, 0x72, 0x12, 0x22, 0x2e, 0x70, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x12, 0x6e, 0x0a, 0x13, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x4f, 0x6e, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x12, 0x29, 0x2e, 0x70, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x4f, 0x6e, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x70, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x4f,
	0x6e, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x6b, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x46, 0x72, 0x6f,
	0x6d, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x12, 0x28, 0x2e, 0x70, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x46, 0x72, 0x6f, 0x6d, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x29, 0x2e, 0x70, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x47, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x46, 0x72, 0x6f, 0x6d, 0x4c, 0x65,
	0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4c,
	0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x12, 0x22,
	0x2e, 0x70, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x59, 0x0a, 0x0c,
	0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x22, 0x2e, 0x70,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x65, 0x64,
	0x67, 0x65, 0x72, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x23, 0x2e, 0x70, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x4c,
	0x65, 0x64, 0x67, 0x65, 0x72, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x22,
	0x2e, 0x70, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x4c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x22, 0x2e, 0x70, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x16, 0x5a, 0x14, 0x70, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_service_proto_rawDescOnce sync.Once
	file_service_proto_rawDescData = file_service_proto_rawDesc
)

func file_service_proto_rawDescGZIP() []byte {
	file_service_proto_rawDescOnce.Do(func() {
		file_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_proto_rawDescData)
	})
	return file_service_proto_rawDescData
}

var file_service_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_service_proto_goTypes = []any{
	(*CreateLedgerRequest)(nil),         // 0: porageservice.CreateLedgerRequest
	(*AppendEntryOnLedgerRequest)(nil),  // 1: porageservice.AppendEntryOnLedgerRequest
	(*AppendEntryOnLedgerResponse)(nil), // 2: porageservice.AppendEntryOnLedgerResponse
	(*GetEntryFromLedgerRequest)(nil),   // 3: porageservice.GetEntryFromLedgerRequest
	(*GetEntryFromLedgerResponse)(nil),  // 4: porageservice.GetEntryFromLedgerResponse
	(*DeleteLedgerRequest)(nil),         // 5: porageservice.DeleteLedgerRequest
	(*ListLedgersResponse)(nil),         // 6: porageservice.ListLedgersResponse
	(*ListWorkersResponse)(nil),         // 7: porageservice.ListWorkersResponse
	(*WorkerDescription)(nil),           // 8: porageservice.WorkerDescription
	(*LedgerLengthRequest)(nil),         // 9: porageservice.LedgerLengthRequest
	(*LedgerLengthResponse)(nil),        // 10: porageservice.LedgerLengthResponse
	nil,                                 // 11: porageservice.ListWorkersResponse.WorkersEntry
	(*emptypb.Empty)(nil),               // 12: google.protobuf.Empty
}
var file_service_proto_depIdxs = []int32{
	11, // 0: porageservice.ListWorkersResponse.workers:type_name -> porageservice.ListWorkersResponse.WorkersEntry
	8,  // 1: porageservice.ListWorkersResponse.WorkersEntry.value:type_name -> porageservice.WorkerDescription
	0,  // 2: porageservice.PorageService.CreateLedger:input_type -> porageservice.CreateLedgerRequest
	1,  // 3: porageservice.PorageService.AppendEntryOnLedger:input_type -> porageservice.AppendEntryOnLedgerRequest
	3,  // 4: porageservice.PorageService.GetEntryFromLedger:input_type -> porageservice.GetEntryFromLedgerRequest
	5,  // 5: porageservice.PorageService.DeleteLedger:input_type -> porageservice.DeleteLedgerRequest
	9,  // 6: porageservice.PorageService.LedgerLength:input_type -> porageservice.LedgerLengthRequest
	12, // 7: porageservice.PorageService.ListLedgers:input_type -> google.protobuf.Empty
	12, // 8: porageservice.PorageService.ListWorkers:input_type -> google.protobuf.Empty
	12, // 9: porageservice.PorageService.CreateLedger:output_type -> google.protobuf.Empty
	2,  // 10: porageservice.PorageService.AppendEntryOnLedger:output_type -> porageservice.AppendEntryOnLedgerResponse
	4,  // 11: porageservice.PorageService.GetEntryFromLedger:output_type -> porageservice.GetEntryFromLedgerResponse
	12, // 12: porageservice.PorageService.DeleteLedger:output_type -> google.protobuf.Empty
	10, // 13: porageservice.PorageService.LedgerLength:output_type -> porageservice.LedgerLengthResponse
	6,  // 14: porageservice.PorageService.ListLedgers:output_type -> porageservice.ListLedgersResponse
	7,  // 15: porageservice.PorageService.ListWorkers:output_type -> porageservice.ListWorkersResponse
	9,  // [9:16] is the sub-list for method output_type
	2,  // [2:9] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CreateLedgerRequest); i {
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
		file_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*AppendEntryOnLedgerRequest); i {
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
		file_service_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*AppendEntryOnLedgerResponse); i {
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
		file_service_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GetEntryFromLedgerRequest); i {
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
		file_service_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*GetEntryFromLedgerResponse); i {
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
		file_service_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteLedgerRequest); i {
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
		file_service_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*ListLedgersResponse); i {
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
		file_service_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*ListWorkersResponse); i {
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
		file_service_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*WorkerDescription); i {
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
		file_service_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*LedgerLengthRequest); i {
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
		file_service_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*LedgerLengthResponse); i {
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
			RawDescriptor: file_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_proto_goTypes,
		DependencyIndexes: file_service_proto_depIdxs,
		MessageInfos:      file_service_proto_msgTypes,
	}.Build()
	File_service_proto = out.File
	file_service_proto_rawDesc = nil
	file_service_proto_goTypes = nil
	file_service_proto_depIdxs = nil
}
