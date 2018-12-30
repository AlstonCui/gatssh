package models

import (
	"time"
)

type User struct {
	Id            int
	Uid          string
	Username      string
	Password      string
	Rank          int
	Group         string
	CreateTime    time.Time
	LastLoginTime time.Time
}

