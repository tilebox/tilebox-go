// Proto messages and service definition for the APIs related to loading and querying data.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: datasets/v1/data_access.proto

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

// Legacy message, to be removed in the future.
// GetDatasetForIntervalRequest contains the request parameters for retrieving data for a time interval.
type GetDatasetForIntervalRequest struct {
	state        protoimpl.MessageState `protogen:"open.v1"`
	CollectionId string                 `protobuf:"bytes,1,opt,name=collection_id,json=collectionId,proto3" json:"collection_id,omitempty"` // The collection id.
	// Either time interval or datapoint interval must be set, but not both.
	TimeInterval      *TimeInterval      `protobuf:"bytes,2,opt,name=time_interval,json=timeInterval,proto3" json:"time_interval,omitempty"`                // The time interval for which data is requested.
	DatapointInterval *DatapointInterval `protobuf:"bytes,6,opt,name=datapoint_interval,json=datapointInterval,proto3" json:"datapoint_interval,omitempty"` // The datapoint interval for which data is requested.
	Page              *LegacyPagination  `protobuf:"bytes,3,opt,name=page,proto3,oneof" json:"page,omitempty"`                                              // The pagination parameters for this request.
	SkipData          bool               `protobuf:"varint,4,opt,name=skip_data,json=skipData,proto3" json:"skip_data,omitempty"`                           // If true, the datapoint data is not returned.
	// If true, the datapoint metadata is not returned.
	// If both skip_data and skip_meta are true,
	// the response will only consist of a list of datapoint ids without any additional data or metadata.
	SkipMeta      bool `protobuf:"varint,5,opt,name=skip_meta,json=skipMeta,proto3" json:"skip_meta,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDatasetForIntervalRequest) Reset() {
	*x = GetDatasetForIntervalRequest{}
	mi := &file_datasets_v1_data_access_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDatasetForIntervalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDatasetForIntervalRequest) ProtoMessage() {}

func (x *GetDatasetForIntervalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_data_access_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDatasetForIntervalRequest.ProtoReflect.Descriptor instead.
func (*GetDatasetForIntervalRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_data_access_proto_rawDescGZIP(), []int{0}
}

func (x *GetDatasetForIntervalRequest) GetCollectionId() string {
	if x != nil {
		return x.CollectionId
	}
	return ""
}

func (x *GetDatasetForIntervalRequest) GetTimeInterval() *TimeInterval {
	if x != nil {
		return x.TimeInterval
	}
	return nil
}

func (x *GetDatasetForIntervalRequest) GetDatapointInterval() *DatapointInterval {
	if x != nil {
		return x.DatapointInterval
	}
	return nil
}

func (x *GetDatasetForIntervalRequest) GetPage() *LegacyPagination {
	if x != nil {
		return x.Page
	}
	return nil
}

func (x *GetDatasetForIntervalRequest) GetSkipData() bool {
	if x != nil {
		return x.SkipData
	}
	return false
}

func (x *GetDatasetForIntervalRequest) GetSkipMeta() bool {
	if x != nil {
		return x.SkipMeta
	}
	return false
}

// Legacy message, to be removed in the future.
// GetDatapointByIdRequest contains the request parameters for retrieving a single data point in a collection by its id.
type GetDatapointByIdRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CollectionId  string                 `protobuf:"bytes,1,opt,name=collection_id,json=collectionId,proto3" json:"collection_id,omitempty"` // The collection id.
	Id            string                 `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`                                         // The id of the requested data point.
	SkipData      bool                   `protobuf:"varint,3,opt,name=skip_data,json=skipData,proto3" json:"skip_data,omitempty"`            // If true, only the datapoint metadata is returned.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetDatapointByIdRequest) Reset() {
	*x = GetDatapointByIdRequest{}
	mi := &file_datasets_v1_data_access_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetDatapointByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDatapointByIdRequest) ProtoMessage() {}

func (x *GetDatapointByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_data_access_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDatapointByIdRequest.ProtoReflect.Descriptor instead.
func (*GetDatapointByIdRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_data_access_proto_rawDescGZIP(), []int{1}
}

func (x *GetDatapointByIdRequest) GetCollectionId() string {
	if x != nil {
		return x.CollectionId
	}
	return ""
}

func (x *GetDatapointByIdRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetDatapointByIdRequest) GetSkipData() bool {
	if x != nil {
		return x.SkipData
	}
	return false
}

// QueryByIDRequest contains the request parameters for retrieving a single data point by its id.
type QueryByIDRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CollectionIds []*ID                  `protobuf:"bytes,1,rep,name=collection_ids,json=collectionIds,proto3" json:"collection_ids,omitempty"` // collection ids to query.
	Id            *ID                    `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`                                            // The id of the requested data point.
	SkipData      bool                   `protobuf:"varint,3,opt,name=skip_data,json=skipData,proto3" json:"skip_data,omitempty"`               // If true, only the datapoint metadata is returned.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueryByIDRequest) Reset() {
	*x = QueryByIDRequest{}
	mi := &file_datasets_v1_data_access_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryByIDRequest) ProtoMessage() {}

func (x *QueryByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_data_access_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryByIDRequest.ProtoReflect.Descriptor instead.
func (*QueryByIDRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_data_access_proto_rawDescGZIP(), []int{2}
}

func (x *QueryByIDRequest) GetCollectionIds() []*ID {
	if x != nil {
		return x.CollectionIds
	}
	return nil
}

func (x *QueryByIDRequest) GetId() *ID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *QueryByIDRequest) GetSkipData() bool {
	if x != nil {
		return x.SkipData
	}
	return false
}

// QueryFilters contains the filters to apply to a query.
type QueryFilters struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Either time interval or datapoint interval must be set, but not both.
	//
	// Types that are valid to be assigned to TemporalInterval:
	//
	//	*QueryFilters_TimeInterval
	//	*QueryFilters_DatapointInterval
	TemporalInterval isQueryFilters_TemporalInterval `protobuf_oneof:"temporal_interval"`
	AreaOfInterest   *Geometry                       `protobuf:"bytes,3,opt,name=area_of_interest,json=areaOfInterest,proto3" json:"area_of_interest,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *QueryFilters) Reset() {
	*x = QueryFilters{}
	mi := &file_datasets_v1_data_access_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryFilters) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryFilters) ProtoMessage() {}

