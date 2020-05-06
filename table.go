package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"path/filepath"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/huandu/go-sqlbuilder"
	"github.com/lipeining/fakesql/config"
	"github.com/lipeining/fakesql/database"
	"github.com/lipeining/fakesql/model"
)
func init() {
	gofakeit.Seed(0)
}

// LOWPRORITY use by load data
const LOWPRORITY bool = false

// CONCURRENT use by load data
const CONCURRENT bool = false
// LOCAL use by load data
const LOCAL bool = true

// Chareset use by load data
const Chareset string = "utf8mb4"
// FieldsTerminatedBy use by load data
const FieldsTerminatedBy string = ","
// FieldsEnclosedBy use by load data
const FieldsEnclosedBy string = ""
// FieldsEscapedBy use by load data
const FieldsEscapedBy string = "\\"

// LinesTerminatedBy use by load data
const LinesTerminatedBy string = "\n"
// const LinesStartingBy string = '\n'

// MakeUUIDString make a select column of uuid with length
func MakeUUIDString(length int) string {
	if length > 0 {
		return "substring(uuid(), 1, " + strconv.Itoa(length) + ")"
	}
	return "uuid()"
}
// MakeIntString make a int number between min and max
func MakeIntString(min, max int) string {
	return strconv.Itoa(min) + " + " + "FLOOR(Rand() * " + strconv.Itoa(max*10) + ")";
}
// MakeDateString make a date
func MakeDateString(startDate string) string {
	str := "date_add("
	if startDate != "" {
		str += startDate
	} else {
		str += "NOW()"
	}
	return str+ ", interval FLOOR(1+(RAND()*10)) month)"
}
// MakeChineseNameString make a name of chinese
func MakeChineseNameString() (string, string) {
	base := "set @SURNAME = '王李张刘陈杨黄赵吴周徐孙马朱胡郭何高林罗郑梁谢宋唐位许韩冯邓曹彭曾萧田董潘袁于蒋蔡余杜叶程苏魏吕丁任沈姚卢姜崔钟谭陆汪范金石廖贾夏韦傅方白邹孟熊秦邱江尹薛阎段雷侯龙史陶黎贺顾毛郝龚邵万钱严覃武戴莫孔向汤'; set @NAME = '丹举义之乐书乾云亦从代以伟佑俊修健傲儿元光兰冬冰冷凌凝凡凯初力勤千卉半华南博又友同向君听和哲嘉国坚城夏夜天奇奥如妙子存季孤宇安宛宸寒寻尔尧山岚峻巧平幼康建开弘强彤彦彬彭心忆志念怀怜恨惜慕成擎敏文新旋旭昊明易昕映春昱晋晓晗晟景晴智曼朋朗杰松枫柏柔柳格桃梦楷槐正水沛波泽洁洋济浦浩海涛润涵渊源溥濮瀚灵灿炎烟烨然煊煜熙熠玉珊珍理琪琴瑜瑞瑶瑾璞痴皓盼真睿碧磊祥祺秉程立竹笑紫绍经绿群翠翰致航良芙芷苍苑若茂荣莲菡菱萱蓉蓝蕊蕾薇蝶觅访诚语谷豪赋超越轩辉达远邃醉金鑫锦问雁雅雨雪霖霜露青靖静风飞香驰骞高鸿鹏鹤黎';"	
	return base, "concat(substr(@surname,floor(rand()*length(@surname)/3+1),1), substr(@NAME,floor(rand()*length(@NAME)/3+1),1), substr(@NAME,floor(rand()*length(@NAME)/3+1),1))"
}
// MakeBaseString make a string with base string, prepare id from table
func MakeBaseString(base string) string {
	return "CONCAT('"+base+"', id)";
}
// PathExists use os.stat to check path
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
// MakeTmpTableName general make the table name
func MakeTmpTableName(num int) string {
	return "base_" + strconv.Itoa(num)
}
// WriteBaseFile create a base id file
func WriteBaseFile(num int) error {
	fileName := MakeTmpTableName(num) + ".txt"
	filePath := filepath.Join(config.Config.Xorm.SecurePivFile, fileName)
	exists, err := PathExists(filePath)
	if err != nil {
		fmt.Println("An error stat with file \n", filePath, err)
		return err
	}
	if exists {
		return nil
	}
	// 不存在时，创建文件
	outputFile, outputError := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Println("An error occurred with file opening or creation\n", outputError)
		return outputError
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	// 每次写 10000 个数字
	every := 10000
	total := num / every
	if num%every != 0 {
		total++
	}
	current := 0
	for i := 0; i < total; i++ {
		for j := 0; j < every; j++ {
			current++
			outputString := strconv.Itoa(current)
			if current != 1 {
				outputString = "\n" + outputString
			}
			outputWriter.WriteString(outputString)
		}
		outputWriter.Flush()
	}
	return nil
}
// CreateBaseTable create a base number table
func CreateBaseTable(num int) error {
	tmpTableName := MakeTmpTableName(num)
	ctb := sqlbuilder.NewCreateTableBuilder()
	fullTmpTableName := config.Config.Xorm.Database + "." + tmpTableName
	ctb.CreateTable(fullTmpTableName).IfNotExists()
	ctb.Define("id", "BIGINT(20)", "NOT NULL", "PRIMARY KEY", `COMMENT "id"`)
	ctb.Option("DEFAULT CHARACTER SET", "utf8mb4")
	fmt.Println(ctb)
	insertCommand := ctb.String()
	results, err := database.Xorm.Query(insertCommand)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("create table results", results)
	return nil
}
// LoadIntoBaseTable load file data into tbl table
func LoadIntoBaseTable(num int) error {
	tmpTableName := MakeTmpTableName(num)
	fileName := tmpTableName + ".txt"
	fullTmpTableName := config.Config.Xorm.Database + "." + tmpTableName
	filePath := config.Config.Xorm.SecurePivFile + "/" + fileName
	loadDataCommand := "load data infile " + "'" + filePath + "'" + " replace into table " + fullTmpTableName + ";"
	results, err := database.Xorm.Query(loadDataCommand)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("load data results", results)
	return nil
}
// CreateTableByloadFile use cols and load file table to create table
func CreateTableByloadFile(tblName string, num int, cols []model.Column) error {
	tmpTableName := MakeTmpTableName(num)
	fullTmpTableName := config.Config.Xorm.Database + "." + tmpTableName
	fullTblName := config.Config.Xorm.Database + "." + tblName
	insertList := make([]string, 0)
	insertList = append(insertList, "INSERT INTO "+fullTblName+" SELECT ")
	insertCols := make([]string, 0)
	insertCols = append(insertCols, "id")
	ctb := sqlbuilder.NewCreateTableBuilder()
	ctb.CreateTable(fullTblName).IfNotExists()
	ctb.Define("id", "BIGINT(20)", "NOT NULL", "PRIMARY KEY", "AUTO_INCREMENT", `COMMENT "id"`)
	// todo 需要检查正确性
	// todo 丰富的 column 属性 默认值，大小，
	for _, column := range cols {
		name := column.Name
		if column.T == "string" {
			insertCols = append(insertCols, MakeUUIDString(0))
			ctb.Define(name, "VARCHAR(255)", "NOT NULL", `COMMENT "`+name+`"`)
		} else if column.T == "datetime" {
			insertCols = append(insertCols, MakeDateString(""))
			ctb.Define(name, "datetime", "NOT NULL", `COMMENT "`+name+`"`)
		} else if column.T == "int64" {
			insertCols = append(insertCols, MakeIntString(1, 1000))
			ctb.Define(name, "BIGINT(20)", "NOT NULL", `COMMENT "`+name+`"`)
		} else if column.T == "int32" {
			insertCols = append(insertCols, MakeIntString(1, 1000))
			ctb.Define(name, "INT(11)", "NOT NULL", `COMMENT "`+name+`"`)
		}
	}
	insertList = append(insertList, strings.Join(insertCols, ","))
	insertList = append(insertList, " FROM "+fullTmpTableName+";")
	ctb.Option("DEFAULT CHARACTER SET", "utf8mb4")
	fmt.Println(ctb)
	createTableCommand := ctb.String()
	insertDataCommand := strings.Join(insertList, "")
	results, err := database.Xorm.Query(createTableCommand)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("create table ", tblName, results)
	results, err = database.Xorm.Query(insertDataCommand)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("insert table ", tblName, results)
	return nil
}

