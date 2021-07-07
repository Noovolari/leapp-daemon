package repository

import (
	"fmt"
	"leapp_daemon/test"
	"leapp_daemon/test/mock"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	expectedHomeDirPath                string
	expectedConfigDirPath              string
	expectedConfigFilePath             string
	expectedActiveConfigFilePath       string
	expectedDefaultCredentialsFilePath string
	expectedCredentialsDbFilePath      string
	expectedAccessTokensDbFilePath     string
	gcpEnvMock                         mock.EnvironmentMock
	gcpFsMock                          mock.FileSystemMock
	credentialsTableMock               mock.GcpCredentialsDbTableMock
	accessTokensTableMock              mock.GcpAccessTokensDbTableMock
	repo                               *GcpConfigurationRepository
)

func gcpConfigurationRepositorySetup() {
	expectedHomeDirPath = filepath.Join("c:/", "appdata")
	expectedConfigDirPath = filepath.Join(expectedHomeDirPath, "gcloud")
	expectedConfigFilePath = filepath.Join(expectedConfigDirPath, "configurations", "config_leapp")
	expectedActiveConfigFilePath = filepath.Join(expectedConfigDirPath, "active_config")
	expectedDefaultCredentialsFilePath = filepath.Join(expectedConfigDirPath, "application_default_credentials.json")
	expectedCredentialsDbFilePath = filepath.Join(expectedConfigDirPath, "credentials.db")
	expectedAccessTokensDbFilePath = filepath.Join(expectedConfigDirPath, "access_tokens.db")

	gcpEnvMock = mock.NewEnvironmentMock()
	gcpFsMock = mock.NewFileSystemMock()
	credentialsTableMock = mock.NewGcpCredentialsDbTableMock()
	accessTokensTableMock = mock.NewGcpAccessTokensDbTableMock()

	repo = &GcpConfigurationRepository{
		Environment:       &gcpEnvMock,
		FileSystem:        &gcpFsMock,
		CredentialsTable:  &credentialsTableMock,
		AccessTokensTable: &accessTokensTableMock,
	}
}

func verifyExpectedCalls(t *testing.T, envMockCalls []string, fsMockCalls []string, credentialsTableMockCalls []string, accessTokensTableMockCalls []string) {
	if !reflect.DeepEqual(gcpEnvMock.GetCalls(), envMockCalls) {
		t.Fatalf("gcpEnvMock expectation violation.\nMock calls: %v", gcpEnvMock.GetCalls())
	}
	if !reflect.DeepEqual(gcpFsMock.GetCalls(), fsMockCalls) {
		t.Fatalf("gcpFsMock expectation violation.\nMock calls: %v", gcpFsMock.GetCalls())
	}
	if !reflect.DeepEqual(credentialsTableMock.GetCalls(), credentialsTableMockCalls) {
		t.Fatalf("credentialsTableMockCalls expectation violation.\nMock calls: %v", credentialsTableMock.GetCalls())
	}
	if !reflect.DeepEqual(accessTokensTableMock.GetCalls(), accessTokensTableMockCalls) {
		t.Fatalf("accessTokensTableMockCalls expectation violation.\nMock calls: %v", accessTokensTableMock.GetCalls())
	}
}

func TestIsGcloudCliAvailable_available(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsCommandAvailable = true

	if !repo.isGcloudCliAvailable() {
		t.Fatalf("should be available")
	}
}

func TestIsGcloudCliAvailable_unavailable(t *testing.T) {
	gcpConfigurationRepositorySetup()
	if repo.isGcloudCliAvailable() {
		t.Fatalf("should be unavailable")
	}
}

func TestDoesGcloudConfigFolderExist_exists(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpDoesFileExist = true
	gcpEnvMock.ExpIsWindows = true

	exists, err := repo.DoesGcloudConfigFolderExist()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !exists {
		t.Fatalf("directory should exist")
	}

	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{fmt.Sprintf("DoesFileExist(%v)", expectedConfigDirPath)}, []string{},
		[]string{})
}

func TestDoesGcloudConfigFolderExist_not_exists(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true

	exists, err := repo.DoesGcloudConfigFolderExist()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exists {
		t.Fatalf("directory should not exist")
	}

	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{fmt.Sprintf("DoesFileExist(%v)", expectedConfigDirPath)}, []string{},
		[]string{})
}

