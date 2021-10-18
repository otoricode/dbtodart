package main

import "gorm.io/gorm"

func getTables(db *gorm.DB, database string) []Table {
	_db, _ := db.DB()

	res, err := _db.Query("SHOW TABLES")

	if err != nil {
		panic(err)
	}
	var tables []Table
	for res.Next() {
		var tableName string
		res.Scan(&tableName)
		tables = append(tables, getTable(db, database, tableName))
	}
	for i, v1 := range tables {
		for _, v2 := range tables {
			for _, c := range v2.Columns {
				if plur.Singular(v1.TableName)+"_id" == c.ColumnName {
					tables[i].Columns = append(tables[i].Columns, Column{
						ColumnName: v2.TableName,
						DataType:   v2.TableName,
						IsNullable: "YES",
						IsList:     true,
						ColumnKey:  "",
					})
					tables[i].HasList = true
					tables[i].Relations = append(tables[i].Relations, v2.TableName)
				}
			}
		}
	}
	return tables
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
