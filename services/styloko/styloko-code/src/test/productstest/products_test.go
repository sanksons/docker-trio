package productstest

import (
	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
	"test/servicetest"
	testUtil "test/utils"
	"testing"
)

func TestProducts(t *testing.T) {
	gm.RegisterFailHandler(gk.Fail)
	gk.RunSpecs(t, "Product test Suite")
}

var _ = gk.Describe("Test Product API", func() {
	servicetest.InitializeTestService()

	apiName := "catalog"
	version := "v1"
	baseUrl := apiName + "/" + version + "/products/"
	testProductIdValid := "2"
	testProductIdNotFound := "0"
	getHeaders := map[string]string{"Expanse": "XLarge", "Visibility-Type": "MSKU", "no-cache": "True"}
	//putHeaders := map[string]string{"Update-Type": "Product"}
	testProductSkuValid := "?sku=[MA142WA97AAC,ZI056MA47EDKINDFAS]"

	//Get Product by Id -> Success
	gk.Describe("Testing Get Product by ID API", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testProductIdValid)
		request = testUtil.SetHeadersInRequest(getHeaders, request)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 for Get Products by ID API for valid ID passed", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//Get Product by Id -> Failure (searched Product ID doesnot exists)
	gk.Describe("Testing Get Product by Id API", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testProductIdNotFound)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 404 - data not found for ID passed", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(404))
			})
		})
	})

	//Get Product by SKU -> Success
	gk.Describe("Testing Get Product by SKU API", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testProductSkuValid)
		request = testUtil.SetHeadersInRequest(getHeaders, request)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 for Get Products by SKU API for valid SKU passed", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//TEST CASES FOR Product POST API
	//Post (Create Product) -> Success
	gk.Describe("Testing Product Post API for Valid data sent", func() {

		body := getPutBodyDataValid(1)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - post data is correct", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//TEST CASES FOR Product PUT API
	//Put -> Success while updating Product (valid data sent)
	gk.Describe("Testing Product Put API for valid updation", func() {

		body := getPutBodyDataValid(1)
		request := testUtil.CreateTestRequestWithBody("PUT", "/"+baseUrl, body)
		request = testUtil.SetHeadersInRequest(getPutHeader(1), request)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - put data is correct", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//TEST CASES FOR Product PUT API
	//Put -> Success while updating Product (valid data sent)
	gk.Describe("Testing Product Put API for valid updation", func() {

		body := getPutBodyDataValid(2)
		request := testUtil.CreateTestRequestWithBody("PUT", "/"+baseUrl, body)
		request = testUtil.SetHeadersInRequest(getPutHeader(2), request)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - put data is correct", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//TEST CASES FOR Product PUT API
	//Put -> Success while updating Product (valid data sent)
	gk.Describe("Testing Product Put API for valid updation", func() {

		body := getPutBodyDataValid(3)
		request := testUtil.CreateTestRequestWithBody("PUT", "/"+baseUrl, body)
		request = testUtil.SetHeadersInRequest(getPutHeader(3), request)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - put data is correct", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//TEST CASES FOR Product PUT API
	//Put -> Success while updating Product (valid data sent)
	gk.Describe("Testing Product Put API for valid updation", func() {

		body := getPutBodyDataValid(4)
		request := testUtil.CreateTestRequestWithBody("PUT", "/"+baseUrl, body)
		request = testUtil.SetHeadersInRequest(getPutHeader(4), request)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - put data is correct", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//TEST CASES FOR Product PUT API
	//Put -> Success while updating Product (valid data sent)
	gk.Describe("Testing Product Put API for valid updation", func() {

		body := getPutBodyDataValid(5)
		request := testUtil.CreateTestRequestWithBody("PUT", "/"+baseUrl, body)
		request = testUtil.SetHeadersInRequest(getPutHeader(5), request)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - put data is correct", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//TEST CASES FOR Product PUT API
	//Put -> Success while updating Product (valid data sent)
	// gk.Describe("Testing Product Put API for valid updation", func() {

	// 	body := getPutBodyDataValid(6)
	// 	request := testUtil.CreateTestRequestWithBody("PUT", "/"+baseUrl, body)
	// 	request = testUtil.SetHeadersInRequest(getPutHeader(6), request)
	// 	response := servicetest.GetResponse(request)

	// 	gk.Context("then the response", func() {

	// 		gk.It("should return HTTP 200 - put data is correct", func() {

	// 			gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
	// 			gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
	// 			gm.Expect(response.Code).To(gm.Equal(200))
	// 		})
	// 	})
	// })

	//TEST CASES FOR Product PUT API
	//Put -> Success while updating Product (valid data sent)
	gk.Describe("Testing Product Put API for valid updation", func() {

		body := getPutBodyDataValid(7)
		request := testUtil.CreateTestRequestWithBody("PUT", "/"+baseUrl, body)
		request = testUtil.SetHeadersInRequest(getPutHeader(7), request)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - put data is correct", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

})
