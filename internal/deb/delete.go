/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package deb

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Delete(pkgName, distro, section, version string) (int, error) {
	sectionPath, err := getDebianSectionPath(pkgName, distro, section)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	cfg, err := NewConfig()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot read configuration file: %s\n", err)
	}
	// check the repository repoName has been configured
	repo := cfg.GetRepo(pkgName)
	if repo == nil {
		return http.StatusBadRequest, fmt.Errorf("invalid repository: %s\n", pkgName)
	}
	var (
		archs    []os.DirEntry
		packages *PackagesData
		pkgPath  string
	)
	archs, err = os.ReadDir(sectionPath)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// loop through  architectures in section and find packages to delete using a version regex
	for _, arch := range archs {
		if !arch.IsDir() {
			continue
		}
		// load Packages info
		packages, err = newPackagesData(filepath.Join(sectionPath, arch.Name(), "Packages"))
		arc := arch.Name()[len("binary-"):]
		pkgPath, err = getDebianPkgPath(pkgName, distro, section, arc)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		// find the right files
		var (
			files   []os.DirEntry
			meta    *PackageMeta
			content []byte
			matched bool
		)
		files, err = os.ReadDir(pkgPath)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		// delete packages
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".deb") {
				continue
			}
			content, err = os.ReadFile(filepath.Join(pkgPath, file.Name()))
			if err != nil {
				return http.StatusInternalServerError, err
			}
			meta, err = getPackageMeta(content)
			if matched, err = regexp.MatchString(version, meta.Version); matched {
				err = os.Remove(filepath.Join(pkgPath, file.Name()))
				if err != nil {
					return http.StatusInternalServerError, err
				}
			}
			if !packages.Remove(pkgName, meta.Version) {
				return http.StatusInternalServerError, fmt.Errorf("cannot remove package meta for:\n%s", meta.Raw)
			}
		}
		//  regenerates Packages metadata
		if err = packages.Save(filepath.Join(pkgPath, "Packages")); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("cannot save Packages file: '%s'\n", err)
		}
		if err = packages.SaveGz(filepath.Join(pkgPath, "Packages.gz")); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("cannot save Packages.gz file: '%s'\n", err)
		}
	}
	// update Release files
	if err = CreateRelease(*repo); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot update Release files: '%s'\n", err)
	}
	return 0, nil
}
