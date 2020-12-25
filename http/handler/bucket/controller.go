// Copyright @ 2020 OPS Inc.
//
// Author: Jinlong Yang
//

package bucket

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"s3/http/handler"
)

func GetBucket(c *gin.Context) {
	c.JSON(http.StatusOK, handler.Talk(BkList, c))
}

func PutBucket(c *gin.Context) {
	c.JSON(http.StatusOK, handler.Talk(BkPut, c))
}

func DelBucket(c *gin.Context) {
	c.JSON(http.StatusOK, handler.Talk(BkDel, c))
}
