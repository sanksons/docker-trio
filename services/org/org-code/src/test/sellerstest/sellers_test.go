package sellerstest

import (
	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
	"test/servicetest"
	testUtil "test/utils"
	"testing"
)

func TestSellers(t *testing.T) {
	gm.RegisterFailHandler(gk.Fail)
	gk.RunSpecs(t, "Sellers test Suite")
}

var _ = gk.Describe("Test Sellers API", func() {
	servicetest.InitializeTestService()

	apiName := "org"
	version := "v1"
	baseUrl := apiName + "/" + version + "/sellers/"
	testSearchIdValid := "4"
	testSearchIdInvalid := "random"
	testSearchIdNotFound := "0"
	testSearchQuery := "?q=id.in~[123,234]"
	testLimitQuery := "?limit=0&offset=10"

	//GET
	gk.Describe("GET ALL API", func() {
		request := testUtil.CreateTestRequest("GET", "/"+baseUrl)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 200 - all the data is displayed", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	gk.Describe("GET by Id API", func() {
		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testSearchIdValid)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 200 - id data is displayed", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	gk.Describe("GET by Id API", func() {
		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testSearchIdInvalid)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 400 - id is invalid, string was passed", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})

	gk.Describe("GET by Id API", func() {
		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testSearchIdNotFound)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 404 - data not found for this id", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(404))
			})
		})
	})

	gk.Describe("GET by SEARCH API", func() {
		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testSearchQuery)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 200 - query data is displayed", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	gk.Describe("GET by SEARCH API", func() {
		request := testUtil.CreateTestRequest("GET", "/"+baseUrl+testLimitQuery)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 200 - query data is displayed after the offset", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	//POST
	gk.Describe("POST API", func() {
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

	gk.Describe("POST API", func() {
		body := getPostBodyData(2)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 400 - mandatory field is missing", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})

	/* Will be uncommented later when seller id is mandatory
	gk.Describe("POST API", func() {
		body := getPostBodyData(3)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 400 - seller id is missing", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})*/

	gk.Describe("POST API", func() {
		body := getPostBodyData(4)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 500 - post data is empty", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(500))
			})
		})
	})

	//PUT
	gk.Describe("PUT API", func() {
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

	gk.Describe("PUT API", func() {
		body := getPutBodyData(3)
		request := testUtil.CreateTestRequestWithBody("PUT", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 500 - put data is empty", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(500))
			})
		})
	})
})
