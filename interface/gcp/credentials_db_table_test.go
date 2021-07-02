package gcp

import (
	"database/sql"
	"io/ioutil"
	"leapp_daemon/test"
	"os"
	"path"
	"reflect"
	"testing"
)

type CredentialsRow struct {
	accountId string
	value     string
}

var (
	tempCredentialsDbPath string
)

func credentialsDbTableSetup() {
	tempCredentialsDbPath = path.Join(os.TempDir(), "sql_lite_test.db")
}

func expectCredentials(t *testing.T, expectedRows []CredentialsRow) {
	database, _ := sql.Open("sqlite3", tempCredentialsDbPath)
	defer database.Close()

	row, _ := database.Query("SELECT * FROM credentials ORDER BY account_id ASC")
	defer row.Close()

	foundRows := []CredentialsRow{}
	var dbAccountId string
	var dbValue string
	for row.Next() {
		row.Scan(&dbAccountId, &dbValue)
		foundRows = append(foundRows, CredentialsRow{accountId: dbAccountId, value: dbValue})
	}

	if !reflect.DeepEqual(foundRows, expectedRows) {
		t.Errorf("expected rows:\n%v", expectedRows)
	}
}

func TestWriteCredentials(t *testing.T) {
	credentialsDbTableSetup()

	accountId := "account_id@domain.com"
	anotherAccountId := "another_account@domain.com"
	value := "jsonCredentials"
	anotherValue := "anotherJsonCredentials"

	createdFile, _ := os.Create(tempCredentialsDbPath)
	createdFile.Close()

	table := GcpCredentialsTable{}
	table.WriteCredentials(tempCredentialsDbPath, accountId, value)
	expectCredentials(t, []CredentialsRow{{accountId: accountId, value: value}})

	table.WriteCredentials(tempCredentialsDbPath, anotherAccountId, value)
	expectCredentials(t, []CredentialsRow{{accountId: accountId, value: value}, {accountId: anotherAccountId, value: value}})

	table.WriteCredentials(tempCredentialsDbPath, accountId, anotherValue)
	expectCredentials(t, []CredentialsRow{{accountId: accountId, value: anotherValue}, {accountId: anotherAccountId, value: value}})
}

func TestWriteCredentials_CorruptedDbFile(t *testing.T) {
	credentialsDbTableSetup()
	os.Remove(tempCredentialsDbPath)
	ioutil.WriteFile(tempCredentialsDbPath, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0644)

	accountId := "account_id@domain.com"
	value := "jsonCredentials"

	table := GcpCredentialsTable{}
	err := table.WriteCredentials(tempCredentialsDbPath, accountId, value)
	test.ExpectHttpError(t, err, 500, "file is not a database")
}

func TestRemoveCredentials(t *testing.T) {
	credentialsDbTableSetup()

	accountId := "account_id@domain.com"
	anotherAccountId := "another_account@domain.com"
	value := "jsonCredentials"
	anotherValue := "anotherJsonCredentials"

	createdFile, _ := os.Create(tempCredentialsDbPath)
	createdFile.Close()

	table := GcpCredentialsTable{}
	table.WriteCredentials(tempCredentialsDbPath, accountId, value)
	table.WriteCredentials(tempCredentialsDbPath, anotherAccountId, anotherValue)

	table.RemoveCredentials(tempCredentialsDbPath, anotherAccountId)
	expectCredentials(t, []CredentialsRow{{accountId: accountId, value: value}})
}

func TestRemoveCredentials_MissingDbFile(t *testing.T) {
	credentialsDbTableSetup()
	os.Remove(tempCredentialsDbPath)

	table := GcpCredentialsTable{}
	err := table.RemoveCredentials(tempCredentialsDbPath, "account_id@domain.com")
	test.ExpectHttpError(t, err, 400, "no such table: credentials")
}
