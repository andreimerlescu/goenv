package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/andreimerlescu/checkfs"
	"github.com/andreimerlescu/checkfs/file"
	"github.com/andreimerlescu/figtree/v2"
	"github.com/andreimerlescu/goenv/env"
)

// NewConfiguration returns a new figtree.Plant that contains each configurable that begins with argEnvFile
func NewConfiguration() figtree.Plant {
	love := figtree.Options{
		IgnoreEnvironment: true,
		Germinate:         true,
	}

	defaultConfigFile := filepath.Join(env.User().HomeDir, ".config", "goenv", "config.yml")
	configFile := env.String("AM_GO_ENV_CONFIG_FILE", defaultConfigFile)
	if err := checkfs.File(configFile, file.Options{Exists: true}); err != nil {
		figtree.ConfigFilePath = defaultConfigFile
		love.ConfigFile = configFile
	}

	figs := figtree.With(love)

	figs = figs.NewString(argEnvFile, Initial(), "Path to env file to process")
	figs = figs.NewString(argEnv, "", "Check for an environment variable name")
	figs = figs.NewString(argValue, "", "Check for an environment variable value. Use with -"+argEnv)
	figs = figs.NewBool(argAdd, false, "Add a new environment variable")
	figs = figs.NewBool(argRm, false, "Remove a new environment variable")
	figs = figs.NewBool(argProd, false, "Allow interaction with .env.production")
	figs = figs.NewBool(argHas, false, "Check for a new environment variable")
	figs = figs.NewBool(argVerbose, false, "Show verbose output")
	figs = figs.NewBool(argVersion, false, "Show version")
	figs = figs.NewBool(argIs, false, "Logic for checking a value in the file")
	figs = figs.NewBool(argWrite, env.Bool(AmGoEnvAlwaysWrite, false), "Write to the -"+argEnvFile)
	figs = figs.NewBool(argJson, env.Bool(AmGoEnvAlwaysUseJson, false), "Output in JSON format")
	figs = figs.NewBool(argYaml, env.Bool(AmGoEnvAlwaysUseYaml, false), "Output in YAML format")
	figs = figs.NewBool(argXml, env.Bool(AmGoEnvAlwaysUseXml, false), "Output in XML format")
	figs = figs.NewBool(argToml, env.Bool(AmGoEnvAlwaysUseToml, false), "Output in TOML format")
	figs = figs.NewBool(argIni, env.Bool(AmGoEnvAlwaysUseIni, false), "Output in INI format")
	figs = figs.NewBool(argPrint, env.Bool(AmGoEnvAlwaysPrint, false), "Always print the contents of the env before exiting upon success")
	figs = figs.NewBool(argNot, false, "Negates -has or -is")
	figs = figs.NewBool(argInit, false, "Create the -file if it does not exist")
	figs = figs.NewBool(argMkAll, false, "Will create all -json -xml -toml -ini -yaml output formats of -file")
	figs = figs.NewBool(argCleanAll, false, "Remove all -json -xml -toml -ini")

	if err := figs.Load(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	return figs
}
