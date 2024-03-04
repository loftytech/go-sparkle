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

var newTableSlice []TableStruct
var newTableMap []map[string]string

var savedTableSlice []TableStruct
var savedTableMap []map[string]string

var savedTableJsonString string
var newTableJsonString string

type Schema struct {
	TableName   string
	columns     []Column
	schema_sql  string
	primary_key string
	skippedMigration bool
	modificationList []map[string]string
	columnDropList []map[string]string
	replacedColumnList []string
}

type Column struct {
	Name string
	Size int
	Type string
}

func (s *Schema) Id() {
	s.columns = append(s.columns, Column{
		Name: "id",
		Type: "bigint(20)",
		Size: 11,
	})

	append_comma := ""

	if len(s.columns) > 1 {
		append_comma = ","
	}

	newTableSlice = append(newTableSlice, TableStruct{
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

func (s *Schema) String(name string, size int) {
	s.columns = append(s.columns, Column{
		Name: name,
		Type: "varchar",
		Size: size,
	})

	append_comma := ""

	if len(s.columns) > 1 {
		append_comma = ","
	}

	newTableSlice = append(newTableSlice, TableStruct{
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

func (s *Schema) Create() {
	database = db.GetDBInstance()
	is_exists := s.checkTabeExist()

	if !is_exists {
		s.createTable()
	} else {
		if newTableJsonString == savedTableJsonString {
			utility.LogNeutral("Skipping "+s.TableName+" no changes")
		} else {
			s.compareColumns()
		}
	}
}

func (s *Schema) createTable() {
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
		savedTableSlice = append(savedTableSlice, _table)
	}

	a, _ := json.Marshal(&savedTableSlice)
	_ = json.Unmarshal(a, &savedTableMap)
	
	b, _ := json.Marshal(&newTableSlice)
	_ = json.Unmarshal(b, &newTableMap)

	savedTableJsonString = string(a)
	newTableJsonString = string(b)

	// fmt.Println(savedTableSlice)
	// fmt.Println(savedTableMap)

	// utility.LogSuccess(savedTableJsonString)
	// utility.LogNeutral(newTableJsonString)
	return true
}


func (s *Schema) compareColumns() {
	if len(savedTableMap) == len(newTableMap) {
		s.compareAndAlterColumns()
	} else {
		s.getNewColumns()
	}
	s.migrateSchema()
}

func (s *Schema) compareAndAlterColumns() {
	is_columns_match := false
	if savedTableJsonString == newTableJsonString {
		is_columns_match = true
	}

	if is_columns_match {
		s.skippedMigration = true
		utility.LogNeutral("Skipping " + s.TableName + " table because no changes was found")
	} else {
		for newColumnkey, column := range newTableMap {
			new_column := s.getColumnType(column["Type"])

			add_primary_key := "";
			auto_increment_sub_query := "";

			saved_column_to_alter := savedTableMap[newColumnkey]

			
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
				mod_sql = "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + column["Field"] + "` DATETIME NOT NULL DEFAULT " + column["Default"] + "; "
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

	saved_table_length := len(savedTableMap);
	for newColumnkey, column := range newTableMap {
		new_column := s.getColumnType(column["Type"])

		auto_increment_sub_query := ""
		add_primary_key := ""

		if newColumnkey < saved_table_length {
			saved_column_to_alter := savedTableMap[newColumnkey]
			
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
				mod_sql = "ALTER TABLE `" + s.TableName + "` CHANGE `" + saved_column_to_alter["Field"] + "` `" + column["Field"] + "` DATETIME NOT NULL DEFAULT " + column["Default"] + "; "
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
				previous_column_field := s.getPreviousColumn(newTableMap, newColumnkey)["Field"];
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
	field_arry := strings.Split("(", field);
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
	new_table_length := len(newTableMap)
	if new_table_length < len(savedTableMap) {
		utility.LogWarning("Some columns would be dropped");
		colums_to_be_dropped := savedTableMap;
		filtered_column_list := []map[string]string{};

		for savedColumnKey, savedColumn := range savedTableMap {
			/*
			*
			*   Check for columns that already exists in the table
			*
			*/

			if savedColumnKey < new_table_length {
				// Log::success($savedColumn["Field"] ." will not be droped");
				filtered_column_list = append(filtered_column_list, savedColumn)
				colums_to_be_dropped = append(colums_to_be_dropped[:savedColumnKey], colums_to_be_dropped[savedColumnKey+1:]...)
			}
		}

		for _, colum_to_be_dropped := range colums_to_be_dropped {
			s.modificationList = append(s.modificationList, map[string]string{
				"sql": "ALTER TABLE `" + s.TableName + "` DROP `" + colum_to_be_dropped["Field"] + "; ",
				"operation_type": "ALTER_DROP",
			})
		}

		savedTableMap = filtered_column_list
	}
}

func (s *Schema) getPreviousColumn(arr []map[string]string, key int) map[string]string {
	return arr[key-1];
}

func (s *Schema) migrateSchema()  {
	for _, column_drop_query := range s.columnDropList {
		// Log::warning($column_drop_query["sql"]);
		db.Exec(column_drop_query["sql"])
	}
	for _, modification := range s.modificationList {
		if modification["operation_type"] == "ALTER_CHANGE_COLUMN" {
			if in_array(modification["from"], s.replacedColumnList) {
				checkedColumn := s.checkColumnExists(modification["to"])
				if checkedColumn["exists"] {
					if modification["from"] != modification["to"] {
						s.replacedColumnList = append(s.replacedColumnList, modification["to"])
						column = s.getColumnType(checkedColumn["data"].Type)
						s.buildAndRunQuery(QueryBuildStrct{
							QueryType: "ALTER_CHANGE_COLUMN", 
							From: modification["to"], 
							To: modification["to"] + "_ALTERED", 
							Type: column.Type, 
							Limit: column.Limit, 
							FromType: modification["from_type"], 
							ToType: modification["to_type"], 
							Extras: checkedColumn["data"].Extra,
						})
					}
				}
				// strings.Replace(field_arry[1], ")", "", -1),
				new_sql := strings.Replace( "CHANGE `" + modification["from"] + "_ALTERED`", modification["sql"], "CHANGE `" + modification["from"] + "`", -1);
				db.Exec(new_sql)
			} else {
				checkedColumn = s.checkColumnExists(modification["to"])
				if checkedColumn["exists"] {
					if modification["from"] != modification["to"] {
						s.replacedColumnList = append(s.replacedColumnList, modification["to"])
						column = s.getColumnType(checkedColumn.data.Type);

						s.buildAndRunQuery(QueryBuildStrct{
							QueryType: "ALTER_CHANGE_COLUMN", 
							From: modification["to"], 
							To: modification["to"] + "_ALTERED", 
							Type: column.Type, 
							Limit: column.Limit, 
							FromType: modification["from_type"], 
							ToType: modification["to_type"], 
							Extras: checkedColumn["data"].Extra,
						})

						s.clearColumn(ClearFormStruct{
							from: modification["from"],
							to: modification["to"],
							fromType: modification["from_type"],
							toType: modification["to_type"],
						});
						// Log::neutral(modification["sql"]);
						db.Exec(modification["sql"]);
					} else {
						s.clearColumn(ClearFormStruct{
							from: modification["from"],
							to: modification["to"],
							fromType: modification["from_type"],
							toType: modification["to_type"],
						});
						// Log::neutral(modification["sql"]);
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

			// Log::neutral(modification["sql"]);
		} else {
			// Log::neutral(modification["sql"]);
			db.Exec(modification["sql"]);
		}

	}

	if !s.skippedMigration {
		utility.LogSuccess(s.TableName + " table migrated successfully")
	}
}

func in_array(target map[string]string, slice []map[string]string) bool {
    for _, m := range slice {
        if mapsEqual(m, target) {
            return true
        }
    }
    return false
}

func mapsEqual(a, b map[string]string) bool {
    if len(a) != len(b) {
        return false
    }
    for key, valA := range a {
        if valB, ok := b[key]; !ok || valA != valB {
            return false
        }
    }
    return true
}

type CheckColumeReturnStruct struct {
	exists bool
	data TableStruct
}

func (s *Schema) checkColumnExists (column map[string]string) CheckColumeReturnStruct {
	query := "DESCRIBE " + s.TableName;
	
	var _tableSlice []TableStruct
	var _tableMap []map[string]string
	colunm_data := TableStruct{};
	column_exists := false;

	rows, error := database.Query(query)
	if error != nil {
		// return false
	}

	for rows.Next() {
		var _table TableStruct
		rows.Scan(&_table.Field, &_table.Type, &_table.Null, &_table.Key, &_table.Default, &_table.Extra)
		_tableSlice = append(_tableSlice, _table)
	}
	
	
	b, _ := json.Marshal(&_tableSlice)
	_ = json.Unmarshal(b, &_tableMap)

	for savedColumnKey, savedColumn := range _tableMap {
		if column["Field"] == savedColumn["field"] {
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