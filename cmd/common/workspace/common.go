package workspace

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
)

var isWorkspaceNameValid = regexp.MustCompile(`^[\w-]*.$`).MatchString
var errorNameNotAllowed = fmt.Errorf("%s is not an allowed name", currentWorkspaceNameFile)
var errorInvalidName = errors.New("name should only contain alphanumric characters, underscores or hypens")
var errorWorkspaceAlreadyExists = errors.New("workspace already exists")
var errorDeleteCurrentNotAllowed = errors.New("cannot delete current workspace")
var errorCopyToSelf = errors.New("cannot copy to itself")

const currentWorkspaceNameFile = "current"
const defaultWorkspaceName = "default"
const envWorkspaceVar = "MGC_WORKSPACE"
const FILE_PERMISSION = 0644
const DIR_PERMISSION = 0744

func buildMGCPath() (string, error) {
	dir := ""
	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("AppData")
		if dir == "" {
			return "", errors.New("%AppData% is not defined")
		}

	default: // Unix
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir == "" {
			home := os.Getenv("HOME")
			if home != "" {
				dir = path.Join(home, ".config")
			}
		}
		if dir == "" {
			return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
		}
	}
	mgcDir := path.Join(dir, "mgc")
	if err := os.MkdirAll(mgcDir, DIR_PERMISSION); err != nil {
		return "", fmt.Errorf("Error creating mgc dir at %s: %w", mgcDir, err)
	}
	return mgcDir, nil
}

func sanitizePath(p string) string {
	pathEntries := strings.Split(p, string(os.PathSeparator))

	if len(pathEntries) == 1 {
		return pathEntries[0]
	}

	result := []string{}
	for _, entry := range pathEntries {
		if entry != "" && entry != ".." && entry != "." {
			result = append(result, entry)
		}
	}

	return strings.Join(result, "/")
}

func checkWorkspaceName(name string) error {
	if name == currentWorkspaceNameFile {
		return errorNameNotAllowed
	}
	if !isWorkspaceNameValid(name) {
		return errorInvalidName
	}

	return nil
}

func read(name string) ([]byte, error) {
	return os.ReadFile(name)
}
