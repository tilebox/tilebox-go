// This file contains some well-known types that are available to use for all dataset protobuf messages.
// They are specially handled in our client libraries and converted to the corresponding representations.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: datasets/v1/well_known_types.proto

package datasetsv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

// Satellite flight direction
type FlightDirection int32

const (
	FlightDirection_FLIGHT_DIRECTION_UNSPECIFIED FlightDirection = 0
	FlightDirection_FLIGHT_DIRECTION_ASCENDING   FlightDirection = 1
	FlightDirection_FLIGHT_DIRECTION_DESCENDING  FlightDirection = 2
)

// Enum value maps for FlightDirection.
var (
	FlightDirection_name = map[int32]string{
		0: "FLIGHT_DIRECTION_UNSPECIFIED",
		1: "FLIGHT_DIRECTION_ASCENDING",
		2: "FLIGHT_DIRECTION_DESCENDING",
	}
	FlightDirection_value = map[string]int32{
		"FLIGHT_DIRECTION_UNSPECIFIED": 0,
		"FLIGHT_DIRECTION_ASCENDING":   1,
		"FLIGHT_DIRECTION_DESCENDING":  2,
	}
)

func (x FlightDirection) Enum() *FlightDirection {
	p := new(FlightDirection)
	*p = x
	return p
}

func (x FlightDirection) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FlightDirection) Descriptor() protoreflect.EnumDescriptor {
	return file_datasets_v1_well_known_types_proto_enumTypes[0].Descriptor()
}

func (FlightDirection) Type() protoreflect.EnumType {
	return &file_datasets_v1_well_known_types_proto_enumTypes[0]
}

func (x FlightDirection) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FlightDirection.Descriptor instead.
func (FlightDirection) EnumDescriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{0}
}

// Observation direction
type ObservationDirection int32

const (
	ObservationDirection_OBSERVATION_DIRECTION_UNSPECIFIED ObservationDirection = 0
	ObservationDirection_OBSERVATION_DIRECTION_LEFT        ObservationDirection = 1
	ObservationDirection_OBSERVATION_DIRECTION_RIGHT       ObservationDirection = 2
)

// Enum value maps for ObservationDirection.
var (
	ObservationDirection_name = map[int32]string{
		0: "OBSERVATION_DIRECTION_UNSPECIFIED",
		1: "OBSERVATION_DIRECTION_LEFT",
		2: "OBSERVATION_DIRECTION_RIGHT",
	}
	ObservationDirection_value = map[string]int32{
		"OBSERVATION_DIRECTION_UNSPECIFIED": 0,
		"OBSERVATION_DIRECTION_LEFT":        1,
		"OBSERVATION_DIRECTION_RIGHT":       2,
	}
)

func (x ObservationDirection) Enum() *ObservationDirection {
	p := new(ObservationDirection)
	*p = x
	return p
}

func (x ObservationDirection) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ObservationDirection) Descriptor() protoreflect.EnumDescriptor {
	return file_datasets_v1_well_known_types_proto_enumTypes[1].Descriptor()
}

func (ObservationDirection) Type() protoreflect.EnumType {
	return &file_datasets_v1_well_known_types_proto_enumTypes[1]
}

func (x ObservationDirection) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ObservationDirection.Descriptor instead.
func (ObservationDirection) EnumDescriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{1}
}

// The open data provider.
type OpendataProvider int32

const (
	OpendataProvider_OPENDATA_PROVIDER_UNSPECIFIED          OpendataProvider = 0
	OpendataProvider_OPENDATA_PROVIDER_ASF                  OpendataProvider = 1 // Alaska Satellite Facility
	OpendataProvider_OPENDATA_PROVIDER_COPERNICUS_DATASPACE OpendataProvider = 2 // Copernicus Dataspace
	OpendataProvider_OPENDATA_PROVIDER_UMBRA                OpendataProvider = 3 // Umbra Space
)

// Enum value maps for OpendataProvider.
var (
	OpendataProvider_name = map[int32]string{
		0: "OPENDATA_PROVIDER_UNSPECIFIED",
		1: "OPENDATA_PROVIDER_ASF",
		2: "OPENDATA_PROVIDER_COPERNICUS_DATASPACE",
		3: "OPENDATA_PROVIDER_UMBRA",
	}
	OpendataProvider_value = map[string]int32{
		"OPENDATA_PROVIDER_UNSPECIFIED":          0,
		"OPENDATA_PROVIDER_ASF":                  1,
		"OPENDATA_PROVIDER_COPERNICUS_DATASPACE": 2,
		"OPENDATA_PROVIDER_UMBRA":                3,
	}
)

