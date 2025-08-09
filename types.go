package main

import (
	"os"
	"time"
)

type (
	fileInfo struct {
		Name    string      `json:"name" yaml:"name" toml:"name" xml:"name" ini:"name"`
		Size    int64       `json:"size" yaml:"size" toml:"size" xml:"size" ini:"size"`
		Mode    os.FileMode `json:"mode" yaml:"mode" toml:"mode" xml:"mode" ini:"mode"`
		IsDir   bool        `json:"isdir" yaml:"isdir" toml:"isdir" xml:"isdir" ini:"isdir"`
		ModTime time.Time   `json:"modtime" yaml:"modtime" toml:"modtime" xml:"modtime" ini:"modtime"`
	}

	stateful struct {
		Path string   `json:"path" yaml:"path" toml:"path" xml:"path" ini:"path"`
		Info fileInfo `json:"info" yaml:"info" toml:"info" xml:"info" ini:"info"`
		Envs []string `json:"envs" yaml:"envs" toml:"envs" xml:"envs" ini:"envs"`

		env, value                                                         string
		is, not, has, printer, add, rm, write, prod, isProd, prodProtected bool
		toJson, toYaml, toXml, toIni, toToml                               bool
	}
)
