/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package deb

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Repositories []Repository
	Keys         []Key
}

func (c *Config) GetRepo(name string) *Repository {
	for _, repository := range c.Repositories {
		if strings.EqualFold(repository.Name, name) {
			repository.cfg = c
			return &repository
		}
	}
	return nil
}

func (c *Config) GetKey(ref string) (*Key, bool) {
	for _, key := range c.Keys {
		if strings.EqualFold(key.Ref, ref) {
			return &key, true
		}
	}
	return nil, false
}

type Repository struct {
	Name          string   `json:"name" yaml:"name"`
	Distribution  string   `json:"distribution" yaml:"distribution"`
	Architectures []string `json:"architectures" yaml:"architectures"`
	Sections      []string `json:"sections" yaml:"sections"`
	KeyRef        string   `json:"key_ref" yaml:"key_ref"`
	cfg           *Config
}

type Key struct {
	Ref      string `json:"ref" yaml:"ref"`
	Private  string `json:"private" yaml:"private"`
	Public   string `json:"public" yaml:"public"`
	Passcode string `json:"passcode" yaml:"passcode"`
}

func NewConfig() (cfg *Config, err error) {
	cfg = new(Config)
	path, err := getConfigPath()
	if err != nil {
		return
	}
	var b []byte
	b, err = os.ReadFile(filepath.Join(path, "config.yaml"))
	if err != nil {
		return
	}
	err = yaml.Unmarshal(b, cfg)
	return
}

func pkgName(repo, version, release, arch string) string {
	return fmt.Sprintf("%s-%s-%s-%s.deb", repo, version, release, arch)
}
