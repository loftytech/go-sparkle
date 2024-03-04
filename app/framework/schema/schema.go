package schema

import (
	"coralscale/app/framework/db"
	"coralscale/app/framework/utility"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
)

var database *sql.DB

type TableStruct struct {
	Field   string         `json:"Field"`
	Type    string         `json:"Type"`
	Null    string         `json:"Null"`
	Key     string         `json:"Key"`
	Default sql.NullString `json:"Default"`
	Extra   string         `json:"Extra"`
}



type Schema struct {
	TableName   string
	columns     []Column
	schema_sql  string
	newTableSlice []TableStruct
	newTableMap []map[string]string
	savedTableSlice []TableStruct
	savedTableMap []map[string]string
	savedTableJsonString string
	newTableJsonString string
	primary_key string
	skippedMigration bool
	useTimestamps bool
	modificationList []map[string]string
	columnDropList []map[string]string
	replacedColumnList []string
	current_column_index int
	num_primary_keys int
	current_column string
}

type Column struct {
	Name string
	Size int
	Type string
}

func (s *Schema) Id() {
	s.current_column = "id"
	s.columns = append(s.columns, Column{
		Name: "id",
		Type: "bigint(20)",
		Size: 11,
	})

	append_comma := ""

	if len(s.columns) > 1 {
		append_comma = ","
		s.current_column_index++
	} else {
		s.current_column_index = 0
	}

	s.newTableSlice = append(s.newTableSlice, TableStruct{
		Field: "id",
		Type:  "bigint(20)",
		Null:  "NO",
		Key:   "PRI",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "auto_increment",
	})

	s.schema_sql = s.schema_sql + append_comma + " `id` bigint(20) NOT NULL AUTO_INCREMENT"
	s.primary_key = ", PRIMARY KEY (id)"
}

