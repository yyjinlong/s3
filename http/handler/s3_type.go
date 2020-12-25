// Copyright @ 2020 OPS Inc.
//
// Author: Jinlong Yang
//

package handler

// 状态码
const (
	SUCCESS = 0
	FAILED  = 1
)

// 返回信息
const (
	ERROR = "success"
)

type S3Request struct {
	HTTPMethod string
	URL        string
	Query      string
	Headers    map[string]string
	RequestID  string
	UID        int64
}

type S3Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
