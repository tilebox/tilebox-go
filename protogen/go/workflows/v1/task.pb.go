// The internal API a task runner uses to communicate with a workflows-service.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: workflows/v1/task.proto

package workflowsv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// NextTaskRequest is the request for requesting the next task to run and marking a task as computed.
type NextTaskRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ComputedTask *ComputedTask `protobuf:"bytes,1,opt,name=computed_task,json=computedTask,proto3,oneof" json:"computed_task,omitempty"` // The task that has been computed. If not set, the next task will
	// The capabilities of the task runner, and therefore the potential tasks that can be run by that task runner.
	NextTaskToRun *NextTaskToRun `protobuf:"bytes,2,opt,name=next_task_to_run,json=nextTaskToRun,proto3,oneof" json:"next_task_to_run,omitempty"`
}

func (x *NextTaskRequest) Reset() {
	*x = NextTaskRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_task_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NextTaskRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NextTaskRequest) ProtoMessage() {}

func (x *NextTaskRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_task_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NextTaskRequest.ProtoReflect.Descriptor instead.
func (*NextTaskRequest) Descriptor() ([]byte, []int) {
	return file_workflows_v1_task_proto_rawDescGZIP(), []int{0}
}

func (x *NextTaskRequest) GetComputedTask() *ComputedTask {
	if x != nil {
		return x.ComputedTask
	}
	return nil
}

func (x *NextTaskRequest) GetNextTaskToRun() *NextTaskToRun {
	if x != nil {
		return x.NextTaskToRun
	}
	return nil
}

// NextTaskToRun is a message specifying the capabilities of the task runner, and therefore the potential
// tasks that can be run by that task runner.
type NextTaskToRun struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClusterSlug string            `protobuf:"bytes,1,opt,name=cluster_slug,json=clusterSlug,proto3" json:"cluster_slug,omitempty"` // The cluster that this task runner is running on.
	Identifiers []*TaskIdentifier `protobuf:"bytes,2,rep,name=identifiers,proto3" json:"identifiers,omitempty"`                    // The task identifiers that this task runner can run.
}

func (x *NextTaskToRun) Reset() {
	*x = NextTaskToRun{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_task_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NextTaskToRun) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NextTaskToRun) ProtoMessage() {}

func (x *NextTaskToRun) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_task_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NextTaskToRun.ProtoReflect.Descriptor instead.
func (*NextTaskToRun) Descriptor() ([]byte, []int) {
	return file_workflows_v1_task_proto_rawDescGZIP(), []int{1}
}

func (x *NextTaskToRun) GetClusterSlug() string {
	if x != nil {
		return x.ClusterSlug
	}
	return ""
}

func (x *NextTaskToRun) GetIdentifiers() []*TaskIdentifier {
	if x != nil {
		return x.Identifiers
	}
	return nil
}

// ComputedTask is a message specifying a task that has been computed by the task runner.
type ComputedTask struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id *UUID `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // The id of the task that has been computed.
	// A display name for the task that has been computed for visualization purposes.
	// If not set, the display message specified upon task submission will be kept.
	Display  string            `protobuf:"bytes,2,opt,name=display,proto3" json:"display,omitempty"`
	SubTasks []*TaskSubmission `protobuf:"bytes,3,rep,name=sub_tasks,json=subTasks,proto3" json:"sub_tasks,omitempty"` // A list of sub-tasks that the just computed task spawned.
}

func (x *ComputedTask) Reset() {
	*x = ComputedTask{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_task_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ComputedTask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComputedTask) ProtoMessage() {}

func (x *ComputedTask) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_task_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComputedTask.ProtoReflect.Descriptor instead.
func (*ComputedTask) Descriptor() ([]byte, []int) {
	return file_workflows_v1_task_proto_rawDescGZIP(), []int{2}
}

func (x *ComputedTask) GetId() *UUID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *ComputedTask) GetDisplay() string {
	if x != nil {
		return x.Display
	}
	return ""
}

func (x *ComputedTask) GetSubTasks() []*TaskSubmission {
	if x != nil {
		return x.SubTasks
	}
	return nil
}

// NextTaskResponse is the response to the NextTask request.
// Right now it only contains the next task to run. Wrapped in a message to allow adding more fields later if needed.
type NextTaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NextTask *Task `protobuf:"bytes,1,opt,name=next_task,json=nextTask,proto3" json:"next_task,omitempty"`
}

func (x *NextTaskResponse) Reset() {
	*x = NextTaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_task_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NextTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NextTaskResponse) ProtoMessage() {}

func (x *NextTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_task_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NextTaskResponse.ProtoReflect.Descriptor instead.
func (*NextTaskResponse) Descriptor() ([]byte, []int) {
	return file_workflows_v1_task_proto_rawDescGZIP(), []int{3}
}

func (x *NextTaskResponse) GetNextTask() *Task {
	if x != nil {
		return x.NextTask
	}
	return nil
}

// TaskFailedRequest is the request for marking a task as failed.
type TaskFailedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId    *UUID  `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Display   string `protobuf:"bytes,2,opt,name=display,proto3" json:"display,omitempty"`
	CancelJob bool   `protobuf:"varint,3,opt,name=cancel_job,json=cancelJob,proto3" json:"cancel_job,omitempty"`
}

