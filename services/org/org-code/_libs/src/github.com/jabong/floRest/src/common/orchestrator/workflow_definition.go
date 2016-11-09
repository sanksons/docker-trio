package orchestrator

import (
	"errors"
	"fmt"
	"github.com/jabong/floRest/src/common/collections"
)

/*
Work flow definition storage
*/
type WorkFlowDefinition struct {

	//Node of the Work Flow Definition
	nodes map[string]WorkFlowNodeInterface

	//The connections of the Work flow Definition
	edges map[string][]string

	//Mapping of the join node for the corresponding fork node
	//Every Fork Node should have a Join node
	joinFork map[string]string

	startNodeId string
}

/*
TODO: Create the workflow definition from configuration file
*/
func (d *WorkFlowDefinition) CreateFromConfig(filename string) error {
	return nil
}

/*
TODO: Create the workflow definition from configuration file
*/
func (d *WorkFlowDefinition) Create() {
	d.nodes = make(map[string]WorkFlowNodeInterface)
	d.edges = make(map[string][]string)
	d.joinFork = make(map[string]string)
}

/*
Set the start Node
*/
func (d *WorkFlowDefinition) SetStartNode(node WorkFlowNodeInterface) error {
	id, iderr := node.GetID()
	if iderr != nil {
		return iderr
	}

	_, found := d.nodes[id]
	if !found {
		errString := fmt.Sprintln("Node with provide Id: ", id, " is not added")
		return errors.New(errString)
	}

	d.startNodeId = id
	return nil

}

/*
Get the start node
*/
func (d *WorkFlowDefinition) GetStartNode() (*WorkFlowNodeInterface, error) {
	if d.startNodeId == "" {
		errString := fmt.Sprintf("Start Node is not set")
		return nil, errors.New(errString)
	}

	startNode := d.nodes[d.startNodeId]

	return &startNode, nil
}

/*
Add a execution node
*/
func (d *WorkFlowDefinition) AddExecutionNode(execNode WorkFlowExecuteNodeInterface) error {

	id, iderr := execNode.GetID()
	if iderr != nil {
		return iderr
	}

	_, found := d.nodes[id]
	if found {
		errString := fmt.Sprintln("Node with provide Id: ", id, " is already added")
		return errors.New(errString)
	}

	d.nodes[id] = execNode

	return nil
}

/*
Add a decision node
*/

func (d *WorkFlowDefinition) AddDecisionNode(decisionNode WorkFlowDecisionNodeInterface,
	yesNode WorkFlowNodeInterface, noNode WorkFlowNodeInterface) error {

	id, iderr := decisionNode.GetID()
	if iderr != nil {
		return iderr
	}

	yesNodeId, yesNodeIderr := yesNode.GetID()
	if yesNodeIderr != nil {
		return yesNodeIderr
	}

	noNodeId, noNodeIderr := noNode.GetID()
	if noNodeIderr != nil {
		return noNodeIderr
	}

	_, found := d.nodes[id]
	if found {
		errString := fmt.Sprintln("Node with provide Id: ", id, " is already added")
		return errors.New(errString)
	}

	d.nodes[id] = decisionNode
	d.nodes[yesNodeId] = yesNode
	d.nodes[noNodeId] = noNode

	d.edges[id] = []string{yesNodeId, noNodeId}

	return nil
}

/*
Add a Fork Node
*/
func (d *WorkFlowDefinition) AddForkNode(forkNode WorkFlowForkNodeInterface,
	forkNodes []WorkFlowNodeInterface) error {

	id, iderr := forkNode.GetID()
	if iderr != nil {
		return iderr
	}

	forkNodesIds, forkNodeIdserr := d.getNodeIds(forkNodes)
	if forkNodeIdserr != nil {
		return forkNodeIdserr
	}

	_, found := d.nodes[id]
	if found {
		errString := fmt.Sprintln("Node with provide Id: ", id, " is already added")
		return errors.New(errString)
	}

	d.nodes[id] = forkNode
	for index, forkNodeId := range forkNodesIds {
		d.nodes[forkNodeId] = forkNodes[index]
	}
	d.edges[id] = forkNodesIds

	return nil
}

func (d *WorkFlowDefinition) getNodeIds(nodes []WorkFlowNodeInterface) ([]string, error) {
	var ids []string

	for _, node := range nodes {
		id, iderr := node.GetID()
		if iderr != nil {
			return []string{}, iderr
		}
		ids = append(ids, id)
	}

	return ids, nil
}

/*
Add a Join Node
*/
func (d *WorkFlowDefinition) AddJoinNode(joinNode WorkFlowJoinNodeInterface) error {
	id, iderr := joinNode.GetID()
	if iderr != nil {
		return iderr
	}

	_, found := d.nodes[id]
	if found {
		errString := fmt.Sprintln("Node with provide Id: ", id, " is already added")
		return errors.New(errString)
	}

	d.nodes[id] = joinNode

	return nil
}

/*
Add connection between 2 nodes
*/
func (d *WorkFlowDefinition) AddConnection(fromNode WorkFlowNodeInterface,
	toNode WorkFlowNodeInterface) error {

	fromNodeid, fromNodeiderr := fromNode.GetID()
	if fromNodeiderr != nil {
		return fromNodeiderr
	}

	toNodeid, toNodeiderr := toNode.GetID()
	if toNodeiderr != nil {
		return toNodeiderr
	}

	_, fromNodefound := d.nodes[fromNodeid]
	if !fromNodefound {
		errString := fmt.Sprintln("Node with provide Id: ", fromNodeid, " is not present. Add the node")
		return errors.New(errString)
	}

	_, toNodefound := d.nodes[toNodeid]
	if !toNodefound {
		errString := fmt.Sprintln("Node with provide Id: ", toNodeid, " is not present. Add the node")
		return errors.New(errString)
	}

	edgeList, edgeListfound := d.edges[fromNodeid]
	if !edgeListfound {
		d.edges[fromNodeid] = []string{toNodeid}
		return nil
	}

	edgeList = append(edgeList, toNodeid)
	d.edges[fromNodeid] = edgeList

	return nil

}

//Create mapping of the fork and corresponding join nodes
func (d *WorkFlowDefinition) createJoinForkMapping() error {

	sNode, serr := d.GetStartNode()
	if serr != nil {
		return serr
	}

	return d.createJoinForkMappingHelper(*sNode, d.startNodeId, &collections.Stack{})
}

func (d *WorkFlowDefinition) String() string {
	var res string
	res = fmt.Sprintf("nodes %v: \n", d.nodes)
	res = res + fmt.Sprintf("edges  %v: \n", d.edges)
	res = res + fmt.Sprintf("join-fork  %v: \n", d.joinFork)
	res = res + fmt.Sprintf("startNodeId  %v: \n", d.startNodeId)
	return res
}
