package workflows // import "github.com/tilebox/tilebox-go/workflows/v1"

import (
	"encoding/binary"
	"slices"

	"github.com/cespare/xxhash/v2"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/samber/lo"
	workflowsv1 "github.com/tilebox/tilebox-go/protogen/workflows/v1"
)

// futureTask represents a prepared task ready for submission.
type futureTask struct {
	clusterSlug  string
	identifier   TaskIdentifier
	input        []byte
	dependencies []uint32
	maxRetries   int64
	optional     bool
}

// taskGroupKey is used to group tasks by their dependencies and dependants.
// Tasks with identical keys can be merged into a single TaskSubmissionGroup.
type taskGroupKey struct {
	dependencies uint64
	dependants   uint64
}

type taskSubmissionGroup struct {
	dependenciesOnOtherGroups []uint32
	inputs                    [][]byte
	identifierPointers        []uint64
	clusterSlugPointers       []uint64
	displayPointers           []uint64
	maxRetriesValues          []int64
	optionalValues            []bool
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

// fastIndexLookup provides O(1) index lookup for unique values in a slice while preserving insertion order.
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

// mergeFutureTasksToSubmissions converts a slice of futureTask to the TaskSubmissions protobuf format.
// Tasks are grouped by their dependencies and dependants to optimize serialization.
func mergeFutureTasksToSubmissions(submissions []*futureTask) *workflowsv1.TaskSubmissions {
	if len(submissions) == 0 {
		return nil
	}

	// Build dependency and dependant graphs
	dependants := make([][]uint32, len(submissions)) // reverse graph relationship of dependencies

	for i := range uint32(len(submissions)) { //nolint:gosec // we validate len(submissions) doesn't overflow uint32
		for _, dependency := range submissions[i].dependencies {
			dependants[dependency] = append(dependants[dependency], i)
		}
	}

	// Group tasks by their (dependencies, dependants) key
	groupKeys := newFastIndexLookup[taskGroupKey]()
	groups := make([]*taskSubmissionGroup, 0)

	// we keep track of which task ends up in which group, so we can convert task dependencies to group dependencies
	taskIndexToGroupIndex := make([]uint32, 0, len(submissions))

	for i := range uint32(len(submissions)) { //nolint:gosec // we validate len(submissions) doesn't overflow uint32
		groupKey := taskGroupKey{
			dependencies: comparableIntSetKey(submissions[i].dependencies...),
			dependants:   comparableIntSetKey(dependants[i]...),
		}
		groupIndex := groupKeys.appendIfUnique(groupKey)
		if groupIndex == uint64(len(groups)) { // it was a new unique group
			groups = append(groups, &taskSubmissionGroup{
				dependenciesOnOtherGroups: submissions[i].dependencies,
			})
		}
		taskIndexToGroupIndex = append(taskIndexToGroupIndex, uint32(groupIndex)) //nolint:gosec // groupIndex is <= len(submissions) so it's safe
	}

	// Now convert our dependencies from task dependencies to group dependencies
	for _, group := range groups {
		if len(group.dependenciesOnOtherGroups) == 0 {
			group.dependenciesOnOtherGroups = nil // group has no dependencies, so set to nil (rather than empty []uint32)
			continue
		}
		uniqueGroupDependencies := mapset.NewSet[uint32]()
		for _, dependencyTaskIndex := range group.dependenciesOnOtherGroups {
			uniqueGroupDependencies.Add(taskIndexToGroupIndex[dependencyTaskIndex])
		}
		group.dependenciesOnOtherGroups = uniqueGroupDependencies.ToSlice()
		slices.Sort(group.dependenciesOnOtherGroups)
	}

	clusterSlugs := newFastIndexLookup[string]()
	identifiers := newFastIndexLookup[TaskIdentifier]()
	displays := newFastIndexLookup[string]()

	for i, submission := range submissions {
		groupIndex := taskIndexToGroupIndex[i]
		group := groups[groupIndex]

		group.inputs = append(group.inputs, submission.input)
		group.identifierPointers = append(group.identifierPointers, identifiers.appendIfUnique(submission.identifier))
		group.clusterSlugPointers = append(group.clusterSlugPointers, clusterSlugs.appendIfUnique(submission.clusterSlug))
		group.displayPointers = append(group.displayPointers, displays.appendIfUnique(submission.identifier.Display()))
		group.maxRetriesValues = append(group.maxRetriesValues, submission.maxRetries)
		group.optionalValues = append(group.optionalValues, submission.optional)
	}

	groupsProto := lo.Map(groups, func(group *taskSubmissionGroup, _ int) *workflowsv1.TaskSubmissionGroup {
		return workflowsv1.TaskSubmissionGroup_builder{
			DependenciesOnOtherGroups: group.dependenciesOnOtherGroups,
			Inputs:                    group.inputs,
			IdentifierPointers:        group.identifierPointers,
			ClusterSlugPointers:       group.clusterSlugPointers,
			DisplayPointers:           group.displayPointers,
			MaxRetriesValues:          group.maxRetriesValues,
			OptionalValues:            group.optionalValues,
		}.Build()
	})

	protoIdentifiers := lo.Map(identifiers.values, func(id TaskIdentifier, _ int) *workflowsv1.TaskIdentifier {
		return workflowsv1.TaskIdentifier_builder{
			Name:    id.Name(),
			Version: id.Version(),
		}.Build()
	})

	return workflowsv1.TaskSubmissions_builder{
		TaskGroups:        groupsProto,
		ClusterSlugLookup: clusterSlugs.values,
		IdentifierLookup:  protoIdentifiers,
		DisplayLookup:     displays.values,
	}.Build()
}
