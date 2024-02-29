package schema

import (
	 "coralscale/app/framework/db"
	"coralscale/app/framework/utility"
	"strconv"
)

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
		Type: "bigint",
		Size: 11,
	})

	append_comma := ""

	if len(s.columns) > 1 {
		append_comma = ","
	}

	s.schema_sql = s.schema_sql + append_comma + " `id` bigint NOT NULL AUTO_INCREMENT"
	s.primary_key = ", PRIMARY KEY (id)"
}

func (s *Schema) String(name string, size int) {
	s.columns = append(s.columns, Column{
		Name: name,
		Type: "string",
		Size: size,
	})

	append_comma := ""

	if len(s.columns) > 1 {
		append_comma = ","
	}

	s.schema_sql = s.schema_sql + append_comma + " `" + name + "` varchar(" + strconv.Itoa(size) + ") NOT NULL"
}

func (s *Schema) Create() {
	create_table_query := "CREATE TABLE IF NOT EXISTS " + s.TableName + " (" + s.schema_sql + s.primary_key + ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"
	// utility.LogNeutral(create_table_query)
	db.Exec(create_table_query)

	utility.LogSuccess(s.TableName + " created successfully")
}
