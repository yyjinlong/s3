// Copyright @ 2020 OPS Inc.
//
// Author: Jinlong Yang
//

package objects

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"s3/ops/db"
)

func GetUserInfo(uid int64) (*db.User, error) {
	engine, err := db.GetSlaveEngine()
	if err != nil {
		return nil, err
	}
	user := new(db.User)
	if has, err := engine.ID(uid).Get(user); err != nil {
		log.Errorf("Query user id: %d info error: %s", uid, err)
		return nil, err
	} else if !has {
		log.Errorf("Query uid: %d from user is not existed.", uid)
		return nil, fmt.Errorf("no data found.")
	}
	return user, nil
}

func GetAuthMetaByAk(ak string) (*db.AuthMeta, error) {
	engine, err := db.GetSlaveEngine()
	if err != nil {
		return nil, err
	}
	authMeta := new(db.AuthMeta)
	if has, err := engine.Where("ak=?", ak).Get(authMeta); err != nil {
		log.Errorf("Query auth_meta by ak: %s error: %s", ak, err)
		return nil, err
	} else if !has {
		log.Errorf("Query ak: %s from auth_meta is not existed.", ak)
		return nil, fmt.Errorf("no data found.")
	}
	return authMeta, nil
}
