// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: workflows/v1/core.proto

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

// The state of a task.
type TaskState int32

const (
	TaskState_UNDEFINED TaskState = 0
	TaskState_QUEUED    TaskState = 1 // The task is queued and waiting to be run.
	TaskState_RUNNING   TaskState = 2 // The task is currently running on some task runner.
	TaskState_COMPUTED  TaskState = 3 // The task has been computed and the output is available. If the task also has no more outstanding children, it is considered COMPLETED.
	TaskState_FAILED    TaskState = 4 // The task has failed.
	TaskState_CANCELLED TaskState = 5 // The task has been cancelled due to user request.
)

// Enum value maps for TaskState.
var (
	TaskState_name = map[int32]string{
		0: "UNDEFINED",
		1: "QUEUED",
		2: "RUNNING",
		3: "COMPUTED",
		4: "FAILED",
		5: "CANCELLED",
	}
	TaskState_value = map[string]int32{
		"UNDEFINED": 0,
		"QUEUED":    1,
		"RUNNING":   2,
		"COMPUTED":  3,
		"FAILED":    4,
		"CANCELLED": 5,
	}
)

func (x TaskState) Enum() *TaskState {
	p := new(TaskState)
	*p = x
	return p
}

func (x TaskState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TaskState) Descriptor() protoreflect.EnumDescriptor {
	return file_workflows_v1_core_proto_enumTypes[0].Descriptor()
}

func (TaskState) Type() protoreflect.EnumType {
	return &file_workflows_v1_core_proto_enumTypes[0]
}

func (x TaskState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TaskState.Descriptor instead.
func (TaskState) EnumDescriptor() ([]byte, []int) {
	return file_workflows_v1_core_proto_rawDescGZIP(), []int{0}
}

// A cluster is a grouping of tasks that are related.
type Cluster struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 1 is reserved for a potential id field in the future.
	Slug        string `protobuf:"bytes,2,opt,name=slug,proto3" json:"slug,omitempty"`                                  // The unique slug of the cluster within the namespace.
	DisplayName string `protobuf:"bytes,3,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"` // The display name of the cluster.
}

func (x *Cluster) Reset() {
	*x = Cluster{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_core_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cluster) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cluster) ProtoMessage() {}

func (x *Cluster) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_core_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cluster.ProtoReflect.Descriptor instead.
func (*Cluster) Descriptor() ([]byte, []int) {
	return file_workflows_v1_core_proto_rawDescGZIP(), []int{0}
}

func (x *Cluster) GetSlug() string {
	if x != nil {
		return x.Slug
	}
	return ""
}

func (x *Cluster) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

// A job is a logical grouping of tasks that are related.
type Job struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          *UUID  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	TraceParent string `protobuf:"bytes,3,opt,name=trace_parent,json=traceParent,proto3" json:"trace_parent,omitempty"`
	Completed   bool   `protobuf:"varint,4,opt,name=completed,proto3" json:"completed,omitempty"`
}

func (x *Job) Reset() {
	*x = Job{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_core_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Job) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Job) ProtoMessage() {}

func (x *Job) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_core_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Job.ProtoReflect.Descriptor instead.
func (*Job) Descriptor() ([]byte, []int) {
	return file_workflows_v1_core_proto_rawDescGZIP(), []int{1}
}

func (x *Job) GetId() *UUID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Job) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Job) GetTraceParent() string {
	if x != nil {
		return x.TraceParent
	}
	return ""
}

func (x *Job) GetCompleted() bool {
	if x != nil {
		return x.Completed
	}
	return false
}