func (x *TaskFailedRequest) Reset() {
	*x = TaskFailedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_task_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskFailedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskFailedRequest) ProtoMessage() {}

func (x *TaskFailedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_task_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskFailedRequest.ProtoReflect.Descriptor instead.
func (*TaskFailedRequest) Descriptor() ([]byte, []int) {
	return file_workflows_v1_task_proto_rawDescGZIP(), []int{4}
}

func (x *TaskFailedRequest) GetTaskId() *UUID {
	if x != nil {
		return x.TaskId
	}
	return nil
}

func (x *TaskFailedRequest) GetDisplay() string {
	if x != nil {
		return x.Display
	}
	return ""
}

func (x *TaskFailedRequest) GetCancelJob() bool {
	if x != nil {
		return x.CancelJob
	}
	return false
}

// TaskStateResponse is the response to the TaskFailed request,
// indicating the current state of the task marked as failed.
type TaskStateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State TaskState `protobuf:"varint,1,opt,name=state,proto3,enum=workflows.v1.TaskState" json:"state,omitempty"`
}

func (x *TaskStateResponse) Reset() {
	*x = TaskStateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_task_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskStateResponse) ProtoMessage() {}

func (x *TaskStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_task_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskStateResponse.ProtoReflect.Descriptor instead.
func (*TaskStateResponse) Descriptor() ([]byte, []int) {
	return file_workflows_v1_task_proto_rawDescGZIP(), []int{5}
}

func (x *TaskStateResponse) GetState() TaskState {
	if x != nil {
		return x.State
	}
	return TaskState_TASK_STATE_UNSPECIFIED
}

// TaskLease is a message specifying the new lease expiration time of a task.
type TaskLeaseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId         *UUID                `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	RequestedLease *durationpb.Duration `protobuf:"bytes,2,opt,name=requested_lease,json=requestedLease,proto3,oneof" json:"requested_lease,omitempty"`
}

func (x *TaskLeaseRequest) Reset() {
	*x = TaskLeaseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_task_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskLeaseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskLeaseRequest) ProtoMessage() {}

func (x *TaskLeaseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_task_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskLeaseRequest.ProtoReflect.Descriptor instead.
func (*TaskLeaseRequest) Descriptor() ([]byte, []int) {
	return file_workflows_v1_task_proto_rawDescGZIP(), []int{6}
}

func (x *TaskLeaseRequest) GetTaskId() *UUID {
	if x != nil {
		return x.TaskId
	}
	return nil
}

func (x *TaskLeaseRequest) GetRequestedLease() *durationpb.Duration {
	if x != nil {
		return x.RequestedLease
	}
	return nil
}

var File_workflows_v1_task_proto protoreflect.FileDescriptor

var file_workflows_v1_task_proto_rawDesc = []byte{
	0x0a, 0x17, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x74,
	0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f,
	0x77, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xc9, 0x01, 0x0a, 0x0f, 0x4e, 0x65, 0x78, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x44, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x64,
	0x5f, 0x74, 0x61, 0x73, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x75,
	0x74, 0x65, 0x64, 0x54, 0x61, 0x73, 0x6b, 0x48, 0x00, 0x52, 0x0c, 0x63, 0x6f, 0x6d, 0x70, 0x75,
	0x74, 0x65, 0x64, 0x54, 0x61, 0x73, 0x6b, 0x88, 0x01, 0x01, 0x12, 0x49, 0x0a, 0x10, 0x6e, 0x65,
	0x78, 0x74, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x74, 0x6f, 0x5f, 0x72, 0x75, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x78, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x54, 0x6f, 0x52, 0x75,
	0x6e, 0x48, 0x01, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x54, 0x6f, 0x52,
	0x75, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x75, 0x74,
	0x65, 0x64, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x6e, 0x65, 0x78, 0x74,
	0x5f, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x74, 0x6f, 0x5f, 0x72, 0x75, 0x6e, 0x22, 0x72, 0x0a, 0x0d,
	0x4e, 0x65, 0x78, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x54, 0x6f, 0x52, 0x75, 0x6e, 0x12, 0x21, 0x0a,
	0x0c, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x6c, 0x75, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x53, 0x6c, 0x75, 0x67,
	0x12, 0x3e, 0x0a, 0x0b, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x52, 0x0b, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x73,
	0x22, 0x87, 0x01, 0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x70, 0x75, 0x74, 0x65, 0x64, 0x54, 0x61, 0x73,
	0x6b, 0x12, 0x22, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x55, 0x49,
	0x44, 0x52, 0x02, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x12,
	0x39, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x08, 0x73, 0x75, 0x62, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x22, 0x43, 0x0a, 0x10, 0x4e, 0x65,
	0x78, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f,
	0x0a, 0x09, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x08, 0x6e, 0x65, 0x78, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x22,
	0x79, 0x0a, 0x11, 0x54, 0x61, 0x73, 0x6b, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x63,
	0x61, 0x6e, 0x63, 0x65, 0x6c, 0x5f, 0x6a, 0x6f, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x09, 0x63, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x4a, 0x6f, 0x62, 0x22, 0x42, 0x0a, 0x11, 0x54, 0x61,
	0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2d, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17,
	0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61,
	0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x9c,
	0x01, 0x0a, 0x10, 0x54, 0x61, 0x73, 0x6b, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64,
	0x12, 0x47, 0x0a, 0x0f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x5f, 0x6c, 0x65,
	0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x65,
	0x64, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x88, 0x01, 0x01, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x65, 0x64, 0x5f, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x32, 0xf4, 0x01,
	0x0a, 0x0b, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x49, 0x0a,
	0x08, 0x4e, 0x65, 0x78, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x1d, 0x2e, 0x77, 0x6f, 0x72, 0x6b,
	0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x78, 0x74, 0x54, 0x61, 0x73,
	0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x65, 0x78, 0x74, 0x54, 0x61, 0x73, 0x6b,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4e, 0x0a, 0x0a, 0x54, 0x61, 0x73, 0x6b,
	0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x12, 0x1f, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f,
	0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0f, 0x45, 0x78, 0x74, 0x65,
	0x6e, 0x64, 0x54, 0x61, 0x73, 0x6b, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x1e, 0x2e, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x4c,
	0x65, 0x61, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x4c,
	0x65, 0x61, 0x73, 0x65, 0x42, 0xb2, 0x01, 0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x6f, 0x72,
	0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x54, 0x61, 0x73, 0x6b, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f, 0x78, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62,
	0x6f, 0x78, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x67,
	0x6f, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x77,
	0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x57, 0x58, 0x58,
	0xaa, 0x02, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x18, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x57, 0x6f, 0x72, 0x6b,
	0x66, 0x6c, 0x6f, 0x77, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_workflows_v1_task_proto_rawDescOnce sync.Once
	file_workflows_v1_task_proto_rawDescData = file_workflows_v1_task_proto_rawDesc
)

func file_workflows_v1_task_proto_rawDescGZIP() []byte {
	file_workflows_v1_task_proto_rawDescOnce.Do(func() {
		file_workflows_v1_task_proto_rawDescData = protoimpl.X.CompressGZIP(file_workflows_v1_task_proto_rawDescData)
	})
	return file_workflows_v1_task_proto_rawDescData
}

var file_workflows_v1_task_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_workflows_v1_task_proto_goTypes = []interface{}{
	(*NextTaskRequest)(nil),     // 0: workflows.v1.NextTaskRequest
	(*NextTaskToRun)(nil),       // 1: workflows.v1.NextTaskToRun
	(*ComputedTask)(nil),        // 2: workflows.v1.ComputedTask
	(*NextTaskResponse)(nil),    // 3: workflows.v1.NextTaskResponse
	(*TaskFailedRequest)(nil),   // 4: workflows.v1.TaskFailedRequest
	(*TaskStateResponse)(nil),   // 5: workflows.v1.TaskStateResponse
	(*TaskLeaseRequest)(nil),    // 6: workflows.v1.TaskLeaseRequest
	(*TaskIdentifier)(nil),      // 7: workflows.v1.TaskIdentifier
	(*UUID)(nil),                // 8: workflows.v1.UUID
	(*TaskSubmission)(nil),      // 9: workflows.v1.TaskSubmission
	(*Task)(nil),                // 10: workflows.v1.Task
	(TaskState)(0),              // 11: workflows.v1.TaskState
	(*durationpb.Duration)(nil), // 12: google.protobuf.Duration
	(*TaskLease)(nil),           // 13: workflows.v1.TaskLease
}
var file_workflows_v1_task_proto_depIdxs = []int32{
	2,  // 0: workflows.v1.NextTaskRequest.computed_task:type_name -> workflows.v1.ComputedTask
	1,  // 1: workflows.v1.NextTaskRequest.next_task_to_run:type_name -> workflows.v1.NextTaskToRun
	7,  // 2: workflows.v1.NextTaskToRun.identifiers:type_name -> workflows.v1.TaskIdentifier
	8,  // 3: workflows.v1.ComputedTask.id:type_name -> workflows.v1.UUID
	9,  // 4: workflows.v1.ComputedTask.sub_tasks:type_name -> workflows.v1.TaskSubmission
	10, // 5: workflows.v1.NextTaskResponse.next_task:type_name -> workflows.v1.Task
	8,  // 6: workflows.v1.TaskFailedRequest.task_id:type_name -> workflows.v1.UUID
	11, // 7: workflows.v1.TaskStateResponse.state:type_name -> workflows.v1.TaskState
	8,  // 8: workflows.v1.TaskLeaseRequest.task_id:type_name -> workflows.v1.UUID
	12, // 9: workflows.v1.TaskLeaseRequest.requested_lease:type_name -> google.protobuf.Duration
	0,  // 10: workflows.v1.TaskService.NextTask:input_type -> workflows.v1.NextTaskRequest
	4,  // 11: workflows.v1.TaskService.TaskFailed:input_type -> workflows.v1.TaskFailedRequest
	6,  // 12: workflows.v1.TaskService.ExtendTaskLease:input_type -> workflows.v1.TaskLeaseRequest
	3,  // 13: workflows.v1.TaskService.NextTask:output_type -> workflows.v1.NextTaskResponse
	5,  // 14: workflows.v1.TaskService.TaskFailed:output_type -> workflows.v1.TaskStateResponse
	13, // 15: workflows.v1.TaskService.ExtendTaskLease:output_type -> workflows.v1.TaskLease
	13, // [13:16] is the sub-list for method output_type
	10, // [10:13] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_workflows_v1_task_proto_init() }
func file_workflows_v1_task_proto_init() {
	if File_workflows_v1_task_proto != nil {
		return
	}
	file_workflows_v1_core_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_workflows_v1_task_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NextTaskRequest); i {
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
		file_workflows_v1_task_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NextTaskToRun); i {
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
		file_workflows_v1_task_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ComputedTask); i {
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
		file_workflows_v1_task_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NextTaskResponse); i {
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
		file_workflows_v1_task_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskFailedRequest); i {
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
		file_workflows_v1_task_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskStateResponse); i {
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
		file_workflows_v1_task_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskLeaseRequest); i {
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
	file_workflows_v1_task_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_workflows_v1_task_proto_msgTypes[6].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_workflows_v1_task_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_workflows_v1_task_proto_goTypes,
		DependencyIndexes: file_workflows_v1_task_proto_depIdxs,
		MessageInfos:      file_workflows_v1_task_proto_msgTypes,
	}.Build()
	File_workflows_v1_task_proto = out.File
	file_workflows_v1_task_proto_rawDesc = nil
	file_workflows_v1_task_proto_goTypes = nil
	file_workflows_v1_task_proto_depIdxs = nil
}