func (x OpendataProvider) Enum() *OpendataProvider {
	p := new(OpendataProvider)
	*p = x
	return p
}

func (x OpendataProvider) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OpendataProvider) Descriptor() protoreflect.EnumDescriptor {
	return file_datasets_v1_well_known_types_proto_enumTypes[2].Descriptor()
}

func (OpendataProvider) Type() protoreflect.EnumType {
	return &file_datasets_v1_well_known_types_proto_enumTypes[2]
}

func (x OpendataProvider) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OpendataProvider.Descriptor instead.
func (OpendataProvider) EnumDescriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{2}
}

// https://www.earthdata.nasa.gov/engage/open-data-services-and-software/data-information-policy/data-levels
type ProcessingLevel int32

const (
	ProcessingLevel_PROCESSING_LEVEL_UNSPECIFIED    ProcessingLevel = 0
	ProcessingLevel_PROCESSING_LEVEL_L0             ProcessingLevel = 12 // Raw data
	ProcessingLevel_PROCESSING_LEVEL_L1             ProcessingLevel = 10
	ProcessingLevel_PROCESSING_LEVEL_L1A            ProcessingLevel = 1
	ProcessingLevel_PROCESSING_LEVEL_L1B            ProcessingLevel = 2
	ProcessingLevel_PROCESSING_LEVEL_L1C            ProcessingLevel = 3
	ProcessingLevel_PROCESSING_LEVEL_L2             ProcessingLevel = 4
	ProcessingLevel_PROCESSING_LEVEL_L2A            ProcessingLevel = 5
	ProcessingLevel_PROCESSING_LEVEL_L2B            ProcessingLevel = 6
	ProcessingLevel_PROCESSING_LEVEL_L3             ProcessingLevel = 7
	ProcessingLevel_PROCESSING_LEVEL_L3A            ProcessingLevel = 8
	ProcessingLevel_PROCESSING_LEVEL_L4             ProcessingLevel = 9
	ProcessingLevel_PROCESSING_LEVEL_NOT_APPLICABLE ProcessingLevel = 11
)

// Enum value maps for ProcessingLevel.
var (
	ProcessingLevel_name = map[int32]string{
		0:  "PROCESSING_LEVEL_UNSPECIFIED",
		12: "PROCESSING_LEVEL_L0",
		10: "PROCESSING_LEVEL_L1",
		1:  "PROCESSING_LEVEL_L1A",
		2:  "PROCESSING_LEVEL_L1B",
		3:  "PROCESSING_LEVEL_L1C",
		4:  "PROCESSING_LEVEL_L2",
		5:  "PROCESSING_LEVEL_L2A",
		6:  "PROCESSING_LEVEL_L2B",
		7:  "PROCESSING_LEVEL_L3",
		8:  "PROCESSING_LEVEL_L3A",
		9:  "PROCESSING_LEVEL_L4",
		11: "PROCESSING_LEVEL_NOT_APPLICABLE",
	}
	ProcessingLevel_value = map[string]int32{
		"PROCESSING_LEVEL_UNSPECIFIED":    0,
		"PROCESSING_LEVEL_L0":             12,
		"PROCESSING_LEVEL_L1":             10,
		"PROCESSING_LEVEL_L1A":            1,
		"PROCESSING_LEVEL_L1B":            2,
		"PROCESSING_LEVEL_L1C":            3,
		"PROCESSING_LEVEL_L2":             4,
		"PROCESSING_LEVEL_L2A":            5,
		"PROCESSING_LEVEL_L2B":            6,
		"PROCESSING_LEVEL_L3":             7,
		"PROCESSING_LEVEL_L3A":            8,
		"PROCESSING_LEVEL_L4":             9,
		"PROCESSING_LEVEL_NOT_APPLICABLE": 11,
	}
)

func (x ProcessingLevel) Enum() *ProcessingLevel {
	p := new(ProcessingLevel)
	*p = x
	return p
}

func (x ProcessingLevel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProcessingLevel) Descriptor() protoreflect.EnumDescriptor {
	return file_datasets_v1_well_known_types_proto_enumTypes[3].Descriptor()
}

