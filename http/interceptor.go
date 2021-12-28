package http

import "github.com/gin-gonic/gin"

type Interceptor interface {
	OnExecuting(uri string, method string, ctx *gin.Context) error

	OnExecuted(uri string, method string, result interface{}, err error, ctx *gin.Context) error
}
