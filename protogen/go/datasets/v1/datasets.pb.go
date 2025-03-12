// Proto messages and service definition for the APIs related to creating and managing datasets.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: datasets/v1/datasets.proto

package datasetsv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// CreateDatasetRequest is used to create a new dataset.
type CreateDatasetRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// name of the dataset to create.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// message type of the dataset to create.
	Type *DatasetType `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	// short text summary of the dataset to create.
	Summary string `protobuf:"bytes,3,opt,name=summary,proto3" json:"summary,omitempty"`
	// normalized snake case name of the dataset to create.
	CodeName      string `protobuf:"bytes,4,opt,name=code_name,json=codeName,proto3" json:"code_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateDatasetRequest) Reset() {
	*x = CreateDatasetRequest{}
	mi := &file_datasets_v1_datasets_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateDatasetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateDatasetRequest) ProtoMessage() {}

func (x *CreateDatasetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_datasets_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateDatasetRequest.ProtoReflect.Descriptor instead.
func (*CreateDatasetRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_datasets_proto_rawDescGZIP(), []int{0}
}

func (x *CreateDatasetRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateDatasetRequest) GetType() *DatasetType {
	if x != nil {
		return x.Type
	}
	return nil
}

func (x *CreateDatasetRequest) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *CreateDatasetRequest) GetCodeName() string {
	if x != nil {
		return x.CodeName
	}
	return ""
}

// GetDatasetRequest is the request message for the GetDataset RPC method for fetching a single dataset
type GetDatasetRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// slug of the dataset to be returned, e.g. "open_data.copernicus.sentinel1_sar"
	Slug string `protobuf:"bytes,1,opt,name=slug,proto3" json:"slug,omitempty"`
	// or alternatively a dataset id
	Id            string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDatasetRequest) Reset() {
	*x = GetDatasetRequest{}
	mi := &file_datasets_v1_datasets_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDatasetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDatasetRequest) ProtoMessage() {}

func (x *GetDatasetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_datasets_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDatasetRequest.ProtoReflect.Descriptor instead.
func (*GetDatasetRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_datasets_proto_rawDescGZIP(), []int{1}
}

func (x *GetDatasetRequest) GetSlug() string {
	if x != nil {
		return x.Slug
	}
	return ""
}

func (x *GetDatasetRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// UpdateDatasetRequest is used to update a dataset.
type UpdateDatasetRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// id of the dataset to update.
	Id *ID `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// updated name of the dataset.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// updated type of the dataset.
	Type *DatasetType `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	// updated summary of the dataset.
	Summary       string `protobuf:"bytes,4,opt,name=summary,proto3" json:"summary,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateDatasetRequest) Reset() {
	*x = UpdateDatasetRequest{}
	mi := &file_datasets_v1_datasets_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateDatasetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDatasetRequest) ProtoMessage() {}

func (x *UpdateDatasetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_datasets_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDatasetRequest.ProtoReflect.Descriptor instead.
func (*UpdateDatasetRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_datasets_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateDatasetRequest) GetId() *ID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *UpdateDatasetRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateDatasetRequest) GetType() *DatasetType {
	if x != nil {
		return x.Type
	}
	return nil
}

func (x *UpdateDatasetRequest) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

// ClientInfo contains information about the client requesting datasets, useful for us to gather usage data
type ClientInfo struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// name of the client, e.g. "python"
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// environment the client is running in, e.g. "JupyterLab using python 3.11.4"
	Environment string `protobuf:"bytes,2,opt,name=environment,proto3" json:"environment,omitempty"`
	// list of packages installed on the client
	Packages      []*Package `protobuf:"bytes,3,rep,name=packages,proto3" json:"packages,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ClientInfo) Reset() {
	*x = ClientInfo{}
	mi := &file_datasets_v1_datasets_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClientInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientInfo) ProtoMessage() {}

func (x *ClientInfo) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_datasets_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientInfo.ProtoReflect.Descriptor instead.
func (*ClientInfo) Descriptor() ([]byte, []int) {
	return file_datasets_v1_datasets_proto_rawDescGZIP(), []int{3}
}

func (x *ClientInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ClientInfo) GetEnvironment() string {
	if x != nil {
		return x.Environment
	}
	return ""
}

func (x *ClientInfo) GetPackages() []*Package {
	if x != nil {
		return x.Packages
	}
	return nil
}

// Package contains information about the installed version of a given package on the client
type Package struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// package name
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// package version
	Version       string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Package) Reset() {
	*x = Package{}
	mi := &file_datasets_v1_datasets_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Package) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Package) ProtoMessage() {}

