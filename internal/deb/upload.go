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
	"strings"
)

func Upload(repoName, section string, pkgBytes []byte) (errorCode int, err error) {
	// read the package metadata
	meta, err := getPackageMeta(pkgBytes)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("cannot retrieve package metadata: %s\n", err)
	}
	cfg, err := NewConfig()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot read configuration file: %s\n", err)
	}
	// check the repository repoName has been configured
	repo := cfg.GetRepo(repoName)
	if repo == nil {
		return http.StatusBadRequest, fmt.Errorf("invalid repository: %s\n", repoName)
	}
	// check that the package repoName matches the repository repoName
	if !strings.EqualFold(meta.Package, repo.Name) {
		return http.StatusBadRequest, fmt.Errorf("invalid package '%s' for this repository '%s'\n", meta.Package, repo.Name)
	}
	// works out the path where the package should be saved
	pkgPath, err := checkPkgPath(repo.Name, repo.Distribution, section, meta.Architecture)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot configure package path: '%s'\n", err)
	}
	// now loads the packages metadata to check if the package is already in the repository
	packages, err := newPackagesData(filepath.Join(pkgPath, "Packages"))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot load packages metadata: '%s'\n", err)
	}
	// check if package already exists
	if packages.Find(repo.Name, meta.Version, meta.Release, meta.Architecture) != nil {
		return http.StatusBadRequest, fmt.Errorf("package already exists\n")
	}
	// save the debian package to disk
	err = os.WriteFile(filepath.Join(pkgPath, pkgName(repo.Name, meta.Version, meta.Release, meta.Architecture)), pkgBytes, 0755)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot persist package: '%s'\n", err)
	}
	// gets package checksums
	var md5Sum, sha1Sum, sha256Sum string
	md5Sum, sha1Sum, sha256Sum, err = getChecksums(pkgBytes)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot calculate package checksums: '%s'\n", err)
	}
	// adds package to metadata
	packages.Add(PackageData{
		Package:      meta.Package,
		Version:      meta.Version,
		Release:      meta.Release,
		Architecture: meta.Architecture,
		Maintainer:   meta.Maintainer,
		Homepage:     meta.Homepage,
		Depends:      meta.Depends,
		Description:  meta.Description,
		Filename: filepath.Join(
			"dists",
			repo.Distribution, section,
			fmt.Sprintf("binary-%s", meta.Architecture),
			pkgName(repo.Name, meta.Version, meta.Release, meta.Architecture)),
		Size:   fmt.Sprintf("%d", len(pkgBytes)),
		MD5sum: md5Sum,
		SHA1:   sha1Sum,
		SHA256: sha256Sum,
	})
	// saves metadata
	if err = packages.Save(filepath.Join(pkgPath, "Packages")); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot save Packages file: '%s'\n", err)
	}
	if err = packages.SaveGz(filepath.Join(pkgPath, "Packages.gz")); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot save Packages.gz file: '%s'\n", err)
	}
	if err = CreateRelease(*repo); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot update Release files: '%s'\n", err)
	}
	return 0, nil
}
