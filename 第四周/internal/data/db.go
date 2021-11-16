package data

import (
	"github.com/jinzhu/gorm"
)

func NewDB() *gorm.DB {
	// 从配置文件中读取信息，连接数据库
	//db, err := gorm.Open("")
	//if err != nil {
	//	panic("failed to connect database")
	//}
	return &gorm.DB{}
}