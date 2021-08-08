package router

import (
	"net/http"
	"strings"
)

// Map of all route parameters.
// If we have defined a route /cars/{id} and the server gets a request
// on the url /cars/123, the RouteParameters map would have a key `id`
// with a value `123`.
type RouteParameters map[string]interface{}

// A function that handles the request.
type HandlerFunc func(http.ResponseWriter, *http.Request, RouteParameters)

// Contains all the routes the server has defined.
// It is also responsible for despatching request to the appropriate handler.
type Router struct {
	root *routeNode
}

// Factory to construct a Router instance.
func NewRouter() *Router {
	root := &routeNode{"", make([]*routeNode, 0), make(map[string]HandlerFunc)}

	return &Router{root}
}

// Dispatches the request to the appropriate handler.
// This function should be passed to the `http.HandleFunc`
// on the root url: `http.HandleFunc("/", r.Dispatch)
func (d *Router) Dispatch(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	url := r.URL.Path
	routeParams := make(RouteParameters)

	var urlParts []string
	if url == "/" {
		urlParts = make([]string, 0)
	} else {
		urlParts = strings.Split(url, "/")[1:]
	}

	t := d.root
	for _, urlPart := range urlParts {
		found := false

		for _, subRoute := range t.subRoutes {
			if subRoute.routePart == urlPart {
				t = subRoute
				found = true
				break
			}

			if key, value, matches := subRoute.extractRouteParam(urlPart); matches {
				t = subRoute
				routeParams[key] = value
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "Not found.", http.StatusNotFound)
			return
		}
	}

	if _, ok := t.handlers[method]; !ok {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	t.handlers[method](w, r, routeParams)
}

// Registers `f` as the handler for all GET requests on the URL `url`
func (d *Router) Get(url string, f HandlerFunc) {
	d.add("GET", url, f)
}

// Registers `f` as the handler for all HEAD requests on the URL `url`
func (d *Router) Head(url string, f HandlerFunc) {
	d.add("HEAD", url, f)
}

// Registers `f` as the handler for all POST requests on the URL `url`
func (d *Router) Post(url string, f HandlerFunc) {
	d.add("POST", url, f)
}

// Registers `f` as the handler for all PUT requests on the URL `url`
func (d *Router) Put(url string, f HandlerFunc) {
	d.add("PUT", url, f)
}

// Registers `f` as the handler for all DELETE requests on the URL `url`
func (d *Router) Delete(url string, f HandlerFunc) {
	d.add("DELETE", url, f)
}

// Registers `f` as the handler for all CONNECT requests on the URL `url`
func (d *Router) Connect(url string, f HandlerFunc) {
	d.add("CONNECT", url, f)
}

// Registers `f` as the handler for all OPTIONS requests on the URL `url`
func (d *Router) Options(url string, f HandlerFunc) {
	d.add("OPTIONS", url, f)
}

// Registers `f` as the handler for all TRACE requests on the URL `url`
func (d *Router) Trace(url string, f HandlerFunc) {
	d.add("TRACE", url, f)
}

// Registers `f` as the handler for all PATCH requests on the URL `url`
func (d *Router) Patch(url string, f HandlerFunc) {
	d.add("PATCH", url, f)
}

func (d *Router) add(method string, url string, f HandlerFunc) {
	if len(url) == 0 || !strings.HasPrefix(url, "/") {
		url = addSlashPrefixToUrl(url)
	}

	if url == "/" {
		d.root.handlers[method] = f
	}

	urlParts := strings.Split(url, "/")[1:]
	t := d.root
	for i, urlPart := range urlParts {
		addNew := true
		for _, subRoute := range t.subRoutes {
			if subRoute.match(urlPart) {
				t = subRoute
				addNew = false
				break
			}
		}

		if addNew {
			t.subRoutes = append(t.subRoutes, newRouteNode("/"+strings.Join(urlParts[i:], "/"), method, f))
			break
		}
	}
}

func addSlashPrefixToUrl(url string) string {
	var sb strings.Builder
	sb.WriteString("/")
	sb.WriteString(url)
	return sb.String()
}
