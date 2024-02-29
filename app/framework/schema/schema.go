package schema

import (
	"coralscale/app/framework/db"
	"coralscale/app/framework/utility"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
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
		// var _savedTableMap map[string]string
		rows.Scan(&_table.Field, &_table.Type, &_table.Null, &_table.Key, &_table.Default, &_table.Extra)
		savedTableSlice = append(savedTableSlice, _table)

		// b, _ := json.Marshal(&_table)
		// _ = json.Unmarshal(b, &_savedTableMap)

		// savedTableMap = append(savedTableMap, _savedTableMap)
	}

	a, _ := json.Marshal(&savedTableSlice)
	_ = json.Unmarshal(a, &savedTableMap)
	
	b, _ := json.Marshal(&newTableSlice)
	_ = json.Unmarshal(b, &newTableMap)

	savedTableJsonString = string(a)
	newTableJsonString = string(b)

	fmt.Println(savedTableSlice)
	fmt.Println(savedTableMap)

	utility.LogSuccess(savedTableJsonString)
	utility.LogNeutral(newTableJsonString)
	return true
}
