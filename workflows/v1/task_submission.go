package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"encoding/binary"
	"slices"

	"github.com/cespare/xxhash/v2"
	mapset "github.com/deckarep/golang-set/v2"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
)

// taskSubmission represents a prepared task ready for submission.
// It contains all the data needed to submit a task in the grouped format.
type taskSubmission struct {
	clusterSlug  string
	identifier   TaskIdentifier
	input        []byte
	dependencies []int64
	maxRetries   int64
}

// taskGroupKey is used to group tasks by their dependencies and dependants.
// Tasks with identical keys can be merged into a single TaskSubmissionGroup.
type taskGroupKey struct {
	dependencies uint64
	dependants   uint64
}

// comparableIntSetKey returns a hash key for a set of integers.
// The key is computed by sorting, deduplicating, serializing to binary, and hashing with xxhash.
func comparableIntSetKey(values ...uint32) uint64 {
	if len(values) == 0 {
		return 0
	}

	deduplicated := mapset.NewSet(values...).ToSlice()
	slices.Sort(deduplicated)

	bin := make([]byte, len(deduplicated)*4)
	for i, value := range deduplicated {
		binary.LittleEndian.PutUint32(bin[i*4:], value)
	}

	return xxhash.Sum64(bin)
}

// fastIndexLookup provides O(1) lookup for unique values while preserving insertion order.
type fastIndexLookup[K comparable] struct {
	values  []K
	indices map[K]uint64
}

func newFastIndexLookup[K comparable]() *fastIndexLookup[K] {
	return &fastIndexLookup[K]{
		values:  make([]K, 0),
		indices: make(map[K]uint64),
	}
}

func (l *fastIndexLookup[K]) appendIfUnique(value K) uint64 {
	if index, ok := l.indices[value]; ok {
		return index
	}
	newIndex := uint64(len(l.values))
	l.indices[value] = newIndex
	l.values = append(l.values, value)
	return newIndex
}

// mergeTasksToSubmissions converts a slice of taskSubmission to the TaskSubmissions protobuf format.
// Tasks are grouped by their dependencies and dependants to optimize serialization.
func mergeTasksToSubmissions(submissions []*taskSubmission) *workflowsv1.TaskSubmissions {
	if len(submissions) == 0 {
		return nil
	}

	n := len(submissions)

	// Build dependency and dependant graphs
	deps := make([][]uint32, n)
	dependants := make([][]uint32, n)

	for i, sub := range submissions {
		for _, d := range sub.dependencies {
			di := uint32(d) //nolint:gosec // dependencies are validated to be in bounds
			deps[i] = append(deps[i], di)
			dependants[di] = append(dependants[di], uint32(i)) //nolint:gosec // index i is within bounds
		}
	}

	// Group tasks by their (dependencies, dependants) key
	groupKeys := newFastIndexLookup[taskGroupKey]()
	taskToGroup := make([]uint32, n)
	groupTasks := make([][]int, 0)

	for i := range submissions {
		key := taskGroupKey{
			dependencies: comparableIntSetKey(deps[i]...),
			dependants:   comparableIntSetKey(dependants[i]...),
		}
		groupIndex := groupKeys.appendIfUnique(key)
		if groupIndex == uint64(len(groupTasks)) {
			groupTasks = append(groupTasks, nil)
		}
		taskToGroup[i] = uint32(groupIndex) //nolint:gosec // groupIndex is bounded by number of tasks
		groupTasks[groupIndex] = append(groupTasks[groupIndex], i)
	}

	// Compute group-level dependencies
	groupDeps := make([][]uint32, len(groupTasks))
	for groupIndex, tasks := range groupTasks {
		depSet := mapset.NewSet[uint32]()
		for _, taskIndex := range tasks {
			for _, depTask := range deps[taskIndex] {
				depGroup := taskToGroup[depTask]
				if depGroup == uint32(groupIndex) { //nolint:gosec // groupIndex is bounded
					continue // skip self-dependency within group
				}
				depSet.Add(depGroup)
			}
		}
		if depSet.Cardinality() > 0 {
			sortedDeps := depSet.ToSlice()
			slices.Sort(sortedDeps)
			groupDeps[groupIndex] = sortedDeps
		}
		// leave groupDeps[groupIndex] as nil if there are no dependencies
	}

	// Build lookup tables and task groups
	clusterSlugLookup := []string{}
	clusterSlugIndex := make(map[string]uint64)
	identifierLookup := []TaskIdentifier{}
	identifierIndex := make(map[string]uint64)
	displayLookup := []string{}
	displayIndex := make(map[string]uint64)

	taskGroups := make([]*workflowsv1.TaskSubmissionGroup, 0, len(groupTasks))

	for groupIndex, tasks := range groupTasks {
		inputs := make([][]byte, 0, len(tasks))
		identifierPointers := make([]uint64, 0, len(tasks))
		clusterPointers := make([]uint64, 0, len(tasks))
		displayPointers := make([]uint64, 0, len(tasks))
		maxRetriesValues := make([]int64, 0, len(tasks))

		for _, taskIndex := range tasks {
			submission := submissions[taskIndex]

			// clusterSlug
			clusterLookupIndex, ok := clusterSlugIndex[submission.clusterSlug]
			if !ok {
				clusterLookupIndex = uint64(len(clusterSlugLookup))
				clusterSlugLookup = append(clusterSlugLookup, submission.clusterSlug)
				clusterSlugIndex[submission.clusterSlug] = clusterLookupIndex
			}

			// identifier
			identifierKey := submission.identifier.Name() + "@" + submission.identifier.Version()
			identifierLookupIndex, ok := identifierIndex[identifierKey]
			if !ok {
				identifierLookupIndex = uint64(len(identifierLookup))
				identifierLookup = append(identifierLookup, submission.identifier)
				identifierIndex[identifierKey] = identifierLookupIndex
			}

			// display
			display := submission.identifier.Display()
			displayLookupIndex, ok := displayIndex[display]
			if !ok {
				displayLookupIndex = uint64(len(displayLookup))
				displayLookup = append(displayLookup, display)
				displayIndex[display] = displayLookupIndex
			}

			inputs = append(inputs, submission.input)
			identifierPointers = append(identifierPointers, identifierLookupIndex)
			clusterPointers = append(clusterPointers, clusterLookupIndex)
			displayPointers = append(displayPointers, displayLookupIndex)
			maxRetriesValues = append(maxRetriesValues, submission.maxRetries)
		}

		taskGroups = append(taskGroups, workflowsv1.TaskSubmissionGroup_builder{
			DependenciesOnOtherGroups: groupDeps[groupIndex],
			Inputs:                    inputs,
			IdentifierPointers:        identifierPointers,
			ClusterSlugPointers:       clusterPointers,
			DisplayPointers:           displayPointers,
			MaxRetriesValues:          maxRetriesValues,
		}.Build())
	}

	// Convert TaskIdentifier to protobuf format
	protoIdentifiers := make([]*workflowsv1.TaskIdentifier, len(identifierLookup))
	for i, id := range identifierLookup {
		protoIdentifiers[i] = workflowsv1.TaskIdentifier_builder{
			Name:    id.Name(),
			Version: id.Version(),
		}.Build()
	}

	return workflowsv1.TaskSubmissions_builder{
		TaskGroups:        taskGroups,
		ClusterSlugLookup: clusterSlugLookup,
		IdentifierLookup:  protoIdentifiers,
		DisplayLookup:     displayLookup,
	}.Build()
}
