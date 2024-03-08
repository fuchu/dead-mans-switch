package controller

import "github.com/gin-gonic/gin"

type JsonStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
type JsonErrStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

func ReturnSuccess(c *gin.Context, code int, msg interface{}, data interface{}) {
	json := &JsonStruct{Code: code, Msg: msg, Data: data}
	c.JSON(200, json)
}

func ReturnError(c *gin.Context, code int, msg interface{}) {
	json := &JsonErrStruct{Code: code, Msg: msg}
	c.JSON(200, json)
}
