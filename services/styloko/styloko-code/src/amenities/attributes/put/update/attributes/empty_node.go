package attributes

import workflow "github.com/jabong/floRest/src/common/orchestrator"

// EmptyNode -> struct for node based data
type EmptyNode struct {
	id string
}

// SetID -> set ID for current node from orchestrator
func (cs *EmptyNode) SetID(id string) {
	cs.id = id
}

// GetID -> returns current node ID to orchestrator
func (cs EmptyNode) GetID() (id string, err error) {
	return cs.id, nil
}

// Name -> Returns node name to orchestrator
func (cs EmptyNode) Name() string {
	return "EmptyNode"
}

// Execute -> Starts node execution.
// If no query params are found, then all active categories are returned.
// TODO: Caching strategy.
func (cs EmptyNode) Execute(io workflow.WorkFlowData) (workflow.WorkFlowData, error) {
	return io, nil
}
