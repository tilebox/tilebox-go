package workflows

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
)

func TestComparableIntSetKey(t *testing.T) {
	tests := []struct {
		name    string
		a       []uint32
		b       []uint32
		equal   bool
		aIsZero bool
	}{
		{
			name:    "empty input returns zero",
			a:       []uint32{},
			aIsZero: true,
		},
		{
			name:  "order insensitive",
			a:     []uint32{1, 2, 3},
			b:     []uint32{3, 2, 1},
			equal: true,
		},
		{
			name:  "duplicate insensitive",
			a:     []uint32{1, 2, 3},
			b:     []uint32{1, 2, 3, 2, 3, 1},
			equal: true,
		},
		{
			name:  "order and duplicate insensitive combined",
			a:     []uint32{1, 2, 3},
			b:     []uint32{3, 2, 1, 2, 3},
			equal: true,
		},
		{
			name:  "different sets produce different keys",
			a:     []uint32{1, 2, 3},
			b:     []uint32{1, 2, 4},
			equal: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyA := comparableIntSetKey(tt.a...)
			if tt.aIsZero {
				assert.Equal(t, uint64(0), keyA)
				return
			}
			keyB := comparableIntSetKey(tt.b...)
			if tt.equal {
				assert.Equal(t, keyA, keyB)
			} else {
				assert.NotEqual(t, keyA, keyB)
			}
		})
	}
}