func (ProcessingLevel) Type() protoreflect.EnumType {
	return &file_datasets_v1_well_known_types_proto_enumTypes[3]
}

func (x ProcessingLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProcessingLevel.Descriptor instead.
func (ProcessingLevel) EnumDescriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{3}
}

// Polarization of the radar signal.
type Polarization int32

const (
	// Horizontal transmit, Horizontal receive
	Polarization_POLARIZATION_UNSPECIFIED Polarization = 0
	Polarization_POLARIZATION_HH          Polarization = 1
	Polarization_POLARIZATION_HV          Polarization = 2  // Horizontal transmit, Vertical receive
	Polarization_POLARIZATION_VH          Polarization = 3  // Vertical transmit, Horizontal receive
	Polarization_POLARIZATION_VV          Polarization = 4  // Vertical transmit, Vertical receive
	Polarization_POLARIZATION_DUAL_HH     Polarization = 5  // HH+HH
	Polarization_POLARIZATION_DUAL_HV     Polarization = 6  // HV+HV
	Polarization_POLARIZATION_DUAL_VH     Polarization = 7  // VH+VH
	Polarization_POLARIZATION_DUAL_VV     Polarization = 8  // VV+VV
	Polarization_POLARIZATION_HH_HV       Polarization = 9  // HH+HV
	Polarization_POLARIZATION_VV_VH       Polarization = 10 // VV+VH
)

// Enum value maps for Polarization.
var (
	Polarization_name = map[int32]string{
		0:  "POLARIZATION_UNSPECIFIED",
		1:  "POLARIZATION_HH",
		2:  "POLARIZATION_HV",
		3:  "POLARIZATION_VH",
		4:  "POLARIZATION_VV",
		5:  "POLARIZATION_DUAL_HH",
		6:  "POLARIZATION_DUAL_HV",
		7:  "POLARIZATION_DUAL_VH",
		8:  "POLARIZATION_DUAL_VV",
		9:  "POLARIZATION_HH_HV",
		10: "POLARIZATION_VV_VH",
	}
	Polarization_value = map[string]int32{
		"POLARIZATION_UNSPECIFIED": 0,
		"POLARIZATION_HH":          1,
		"POLARIZATION_HV":          2,
		"POLARIZATION_VH":          3,
		"POLARIZATION_VV":          4,
		"POLARIZATION_DUAL_HH":     5,
		"POLARIZATION_DUAL_HV":     6,
		"POLARIZATION_DUAL_VH":     7,
		"POLARIZATION_DUAL_VV":     8,
		"POLARIZATION_HH_HV":       9,
		"POLARIZATION_VV_VH":       10,
	}
)

func (x Polarization) Enum() *Polarization {
	p := new(Polarization)
	*p = x
	return p
}

func (x Polarization) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Polarization) Descriptor() protoreflect.EnumDescriptor {
	return file_datasets_v1_well_known_types_proto_enumTypes[4].Descriptor()
}

func (Polarization) Type() protoreflect.EnumType {
	return &file_datasets_v1_well_known_types_proto_enumTypes[4]
}

func (x Polarization) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Polarization.Descriptor instead.
func (Polarization) EnumDescriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{4}
}

// Sentinel-1 SAR (beam mode):
// https://sentinels.copernicus.eu/web/sentinel/technical-guides/sentinel-1-sar/sar-instrument/acquisition-modes
type AcquisitionMode int32

const (
	AcquisitionMode_ACQUISITION_MODE_UNSPECIFIED AcquisitionMode = 0 // In case it is not set for a dataset
	// used by Sentinel-1 SAR:
	AcquisitionMode_ACQUISITION_MODE_SM AcquisitionMode = 1 // Strip Map
	AcquisitionMode_ACQUISITION_MODE_EW AcquisitionMode = 2 // Extra Wide Swath
	AcquisitionMode_ACQUISITION_MODE_IW AcquisitionMode = 3 // Interferometric Wide Swath
	AcquisitionMode_ACQUISITION_MODE_WV AcquisitionMode = 4 // Wave
	// used by Umbra SAR:
	AcquisitionMode_ACQUISITION_MODE_SPOTLIGHT AcquisitionMode = 10 // Spotlight
	// used by Sentinel 2 MSI:
	AcquisitionMode_ACQUISITION_MODE_NOBS AcquisitionMode = 20 // Nominal Observation
	AcquisitionMode_ACQUISITION_MODE_EOBS AcquisitionMode = 21 // Extended Observation
	AcquisitionMode_ACQUISITION_MODE_DASC AcquisitionMode = 22 // Dark Signal Calibration
	AcquisitionMode_ACQUISITION_MODE_ABSR AcquisitionMode = 23 // Absolute Radiometry Calibration
	AcquisitionMode_ACQUISITION_MODE_VIC  AcquisitionMode = 24 // Vicarious Calibration
	AcquisitionMode_ACQUISITION_MODE_RAW  AcquisitionMode = 25 // Raw Measurement
	AcquisitionMode_ACQUISITION_MODE_TST  AcquisitionMode = 26 // Test Mode
)

