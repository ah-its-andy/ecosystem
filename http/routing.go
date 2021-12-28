package http

import "github.com/gin-gonic/gin"

type RouteMap struct {
	areaItems []*areaItem
}

type routeItem struct {
	uri    string
	f      func(*gin.Context) (interface{}, error)
	method string
}

type areaItem struct {
	uri        string
	routeItems []*routeItem
}
