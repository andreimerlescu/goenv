package main

import (
	"fmt"
	"os"

	"github.com/andreimerlescu/figtree/v2"
)

// Sanity ensures that you're not exporting to multiple formats at once and uses argVerbose to print statements to STDOUT
//
// Parameters:
//   	figs: The CLI state of arguments and inputs to the runtime
// 	 	state: Read-Only verification on export options being singular in choice
//
// Usage:
// 		Sanity(figs, state)
//
// Panics:
// 		- When figs or state are nil, no sanity verification can take place
func Sanity(figs figtree.Plant, state *stateful) {
	if figs == nil || state == nil {
		panic("Sanity called with nil figtree or state!")
	}
	// #begin
	using := ""
	selectedOut := false

	// -json
	if state.toJson && *figs.Bool(argVerbose) {
		fmt.Println("Using JSON environment file")
		if state.write && state.isProd && state.prodProtected {
			fmt.Printf("We'll write to %s%s for you!", state.Path, outFormatJson)
		}
	}
	if state.toJson { // skip selectedOut here
		using = "json"
		selectedOut = true
	}

	// -ini
	if state.toIni && *figs.Bool(argVerbose) {
		fmt.Println("Using INI environment file")
		if state.write && state.isProd && state.prodProtected {
			fmt.Printf("We'll write to %s%s for you!", state.Path, outFormatIni)
		}
	}
	if state.toIni && selectedOut {
		_, _ = fmt.Fprintf(os.Stderr, "using %s write to %s.%s. ERROR CANNOT COMBINE -json -ini -toml -yaml \n", using, state.Path, using)
		os.Exit(1)
	} else if state.toIni && !selectedOut {
		selectedOut = true
		using = "ini"
	}

	// -yaml
	if state.toYaml && *figs.Bool(argVerbose) {
		fmt.Println("Using YAML environment file")
		if state.write && state.isProd && state.prodProtected {
			fmt.Printf("We'll write to %s%s for you!\n", state.Path, outFormatYaml)
		}
	}
	if state.toYaml && selectedOut {
		_, _ = fmt.Fprintf(os.Stderr, "using %s write to %s.%s. ERROR CANNOT COMBINE -json -ini -toml -yaml \n", using, state.Path, using)
		os.Exit(1)
	} else if state.toYaml && !selectedOut {
		selectedOut = true
		using = "yaml"
	}

	// -xml
	if state.toXml && *figs.Bool(argVerbose) {
		fmt.Println("Using XML environment file")
		if state.write && state.isProd && state.prodProtected {
			fmt.Printf("We'll write to %s%s for you!\n", state.Path, outFormatXml)
		}
	}
	if state.toXml && selectedOut {
		_, _ = fmt.Fprintf(os.Stderr, "using %s write to %s.%s. ERROR CANNOT COMBINE -json -ini -toml -yaml \n", using, state.Path, using)
		os.Exit(1)
	} else if state.toXml && !selectedOut {
		selectedOut = true
		using = "xml"
	}

	// -toml
	if state.toToml && *figs.Bool(argVerbose) {
		fmt.Println("Using TOML environment file")
		if state.write && state.isProd && state.prodProtected {
			fmt.Printf("We'll write to %s%s for you!\n", state.Path, outFormatToml)
		}
	}
	if state.toToml && selectedOut {
		_, _ = fmt.Fprintf(os.Stderr, "using %s write to %s.%s. ERROR CANNOT COMBINE -json -ini -toml -yaml \n", using, state.Path, using)
		os.Exit(1)
	} else if state.toToml && !selectedOut {
		selectedOut = true
		using = "toml"
	}

	// #done
	if !selectedOut && *figs.Bool(argVerbose) {
		fmt.Println("Not exporting the environment to a new file")
	}
}
