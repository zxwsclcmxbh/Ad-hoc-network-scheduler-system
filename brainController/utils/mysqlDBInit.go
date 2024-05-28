package utils

import (
	"cloud/brainController/model"

	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func MysqlDBInit() {
	// dsn := "root:Aa123456!@tcp(127.0.0.1:3306)/brain_controller?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := Config.DB
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = DB.AutoMigrate(&model.Task{}, &model.IaiImage{}, &model.Pod{}, &model.QueueRoute{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("MysqlDBInit successfully")
}