func (s *Schema) Integer(name string) *Schema {
	s.current_column = name
	s.columns = append(s.columns, Column{
		Name: name,
		Type: "bigint(20)",
	})

	append_comma := ""

	if len(s.columns) > 1 {
		append_comma = ","
		s.current_column_index++
	} else {
		s.current_column_index = 0
	}

	s.newTableSlice = append(s.newTableSlice, TableStruct{
		Field: name,
		Type:  "bigint(20)",
		Null:  "NO",
		Key:   "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	})

	s.schema_sql = s.schema_sql + append_comma + " `" + name + "` bigint(20) NOT NULL"
	return s
}


func (s *Schema) Double(name string) {
	s.current_column = name
	s.columns = append(s.columns, Column{
		Name: name,
		Type: "double",
	})

	append_comma := ""

	if len(s.columns) > 1 {
		append_comma = ","
		s.current_column_index++
	} else {
		s.current_column_index = 0
	}

	s.newTableSlice = append(s.newTableSlice, TableStruct{
		Field: name,
		Type:  "double",
		Null:  "NO",
		Key:   "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	})

	s.schema_sql = s.schema_sql + append_comma + " `" + name + "` double NOT NULL"
}

func (s *Schema) String(name string, size int) {
	s.current_column = name
	s.columns = append(s.columns, Column{
		Name: name,
		Type: "varchar",
		Size: size,
	})

	append_comma := ""

	if len(s.columns) > 1 {
		append_comma = ","
		s.current_column_index++
	} else {
		s.current_column_index = 0
	}

	s.newTableSlice = append(s.newTableSlice, TableStruct{
		Field: name,
		Type:  "varchar("+strconv.Itoa(size)+")",
		Null:  "NO",
		Key:   "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	})

	s.schema_sql = s.schema_sql + append_comma + " `" + name + "` varchar(" + strconv.Itoa(size) + ") NOT NULL"
}

func (s *Schema) Text(name string) {
	s.current_column = name
	s.columns = append(s.columns, Column{
		Name: name,
		Type: "text",
	})

	append_comma := ""

	if len(s.columns) > 1 {
		append_comma = ","
		s.current_column_index++
	} else {
		s.current_column_index = 0
	}

	s.newTableSlice = append(s.newTableSlice, TableStruct{
		Field: name,
		Type:  "text",
		Null:  "NO",
		Key:   "",
		Default: sql.NullString{
			String: "",
			Valid:  false,
		},
		Extra: "",
	})

	s.schema_sql = s.schema_sql + append_comma + " `" + name + "` text NOT NULL"
}

func (s *Schema) AutoIncrement() *Schema {
	s.newTableSlice[s.current_column_index].Extra = "auto_increment";
	s.schema_sql = s.schema_sql + " AUTO_INCREMENT";
	return s
}

func (s *Schema) Primary() {
	if (s.num_primary_keys > 0) {
		utility.LogError("Aborting migration: Duplicate primary keys detected in " + s.TableName + " table there can only be one auto column and it must be defined as a key")
	}
	s.newTableSlice[s.current_column_index].Key = "PRI";
	s.primary_key = ", PRIMARY KEY (" + s.current_column + ")";
	s.num_primary_keys = s.num_primary_keys+1;
}


func (s *Schema) Create() {
	s.useTimestamps = true
	database = db.GetDBInstance()
	is_exists := s.checkTabeExist()

	if !is_exists {
		s.createTable()
	} else {
		if (s.useTimestamps) {
			s.enableTimestamps()
		}
		s.compareColumns()
	}
}

func (s *Schema) enableTimestamps() {
	s.checkTableExists("date_created");
	s.checkTableExists("date_updated");
	// s.current_column_index = s.current_column_index + 1;
	s.newTableSlice = append(s.newTableSlice, TableStruct{
		Field: "date_created",
		Type:  "datetime",
		Null:  "NO",
		Key:   "",
		Default: sql.NullString{
			String: "current_timestamp()",
			Valid:  true,
		},
		Extra: "",
	})
	
	s.newTableSlice = append(s.newTableSlice, TableStruct{
		Field: "date_updated",
		Type:  "datetime",
		Null:  "NO",
		Key:   "",
		Default: sql.NullString{
			String: "current_timestamp()",
			Valid:  true,
		},
		Extra: "",
	})
}

func (s *Schema) checkTableExists (newColumn string) {
	trimed_str := strings.TrimSpace(newColumn)
	if (strlen(&trimed_str) < 1) {
		panic("Aborting migration: Empty column detected in " + s.TableName + " table. Columns can't be empty")
	}
	for _, column := range s.newTableMap {
		if exists_in_map(newColumn, column) {
			panic("Aborting migration: Duplicate colunm `" + newColumn + "` detected in " + s.TableName + " table. Columns must be unique")
		}
	}
}

func strlen(value *string) int {
	len := 0
    for key := range *value {
        len = key+1
    }

	return len
}

func (s *Schema) createTable() {
	if (s.useTimestamps) {
		s.schema_sql =  s.schema_sql + ", `date_created` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, `date_updated` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP"
	}

	create_table_query := "CREATE TABLE IF NOT EXISTS " + s.TableName + " (" + s.schema_sql + s.primary_key + ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"

	db.Exec(create_table_query)
	utility.LogSuccess(s.TableName + " created successfully")
}

func (s *Schema) checkTabeExist() bool {
	create_table_query := "DESCRIBE " + s.TableName
	rows, error := database.Query(create_table_query)
	if error != nil {
		return false
	}
	for rows.Next() {
		var _table TableStruct
		rows.Scan(&_table.Field, &_table.Type, &_table.Null, &_table.Key, &_table.Default, &_table.Extra)
		s.savedTableSlice = append(s.savedTableSlice, _table)
	}

	// a, _ := json.Marshal(&s.savedTableSlice)
	// _ = json.Unmarshal(a, &s.savedTableMap)
	
	b, _ := json.Marshal(&s.newTableSlice)
	_ = json.Unmarshal(b, &s.newTableMap)

	// s.savedTableJsonString = string(a)
	// s.newTableJsonString = string(b)

	// fmt.Println(s.savedTableSlice)
	// fmt.Println(s.savedTableMap)

	// utility.LogSuccess(s.savedTableJsonString)
	// utility.LogNeutral(s.newTableJsonString)

	// s.removeTimestampsFromNewColumns(&s.newTableSlice)
	s.removeTimestampsFromSavedColumns(&s.savedTableSlice)
	return true
}


func (s *Schema) compareColumns() {
	s.removeTimestampsFromNewColumns(&s.newTableSlice)
	if len(s.savedTableMap) == len(s.newTableMap) {
		s.compareAndAlterColumns()
	} else {
		s.getNewColumns()
	}
	s.migrateSchema()
}

func (s *Schema) compareAndAlterColumns() {
	is_columns_match := false
	if s.savedTableJsonString == s.newTableJsonString {
		is_columns_match = true
	}

	if is_columns_match {
		s.skippedMigration = true
		utility.LogNeutral("Skipping " + s.TableName + " table because no changes was found")
	} else {
		for newColumnkey, column := range s.newTableMap {
			new_column := s.getColumnType(column["Type"])

			add_primary_key := "";
			auto_increment_sub_query := "";

			saved_column_to_alter := s.savedTableMap[newColumnkey]

			
			if column["Extra"] == "auto_increment" && saved_column_to_alter["Extra"] != "auto_increment" && saved_column_to_alter["Key"] != "PRI" && column["Key"] == "PRI" {
				add_primary_key = ", add PRIMARY KEY (`" + column["Field"] + "`)"
				auto_increment_sub_query = " AUTO_INCREMENT"
			}

			if column["Extra"] == "auto_increment" && saved_column_to_alter["Extra"] != "auto_increment" && saved_column_to_alter["Key"] == "PRI" && column["Key"] == "PRI" {
				auto_increment_sub_query = " AUTO_INCREMENT"
			}

			if column["Extra"] == "auto_increment" && saved_column_to_alter["Extra"] != "auto_increment" && saved_column_to_alter["Key"] == "PRI" {
				auto_increment_sub_query = " AUTO_INCREMENT"
			}
		
			if column["Extra"] == "auto_increment" && saved_column_to_alter["Extra"] == "auto_increment" {
				auto_increment_sub_query = " AUTO_INCREMENT"
			}

			mod_sql := "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + column["Field"] + "` " + new_column.Type + "(" + new_column.Limit + ") NOT NULL " + auto_increment_sub_query + " " + add_primary_key + "; "
			
			if new_column.Type == "text" {
				mod_sql = "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + column["Field"] + "` " + new_column.Type + " NOT NULL " + auto_increment_sub_query + " " + add_primary_key + "; "
			} else if new_column.Type == "double" {
				mod_sql = "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + column["Field"] + "` " + new_column.Type + " NOT NULL " + auto_increment_sub_query + " " + add_primary_key + "; "
			} else if new_column.Type == "datetime" {
				mod_sql = "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + column["Field"] + "` DATETIME NOT NULL DEFAULT " + s.savedTableSlice[newColumnkey].Default.String + "; "
			}

			s.modificationList = append(s.modificationList, map[string]string{
				"sql": mod_sql,
				"from": saved_column_to_alter["Field"],
				"to": column["Field"],
				"from_type": saved_column_to_alter["Type"],
				"to_type": column["Type"],
				"operation_type": "ALTER_CHANGE_COLUMN",
			})
	
			if column["Extra"] != "auto_increment" && saved_column_to_alter["Extra"] != "auto_increment" && saved_column_to_alter["Key"] != "PRI" && column["Key"] == "PRI" {
				s.modificationList = append(s.modificationList, map[string]string{
					"sql": "ALTER TABLE `" + s.TableName + "` ADD PRIMARY KEY(`" + column["Field"] + "`);",
					"operation_type": "ALTER_ADD_PRIMARY_KEY",
				})
			}

			if column["Extra"] != "auto_increment" && saved_column_to_alter["Extra"] != "auto_increment" && saved_column_to_alter["Key"] == "PRI" && column["Key"] != "PRI" {
				s.columnDropList = append(s.columnDropList, map[string]string{
					"sql": "ALTER TABLE `" + s.TableName + "` DROP PRIMARY KEY;",
					"operation_type": "ALTER_DROP_PRIMARY_KEY",
				})
			}

			/*
			*  Check if new column has a primary key
			*/

			if column["Extra"] != "auto_increment" && saved_column_to_alter["Extra"] == "auto_increment" && saved_column_to_alter["Key"] == "PRI" && column["Key"] != "PRI" {
				/*
				*  Remove auto increment from current column
				*/
				s.columnDropList = append(s.columnDropList, map[string]string{
					"sql": "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + saved_column_to_alter["Field"] + "` " + saved_column_to_alter["Type"] + " NOT NULL; ",
					"from": saved_column_to_alter["Field"],
					"to": saved_column_to_alter["Field"],
					"operation_type": "ALTER_CHANGE_COLUMN",
				})

				/*
				*  Remove primary key from table
				*/

				s.columnDropList = append(s.columnDropList, map[string]string{
					"sql": "ALTER TABLE `" + s.TableName + "` DROP PRIMARY KEY;",
					"operation_type": "ALTER_DROP_PRIMARY_KEY",
				})
			}
		}
	}

}

func (s *Schema) getNewColumns() {
	s.filterColumnsToDrop();

	saved_table_length := len(s.savedTableMap);
	for newColumnkey, column := range s.newTableMap {
		new_column := s.getColumnType(column["Type"])

		auto_increment_sub_query := ""
		add_primary_key := ""

		if newColumnkey < saved_table_length {
			saved_column_to_alter := s.savedTableMap[newColumnkey]
			
			if column["Extra"] == "auto_increment" && saved_column_to_alter["Extra"] != "auto_increment" && saved_column_to_alter["Key"] != "PRI" && column["Key"] == "PRI" {
				add_primary_key = ", add PRIMARY KEY (`" + column["Field"] + "`)"
				auto_increment_sub_query = " AUTO_INCREMENT"
			}

			if column["Extra"] == "auto_increment" && saved_column_to_alter["Extra"] != "auto_increment" && saved_column_to_alter["Key"] == "PRI" && column["Key"] == "PRI" {
				auto_increment_sub_query = " AUTO_INCREMENT";
			}

			if column["Extra"] == "auto_increment" && saved_column_to_alter["Extra"] != "auto_increment" && saved_column_to_alter["Key"] == "PRI" {
				auto_increment_sub_query = " AUTO_INCREMENT";
			}
		
			if column["Extra"] == "auto_increment" && saved_column_to_alter["Extra"] == "auto_increment" {
				auto_increment_sub_query = " AUTO_INCREMENT";
			}

			mod_sql := "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + column["Field"] + "` " + new_column.Type + "(" + new_column.Limit + ") NOT NULL " + auto_increment_sub_query + " " + add_primary_key +"; "
			
			if new_column.Type == "text" {
				mod_sql = "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + column["Field"] + "` " + new_column.Type + " NOT NULL " + auto_increment_sub_query + " " + add_primary_key +"; "
			} else if new_column.Type == "double" {
				mod_sql = "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + column["Field"] + "` " + new_column.Type + " NOT NULL " + auto_increment_sub_query + " " + add_primary_key +"; "
			} else if new_column.Type == "datetime" {
				mod_sql = "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + column["Field"] + "` DATETIME NOT NULL DEFAULT " + s.savedTableSlice[newColumnkey].Default.String + "; "
			}

			s.modificationList = append(s.modificationList, map[string]string{
				"sql": mod_sql,
				"from": saved_column_to_alter["Field"],
				"to": column["Field"],
				"from_type": saved_column_to_alter["Type"],
				"to_type": column["Type"],
				"operation_type": "ALTER_CHANGE_COLUMN",
			})

			if column["Extra"] != "auto_increment" && saved_column_to_alter["Extra"] != "auto_increment" && saved_column_to_alter["Key"] != "PRI" && column["Key"] == "PRI" {
				s.modificationList = append(s.modificationList, map[string]string{
					"sql": "ALTER TABLE `" + s.TableName + "` ADD PRIMARY KEY(`" + column["Field"] + "`);",
					"operation_type": "ADD_PRIMARY_KEY",
				})
			}

			if column["Extra"] != "auto_increment" && saved_column_to_alter["Extra"] != "auto_increment" && saved_column_to_alter["Key"] == "PRI" && column["Key"] != "PRI" {
				s.modificationList = append(s.modificationList, map[string]string{
					"sql": "ALTER TABLE `" + s.TableName + "` DROP PRIMARY KEY;",
					"operation_type": "DROP_PRIMARY_KEY",
				})
			}

			/*
			*  Check if new column has a primary key
			*/

			if column["Extra"] != "auto_increment" && saved_column_to_alter["Extra"] == "auto_increment" && saved_column_to_alter["Key"] == "PRI" && column["Key"] != "PRI" {


				/*
				*  Remove auto increment from current column
				*/
				s.columnDropList = append(s.columnDropList, map[string]string{
					"sql": "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + saved_column_to_alter["Field"] + "` " + saved_column_to_alter["Type"] + " NOT NULL; ",
					"from": saved_column_to_alter["Field"],
					"to": saved_column_to_alter["Field"],
					"operation_type": "ALTER_CHANGE_COLUMN",
				})
				
				/*
				*  Remove primary key from table
				*/
				s.columnDropList = append(s.columnDropList, map[string]string{
					"sql": "ALTER TABLE `" + s.TableName + "` DROP PRIMARY KEY;",
					"operation_type": "ALTER_DROP_PRIMARY_KEY",
				})
			}
		} else {
			alter_after_sub_query := ""
			if newColumnkey != 0 {
				previous_column_field := s.getPreviousColumn(s.newTableMap, newColumnkey)["Field"];
				alter_after_sub_query = "AFTER `" + previous_column_field + "`"
			} else {
				alter_after_sub_query = "FIRST"
			}

			if column["Extra"] == "auto_increment" {
				auto_increment_sub_query = " AUTO_INCREMENT ";
				add_primary_key = " , add PRIMARY KEY (`" + column["Field"] + "`)";
			}

			mod_sql := "ALTER TABLE `" + s.TableName + "` ADD `" + column["Field"] + "` " + new_column.Type + "(" + new_column.Limit + ") NOT NULL " + auto_increment_sub_query + " " + alter_after_sub_query + " " + add_primary_key + "; "
			
			if new_column.Type == "text" {
				mod_sql = "ALTER TABLE `" + s.TableName + "` ADD `" + column["Field"] + "` " + new_column.Type + " NOT NULL " + auto_increment_sub_query + " " + alter_after_sub_query + " " + add_primary_key + "; "
			} else if new_column.Type == "double" {
				mod_sql = "ALTER TABLE `" + s.TableName + "` ADD `" + column["Field"] + "` " + new_column.Type + " NOT NULL " + auto_increment_sub_query + " " + alter_after_sub_query + " " + add_primary_key + "; "
			} else if new_column.Type == "datetime" {
				mod_sql = "ALTER TABLE `" + s.TableName + "` ADD `" + column["Field"] + "` " + new_column.Type + " NOT NULL DEFAULT CURRENT_TIMESTAMP " + alter_after_sub_query + ";"
			}

			s.modificationList = append(s.modificationList, map[string]string{
				"sql": mod_sql,
				"operation_type": "ALTER_ADD_COLUMN",
			})

			if column["Extra"] != "auto_increment" && column["Key"] == "PRI" {
				s.columnDropList = append(s.columnDropList, map[string]string{
					"sql": "ALTER TABLE `" + s.TableName + "` ADD PRIMARY KEY(`" + column["Field"] + "`);",
					"operation_type":"ALTER_ADD_PRIMARY_KEY",
				})
			}
		}
	}
}

type ColumnTypeStruct struct {
	Type string
	Limit string
}

func (s *Schema) getColumnType(field string) ColumnTypeStruct {
	field_arry := strings.Split(field, "(");
	openParenthesis := '('

	if strings.ContainsRune(field, openParenthesis) {
		return ColumnTypeStruct{
			Type: field_arry[0],
			Limit: strings.Replace(field_arry[1], ")", "", -1),
		}
	} else {
		return ColumnTypeStruct{
			Type: field,
			Limit: "1",
		}
	}

}

func (s *Schema) filterColumnsToDrop() {
	new_table_length := len(s.newTableMap)
	if new_table_length < len(s.savedTableMap) {
		utility.LogWarning("Some columns would be dropped");
		colums_to_be_dropped := []map[string]string{};
		filtered_column_list := []map[string]string{};

		for savedColumnKey, savedColumn := range s.savedTableMap {
			/*
			*
			*   Check for columns that already exists in the table
			*
			*/

			if savedColumnKey < new_table_length {
				filtered_column_list = append(filtered_column_list, savedColumn)
			} else {
				colums_to_be_dropped = append(colums_to_be_dropped, savedColumn)
			}
		}

		

		for _, colum_to_be_dropped := range colums_to_be_dropped {
			s.modificationList = append(s.modificationList, map[string]string{
				"sql": "ALTER TABLE `" + s.TableName + "` DROP `" + colum_to_be_dropped["Field"] + "`; ",
				"operation_type": "ALTER_DROP",
			})
		}

		s.savedTableMap = filtered_column_list
	}
}

func (s *Schema) getPreviousColumn(arr []map[string]string, key int) map[string]string {
	return arr[key-1];
}

func (s *Schema) migrateSchema()  {
	// utility.LogNeutral("running migration ...")
	for _, column_drop_query := range s.columnDropList {
		db.Exec(column_drop_query["sql"])
	}

	for _, modification := range s.modificationList {
		if modification["operation_type"] == "ALTER_CHANGE_COLUMN" {
			if exists_in_slice(modification["from"], s.replacedColumnList) {
				checkedColumn := s.checkColumnExists(modification["to"])
				if checkedColumn.exists {
					if modification["from"] != modification["to"] {
						s.replacedColumnList = append(s.replacedColumnList, modification["to"])
						column := s.getColumnType(checkedColumn.data.Type)
						s.buildAndRunQuery(QueryBuildStrct{
							QueryType: "ALTER_CHANGE_COLUMN", 
							From: modification["to"], 
							To: modification["to"] + "_ALTERED", 
							Type: column.Type, 
							Limit: column.Limit, 
							FromType: modification["from_type"], 
							ToType: modification["to_type"], 
							Extras: checkedColumn.data.Extra,
						})
					}
				}
				new_sql := strings.Replace( "CHANGE `" + modification["from"] + "_ALTERED`", modification["sql"], "CHANGE `" + modification["from"] + "`", -1);
				db.Exec(new_sql)
			} else {
				checkedColumn := s.checkColumnExists(modification["to"])
				if checkedColumn.exists {
					if modification["from"] != modification["to"] {
						s.replacedColumnList = append(s.replacedColumnList, modification["to"])
						column := s.getColumnType(checkedColumn.data.Type);

						s.buildAndRunQuery(QueryBuildStrct{
							QueryType: "ALTER_CHANGE_COLUMN", 
							From: modification["to"], 
							To: modification["to"] + "_ALTERED", 
							Type: column.Type, 
							Limit: column.Limit, 
							FromType: modification["from_type"], 
							ToType: modification["to_type"], 
							Extras: checkedColumn.data.Extra,
						})

						s.clearColumn(ClearFormStruct{
							from: modification["from"],
							to: modification["to"],
							fromType: modification["from_type"],
							toType: modification["to_type"],
						});
						db.Exec(modification["sql"]);
					} else {
						s.clearColumn(ClearFormStruct{
							from: modification["from"],
							to: modification["to"],
							fromType: modification["from_type"],
							toType: modification["to_type"],
						});
						db.Exec(modification["sql"]);
					}
				} else {
					s.clearColumn(ClearFormStruct{
						from: modification["from"],
						to: modification["to"],
						fromType: modification["from_type"],
						toType: modification["to_type"],
					});
					db.Exec(modification["sql"]);
				}
			}

		} else {
			db.Exec(modification["sql"]);
		}

	}

	if !s.skippedMigration {
		utility.LogSuccess(s.TableName + " table migrated successfully")
	}
}

// func in_array(target map[string]string, slice []map[string]string) bool {
//     for _, m := range slice {
//         if mapsEqual(m, target) {
//             return true
//         }
//     }
//     return false
// }

func exists_in_slice(target string, slice []string) bool {
    for _, m := range slice {
        if m == target {
            return true
        }
    }
    return false
}

func exists_in_map(target string, slice map[string]string) bool {
    for _, m := range slice {
        if m == target {
            return true
        }
    }
    return false
}

// func mapsEqual(a, b map[string]string) bool {
//     if len(a) != len(b) {
//         return false
//     }
//     for key, valA := range a {
//         if valB, ok := b[key]; !ok || valA != valB {
//             return false
//         }
//     }
//     return true
// }

type CheckColumeReturnStruct struct {
	exists bool
	data TableStruct
}

func (s *Schema) checkColumnExists (column string) CheckColumeReturnStruct {
	query := "DESCRIBE " + s.TableName;
	
	var _tableSlice []TableStruct
	var _tableMap []map[string]string
	colunm_data := TableStruct{};
	column_exists := false;

	rows, error := database.Query(query)
	if error != nil {
		utility.LogError("caling panic")
		panic(error.Error())
	}

	for rows.Next() {
		var _table TableStruct
		rows.Scan(&_table.Field, &_table.Type, &_table.Null, &_table.Key, &_table.Default, &_table.Extra)
		_tableSlice = append(_tableSlice, _table)
	}
	
	
	b, _ := json.Marshal(&_tableSlice)
	_ = json.Unmarshal(b, &_tableMap)

	for savedColumnKey, savedColumn := range _tableMap {
		if exists_in_map(column, savedColumn) {
			colunm_data = _tableSlice[savedColumnKey]
			column_exists = true
			break
		}
	}

	return CheckColumeReturnStruct{
		exists: column_exists,
		data: colunm_data,
	};
}

type QueryBuildStrct struct {
	QueryType string
	From string
	To string
	Type string
	Limit string
	FromType string
	ToType string
	Extras string
	Key string
}

func (s *Schema) buildAndRunQuery(query QueryBuildStrct) {

	auto_increment_sub_query := "";
	add_primary_key := "";

	if (query.Extras == "auto_increment") {
		add_primary_key = ", add PRIMARY KEY (`" + query.To + "`)"
		auto_increment_sub_query = " AUTO_INCREMENT";
	}
	if (query.Key == "PRI") {
		add_primary_key = ", add PRIMARY KEY (`" + query.To + "`)"
		auto_increment_sub_query = " AUTO_INCREMENT";
	}

	if (query.QueryType == "ALTER_CHANGE_COLUMN") {
		sql := "ALTER TABLE `" + s.TableName + "` CHANGE `" + query.From + "` `" + query.To + "` " + query.Type + "(" + query.Limit + ") NOT NULL " + auto_increment_sub_query + " " + add_primary_key + "; "

		if Contains(query.Type, []string{"text", "datetime", "double"} ) {
			sql = "ALTER TABLE `" + s.TableName + "` CHANGE `" + query.From + "` `" + query.To + "` " + query.Type + " NOT NULL " + auto_increment_sub_query + " " + add_primary_key + "; "
		}
		
		db.Exec(sql);
	}
}

type ClearFormStruct struct {
	from string
	to string
	fromType string
	toType string
}

func (s *Schema) clearColumn(clear ClearFormStruct) {
	raw_fromType := s.getColumnType(clear.fromType)
	raw_toType := s.getColumnType(clear.toType)
	// echo "to: $to (raw_toType.Type) from: $from (raw_fromType.Type) \n";
	if (raw_fromType.Type != "" && raw_toType.Type != "" && raw_fromType.Type != raw_toType.Type) {
		if (Contains(strings.ToLower(raw_toType.Type), []string{"int", "bigint", "double"}) && Contains(strings.ToLower(raw_fromType.Type), []string{"varchar", "text", "datetime"})) {
			update_sql := "UPDATE `" + s.TableName + "` SET `" + clear.from + "`=0"
			// echo "$$update_sql\n";
			// Log::warning($update_sql);
			db.Exec(update_sql);
		} else {
			if Contains(strings.ToLower(raw_toType.Type), []string{"datetime"}) {
				update_sql := "UPDATE `" + s.TableName + "` SET `" + clear.from + "`='0000-00-00 00:00:00'"
				// echo "$$update_sql\n";
				// Log::warning($update_sql);
				db.Exec(update_sql);
			}
		}
	}
}

func Contains(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func (s *Schema) removeTimestampsFromNewColumns (tableSlice *[]TableStruct) {	
	var _tableSlice []TableStruct = *tableSlice

	table_columns_length := len(_tableSlice)

	if s.useTimestamps {
		date_created_index := table_columns_length-2
		date_updated_index := table_columns_length-1
		if  len(_tableSlice) > 2 && _tableSlice[date_created_index].Field == "date_created" && _tableSlice[date_updated_index].Field == "date_updated" {
			// _tableSlice = append(_tableSlice[:date_created_index], _tableSlice[date_created_index+2:]...)
			utility.LogBlue("preserving timestamps if exist")
		}
	}

	a, _ := json.Marshal(&_tableSlice)
	_ = json.Unmarshal(a, &s.newTableMap)
	s.newTableJsonString = string(a)
}

func (s *Schema) removeTimestampsFromSavedColumns (tableSlice *[]TableStruct) {
	var _tableSlice []TableStruct = *tableSlice

	table_columns_length := len(_tableSlice)

	if s.useTimestamps {
		date_created_index := table_columns_length-2
		date_updated_index := table_columns_length-1
		if len(_tableSlice) > 2 && _tableSlice[date_created_index].Field == "date_created" && _tableSlice[date_updated_index].Field == "date_updated" {
			// _tableSlice = append(_tableSlice[:date_created_index], _tableSlice[date_created_index+2:]...)
			utility.LogBlue("preserving timestamps if exist")
		}
	}
	
	a, _ := json.Marshal(&_tableSlice)
	_ = json.Unmarshal(a, &s.savedTableMap)
	s.savedTableJsonString = string(a)
}