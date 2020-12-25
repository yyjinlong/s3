// Copyright @ 2020 OPS Inc.
//
// Author: Jinlong Yang
//

package db

import (
	"time"
)

type User struct {
	ID       int64
	Name     string    `xorm:"varchar(128) notnull"`
	Status   int       `xorm:"int notnull"`
	CreateAt time.Time `xorm:"timestamp notnull created"`
	UpdateAt time.Time `xorm:"timestamp notnull updated"`
}

type AuthMeta struct {
	ID       int64
	UserID   int64     `xorm:"bigint notnull"`
	AK       string    `xorm:"varchar(32) notnull"`
	SK       string    `xorm:"varchar(32) notnull"`
	IsValid  bool      `xorm:"bool"`
	CreateAt time.Time `xorm:"timestamp notnull created"`
	UpdateAt time.Time `xorm:"timestamp notnull updated"`
}

type Bucket struct {
	ID            int64
	UserID        int64     `xorm:"bigint notnull"`            // 用户id
	SlaveBucketID int64     `xorm:"bigint not null"`           // 备份桶id
	Name          string    `xorm:"varchar(64) notnull"`       // 桶名称
	Status        int       `xorm:"int notnull"`               // 0: 有效 1: 无效
	IsEncrypted   int       `xorm:"int notnull"`               // 0: 不开启 1: 开启
	ACLType       int       `xorm:"int notnull"`               // 0: private 1: public-read 2: public-read-write
	ObjCacheOpen  int       `xorm:"int notnull"`               // 0: 不开启 1: 开启
	CreateAt      time.Time `xorm:"timestamp notnull created"` // 创建时间
	UpdateAt      time.Time `xorm:"timestamp notnull updated"` // 更新时间
}

type Meta struct {
	ID             int64
	UserID         int64     `xorm:"bigint notnull"`
	BucketID       int64     `xorm:"bigint notnull"`
	Object         string    `xorm:"varchar(512) notnull"`
	Version        string    `xorm:"varchar(32) notnull"`
	Size           int64     `xorm:"bigint notnull"`
	UserMeta       string    `xorm:"varchar(512)"`
	ReferObject    string    `xorm:"varchar(512)"`
	RefertVersion  string    `xorm:"varchar(32)"`
	GrantedRUsers  string    `xorm:"varchar(256) notnull"`
	GrantedRWUsers string    `xorm:"varchar(256) notnull"`
	IsDeleted      bool      `xorm:"bool notnull"`
	CreateAt       time.Time `xorm:"timestamp notnull created"`
	UpdateAt       time.Time `xorm:"timestamp notnull updated"`
	DeleteAt       time.Time `xorm:"timestamp"`
}

type Recoup struct {
	ID            int64
	BucketName    string    `xorm:"varchar(64)"`
	ObjectName    string    `xorm:"varchar(512)"`
	ObjectVersion string    `xorm:"varchar(32)"`
	UserName      string    `xorm:"varchar(128)"`
	Status        int       `xorm:"int"`
	RetryNum      int       `xorm:"int"`
	CreateAt      time.Time `xorm:"timestamp notnull created"`
	UpdateAt      time.Time `xorm:"timestamp notnull updated"`
	FinishAt      time.Time `xorm:"timestamp"`
}
