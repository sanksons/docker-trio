package categoriestest

var invalidPutData = `
{
  "sizechartActive": 0,
  "pdfActive":3,
  "sizeChartApplicable":1
}
`
var invalidPostData = `
{
    "pdfActive":1,
    "status":"active",
    "parent":423,
    "urlKey":"/fa"
}
`

// GetInvalidData -> Returns invalid data for each request method
func GetInvalidData(ty string) string {
	switch ty {
	case "POST":
		return invalidPostData
	case "PUT":
		return invalidPutData
	default:
		return ``
	}
}