// A task is a single unit of work.
type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         *UUID           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`                                     // The id of the task instance. Contains the submission timestamp as the time part of the ULID.
	Identifier *TaskIdentifier `protobuf:"bytes,2,opt,name=identifier,proto3" json:"identifier,omitempty"`                     // Unique identifier for the task. Used by runners to match tasks to specific functions.
	State      TaskState       `protobuf:"varint,3,opt,name=state,proto3,enum=workflows.v1.TaskState" json:"state,omitempty"`  // The current state of the task.
	Input      []byte          `protobuf:"bytes,4,opt,name=input,proto3,oneof" json:"input,omitempty"`                         // The serialized input parameters for the task in the format that this task expects.
	Display    *string         `protobuf:"bytes,5,opt,name=display,proto3,oneof" json:"display,omitempty"`                     // Display is a human readable representation of the Task used for printing or visualizations
	Job        *Job            `protobuf:"bytes,6,opt,name=job,proto3" json:"job,omitempty"`                                   // The job that this task belongs to.
	ParentId   *UUID           `protobuf:"bytes,7,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"`         // The id of the parent task.
	DependsOn  []*UUID         `protobuf:"bytes,8,rep,name=depends_on,json=dependsOn,proto3" json:"depends_on,omitempty"`      // The ids of the tasks that this task depends on.
	Lease      *TaskLease      `protobuf:"bytes,9,opt,name=lease,proto3" json:"lease,omitempty"`                               // The lease of the task.
	RetryCount int64           `protobuf:"varint,10,opt,name=retry_count,json=retryCount,proto3" json:"retry_count,omitempty"` // The number of times this task has been retried.
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_core_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_core_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_workflows_v1_core_proto_rawDescGZIP(), []int{2}
}

func (x *Task) GetId() *UUID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Task) GetIdentifier() *TaskIdentifier {
	if x != nil {
		return x.Identifier
	}
	return nil
}

func (x *Task) GetState() TaskState {
	if x != nil {
		return x.State
	}
	return TaskState_UNDEFINED
}

func (x *Task) GetInput() []byte {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *Task) GetDisplay() string {
	if x != nil && x.Display != nil {
		return *x.Display
	}
	return ""
}

func (x *Task) GetJob() *Job {
	if x != nil {
		return x.Job
	}
	return nil
}

func (x *Task) GetParentId() *UUID {
	if x != nil {
		return x.ParentId
	}
	return nil
}

func (x *Task) GetDependsOn() []*UUID {
	if x != nil {
		return x.DependsOn
	}
	return nil
}

func (x *Task) GetLease() *TaskLease {
	if x != nil {
		return x.Lease
	}
	return nil
}

func (x *Task) GetRetryCount() int64 {
	if x != nil {
		return x.RetryCount
	}
	return 0
}

// An identifier for a task.
type TaskIdentifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`       // A unique name of a task (unique within a namespace).
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"` // Version of the task.
}

func (x *TaskIdentifier) Reset() {
	*x = TaskIdentifier{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_core_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskIdentifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskIdentifier) ProtoMessage() {}

func (x *TaskIdentifier) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_core_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskIdentifier.ProtoReflect.Descriptor instead.
func (*TaskIdentifier) Descriptor() ([]byte, []int) {
	return file_workflows_v1_core_proto_rawDescGZIP(), []int{3}
}

func (x *TaskIdentifier) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TaskIdentifier) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

// A list of tasks.
type Tasks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tasks []*Task `protobuf:"bytes,1,rep,name=tasks,proto3" json:"tasks,omitempty"`
}

func (x *Tasks) Reset() {
	*x = Tasks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_core_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tasks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tasks) ProtoMessage() {}

func (x *Tasks) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_core_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tasks.ProtoReflect.Descriptor instead.
func (*Tasks) Descriptor() ([]byte, []int) {
	return file_workflows_v1_core_proto_rawDescGZIP(), []int{4}
}

func (x *Tasks) GetTasks() []*Task {
	if x != nil {
		return x.Tasks
	}
	return nil
}