func TestDoesGcloudConfigFolderExist_errorGettingHomeDir(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnGetHomeDir = true

	_, err := repo.DoesGcloudConfigFolderExist()
	test.ExpectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{}, []string{})
}

func TestGcloudConfigurationRepository_onWindows(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true

	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedConfigDir := filepath.Join("c:/", "appdata", "gcloud")
	if configDir != expectedConfigDir {
		t.Fatalf("expected gcloudConfigDir: %v", expectedConfigDir)
	}

	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{}, []string{}, []string{})
}

func TestGcloudConfigurationRepository_onUnix(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpHomeDir = "/user/home"
	configDir, err := repo.gcloudConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedConfigDir := filepath.Join("/", "user", "home", ".config", "gcloud")
	if configDir != expectedConfigDir {
		t.Fatalf("expected gcloudConfigDir: %v", expectedConfigDir)
	}

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{},
		[]string{})
}

func TestGcloudConfigurationRepository_error(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnGetHomeDir = true

	_, err := repo.gcloudConfigDir()
	test.ExpectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{}, []string{})
}

func TestCreateConfiguration(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true

	err := repo.CreateConfiguration("accountId", "projectId")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedFileContent := "[core]\naccount = accountId\nproject = projectId\n"
	expectedFile := fmt.Sprintf("WriteToFile(%v, %v)", expectedConfigFilePath, []byte(expectedFileContent))

	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedFile}, []string{}, []string{})
}

func TestCreateConfiguration_errorGettingHomeDir(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnGetHomeDir = true

	err := repo.CreateConfiguration("accountId", "projectId")
	test.ExpectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{}, []string{})
}

func TestCreateConfiguration_errorWritingFile(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnWriteToFile = true
	gcpEnvMock.ExpIsWindows = true

	err := repo.CreateConfiguration("accountId", "projectId")
	test.ExpectHttpError(t, err, 500, "error writing file")

	expectedFileContent := "[core]\naccount = accountId\nproject = projectId\n"
	expectedWriteFile := fmt.Sprintf("WriteToFile(%v, %v)", expectedConfigFilePath, []byte(expectedFileContent))
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedWriteFile}, []string{}, []string{})
}

func TestRemoveConfiguration(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true

	err := repo.RemoveConfiguration()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedFile := fmt.Sprintf("RemoveFile(%v)", expectedConfigFilePath)
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedFile}, []string{}, []string{})
}

func TestRemoveConfiguration_errorGettingHomeDir(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnGetHomeDir = true

	err := repo.RemoveConfiguration()
	test.ExpectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{}, []string{})
}

func TestRemoveConfiguration_errorRemovingFile(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true
	gcpFsMock.ExpErrorOnRemoveFile = true

	err := repo.RemoveConfiguration()
	test.ExpectHttpError(t, err, 500, "error removing file")

	expectedFile := fmt.Sprintf("RemoveFile(%v)", expectedConfigFilePath)
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedFile}, []string{}, []string{})
}

func TestActivateConfiguration(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true

	err := repo.ActivateConfiguration()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	fullActiveConfigurationName := fmt.Sprintf("leapp")
	expectedFile := fmt.Sprintf("WriteToFile(%v, %v)", expectedActiveConfigFilePath, []byte(fullActiveConfigurationName))
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedFile}, []string{}, []string{})
}

func TestActivateConfiguration_errorGettingHomeDir(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnGetHomeDir = true

	err := repo.ActivateConfiguration()
	test.ExpectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{}, []string{})
}

func TestActivateConfiguration_errorWritingFile(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true
	gcpFsMock.ExpErrorOnWriteToFile = true

	err := repo.ActivateConfiguration()
	test.ExpectHttpError(t, err, 500, "error writing file")

	expectedFile := fmt.Sprintf("WriteToFile(%v, %v)", expectedActiveConfigFilePath,
		[]byte("leapp"))
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedFile}, []string{}, []string{})
}

func TestDeactivateConfiguration(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true

	err := repo.DeactivateConfiguration()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedFile := fmt.Sprintf("RemoveFile(%v)", expectedActiveConfigFilePath)
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedFile}, []string{}, []string{})
}