func (x *QueryFilters) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_data_access_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryFilters.ProtoReflect.Descriptor instead.
func (*QueryFilters) Descriptor() ([]byte, []int) {
	return file_datasets_v1_data_access_proto_rawDescGZIP(), []int{3}
}

func (x *QueryFilters) GetTemporalInterval() isQueryFilters_TemporalInterval {
	if x != nil {
		return x.TemporalInterval
	}
	return nil
}

func (x *QueryFilters) GetTimeInterval() *TimeInterval {
	if x != nil {
		if x, ok := x.TemporalInterval.(*QueryFilters_TimeInterval); ok {
			return x.TimeInterval
		}
	}
	return nil
}

func (x *QueryFilters) GetDatapointInterval() *DatapointInterval {
	if x != nil {
		if x, ok := x.TemporalInterval.(*QueryFilters_DatapointInterval); ok {
			return x.DatapointInterval
		}
	}
	return nil
}

func (x *QueryFilters) GetAreaOfInterest() *Geometry {
	if x != nil {
		return x.AreaOfInterest
	}
	return nil
}

type isQueryFilters_TemporalInterval interface {
	isQueryFilters_TemporalInterval()
}

type QueryFilters_TimeInterval struct {
	TimeInterval *TimeInterval `protobuf:"bytes,1,opt,name=time_interval,json=timeInterval,proto3,oneof"`
}

type QueryFilters_DatapointInterval struct {
	DatapointInterval *DatapointInterval `protobuf:"bytes,2,opt,name=datapoint_interval,json=datapointInterval,proto3,oneof"`
}

func (*QueryFilters_TimeInterval) isQueryFilters_TemporalInterval() {}

func (*QueryFilters_DatapointInterval) isQueryFilters_TemporalInterval() {}

// QueryRequest contains the request parameters for retrieving data from a Tilebox dataset.
type QueryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CollectionIds []*ID                  `protobuf:"bytes,1,rep,name=collection_ids,json=collectionIds,proto3" json:"collection_ids,omitempty"` // collection ids to query.
	Filters       *QueryFilters          `protobuf:"bytes,2,opt,name=filters,proto3" json:"filters,omitempty"`                                  // Filters to apply to the query.
	Page          *Pagination            `protobuf:"bytes,3,opt,name=page,proto3,oneof" json:"page,omitempty"`                                  // The pagination parameters for this request.
	SkipData      bool                   `protobuf:"varint,4,opt,name=skip_data,json=skipData,proto3" json:"skip_data,omitempty"`               // If true, only datapoint metadata, such as id, time and ingestion_time are returned.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueryRequest) Reset() {
	*x = QueryRequest{}
	mi := &file_datasets_v1_data_access_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRequest) ProtoMessage() {}

