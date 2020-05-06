package model

// TblTmpID 临时表结构
type TblTmpID struct {
	ID int64 `xorm:"pk notnull"`
}

// Column 输入列属性
type Column struct {
	Name string
	T    string
	Def  string
}
