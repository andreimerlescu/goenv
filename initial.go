package main

import (
	"path/filepath"

	"github.com/andreimerlescu/checkfs"
	"github.com/andreimerlescu/checkfs/file"
)

func Initial() string {
	check := []string{
		filepath.Join(".", envFileDefault),
		filepath.Join(".", envFileLocal),
		filepath.Join(".", envFileDevelopment),
		filepath.Join(".", envFileProduction),
	}
	for _, c := range check {
		if err := checkfs.File(c, file.Options{Exists: true}); err == nil {
			return c
		}
	}
	return ""
}
