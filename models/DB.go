package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "gatlin.db")
	if err != nil {
		log.Fatal("conn",err)
	}

	db.AutoMigrate(&Host{}).AddUniqueIndex("idx_ip_owner", "ip", "owner")
	db.AutoMigrate(&TaskRecord{})
	db.AutoMigrate(&TaskDetail{})

}
