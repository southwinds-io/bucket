/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package deb

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CreateRelease(repo Repository) error {
	workingDirectory, err := getReleasePath(repo)
	if err != nil {
		return err
	}
	outFile, err := os.Create(filepath.Join(workingDirectory, "Release"))
	if err != nil {
		return fmt.Errorf("failed to create Release: %s", err)
	}
	defer func(outFile *os.File) {
		err = outFile.Close()
		if err != nil {
			return
		}
	}(outFile)

	currentTime := time.Now().UTC()
	_, err = fmt.Fprintf(outFile, "Suite: %s\n", repo.Distribution)
	if err != nil {
		return fmt.Errorf("can not write file: %v", err)
	}
	_, err = fmt.Fprintf(outFile, "Codename: %s\n", repo.Distribution)
	if err != nil {
		return fmt.Errorf("can not write file: %v", err)
	}
	_, err = fmt.Fprintf(outFile, "Components: %s\n", strings.Join(repo.Sections, " "))
	if err != nil {
		return fmt.Errorf("can not write file: %v", err)
	}
	_, err = fmt.Fprintf(outFile, "Architectures: %s\n", strings.Join(repo.Architectures, " "))
	if err != nil {
		return fmt.Errorf("can not write file: %v", err)
	}
	_, err = fmt.Fprintf(outFile, "Date: %s\n", currentTime.Format("Mon, 02 Jan 2006 15:04:05 UTC"))
	if err != nil {
		return fmt.Errorf("can not write file: %v", err)
	}

	var md5Sums strings.Builder
	var sha1Sums strings.Builder
	var sha256Sums strings.Builder

	err = filepath.Walk(workingDirectory, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "Package.gz") || strings.HasSuffix(path, "Packages") {
			var (
				md5hash    = md5.New()
				sha1hash   = sha1.New()
				sha256hash = sha256.New()
			)
			relPath, _ := filepath.Rel(workingDirectory, path)
			slashPath := filepath.ToSlash(relPath)
			f, err := os.Open(path)
			if err != nil {
				log.Println("error opening the packages file for reading", err)
			}
			if _, err = io.Copy(io.MultiWriter(md5hash, sha1hash, sha256hash), f); err != nil {
				return fmt.Errorf("error hashing file for release list: %s", err)
			}
			_, err = fmt.Fprintf(&md5Sums, " %s %d %s\n", hex.EncodeToString(md5hash.Sum(nil)), file.Size(), slashPath)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(&sha1Sums, " %s %d %s\n", hex.EncodeToString(sha1hash.Sum(nil)), file.Size(), slashPath)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(&sha256Sums, " %s %d %s\n", hex.EncodeToString(sha256hash.Sum(nil)), file.Size(), slashPath)
			if err != nil {
				return err
			}

			f = nil
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error scaning for packages files: %s", err)
	}

	_, err = outFile.WriteString("MD5Sum:\n")
	if err != nil {
		return fmt.Errorf("can not write file: %v", err)
	}
	_, err = outFile.WriteString(md5Sums.String())
	if err != nil {
		return fmt.Errorf("can not write file: %v", err)
	}
	_, err = outFile.WriteString("SHA1:\n")
	if err != nil {
		return fmt.Errorf("can not write file: %v", err)
	}
	_, err = outFile.WriteString(sha1Sums.String())
	if err != nil {
		return fmt.Errorf("can not write file: %v", err)
	}
	_, err = outFile.WriteString("SHA256:\n")
	if err != nil {
		return fmt.Errorf("can not write file: %v", err)
	}
	_, err = outFile.WriteString(sha256Sums.String())
	if err != nil {
		return fmt.Errorf("can not write file: %v", err)
	}

	key, ok := repo.cfg.GetKey(repo.KeyRef)
	if !ok {
		return fmt.Errorf("cannot find signing key for repository %s, check the service configuration", repo.Name)
	}
	if err = signRelease([]byte(key.Private), []byte(key.Passcode), outFile.Name()); err != nil {
		return fmt.Errorf("error signing Release file: %s", err)
	}
	return nil
}
