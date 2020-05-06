package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lipeining/fakesql/config"
	"github.com/lipeining/fakesql/database"
)
func main() {
	fmt.Println(config.Config)
	database.NewXorm(config.Config.Xorm.User, config.Config.Xorm.Passwd, config.Config.Xorm.Database, config.Config.Xorm.SecurePivFile)
	database.NewGorm(config.Config.Gorm.Mysql)
}
