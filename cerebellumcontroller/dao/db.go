package dao

import (
	"errors"
	"iai/cerebellumController/config"
	"iai/cerebellumController/dao/model"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

func Init() {
	var err error
	DBConn, err = gorm.Open(sqlite.Open(config.Config.DB.File), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Panicln(err)
	}
	DBConn.AutoMigrate(&model.RouteItem{})
}

func GetRoute(routeid string) (model.RouteItem, error) {
	var result model.RouteItem
	r := DBConn.Where(&model.RouteItem{RouteId: routeid}).First(&result)
	return result, r.Error
}

func AddRoute(item model.RouteItem) error {
	item.RouteId = item.Dst
	return DBConn.Create(&item).Error

}

func GetRouteList(taskId string) (routeList []model.RouteItem, err error) {
	err = DBConn.Where(&model.RouteItem{Taskid: taskId}).Find(&routeList).Error
	return
}

func ModifyRoute(item model.RouteItem) (err error) {
	var tmp model.RouteItem
	row := DBConn.Where(&model.RouteItem{RouteId: item.RouteId}).First(&tmp).RowsAffected
	if row != 1 {
		return errors.New("no such route")
	} else {
		tmp.Dst = item.Dst
		tmp.Src = item.Dst
		tmp.Svc = item.Svc
		tmp.Node = item.Node
		tmp.Taskid = item.Taskid
		return DBConn.Save(&tmp).Error
	}
}

func DeleteRoute(routeid string) {
	DBConn.Where(&model.RouteItem{RouteId: routeid}).Delete(&model.RouteItem{})
}

func GetAllRouteList() (routeList []model.RouteItem, err error) {
	err = DBConn.Find(&routeList).Error
	return
}
