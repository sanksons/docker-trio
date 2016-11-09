package delete

//
// type DeleteQuery struct {
// 	Id  int
// 	Sku string
// }

// Query -> Struct for deletion
type Query struct {
	Ids  []int
	Skus []string
}