func (x *Package) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_datasets_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Package.ProtoReflect.Descriptor instead.
func (*Package) Descriptor() ([]byte, []int) {
	return file_datasets_v1_datasets_proto_rawDescGZIP(), []int{4}
}

func (x *Package) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Package) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

// UpdateDatasetDescriptionRequest is used to update a dataset description
type UpdateDatasetDescriptionRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// dataset id
	Id *ID `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// description of the dataset, in markdown format
	Description   string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateDatasetDescriptionRequest) Reset() {
	*x = UpdateDatasetDescriptionRequest{}
	mi := &file_datasets_v1_datasets_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateDatasetDescriptionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateDatasetDescriptionRequest) ProtoMessage() {}

func (x *UpdateDatasetDescriptionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_datasets_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateDatasetDescriptionRequest.ProtoReflect.Descriptor instead.
func (*UpdateDatasetDescriptionRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_datasets_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateDatasetDescriptionRequest) GetId() *ID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *UpdateDatasetDescriptionRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

// ListDatasetsRequest is used to request a list of datasets
type ListDatasetsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// some information about the Tilebox client requesting the datasets
	ClientInfo    *ClientInfo `protobuf:"bytes,1,opt,name=client_info,json=clientInfo,proto3" json:"client_info,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListDatasetsRequest) Reset() {
	*x = ListDatasetsRequest{}
	mi := &file_datasets_v1_datasets_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListDatasetsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDatasetsRequest) ProtoMessage() {}

func (x *ListDatasetsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_datasets_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDatasetsRequest.ProtoReflect.Descriptor instead.
func (*ListDatasetsRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_datasets_proto_rawDescGZIP(), []int{6}
}

func (x *ListDatasetsRequest) GetClientInfo() *ClientInfo {
	if x != nil {
		return x.ClientInfo
	}
	return nil
}

// A list of datasets and dataset groups
type ListDatasetsResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// list of datasets a user has access to
	Datasets []*Dataset `protobuf:"bytes,1,rep,name=datasets,proto3" json:"datasets,omitempty"`
	// list of dataset groups a user has access to
	Groups []*DatasetGroup `protobuf:"bytes,2,rep,name=groups,proto3" json:"groups,omitempty"`
	// an optional message to be displayed to the user when they access a list of datasets
	ServerMessage string `protobuf:"bytes,3,opt,name=server_message,json=serverMessage,proto3" json:"server_message,omitempty"`
	// Number of datasets created by the organization
	NbCreateDatasets int64 `protobuf:"varint,4,opt,name=nb_create_datasets,json=nbCreateDatasets,proto3" json:"nb_create_datasets,omitempty"`
	// Maximum number of datasets created by the organization
	NbCreatedDatasetsLimit int64 `protobuf:"varint,5,opt,name=nb_created_datasets_limit,json=nbCreatedDatasetsLimit,proto3" json:"nb_created_datasets_limit,omitempty"`
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *ListDatasetsResponse) Reset() {
	*x = ListDatasetsResponse{}
	mi := &file_datasets_v1_datasets_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListDatasetsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListDatasetsResponse) ProtoMessage() {}

