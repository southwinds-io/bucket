/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package deb

import (
	"fmt"
	"os"
	"path/filepath"
)

const repoFolder = ".bucket"

func GetDataPath() (string, error) {
	path := os.Getenv("BUCKET_DATA_PATH")
	if len(path) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, repoFolder), nil
	}
	return path, nil
}

func getConfigPath() (string, error) {
	path := os.Getenv("BUCKET_CONFIG_PATH")
	if len(path) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, repoFolder), nil
	}
	return path, nil
}

func getReleasePath(repo Repository) (path string, err error) {
	var root string
	root, err = GetDataPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, repo.Name, "dists", repo.Distribution), nil
}

func getSectionPath(repo, distro, section string) (path string, err error) {
	var root string
	root, err = GetDataPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, repo, "dists", distro, section), nil
}

func getPkgPath(repo, distro, section, arch string) (path string, err error) {
	var root string
	root, err = GetDataPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, repo, "dists", distro, section, fmt.Sprintf("binary-%s", arch)), nil
}

func checkPkgPath(repo, distro, section, arch string) (path string, err error) {
	var root string
	root, err = GetDataPath()
	if err != nil {
		return "", err
	}
	pkgPath := filepath.Join(root, repo, "dists", distro, section, fmt.Sprintf("binary-%s", arch))
	if _, err = os.Stat(pkgPath); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(pkgPath, 0755); err != nil {
				return pkgPath, err
			}
		} else {
			return "", fmt.Errorf("error inspecting %s (%s): %s", distro, arch, err)
		}
	}
	return pkgPath, nil
}
