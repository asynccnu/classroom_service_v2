package model

import (
	"time"
)

var DBName = "classroom_service"

const (
	ClassroomCol = "classrooms"
)

// 数据库表的结构
type Classrooms struct {
	Week                int    `bson:"week"`
	Weekday             int    `bson:"weekday"`
	Building            string `bson:"building"`
	AvailableClassrooms Rooms  `bson:"available_classrooms"`
}

// 参照上一版接口文档写的
type Rooms struct {
	One    []string `bson:"one"`
	Two    []string `bson:"two"`
	Three  []string `bson:"three"`
	Four   []string `bson:"four"`
	Five   []string `bson:"five"`
	Six    []string `bson:"six"`
	Seven  []string `bson:"seven"`
	Eight  []string `bson:"eight"`
	Nine   []string `bson:"nine"`
	Ten    []string `bson:"ten"`
	Eleven []string `bson:"eleven"`
	Twelve []string `bson:"twelve"`
}

type BaseModel struct {
	Id        uint64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"-"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"-"`
	DeletedAt *time.Time `gorm:"column:deletedAt" sql:"index" json:"-"`
}

// Token represents a JSON web token.
type Token struct {
	Token string `json:"token"`
}
