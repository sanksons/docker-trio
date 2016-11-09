package orchestrator

import (
	"errors"
	"github.com/jabong/floRest/src/common/collections"
)

//Node Element for the Nodes Stack
type nodeElement struct {
	value WorkFlowNodeInterface
	next  *nodeElement
}

func (d *WorkFlowDefinition) forkTypeNodeHelper(forkNode WorkFlowForkNodeInterface,
	forkNodeId string,
	stck *collections.Stack) error {

	stck.Push(forkNode)
	forkedEdges, found := d.edges[forkNodeId]
	if !found {
		return errors.New("No outgoing edges from the fork node")
	}
	for _, forkedNodeId := range forkedEdges {
		forkedNode, fNodeFound := d.nodes[forkedNodeId]
		if !fNodeFound {
			return errors.New("Forked Node with the node Id not found")
		}
		//clone the stack and send to all the forked nodes for this fork node
		fStck := stck.Clone()

		ferr := d.createJoinForkMappingHelper(forkedNode, forkedNodeId, fStck)
		if ferr != nil {
			return ferr
		}
	}
	return nil
}

func (d *WorkFlowDefinition) joinTypeNodeHelper(joinNode WorkFlowJoinNodeInterface,
	joinNodeId string,
	stck *collections.Stack) error {

	if stck.IsEmpty() {
		return errors.New("Join node does not have a corresponding fork node")
	}
	stckValue := stck.Pop()

	forkNode, ok := stckValue.(WorkFlowForkNodeInterface)
	if !ok {
		return errors.New("Incorrect Node type popped from stack")
	}
	forkNodeId, ferr := forkNode.GetID()
	if ferr != nil {
		return errors.New("Fork Node does not have ID set")
	}
	d.joinFork[forkNodeId] = joinNodeId
	edges, found := d.edges[joinNodeId]
	if !found || len(edges) == 0 {
		return nil
	}
	nextNodeId := edges[0]
	nextNode := d.nodes[nextNodeId]
	return d.createJoinForkMappingHelper(nextNode, nextNodeId, stck)
}

func (d *WorkFlowDefinition) executeTypeNodeHelper(execNode WorkFlowExecuteNodeInterface,
	execNodeId string,
	stck *collections.Stack) error {

	edges, found := d.edges[execNodeId]
	if !found || len(edges) == 0 {
		return nil
	}
	nextNodeId := edges[0]
	nextNode := d.nodes[nextNodeId]
	return d.createJoinForkMappingHelper(nextNode, nextNodeId, stck)
}

func (d *WorkFlowDefinition) decisionTypeNodeHelper(decisionNode WorkFlowDecisionNodeInterface,
	decisionNodeId string,
	stck *collections.Stack) error {

	edges, found := d.edges[decisionNodeId]
	if !found || len(edges) == 0 {
		return nil
	}

	for _, nextNodeId := range edges {
		nextNode, nextNodeFound := d.nodes[nextNodeId]
		if !nextNodeFound {
			return errors.New("Node with the node Id not found")
		}
		//clone the stack and send to all the nodes for this decision node
		dStck := stck.Clone()

		derr := d.createJoinForkMappingHelper(nextNode, nextNodeId, dStck)
		if derr != nil {
			return derr
		}
	}
	return nil
}

func (d *WorkFlowDefinition) createJoinForkMappingHelper(currNode WorkFlowNodeInterface,
	currNodeId string,
	stck *collections.Stack) error {

	//Error condition
	if (currNode == nil || currNodeId == "") && !stck.IsEmpty() {
		return errors.New("Error in creating Fork Join Node mapping")
	}

	//Recursion termination
	if currNode == nil || currNodeId == "" {
		return nil
	}

	if forkNode, ok := currNode.(WorkFlowForkNodeInterface); ok {
		//Current Node is a forked node
		return d.forkTypeNodeHelper(forkNode, currNodeId, stck)
	}

	if joinNode, ok := currNode.(WorkFlowJoinNodeInterface); ok {
		//Current Node is a join node
		return d.joinTypeNodeHelper(joinNode, currNodeId, stck)
	}

	if execNode, ok := currNode.(WorkFlowExecuteNodeInterface); ok {
		//Current Node is an execution node
		return d.executeTypeNodeHelper(execNode, currNodeId, stck)
	}

	if decisionNode, ok := currNode.(WorkFlowDecisionNodeInterface); ok {
		//Current Node is a decision node
		return d.decisionTypeNodeHelper(decisionNode, currNodeId, stck)
	}

	return errors.New("Unknown Node type")
}
