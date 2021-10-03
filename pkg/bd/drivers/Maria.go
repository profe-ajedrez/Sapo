package drivers

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/profe-ajedrez/sapo/pkg/bd/groundwork"
	"github.com/profe-ajedrez/sapo/pkg/json/formatter"
	"github.com/profe-ajedrez/sapo/pkg/json/parser/file"
)

const SQL_SHOW_COLUMNS = "SHOW COLUMNS FROM %TABLE%"
const SQL_REFERENCERS = `SELECT referenced_table_name as table_name
	FROM information_schema.key_column_usage
	WHERE referenced_table_name IS NOT NULL
	  AND table_name = '%TABLE%' and table_schema = '%SCHEME%';`

const SQL_REFERENCED = `SELECT DISTINCT table_name AS table_name
	FROM information_schema.key_column_usage
	WHERE referenced_table_name = '%TABLE%' AND table_schema = '%SCHEME%';
    `

type Maria struct {
	host    string `default:""`
	user    string `default:""`
	pass    string `default:""`
	port    string `default:""`
	schema  string `default:""`
	options string `default:""`
	db      *sql.DB
}

type Columns struct {
	Field   string
	Type    string
	Null    sql.NullString
	Key     string
	Default sql.NullString
	Auto    sql.NullString
}

var ExtendedQuery bool
var ColumnToSearch string

func init() {
	ExtendedQuery = false
	ColumnToSearch = ""
}

func (db *Maria) Connect(configFilePath string, maxOpenConns int, maxIdleConns int, maxLifetime time.Duration) error {
	jsonMap, err := file.ReadJSON(configFilePath)

	if err == nil {
		db.host = string(jsonMap["host"])
		db.user = string(jsonMap["username"])
		db.pass = string(jsonMap["password"])
		db.port = string(jsonMap["port"])
		db.schema = string(jsonMap["schema"])
		db.options = string(jsonMap["options"])

		var dsn string = strings.TrimSpace(db.user) + ":" +
			strings.TrimSpace(db.pass) + "@tcp(" +
			strings.TrimSpace(db.host) + ":" +
			strings.TrimSpace(db.port) + ")/" +
			strings.TrimSpace(db.schema) + "?" +
			strings.TrimSpace(db.options)

		var err error
		db.db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}

		err = db.db.Ping()
		if err == nil {
			db.db.SetConnMaxLifetime(maxLifetime)
			db.db.SetMaxOpenConns(maxOpenConns)
			db.db.SetMaxIdleConns(maxIdleConns)
		}
	}

	return err
}

func (db *Maria) Close() error {
	return db.db.Close()
}

func (db *Maria) Relations(table string) (groundwork.Relation, error) {
	referenced, err := referenced(db, table)

	if err == nil {
		referencers, err := referencers(db, table)
		if err == nil {
			relation := groundwork.Relation{
				Table:    table,
				Refered:  referenced,
				Refering: referencers,
			}
			return relation, err
		}
	}

	return groundwork.Relation{}, err
}

func (db *Maria) GetStructure(table string) (string, error) {
	result, err := getStructureResult(db, table)
	jsonStruct := "\"\""

	if err == nil {
		jsonStruct, _ = formatter.ToPrettyJson(result)
	}
	return jsonStruct, err
}

func getStructureResult(db *Maria, table string) ([]Columns, error) {
	strSql := strings.ReplaceAll(SQL_SHOW_COLUMNS, "%TABLE%", table)
	results, err := db.db.Query(strSql)
	columns := []Columns{}

	if err == nil {
		row := Columns{}
		for results.Next() {
			err = results.Scan(&row.Field, &row.Type, &row.Null, &row.Key, &row.Default, &row.Auto)
			columns = append(columns, row)
		}

	}

	return columns, err
}

func referenced(db *Maria, table string) ([]string, error) {
	strSql := strings.ReplaceAll(strings.ReplaceAll(SQL_REFERENCED, "%SCHEME%", db.schema), "%TABLE%", table)
	return querier(db, strSql)
}

func referencers(db *Maria, table string) ([]string, error) {
	var strSql string
	if ExtendedQuery && len(ColumnToSearch) > 0 {
		strSql = sqlExtended(db.schema)
	} else {
		strSql = strings.ReplaceAll(strings.ReplaceAll(SQL_REFERENCERS, "%SCHEME%", db.schema), "%TABLE%", table)
	}
	referencers, err := querier(db, strSql)
	if err == nil {
		return referencers, nil
	}

	return referencers, err
}

func querier(db *Maria, strSql string) ([]string, error) {
	results, err := db.db.Query(strSql)
	var relatives []string

	relative := ""
	for results.Next() {
		err = results.Scan(&relative)
		if err != nil {
			break
		}
		relatives = append(relatives, relative)
	}

	return relatives, err
}

func sqlExtended(scheme string) string {
	return fmt.Sprintf(`SELECT DISTINCT table_name 
	FROM INFORMATION_SCHEMA.COLUMNS
	WHERE COLUMN_NAME IN ('%s')
	AND TABLE_SCHEMA='%s';`, ColumnToSearch, scheme)
}
