package environment

import (
	"github.com/google/uuid"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Environment struct {
}

func (env *Environment) GetEnvironmentVariable(variable string) string {
	return os.Getenv(variable)
}

func (env *Environment) IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func (env *Environment) GetOs() string {
	return runtime.GOOS
}

func (env *Environment) IsWindows() bool {
	return env.GetOs() == "windows"
}

func (env *Environment) GenerateUuid() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func (env *Environment) GetTime() string {
	return time.Now().Format(time.RFC3339)
}

func (env *Environment) GetFormattedTime(format string) string {
	return time.Now().Format(format)
}
