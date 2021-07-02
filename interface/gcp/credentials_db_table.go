package gcp

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"leapp_daemon/infrastructure/http/http_error"
)

var credentialsTableName = "credentials"

type GcpCredentialsTable struct {
}

func (table *GcpCredentialsTable) WriteCredentials(sqlFilePath string, accountId string, value string) error {
	database, err := getSqliteDatabase(sqlFilePath)
	if err != nil {
		return err
	}
	defer database.Close()

	err = table.createTable(database)
	if err != nil {
		return err
	}

	err = table.insertQuery(database, accountId, value)
	if err != nil {
		return err
	}

	return nil
}

func (table *GcpCredentialsTable) RemoveCredentials(sqlFilePath string, accountId string) error {
	database, err := getSqliteDatabase(sqlFilePath)
	if err != nil {
		return err
	}
	defer database.Close()

	err = table.deleteQuery(database, accountId)
	if err != nil {
		return err
	}

	return nil
}

func (table *GcpCredentialsTable) createTable(database *sql.DB) error {
	createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%v" (
	  "account_id"	TEXT,
	  "value"	BLOB,
	  PRIMARY KEY("account_id")
  );`, credentialsTableName)

	statement, err := database.Prepare(createTableQuery)
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	_, err = statement.Exec()
	if err != nil {
		return http_error.NewInternalServerError(err)
	}

	return nil
}

func (table *GcpCredentialsTable) insertQuery(database *sql.DB, accountId string, value string) error {
	insertQuery := fmt.Sprintf("INSERT INTO %v (account_id, value) VALUES(?, ?) ON CONFLICT(account_id) DO UPDATE SET value=excluded.value", credentialsTableName)
	statement, err := database.Prepare(insertQuery)
	if err != nil {
		return http_error.NewBadRequestError(err)
	}

	_, err = statement.Exec(accountId, value)
	if err != nil {
		return http_error.NewBadRequestError(err)
	}

	return nil
}

func (table *GcpCredentialsTable) deleteQuery(database *sql.DB, accountId string) error {
	deleteQuery := fmt.Sprintf("DELETE FROM %v WHERE account_id=?", credentialsTableName)
	statement, err := database.Prepare(deleteQuery)
	if err != nil {
		return http_error.NewBadRequestError(err)
	}

	_, err = statement.Exec(accountId)
	if err != nil {
		return http_error.NewBadRequestError(err)
	}

	return nil
}
