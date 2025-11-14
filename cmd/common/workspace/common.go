package workspace

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
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
var errorWorkspaceNotFound = errors.New("workspace not found")

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

func checkWorkspaceName(dirConfigWorkspace string, name string) error {
	if name == currentWorkspaceNameFile {
		return errorNameNotAllowed
	}
	if !isWorkspaceNameValid(name) {
		return errorInvalidName
	}

	if _, err := os.Stat(path.Join(dirConfigWorkspace, name)); os.IsNotExist(err) {
		return errorWorkspaceNotFound
	} else if err != nil {
		return err
	}

	return nil // Workspace found
}

func read(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}
