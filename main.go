package main

import (
	"fmt"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lipeining/fakesql/config"
	"github.com/lipeining/fakesql/database"
)
var tblName string
var num int
var jsonPath string

func init() {
    const (
        // defaultNum = 1000000
        defaultNum = 10000
    )
    flag.StringVar(&tblName, "tblName", "user", "tblName")
    flag.StringVar(&jsonPath, "jsonPath", "./tables/user.json", "jsonPath absolute or relative path")
    flag.IntVar(&num, "num", defaultNum, "generate of num rows")
}
func main() {
	fmt.Println(config.Config)
	database.NewXorm(config.Config.Xorm.User, config.Config.Xorm.Passwd, config.Config.Xorm.Database, config.Config.Xorm.SecurePivFile)
	r := setupRouter()
	r.Run(":8000")
	flag.Parse()
	if tblName == "" || jsonPath == "" {
		fmt.Println("missed tblName || jsonPath")
	} else {
		fmt.Println(tblName, num, jsonPath)
		cols, err := ParseJSONColumn(jsonPath)
		if err != nil {
			fmt.Println("parse json file error ", err)
			return
		}
		makeColumn := MakeColumnFuncFactory(cols)
		err = Writetxt(tblName, num, makeColumn)
		if err != nil {
			fmt.Println("wrtite txt file error ", err)
			return
		}
		err = CreateTable(tblName, num, cols)
		if err != nil {
			fmt.Println("create table error ", err)
			return
		}
	}
}