// Enum value maps for AcquisitionMode.
var (
	AcquisitionMode_name = map[int32]string{
		0:  "ACQUISITION_MODE_UNSPECIFIED",
		1:  "ACQUISITION_MODE_SM",
		2:  "ACQUISITION_MODE_EW",
		3:  "ACQUISITION_MODE_IW",
		4:  "ACQUISITION_MODE_WV",
		10: "ACQUISITION_MODE_SPOTLIGHT",
		20: "ACQUISITION_MODE_NOBS",
		21: "ACQUISITION_MODE_EOBS",
		22: "ACQUISITION_MODE_DASC",
		23: "ACQUISITION_MODE_ABSR",
		24: "ACQUISITION_MODE_VIC",
		25: "ACQUISITION_MODE_RAW",
		26: "ACQUISITION_MODE_TST",
	}
	AcquisitionMode_value = map[string]int32{
		"ACQUISITION_MODE_UNSPECIFIED": 0,
		"ACQUISITION_MODE_SM":          1,
		"ACQUISITION_MODE_EW":          2,
		"ACQUISITION_MODE_IW":          3,
		"ACQUISITION_MODE_WV":          4,
		"ACQUISITION_MODE_SPOTLIGHT":   10,
		"ACQUISITION_MODE_NOBS":        20,
		"ACQUISITION_MODE_EOBS":        21,
		"ACQUISITION_MODE_DASC":        22,
		"ACQUISITION_MODE_ABSR":        23,
		"ACQUISITION_MODE_VIC":         24,
		"ACQUISITION_MODE_RAW":         25,
		"ACQUISITION_MODE_TST":         26,
	}
)

func (x AcquisitionMode) Enum() *AcquisitionMode {
	p := new(AcquisitionMode)
	*p = x
	return p
}

func (x AcquisitionMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AcquisitionMode) Descriptor() protoreflect.EnumDescriptor {
	return file_datasets_v1_well_known_types_proto_enumTypes[5].Descriptor()
}

func (AcquisitionMode) Type() protoreflect.EnumType {
	return &file_datasets_v1_well_known_types_proto_enumTypes[5]
}

func (x AcquisitionMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AcquisitionMode.Descriptor instead.
func (AcquisitionMode) EnumDescriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{5}
}

// Bytes field (in message)
type UUID struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Uuid          []byte                 `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UUID) Reset() {
	*x = UUID{}
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UUID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UUID) ProtoMessage() {}

func (x *UUID) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UUID.ProtoReflect.Descriptor instead.
func (*UUID) Descriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{0}
}

func (x *UUID) GetUuid() []byte {
	if x != nil {
		return x.Uuid
	}
	return nil
}

// A 3D vector
type Vec3 struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	X             float64                `protobuf:"fixed64,1,opt,name=x,proto3" json:"x,omitempty"`
	Y             float64                `protobuf:"fixed64,2,opt,name=y,proto3" json:"y,omitempty"`
	Z             float64                `protobuf:"fixed64,3,opt,name=z,proto3" json:"z,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Vec3) Reset() {
	*x = Vec3{}
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Vec3) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vec3) ProtoMessage() {}

func (x *Vec3) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vec3.ProtoReflect.Descriptor instead.
func (*Vec3) Descriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{1}
}

func (x *Vec3) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *Vec3) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}

func (x *Vec3) GetZ() float64 {
	if x != nil {
		return x.Z
	}
	return 0
}

