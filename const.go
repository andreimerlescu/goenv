package main

const (
	EnvNeverWriteProduction        = "GOENV_NEVER_WRITE_PRODUCTION"
	AmGoEnvAlwaysWrite      string = "AM_GO_ENV_ALWAYS_WRITE"
	AmGoEnvAlwaysUseJson    string = "AM_GO_ENV_ALWAYS_USE_JSON"
	AmGoEnvAlwaysUseYaml    string = "AM_GO_ENV_ALWAYS_USE_YAML"
	AmGoEnvAlwaysUseXml     string = "AM_GO_ENV_ALWAYS_USE_XML"
	AmGoEnvAlwaysUseToml    string = "AM_GO_ENV_ALWAYS_USE_TOML"
	AmGoEnvAlwaysUseIni     string = "AM_GO_ENV_ALWAYS_USE_INI"
	AmGoEnvAlwaysPrint      string = "AM_GO_ENV_ALWAYS_PRINT"

	envFileDefault     string = ".env"
	envFileLocal       string = ".env.local"
	envFileDevelopment string = ".env.development"
	envFileProduction  string = ".env.production"
)
