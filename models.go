package main

type Table struct {
	TableName string   `json:"name,omitempty"`
	Columns   []Column `json:"columns,omitempty"`
	HasList   bool
	Relations []string
}

type Column struct {
	ColumnName string `gorm:"column:COLUMN_NAME" json:"column_name,omitempty"`
	DataType   string `gorm:"column:DATA_TYPE" json:"data_type,omitempty"`
	IsNullable string `gorm:"column:IS_NULLABLE" json:"is_nullable,omitempty"`
	ColumnKey  string `gorm:"column:COLUMN_KEY" json:"column_key,omitempty"`
	IsList     bool
}
type Class struct {
	Name       string
	Properties []ClassProperties
	Import     []string
	HasList    bool
}

type ClassProperties struct {
	VarName             string
	VarType             string
	IsNullable          bool
	Declaration         string
	ConstructorLine     string
	JsonKey             string
	FromMapLine         string
	CopyWithParamLine   string
	CopyWithPassingLine string
	ToMapLine           string
	ToStringLine        string
	EqualityLine        string
	HashCodeLine        string
}