func TestDeactivateConfiguration_errorGettingHomeDir(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnGetHomeDir = true

	err := repo.DeactivateConfiguration()
	test.ExpectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{}, []string{})
}

func TestDeactivateConfiguration_errorRemovingFile(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true
	gcpFsMock.ExpErrorOnRemoveFile = true

	err := repo.DeactivateConfiguration()
	test.ExpectHttpError(t, err, 500, "error removing file")

	expectedFile := fmt.Sprintf("RemoveFile(%v)", expectedActiveConfigFilePath)
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedFile}, []string{}, []string{})
}

func TestWriteDefaultCredentials(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true

	defaultCredentialJson := "{\"credentials\":\"json\"}"
	err := repo.WriteDefaultCredentials(defaultCredentialJson)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedFile := fmt.Sprintf("WriteToFile(%v, %v)", expectedDefaultCredentialsFilePath, []byte(defaultCredentialJson))
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedFile}, []string{}, []string{})
}

func TestWriteDefaultCredentials_errorGettingHomeDir(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnGetHomeDir = true

	err := repo.WriteDefaultCredentials("{\"credentials\":\"json\"}")
	test.ExpectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{}, []string{})
}

func TestWriteDefaultCredentials_errorWritingFile(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnWriteToFile = true
	gcpEnvMock.ExpIsWindows = true

	err := repo.WriteDefaultCredentials("{\"credentials\":\"json\"}")
	test.ExpectHttpError(t, err, 500, "error writing file")

	defaultCredentialJson := "{\"credentials\":\"json\"}"
	expectedFile := fmt.Sprintf("WriteToFile(%v, %v)", expectedDefaultCredentialsFilePath,
		[]byte(defaultCredentialJson))
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedFile}, []string{}, []string{})
}

func TestRemoveDefaultCredentials(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true

	err := repo.RemoveDefaultCredentials()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	expectedFile := fmt.Sprintf("RemoveFile(%v)", expectedDefaultCredentialsFilePath)
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedFile}, []string{}, []string{})
}

func TestRemoveDefaultCredentials_errorGettingHomeDir(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnGetHomeDir = true

	err := repo.RemoveDefaultCredentials()
	test.ExpectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{}, []string{})
}

func TestRemoveDefaultCredentials_errorRemovingFile(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true
	gcpFsMock.ExpErrorOnRemoveFile = true

	err := repo.RemoveDefaultCredentials()
	test.ExpectHttpError(t, err, 500, "error removing file")

	expectedRemoveFileCall := fmt.Sprintf("RemoveFile(%v)", expectedDefaultCredentialsFilePath)
	verifyExpectedCalls(t, []string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{expectedRemoveFileCall}, []string{}, []string{})
}

func TestWriteCredentialsToDb(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true

	accountId := "account_id@domain.com"
	defaultCredentialJson := "{\"credentials\":\"json\"}"
	err := repo.WriteCredentialsToDb(accountId, defaultCredentialJson)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	verifyExpectedCalls(t,
		[]string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"}, []string{},
		[]string{fmt.Sprintf("WriteCredentials(%v, %v, %v)", expectedCredentialsDbFilePath, accountId, defaultCredentialJson)},
		[]string{})
}

func TestWriteCredentialsToDb_errorGettingHomeDir(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnGetHomeDir = true

	accountId := "account_id@domain.com"
	defaultCredentialJson := "{\"credentials\":\"json\"}"
	err := repo.WriteCredentialsToDb(accountId, defaultCredentialJson)
	test.ExpectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{}, []string{})
}

func TestWriteCredentialsToDb_errorWritingCredentials(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true
	credentialsTableMock.ExpErrorOnExecInsertQuery = true

	accountId := "account_id@domain.com"
	defaultCredentialJson := "{\"credentials\":\"json\"}"
	err := repo.WriteCredentialsToDb(accountId, defaultCredentialJson)
	test.ExpectHttpError(t, err, 500, "error executing insert query")
	expectedWriteCall := fmt.Sprintf("WriteCredentials(%v, %v, %v)", expectedCredentialsDbFilePath, accountId, defaultCredentialJson)

	verifyExpectedCalls(t,
		[]string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{},
		[]string{expectedWriteCall},
		[]string{})
}

