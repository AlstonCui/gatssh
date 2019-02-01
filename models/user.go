package models

import (
	"time"
	"github.com/satori/go.uuid"
	"gatssh/utils"
)

type User struct {
	Id            int
	Uid           string `gorm:";type:varchar(100);not null;primary_key;"`
	Username      string `gorm:";type:varchar(100);not null;primary_key;"`
	Password      string `gorm:";type:varchar(200);not null;"`
	Rank          int    `gorm:";type:int(8);not null;"`
	Group         string `gorm:";type:varchar(200);not null;"`
	CreateTime    time.Time
	LastLoginTime time.Time
	PrivateKey    string `gorm:"type:TEXT;"`
	PublicKey     string `gorm:"type:TEXT;"`
}

func init() {
	var user User
	rows := db.First(&user)
	if rows.RowsAffected == 0 {
		uid := uuid.NewV4().String()
		user := &User{Username: "admin", Uid: uid, Password: utils.Md5Sum("123456"), Rank: 1000, Group: "admin", CreateTime: time.Now()}
		db.Create(&user)
	}
}

func GetUid(username string) string {
	var user User
	db.Where("username =?", username).Find(&user)
	return user.Uid
}

func (this *User) AuthUserAndPass() (uid string, ok bool) {

	var user User

	pass := utils.Md5Sum(this.Password)

	row := db.Where("username = ? AND password = ?", this.Username, pass).Find(&user)

	if row.RowsAffected != 0 {
		ok = true
		uid = user.Uid
	} else {
		ok = false
	}

	return
}

func (this *User) UpdatePassword(uid string, newPass string) (err error) {

	err = db.Model(this).Where("uid = ?", uid).Update("password", utils.Md5Sum(newPass)).Error
	if err != nil {
		return
	}

	return
}
