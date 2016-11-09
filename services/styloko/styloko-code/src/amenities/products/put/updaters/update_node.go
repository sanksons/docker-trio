package updaters

import (
	proUtil "amenities/products/common"
	put "amenities/products/put"
	validator "gopkg.in/go-playground/validator.v8"
)

// Update-Type: Node
// Original Add Node was a utility method to add
// data to product, but in later stages we decided to use
// mongo and mysql both for writes, due to this the future
// of this method is dicey.
// This method only supports mongo adapter
type Node struct {
	Type     string      `json:"type" validate:"required,eq=addNode|eq=deleteNode"`
	NodeName string      `json:"nodeName" validate:"required"`
	SKU      string      `json:"sku" validate:"required"`
	NodeData interface{} `json:"data" validate:"required"`
}

func (n *Node) getAllowedNodes() map[string]string {
	return map[string]string{
		"ShopLook":       "shopLook",
		"ShopCollection": "shopCollection",
		"SizeChart":      "sizeChart",
	}
}

func (n *Node) Response(p *proUtil.Product) interface{} {
	response := make(map[string]interface{})
	response["configId"] = p.SeqId
	response["sku"] = p.SKU
	response["updatedAt"] = p.UpdatedAt
	return response
}

func (n *Node) isNodeAllowed() bool {
	allowedNodes := n.getAllowedNodes()
	if _, ok := allowedNodes[n.NodeName]; ok {
		return true
	}
	return false
}

func (n *Node) Validate() []string {
	errs := put.Validate.Struct(n)
	if errs != nil {
		validationErrors := errs.(validator.ValidationErrors)
		msgs := proUtil.PrepareErrorMessages(validationErrors)
		return msgs
	}
	if !n.isNodeAllowed() {
		return []string{"Unsupported Node"}
	}
	return nil
}

func (n *Node) Update() (proUtil.Product, error) {

	p := &proUtil.Product{}
	allowedNodes := n.getAllowedNodes()
	nodeMongoName := allowedNodes[n.NodeName]

	var err error
	if n.Type == proUtil.ADD_NODE {
		err = proUtil.GetAdapter(put.DbAdapterName).AddNode(
			n.SKU, nodeMongoName, n.NodeData,
		)
	} else {
		err = proUtil.GetAdapter(put.DbAdapterName).DeleteNode(
			n.SKU, nodeMongoName,
		)
	}
	if err != nil {
		return *p, err
	}
	p.LoadBySku(n.SKU, put.DbAdapterName)
	return *p, err
}

func (n *Node) InvalidateCache() error {
	go func() {
		defer proUtil.RecoverHandler("Node#Invaliddate Cache")
		put.CacheMngr.DeleteBySku([]string{n.SKU}, true)
	}()
	return nil
}

func (n *Node) Publish() error {
	pro, err := proUtil.GetAdapter(proUtil.DB_READ_ADAPTER).GetBySku(
		n.SKU,
	)
	if err != nil {
		return err
	}
	go func() {
		defer proUtil.RecoverHandler("Node#Publish")
		pro.Publish("", true)
	}()
	return nil
}

//
// Acquire Lock
//
func (n *Node) Lock() bool {
	return true
}

//
// Release Lock
//
func (n *Node) UnLock() bool {
	return true
}