func (x *QueryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_data_access_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRequest.ProtoReflect.Descriptor instead.
func (*QueryRequest) Descriptor() ([]byte, []int) {
	return file_datasets_v1_data_access_proto_rawDescGZIP(), []int{4}
}

func (x *QueryRequest) GetCollectionIds() []*ID {
	if x != nil {
		return x.CollectionIds
	}
	return nil
}

func (x *QueryRequest) GetFilters() *QueryFilters {
	if x != nil {
		return x.Filters
	}
	return nil
}

func (x *QueryRequest) GetPage() *Pagination {
	if x != nil {
		return x.Page
	}
	return nil
}

func (x *QueryRequest) GetSkipData() bool {
	if x != nil {
		return x.SkipData
	}
	return false
}

// QueryResultPage is a single page of data points of a Tilebox dataset
type QueryResultPage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          *RepeatedAny           `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`                               // The datapoints.
	NextPage      *Pagination            `protobuf:"bytes,2,opt,name=next_page,json=nextPage,proto3,oneof" json:"next_page,omitempty"` // The pagination parameters for the next page.
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueryResultPage) Reset() {
	*x = QueryResultPage{}
	mi := &file_datasets_v1_data_access_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryResultPage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryResultPage) ProtoMessage() {}

func (x *QueryResultPage) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_data_access_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryResultPage.ProtoReflect.Descriptor instead.
func (*QueryResultPage) Descriptor() ([]byte, []int) {
	return file_datasets_v1_data_access_proto_rawDescGZIP(), []int{5}
}

func (x *QueryResultPage) GetData() *RepeatedAny {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *QueryResultPage) GetNextPage() *Pagination {
	if x != nil {
		return x.NextPage
	}
	return nil
}

var File_datasets_v1_data_access_proto protoreflect.FileDescriptor