func (x *ListDatasetsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_datasets_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListDatasetsResponse.ProtoReflect.Descriptor instead.
func (*ListDatasetsResponse) Descriptor() ([]byte, []int) {
	return file_datasets_v1_datasets_proto_rawDescGZIP(), []int{7}
}

func (x *ListDatasetsResponse) GetDatasets() []*Dataset {
	if x != nil {
		return x.Datasets
	}
	return nil
}

func (x *ListDatasetsResponse) GetGroups() []*DatasetGroup {
	if x != nil {
		return x.Groups
	}
	return nil
}

func (x *ListDatasetsResponse) GetServerMessage() string {
	if x != nil {
		return x.ServerMessage
	}
	return ""
}

func (x *ListDatasetsResponse) GetNbCreateDatasets() int64 {
	if x != nil {
		return x.NbCreateDatasets
	}
	return 0
}

func (x *ListDatasetsResponse) GetNbCreatedDatasetsLimit() int64 {
	if x != nil {
		return x.NbCreatedDatasetsLimit
	}
	return 0
}

var File_datasets_v1_datasets_proto protoreflect.FileDescriptor

var file_datasets_v1_datasets_proto_rawDesc = string([]byte{
	0x0a, 0x1a, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x61,
	0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x64, 0x61,
	0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x16, 0x64, 0x61, 0x74, 0x61, 0x73,
	0x65, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x64,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x8f, 0x01, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2c,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x73,
	0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73,
	0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6f, 0x64, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f, 0x64, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0x37, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6c, 0x75, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x93, 0x01, 0x0a,
	0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x49, 0x44, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2c, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73,
	0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d,
	0x61, 0x72, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61,
	0x72, 0x79, 0x22, 0x74, 0x0a, 0x0a, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72,
	0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x30, 0x0a, 0x08, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67,
	0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73,
	0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x52, 0x08,
	0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x22, 0x37, 0x0a, 0x07, 0x50, 0x61, 0x63, 0x6b,
	0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x22, 0x64, 0x0a, 0x1f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x73,
	0x65, 0x74, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49,
	0x44, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x4f, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x44,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x38,
	0x0a, 0x0b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x8b, 0x02, 0x0a, 0x14, 0x4c, 0x69, 0x73,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x30, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x73,
	0x65, 0x74, 0x73, 0x12, 0x31, 0x0a, 0x06, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x06,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2c, 0x0a,
	0x12, 0x6e, 0x62, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x73,
	0x65, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x6e, 0x62, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x12, 0x39, 0x0a, 0x19, 0x6e,
	0x62, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65,
	0x74, 0x73, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x16,
	0x6e, 0x62, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74,
	0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x32, 0xa7, 0x03, 0x0a, 0x0e, 0x44, 0x61, 0x74, 0x61, 0x73,
	0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x12, 0x21, 0x2e, 0x64, 0x61, 0x74,
	0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e,
	0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x12, 0x1e, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x0d, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x12, 0x21, 0x2e, 0x64,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x14, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x73, 0x65, 0x74, 0x22, 0x00, 0x12, 0x60, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x44,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x14, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x0c, 0x4c, 0x69, 0x73,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x12, 0x20, 0x2e, 0x64, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x64, 0x61,
	0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0xaf, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74,
	0x73, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f, 0x78, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f,
	0x78, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f,
	0x2f, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x64, 0x61, 0x74,
	0x61, 0x73, 0x65, 0x74, 0x73, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x44, 0x58, 0x58, 0xaa, 0x02, 0x0b,
	0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b, 0x44, 0x61,
	0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x17, 0x44, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x0c, 0x44, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x73, 0x3a, 0x3a,
	0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_datasets_v1_datasets_proto_rawDescOnce sync.Once
	file_datasets_v1_datasets_proto_rawDescData []byte
)

func file_datasets_v1_datasets_proto_rawDescGZIP() []byte {
	file_datasets_v1_datasets_proto_rawDescOnce.Do(func() {
		file_datasets_v1_datasets_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_datasets_v1_datasets_proto_rawDesc), len(file_datasets_v1_datasets_proto_rawDesc)))
	})
	return file_datasets_v1_datasets_proto_rawDescData
}

