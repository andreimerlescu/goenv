package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/andreimerlescu/figtree/v2"
	"github.com/andreimerlescu/goenv/env"
)

func Result(
	figs figtree.Plant,
	envs map[string]string,
	state *stateful,
) {

	for e, v := range envs {
		e, v = strings.TrimSpace(e), strings.TrimSpace(v)
		state.Envs = append(state.Envs, fmt.Sprintf("%s=%s", e, v))
	}

	if state.toJson {
		processJson(figs, envs, state)
	}

	if state.toIni {
		processIni(figs, envs, state)
	}

	if state.toYaml {
		processYaml(figs, envs, state)
	}

	if state.toToml {
		processToml(figs, envs, state)
	}

	if state.toXml {
		processXml(figs, envs, state)
	}

	var out bytes.Buffer
	for e, v := range envs {
		e, v = strings.TrimSpace(e), strings.TrimSpace(v)
		out.WriteString(fmt.Sprintf("%s=%s\n", e, v))
	}
	if state.write {
		if env.Bool(EnvNeverWriteProduction, state.prodProtected) {
			_, _ = fmt.Fprintln(os.Stderr, "HALT: PRODUCTION IS PROTECTED! WRITE OPERATION CANCELED.")
			os.Exit(1)
		}
		if writeErr := os.WriteFile(state.Path, out.Bytes(), 0644); writeErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error writing file %s: %s", state.Path, writeErr)
			os.Exit(1)
		}
		os.Exit(0)
	} else if *figs.Bool(argPrint) {
		fmt.Println(out.String())
		os.Exit(0)
	}
	if *figs.Bool(argVerbose) {
		fmt.Println("Finished executing!")
		os.Exit(0)
	}
	os.Exit(0)
}
