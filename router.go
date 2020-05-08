package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)
// TableForm use to bind post /table form
type TableForm struct {
	TblName    string `json:"tblName" form:"tblName"`
	// JSONPath   string `json:"jsonPath" form:"jsonPath"`
	Cols   string `json:"cols" form:"cols"`
	Num        int    `json:"num" form:"num"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context){
		c.String(http.StatusOK, "pong")
	})
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", func(c *gin.Context){
		c.String(http.StatusOK, "home")
	})
	r.POST("/tables", func(c *gin.Context){
		var tableForm TableForm
		if err := c.ShouldBindJSON(&tableForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tblName, num, defCols := tableForm.TblName, tableForm.Num, tableForm.Cols
		fmt.Println(tblName, num, defCols)
		cols, err := ParseJSONColumn(tblName)
		if err != nil {
			fmt.Println("parse json file error ", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		makeColumn := MakeColumnFuncFactory(cols)
		err = Writetxt(tblName, num, makeColumn)
		if err != nil {
			fmt.Println("wrtite txt file error ", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		err = CreateTable(tblName, num, cols)
		if err != nil {
			fmt.Println("create table error ", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.String(http.StatusOK, "ok")
	})
	return r
}