// TaskSubmission is a message of a task that is just about to be submitted, either by submitting a job or as a subtask.
type TaskSubmission struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClusterSlug  string          `protobuf:"bytes,1,opt,name=cluster_slug,json=clusterSlug,proto3" json:"cluster_slug,omitempty"` // The cluster that this task should be run on
	Identifier   *TaskIdentifier `protobuf:"bytes,2,opt,name=identifier,proto3" json:"identifier,omitempty"`                      // The task identifier
	Input        []byte          `protobuf:"bytes,3,opt,name=input,proto3" json:"input,omitempty"`                                // The serialized task instance
	Display      string          `protobuf:"bytes,4,opt,name=display,proto3" json:"display,omitempty"`                            // A human-readable description of the task
	Dependencies []int64         `protobuf:"varint,5,rep,packed,name=dependencies,proto3" json:"dependencies,omitempty"`          // A list of indices, corresponding to tasks in the list of sub_tasks that this SubTask is part of.
	MaxRetries   int64           `protobuf:"varint,6,opt,name=max_retries,json=maxRetries,proto3" json:"max_retries,omitempty"`   // The maximum number of retries for this task.
}

func (x *TaskSubmission) Reset() {
	*x = TaskSubmission{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_core_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskSubmission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskSubmission) ProtoMessage() {}

func (x *TaskSubmission) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_core_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskSubmission.ProtoReflect.Descriptor instead.
func (*TaskSubmission) Descriptor() ([]byte, []int) {
	return file_workflows_v1_core_proto_rawDescGZIP(), []int{5}
}

func (x *TaskSubmission) GetClusterSlug() string {
	if x != nil {
		return x.ClusterSlug
	}
	return ""
}

func (x *TaskSubmission) GetIdentifier() *TaskIdentifier {
	if x != nil {
		return x.Identifier
	}
	return nil
}

func (x *TaskSubmission) GetInput() []byte {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *TaskSubmission) GetDisplay() string {
	if x != nil {
		return x.Display
	}
	return ""
}

func (x *TaskSubmission) GetDependencies() []int64 {
	if x != nil {
		return x.Dependencies
	}
	return nil
}

func (x *TaskSubmission) GetMaxRetries() int64 {
	if x != nil {
		return x.MaxRetries
	}
	return 0
}

// Bytes field (in message)
type UUID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid []byte `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
}

func (x *UUID) Reset() {
	*x = UUID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_core_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UUID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UUID) ProtoMessage() {}

func (x *UUID) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_core_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
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
	return file_workflows_v1_core_proto_rawDescGZIP(), []int{6}
}

func (x *UUID) GetUuid() []byte {
	if x != nil {
		return x.Uuid
	}
	return nil
}

type TaskLease struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lease                             *durationpb.Duration `protobuf:"bytes,1,opt,name=lease,proto3" json:"lease,omitempty"`
	RecommendedWaitUntilNextExtension *durationpb.Duration `protobuf:"bytes,2,opt,name=recommended_wait_until_next_extension,json=recommendedWaitUntilNextExtension,proto3" json:"recommended_wait_until_next_extension,omitempty"`
}

func (x *TaskLease) Reset() {
	*x = TaskLease{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflows_v1_core_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskLease) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskLease) ProtoMessage() {}

func (x *TaskLease) ProtoReflect() protoreflect.Message {
	mi := &file_workflows_v1_core_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskLease.ProtoReflect.Descriptor instead.
func (*TaskLease) Descriptor() ([]byte, []int) {
	return file_workflows_v1_core_proto_rawDescGZIP(), []int{7}
}

func (x *TaskLease) GetLease() *durationpb.Duration {
	if x != nil {
		return x.Lease
	}
	return nil
}

func (x *TaskLease) GetRecommendedWaitUntilNextExtension() *durationpb.Duration {
	if x != nil {
		return x.RecommendedWaitUntilNextExtension
	}
	return nil
}

var File_workflows_v1_core_proto protoreflect.FileDescriptor

var file_workflows_v1_core_proto_rawDesc = []byte{
	0x0a, 0x17, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x40, 0x0a, 0x07, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x73, 0x6c, 0x75, 0x67, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61,
	0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x69,
	0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x7e, 0x0a, 0x03, 0x4a, 0x6f, 0x62,
	0x12, 0x22, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77,
	0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x55, 0x49, 0x44,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x63,
	0x65, 0x5f, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x74, 0x72, 0x61, 0x63, 0x65, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x63,
	0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09,
	0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x22, 0xc0, 0x03, 0x0a, 0x04, 0x54, 0x61,
	0x73, 0x6b, 0x12, 0x22, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x55,
	0x49, 0x44, 0x52, 0x02, 0x69, 0x64, 0x12, 0x3c, 0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x77, 0x6f, 0x72,
	0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x12, 0x19, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0c, 0x48, 0x00, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x88, 0x01, 0x01, 0x12, 0x1d,
	0x0a, 0x07, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x01, 0x52, 0x07, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a,
	0x03, 0x6a, 0x6f, 0x62, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x77, 0x6f, 0x72,
	0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x52, 0x03, 0x6a,
	0x6f, 0x62, 0x12, 0x2f, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x08, 0x70, 0x61, 0x72, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x0a, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x5f, 0x6f,
	0x6e, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c,
	0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x09, 0x64, 0x65, 0x70,
	0x65, 0x6e, 0x64, 0x73, 0x4f, 0x6e, 0x12, 0x2d, 0x0a, 0x05, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x4c, 0x65, 0x61, 0x73, 0x65, 0x52, 0x05,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x74, 0x72, 0x79, 0x5f, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x72, 0x65, 0x74, 0x72,
	0x79, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x69, 0x6e, 0x70, 0x75, 0x74,
	0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x22, 0x3e, 0x0a, 0x0e,
	0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x31, 0x0a, 0x05,
	0x54, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x28, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x22,
	0xe6, 0x01, 0x0a, 0x0e, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x6c,
	0x75, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x53, 0x6c, 0x75, 0x67, 0x12, 0x3c, 0x0a, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x77, 0x6f, 0x72, 0x6b,
	0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x0a, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66,
	0x69, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x69, 0x73,
	0x70, 0x6c, 0x61, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x69, 0x73, 0x70,
	0x6c, 0x61, 0x79, 0x12, 0x22, 0x0a, 0x0c, 0x64, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x65, 0x6e, 0x63,
	0x69, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0c, 0x64, 0x65, 0x70, 0x65, 0x6e,
	0x64, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x61, 0x78, 0x5f, 0x72,
	0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6d, 0x61,
	0x78, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x1a, 0x0a, 0x04, 0x55, 0x55, 0x49, 0x44,
	0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x75, 0x75, 0x69, 0x64, 0x22, 0xa9, 0x01, 0x0a, 0x09, 0x54, 0x61, 0x73, 0x6b, 0x4c, 0x65, 0x61,
	0x73, 0x65, 0x12, 0x2f, 0x0a, 0x05, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x6c, 0x65,
	0x61, 0x73, 0x65, 0x12, 0x6b, 0x0a, 0x25, 0x72, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64,
	0x65, 0x64, 0x5f, 0x77, 0x61, 0x69, 0x74, 0x5f, 0x75, 0x6e, 0x74, 0x69, 0x6c, 0x5f, 0x6e, 0x65,
	0x78, 0x74, 0x5f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x21, 0x72,
	0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x57, 0x61, 0x69, 0x74, 0x55, 0x6e,
	0x74, 0x69, 0x6c, 0x4e, 0x65, 0x78, 0x74, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e,
	0x2a, 0x5c, 0x0a, 0x09, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0d, 0x0a,
	0x09, 0x55, 0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06,
	0x51, 0x55, 0x45, 0x55, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x55, 0x4e, 0x4e,
	0x49, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x4f, 0x4d, 0x50, 0x55, 0x54, 0x45,
	0x44, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x04, 0x12,
	0x0d, 0x0a, 0x09, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x4c, 0x45, 0x44, 0x10, 0x05, 0x42, 0xc4,
	0x01, 0x0a, 0x10, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73,
	0x2e, 0x76, 0x31, 0x42, 0x09, 0x43, 0x6f, 0x72, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x54, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x6c,
	0x65, 0x62, 0x6f, 0x78, 0x2f, 0x74, 0x69, 0x6c, 0x65, 0x62, 0x6f, 0x78, 0x2d, 0x67, 0x6f, 0x2f,
	0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x77, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x77, 0x6f, 0x72, 0x6b, 0x66,
	0x6c, 0x6f, 0x77, 0x73, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x57, 0x58, 0x58, 0xaa, 0x02, 0x0c, 0x57,
	0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0c, 0x57, 0x6f,
	0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x18, 0x57, 0x6f, 0x72,
	0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x73, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77,
	0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_workflows_v1_core_proto_rawDescOnce sync.Once
	file_workflows_v1_core_proto_rawDescData = file_workflows_v1_core_proto_rawDesc
)

func file_workflows_v1_core_proto_rawDescGZIP() []byte {
	file_workflows_v1_core_proto_rawDescOnce.Do(func() {
		file_workflows_v1_core_proto_rawDescData = protoimpl.X.CompressGZIP(file_workflows_v1_core_proto_rawDescData)
	})
	return file_workflows_v1_core_proto_rawDescData
}

var file_workflows_v1_core_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_workflows_v1_core_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_workflows_v1_core_proto_goTypes = []interface{}{
	(TaskState)(0),              // 0: workflows.v1.TaskState
	(*Cluster)(nil),             // 1: workflows.v1.Cluster
	(*Job)(nil),                 // 2: workflows.v1.Job
	(*Task)(nil),                // 3: workflows.v1.Task
	(*TaskIdentifier)(nil),      // 4: workflows.v1.TaskIdentifier
	(*Tasks)(nil),               // 5: workflows.v1.Tasks
	(*TaskSubmission)(nil),      // 6: workflows.v1.TaskSubmission
	(*UUID)(nil),                // 7: workflows.v1.UUID
	(*TaskLease)(nil),           // 8: workflows.v1.TaskLease
	(*durationpb.Duration)(nil), // 9: google.protobuf.Duration
}
var file_workflows_v1_core_proto_depIdxs = []int32{
	7,  // 0: workflows.v1.Job.id:type_name -> workflows.v1.UUID
	7,  // 1: workflows.v1.Task.id:type_name -> workflows.v1.UUID
	4,  // 2: workflows.v1.Task.identifier:type_name -> workflows.v1.TaskIdentifier
	0,  // 3: workflows.v1.Task.state:type_name -> workflows.v1.TaskState
	2,  // 4: workflows.v1.Task.job:type_name -> workflows.v1.Job
	7,  // 5: workflows.v1.Task.parent_id:type_name -> workflows.v1.UUID
	7,  // 6: workflows.v1.Task.depends_on:type_name -> workflows.v1.UUID
	8,  // 7: workflows.v1.Task.lease:type_name -> workflows.v1.TaskLease
	3,  // 8: workflows.v1.Tasks.tasks:type_name -> workflows.v1.Task
	4,  // 9: workflows.v1.TaskSubmission.identifier:type_name -> workflows.v1.TaskIdentifier
	9,  // 10: workflows.v1.TaskLease.lease:type_name -> google.protobuf.Duration
	9,  // 11: workflows.v1.TaskLease.recommended_wait_until_next_extension:type_name -> google.protobuf.Duration
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_workflows_v1_core_proto_init() }
func file_workflows_v1_core_proto_init() {
	if File_workflows_v1_core_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_workflows_v1_core_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cluster); i {
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
		file_workflows_v1_core_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Job); i {
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
		file_workflows_v1_core_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
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
		file_workflows_v1_core_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskIdentifier); i {
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
		file_workflows_v1_core_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tasks); i {
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
		file_workflows_v1_core_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskSubmission); i {
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
		file_workflows_v1_core_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UUID); i {
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
		file_workflows_v1_core_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskLease); i {
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
	file_workflows_v1_core_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_workflows_v1_core_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_workflows_v1_core_proto_goTypes,
		DependencyIndexes: file_workflows_v1_core_proto_depIdxs,
		EnumInfos:         file_workflows_v1_core_proto_enumTypes,
		MessageInfos:      file_workflows_v1_core_proto_msgTypes,
	}.Build()
	File_workflows_v1_core_proto = out.File
	file_workflows_v1_core_proto_rawDesc = nil
	file_workflows_v1_core_proto_goTypes = nil
	file_workflows_v1_core_proto_depIdxs = nil
}