// CreateTable 创建表格
func CreateTable(tblName string, num int, cols []model.Column) error {
	fullTblName := config.Config.Xorm.Database + "." + tblName
	fileName :=  tblName + "_" + strconv.Itoa(num) + ".txt"
	filePath := config.Config.Xorm.SecurePivFile + "/" + fileName
	loadDataCommand := "load data infile " +
	  "'" + filePath + "'" +
	  " replace into table " + fullTblName +
	  " character set " + Chareset +
	  " fields terminated by " + "'" + FieldsTerminatedBy + "' ";
	//   +	  " lines terminated by " + "'" + LinesTerminatedBy + "'";
	insertCols := make([]string, 0);
	ctb := sqlbuilder.NewCreateTableBuilder()
	ctb.CreateTable(fullTblName).IfNotExists()
	timeCols := make([]string, 0);
	// todo 需要检查正确性
	// todo 丰富的 column 属性 默认值，大小，
	for _, column := range cols {
		name := column.Name
		if column.T == "string" {
			insertCols = append(insertCols, "`" + name + "`")
			ctb.Define(name, "VARCHAR(255)", "NOT NULL", `COMMENT "`+name+`"`)
		} else if column.T == "datetime" {
			insertCols = append(insertCols, "`" + name + "`")
			// insertCols = append(insertCols, "@" + name)
			// timeCols = append(timeCols, name)
			ctb.Define(name, "datetime", "NOT NULL", `COMMENT "`+name+`"`)
		} else if column.T == "int64" {
			insertCols = append(insertCols, "`" + name + "`")
			ctb.Define(name, "BIGINT(20)", "NOT NULL", `COMMENT "`+name+`"`)
		} else if column.T == "int32" {
			insertCols = append(insertCols, "`" + name + "`")
			ctb.Define(name, "INT(11)", "NOT NULL", `COMMENT "`+name+`"`)
		} else if column.T == "id" {
			insertCols = append(insertCols, "`" + name + "`")
			ctb.Define("id", "BIGINT(20)", "NOT NULL", "PRIMARY KEY", "AUTO_INCREMENT", `COMMENT "id"`)
		}
	}
	ctb.Option("DEFAULT CHARACTER SET", "utf8mb4")
	fmt.Println(ctb)
	createTableCommand := ctb.String()
	// load data command 暂时不处理 set 操作
	results, err := database.Xorm.Query(createTableCommand)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// fmt.Println("create table results", results)
	loadDataCommand += "("+strings.Join(insertCols, ",")+")"
	// // 处理 time column
	// 如果是标准的 %Y-%m-%d %H:%i:%s ， 不需要 set 操作
	setStr := ""
	for _, col := range timeCols {
		setStr += " set " + col + " = STR_TO_DATE(@"+col+", '%Y-%m-%d %H:%i:%s') "
	}
	loadDataCommand += setStr
	loadDataCommand += ";"
	fmt.Println(loadDataCommand)
	results, err = database.Xorm.Query(loadDataCommand)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("load table results", results)
	return nil
}
// MakeColumnFuncFactory 通用生成一行数据的回调函数
func MakeColumnFuncFactory(cols []model.Column) func(int) []string {
	return func(current int) []string {
		// todo 需要检查正确性
		// todo 丰富的 column 属性 默认值，大小，
		insertCols := make([]string, 0)
		for _, column := range cols {
			if column.T == "string" {
				insertCols = append(insertCols, gofakeit.Word())
			} else if column.T == "datetime" {
				d := gofakeit.Date().Format("2006-01-02 15:04:05")
				insertCols = append(insertCols, d)
			} else if column.T == "int64" {
				insertCols = append(insertCols, strconv.FormatInt(int64(gofakeit.Uint32()), 10))
			} else if column.T == "int32" {
				insertCols = append(insertCols, strconv.FormatInt(int64(gofakeit.Uint32()), 10))
			} else if column.T == "id" {
				insertCols = append(insertCols, strconv.Itoa(current))
			}
		}
		return insertCols
	}
}
// Writetxt write 1,000,000
func Writetxt(tblName string, num int, makeColumn func(int) []string ) error {
	fileName :=  tblName + "_" + strconv.Itoa(num) + ".txt"
	filePath := filepath.Join(config.Config.Xorm.SecurePivFile, fileName)
	exists, err := PathExists(filePath)
	if err != nil {
		fmt.Println("An error stat with file \n", filePath, err)
		return err
	}
	if exists {
		return nil
	}
	outputFile, outputError := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Println("An error occurred with file opening or creation\n", filePath, outputError)
		return outputError
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	// 每次写 10000 行
	every := 10000
	total := num / every
	if num%every != 0 {
		total++
	}
	current := 0
	for i := 0; i < total; i++ {
		for j := 0; j < every; j++ {
			current++
			outputString := strings.Join(makeColumn(current), FieldsTerminatedBy)
			if current != 1 {
				outputString = LinesTerminatedBy + outputString
			}
			outputWriter.WriteString(outputString)
		}
		outputWriter.Flush()
	}
	return nil
}