var file_datasets_v1_datasets_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_datasets_v1_datasets_proto_goTypes = []any{
	(*CreateDatasetRequest)(nil),            // 0: datasets.v1.CreateDatasetRequest
	(*GetDatasetRequest)(nil),               // 1: datasets.v1.GetDatasetRequest
	(*UpdateDatasetRequest)(nil),            // 2: datasets.v1.UpdateDatasetRequest
	(*ClientInfo)(nil),                      // 3: datasets.v1.ClientInfo
	(*Package)(nil),                         // 4: datasets.v1.Package
	(*UpdateDatasetDescriptionRequest)(nil), // 5: datasets.v1.UpdateDatasetDescriptionRequest
	(*ListDatasetsRequest)(nil),             // 6: datasets.v1.ListDatasetsRequest
	(*ListDatasetsResponse)(nil),            // 7: datasets.v1.ListDatasetsResponse
	(*DatasetType)(nil),                     // 8: datasets.v1.DatasetType
	(*ID)(nil),                              // 9: datasets.v1.ID
	(*Dataset)(nil),                         // 10: datasets.v1.Dataset
	(*DatasetGroup)(nil),                    // 11: datasets.v1.DatasetGroup
}
var file_datasets_v1_datasets_proto_depIdxs = []int32{
	8,  // 0: datasets.v1.CreateDatasetRequest.type:type_name -> datasets.v1.DatasetType
	9,  // 1: datasets.v1.UpdateDatasetRequest.id:type_name -> datasets.v1.ID
	8,  // 2: datasets.v1.UpdateDatasetRequest.type:type_name -> datasets.v1.DatasetType
	4,  // 3: datasets.v1.ClientInfo.packages:type_name -> datasets.v1.Package
	9,  // 4: datasets.v1.UpdateDatasetDescriptionRequest.id:type_name -> datasets.v1.ID
	3,  // 5: datasets.v1.ListDatasetsRequest.client_info:type_name -> datasets.v1.ClientInfo
	10, // 6: datasets.v1.ListDatasetsResponse.datasets:type_name -> datasets.v1.Dataset
	11, // 7: datasets.v1.ListDatasetsResponse.groups:type_name -> datasets.v1.DatasetGroup
	0,  // 8: datasets.v1.DatasetService.CreateDataset:input_type -> datasets.v1.CreateDatasetRequest
	1,  // 9: datasets.v1.DatasetService.GetDataset:input_type -> datasets.v1.GetDatasetRequest
	2,  // 10: datasets.v1.DatasetService.UpdateDataset:input_type -> datasets.v1.UpdateDatasetRequest
	5,  // 11: datasets.v1.DatasetService.UpdateDatasetDescription:input_type -> datasets.v1.UpdateDatasetDescriptionRequest
	6,  // 12: datasets.v1.DatasetService.ListDatasets:input_type -> datasets.v1.ListDatasetsRequest
	10, // 13: datasets.v1.DatasetService.CreateDataset:output_type -> datasets.v1.Dataset
	10, // 14: datasets.v1.DatasetService.GetDataset:output_type -> datasets.v1.Dataset
	10, // 15: datasets.v1.DatasetService.UpdateDataset:output_type -> datasets.v1.Dataset
	10, // 16: datasets.v1.DatasetService.UpdateDatasetDescription:output_type -> datasets.v1.Dataset
	7,  // 17: datasets.v1.DatasetService.ListDatasets:output_type -> datasets.v1.ListDatasetsResponse
	13, // [13:18] is the sub-list for method output_type
	8,  // [8:13] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_datasets_v1_datasets_proto_init() }
func file_datasets_v1_datasets_proto_init() {
	if File_datasets_v1_datasets_proto != nil {
		return
	}
	file_datasets_v1_core_proto_init()
	file_datasets_v1_dataset_type_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_datasets_v1_datasets_proto_rawDesc), len(file_datasets_v1_datasets_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_datasets_v1_datasets_proto_goTypes,
		DependencyIndexes: file_datasets_v1_datasets_proto_depIdxs,
		MessageInfos:      file_datasets_v1_datasets_proto_msgTypes,
	}.Build()
	File_datasets_v1_datasets_proto = out.File
	file_datasets_v1_datasets_proto_goTypes = nil
	file_datasets_v1_datasets_proto_depIdxs = nil
}
