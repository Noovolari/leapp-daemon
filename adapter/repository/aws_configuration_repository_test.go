package repository

import (
	"fmt"
	"leapp_daemon/test"
	"leapp_daemon/test/mock"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	tempDirPath                string
	tempAwsDirPath             string
	tempCredentialsFilePath    string
	awsCredentials             AwsTempCredentials
	anotherAwsCredentials      AwsTempCredentials
	awsCredentialsFileContent1 string
	awsCredentialsFileContent2 string
	awsFsMock                  mock.FileSystemMock
	awsEnvMock                 mock.EnvironmentMock
	awsRepo                    *AwsConfigurationRepository
)

func awsConfigurationRepositorySetup() {
	tempDirPath = filepath.Join(os.TempDir(), "temp-aws")
	_ = os.RemoveAll(tempDirPath)
	tempAwsDirPath = filepath.Join(tempDirPath, ".aws")
	_ = os.MkdirAll(tempAwsDirPath, 0600)
	tempCredentialsFilePath = filepath.Join(tempAwsDirPath, "credentials")

	awsCredentials = AwsTempCredentials{
		ProfileName:  "profile-name",
		AccessKeyId:  "access-key-id",
		SecretKey:    "secret-key",
		SessionToken: "session-token",
		Region:       "region",
	}

	anotherAwsCredentials = AwsTempCredentials{
		ProfileName:  "another-profile-name",
		AccessKeyId:  "another-access-key-id",
		SecretKey:    "another-secret-key",
		SessionToken: "another-session-token",
		Region:       "another-region",
	}

	awsCredentialsFileContent1 = "; Credentials managed by Leapp, do not manually edit this file\n" +
		"[profile-name]\n" +
		"aws_access_key_id     = access-key-id\n" +
		"aws_secret_access_key = secret-key\n" +
		"aws_session_token     = session-token\n" +
		"region                = region\n" +
		"\n"

	awsCredentialsFileContent2 = "; Credentials managed by Leapp, do not manually edit this file\n" +
		"[another-profile-name]\n" +
		"aws_access_key_id     = another-access-key-id\n" +
		"aws_secret_access_key = another-secret-key\n" +
		"aws_session_token     = another-session-token\n" +
		"region                = another-region\n" +
		"\n"

	awsFsMock = mock.NewFileSystemMock()
	awsEnvMock = mock.NewEnvironmentMock()
	awsRepo = &AwsConfigurationRepository{
		FileSystem:  &awsFsMock,
		Environment: &awsEnvMock,
	}
}

func awsConfigurationRepositoryVerifyExpectedCalls(t *testing.T, fsMockCalls []string, envMockCalls []string) {
	if !reflect.DeepEqual(awsFsMock.GetCalls(), fsMockCalls) {
		t.Fatalf("awsFsMock expectation violation.\nMock calls: %v", awsFsMock.GetCalls())
	}
	if !reflect.DeepEqual(awsEnvMock.GetCalls(), envMockCalls) {
		t.Fatalf("awsEnvMock expectation violation.\nMock calls: %v", awsEnvMock.GetCalls())
	}
}

func TestWriteCredentials_ErrorOnGetCredentialsFilePath(t *testing.T) {
	awsConfigurationRepositorySetup()
	awsFsMock.ExpErrorOnGetHomeDir = true

	err := awsRepo.WriteCredentials([]AwsTempCredentials{awsCredentials})
	test.ExpectHttpError(t, err, 500, "error getting home dir")
	awsConfigurationRepositoryVerifyExpectedCalls(t, []string{"GetHomeDir()"}, []string{})
}

func TestWriteCredentials_CredentialsFileDoesNotExist(t *testing.T) {
	awsConfigurationRepositorySetup()
	awsFsMock.ExpHomeDir = tempDirPath

	err := awsRepo.WriteCredentials([]AwsTempCredentials{awsCredentials})
	if err != nil {
		t.Fatalf("unexpected error")
	}

	fileData, err := os.ReadFile(tempCredentialsFilePath)
	if err != nil {
		t.Fatalf("credentials file has not been created")
	}

	actualFileContent := string(fileData)
	if actualFileContent != awsCredentialsFileContent1 {
		t.Fatalf("unexpected credentials file content:\n%v", actualFileContent)
	}

	expectedFileCheck := fmt.Sprintf("DoesFileExist(%v)", tempCredentialsFilePath)
	awsConfigurationRepositoryVerifyExpectedCalls(t, []string{"GetHomeDir()", expectedFileCheck}, []string{})
}

func TestWriteCredentials_CredentialsFileDoesNotExist_NoRegion(t *testing.T) {
	awsConfigurationRepositorySetup()
	awsFsMock.ExpHomeDir = tempDirPath

	awsCredentials.Region = ""
	err := awsRepo.WriteCredentials([]AwsTempCredentials{awsCredentials})
	if err != nil {
		t.Fatalf("unexpected error")
	}

	fileData, err := os.ReadFile(tempCredentialsFilePath)
	if err != nil {
		t.Fatalf("credentials file has not been created")
	}

	expectedFileContent := "; Credentials managed by Leapp, do not manually edit this file\n" +
		"[profile-name]\n" +
		"aws_access_key_id     = access-key-id\n" +
		"aws_secret_access_key = secret-key\n" +
		"aws_session_token     = session-token\n" +
		"\n"

	actualFileContent := string(fileData)
	if actualFileContent != expectedFileContent {
		t.Fatalf("unexpected credentials file content:\n%v", actualFileContent)
	}

	expectedFileCheck := fmt.Sprintf("DoesFileExist(%v)", tempCredentialsFilePath)
	awsConfigurationRepositoryVerifyExpectedCalls(t, []string{"GetHomeDir()", expectedFileCheck}, []string{})
}

