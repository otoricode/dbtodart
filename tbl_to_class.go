package main

import "github.com/gobeam/stringy"

func tablesToClasses(tables []Table) []Class {
	var classes []Class
	for i := 0; i < len(tables); i++ {
		classes = append(classes, tabelToClass(tables[i]))
	}

	return classes
}

func tabelToClass(tbl Table) Class {

	var props []ClassProperties
	var class Class
	nam := stringy.New(tbl.TableName)

	class.Name = nam.CamelCase()
	class.HasList = tbl.HasList
	class.Import = tbl.Relations

	for _, v := range tbl.Columns {
		props = append(props, convertColumnToProperties(v))
	}
	class.Properties = props
	return class
}
