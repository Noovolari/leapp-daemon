package environment

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestGetEnvironmentVariable(t *testing.T) {
	variableName := "testVariable"
	variableValue := "check123"
	defer func(key string) {
		err := os.Unsetenv(key)
		if err != nil {
			t.Fatalf("Unable to unset the environment variable used for the test")
		}
	}(variableName)

	env := Environment{}
	if env.GetEnvironmentVariable(variableName) != "" {
		t.Fatalf("Unexpected variableValue")
	}

	err := os.Setenv(variableName, variableValue)
	if err != nil {
		t.Fatalf("Unable to set environment variable for the test")
	}

	if env.GetEnvironmentVariable(variableName) != variableValue {
		t.Fatalf("Unexpected variableValue")
	}
}

func TestIsCommandAvailable_available(t *testing.T) {
	env := Environment{}

	if runtime.GOOS == "windows" {
		if !env.IsCommandAvailable("cmd.exe") {
			t.Fatalf("cmd.exe should exist on Windows")
		}
	} else {
		if !env.IsCommandAvailable("ls") {
			t.Fatalf("ls should exist on Unix-Like system")
		}
	}
}

func TestIsCommandAvailable_unavailable(t *testing.T) {
	env := Environment{}

	if env.IsCommandAvailable("unavailable_command_23465236") {
		t.Fatalf("command should not exist")
	}
}

func TestGetOs(t *testing.T) {
	env := Environment{}

	if env.GetOs() == "" {
		t.Fatalf("Unable to define the current OS")
	}

	if env.GetOs() != runtime.GOOS {
		t.Fatalf("Unexpected OS returned")
	}
}

func TestIsWindows(t *testing.T) {
	env := Environment{}

	if runtime.GOOS == "windows" && !env.IsWindows() {
		t.Fatalf("Environment is Windows")
	} else if runtime.GOOS != "windows" && env.IsWindows() {
		t.Fatalf("Environment is not Windows")
	}
}

func TestGenerateUuid_randomness(t *testing.T) {
	env := Environment{}

	previousUuid := env.GenerateUuid()
	for i := 0; i < 100; i++ {
		currentUuid := env.GenerateUuid()
		if previousUuid == currentUuid {
			t.Fatalf("expected different uuids")
		}

		previousUuid = currentUuid
	}
}

func TestGenerateUuid_withoutDashes(t *testing.T) {
	env := Environment{}

	uuid := env.GenerateUuid()
	if strings.Contains(uuid, "-") {
		t.Fatalf("uuid should not contain dashes")
	}
}

func TestGetTime(t *testing.T) {
	env := Environment{}

	timestamp := env.GetTime()
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Fatalf("cannot parse timestamp")
	}

	if time.Now().Before(parsedTime) {
		t.Fatalf("getTime is returning a time after now")
	}
}
