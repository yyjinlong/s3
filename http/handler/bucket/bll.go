// Copyright @ 2020 OPS Inc.
//
// Author: Jinlong Yang
//

package bucket

import (
	"fmt"
	"sync"

	"s3/http/handler"
	"s3/ops/objects"
)

var (
	BkList *BucketList
	BkPut  *BucketPut
	BkDel  *BucketDel
	once   sync.Once
)

func init() {
	once.Do(func() {
		BkList = &BucketList{}
		BkPut = &BucketPut{}
		BkDel = &BucketDel{}
	})
}

type BucketList struct {
}

func (bl *BucketList) Listen(r *handler.S3Request) (interface{}, error) {
	bucketList, err := objects.GetBucketList(r.UID)
	if err != nil {
		return nil, err
	}
	userInfo, err := objects.GetUserInfo(r.UID)
	if err != nil {
		return nil, err
	}

	ownerInfo := make(map[string]interface{})
	ownerInfo["uid"] = r.UID
	ownerInfo["name"] = userInfo.Name

	bucketArray := make([]map[string]string, 0)
	for _, item := range bucketList {
		bucketInfo := map[string]string{
			"name":        item.Name,
			"create_time": handler.ConvertTime(item.CreateAt),
		}
		bucketArray = append(bucketArray, bucketInfo)
	}

	result := make(map[string]interface{})
	result["owner"] = ownerInfo
	result["buckets"] = bucketArray
	return result, nil
}

type BucketPut struct {
}

func (bp *BucketPut) Listen(r *handler.S3Request) (interface{}, error) {
	// TODO: 数据录入
	// TODO: 调用proxy-server接口
	fmt.Println("------put bucket")
	return nil, nil
}

type BucketDel struct {
}

func (bd *BucketDel) Listen(r *handler.S3Request) (interface{}, error) {
	// TODO: 数据软删除
	// TODO: 调用proxy-server接口
	fmt.Println("------del bucket")
	return nil, nil
}
