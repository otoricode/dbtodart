package main

import (
	"fmt"
	"strings"

	"github.com/gobeam/stringy"
)

func toVarDartType(dt string) string {

	switch dt {
	case "bigint":
		return "int"
	case "char":
		return "String"
	case "date":
		return "DateTime"
	case "datetime":
		return "DateTime"
	case "double":
		return "double"
	case "int":
		return "int"
	case "longtext":
		return "String"
	case "mediumtext":
		return "String"
	case "text":
		return "String"
	case "timestamp":
		return "DateTime"
	case "tinyint":
		return "int"
	case "varchar":
		return "String"
	default:
		s := stringy.New(dt)
		return s.CamelCase()
	}
}

func declaration(isnull, isList bool, _type, varName string) string {
	if isList {
		if isnull {
			return fmt.Sprintf("final List<%s>? %s;", _type, varName)
		} else {
			return fmt.Sprintf("final List<%s> %s;", _type, varName)
		}
	}
	if isnull {
		return fmt.Sprintf("final %s? %s;", _type, varName)
	} else {
		return fmt.Sprintf("final %s %s;", _type, varName)
	}
}

func constructorLine(isnull bool, varName string) string {
	if isnull {
		return fmt.Sprintf("this.%s,", varName)
	} else {
		return fmt.Sprintf("required this.%s,", varName)
	}
}

func formMapLine(varName, jsonKey, varType string, isNull, isList bool) string {
	if !isList {
		if isNull {
			switch varType {
			case "int":
				return fmt.Sprintf("%s: map['%s'],", varName, jsonKey)
			case "double":
				return fmt.Sprintf("%s: map['%s']?.toDouble(),", varName, jsonKey)
			case "String":
				return fmt.Sprintf("%s: map['%s'],", varName, jsonKey)
			case "DateTime":
				return fmt.Sprintf("%s: map['%s'] == null ? null : DateTime.parse(map['%s']),", varName, jsonKey, jsonKey)
			default:
				return fmt.Sprintf("%s: %s.fromMap(map['%s'])", varName, varType, jsonKey)
			}
		} else {
			switch varType {
			case "int":
				return fmt.Sprintf("%s: map['%s'],", varName, jsonKey)
			case "double":
				return fmt.Sprintf("%s: map['%s'].toDouble(),", varName, jsonKey)
			case "String":
				return fmt.Sprintf("%s: map['%s'],", varName, jsonKey)
			case "DateTime":
				return fmt.Sprintf("%s: DateTime.parse(map['%s']),", varName, jsonKey)
			default:
				return fmt.Sprintf("%s: %s.fromMap(map['%s'])", varName, varType, jsonKey)
			}
		}
	}

	switch varType {
	case "int", "double", "String", "DateTime":
		panic("unimplemented")
	default:
		return fmt.Sprintf("%s: map['%s'] == null ? null : List<%s>.from(map['%s']?.map((x) => %s.fromMap(x))),", varName, jsonKey, varType, jsonKey, varType)
	}
	//
}

func toMapLine(varType, jsonKey, varName string, isnull, isList bool) string {
	if isList {
		switch varType {
		case "int", "double", "String", "DateTime":
			panic("unimplemented")
		default:
			return fmt.Sprintf("'%s': %s?.map((x) => x.toMap()).toList(),", jsonKey, varName)
		}
	}
	switch varType {
	case "int", "double", "String":
		return fmt.Sprintf("'%s': %s,", jsonKey, varName)
	case "DateTime":
		if isnull {
			return fmt.Sprintf("'%s': %s?.toIso8601String(),", jsonKey, varName)
		} else {
			return fmt.Sprintf("'%s': %s.toIso8601String(),", jsonKey, varName)
		}
	default:
		if isnull {
			return fmt.Sprintf("'%s': %s?.toMap(),", jsonKey, varName)
		} else {
			return fmt.Sprintf("'%s': %s.toMap(),", jsonKey, varName)
		}
	}
}
func convertColumnToProperties(dbCol Column) ClassProperties {
	var cIn ClassProperties
	//  membuat varName
	str2 := stringy.New(dbCol.ColumnName)
	str2 = stringy.New(str2.CamelCase())
	cIn.VarName = str2.LcFirst()
	cIn.VarName = strings.ReplaceAll(cIn.VarName, "Id", "ID")

	if dbCol.IsList {
		cIn.VarName = plur.Plural(cIn.VarName)

	} else {
		cIn.VarName = plur.Singular(cIn.VarName)
	}

	// Membuat JSONKey
	cIn.JsonKey = cIn.VarName

	// membuat isnullable
	cIn.IsNullable = true
	if dbCol.IsNullable == "NO" {
		cIn.IsNullable = false
	}
	// konversi varType
	cIn.VarType = toVarDartType(dbCol.DataType)

	// Declaration
	cIn.Declaration = declaration(cIn.IsNullable, dbCol.IsList, cIn.VarType, cIn.VarName)

	// Constructor Line
	cIn.ConstructorLine = constructorLine(cIn.IsNullable, cIn.VarName)

	// From Map line
	cIn.FromMapLine = formMapLine(cIn.VarName, cIn.JsonKey, cIn.VarType, cIn.IsNullable, dbCol.IsList)
	// copyWith
	if dbCol.IsList {
		cIn.CopyWithParamLine = fmt.Sprintf("List<%s>? %s,", cIn.VarType, cIn.VarName)
		cIn.CopyWithPassingLine = fmt.Sprintf("%s: %s ?? this.%s,", cIn.VarName, cIn.VarName, cIn.VarName)
	} else {
		cIn.CopyWithParamLine = fmt.Sprintf("%s? %s,", cIn.VarType, cIn.VarName)
		cIn.CopyWithPassingLine = fmt.Sprintf("%s: %s ?? this.%s,", cIn.VarName, cIn.VarName, cIn.VarName)

	}

	// toMap
	cIn.ToMapLine = toMapLine(cIn.VarType, cIn.JsonKey, cIn.VarName, cIn.IsNullable, dbCol.IsList)

	// toString
	cIn.ToStringLine = fmt.Sprintf("%s: $%s", cIn.VarName, cIn.VarName)

	// equality
	if dbCol.IsList {
		cIn.EqualityLine = fmt.Sprintf("DeepCollectionEquality().equals(other.%s, %s)", cIn.VarName, cIn.VarName)
	} else {
		cIn.EqualityLine = fmt.Sprintf("other.%s == %s", cIn.VarName, cIn.VarName)
	}

	// hasCode
	cIn.HashCodeLine = fmt.Sprintf("%s.hashCode", cIn.VarName)
	return cIn
}