// A quaternion
type Quaternion struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Q1            float64                `protobuf:"fixed64,1,opt,name=q1,proto3" json:"q1,omitempty"`
	Q2            float64                `protobuf:"fixed64,2,opt,name=q2,proto3" json:"q2,omitempty"`
	Q3            float64                `protobuf:"fixed64,3,opt,name=q3,proto3" json:"q3,omitempty"`
	Q4            float64                `protobuf:"fixed64,4,opt,name=q4,proto3" json:"q4,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Quaternion) Reset() {
	*x = Quaternion{}
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Quaternion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Quaternion) ProtoMessage() {}

func (x *Quaternion) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Quaternion.ProtoReflect.Descriptor instead.
func (*Quaternion) Descriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{2}
}

func (x *Quaternion) GetQ1() float64 {
	if x != nil {
		return x.Q1
	}
	return 0
}

func (x *Quaternion) GetQ2() float64 {
	if x != nil {
		return x.Q2
	}
	return 0
}

func (x *Quaternion) GetQ3() float64 {
	if x != nil {
		return x.Q3
	}
	return 0
}

func (x *Quaternion) GetQ4() float64 {
	if x != nil {
		return x.Q4
	}
	return 0
}

// LatLon is a pair of latitude and longitude values representing a point on the Earth's surface
type LatLon struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Latitude      float64                `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude     float64                `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LatLon) Reset() {
	*x = LatLon{}
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LatLon) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LatLon) ProtoMessage() {}

func (x *LatLon) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LatLon.ProtoReflect.Descriptor instead.
func (*LatLon) Descriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{3}
}

func (x *LatLon) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *LatLon) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

// LatLonAlt is a pair of latitude, longitude and altitude values representing a point on the Earth's surface
// While we could also use a Vec3 for this, we use a separate message to make the intent of the data clearer.
type LatLonAlt struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Latitude      float64                `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude     float64                `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
	Altitude      float64                `protobuf:"fixed64,3,opt,name=altitude,proto3" json:"altitude,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LatLonAlt) Reset() {
	*x = LatLonAlt{}
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LatLonAlt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LatLonAlt) ProtoMessage() {}

func (x *LatLonAlt) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LatLonAlt.ProtoReflect.Descriptor instead.
func (*LatLonAlt) Descriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{4}
}

func (x *LatLonAlt) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *LatLonAlt) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

func (x *LatLonAlt) GetAltitude() float64 {
	if x != nil {
		return x.Altitude
	}
	return 0
}

// Geometry is of a particular type, e.g. POINT, POLYGON, MULTIPOLYGON, encoded in a binary format
type Geometry struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Wkb           []byte                 `protobuf:"bytes,1,opt,name=wkb,proto3" json:"wkb,omitempty"` // well-known binary representation of a geometry
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Geometry) Reset() {
	*x = Geometry{}
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Geometry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Geometry) ProtoMessage() {}

