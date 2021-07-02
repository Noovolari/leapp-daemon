package gcp

import (
	"fmt"
	"leapp_daemon/test"
	"os"
	"path"
	"testing"
)

var (
	tempAccessTokensDbPath string
)

func accessTokensTableSetup() {
	tempAccessTokensDbPath = path.Join(os.TempDir(), "sql_lite_test.db")
}

func TestRemoveAccessToken(t *testing.T) {
	accessTokensTableSetup()

	createdFile, _ := os.Create(tempAccessTokensDbPath)
	createdFile.Close()

	database, err := getSqliteDatabase(tempAccessTokensDbPath)
	if err != nil {
		return
	}

	createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%v" (
    "account_id"	TEXT,
    "access_token"	TEXT,
    "token_expiry"	TIMESTAMP,
    "rapt_token"	TEXT,
    "id_token"	TEXT,
    PRIMARY KEY("account_id")
  );`, accessTokensTableName)

	statement, err := database.Prepare(createTableQuery)
	if err != nil {
		t.Fatalf("Unable to prepare table for the test 1/2")
	}

	_, err = statement.Exec()
	if err != nil {
		t.Fatalf("Unable to prepare table for the test 2/2")

	}

	insertQuery := fmt.Sprintf("INSERT INTO %v (account_id, access_token) VALUES(?, ?)", accessTokensTableName)
	statement, err = database.Prepare(insertQuery)
	if err != nil {
		t.Fatalf("Unable to insert data into the table for the test 1/2")
	}

	accountId := "account_id@domain.com"
	_, err = statement.Exec(accountId, "access-token")
	if err != nil {
		t.Fatalf("Unable to insert data into the table for the test 2/2")
	}
	database.Close()

	table := GcpAccessTokensTable{}
	table.RemoveAccessToken(tempAccessTokensDbPath, accountId)

	database, err = getSqliteDatabase(tempAccessTokensDbPath)
	defer database.Close()

	row, _ := database.Query(fmt.Sprintf("SELECT * FROM %v", accessTokensTableName))
	defer row.Close()

	if row.Next() {
		t.Errorf("expected no rows in the table")
	}
}

func TestRemoveAccessToken_MissingDbFile(t *testing.T) {
	accessTokensTableSetup()
	os.Remove(tempAccessTokensDbPath)

	table := GcpAccessTokensTable{}
	err := table.RemoveAccessToken(tempAccessTokensDbPath, "account_id@domain.com")
	test.ExpectHttpError(t, err, 400, "no such table: access_tokens")
}
