package orchestrator

import (
	"fmt"
)

/*
Instance of the workflow
*/
type Orchestrator struct {
	//Workflow definition
	workflow *WorkFlowDefinition
}

/*
TODO: This should read the workflow configuration file
and create the appropriate orchestrator
*/
func (o *Orchestrator) CreateFromConfig(fileName string) error {
	return nil
}

/*
Constructor function for the pipeline creation
create(workflow)
*/
func (o *Orchestrator) Create(workflowdefinition *WorkFlowDefinition) error {
	jfErr := workflowdefinition.createJoinForkMapping()
	if jfErr != nil {
		return jfErr
	}
	o.workflow = workflowdefinition
	return nil
}

//Helper function to execute the Execution Node
func execExecuteNode(execNodeId string,
	execNode WorkFlowExecuteNodeInterface,
	wfData *WorkFlowData,
	wfDefinition *WorkFlowDefinition) (nextNodeId string, nextwfData *WorkFlowData) {

	nextwfData = wfData

	outputData, err := execNode.Execute(*wfData)
	if err != nil {
		nextwfData.setWorkflowState(execNode.Name(), err)
		return "", nextwfData
	} else {
		nextwfData = &outputData
	}

	nextNodeIds, found := wfDefinition.edges[execNodeId]

	if !found {
		return "", nextwfData
	}

	nextNodeId = nextNodeIds[0]
	return nextNodeId, nextwfData
}

//Helper function to execute the Decision Node
func execDecisionNode(decisionNodeId string,
	decisionNode WorkFlowDecisionNodeInterface,
	wfData *WorkFlowData,
	wfDefinition *WorkFlowDefinition) (nextNodeId string, nextwfData *WorkFlowData) {

	nextwfData = wfData

	yes, err := decisionNode.GetDecision(*wfData)
	if err != nil {
		nextwfData.setWorkflowState(decisionNode.Name(), err)
		return "", nextwfData
	}

	nextNodeIds, found := wfDefinition.edges[decisionNodeId]

	if !found {
		return "", nextwfData
	}

	if yes {
		//Yes Node
		nextNodeId = nextNodeIds[0]
	}

	if !yes {
		//No Node
		nextNodeId = nextNodeIds[1]
	}

	return nextNodeId, nextwfData
}

func execForkWorkFlow(forkNodeId string,
	wfDefinition *WorkFlowDefinition,
	wfData *WorkFlowData,
	joinNodeId string,
	wfDataChannel chan *WorkFlowData) {

	//Fork Workflow path data passed into workflow data channel
	wfDataChannel <- run(forkNodeId, wfDefinition, wfData, joinNodeId)
}

//Helper function to execute the Fork Node
func execForkNode(forkNodeId string,
	forkNode WorkFlowForkNodeInterface,
	wfData *WorkFlowData,
	wfDefinition *WorkFlowDefinition) (nextNodeId string, nextwfData *WorkFlowData) {

	wfDataChannel := make(chan *WorkFlowData)

	joinNodeId := wfDefinition.joinFork[forkNodeId]
	forkNodesId := wfDefinition.edges[forkNodeId]

	//Execute concurrently the fork node paths
	for _, forkedNodeId := range forkNodesId {
		clonedWfData := wfData.Clone()
		go execForkWorkFlow(forkedNodeId, wfDefinition, &clonedWfData, joinNodeId, wfDataChannel)
	}

	var joinNodeWfData []*WorkFlowData
	for i := 0; i < len(forkNodesId); i++ {
		wfDataFromChannel := <-wfDataChannel
		joinNodeWfData = append(joinNodeWfData, wfDataFromChannel)
	}

	//Pass the data to Join Node
	jNode, found := wfDefinition.nodes[joinNodeId]
	if !found {
		errString := fmt.Sprintln("Node id ", joinNodeId, " not present for execution")
		wfData.setWorkflowState("WORKFLOW_ERROR", errString)
		return "", wfData
	}
	logInfo("Current Node : " + jNode.Name())
	if joinNode, ok := jNode.(WorkFlowJoinNodeInterface); ok {
		logInfo("Execute Node")
		return executeJoinNode(joinNodeId, joinNode, wfData, joinNodeWfData, wfDefinition)
	}

	return "", wfData
}

//Helper function to execute Join Node
func executeJoinNode(joinNodeId string,
	joinNode WorkFlowJoinNodeInterface,
	forkWfData *WorkFlowData,
	joinWfData []*WorkFlowData,
	wfDefinition *WorkFlowDefinition) (nextNodeId string, nextwfData *WorkFlowData) {

	nextwfData = forkWfData
	outputData, err := joinNode.Join(joinWfData)
	if err != nil {
		nextwfData.setWorkflowState(joinNode.Name(), err)
		return "", nextwfData
	} else {
		nextwfData = &outputData
	}

	nextNodeIds, found := wfDefinition.edges[joinNodeId]

	if !found {
		return "", nextwfData
	}

	nextNodeId = nextNodeIds[0]
	return nextNodeId, nextwfData
}

//Helper function to run the pipeline
func run(currNodeId string,
	wfDefinition *WorkFlowDefinition,
	wfData *WorkFlowData,
	terminateNodeId string) *WorkFlowData {

	logInfo("Current Node id: " + currNodeId)

	//No workflow definition
	if wfDefinition == nil {
		return wfData
	}

	node, found := wfDefinition.nodes[currNodeId]
	if !found {
		errString := fmt.Sprintln("Node id ", currNodeId, " not present for execution")
		wfData.setWorkflowState("WORKFLOW_ERROR", errString)
		return wfData
	}
	logInfo("Current Node : " + node.Name())

	var nextNodeId string
	var nextwfData *WorkFlowData = wfData

	if execNode, ok := node.(WorkFlowExecuteNodeInterface); ok {
		logInfo("Execute Node")
		nextNodeId, nextwfData = execExecuteNode(currNodeId, execNode,
			wfData, wfDefinition)
	}

	if decisionNode, ok := node.(WorkFlowDecisionNodeInterface); ok {
		logInfo("Decision Node")
		nextNodeId, nextwfData = execDecisionNode(currNodeId, decisionNode,
			wfData, wfDefinition)
	}

	if forkedNode, ok := node.(WorkFlowForkNodeInterface); ok {
		logInfo("Fork Node")
		nextNodeId, nextwfData = execForkNode(currNodeId, forkedNode,
			wfData, wfDefinition)
	}

	if nextNodeId == terminateNodeId {
		//End of execution
		return nextwfData
	}
	logInfo("Next Node id: " + nextNodeId)
	return run(nextNodeId, wfDefinition, nextwfData, terminateNodeId)
}

/*
Workflow execution begins here
the caller should create the Work Flow State
which has the InputOutput, ExecutionContext
*/
func (o *Orchestrator) Start(wfData *WorkFlowData) *WorkFlowData {

	if o.workflow == nil {
		logError("Error Empty workflow definition passed for execution")
		return new(WorkFlowData)
	}
	return run(o.workflow.startNodeId, o.workflow, wfData, "")
}

func (d *Orchestrator) String() string {
	return d.workflow.String()
}

/*
Orchestrator implements the version manager GetInstance
*/
func (o Orchestrator) GetInstance() interface{} {
	return o
}