func (x *Geometry) ProtoReflect() protoreflect.Message {
	mi := &file_datasets_v1_well_known_types_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Geometry.ProtoReflect.Descriptor instead.
func (*Geometry) Descriptor() ([]byte, []int) {
	return file_datasets_v1_well_known_types_proto_rawDescGZIP(), []int{5}
}

func (x *Geometry) GetWkb() []byte {
	if x != nil {
		return x.Wkb
	}
	return nil
}

var File_datasets_v1_well_known_types_proto protoreflect.FileDescriptor

const file_datasets_v1_well_known_types_proto_rawDesc = "" +
	"\n" +
	"\"datasets/v1/well_known_types.proto\x12\vdatasets.v1\x1a\x1egoogle/protobuf/duration.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\x1a\n" +
	"\x04UUID\x12\x12\n" +
	"\x04uuid\x18\x01 \x01(\fR\x04uuid\"0\n" +
	"\x04Vec3\x12\f\n" +
	"\x01x\x18\x01 \x01(\x01R\x01x\x12\f\n" +
	"\x01y\x18\x02 \x01(\x01R\x01y\x12\f\n" +
	"\x01z\x18\x03 \x01(\x01R\x01z\"L\n" +
	"\n" +
	"Quaternion\x12\x0e\n" +
	"\x02q1\x18\x01 \x01(\x01R\x02q1\x12\x0e\n" +
	"\x02q2\x18\x02 \x01(\x01R\x02q2\x12\x0e\n" +
	"\x02q3\x18\x03 \x01(\x01R\x02q3\x12\x0e\n" +
	"\x02q4\x18\x04 \x01(\x01R\x02q4\"B\n" +
	"\x06LatLon\x12\x1a\n" +
	"\blatitude\x18\x01 \x01(\x01R\blatitude\x12\x1c\n" +
	"\tlongitude\x18\x02 \x01(\x01R\tlongitude\"a\n" +
	"\tLatLonAlt\x12\x1a\n" +
	"\blatitude\x18\x01 \x01(\x01R\blatitude\x12\x1c\n" +
	"\tlongitude\x18\x02 \x01(\x01R\tlongitude\x12\x1a\n" +
	"\baltitude\x18\x03 \x01(\x01R\baltitude\"\x1c\n" +
	"\bGeometry\x12\x10\n" +
	"\x03wkb\x18\x01 \x01(\fR\x03wkb*t\n" +
	"\x0fFlightDirection\x12 \n" +
	"\x1cFLIGHT_DIRECTION_UNSPECIFIED\x10\x00\x12\x1e\n" +
	"\x1aFLIGHT_DIRECTION_ASCENDING\x10\x01\x12\x1f\n" +
	"\x1bFLIGHT_DIRECTION_DESCENDING\x10\x02*~\n" +
	"\x14ObservationDirection\x12%\n" +
	"!OBSERVATION_DIRECTION_UNSPECIFIED\x10\x00\x12\x1e\n" +
	"\x1aOBSERVATION_DIRECTION_LEFT\x10\x01\x12\x1f\n" +
	"\x1bOBSERVATION_DIRECTION_RIGHT\x10\x02*\x99\x01\n" +
	"\x10OpendataProvider\x12!\n" +
	"\x1dOPENDATA_PROVIDER_UNSPECIFIED\x10\x00\x12\x19\n" +
	"\x15OPENDATA_PROVIDER_ASF\x10\x01\x12*\n" +
	"&OPENDATA_PROVIDER_COPERNICUS_DATASPACE\x10\x02\x12\x1b\n" +
	"\x17OPENDATA_PROVIDER_UMBRA\x10\x03*\xf1\x02\n" +
	"\x0fProcessingLevel\x12 \n" +
	"\x1cPROCESSING_LEVEL_UNSPECIFIED\x10\x00\x12\x17\n" +
	"\x13PROCESSING_LEVEL_L0\x10\f\x12\x17\n" +
	"\x13PROCESSING_LEVEL_L1\x10\n" +
	"\x12\x18\n" +
	"\x14PROCESSING_LEVEL_L1A\x10\x01\x12\x18\n" +
	"\x14PROCESSING_LEVEL_L1B\x10\x02\x12\x18\n" +
	"\x14PROCESSING_LEVEL_L1C\x10\x03\x12\x17\n" +
	"\x13PROCESSING_LEVEL_L2\x10\x04\x12\x18\n" +
	"\x14PROCESSING_LEVEL_L2A\x10\x05\x12\x18\n" +
	"\x14PROCESSING_LEVEL_L2B\x10\x06\x12\x17\n" +
	"\x13PROCESSING_LEVEL_L3\x10\a\x12\x18\n" +
	"\x14PROCESSING_LEVEL_L3A\x10\b\x12\x17\n" +
	"\x13PROCESSING_LEVEL_L4\x10\t\x12#\n" +
	"\x1fPROCESSING_LEVEL_NOT_APPLICABLE\x10\v*\x98\x02\n" +
	"\fPolarization\x12\x1c\n" +
	"\x18POLARIZATION_UNSPECIFIED\x10\x00\x12\x13\n" +
	"\x0fPOLARIZATION_HH\x10\x01\x12\x13\n" +
	"\x0fPOLARIZATION_HV\x10\x02\x12\x13\n" +
	"\x0fPOLARIZATION_VH\x10\x03\x12\x13\n" +
	"\x0fPOLARIZATION_VV\x10\x04\x12\x18\n" +
	"\x14POLARIZATION_DUAL_HH\x10\x05\x12\x18\n" +
	"\x14POLARIZATION_DUAL_HV\x10\x06\x12\x18\n" +
	"\x14POLARIZATION_DUAL_VH\x10\a\x12\x18\n" +
	"\x14POLARIZATION_DUAL_VV\x10\b\x12\x16\n" +
	"\x12POLARIZATION_HH_HV\x10\t\x12\x16\n" +
	"\x12POLARIZATION_VV_VH\x10\n" +
	"*\xf1\x02\n" +
	"\x0fAcquisitionMode\x12 \n" +
	"\x1cACQUISITION_MODE_UNSPECIFIED\x10\x00\x12\x17\n" +
	"\x13ACQUISITION_MODE_SM\x10\x01\x12\x17\n" +
	"\x13ACQUISITION_MODE_EW\x10\x02\x12\x17\n" +
	"\x13ACQUISITION_MODE_IW\x10\x03\x12\x17\n" +
	"\x13ACQUISITION_MODE_WV\x10\x04\x12\x1e\n" +
	"\x1aACQUISITION_MODE_SPOTLIGHT\x10\n" +
	"\x12\x19\n" +
	"\x15ACQUISITION_MODE_NOBS\x10\x14\x12\x19\n" +
	"\x15ACQUISITION_MODE_EOBS\x10\x15\x12\x19\n" +
	"\x15ACQUISITION_MODE_DASC\x10\x16\x12\x19\n" +
	"\x15ACQUISITION_MODE_ABSR\x10\x17\x12\x18\n" +
	"\x14ACQUISITION_MODE_VIC\x10\x18\x12\x18\n" +
	"\x14ACQUISITION_MODE_RAW\x10\x19\x12\x18\n" +
	"\x14ACQUISITION_MODE_TST\x10\x1aB\xb5\x01\n" +
	"\x0fcom.datasets.v1B\x13WellKnownTypesProtoP\x01Z@github.com/tilebox/tilebox-go/protogen/go/datasets/v1;datasetsv1\xa2\x02\x03DXX\xaa\x02\vDatasets.V1\xca\x02\vDatasets\\V1\xe2\x02\x17Datasets\\V1\\GPBMetadata\xea\x02\fDatasets::V1b\x06proto3"

var (
	file_datasets_v1_well_known_types_proto_rawDescOnce sync.Once
	file_datasets_v1_well_known_types_proto_rawDescData []byte
)

func file_datasets_v1_well_known_types_proto_rawDescGZIP() []byte {
	file_datasets_v1_well_known_types_proto_rawDescOnce.Do(func() {
		file_datasets_v1_well_known_types_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_datasets_v1_well_known_types_proto_rawDesc), len(file_datasets_v1_well_known_types_proto_rawDesc)))
	})
	return file_datasets_v1_well_known_types_proto_rawDescData
}

