# GoDBMSManager

GoDBMSManager is a lightweight and extensible database management system (DBMS) manager written in Golang. It simplifies common database operations, such as creating and managing databases and tables, while providing a structured approach to defining and working with table structures. This tool is specifically designed to work with MySQL databases.

## Table of Contents

- [Features](#features)
  - [Implemented Features](#implemented-features)
  - [Unimplemented Features](#unimplemented-features)
- [Getting Started](#getting-started)
- [Usage](#usage)
  - [Creating a Database](#creating-a-database)
  - [Creating a Table](#creating-a-table)
  - [Displaying Table Structure](#displaying-table-structure)
  - [Dropping a Table](#dropping-a-table)
  -
## Features

### Implemented Features

1. **Database Operations**
   - Create a database if it doesn't exist.
   - Drop an existing database.

   **Example:**
   ```go
   createDatabase("mydatabase", db)
   dropDatabase("mydatabase", db)
   ```

2. **Table Operations**
   - Create a table based on a provided JSON structure.
   - Display the structure of a table in both JSON and SQL formats.
   - Drop an existing table.

   **Example:**
   ```go
    roleTable := Table{
    name: "roles",
    columns: []TableColumn{
        {Name: "role_id", KeyType: PRIMARY_KEY, ValueType: INT_TYPE, AutoIncrement: true},
        {Name: "role_name", ValueType: VARCHAR_TYPE, Unique: true},
        // ... other columns
      },
    }

    createTable(roleTable, "mydatabase", db)

   displayTable(roleTable)      // Display JSON structure
   displayTableSQL(roleTable)   // Display SQL structure
   ```

3.  **Column Types**

    - Supports various column types, including INT, VARCHAR, CHAR, TEXT, FLOAT, DOUBLE, DECIMAL, DATE, and DATETIME.
    - Auto-increment for primary key columns.
    - Unique constraint for columns.
    - Default values for columns.
    - Foreign key relationships between tables.


4. **JSON Structure**

    - Read and parse table structure from JSON files.
    - Write table structure to JSON files.


  ### Unimplemented Features

 -  [ ] Alter Table
 -  [ ] Modify the structure of an existing table (e.g., add, drop, or modify columns).
 -  [ ] Data Operations
 -  [ ] Insert, update, and delete data from tables.
 -  [ ] Indexing
 -  [ ] Create and manage indexes on table columns.
 -  [ ] Query Execution
 -  [ ] Execute custom SQL queries on the database.

## Getting Started
### Clone the repository:

```bash
git clone https://github.com/yourusername/godbmsmanager.git
Ensure you have Go installed: https://golang.org/dl/
```

### Install dependencies:

```bash
go get -u github.com/go-sql-driver/mysql
Configure the database connection details in the main.go file (DB_USER, DB_PASS, DB_NAME).
```

### Run the application:

```bash
go run main.go
```

## Usage

### Creating a Database
To create a new database, use the createDatabase function and provide the desired database name:

```go
createDatabase("mydatabase", db)
```

### Creating a Table
Define your table structure using the Table struct and then use the createTable function to create the table:

```go
roleTable := Table{
    name: "roles",
    columns: []TableColumn{
        {Name: "role_id", KeyType: PRIMARY_KEY, ValueType: INT_TYPE, AutoIncrement: true},
        {Name: "role_name", ValueType: VARCHAR_TYPE, Unique: true},
        // ... other columns
    },
}

createTable(roleTable, "mydatabase", db)
```

### Displaying Table Structure

#### Display the JSON and SQL structure of a table:

```go
displayTable(roleTable)      // Display JSON structure
displayTableSQL(roleTable)   // Display SQL structure
```

Dropping a Table
To drop a table, use the dropTable function:

```go
dropTable("roles", db)
```