func TestWriteCredentials_CredentialsFileDoesNotExist_MultipleCredentials(t *testing.T) {
	awsConfigurationRepositorySetup()
	awsFsMock.ExpHomeDir = tempDirPath

	err := awsRepo.WriteCredentials([]AwsTempCredentials{awsCredentials, anotherAwsCredentials})
	if err != nil {
		t.Fatalf("unexpected error")
	}

	fileData, err := os.ReadFile(tempCredentialsFilePath)
	if err != nil {
		t.Fatalf("credentials file has not been created")
	}

	expectedFileContent := awsCredentialsFileContent1 + awsCredentialsFileContent2

	actualFileContent := string(fileData)
	if actualFileContent != expectedFileContent {
		t.Fatalf("unexpected credentials file content:\n%v", actualFileContent)
	}

	expectedFileCheck := fmt.Sprintf("DoesFileExist(%v)", tempCredentialsFilePath)
	awsConfigurationRepositoryVerifyExpectedCalls(t, []string{"GetHomeDir()", expectedFileCheck}, []string{})
}

func TestWriteCredentials_CredentialsFileAlreadyExistAndManagedByLeapp(t *testing.T) {
	awsConfigurationRepositorySetup()
	awsFsMock.ExpHomeDir = tempDirPath
	awsFsMock.ExpDoesFileExist = true

	os.WriteFile(tempCredentialsFilePath, []byte(awsCredentialsFileContent1), 0666)
	err := awsRepo.WriteCredentials([]AwsTempCredentials{anotherAwsCredentials})
	if err != nil {
		t.Fatalf("unexpected error")
	}

	fileData, err := os.ReadFile(tempCredentialsFilePath)
	if err != nil {
		t.Fatalf("credentials file has not been created")
	}

	actualFileContent := string(fileData)
	if actualFileContent != awsCredentialsFileContent2 {
		t.Fatalf("unexpected credentials file content:\n%v", actualFileContent)
	}

	expectedFileCheck := fmt.Sprintf("DoesFileExist(%v)", tempCredentialsFilePath)
	awsConfigurationRepositoryVerifyExpectedCalls(t, []string{"GetHomeDir()", expectedFileCheck}, []string{})
}

func TestWriteCredentials_CredentialsFileAlreadyExistButNotManagedByLeapp(t *testing.T) {
	awsConfigurationRepositorySetup()
	awsFsMock.ExpHomeDir = tempDirPath
	awsFsMock.ExpDoesFileExist = true
	awsEnvMock.ExpFormattedTime = "2021-05-02 15-00-01"

	previousData := "[profile-name]\n" +
		"aws_access_key_id     = access-key-id\n" +
		"aws_secret_access_key = secret-key\n" +
		"aws_session_token     = session-token\n" +
		"region                = region\n" +
		"\n"

	os.WriteFile(tempCredentialsFilePath, []byte(previousData), 0666)
	err := awsRepo.WriteCredentials([]AwsTempCredentials{anotherAwsCredentials})
	if err != nil {
		t.Fatalf("unexpected error")
	}

	fileData, err := os.ReadFile(tempCredentialsFilePath)
	if err != nil {
		t.Fatalf("credentials file has not been created")
	}

	actualFileContent := string(fileData)
	if actualFileContent != awsCredentialsFileContent2 {
		t.Fatalf("unexpected credentials file content:\n%v", actualFileContent)
	}

	expectedCheckFileExist := fmt.Sprintf("DoesFileExist(%v)", tempCredentialsFilePath)
	expectedRenameFile := fmt.Sprintf("RenameFile(%v, %v_2021-05-02 15-00-01.backup)", tempCredentialsFilePath, tempCredentialsFilePath)
	awsConfigurationRepositoryVerifyExpectedCalls(t, []string{"GetHomeDir()", expectedCheckFileExist, "GetHomeDir()",
		expectedRenameFile}, []string{"GetFormattedTime(2006-01-02 15-04-05)"})
}

func TestWriteCredentials_CredentialsFileAlreadyExistButNotManagedByLeapp_BackupFailed(t *testing.T) {
	awsConfigurationRepositorySetup()
	awsFsMock.ExpHomeDir = tempDirPath
	awsFsMock.ExpDoesFileExist = true
	awsFsMock.ExpErrorOnRenamingFile = true
	awsEnvMock.ExpFormattedTime = "2021-05-02 15-00-01"

	previousData := "[profile-name]\n" +
		"aws_access_key_id     = access-key-id\n" +
		"aws_secret_access_key = secret-key\n" +
		"aws_session_token     = session-token\n" +
		"region                = region\n" +
		"\n"

	os.WriteFile(tempCredentialsFilePath, []byte(previousData), 0666)
	err := awsRepo.WriteCredentials([]AwsTempCredentials{anotherAwsCredentials})
	test.ExpectHttpError(t, err, 500, "error renaming file")

	fileData, err := os.ReadFile(tempCredentialsFilePath)
	if err != nil {
		t.Fatalf("credentials file has not been created")
	}

	actualFileContent := string(fileData)
	if actualFileContent != previousData {
		t.Fatalf("unexpected credentials file content:\n%v", actualFileContent)
	}

	expectedCheckFileExist := fmt.Sprintf("DoesFileExist(%v)", tempCredentialsFilePath)
	expectedRenameFile := fmt.Sprintf("RenameFile(%v, %v_2021-05-02 15-00-01.backup)", tempCredentialsFilePath, tempCredentialsFilePath)
	awsConfigurationRepositoryVerifyExpectedCalls(t, []string{"GetHomeDir()", expectedCheckFileExist, "GetHomeDir()",
		expectedRenameFile}, []string{"GetFormattedTime(2006-01-02 15-04-05)"})
}