func TestMergeTasksToSubmissions(t *testing.T) {
	tests := []struct {
		name        string
		submissions []*taskSubmission
		validate    func(t *testing.T, result *workflowsv1.TaskSubmissions)
	}{
		{
			name:        "nil input returns nil",
			submissions: nil,
			validate: func(t *testing.T, result *workflowsv1.TaskSubmissions) {
				assert.Nil(t, result)
			},
		},
		{
			name:        "empty slice returns nil",
			submissions: []*taskSubmission{},
			validate: func(t *testing.T, result *workflowsv1.TaskSubmissions) {
				assert.Nil(t, result)
			},
		},
		{
			name: "single task",
			submissions: []*taskSubmission{
				{
					clusterSlug:  "cluster1",
					identifier:   NewTaskIdentifier("TaskA", "v1.0"),
					input:        []byte("payload"),
					dependencies: nil,
					maxRetries:   3,
				},
			},
			validate: func(t *testing.T, result *workflowsv1.TaskSubmissions) {
				require.NotNil(t, result)
				assert.Len(t, result.GetTaskGroups(), 1)

				group := result.GetTaskGroups()[0]
				assert.Nil(t, group.GetDependenciesOnOtherGroups())
				assert.Equal(t, [][]byte{[]byte("payload")}, group.GetInputs())
				assert.Equal(t, []uint64{0}, group.GetIdentifierPointers())
				assert.Equal(t, []uint64{0}, group.GetClusterSlugPointers())
				assert.Equal(t, []uint64{0}, group.GetDisplayPointers())
				assert.Equal(t, []int64{3}, group.GetMaxRetriesValues())

				assert.Equal(t, []string{"cluster1"}, result.GetClusterSlugLookup())
				assert.Len(t, result.GetIdentifierLookup(), 1)
				assert.Equal(t, "TaskA", result.GetIdentifierLookup()[0].GetName())
				assert.Equal(t, "v1.0", result.GetIdentifierLookup()[0].GetVersion())
				assert.Equal(t, []string{"TaskA"}, result.GetDisplayLookup())
			},
		},
		{
			name: "multiple independent tasks",
			submissions: []*taskSubmission{
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskA", "v1.0"), input: []byte("input1"), dependencies: nil, maxRetries: 0},
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskA", "v1.0"), input: []byte("input2"), dependencies: nil, maxRetries: 0},
				{clusterSlug: "other", identifier: NewTaskIdentifier("TaskB", "v1.0"), input: []byte("input3"), dependencies: nil, maxRetries: 1},
			},
			validate: func(t *testing.T, result *workflowsv1.TaskSubmissions) {
				require.NotNil(t, result)
				// All tasks have no dependencies and no dependants, so they should merge into one group
				assert.Len(t, result.GetTaskGroups(), 1)

				group := result.GetTaskGroups()[0]
				assert.Nil(t, group.GetDependenciesOnOtherGroups())
				assert.Len(t, group.GetInputs(), 3)
				assert.Equal(t, [][]byte{[]byte("input1"), []byte("input2"), []byte("input3")}, group.GetInputs())

				// Verify lookup tables are deduplicated
				assert.Equal(t, []string{"test", "other"}, result.GetClusterSlugLookup())
				assert.Len(t, result.GetIdentifierLookup(), 2)
				assert.Equal(t, []string{"TaskA", "TaskB"}, result.GetDisplayLookup())

				// Verify pointers reference correct lookup indices
				assert.Equal(t, []uint64{0, 0, 1}, group.GetClusterSlugPointers())
				assert.Equal(t, []uint64{0, 0, 1}, group.GetIdentifierPointers())
				assert.Equal(t, []uint64{0, 0, 1}, group.GetDisplayPointers())
				assert.Equal(t, []int64{0, 0, 1}, group.GetMaxRetriesValues())
			},
		},
		{
			// Mirrors Python test: test_merge_future_tasks_to_submissions_dependencies
			// tasks_1 (indices 0,1): no deps
			// tasks_2 (indices 2,3): no deps
			// task_3 (index 4): depends on tasks_1 (0,1)
			// task_4 (index 5): no deps
			// task_5 (index 6): depends on tasks_2 (2,3)
			name: "with dependencies",
			submissions: []*taskSubmission{
				// tasks_1
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskA", "v1.0"), input: []byte("task0"), dependencies: nil},
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskA", "v1.0"), input: []byte("task1"), dependencies: nil},
				// tasks_2
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskA", "v1.0"), input: []byte("task2"), dependencies: nil},
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskA", "v1.0"), input: []byte("task3"), dependencies: nil},
				// task_3 depends on tasks_1
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskB", "v1.0"), input: []byte("task4"), dependencies: []int64{0, 1}},
				// task_4 no deps
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskB", "v1.0"), input: []byte("task5"), dependencies: nil},
				// task_5 depends on tasks_2
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskB", "v1.0"), input: []byte("task6"), dependencies: []int64{2, 3}},
			},
			validate: func(t *testing.T, result *workflowsv1.TaskSubmissions) {
				require.NotNil(t, result)
				// tasks_1, tasks_2, and task_4 should NOT be merged because they have different dependants
				// task_3 and task_5 should NOT be merged because they have different dependencies
				assert.Len(t, result.GetTaskGroups(), 5)

				// Group 0: tasks_1 (indices 0,1) - no deps, dependants are {4}
				assert.Nil(t, result.GetTaskGroups()[0].GetDependenciesOnOtherGroups())
				assert.Equal(t, [][]byte{[]byte("task0"), []byte("task1")}, result.GetTaskGroups()[0].GetInputs())

				// Group 1: tasks_2 (indices 2,3) - no deps, dependants are {6}
				assert.Nil(t, result.GetTaskGroups()[1].GetDependenciesOnOtherGroups())
				assert.Equal(t, [][]byte{[]byte("task2"), []byte("task3")}, result.GetTaskGroups()[1].GetInputs())

				// Group 2: task_3 (index 4) - depends on group 0
				assert.Equal(t, []uint32{0}, result.GetTaskGroups()[2].GetDependenciesOnOtherGroups())
				assert.Equal(t, [][]byte{[]byte("task4")}, result.GetTaskGroups()[2].GetInputs())

				// Group 3: task_4 (index 5) - no deps, no dependants
				assert.Nil(t, result.GetTaskGroups()[3].GetDependenciesOnOtherGroups())
				assert.Equal(t, [][]byte{[]byte("task5")}, result.GetTaskGroups()[3].GetInputs())

				// Group 4: task_5 (index 6) - depends on group 1
				assert.Equal(t, []uint32{1}, result.GetTaskGroups()[4].GetDependenciesOnOtherGroups())
				assert.Equal(t, [][]byte{[]byte("task6")}, result.GetTaskGroups()[4].GetInputs())
			},
		},
		{
			// Mirrors Python test: test_merge_future_tasks_two_separate_branches
			//
			//           task_a (0)
			//          /         \
			//   task_b_left(1)   task_b_right(3)
			//        |                |
			//   task_c_left(2)   task_c_right(4)
			name: "two separate branches",
			submissions: []*taskSubmission{
				// task_a (index 0)
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskA", "v1.0"), input: []byte("a"), dependencies: nil},
				// task_b_left (index 1) depends on task_a
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskB", "v1.0"), input: []byte("b_left"), dependencies: []int64{0}},
				// task_c_left (index 2) depends on task_b_left
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskB", "v1.0"), input: []byte("c_left"), dependencies: []int64{1}},
				// task_b_right (index 3) depends on task_a
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskB", "v1.0"), input: []byte("b_right"), dependencies: []int64{0}},
				// task_c_right (index 4) depends on task_b_right
				{clusterSlug: "test", identifier: NewTaskIdentifier("TaskB", "v1.0"), input: []byte("c_right"), dependencies: []int64{3}},
			},
			validate: func(t *testing.T, result *workflowsv1.TaskSubmissions) {
				require.NotNil(t, result)
				// 5 groups: task_a, task_b_left, task_c_left, task_b_right, task_c_right
				// task_b_left and task_b_right have same deps {0} but different dependants ({2} vs {4})
				assert.Len(t, result.GetTaskGroups(), 5)

				// Group 0: task_a - no deps
				assert.Nil(t, result.GetTaskGroups()[0].GetDependenciesOnOtherGroups())
				assert.Equal(t, [][]byte{[]byte("a")}, result.GetTaskGroups()[0].GetInputs())

				// Group 1: task_b_left - depends on group 0
				assert.Equal(t, []uint32{0}, result.GetTaskGroups()[1].GetDependenciesOnOtherGroups())
				assert.Equal(t, [][]byte{[]byte("b_left")}, result.GetTaskGroups()[1].GetInputs())

				// Group 2: task_c_left - depends on group 1
				assert.Equal(t, []uint32{1}, result.GetTaskGroups()[2].GetDependenciesOnOtherGroups())
				assert.Equal(t, [][]byte{[]byte("c_left")}, result.GetTaskGroups()[2].GetInputs())

				// Group 3: task_b_right - depends on group 0
				assert.Equal(t, []uint32{0}, result.GetTaskGroups()[3].GetDependenciesOnOtherGroups())
				assert.Equal(t, [][]byte{[]byte("b_right")}, result.GetTaskGroups()[3].GetInputs())

				// Group 4: task_c_right - depends on group 3
				assert.Equal(t, []uint32{3}, result.GetTaskGroups()[4].GetDependenciesOnOtherGroups())
				assert.Equal(t, [][]byte{[]byte("c_right")}, result.GetTaskGroups()[4].GetInputs())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeTasksToSubmissions(tt.submissions)
			tt.validate(t, result)
		})
	}
}

