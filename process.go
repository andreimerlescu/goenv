package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/andreimerlescu/figtree/v2"
)

// processJson renders the argEnvFile with an ext of outFormatJson
//
// Parameters:
//   	figs: The CLI state of arguments and inputs to the runtime
// 	 	state: Read-Only verification on export options being singular in choice
// 		envs: map of environment variables as key=value pairs
func processJson(figs figtree.Plant, envs map[string]string, state *stateful) {
	output, err := json.MarshalIndent(envs, "", "  ")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error marshalling environment variable:", state.env)
		os.Exit(1)
	}
	var bb bytes.Buffer
	bb.Write(output)
	writeProcessed(figs, &bb, outFormatJson, state)
}

// processIni renders the argEnvFile with an ext of outFormatIni
//
// Parameters:
//   	figs: The CLI state of arguments and inputs to the runtime
// 	 	state: Read-Only verification on export options being singular in choice
// 		envs: map of environment variables as key=value pairs
func processIni(figs figtree.Plant, envs map[string]string, state *stateful) {
	var ini bytes.Buffer
	ini.WriteString("[default]\n")
	for e, v := range envs {
		e, v = strings.TrimSpace(e), strings.TrimSpace(v)
		ini.WriteString(fmt.Sprintf("%s = %s\n", e, v))
	}
	writeProcessed(figs, &ini, outFormatIni, state)
}

// processToml renders the argEnvFile with an ext of outFormatToml
//
// Parameters:
//   	figs: The CLI state of arguments and inputs to the runtime
// 	 	state: Read-Only verification on export options being singular in choice
// 		envs: map of environment variables as key=value pairs
func processToml(figs figtree.Plant, envs map[string]string, state *stateful) {
	var toml bytes.Buffer
	for e, v := range envs {
		e, v = strings.TrimSpace(e), strings.TrimSpace(v)
		toml.WriteString(fmt.Sprintf("%s: \"%s\" \n", e, v))
	}
	writeProcessed(figs, &toml, outFormatToml, state)
}

// processYaml renders the argEnvFile with an ext of outFormatYaml
//
// Parameters:
//   	figs: The CLI state of arguments and inputs to the runtime
// 	 	state: Read-Only verification on export options being singular in choice
// 		envs: map of environment variables as key=value pairs
func processYaml(figs figtree.Plant, envs map[string]string, state *stateful) {
	var yaml bytes.Buffer
	yaml.WriteString("---\n")
	for e, v := range envs {
		e, v = strings.TrimSpace(e), strings.TrimSpace(v)
		yaml.WriteString(fmt.Sprintf("%s: \"%s\" \n", e, v))
	}
	writeProcessed(figs, &yaml, outFormatYaml, state)
}

// processXml renders the argEnvFile with an ext of outFormatXml
//
// Parameters:
//   	figs: The CLI state of arguments and inputs to the runtime
// 	 	state: Read-Only verification on export options being singular in choice
// 		envs: map of environment variables as key=value pairs
func processXml(figs figtree.Plant, envs map[string]string, state *stateful) {
	var xml bytes.Buffer
	xml.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	xml.WriteString("<env>\n")
	for e, v := range envs {
		i := "   "
		e, v = strings.TrimSpace(e), strings.TrimSpace(v)

		xml.WriteString(i + "<")
		xml.WriteString(e)
		xml.WriteString(">")
		xml.WriteString(v)
		xml.WriteString("</")
		xml.WriteString(e)
		xml.WriteString(">\n")
	}
	xml.WriteString("</env>\n")
	writeProcessed(figs, &xml, outFormatXml, state)
}

// writeProcessed renders the argEnvFile + extension with the buffered bytes
//
// Parameters:
//   	figs: The CLI state of arguments and inputs to the runtime
// 	 	state: Read-Only verification on export options being singular in choice
// 		envs: map of environment variables as key=value pairs
func writeProcessed(figs figtree.Plant, buf *bytes.Buffer, ext string, state *stateful) {
	if state.write {
		path := fmt.Sprintf("%s%s", state.Path, ext)
		if writeErr := os.WriteFile(path, buf.Bytes(), 0644); writeErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error writing file %s: %s", path, writeErr)
			os.Exit(1)
		}
		if *figs.Bool(argVerbose) {
			fmt.Printf("Writing file %s\n", path)
		}
	}
	if !state.mkAll {
		fmt.Println(buf.String())
		os.Exit(0)
	}
}
