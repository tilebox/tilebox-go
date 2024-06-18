// The externally facing API allowing users to interact with jobs.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: workflows/v1/job.proto

package workflowsv1

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

// SubmitJobRequest submits and schedules a job for execution. The job can have multiple root tasks.
type SubmitJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The root tasks for the job.
	Tasks []*TaskSubmission `protobuf:"bytes,1,rep,name=tasks,proto3" json:"tasks,omitempty"`
	// The name of the job.
	JobName string `protobuf:"bytes,2,opt,name=job_name,json=jobName,proto3" json:"job_name,omitempty"`
	// Tracing information for the job. This is used to propagate tracing information to the workers that execute the job.
	TraceParent string `protobuf:"bytes,3,opt,name=trace_parent,json=traceParent,proto3" json:"trace_parent,omitempty"`
}

func (x *SubmitJobRequest) Reset() {
	*x = SubmitJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_job_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitJobRequest) ProtoMessage() {}

func (x *SubmitJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_job_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitJobRequest.ProtoReflect.Descriptor instead.
func (*SubmitJobRequest) Descriptor() ([]byte, []int) {
	return file_workflows_v1_job_proto_rawDescGZIP(), []int{0}
}

func (x *SubmitJobRequest) GetTasks() []*TaskSubmission {
	if x != nil {
		return x.Tasks
	}
	return nil
}

func (x *SubmitJobRequest) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *SubmitJobRequest) GetTraceParent() string {
	if x != nil {
		return x.TraceParent
	}
	return ""
}

// GetJobRequest requests details for a job.
type GetJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The ID of the job to get details for.
	JobId *UUID `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (x *GetJobRequest) Reset() {
	*x = GetJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_job_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetJobRequest) ProtoMessage() {}

func (x *GetJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_job_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetJobRequest.ProtoReflect.Descriptor instead.
func (*GetJobRequest) Descriptor() ([]byte, []int) {
	return file_workflows_v1_job_proto_rawDescGZIP(), []int{1}
}

func (x *GetJobRequest) GetJobId() *UUID {
	if x != nil {
		return x.JobId
	}
	return nil
}

// RetryJobRequest requests a retry of a job that has failed.
type RetryJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The job to retry.
	JobId *UUID `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (x *RetryJobRequest) Reset() {
	*x = RetryJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_job_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetryJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetryJobRequest) ProtoMessage() {}

func (x *RetryJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_job_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetryJobRequest.ProtoReflect.Descriptor instead.
func (*RetryJobRequest) Descriptor() ([]byte, []int) {
	return file_workflows_v1_job_proto_rawDescGZIP(), []int{2}
}

func (x *RetryJobRequest) GetJobId() *UUID {
	if x != nil {
		return x.JobId
	}
	return nil
}

// RetryJobResponse is the response to a RetryJobRequest.
type RetryJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The number of tasks that were rescheduled.
	NumTasksRescheduled int64 `protobuf:"varint,1,opt,name=num_tasks_rescheduled,json=numTasksRescheduled,proto3" json:"num_tasks_rescheduled,omitempty"`
}

func (x *RetryJobResponse) Reset() {
	*x = RetryJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_job_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RetryJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RetryJobResponse) ProtoMessage() {}

func (x *RetryJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_job_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RetryJobResponse.ProtoReflect.Descriptor instead.
func (*RetryJobResponse) Descriptor() ([]byte, []int) {
	return file_workflows_v1_job_proto_rawDescGZIP(), []int{3}
}

func (x *RetryJobResponse) GetNumTasksRescheduled() int64 {
	if x != nil {
		return x.NumTasksRescheduled
	}
	return 0
}

// CancelJobRequest requests a cancel of a job.
type CancelJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The job to cancel.
	JobId *UUID `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (x *CancelJobRequest) Reset() {
	*x = CancelJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_job_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelJobRequest) ProtoMessage() {}

func (x *CancelJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_job_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelJobRequest.ProtoReflect.Descriptor instead.
func (*CancelJobRequest) Descriptor() ([]byte, []int) {
	return file_workflows_v1_job_proto_rawDescGZIP(), []int{4}
}

func (x *CancelJobRequest) GetJobId() *UUID {
	if x != nil {
		return x.JobId
	}
	return nil
}

// CancelJobResponse is the response to a CancelJobRequest.
type CancelJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CancelJobResponse) Reset() {
	*x = CancelJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_job_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelJobResponse) ProtoMessage() {}

func (x *CancelJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_job_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelJobResponse.ProtoReflect.Descriptor instead.
func (*CancelJobResponse) Descriptor() ([]byte, []int) {
	return file_workflows_v1_job_proto_rawDescGZIP(), []int{5}
}

// VisualizeJobRequest requests a cancel of a job.
type VisualizeJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The job to cancel.
	JobId *UUID `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	// The options for rendering the diagram
	RenderOptions *RenderOptions `protobuf:"bytes,2,opt,name=render_options,json=renderOptions,proto3" json:"render_options,omitempty"`
}

