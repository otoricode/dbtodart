package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	var classes []Class
	for res.Next() {
		res.Scan(&table)
		tables = append(tables, getTable(db, database, table))
	}
	for i := 0; i < len(tables); i++ {
		classes = append(classes, tabelToClass(tables[i]))
	}
	for i, v := range classes {
		os.WriteFile("./out/"+tables[i].TableName+".dart", []byte(classToDart(v)), 0644)
	}
}

func getTable(db *gorm.DB, database, tableName string) Table {
	query := `SELECT 	COLUMN_NAME, 
	DATA_TYPE, 
	IS_NULLABLE, 
	COLUMN_KEY
	FROM 		INFORMATION_SCHEMA.COLUMNS 
	WHERE 		TABLE_SCHEMA = ? 
	AND 		TABLE_NAME = ?`
	var colum []Column
	if err := db.Raw(query, database, tableName).Find(&colum).Error; err != nil {
		panic(err)
	}

	return Table{
		TableName: tableName,
		Columns:   colum,
	}
}

func ConvertColumnToProperties(dbCol Column) ClassProperties {
	var cIn ClassProperties
	//  membuat varName
	str2 := stringy.New(dbCol.ColumnName)
	str2 = stringy.New(str2.CamelCase())
	cIn.VarName = str2.LcFirst()
	cIn.VarName = strings.ReplaceAll(cIn.VarName, "Id", "ID")

	// Membuat JSONKey
	cIn.JsonKey = cIn.VarName

	// membuat isnullable
	cIn.IsNullable = true
	if dbCol.IsNullable == "NO" {
		cIn.IsNullable = false
	}
	// fmt.Println(dbCol.IsNullable == "NO", dbCol.IsNullable == "YES")
	// konversi varType
	switch dbCol.DataType {
	case "bigint":
		cIn.VarType = "int"
	case "char":
		cIn.VarType = "String"
	case "date":
		cIn.VarType = "DateTime"
	case "datetime":
		cIn.VarType = "DateTime"
	case "double":
		cIn.VarType = "double"
	case "int":
		cIn.VarType = "int"
	case "longtext":
		cIn.VarType = "String"
	case "mediumtext":
		cIn.VarType = "String"
	case "text":
		cIn.VarType = "String"
	case "timestamp":
		cIn.VarType = "DateTime"
	case "tinyint":
		cIn.VarType = "int"
	case "varchar":
		cIn.VarType = "String"
	default:
		s := stringy.New(dbCol.DataType)
		cIn.VarType = s.CamelCase()
	}

	fmt.Println(cIn)
	// Declaration
	if cIn.IsNullable {
		cIn.Declaration = fmt.Sprintf("final %s? %s;", cIn.VarType, cIn.VarName)
	} else {
		cIn.Declaration = fmt.Sprintf("final %s %s;", cIn.VarType, cIn.VarName)
	}
	// Constructor Line
	if cIn.IsNullable {
		cIn.ConstructorLine = fmt.Sprintf("this.%s,", cIn.VarName)
	} else {
		cIn.ConstructorLine = fmt.Sprintf("required this.%s,", cIn.VarName)
	}

	// From Map line
	if !cIn.IsNullable {
		switch cIn.VarType {
		case "int":
			cIn.FromMapLine = fmt.Sprintf("%s: map['%s'] ?? 0,", cIn.VarName, cIn.JsonKey)
		case "double":
			cIn.FromMapLine = fmt.Sprintf("%s: (map['%s'] ?? 0.0).toDouble(),", cIn.VarName, cIn.JsonKey)
		case "String":
			cIn.FromMapLine = fmt.Sprintf("%s: map['%s'] ?? '',", cIn.VarName, cIn.JsonKey)
		case "DateTime":
			cIn.FromMapLine = fmt.Sprintf("%s: DateTime.parse(map['%s']),", cIn.VarName, cIn.JsonKey)
		default:
			cIn.FromMapLine = fmt.Sprintf("%s: %s.fromMap(map['%s'])", cIn.VarName, cIn.VarType, cIn.JsonKey)
		}
	} else {
		switch cIn.VarType {
		case "int":
			cIn.FromMapLine = fmt.Sprintf("%s: map['%s'],", cIn.VarName, cIn.JsonKey)
		case "double":
			cIn.FromMapLine = fmt.Sprintf("%s: map['%s']?.toDouble(),", cIn.VarName, cIn.JsonKey)
		case "String":
			cIn.FromMapLine = fmt.Sprintf("%s: map['%s'],", cIn.VarName, cIn.JsonKey)
		case "DateTime":
			cIn.FromMapLine = fmt.Sprintf("%s: map['%s'] == null ? null : DateTime.parse(map['%s']),", cIn.JsonKey, cIn.VarName, cIn.JsonKey)
		default:
			cIn.FromMapLine = fmt.Sprintf("%s: %s.fromMap(map['%s'])", cIn.VarName, cIn.VarType, cIn.JsonKey)
		}
	}
	// copyWith
	cIn.CopyWithParamLine = fmt.Sprintf("%s? %s,", cIn.VarType, cIn.VarName)
	cIn.CopyWithPassingLine = fmt.Sprintf("%s: %s ?? this.%s,", cIn.VarName, cIn.VarName, cIn.VarName)

	// toMap
	switch cIn.VarType {
	case "int", "double", "String":
		cIn.ToMapLine = fmt.Sprintf("'%s': %s,", cIn.JsonKey, cIn.VarName)
	case "DateTime":
		if cIn.IsNullable {
			cIn.ToMapLine = fmt.Sprintf("'%s': %s?.toIso8601String(),", cIn.JsonKey, cIn.VarName)
		} else {
			cIn.ToMapLine = fmt.Sprintf("'%s': %s.toIso8601String(),", cIn.JsonKey, cIn.VarName)
		}
	default:
		if cIn.IsNullable {
			cIn.ToMapLine = fmt.Sprintf("'%s': %s?.toMap(),", cIn.JsonKey, cIn.VarName)
		} else {
			cIn.ToMapLine = fmt.Sprintf("'%s': %s.toMap(),", cIn.JsonKey, cIn.VarName)
		}
	}

	// toString
	cIn.ToStringLine = fmt.Sprintf("%s: $%s", cIn.VarName, cIn.VarName)

	// equality
	cIn.EqualityLine = fmt.Sprintf("other.%s == %s", cIn.VarName, cIn.VarName)

	// hasCode
	cIn.HashCodeLine = fmt.Sprintf("%s.hashCode", cIn.VarName)
	return cIn
}

