package workspace

import (
	"os"
	"path"
)

type Workspace interface {
	Create(name string) error
	Delete(name string) error
	Get() Workspace
	List() ([]Workspace, error)
	Set(name string) error
	Name() string
	Current() Workspace
	Dir() string
}

type workspace struct {
	current   string
	dirConfig string
}

func NewWorkspace() Workspace {
	dirConfig, err := buildMGCPath()
	if err != nil {
		panic(err)
	}

	return &workspace{
		current:   "",
		dirConfig: dirConfig,
	}
}

func (w *workspace) Create(name string) error {
	err := os.Mkdir(path.Join(w.dirConfig, name), 0755)
	if err != nil {
		return err
	}
	return nil
}

func (w *workspace) Delete(name string) error {
	err := os.Remove(path.Join(w.dirConfig, name))
	if err != nil {
		return err
	}
	return nil
}

func (w *workspace) Get() Workspace {
	if w.current == "" {
		return w.Current()
	}
	return w
}

func (w *workspace) List() ([]Workspace, error) {

	files, err := os.ReadDir(path.Join(w.dirConfig, "workspaces"))
	if err != nil {
		return nil, err
	}

	workspaces := []Workspace{}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		workspaces = append(workspaces, &workspace{current: file.Name()})
	}

	return workspaces, nil
}

func (w *workspace) Set(name string) error {
	if err := checkWorkspaceName(name); err != nil {
		return err
	}

	currentFile := path.Join(w.dirConfig, currentWorkspaceNameFile)
	if err := os.WriteFile(currentFile, []byte(name), FILE_PERMISSION); err != nil {
		return err
	}
	w.current = name
	return nil
}

func (w *workspace) Current() Workspace {
	name := defaultWorkspaceName

	data, err := read(path.Join(w.dirConfig, currentWorkspaceNameFile))
	if err == nil && len(data) > 0 {
		name = string(data)
	}

	w.current = name

	return w
}

func (w *workspace) Name() string {
	return w.current
}

func (w *workspace) Dir() string {
	return path.Join(w.dirConfig, w.current)
}
