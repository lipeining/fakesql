package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lipeining/fakesql/model"
)

// TableForm use to bind post /table form
type TableForm struct {
	TblName string `form:"tblName" binding:"required"`
	// JSONPath   string `json:"jsonPath" form:"jsonPath"`
	Cols string `form:"cols" binding:"required"`
	Num  int    `form:"num" binding:"required"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "home")
	})
	r.POST("/tables", func(c *gin.Context) {
		// var tableForm TableForm
		// if err := c.ShouldBind(&tableForm); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		// tblName, num, defCols := tableForm.TblName, tableForm.Num, tableForm.Cols
		// cols, err := ParseJSONColumn(tblName)
		// if err != nil {
		// 	fmt.Println("parse json file error ", err)
		// 	c.JSON(http.StatusBadRequest, err)
		// 	return
		// }
		// makeColumn := MakeColumnFuncFactory(cols)
		tblName := c.PostForm("tblName")
		num, _ := strconv.Atoi(c.DefaultPostForm("Num", "10000"))
		colsStr := c.PostForm("cols")
		var cols []model.Column
		err := json.Unmarshal(bytes.NewBufferString(colsStr).Bytes(), &cols)
		if err != nil {
			fmt.Println("un marshal cols error ", err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		fmt.Println(tblName, num, cols)
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