const file_datasets_v1_data_access_proto_rawDesc = "" +
	"\n" +
	"\x1ddatasets/v1/data_access.proto\x12\vdatasets.v1\x1a\x16datasets/v1/core.proto\x1a\"datasets/v1/well_known_types.proto\"\xcd\x02\n" +
	"\x1cGetDatasetForIntervalRequest\x12#\n" +
	"\rcollection_id\x18\x01 \x01(\tR\fcollectionId\x12>\n" +
	"\rtime_interval\x18\x02 \x01(\v2\x19.datasets.v1.TimeIntervalR\ftimeInterval\x12M\n" +
	"\x12datapoint_interval\x18\x06 \x01(\v2\x1e.datasets.v1.DatapointIntervalR\x11datapointInterval\x126\n" +
	"\x04page\x18\x03 \x01(\v2\x1d.datasets.v1.LegacyPaginationH\x00R\x04page\x88\x01\x01\x12\x1b\n" +
	"\tskip_data\x18\x04 \x01(\bR\bskipData\x12\x1b\n" +
	"\tskip_meta\x18\x05 \x01(\bR\bskipMetaB\a\n" +
	"\x05_page\"k\n" +
	"\x17GetDatapointByIdRequest\x12#\n" +
	"\rcollection_id\x18\x01 \x01(\tR\fcollectionId\x12\x0e\n" +
	"\x02id\x18\x02 \x01(\tR\x02id\x12\x1b\n" +
	"\tskip_data\x18\x03 \x01(\bR\bskipData\"\x88\x01\n" +
	"\x10QueryByIDRequest\x126\n" +
	"\x0ecollection_ids\x18\x01 \x03(\v2\x0f.datasets.v1.IDR\rcollectionIds\x12\x1f\n" +
	"\x02id\x18\x02 \x01(\v2\x0f.datasets.v1.IDR\x02id\x12\x1b\n" +
	"\tskip_data\x18\x03 \x01(\bR\bskipData\"\xf7\x01\n" +
	"\fQueryFilters\x12@\n" +
	"\rtime_interval\x18\x01 \x01(\v2\x19.datasets.v1.TimeIntervalH\x00R\ftimeInterval\x12O\n" +
	"\x12datapoint_interval\x18\x02 \x01(\v2\x1e.datasets.v1.DatapointIntervalH\x00R\x11datapointInterval\x12?\n" +
	"\x10area_of_interest\x18\x03 \x01(\v2\x15.datasets.v1.GeometryR\x0eareaOfInterestB\x13\n" +
	"\x11temporal_interval\"\xd3\x01\n" +
	"\fQueryRequest\x126\n" +
	"\x0ecollection_ids\x18\x01 \x03(\v2\x0f.datasets.v1.IDR\rcollectionIds\x123\n" +
	"\afilters\x18\x02 \x01(\v2\x19.datasets.v1.QueryFiltersR\afilters\x120\n" +
	"\x04page\x18\x03 \x01(\v2\x17.datasets.v1.PaginationH\x00R\x04page\x88\x01\x01\x12\x1b\n" +
	"\tskip_data\x18\x04 \x01(\bR\bskipDataB\a\n" +
	"\x05_page\"\x88\x01\n" +
	"\x0fQueryResultPage\x12,\n" +
	"\x04data\x18\x01 \x01(\v2\x18.datasets.v1.RepeatedAnyR\x04data\x129\n" +
	"\tnext_page\x18\x02 \x01(\v2\x17.datasets.v1.PaginationH\x00R\bnextPage\x88\x01\x01B\f\n" +
	"\n" +
	"_next_page2\xcd\x02\n" +
	"\x11DataAccessService\x12`\n" +
	"\x15GetDatasetForInterval\x12).datasets.v1.GetDatasetForIntervalRequest\x1a\x1a.datasets.v1.DatapointPage\"\x00\x12R\n" +
	"\x10GetDatapointByID\x12$.datasets.v1.GetDatapointByIdRequest\x1a\x16.datasets.v1.Datapoint\"\x00\x12>\n" +
	"\tQueryByID\x12\x1d.datasets.v1.QueryByIDRequest\x1a\x10.datasets.v1.Any\"\x00\x12B\n" +
	"\x05Query\x12\x19.datasets.v1.QueryRequest\x1a\x1c.datasets.v1.QueryResultPage\"\x00B\xb1\x01\n" +
	"\x0fcom.datasets.v1B\x0fDataAccessProtoP\x01Z@github.com/tilebox/tilebox-go/protogen/go/datasets/v1;datasetsv1\xa2\x02\x03DXX\xaa\x02\vDatasets.V1\xca\x02\vDatasets\\V1\xe2\x02\x17Datasets\\V1\\GPBMetadata\xea\x02\fDatasets::V1b\x06proto3"

var (
	file_datasets_v1_data_access_proto_rawDescOnce sync.Once
	file_datasets_v1_data_access_proto_rawDescData []byte
)

func file_datasets_v1_data_access_proto_rawDescGZIP() []byte {
	file_datasets_v1_data_access_proto_rawDescOnce.Do(func() {
		file_datasets_v1_data_access_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_datasets_v1_data_access_proto_rawDesc), len(file_datasets_v1_data_access_proto_rawDesc)))
	})
	return file_datasets_v1_data_access_proto_rawDescData
}

