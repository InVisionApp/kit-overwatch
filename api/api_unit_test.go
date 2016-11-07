// +build unit

package api

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/InVisionApp/kit-overwatch/config"
	"github.com/InVisionApp/kit-overwatch/deps"
)

var _ = Describe("API", func() {
	var (
		request  *http.Request
		response *httptest.ResponseRecorder

		cfg *config.Config
		d   *deps.Dependencies
		api *Api

		testListenAddress = "0.0.0.0:8181"
		testVersion       = "1.0.1"
	)

	BeforeEach(func() {
		// Create our fake dependencies
		d = &deps.Dependencies{}

		cfg = config.New()
		cfg.ListenAddress = testListenAddress

		api = New(cfg, d, testVersion)

		response = httptest.NewRecorder()
	})

	Describe("New", func() {
		Context("when instantiating an api", func() {
			It("should have correct attributes", func() {
				Expect(api.Config).ToNot(BeNil())
				Expect(api.Version).To(Equal(testVersion))
			})
		})
	})

	Describe("GET /", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/", nil)
			api.HomeHandler(response, request)
		})

		It("should return 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("should return info about API usage", func() {
			Expect(response.Body).To(ContainSubstring("README.md"))
		})
	})

	Describe("GET /version", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/version", nil)
			api.VersionHandler(response, request)
		})

		It("should return 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("should return version string", func() {
			Expect(response.Body).To(ContainSubstring(testVersion))
		})
	})

	Describe("GET /healthcheck", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/healthcheck", nil)
			api.HealthHandler(response, request)
		})

		It("should return 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("should have a friendly and positive message", func() {
			Expect(response.Body).To(ContainSubstring("peechy"))
		})
	})
})
