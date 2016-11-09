// @APIVersion 1.0.0
// @basePath /newApp/v1
package hello

// @Title list
// @Description get hello
// @Accept  json
// @Router /hello [get]
func list() {}

// @Title geteSku
// @Description get sku
// @Accept  json
// @Param   Limit     query    string     true        "limit"
// @Param   Offset     query    string     true        "offset"
// @Param   X-Jabong-Sessionid     header    string     true        "ssn"
// @Param   X-Jabong-Token     header    string     true        "token"
// @Router /hello/update [get]
func geteSku() {}

// @Title remove
// @Description delete sku
// @Accept  json
// @Param   Sku     path    string     true        "sku to be removed"
// @Param   X-Jabong-Sessionid     header    string     true        "ssn"
// @Param   X-Jabong-Token     header    string     true        "token"
// @Router /hello [delete]
func remove() {}

// @Title add
// @Description add sku
// @Accept  json
// @Param   BodyParam     body    AddParam     true        "body"
// @Param   X-Jabong-Sessionid     header    string     true        "ssn"
// @Param   X-Jabong-Token     header    string     true        "token"
// @Router /hello [post]
func add() {}