var file_datasets_v1_data_access_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_datasets_v1_data_access_proto_goTypes = []any{
	(*GetDatasetForIntervalRequest)(nil), // 0: datasets.v1.GetDatasetForIntervalRequest
	(*GetDatapointByIdRequest)(nil),      // 1: datasets.v1.GetDatapointByIdRequest
	(*QueryByIDRequest)(nil),             // 2: datasets.v1.QueryByIDRequest
	(*QueryFilters)(nil),                 // 3: datasets.v1.QueryFilters
	(*QueryRequest)(nil),                 // 4: datasets.v1.QueryRequest
	(*QueryResultPage)(nil),              // 5: datasets.v1.QueryResultPage
	(*TimeInterval)(nil),                 // 6: datasets.v1.TimeInterval
	(*DatapointInterval)(nil),            // 7: datasets.v1.DatapointInterval
	(*LegacyPagination)(nil),             // 8: datasets.v1.LegacyPagination
	(*ID)(nil),                           // 9: datasets.v1.ID
	(*Geometry)(nil),                     // 10: datasets.v1.Geometry
	(*Pagination)(nil),                   // 11: datasets.v1.Pagination
	(*RepeatedAny)(nil),                  // 12: datasets.v1.RepeatedAny
	(*DatapointPage)(nil),                // 13: datasets.v1.DatapointPage
	(*Datapoint)(nil),                    // 14: datasets.v1.Datapoint
	(*Any)(nil),                          // 15: datasets.v1.Any
}
var file_datasets_v1_data_access_proto_depIdxs = []int32{
	6,  // 0: datasets.v1.GetDatasetForIntervalRequest.time_interval:type_name -> datasets.v1.TimeInterval
	7,  // 1: datasets.v1.GetDatasetForIntervalRequest.datapoint_interval:type_name -> datasets.v1.DatapointInterval
	8,  // 2: datasets.v1.GetDatasetForIntervalRequest.page:type_name -> datasets.v1.LegacyPagination
	9,  // 3: datasets.v1.QueryByIDRequest.collection_ids:type_name -> datasets.v1.ID
	9,  // 4: datasets.v1.QueryByIDRequest.id:type_name -> datasets.v1.ID
	6,  // 5: datasets.v1.QueryFilters.time_interval:type_name -> datasets.v1.TimeInterval
	7,  // 6: datasets.v1.QueryFilters.datapoint_interval:type_name -> datasets.v1.DatapointInterval
	10, // 7: datasets.v1.QueryFilters.area_of_interest:type_name -> datasets.v1.Geometry
	9,  // 8: datasets.v1.QueryRequest.collection_ids:type_name -> datasets.v1.ID
	3,  // 9: datasets.v1.QueryRequest.filters:type_name -> datasets.v1.QueryFilters
	11, // 10: datasets.v1.QueryRequest.page:type_name -> datasets.v1.Pagination
	12, // 11: datasets.v1.QueryResultPage.data:type_name -> datasets.v1.RepeatedAny
	11, // 12: datasets.v1.QueryResultPage.next_page:type_name -> datasets.v1.Pagination
	0,  // 13: datasets.v1.DataAccessService.GetDatasetForInterval:input_type -> datasets.v1.GetDatasetForIntervalRequest
	1,  // 14: datasets.v1.DataAccessService.GetDatapointByID:input_type -> datasets.v1.GetDatapointByIdRequest
	2,  // 15: datasets.v1.DataAccessService.QueryByID:input_type -> datasets.v1.QueryByIDRequest
	4,  // 16: datasets.v1.DataAccessService.Query:input_type -> datasets.v1.QueryRequest
	13, // 17: datasets.v1.DataAccessService.GetDatasetForInterval:output_type -> datasets.v1.DatapointPage
	14, // 18: datasets.v1.DataAccessService.GetDatapointByID:output_type -> datasets.v1.Datapoint
	15, // 19: datasets.v1.DataAccessService.QueryByID:output_type -> datasets.v1.Any
	5,  // 20: datasets.v1.DataAccessService.Query:output_type -> datasets.v1.QueryResultPage
	17, // [17:21] is the sub-list for method output_type
	13, // [13:17] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_datasets_v1_data_access_proto_init() }
func file_datasets_v1_data_access_proto_init() {
	if File_datasets_v1_data_access_proto != nil {
		return
	}
	file_datasets_v1_core_proto_init()
	file_datasets_v1_well_known_types_proto_init()
	file_datasets_v1_data_access_proto_msgTypes[0].OneofWrappers = []any{}
	file_datasets_v1_data_access_proto_msgTypes[3].OneofWrappers = []any{
		(*QueryFilters_TimeInterval)(nil),
		(*QueryFilters_DatapointInterval)(nil),
	}
	file_datasets_v1_data_access_proto_msgTypes[4].OneofWrappers = []any{}
	file_datasets_v1_data_access_proto_msgTypes[5].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_datasets_v1_data_access_proto_rawDesc), len(file_datasets_v1_data_access_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_datasets_v1_data_access_proto_goTypes,
		DependencyIndexes: file_datasets_v1_data_access_proto_depIdxs,
		MessageInfos:      file_datasets_v1_data_access_proto_msgTypes,
	}.Build()
	File_datasets_v1_data_access_proto = out.File
	file_datasets_v1_data_access_proto_goTypes = nil
	file_datasets_v1_data_access_proto_depIdxs = nil
}
