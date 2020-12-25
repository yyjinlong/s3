package main

import (
	"github.com/gin-gonic/gin"

	"s3/http/handler/bucket"
)

func urls(r *gin.Engine) {
	/* bucket */
	r.GET("/bucket", bucket.GetBucket)
	r.PUT("/bucket", bucket.PutBucket)
	r.DELETE("/bucket", bucket.DelBucket)

	/* object */

	/* acl */
}
