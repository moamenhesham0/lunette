package router

import (
	"io"
	"lunette-gateway/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/gin-gonic/gin"
)

func setupReverseProxyTesting() *gin.Engine {
	gin.SetMode(gin.TestMode)
	gateway := gin.New()
	return gateway
}

func createMockService() *httptest.Server {
	mockService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Invalid Path"}`))
			return
		}

		if r.Header.Get("user_id") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	return mockService
}

func TestReverseProxy_ValidRequest(t *testing.T) {
	// Given
	mockService := createMockService()
	defer mockService.Close()

	gateway := setupReverseProxyTesting()
	gatewayServer := httptest.NewServer(gateway)
	defer gatewayServer.Close()

	requestURL := gatewayServer.URL + "/service/test"

	// When
	gateway.Use(func(c *gin.Context) {
		c.Set("user_id", uuid.New().String())
		c.Next()
	})
	gateway.Any("/service/*path", ReverseProxy(mockService.URL))

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		t.Fatalf("Gateway server can't create a request")
	}

	resp, err := http.DefaultClient.Do(req)

	// Then
	if err != nil {
		t.Fatalf("Gateway server can't send send the request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		expectationLog := util.UnitTestExpectation(resp.StatusCode, http.StatusOK)
		serviceMsg, _ := io.ReadAll(resp.Body)
		t.Errorf("%s. %s", serviceMsg, expectationLog)
	}
}

func TestReverseProxy_NoUserID(t *testing.T) {
	// Given
	mockService := createMockService()
	defer mockService.Close()

	gateway := setupReverseProxyTesting()
	gatewayServer := httptest.NewServer(gateway)
	defer gatewayServer.Close()

	requestURL := gatewayServer.URL + "/service/test"

	// When
	gateway.Any("/service/*path", ReverseProxy(mockService.URL))

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		t.Fatalf("Gateway server can't create a request")
	}

	resp, err := http.DefaultClient.Do(req)

	// Then
	if err != nil {
		t.Fatalf("Gateway server can't send send the request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		expectationLog := util.UnitTestExpectation(resp.StatusCode, http.StatusUnauthorized)
		serviceMsg, _ := io.ReadAll(resp.Body)
		t.Errorf("%s. %s", serviceMsg, expectationLog)
	}
}

func TestReverseProxy_BadPathRequest(t *testing.T) {
	// Given
	mockService := createMockService()
	defer mockService.Close()

	gateway := setupReverseProxyTesting()
	gatewayServer := httptest.NewServer(gateway)
	defer gatewayServer.Close()

	requestURL := gatewayServer.URL + "/service/nonexistent"

	// When
	gateway.Any("/service/*path", ReverseProxy(mockService.URL))

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		t.Fatalf("Gateway server can't create a request")
	}

	resp, err := http.DefaultClient.Do(req)

	// Then
	if err != nil {
		t.Fatalf("Gateway server can't send send the request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		expectationLog := util.UnitTestExpectation(resp.StatusCode, http.StatusNotFound)
		serviceMsg, _ := io.ReadAll(resp.Body)
		t.Errorf("%s. %s", serviceMsg, expectationLog)
	}
}
