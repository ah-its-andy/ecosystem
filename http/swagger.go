package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func GetSwaggerHandler(host string, serviceUri string) gin.HandlerFunc {
	url := ginSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", host))
	return ginSwagger.WrapHandler(swaggerFiles.Handler, url)
}
