package main

import (
	"fmt"
	"strings"

	"github.com/gobeam/stringy"
)

func tblToGo(t Table) string {
	/*
			type AffiliateLogs struct {
			ID             int        `db:"id" gorm:"column:id;primaryKey;autoIncrement" json:"id"`
			UserID         *int       `db:"user_id" gorm:"column:user_id" json:"userID"`
			GuestID        *int       `db:"guest_id" gorm:"column:guest_id" json:"guestID"`
			ReferredByUser int        `db:"referred_by_user" gorm:"column:referred_by_user" json:"referredByUser"`
			Amount         *float64   `db:"amount" gorm:"column:amount" json:"amount"`
			OrderID        *int       `db:"order_id" gorm:"column:order_id" json:"orderID"`
			OrderDetailID  *int       `db:"order_detail_id" gorm:"column:order_detail_id" json:"orderDetailID"`
			AffiliateType  string     `db:"affiliate_type" gorm:"column:affiliate_type" json:"affiliateType"`
			Status         int        `db:"status" gorm:"column:status" json:"status"`
			CreatedAt      *time.Time `db:"created_at" gorm:"column:created_at" json:"createdAt"`
			UpdatedAt      *time.Time `db:"updated_at" gorm:"column:updated_at" json:"updatedAt"`
		}
	*/
	structName := stringy.New(plur.Singular(t.TableName)).CamelCase()
	s := "package models\n"
	s += "\n"
	s += "import \"time\"\n"
	s += "\n"
	s += "type " + structName + " struct {\n"
	for _, v := range t.Columns {
		fielName := pascalCase(v.ColumnName)
		if v.IsList {
			fielName = plur.Plural(fielName)
		}
		_type := toVarGoType(v.DataType)
		columnName := v.ColumnName
		prim := ""
		if v.ColumnKey == "PRI" {
			prim = "primaryKey;"
		}
		null := ""
		if v.IsNullable == "YES" {
			null = "*"
		}
		list := ""
		if v.IsList {
			list = "[]"
		}
		_json := camelCase(v.ColumnName)
		if v.IsList {
			_json = plur.Plural(_json)
		}
		foreign := ""
		if v.IsList {
			foreign = "foreignKey:" + stringy.New(plur.Singular(t.TableName)).CamelCase() + "ID;"
		}
		s += fmt.Sprintf("	%s %s%s%s `db:\"%s\" gorm:\"column:%s;%s%s\" json:\"%s\"`\n", fielName, list, null, _type, columnName, columnName, prim, foreign, _json)
	}
	s += "}\n\n"

	s += fmt.Sprintf("func (%s)TableName() string {\n", structName)
	s += "return \"" + t.TableName + "\"\n"
	s += "}"
	return s
}

func camelCase(in string) string {
	out := stringy.New(stringy.New(plur.Singular(in)).CamelCase()).LcFirst()
	out = strings.ReplaceAll(out, "Id", "ID")

	return out
}

func pascalCase(in string) string {
	out := stringy.New(in).CamelCase()
	out = strings.ReplaceAll(out, "Id", "ID")
	return out
}

func toVarGoType(dt string) string {

	switch dt {
	case "bigint":
		return "int"
	case "char":
		return "string"
	case "date":
		return "time.Time"
	case "datetime":
		return "time.Time"
	case "double":
		return "float64"
	case "int":
		return "int"
	case "longtext":
		return "string"
	case "mediumtext":
		return "string"
	case "text":
		return "string"
	case "timestamp":
		return "time.Time"
	case "tinyint":
		return "int"
	case "varchar":
		return "string"
	default:
		s := stringy.New(dt)
		return plur.Singular(s.CamelCase())
	}
}
