package ecosystem

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ActionResult struct {
	Code    string      `json:"code"`
	Msg     string      `json:"msg"`
	TipType string      `json:"tipType"`
	Tip     string      `json:"tip"`
	Data    interface{} `json:"data"`
}

func WriteResponse(c *gin.Context, actionResult *ActionResult) {
	c.JSON(http.StatusOK, *actionResult)
}
