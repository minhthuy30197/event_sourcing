package model

import "github.com/gin-gonic/gin"

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(ctx *gin.Context, status int, err error) {
	e := HTTPError{
		Code:    status,
		Message: err.Error(),
	}

	ctx.JSON(status, e)
}
