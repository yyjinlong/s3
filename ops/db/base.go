// Copyright @ 2020 OPS Inc.
//
// Author: Jinlong Yang
//

package db

import (
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	"s3/http/g"
)

func GetMasterEngine() (*xorm.Engine, error) {
	if eg, err := connect(); err != nil {
		return nil, err
	} else {
		return eg.Master(), nil
	}
}

func GetSlaveEngine() (*xorm.Engine, error) {
	if eg, err := connect(); err != nil {
		return nil, err
	} else {
		return eg.Slave(), nil
	}
}

func connect() (*xorm.EngineGroup, error) {
	master, err := xorm.NewEngine("postgres", g.Config().Postgres.Master)
	if err != nil {
		log.Errorf("Connect master database error: %s", err)
		return nil, err
	}

	slave1, err := xorm.NewEngine("postgres", g.Config().Postgres.Slave1)
	if err != nil {
		log.Errorf("Connect slave1 database error: %s", err)
		return nil, err
	}

	slave2, err := xorm.NewEngine("postgres", g.Config().Postgres.Slave2)
	if err != nil {
		log.Errorf("Connect slave2 database error: %s", err)
		return nil, err
	}

	slaves := []*xorm.Engine{slave1, slave2}
	eg, err := xorm.NewEngineGroup(master, slaves)
	if err != nil {
		log.Errorf("Create engine group failed: %s", err)
		return nil, err
	}

	if err := eg.Ping(); err != nil {
		log.Errorf("Ping connect error: %s", err)
		return nil, err
	}

	eg.ShowSQL(true)
	eg.SetMapper(names.GonicMapper{})
	return eg, nil
}