type Table struct {
	TableName string   `json:"name,omitempty"`
	Columns   []Column `json:"columns,omitempty"`
}
type Column struct {
	ColumnName string `gorm:"column:COLUMN_NAME" json:"column_name,omitempty"`
	DataType   string `gorm:"column:DATA_TYPE" json:"data_type,omitempty"`
	IsNullable string `gorm:"column:IS_NULLABLE" json:"is_nullable,omitempty"`
	ColumnKey  string `gorm:"column:COLUMN_KEY" json:"column_key,omitempty"`
}
type Class struct {
	Name       string
	Properties []ClassProperties
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

func tabelToClass(tbl Table) Class {
	var props []ClassProperties
	var class Class
	nam := stringy.New(tbl.TableName)

	class.Name = nam.CamelCase()
	for _, v := range tbl.Columns {
		props = append(props, ConvertColumnToProperties(v))
	}
	class.Properties = props
	return class
}

func classToDart(c Class) string {
	// import class
	s := "import 'dart:convert';\n\n"
	// class
	s += fmt.Sprintf("class %s {\n", c.Name)
	// declaration
	for _, v := range c.Properties {
		s += "  " + fmt.Sprintln(v.Declaration)
	}
	s += "\n"
	// constructor
	s += fmt.Sprintf("  %s({\n", c.Name)
	for _, v := range c.Properties {
		s += "    " + v.ConstructorLine + "\n"
	}
	s += "  });\n\n"
	// fromMap
	s += fmt.Sprintf(
		"  factory %s.fromMap(Map<String, dynamic> map) {\n    return %s(\n",
		c.Name,
		c.Name,
	)
	for _, v := range c.Properties {
		s += "      " + v.FromMapLine + "\n"
	}
	s += "    );\n  }\n\n"
	// copyWith
	s += fmt.Sprintf("  %s copyWith({\n", c.Name)
	for _, v := range c.Properties {
		s += "    " + v.CopyWithParamLine + "\n"
	}
	s += fmt.Sprintf("  }) {\n    return %s(\n", c.Name)
	for _, v := range c.Properties {
		s += "      " + v.CopyWithPassingLine + "\n"
	}
	s += "    );\n"
	s += "  }\n\n"
	// toMap
	s += "  Map<String, dynamic> toMap() {\n"
	s += "    return {\n"
	for _, v := range c.Properties {
		s += "      " + v.ToMapLine + "\n"
	}

	s += "    };\n"
	s += "  }\n\n"
	// toJson
	s += "  String toJson() => json.encode(toMap());\n\n"
	// fromJson
	s += fmt.Sprintf("  factory %s.fromJson(String source) =>%s.fromMap(json.decode(source));\n\n", c.Name, c.Name)
	// toString

	s += "  @override\n"
	s += "  String toString() {\n"
	s += fmt.Sprintf("    return '%s(", c.Name)
	for i, v := range c.Properties {
		if i != len(c.Properties)-1 {
			s += fmt.Sprintf(" %s,", v.ToStringLine)
		} else {
			s += fmt.Sprintf("%s)';\n", v.ToStringLine)
		}
	}
	s += " }\n\n"

	// equality
	s += "  @override"
	s += "  bool operator ==(Object other) {"
	s += "    if (identical(this, other)) return true;"
	s += "\n"
	s += fmt.Sprintf("    return other is %s &&\n", c.Name)
	for i, v := range c.Properties {
		if i < len(c.Properties)-1 {
			s += "        " + v.EqualityLine + "&&\n"

		} else {
			s += "        " + v.EqualityLine + ";\n"

		}
	}
	s += "  }\n\n"

	// hashCode
	s += "  @override\n"
	s += "  int get hashCode {\n"
	for i, v := range c.Properties {
		if i == 0 {
			s += "    return " + v.HashCodeLine + " ^\n"
		} else if i == len(c.Properties)-1 {
			s += "        " + v.HashCodeLine + ";\n"
		} else {
			s += "        " + v.HashCodeLine + " ^\n"

		}
	}
	s += "  }\n\n"

	// class
	s += "}"
	return s
}
