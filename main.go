package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type TableColumn struct {
	Name          string      `json:"name"`
	KeyType       KEY_TYPE    `json:"key_type"`
	ValueType     VALUE_TYPE  `json:"value_type"`
	ForeignKey    *ForeignKey `json:"foreign_key,omitempty"`
	Unique        bool        `json:"unique,omitempty"`
	Check         string      `json:"check,omitempty"`
	Default       string      `json:"default,omitempty"`
	AutoIncrement bool        `json:"auto_increment,omitempty"`
}

type ForeignKey struct {
	ReferenceTable  string `json:"reference_table"`
	ReferenceColumn string `json:"reference_column"`
}

type Table struct {
	name    string
	columns []TableColumn
}

type STATEMENT string
type KEY_TYPE string
type VALUE_TYPE string

const (
	PRIMARY_KEY KEY_TYPE = "PRIMARY KEY"
	UNIQUE_KEY  KEY_TYPE = "UNIQUE KEY"
	FOREIGN_KEY KEY_TYPE = "FOREIGN KEY"
)

const (
	INT_TYPE      VALUE_TYPE = "INT"
	VARCHAR_TYPE  VALUE_TYPE = "VARCHAR(255)"
	CHAR_TYPE     VALUE_TYPE = "CHAR"
	TEXT_TYPE     VALUE_TYPE = "TEXT"
	FLOAT_TYPE    VALUE_TYPE = "FLOAT"
	DOUBLE_TYPE   VALUE_TYPE = "DOUBLE"
	DECIMAL_TYPE  VALUE_TYPE = "DECIMAL"
	DATE_TYPE     VALUE_TYPE = "DATE"
	DATETIME_TYPE VALUE_TYPE = "DATETIME"
)

const (
	CREATE_DATABASE STATEMENT = "CREATE DATABASE IF NOT EXISTS "
	DROP_DATABASE   STATEMENT = "DROP DATABASE IF EXISTS "
)

const (
	CREATE_TABLE STATEMENT = "CREATE TABLE "
	ALTER_TABLE  STATEMENT = "ALTER TABLE "
	DROP_TABLE   STATEMENT = "DROP TABLE "
	USE          STATEMENT = "USE "
)

const (
	DB_NAME = "rest"
	DB_USER = "root"
	DB_PASS = "admin"
)

const STRUCTURE_PATH = "./jsons/structures/"

func createStatement(statement string, statementType STATEMENT) string {
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println("Statement is:")
	fmt.Printf("%v%s\n", statementType, statement)
	fmt.Println("--------------------------------------------------------------------------------------------")
	return string(statementType) + statement
}

func createTable(table Table, database string, db *sql.DB) {
	use_query := createStatement(database, USE)
	_, err := db.Exec(use_query)
	if err != nil {
		log.Fatal(err)
	}

	createStatement := parseTableJSON(table.name)

	_, error := db.Exec(createStatement)

	if error != nil {
		log.Fatal(error)
	}

}

func displayTable(table Table) {
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println("Displaying Structure For Table: ", table.name)
	fmt.Println("--------------------------------------------------------------------------------------------")
	jsonData, err := json.MarshalIndent(table.columns, "", "    ")
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
	fmt.Println("--------------------------------------------------------------------------------------------")
}

func displayTableSQL(table Table) {
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println("Displaying SQL For Table: ", table.name)
	fmt.Println("--------------------------------------------------------------------------------------------")
	parseTableJSON(table.name)
	fmt.Println("--------------------------------------------------------------------------------------------")
}