func (x *VisualizeJobRequest) Reset() {
	*x = VisualizeJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_job_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VisualizeJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VisualizeJobRequest) ProtoMessage() {}

func (x *VisualizeJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_job_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VisualizeJobRequest.ProtoReflect.Descriptor instead.
func (*VisualizeJobRequest) Descriptor() ([]byte, []int) {
	return file_workflows_v1_job_proto_rawDescGZIP(), []int{6}
}

func (x *VisualizeJobRequest) GetJobId() *UUID {
	if x != nil {
		return x.JobId
	}
	return nil
}

func (x *VisualizeJobRequest) GetRenderOptions() *RenderOptions {
	if x != nil {
		return x.RenderOptions
	}
	return nil
}

var File_workflows_v1_job_proto protoreflect.FileDescriptor

var file_workflows_v1_job_proto_rawDesc = []byte{
	0x0a, 0x16, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x6a,
	0x6f, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x17, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1a, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x69,
	0x61, 0x67, 0x72, 0x61, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x01, 0x0a, 0x10,
	0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x32, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1c, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54,
	0x61, 0x73, 0x6b, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x74,
	0x61, 0x73, 0x6b, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x6a, 0x6f, 0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x63, 0x65, 0x50, 0x61, 0x72, 0x65,
	0x6e, 0x74, 0x22, 0x3a, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x29, 0x0a, 0x06, 0x6a, 0x6f, 0x62, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x05, 0x6a, 0x6f, 0x62, 0x49, 0x64, 0x22, 0x3c,
	0x0a, 0x0f, 0x52, 0x65, 0x74, 0x72, 0x79, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x29, 0x0a, 0x06, 0x6a, 0x6f, 0x62, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x05, 0x6a, 0x6f, 0x62, 0x49, 0x64, 0x22, 0x46, 0x0a, 0x10,
	0x52, 0x65, 0x74, 0x72, 0x79, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x32, 0x0a, 0x15, 0x6e, 0x75, 0x6d, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x5f, 0x72, 0x65,
	0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x13, 0x6e, 0x75, 0x6d, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x64, 0x22, 0x3d, 0x0a, 0x10, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x4a, 0x6f,
	0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a, 0x06, 0x6a, 0x6f, 0x62, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x05, 0x6a, 0x6f,
	0x62, 0x49, 0x64, 0x22, 0x13, 0x0a, 0x11, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x4a, 0x6f, 0x62,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x84, 0x01, 0x0a, 0x13, 0x56, 0x69, 0x73,
	0x75, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x29, 0x0a, 0x06, 0x6a, 0x6f, 0x62, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x55, 0x55, 0x49, 0x44, 0x52, 0x05, 0x6a, 0x6f, 0x62, 0x49, 0x64, 0x12, 0x42, 0x0a, 0x0e, 0x72,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x52, 0x0d, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x32,
	0xe9, 0x02, 0x0a, 0x0a, 0x4a, 0x6f, 0x62, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3e,
	0x0a, 0x09, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4a, 0x6f, 0x62, 0x12, 0x1e, 0x2e, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69,
	0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x12, 0x38,
	0x0a, 0x06, 0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x12, 0x1b, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x12, 0x49, 0x0a, 0x08, 0x52, 0x65, 0x74, 0x72,
	0x79, 0x4a, 0x6f, 0x62, 0x12, 0x1d, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x79, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x52, 0x65, 0x74, 0x72, 0x79, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x09, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x4a, 0x6f, 0x62,
	0x12, 0x1e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1f, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x48, 0x0a, 0x0c, 0x56, 0x69, 0x73, 0x75, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x4a, 0x6f,
	0x62, 0x12, 0x21, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x56, 0x69, 0x73, 0x75, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69, 0x61, 0x67, 0x72, 0x61, 0x6d, 0x42, 0xb1, 0x01, 0x0a, 0x10,
	0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31,
	0x42, 0x08, 0x4a, 0x6f, 0x62, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x42, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f, 0x78,
	0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f, 0x78, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x2f, 0x76, 0x31, 0x3b, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x76, 0x31,
	0xa2, 0x02, 0x03, 0x57, 0x58, 0x58, 0xaa, 0x02, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f,
	0x77, 0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x18, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_workflows_v1_job_proto_rawDescOnce sync.Once
	file_workflows_v1_job_proto_rawDescData = file_workflows_v1_job_proto_rawDesc
)

func file_workflows_v1_job_proto_rawDescGZIP() []byte {
	file_workflows_v1_job_proto_rawDescOnce.Do(func() {
		file_workflows_v1_job_proto_rawDescData = protoimpl.X.CompressGZIP(file_workflows_v1_job_proto_rawDescData)
	})
	return file_workflows_v1_job_proto_rawDescData
}

var file_workflows_v1_job_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_workflows_v1_job_proto_goTypes = []any{
	(*SubmitJobRequest)(nil),    // 0: workflows.v1.SubmitJobRequest
	(*GetJobRequest)(nil),       // 1: workflows.v1.GetJobRequest
	(*RetryJobRequest)(nil),     // 2: workflows.v1.RetryJobRequest
	(*RetryJobResponse)(nil),    // 3: workflows.v1.RetryJobResponse
	(*CancelJobRequest)(nil),    // 4: workflows.v1.CancelJobRequest
	(*CancelJobResponse)(nil),   // 5: workflows.v1.CancelJobResponse
	(*VisualizeJobRequest)(nil), // 6: workflows.v1.VisualizeJobRequest
	(*TaskSubmission)(nil),      // 7: workflows.v1.TaskSubmission
	(*UUID)(nil),                // 8: workflows.v1.UUID
	(*RenderOptions)(nil),       // 9: workflows.v1.RenderOptions
	(*Job)(nil),                 // 10: workflows.v1.Job
	(*Diagram)(nil),             // 11: workflows.v1.Diagram
}
var file_workflows_v1_job_proto_depIdxs = []int32{
	7,  // 0: workflows.v1.SubmitJobRequest.tasks:type_name -> workflows.v1.TaskSubmission
	8,  // 1: workflows.v1.GetJobRequest.job_id:type_name -> workflows.v1.UUID
	8,  // 2: workflows.v1.RetryJobRequest.job_id:type_name -> workflows.v1.UUID
	8,  // 3: workflows.v1.CancelJobRequest.job_id:type_name -> workflows.v1.UUID
	8,  // 4: workflows.v1.VisualizeJobRequest.job_id:type_name -> workflows.v1.UUID
	9,  // 5: workflows.v1.VisualizeJobRequest.render_options:type_name -> workflows.v1.RenderOptions
	0,  // 6: workflows.v1.JobService.SubmitJob:input_type -> workflows.v1.SubmitJobRequest
	1,  // 7: workflows.v1.JobService.GetJob:input_type -> workflows.v1.GetJobRequest
	2,  // 8: workflows.v1.JobService.RetryJob:input_type -> workflows.v1.RetryJobRequest
	4,  // 9: workflows.v1.JobService.CancelJob:input_type -> workflows.v1.CancelJobRequest
	6,  // 10: workflows.v1.JobService.VisualizeJob:input_type -> workflows.v1.VisualizeJobRequest
	10, // 11: workflows.v1.JobService.SubmitJob:output_type -> workflows.v1.Job
	10, // 12: workflows.v1.JobService.GetJob:output_type -> workflows.v1.Job
	3,  // 13: workflows.v1.JobService.RetryJob:output_type -> workflows.v1.RetryJobResponse
	5,  // 14: workflows.v1.JobService.CancelJob:output_type -> workflows.v1.CancelJobResponse
	11, // 15: workflows.v1.JobService.VisualizeJob:output_type -> workflows.v1.Diagram
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_workflows_v1_job_proto_init() }
func file_workflows_v1_job_proto_init() {
	if File_workflows_v1_job_proto != nil {
		return
	}
	file_workflows_v1_core_proto_init()
	file_workflows_v1_diagram_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_workflows_v1_job_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*SubmitJobRequest); i {
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
		file_workflows_v1_job_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*GetJobRequest); i {
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
		file_workflows_v1_job_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*RetryJobRequest); i {
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
		file_workflows_v1_job_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*RetryJobResponse); i {
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
		file_workflows_v1_job_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*CancelJobRequest); i {
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
		file_workflows_v1_job_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*CancelJobResponse); i {
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
		file_workflows_v1_job_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*VisualizeJobRequest); i {
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
			RawDescriptor: file_workflows_v1_job_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_workflows_v1_job_proto_goTypes,
		DependencyIndexes: file_workflows_v1_job_proto_depIdxs,
		MessageInfos:      file_workflows_v1_job_proto_msgTypes,
	}.Build()
	File_workflows_v1_job_proto = out.File
	file_workflows_v1_job_proto_rawDesc = nil
	file_workflows_v1_job_proto_goTypes = nil
	file_workflows_v1_job_proto_depIdxs = nil
}
