package categoriestest

import (
	gk "github.com/onsi/ginkgo"
	gm "github.com/onsi/gomega"
	"path"
	"test/common"
	"test/servicetest"
	"testing"
)

// TestCategories registers test suite
func TestCategories(t *testing.T) {
	gm.RegisterFailHandler(gk.Fail)
	gk.RunSpecs(t, "Categories test Suite")
}

var _ = gk.Describe("Test Categories API", func() {
	servicetest.InitializeTestService()

	apiName := "catalog"
	version := "v1"
	endpoint := "categories/"
	baseURL := path.Join(apiName, version, endpoint)

	validCategoryID := "1"
	invalidCategoryID := "0"

	invalidPath := path.Join(baseURL, validCategoryID, invalidCategoryID)

	activeStatus := "active"
	inactiveStatus := "inactive"
	deletedStatus := "deleted"
	invalidStatus := "invalid"

	// Get All API test
	gk.Describe("GET ALL API", func() {
		response := common.Get(baseURL)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 200", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	// Valid Get by ID
	gk.Describe("GET By ID API Valid ID", func() {
		response := common.Get(path.Join(baseURL, validCategoryID))
		gk.Context("then the response", func() {
			gk.It("should return HTTP 200, valid ID was provided", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	// Invalid Get by ID
	gk.Describe("GET By ID API Invalid ID", func() {
		response := common.Get(path.Join(baseURL, invalidCategoryID))
		gk.Context("then the response", func() {
			gk.It("should return HTTP 404, invalid ID was provided", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(404))
			})
		})
	})

	// Invalid Path GET response
	gk.Describe("GET API invalid Path", func() {
		response := common.Get(path.Join(baseURL, invalidPath))
		gk.Context("then the response", func() {
			gk.It("should return HTTP 501, invalid Path was provided", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(501))
			})
		})
	})

	// Active Status search GET response
	gk.Describe("GET All API active Status", func() {
		response := common.Get(baseURL + "/?status=" + activeStatus)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 200, valid status was provided", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	// Inactive Status search GET response
	gk.Describe("GET All API inactive Status", func() {
		response := common.Get(baseURL + "/?status=" + inactiveStatus)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 200, valid status was provided", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	// Deleted Status search GET response
	gk.Describe("GET All API deleted Status", func() {
		response := common.Get(baseURL + "/?status=" + deletedStatus)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 200, valid status was provided", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(200))
			})
		})
	})

	// Valid Status search GET response
	gk.Describe("GET All API invalid Status", func() {
		response := common.Get(baseURL + "/?status=" + invalidStatus)
		gk.Context("then the response", func() {
			gk.It("should return HTTP 404, invalid status was provided", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(404))
			})
		})
	})

	// PUT API validation failure response
	gk.Describe("PUT API invalid data", func() {
		response := common.Put(path.Join(baseURL, validCategoryID), GetInvalidData("PUT"))
		gk.Context("then the response", func() {
			gk.It("should return HTTP 400, invalid data provided", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})

	// POST API validation failure response
	gk.Describe("POST API invalid data", func() {
		response := common.Post(path.Join(baseURL, validCategoryID), GetInvalidData("POST"))
		gk.Context("then the response", func() {
			gk.It("should return HTTP 400, invalid data provided", func() {
				gm.Expect(response.HeaderMap.Get("Content-Type")).To(gm.Equal("application/json"))
				gm.Expect(response.HeaderMap.Get("Cache-Control")).To(gm.Equal(""))
				gm.Expect(response.Code).To(gm.Equal(400))
			})
		})
	})

})
