package standardsizetest

import (
	"test/servicetest"
	testUtil "test/utils"
	"testing"

	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
)

func TestStandardSize(t *testing.T) {
	gm.RegisterFailHandler(gk.Fail)
	gk.RunSpecs(t, "Standard Size test Suite")
}

var _ = gk.Describe("Test Standard Size API", func() {
	servicetest.InitializeTestService()

	apiName := "catalog"
	version := "v1"
	baseUrl := apiName + "/" + version + "/standardsize/"

	//POST
	gk.Describe("POST API", func() {
		body := getPostBodyData(1)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 400 - post data is empty", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})

	gk.Describe("POST API", func() {
		body := getPostBodyData(2)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 400 - invalid attribute set", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})

	gk.Describe("POST API", func() {
		body := getPostBodyData(3)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 400 - invalid leaf category", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})

	gk.Describe("POST API", func() {
		body := getPostBodyData(4)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 400 - invalid brand", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})

	gk.Describe("POST API", func() {
		body := getPostBodyData(5)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 200 - valid post data, standard size should be inserted", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	gk.Describe("POST API", func() {
		body := getPostBodyData(6)
		request := testUtil.CreateTestRequestWithBody("POST", "/"+baseUrl, body)
		response := servicetest.GetResponse(request)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 200 - standard size doesn't have a mapping", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})
})
