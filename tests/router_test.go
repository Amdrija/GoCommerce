package router_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Amdrija/GoCommerce/router"
)

type mockResponseWriter struct {
	content string
}

func (m *mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	m.content = string(p)
	return len(p), nil
}

func (m *mockResponseWriter) WriteHeader(int) {
}

func testHandler(w http.ResponseWriter, r *http.Request, routeParams router.RouteParameters) {
	fmt.Fprintf(w, "Url: %s | Route parameters: %v", r.URL.Path, routeParams)
}

func TestGetDispatcher(t *testing.T) {
	r := router.NewRouter()
	m1 := &mockResponseWriter{}

	r.Get("/about", testHandler)

	request, _ := http.NewRequest(http.MethodGet, "/about", nil)
	r.Dispatch(m1, request)

	m2 := &mockResponseWriter{}
	testHandler(m2, request, make(router.RouteParameters))
	if m1.content != m2.content {
		t.Error("Get on route /about not working")
	}
}

func TestGetDispatcherWithRouteParams(t *testing.T) {
	r := router.NewRouter()
	m1 := &mockResponseWriter{}

	r.Get("/about/{id}/andy/{word}/{sea}", testHandler)

	request, _ := http.NewRequest(http.MethodGet, "/about/123/andy/string/blue", nil)
	r.Dispatch(m1, request)

	m2 := &mockResponseWriter{}
	routeParameters := make(router.RouteParameters)
	routeParameters["id"] = 123
	routeParameters["word"] = "string"
	routeParameters["sea"] = "blue"
	testHandler(m2, request, routeParameters)
	if m1.content != m2.content {
		t.Error("Get on route /about/{id}/andy/{word}/{sea} not working")
	}
}

func TestGetDispatcherWithRoot(t *testing.T) {
	r := router.NewRouter()
	m1 := &mockResponseWriter{}

	r.Get("/", testHandler)

	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.Dispatch(m1, request)

	m2 := &mockResponseWriter{}
	testHandler(m2, request, make(router.RouteParameters))
	if m1.content != m2.content {
		t.Error("Get on route / not working")
	}
}
