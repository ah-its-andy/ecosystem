package http

import (
	"github.com/ah-its-andy/ecosystem/config"
	"github.com/gin-gonic/gin"
)

type RouteMap struct {
	appConfig config.ConfigSection
	areaItems []*AreaItem
}

type RouteItem struct {
	uri    string
	f      func(*gin.Context) (interface{}, error)
	method string
}

type AreaItem struct {
	uri        string
	routeItems []*RouteItem
}

func (r *RouteMap) Area(areaName string, f func(*AreaItem)) {
	areaItem := &AreaItem{
		uri:        areaName,
		routeItems: make([]*RouteItem, 0),
	}
	f(areaItem)
	r.areaItems = append(r.areaItems, areaItem)
}

func (area *AreaItem) Route(method string, action string, f func(*gin.Context) (interface{}, error)) {
	area.routeItems = append(area.routeItems, &RouteItem{
		uri:    action,
		f:      f,
		method: method,
	})
}
