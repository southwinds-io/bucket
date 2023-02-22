/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package deb

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

// PackagesData manage "Packages" and "Packages.gz" files
type PackagesData struct {
	Items []*PackageData
}

func newPackagesDataFromContent(content string) *PackagesData {
	return &PackagesData{
		Items: parsePackages(content),
	}
}

func NewPackagesData(filename string) (*PackagesData, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			// Packages file does not exist therefore returns an empty instance
			return &PackagesData{
				Items: []*PackageData{},
			}, nil
		}
		return nil, err
	}
	return &PackagesData{
		Items: parsePackages(string(content[:])),
	}, nil
}

func (p *PackagesData) Add(data PackageData) error {
	if p.Find(data.Package, data.Version, data.Release, data.Architecture) != nil {
		return fmt.Errorf("package already exists")
	}
	p.Items = append(p.Items, &data)
	return nil
}

func (p *PackagesData) Remove(packageName, version string) bool {
	for i := 0; i < len(p.Items); i++ {
		if strings.EqualFold(p.Items[i].Package, packageName) &&
			strings.EqualFold(p.Items[i].Version, version) {
			p.Items = remove(p.Items, i)
			return true
		}
	}
	return false
}

func (p *PackagesData) Find(packageName, version, release, arc string) *PackageData {
	for _, item := range p.Items {
		if strings.EqualFold(item.Package, packageName) &&
			strings.EqualFold(item.Version, version) &&
			strings.EqualFold(item.Release, release) &&
			strings.EqualFold(item.Architecture, arc) {
			return item
		}
	}
	return nil
}

func (p *PackagesData) String() string {
	buffer := new(bytes.Buffer)
	for ix, item := range p.Items {
		buffer.WriteString(item.String())
		if ix < len(p.Items)-1 {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}

func (p *PackagesData) Save(path string) error {
	// if the Packages contain no items
	if len(p.Items) == 0 {
		// delete it altogether
		return os.Remove(path)
	}
	// otherwise write it
	return os.WriteFile(path, []byte(p.String()), 0755)
}

func (p *PackagesData) SaveGz(path string) error {
	// if the Packages contain no items
	if len(p.Items) == 0 {
		// delete Packages.gz
		return os.Remove(path)
	}
	packageGzFile, err := os.Create(path)
	if err != nil {
		return err
	}
	gzOut := gzip.NewWriter(packageGzFile)
	defer func(gzOut *gzip.Writer) {
		if err = gzOut.Close(); err != nil {
			return
		}
	}(gzOut)
	_, err = gzOut.Write([]byte(p.String()))
	if err != nil {
		return err
	}
	return gzOut.Flush()
}

type PackageData struct {
	Package      string
	Version      string
	Release      string
	Architecture string
	Maintainer   string
	Homepage     string
	Depends      string
	Description  string
	Filename     string
	Size         string
	MD5sum       string
	SHA1         string
	SHA256       string
}

func (p *PackageData) String() string {
	buffer := new(bytes.Buffer)
	buffer.WriteString(fmt.Sprintf("Package: %s\n", p.Package))
	buffer.WriteString(fmt.Sprintf("Version: %s\n", p.Version))
	buffer.WriteString(fmt.Sprintf("Release: %s\n", p.Release))
	buffer.WriteString(fmt.Sprintf("Architecture: %s\n", p.Architecture))
	buffer.WriteString(fmt.Sprintf("Maintainer: %s\n", p.Maintainer))
	buffer.WriteString(fmt.Sprintf("Homepage: %s\n", p.Homepage))
	buffer.WriteString(fmt.Sprintf("Depends: %s\n", p.Depends))
	buffer.WriteString(fmt.Sprintf("Description: %s\n", p.Description))
	buffer.WriteString(fmt.Sprintf("Filename: %s\n", p.Filename))
	buffer.WriteString(fmt.Sprintf("Size: %s\n", p.Size))
	buffer.WriteString(fmt.Sprintf("MD5sum: %s\n", p.MD5sum))
	buffer.WriteString(fmt.Sprintf("SHA1: %s\n", p.SHA1))
	buffer.WriteString(fmt.Sprintf("SHA256: %s\n", p.SHA256))
	return buffer.String()
}

func parsePackages(content string) []*PackageData {
	result := make([]*PackageData, 0)
	parts := strings.Split(content, "\n")
	var data *PackageData
	var key, value string
	for ix, part := range parts {
		if len(strings.Trim(part, " ")) == 0 {
			continue
		}
		items := strings.SplitN(part, ":", 2)
		if len(items) < 2 {
			continue
		}
		key = strings.Trim(items[0], " ")
		value = strings.Trim(items[1], " ")
		// lookahead
		if ix < len(parts) && len(strings.SplitN(parts[ix+1], ":", 2)) < 2 {
			value += fmt.Sprintf("\n%s", parts[ix+1])
		}
		if strings.EqualFold(key, "package") {
			if data == nil {
				data = new(PackageData)
				data.Package = value
			}
		} else {
			switch key {
			case "Version":
				data.Version = value
				break
			case "Release":
				data.Release = value
				break
			case "Architecture":
				data.Architecture = value
				break
			case "Maintainer":
				data.Maintainer = value
				break
			case "Homepage":
				data.Homepage = value
				break
			case "Depends":
				data.Depends = value
				break
			case "Description":
				data.Description = value
				break
			case "Filename":
				data.Filename = value
				break
			case "Size":
				data.Size = value
				break
			case "MD5sum":
				data.MD5sum = value
				break
			case "SHA1":
				data.SHA1 = value
				break
			case "SHA256":
				data.SHA256 = value
				result = append(result, data)
				break
			}
		}
	}
	return result
}

func remove(slice []*PackageData, s int) []*PackageData {
	return append(slice[:s], slice[s+1:]...)
}

func getChecksums(pkg []byte) (md5Sum, sha1Sum, sha256Sum string, err error) {
	var (
		md5hash    = md5.New()
		sha1hash   = sha1.New()
		sha256hash = sha256.New()
	)
	_, err = io.Copy(io.MultiWriter(md5hash, sha1hash, sha256hash), bytes.NewReader(pkg))
	if err != nil {
		return "", "", "", fmt.Errorf("error hashing debian package: %s", err)
	}
	return hex.EncodeToString(md5hash.Sum(nil)), hex.EncodeToString(sha1hash.Sum(nil)), hex.EncodeToString(sha256hash.Sum(nil)), nil
}
