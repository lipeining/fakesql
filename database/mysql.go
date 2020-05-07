package database

import (
	"fmt"

	"xorm.io/xorm"
)

// Xorm a pointer to xorm.Engine
var Xorm *xorm.Engine

// const user string = "root"
// const passwd string = "root"
// const database string = "test"
// const securePivFile string = "C:/ProgramData/MySQL/MySQL Server 5.7/Uploads/base.txt"

// NewXorm init global DB
func NewXorm(user, passwd, database, securePivFile string) {
	var err error
	url := user + ":" + passwd + "@" + "/" + database
	Xorm, err = xorm.NewEngine("mysql", url)
	if err != nil {
		fmt.Println(err)
	}
	Xorm.ShowSQL(true)
}

// // Gorm a pointer to gorm.DB
// var Gorm  *gorm.DB
// // NewGorm init global DB
// func NewGorm(mysql string) {
// 	var err error
// 	Gorm, err := gorm.Open("mysql", mysql)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	Gorm.LogMode(true)
// }
