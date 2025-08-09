package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/andreimerlescu/figtree/v2"
)

func processJson(figs figtree.Plant, envs map[string]string, state *stateful) {
	output, err := json.MarshalIndent(envs, "", "  ")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error marshalling environment variable:", state.env)
		os.Exit(1)
	}
	var bb bytes.Buffer
	bb.Write(output)
	writeProcessed(figs, &bb, ".json", state)
}

func processIni(figs figtree.Plant, envs map[string]string, state *stateful) {
	var ini bytes.Buffer
	ini.WriteString("[default]\n")
	for e, v := range envs {
		e, v = strings.TrimSpace(e), strings.TrimSpace(v)
		ini.WriteString(fmt.Sprintf("%s = %s\n", e, v))
	}
	writeProcessed(figs, &ini, ".ini", state)
}

func processToml(figs figtree.Plant, envs map[string]string, state *stateful) {
	var toml bytes.Buffer
	for e, v := range envs {
		e, v = strings.TrimSpace(e), strings.TrimSpace(v)
		toml.WriteString(fmt.Sprintf("%s: \"%s\" \n", e, v))
	}
	writeProcessed(figs, &toml, ".toml", state)
}

func processYaml(figs figtree.Plant, envs map[string]string, state *stateful) {
	var yaml bytes.Buffer
	yaml.WriteString("---\n")
	for e, v := range envs {
		e, v = strings.TrimSpace(e), strings.TrimSpace(v)
		yaml.WriteString(fmt.Sprintf("%s: \"%s\" \n", e, v))
	}
	writeProcessed(figs, &yaml, ".yaml", state)
}

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
	writeProcessed(figs, &xml, ".xml", state)
}

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
	fmt.Println(buf.String())
	os.Exit(0)
}