// TestMergeTasksToSubmissions_ManyTasks mirrors Python test: test_merge_future_tasks_to_submissions_many_tasks
func TestMergeTasksToSubmissions_ManyTasks(t *testing.T) {
	n := 100
	submissions := make([]*taskSubmission, 0, n*2)

	// tasks_1: n tasks with no dependencies
	for i := range n {
		submissions = append(submissions, &taskSubmission{
			clusterSlug:  "test",
			identifier:   NewTaskIdentifier("TaskA", "v1.0"),
			input:        []byte{byte(i)},
			dependencies: nil,
		})
	}

	// tasks_2: n tasks that all depend on all of tasks_1
	deps := make([]int64, n)
	for i := range n {
		deps[i] = int64(i)
	}
	for i := range n {
		submissions = append(submissions, &taskSubmission{
			clusterSlug:  "test",
			identifier:   NewTaskIdentifier("TaskB", "v1.0"),
			input:        []byte{byte(n + i)},
			dependencies: deps,
		})
	}

	result := mergeTasksToSubmissions(submissions)

	require.NotNil(t, result)
	// All tasks_1 share same deps (none) and same dependants (all of tasks_2), so they merge into one group
	// All tasks_2 share same deps (all of tasks_1) and same dependants (none), so they merge into one group
	assert.Len(t, result.GetTaskGroups(), 2)

	assert.Nil(t, result.GetTaskGroups()[0].GetDependenciesOnOtherGroups())
	assert.Len(t, result.GetTaskGroups()[0].GetInputs(), n)

	assert.Equal(t, []uint32{0}, result.GetTaskGroups()[1].GetDependenciesOnOtherGroups())
	assert.Len(t, result.GetTaskGroups()[1].GetInputs(), n)
}

// TestMergeTasksToSubmissions_ManyNonMergeableDependencyGroups mirrors Python test:
// test_merge_future_tasks_to_submissions_many_non_mergeable_dependency_groups
func TestMergeTasksToSubmissions_ManyNonMergeableDependencyGroups(t *testing.T) {
	n := 100
	submissions := make([]*taskSubmission, 0, n*2)

	// For each i, create task_1[i] with no deps, then task_2[i] that depends only on task_1[i]
	for i := range n {
		// task_1[i]
		submissions = append(submissions, &taskSubmission{
			clusterSlug:  "test",
			identifier:   NewTaskIdentifier("TaskA", "v1.0"),
			input:        []byte{byte(i)},
			dependencies: nil,
		})
		// task_2[i] depends on task_1[i]
		submissions = append(submissions, &taskSubmission{
			clusterSlug:  "test",
			identifier:   NewTaskIdentifier("TaskB", "v1.0"),
			input:        []byte{byte(n + i)},
			dependencies: []int64{int64(i * 2)}, // task_1[i] is at index i*2
		})
	}

	result := mergeTasksToSubmissions(submissions)

	require.NotNil(t, result)
	// Each pair has unique (deps, dependants), so no merging possible
	assert.Len(t, result.GetTaskGroups(), 2*n)
}

func TestFastIndexLookup(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected []uint64
		wantLen  int
	}{
		{
			name:     "appends unique values",
			values:   []string{"a", "b", "c"},
			expected: []uint64{0, 1, 2},
			wantLen:  3,
		},
		{
			name:     "returns existing index for duplicates",
			values:   []string{"a", "b", "a", "b"},
			expected: []uint64{0, 1, 0, 1},
			wantLen:  2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lookup := newFastIndexLookup[string]()
			indices := make([]uint64, 0, len(tt.values))
			for _, v := range tt.values {
				indices = append(indices, lookup.appendIfUnique(v))
			}
			assert.Equal(t, tt.expected, indices)
			assert.Len(t, lookup.values, tt.wantLen)
		})
	}
}