func TestRemoveCredentialsFromDb(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true
	gcpFsMock.ExpDoesFileExist = true

	accountId := "account_id@domain.com"
	err := repo.RemoveCredentialsFromDb(accountId)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	verifyExpectedCalls(t,
		[]string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{fmt.Sprintf("DoesFileExist(%v)", expectedCredentialsDbFilePath)},
		[]string{fmt.Sprintf("RemoveCredentials(%v, %v)", expectedCredentialsDbFilePath, accountId)},
		[]string{})
}

func TestRemoveCredentialsFromDb_DbFileDoesNotExist(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true

	accountId := "account_id@domain.com"
	err := repo.RemoveCredentialsFromDb(accountId)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	verifyExpectedCalls(t,
		[]string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{fmt.Sprintf("DoesFileExist(%v)", expectedCredentialsDbFilePath)},
		[]string{}, []string{})
}

func TestRemoveCredentialsFromDb_errorGettingHomeDir(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnGetHomeDir = true

	accountId := "account_id@domain.com"
	err := repo.RemoveCredentialsFromDb(accountId)
	test.ExpectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{}, []string{})
}

func TestRemoveCredentialsFromDb_errorRemovingCredentialsFromDb(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true
	gcpFsMock.ExpDoesFileExist = true
	credentialsTableMock.ExpErrorOnExecDeleteQuery = true

	accountId := "account_id@domain.com"
	err := repo.RemoveCredentialsFromDb(accountId)
	test.ExpectHttpError(t, err, 500, "error executing delete query")

	expectedRemoveCall := fmt.Sprintf("RemoveCredentials(%v, %v)", expectedCredentialsDbFilePath, accountId)
	verifyExpectedCalls(t,
		[]string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{fmt.Sprintf("DoesFileExist(%v)", expectedCredentialsDbFilePath)},
		[]string{expectedRemoveCall}, []string{})
}

func TestRemoveAccessTokensFromDb(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true
	gcpFsMock.ExpDoesFileExist = true

	accountId := "account_id@domain.com"
	err := repo.RemoveAccessTokensFromDb(accountId)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	verifyExpectedCalls(t,
		[]string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{fmt.Sprintf("DoesFileExist(%v)", expectedAccessTokensDbFilePath)},
		[]string{},
		[]string{fmt.Sprintf("RemoveAccessToken(%v, %v)", expectedAccessTokensDbFilePath, accountId)})
}

func TestRemoveAccessTokensFromDb_DbFileDoesNotExist(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true

	accountId := "account_id@domain.com"
	err := repo.RemoveAccessTokensFromDb(accountId)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	verifyExpectedCalls(t,
		[]string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{fmt.Sprintf("DoesFileExist(%v)", expectedAccessTokensDbFilePath)},
		[]string{}, []string{})
}

func TestRemoveAccessTokensFromDb_errorGettingHomeDir(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpFsMock.ExpErrorOnGetHomeDir = true

	accountId := "account_id@domain.com"
	err := repo.RemoveAccessTokensFromDb(accountId)
	test.ExpectHttpError(t, err, 500, "error getting home dir")

	verifyExpectedCalls(t, []string{"IsWindows()"}, []string{"GetHomeDir()"}, []string{}, []string{})
}

func TestRemoveAccessTokensFromDb_errorRemovingCredentialsFromDb(t *testing.T) {
	gcpConfigurationRepositorySetup()
	gcpEnvMock.ExpIsWindows = true
	gcpFsMock.ExpDoesFileExist = true
	accessTokensTableMock.ExpErrorOnExecDeleteQuery = true

	accountId := "account_id@domain.com"
	err := repo.RemoveAccessTokensFromDb(accountId)
	test.ExpectHttpError(t, err, 500, "error executing delete query")
	expectedRemoveCall := fmt.Sprintf("RemoveAccessToken(%v, %v)", expectedAccessTokensDbFilePath, accountId)

	verifyExpectedCalls(t,
		[]string{"IsWindows()", "GetEnvironmentVariable(APPDATA)"},
		[]string{fmt.Sprintf("DoesFileExist(%v)", expectedAccessTokensDbFilePath)},
		[]string{}, []string{expectedRemoveCall})
}
