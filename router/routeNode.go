package router

import (
	"regexp"
	"strings"
)

var routeParameterRegex = regexp.MustCompile(`{(.+)}`)

type routeNode struct {
	routePart string
	subRoutes []*routeNode
	handlers  map[string]HandlerFunc
}

func (r *routeNode) match(s string) bool {
	return r.routePart == s || (routeParameterRegex.MatchString(r.routePart) && routeParameterRegex.MatchString(s))
}

func (r *routeNode) extractRouteParam(s string) (string, string, bool) {
	if !routeParameterRegex.MatchString(r.routePart) {
		return "", "", false
	}

	return r.routePart[1 : len(r.routePart)-1], s, true
}

func newRouteNode(url string, method string, f HandlerFunc) *routeNode {
	routePartRegex := regexp.MustCompile(`^/([^\/]+)`)
	urlPart := routePartRegex.FindStringIndex(url)
	if urlPart == nil {
		return nil
	}

	subroutes := make([]*routeNode, 0)
	handler := make(map[string]HandlerFunc)
	if child := newRouteNode(url[urlPart[1]:], method, f); child != nil {
		subroutes = append(subroutes, child)
	} else {
		handler[method] = f
	}

	urlPartStart := urlPart[0] + 1
	urlPartEnd := urlPart[1]
	return &routeNode{url[urlPartStart:urlPartEnd], subroutes, handler}
}

func (r *routeNode) String() (s string) {
	var sb strings.Builder

	queue := make([]*routeNode, 2)
	queue[0] = r
	queue[1] = nil
	for len(queue) > 1 {
		root := queue[0]
		queue = queue[1:]
		if root == nil {
			sb.WriteString("\n")
			if len(queue) == 0 {
				break
			}
			queue = append(queue, nil)
			continue
		} else {
			sb.WriteString(root.routePart + ", ")
		}

		queue = append(queue, root.subRoutes...)
	}

	return sb.String()
}
