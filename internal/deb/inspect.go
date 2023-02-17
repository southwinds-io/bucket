/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package deb

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/blakesmith/ar"
	"github.com/ulikunitz/xz/lzma"
	"io"
	"strings"
)

type Compression int

const (
	LZMA Compression = iota
	GZIP
)

type PackageMeta struct {
	Package      string
	Version      string
	Release      string
	Architecture string
	Maintainer   string
	Homepage     string
	Description  string
	Depends      string
	Raw          string
}

// getPackageMeta get metadata in the specified debian package file
func getPackageMeta(f []byte) (*PackageMeta, error) {
	meta, err := inspectPackage(f)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(meta, "\n")
	result := &PackageMeta{
		Raw: meta,
	}
	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		switch parts[0] {
		case "Package":
			result.Package = strings.Trim(parts[1], " ")
			break
		case "Version":
			result.Version = strings.Trim(parts[1], " ")
			break
		case "Release":
			result.Release = strings.Trim(parts[1], " ")
			break
		case "Architecture":
			result.Architecture = strings.Trim(parts[1], " ")
			break
		case "Maintainer":
			result.Maintainer = strings.Trim(parts[1], " ")
			break
		case "Homepage":
			result.Homepage = strings.Trim(parts[1], " ")
			break
		case "Description":
			result.Description = strings.Trim(parts[1], " ")
			break
		case "Depends":
			result.Depends = strings.Trim(parts[1], " ")
			break
		}
	}
	return result, nil
}

func inspectPackage(f []byte) (string, error) {
	arReader := ar.NewReader(bytes.NewReader(f))
	var (
		controlBuffer bytes.Buffer
		header        *ar.Header
		err           error
	)

	for {
		header, err = arReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error in inspectPackage loop: %s", err)
		}
		if strings.Contains(header.Name, "control.tar") {
			var compression Compression
			if strings.TrimRight(header.Name, "/") == "control.tar.gz" {
				compression = GZIP
			} else if strings.TrimRight(header.Name, "/") == "control.tar.xz" {
				compression = LZMA
			} else {
				return "", fmt.Errorf("%s", "No control file found")
			}
			_, err = io.Copy(&controlBuffer, arReader)
			if err != nil {
				return "", err
			}
			return inspectPackageControl(compression, controlBuffer)
		}
	}
	return "", nil
}

func inspectPackageControl(compression Compression, fileName bytes.Buffer) (result string, err error) {
	var tarReader *tar.Reader
	switch compression {
	case GZIP:
		var compressedFile *gzip.Reader
		compressedFile, err = gzip.NewReader(bytes.NewReader(fileName.Bytes()))
		tarReader = tar.NewReader(compressedFile)
		break
	case LZMA:
		var compressedFile *lzma.Reader
		compressedFile, err = lzma.NewReader(bytes.NewReader(fileName.Bytes()))
		tarReader = tar.NewReader(compressedFile)
		break
	}
	var (
		controlBuffer bytes.Buffer
		header        *tar.Header
	)
	for {
		header, err = tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to inspect package: %s", err)
		}
		name := header.Name
		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			switch name {
			case "control", "./control":
				_, err = io.Copy(&controlBuffer, tarReader)
				if err != nil {
					return "", fmt.Errorf("can not copy file: %v", err)
				}
				return strings.TrimRight(controlBuffer.String(), "\n") + "\n", nil
			}
		default:
			return "", fmt.Errorf("Unable to figure out type : %c in file %s\n", header.Typeflag, name)
		}
	}
	return "", nil
}
