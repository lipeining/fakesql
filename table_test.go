package main

import (
	"testing"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lipeining/fakesql/config"
	"github.com/lipeining/fakesql/database"
	"github.com/lipeining/fakesql/model"
	"github.com/stretchr/testify/assert"
)
func setup() {
	// fmt.Println(config.Config)
	database.NewXorm(config.Config.Xorm.User, config.Config.Xorm.Passwd, config.Config.Xorm.Database, config.Config.Xorm.SecurePivFile)
}
func TestWriteBaseFile(t *testing.T) {
	setup()
	err := WriteBaseFile(20000)
	assert.Equal(t, nil, err)
}
func TestCreateBaseTable(t *testing.T) {
	setup()
	num := 20000
	err := CreateBaseTable(num)
	assert.Equal(t, nil, err)
	err = LoadIntoBaseTable(num)
	assert.Equal(t, nil, err)
	cols := []model.Column{
		model.Column{Name: "city", T: "string",},
		model.Column{Name: "province", T: "string",},
		model.Column{Name: "create_time", T: "datetime",},
	}
	tblName := "car"
	err = CreateTableByloadFile(tblName, num, cols)
	assert.Equal(t, nil, err)
}
func TestMakeColumnFuncFactory(t *testing.T) {
	cols := []model.Column{
		model.Column{Name: "id", T: "id",},
		model.Column{Name: "city", T: "string",},
		model.Column{Name: "province", T: "string",},
		model.Column{Name: "create_time", T: "datetime",},
	}
	makeColumn := MakeColumnFuncFactory(cols)
	oneLine := makeColumn(1)
	fmt.Println(oneLine)
	assert.Equal(t, true, len(oneLine)>0)
}
func TestWritetxt(t *testing.T) {
	setup()
	cols := []model.Column{
		model.Column{Name: "id", T: "id",},
		model.Column{Name: "city", T: "string",},
		model.Column{Name: "province", T: "string",},
		model.Column{Name: "create_time", T: "datetime",},
	}
	makeColumn := MakeColumnFuncFactory(cols)
	tblName := "computer"
	num := 10000
	err := Writetxt(tblName, num, makeColumn)
	assert.Equal(t, nil, err)
}
func TestCreateTable(t *testing.T) {
	setup()
	cols := []model.Column{
		model.Column{Name: "id", T: "id",},
		model.Column{Name: "city", T: "string",},
		model.Column{Name: "province", T: "string",},
		model.Column{Name: "create_time", T: "datetime",},
	}
	makeColumn := MakeColumnFuncFactory(cols)
	tblName := "book"
	num := 10000
	err := Writetxt(tblName, num, makeColumn)
	assert.Equal(t, nil, err)
	err = CreateTable(tblName, num, cols)
	assert.Equal(t, nil, err)
}
func TestCreateTableAndInsertSQL(t *testing.T) {
	setup()
	cols := []model.Column{
		model.Column{Name: "id", T: "id",},
		model.Column{Name: "city", T: "string",},
		model.Column{Name: "province", T: "string",},
		model.Column{Name: "create_time", T: "datetime",},
	}
	tblName := "fork"
	num := 10000
	err := CreateTableAndInsertSQL(tblName, num, cols)
	assert.Equal(t, nil, err)
}
func TestCreateTableAndInsertSQLFile(t *testing.T) {
	setup()
	cols := []model.Column{
		model.Column{Name: "id", T: "id",},
		model.Column{Name: "city", T: "string",},
		model.Column{Name: "province", T: "string",},
		model.Column{Name: "create_time", T: "datetime",},
	}
	tblName := "pen"
	num := 100000
	err := CreateTableAndInsertSQLFile(tblName, num, cols)
	assert.Equal(t, nil, err)
	// err = SourceSQL(tblName, num)
	// assert.Equal(t, nil, err)
}