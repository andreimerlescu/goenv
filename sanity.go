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
	using := ""
	selectedOut := false

	if state.toJson && *figs.Bool(argVerbose) {
		fmt.Println("Using JSON environment file")
		if state.write && state.isProd && state.prodProtected {
			fmt.Printf("We'll write to %s.json for you!", state.Path)
		}
	}
	if state.toJson { // skip selectedOut here
		using = "json"
		selectedOut = true
	}

	if state.toIni && *figs.Bool(argVerbose) {
		fmt.Println("Using INI environment file")
		if state.write && state.isProd && state.prodProtected {
			fmt.Printf("We'll write to %s.ini for you!", state.Path)
		}
	}
	if state.toIni && selectedOut {
		_, _ = fmt.Fprintf(os.Stderr, "using %s write to %s.%s. ERROR CANNOT COMBINE -json -ini -toml -yaml \n", using, state.Path, using)
		os.Exit(1)
	} else if state.toIni && !selectedOut {
		selectedOut = true
		using = "ini"
	}

	if state.toYaml && *figs.Bool(argVerbose) {
		fmt.Println("Using YAML environment file")
		if state.write && state.isProd && state.prodProtected {
			fmt.Printf("We'll write to %s.yaml for you!\n", state.Path)
		}
	}
	if state.toYaml && selectedOut {
		_, _ = fmt.Fprintf(os.Stderr, "using %s write to %s.%s. ERROR CANNOT COMBINE -json -ini -toml -yaml \n", using, state.Path, using)
		os.Exit(1)
	} else if state.toYaml && !selectedOut {
		selectedOut = true
		using = "yaml"
	}

	if state.toToml && *figs.Bool(argVerbose) {
		fmt.Println("Using TOML environment file")
		if state.write && state.isProd && state.prodProtected {
			fmt.Printf("We'll write to %s.toml for you!\n", state.Path)
		}
	}
	if state.toToml && selectedOut {
		_, _ = fmt.Fprintf(os.Stderr, "using %s write to %s.%s. ERROR CANNOT COMBINE -json -ini -toml -yaml \n", using, state.Path, using)
	} else if state.toToml && !selectedOut {
		selectedOut = true
		using = "toml"
	}

	if !selectedOut && *figs.Bool(argVerbose) {
		fmt.Println("Not exporting the environment to a new file")
	}
}
