//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/mysql"
)

var Author = newAuthorTable("posts", "Author", "")

type authorTable struct {
	mysql.Table

	// Columns
	ID   mysql.ColumnInteger
	Name mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type AuthorTable struct {
	authorTable

	NEW authorTable
}

// AS creates new AuthorTable with assigned alias
func (a AuthorTable) AS(alias string) *AuthorTable {
	return newAuthorTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new AuthorTable with assigned schema name
func (a AuthorTable) FromSchema(schemaName string) *AuthorTable {
	return newAuthorTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new AuthorTable with assigned table prefix
func (a AuthorTable) WithPrefix(prefix string) *AuthorTable {
	return newAuthorTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new AuthorTable with assigned table suffix
func (a AuthorTable) WithSuffix(suffix string) *AuthorTable {
	return newAuthorTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newAuthorTable(schemaName, tableName, alias string) *AuthorTable {
	return &AuthorTable{
		authorTable: newAuthorTableImpl(schemaName, tableName, alias),
		NEW:         newAuthorTableImpl("", "new", ""),
	}
}

func newAuthorTableImpl(schemaName, tableName, alias string) authorTable {
	var (
		IDColumn       = mysql.IntegerColumn("ID")
		NameColumn     = mysql.StringColumn("name")
		allColumns     = mysql.ColumnList{IDColumn, NameColumn}
		mutableColumns = mysql.ColumnList{NameColumn}
	)

	return authorTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:   IDColumn,
		Name: NameColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
