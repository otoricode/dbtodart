package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gobeam/stringy"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}
	var table string
	_db, _ := db.DB()

	res, err := _db.Query("SHOW TABLES")

	if err != nil {
		panic(err)
	}
	var tables []Table
	// var tables []string
	for res.Next() {
		res.Scan(&table)
		// getTable(db, table)
		tables = append(tables, getTable(db, table))
	}
	for i := 0; i < len(tables); i++ {
		tabelToClass(tables[i])
	}
	// fmt.Println(tables)
}

func getTable(db *gorm.DB, tableName string) Table {
	query := `SELECT 	COLUMN_NAME, 
	DATA_TYPE, 
	IS_NULLABLE, 
	COLUMN_KEY
	FROM 		INFORMATION_SCHEMA.COLUMNS 
	WHERE 		TABLE_SCHEMA = 'sagara_kosmedis' 
	AND 		TABLE_NAME = ?`
	var colum []*Column
	if err := db.Raw(query, tableName).Find(&colum).Error; err != nil {
		panic(err)
	}
	for i := 0; i < len(colum); i++ {
		colum[i].DataType = dataTypeConvert(colum[i].DataType, colum[i].IsNullable)
		// println(colum[i].DataType, colum[i].ColumnName)
	}
	return Table{
		Name:    tableName,
		Columns: colum,
	}
}

type Column struct {
	TableName  string `gorm:"column:TABLE_NAME" json:"table_name,omitempty"`
	ColumnName string `gorm:"column:COLUMN_NAME" json:"column_name,omitempty"`
	DataType   string `gorm:"column:DATA_TYPE" json:"data_type,omitempty"`
	IsNullable string `gorm:"column:IS_NULLABLE" json:"is_nullable,omitempty"`
	ColumnKey  string `gorm:"column:COLUMN_KEY" json:"column_key,omitempty"`
}

type Table struct {
	Name    string    `json:"name,omitempty"`
	Columns []*Column `json:"columns,omitempty"`
}

func dataTypeConvert(in1, in2 string) (out string) {
	nullable := ""
	if in2 == "YES" {
		nullable = "?"
	}
	switch in1 {
	case "bigint":
		return "int" + nullable
	case "char":
		return "String" + nullable
	case "date":
		return "DateTime" + nullable
	case "datetime":
		return "DateTime" + nullable
	case "double":
		return "double" + nullable
	case "int":
		return "int" + nullable
	case "longtext":
		return "String" + nullable
	case "mediumtext":
		return "String" + nullable
	case "text":
		return "String" + nullable
	case "timestamp":
		return "DateTime" + nullable
	case "tinyint":
		return "int" + nullable
	case "varchar":
		return "String" + nullable
	}
	return
}

func tabelToClass(in Table) string {
	str := stringy.New(in.Name)
	name := str.CamelCase()

	s := fmt.Sprintf("class %v {\n", name)

	for i := 0; i < len(in.Columns); i++ {
		str2 := stringy.New(in.Columns[i].ColumnName)
		str2 = stringy.New(str2.CamelCase())
		in.Columns[i].ColumnName = str2.LcFirst()
		s += fmt.Sprintln("  final", in.Columns[i].DataType, in.Columns[i].ColumnName+";")
	}
	s += "\n   " + name + "({"
	for i := 0; i < len(in.Columns); i++ {
		var required = ""
		if in.Columns[i].IsNullable == "NO" {
			required = "required "
		}
		s += fmt.Sprint("\n      " + required + "this." + in.Columns[i].ColumnName + ",")
	}
	s += "\n   });\n"
	s += fmt.Sprintln("}")
	println(s)
	os.WriteFile("./out/"+in.Name+".dart", []byte(s), 0644)
	return s
}
