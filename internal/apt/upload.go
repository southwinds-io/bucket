/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package apt

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"southwinds.dev/bucket/internal/cfg"
	"strings"
)

func Upload(repoName, dist, section string, pkgBytes []byte) (errorCode int, err error) {
	// read the package metadata
	meta, err := getPackageMeta(pkgBytes)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("cannot retrieve package metadata: %s\n", err)
	}
	conf, err := cfg.NewConfig()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot read configuration file: %s\n", err)
	}
	// check the repository repoName has been configured
	repo := conf.GetRepo(repoName)
	if repo == nil {
		return http.StatusBadRequest, fmt.Errorf("invalid repository: %s\n", repoName)
	}
	// check that the package distribution is supported by the repository
	if !in(dist, repo.Distributions) {
		return http.StatusBadRequest, fmt.Errorf("invalid distribution: %s, not allowed in repository %s\n", dist, repoName)
	}
	// check that the package section is supported by the repository
	if !in(section, repo.Sections) {
		return http.StatusBadRequest, fmt.Errorf("invalid section: %s, not allowed in repository %s\n", section, repoName)
	}
	// check that the package architecture is supported by the repository
	if !in(meta.Architecture, repo.Architectures) {
		return http.StatusBadRequest, fmt.Errorf("invalid package architecture: %s, not allowed in repository %s\n", meta.Architecture, repoName)
	}
	// works out the path where the package should be saved
	pkgPath, err := cfg.CheckAptPkgPath(repo.Name, dist, section, meta.Architecture)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot configure package path: '%s'\n", err)
	}
	// now loads the packages metadata to check if the package is already in the repository
	packages, err := NewPackagesData(filepath.Join(pkgPath, "Packages"))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot load packages metadata: '%s'\n", err)
	}
	// check if package already exists
	if packages.Find(repo.Name, meta.Version, meta.Release, meta.Architecture) != nil {
		return http.StatusBadRequest, fmt.Errorf("package already exists\n")
	}
	// save the debian package to disk
	err = os.WriteFile(filepath.Join(pkgPath, cfg.DebianPkgName(repo.Name, meta.Version, meta.Release, meta.Architecture)), pkgBytes, 0755)
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
			dist,
			section,
			fmt.Sprintf("binary-%s", meta.Architecture),
			cfg.DebianPkgName(repo.Name, meta.Version, meta.Release, meta.Architecture)),
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
	if err = CreateRelease(*repo, dist); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("cannot update Release files: '%s'\n", err)
	}
	return 0, nil
}

func in(value string, values []string) bool {
	for _, v := range values {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}
