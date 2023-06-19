package internal

import (
	"os"
	"os/user"
	"path"
)

// GetWorkspaceDir returns the path to a workspace directory for mear to manipulate files with.
func GetWorkspaceDir(workspace_name string) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	workspace_dir := path.Join(currentUser.HomeDir, ".mear", workspace_name)
	if err := os.MkdirAll(workspace_dir, 0755); err != nil {
		return "", err
	}

	return workspace_dir, nil
}
