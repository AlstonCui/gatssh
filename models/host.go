package models

import (
	"time"
	"errors"
)

type Host struct {
	Id                int
	Ip                string `gorm:";type:varchar(100);not null;"`
	Port              int    `gorm:";type:int(8);not null;"`
	User              string `gorm:";type:varchar(100);not null;"`
	Password          string `gorm:"type:varchar(200);not null;"`
	Owner             string `gorm:"type:varchar(100);not null;"`
	LastOperationTime time.Time
}

func (this *Host) SaveHost() (err error) {

	var host Host
	rows := db.Where(&Host{Ip: this.Ip, Port: this.Port, Owner: this.Owner}).Find(&host)
	if rows.RowsAffected != 0 {
		err = db.Model(&host).Updates(map[string]interface{}{"password": this.Password, "last_operation_time": time.Now()}).Error
		if err != nil {
			return
		}
		return
	}

	err = db.Create(this).Error
	if err != nil {
		return
	}

	return
}

func QueryHost(ip string, port int, gatUser string) (user string, password string, err error) {

	var host Host
	rows := db.Where(&Host{Ip: ip, Port: port, Owner: gatUser}).Find(&host)
	if rows.RowsAffected == 0 {
		err = errors.New("No match password in DB,Please make sure the password has been saved...")
		return
	}
	user = host.User
	password = host.Password

	return
}