var file_datasets_v1_well_known_types_proto_enumTypes = make([]protoimpl.EnumInfo, 6)
var file_datasets_v1_well_known_types_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_datasets_v1_well_known_types_proto_goTypes = []any{
	(FlightDirection)(0),      // 0: datasets.v1.FlightDirection
	(ObservationDirection)(0), // 1: datasets.v1.ObservationDirection
	(OpendataProvider)(0),     // 2: datasets.v1.OpendataProvider
	(ProcessingLevel)(0),      // 3: datasets.v1.ProcessingLevel
	(Polarization)(0),         // 4: datasets.v1.Polarization
	(AcquisitionMode)(0),      // 5: datasets.v1.AcquisitionMode
	(*UUID)(nil),              // 6: datasets.v1.UUID
	(*Vec3)(nil),              // 7: datasets.v1.Vec3
	(*Quaternion)(nil),        // 8: datasets.v1.Quaternion
	(*LatLon)(nil),            // 9: datasets.v1.LatLon
	(*LatLonAlt)(nil),         // 10: datasets.v1.LatLonAlt
	(*Geometry)(nil),          // 11: datasets.v1.Geometry
}
var file_datasets_v1_well_known_types_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_datasets_v1_well_known_types_proto_init() }
func file_datasets_v1_well_known_types_proto_init() {
	if File_datasets_v1_well_known_types_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_datasets_v1_well_known_types_proto_rawDesc), len(file_datasets_v1_well_known_types_proto_rawDesc)),
			NumEnums:      6,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_datasets_v1_well_known_types_proto_goTypes,
		DependencyIndexes: file_datasets_v1_well_known_types_proto_depIdxs,
		EnumInfos:         file_datasets_v1_well_known_types_proto_enumTypes,
		MessageInfos:      file_datasets_v1_well_known_types_proto_msgTypes,
	}.Build()
	File_datasets_v1_well_known_types_proto = out.File
	file_datasets_v1_well_known_types_proto_goTypes = nil
	file_datasets_v1_well_known_types_proto_depIdxs = nil
}