func createDatabase(name string, db *sql.DB) {
	query := createStatement(name, CREATE_DATABASE)
	fmt.Println(query)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func dropDatabase(name string, db *sql.DB) {
	query := createStatement(name, DROP_DATABASE)
	fmt.Println(query)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func dropTable(name string, db *sql.DB) {
	query := createStatement(name, DROP_TABLE)
	fmt.Println(query)
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func parseTableJSON(tableName string) string {
	fmt.Println("--------------------------------------------------------------------------------------------")
	fmt.Println("Parsing Table JSON")
	fmt.Println("--------------------------------------------------------------------------------------------")

	tableData, err := os.ReadFile(STRUCTURE_PATH + tableName + ".json")
	if err != nil {
		log.Fatalf("Error reading file %s.json: %v", tableName, err)
		return "err"
	}

	var columns []TableColumn
	err = json.Unmarshal(tableData, &columns)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON for table %s: %v", tableName, err)
		return "err"
	}

	createTableStatement := createTableSQL(tableName, columns)
	fmt.Println("Create Table Statement:")
	fmt.Println(createTableStatement)

	fmt.Println("--------------------------------------------------------------------------------------------")

	return createTableStatement
}

func createTableSQL(tableName string, columns []TableColumn) string {
	var statements []string

	for _, column := range columns {
		statement := fmt.Sprintf("%s %s", column.Name, column.ValueType)

		if column.KeyType != "" {
			statement += " " + string(column.KeyType)
		}

		if column.AutoIncrement {
			statement += " AUTO_INCREMENT"
		}

		if column.Unique {
			statement += " UNIQUE"
		}

		if column.Default != "" {
			statement += " DEFAULT " + column.Default
		}

		// FOREIGN KEY (PersonID) REFERENCES Persons(PersonID)
		// this is the format for foreign keys fucking pain to do, dont change below code

		if column.ForeignKey != nil {
			statement = strings.Trim(statement, fmt.Sprintf("%s %s", column.Name, column.ValueType))
			statement = strings.Trim(statement, "FOREIGN KEY")
			statement += fmt.Sprintf(" %s %s,\n", column.Name, column.ValueType)
			statement += fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", column.Name, column.ForeignKey.ReferenceTable, column.ForeignKey.ReferenceColumn)
		}

		statements = append(statements, statement)
	}

	createTableStatement := fmt.Sprintf("CREATE TABLE %s (\n\t%s\n);", tableName, strings.Join(statements, ",\n\t"))

	return createTableStatement
}

func writeTableJSON(dir string, table Table) {
	fmt.Println("--------------------------------------------------------------------------------------------")
	filename := table.name + ".json"
	finalDirectory := dir + filename

	f, err := os.Create(finalDirectory)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	jsonData, err := json.MarshalIndent(table.columns, "", "    ")
	if err != nil {
		log.Fatal("Error marshalling JSON:", err)
		return
	}

	err = os.WriteFile(finalDirectory, jsonData, 0666)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
	fmt.Println("--------------------------------------------------------------------------------------------")
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", DB_USER, DB_PASS, DB_NAME))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	roleTable := Table{
		name: "roles",
		columns: []TableColumn{
			{Name: "role_id", KeyType: PRIMARY_KEY, ValueType: INT_TYPE, AutoIncrement: true},
			{Name: "role_name", ValueType: VARCHAR_TYPE, Unique: true},
			{Name: "created_at", ValueType: DATETIME_TYPE},
			{Name: "updated_at", ValueType: DATETIME_TYPE},
		},
	}

	userTable := Table{
		name: "users",
		columns: []TableColumn{
			{Name: "user_id", KeyType: PRIMARY_KEY, ValueType: INT_TYPE, AutoIncrement: true},
			{Name: "username", ValueType: VARCHAR_TYPE, Unique: true},
			{Name: "email", ValueType: VARCHAR_TYPE, Unique: true},
			{Name: "created_at", ValueType: DATETIME_TYPE},
			{Name: "updated_at", ValueType: DATETIME_TYPE},
			{Name: "role_id", KeyType: FOREIGN_KEY, ValueType: INT_TYPE, ForeignKey: &ForeignKey{ReferenceTable: "roles", ReferenceColumn: "role_id"}},
		},
	}

	fmt.Println()

	writeTableJSON("./jsons/structures/", roleTable)
	writeTableJSON("./jsons/structures/", userTable)

	createDatabase("test_db_1", db)

	createTable(roleTable, "test_db_1", db)
	createTable(userTable, "test_db_1", db)

	displayTable(roleTable)
	displayTableSQL(roleTable)
	displayTable(userTable)
	displayTableSQL(userTable)

    dropTable(userTable.name, db)
    dropTable(roleTable.name, db)

	// parseTableJSON("roles")
	// parseTableJSON("users")

}

