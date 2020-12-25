// Copyright @ 2020 OPS Inc.
//
// Author: Jinlong Yang
//

package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

// Talker 约定所有业务逻辑层(bll)的实现
type Talker interface {
	Listen(r *S3Request) (interface{}, error)
}

// Talk 约定controller的调用
func Talk(t Talker, c *gin.Context) *S3Response {
	r := parseRequest(c)

	pass, meta := CheckSign(r)
	if !pass {
		return response(FAILED, "签名验证失败!", nil)
	}
	r.UID = meta.UserID

	data, err := t.Listen(r)
	if err != nil {
		return response(FAILED, err.Error(), nil)
	}

	return response(SUCCESS, ERROR, data)
}

func parseRequest(c *gin.Context) *S3Request {
	r := &S3Request{
		HTTPMethod: c.Request.Method,
		URL:        c.Request.URL.Path,
		Query:      c.Request.URL.RawQuery,
		RequestID:  uuid.NewV4().String(),
	}

	headers := make(map[string]string)
	for k, vList := range c.Request.Header {
		headers[strings.ToLower(k)] = vList[0]
	}
	r.Headers = headers
	return r
}

func response(code int, msg string, data interface{}) *S3Response {
	r := &S3Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	return r
}
