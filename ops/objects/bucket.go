// Copyright @ 2020 OPS Inc.
//
// Author: Jinlong Yang
//

package objects

import (
	log "github.com/sirupsen/logrus"

	"s3/ops/db"
)

func GetBucketList(uid int64) ([]db.Bucket, error) {
	engine, err := db.GetSlaveEngine()
	if err != nil {
		return nil, err
	}

	bucketList := make([]db.Bucket, 0)
	if err := engine.Where("user_id=?", uid).Find(&bucketList); err != nil {
		log.Errorf("Query bucket list by uid: %d error: %s", uid, err)
		return nil, err
	}
	return bucketList, nil
}
