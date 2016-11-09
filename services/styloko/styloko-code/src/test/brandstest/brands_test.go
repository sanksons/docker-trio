package brandstest

import (
	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
	"test/servicetest"
	testUtil "test/utils"
	"testing"
)

func TestSellers(t *testing.T) {
	gm.RegisterFailHandler(gk.Fail)
	gk.RunSpecs(t, "Brands test Suite")
}

var _ = gk.Describe("Test Brands APIs", func() {

	servicetest.InitializeTestService()

	apiName := "catalog"
	version := "v1"
	baseUrl := apiName + "/" + version + "/brands/"
	incorrectUrl := apiName + "/" + version + "/brand"
	testBrandIdValid := "2"
	testBrandIdInvalidSyntax := "kuchBhiId"
	testBrandIdNotFound := "0"
	testSearchActiveBrandsQuery := "?status=active"
	testSearchInactiveBrandsQuery := "?status=inactive"
	testSearchDeletedBrandsQuery := "?status=deleted"
	testSearchBrandInvalidQuery := "*status=active"
	testSearchBrandValidQuery := "?status=inactive"
	testSearchInvalidStatusQuery := "?status=Blah"

	//TEST CASES FOR BRAND GET API
	//Get All Brands -> Success
	gk.Describe("Testing GetAll Brands API", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl)

		response := servicetest.GetResponse(request)
		//fmt.Println(string(response.Body.Bytes()))
		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 for GetAll Brands API", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//Get All Brands -> Failure (incorrect API Url)
	gk.Describe("Testing GetAll Brands API", func() {

		request := testUtil.CreateTestRequest("GET", "/"+incorrectUrl)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 404 for GetAll Brands API for incorrect Url", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(404))
			})
		})
	})

	//Get Brand by Id -> Success
	gk.Describe("Testing Get Brand by ID API", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testBrandIdValid)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 for Get Brands by ID API for valid ID passed", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//Get Brand by Id -> Failure (invalid syntax for Brand ID)
	gk.Describe("Testing Get Brand by Id API", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testBrandIdInvalidSyntax)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 404 - id is invalid, string was passed", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(404))
			})
		})
	})

	//Get Brand by Id -> Failure (searched Brand ID doesnot exists)
	gk.Describe("Testing Get Brand by Id API", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testBrandIdNotFound)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 404 - data not found for ID passed", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(404))
			})
		})
	})

	//Get Brand by Status -> Success
	gk.Describe("Testing Get Brands by Status API for Valid Query fired", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testSearchBrandValidQuery)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - Valid Query fired", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//Get Brand by Status -> Success (Active Brands Searched)
	gk.Describe("Testing Get Brand by Status API for Active Brands", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testSearchActiveBrandsQuery)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - Active Brands are displayed", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//Get Brand by Status -> Success (Inactive Brands Searched)
	gk.Describe("Testing Get Brand by Status API for Invalid Brands", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testSearchInactiveBrandsQuery)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - Inactive Brands are displayed", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//Get Brand by Status -> Success (Deleted Brands Searched)
	gk.Describe("Testing Get Brand by Status API for Deleted Brands", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testSearchDeletedBrandsQuery)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - Deleted Brands are displayed", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//Get Brand by Status -> Failure (Query syntax wrong)
	gk.Describe("Testing Get Brand by Status API for Invalid Query fired", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testSearchBrandInvalidQuery)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 404 - Invalid Query Syntax", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(404))
			})
		})
	})

	//Get Brand by Status -> Failure (Brand Status is invalid)
	gk.Describe("Testing Get Brand by Status API for Invalid Status passed", func() {

		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testSearchInvalidStatusQuery)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - Brand Status is invalid", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})

	//TEST CASES FOR BRAND POST API
	//Post (Create Brand) -> Success
	gk.Describe("Testing Brand Post API for Valid data sent", func() {

		body := getPostBodyData(1)
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

	//Post -> Failure(Empty Post Data)
	gk.Describe("Testing Brand Post API for empty post data body", func() {

		body := getPostBodyData(2)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 501 - post data is missing", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(501))
			})
		})
	})

	//Post -> Failure(Mandatory fields missing while Brand creation)
	gk.Describe("Testing Brand Post API for missing mandatory field", func() {
		body := getPostBodyData(3)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 400 - mandatory is missing", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})

	//Post -> Failure(Invalid value of filed while brand creation)
	gk.Describe("Testing Brand Post API for invalid field value in data body", func() {

		body := getPostBodyData(4)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 501 - invalid field value sent", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(501))
			})
		})
	})

	//TEST CASES FOR BRAND PUT API
	//Put -> Success while updating Brand (valid data sent)
	gk.Describe("Testing Brand Post API for valid updation", func() {

		body := getPutBodyData(1)
		request := testUtil.CreateTestRequestWithBody("PUT", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 200 - put data is correct", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//Put -> Failure (Brand ID doesnot exist)
	gk.Describe("PUT API", func() {

		body := getPutBodyData(2)
		request := testUtil.CreateTestRequestWithBody("PUT", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 400 - sequence id is invalid", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})

	//Put -> Failure(Put Data Body is empty)
	gk.Describe("Testing Brand Put API when put data body is empty", func() {

		body := getPutBodyData(3)
		request := testUtil.CreateTestRequestWithBody("PUT", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)

		gk.Context("then the response", func() {

			gk.It("should return HTTP 400 - put data is empty", func() {

				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})
})
