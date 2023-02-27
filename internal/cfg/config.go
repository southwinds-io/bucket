/*
   Bucket - Debian & RPM Package Repository
   ©2023 SouthWinds Tech Ltd
*/

package cfg

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Debian []DebianRepository
	Keys   []Key
}

func (c *Config) GetRepo(name string) *DebianRepository {
	for _, repository := range c.Debian {
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

type DebianRepository struct {
	Name          string   `json:"name" yaml:"name"`
	Distributions []string `json:"distributions" yaml:"distributions"`
	Architectures []string `json:"architectures" yaml:"architectures"`
	Sections      []string `json:"sections" yaml:"sections"`
	KeyRef        string   `json:"key_ref" yaml:"key_ref"`
	cfg           *Config
}

func (r *DebianRepository) GetConfig() *Config {
	return r.cfg
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

func DebianPkgName(repo, version, release, arch string) string {
	return fmt.Sprintf("%s-%s-%s-%s.deb", repo, version, release, arch)
}
