package controller

import (
	"github.com/gin-gonic/gin"
)

type HealthOk struct {
	Status string `json:"status"`
}
type Health struct{}

func (Health) GetHealth(c *gin.Context) {
	ok := new(HealthOk)
	ok.Status = "ok"
	ReturnSuccess(c, 200, "正常", ok)
}
