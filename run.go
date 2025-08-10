package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/andreimerlescu/figtree/v2"
	"github.com/andreimerlescu/goenv/env"
)

// Run is the primary goenv application
//
// Parameters:
//   	figs: The CLI state of arguments and inputs to the runtime
func Run(figs figtree.Plant) {
	showVersion := *figs.Bool(argVersion)
	if showVersion {
		fmt.Println(Version())
		os.Exit(0)
	}

	state := &stateful{
		Path: *figs.String(argEnvFile),
		Envs: []string{},
		Info: fileInfo{},

		env:   *figs.String(argEnv),
		value: *figs.String(argValue),

		mkAll:   *figs.Bool(argMkAll),
		init:    *figs.Bool(argInit),
		printer: *figs.Bool(argPrint),

		add:   *figs.Bool(argAdd),
		rm:    *figs.Bool(argRm),
		write: *figs.Bool(argWrite),

		is:  *figs.Bool(argIs),
		has: *figs.Bool(argHas),
		not: *figs.Bool(argNot),

		prod:          *figs.Bool(argProd),
		isProd:        strings.Contains(*figs.String(argEnvFile), "prod"),
		prodProtected: strings.Contains(*figs.String(argEnvFile), "prod") == true,

		toIni:  *figs.Bool(argIni),
		toYaml: *figs.Bool(argYaml),
		toXml:  *figs.Bool(argXml),
		toJson: *figs.Bool(argJson),
		toToml: *figs.Bool(argToml),
	}

	if len(state.Path) == 0 && (state.write || state.init) {
		// when no path is provided
		if state.prod {
			state.Path = ".env.production"
			state.isProd = true
		} else {
			state.Path = ".env"
			state.isProd = false
		}
	}
	d, err := os.Stat(state.Path)
	if os.IsNotExist(err) {
		if !state.init && !state.write {
			_, _ = fmt.Fprintf(os.Stderr, "%s does not exists, use -write to create\n", state.Path)
			os.Exit(1)
		}
	}
	if d != nil {
		state.Info.Name = d.Name()
		state.Info.ModTime = d.ModTime()
		state.Info.IsDir = d.IsDir()
		state.Info.Size = d.Size()
		state.Info.Mode = d.Mode()
	}

	state.isProd = strings.Contains(state.Path, ".env.production") || state.prod
	if state.isProd {
		fmt.Println("Using PRODUCTION environment file")
	} else if *figs.Bool(argVerbose) {
		fmt.Printf("Using %s environment file", state.Path)
	}
	state.prodProtected = env.Bool(EnvNeverWriteProduction, state.isProd)

	if *figs.Bool(argCleanAll) {
		for _, ext := range []string{outFormatJson, outFormatYaml, outFormatToml, outFormatXml, outFormatIni} {
			ext := ext
			err = nil
			path := state.Path + ext
			_, err = os.Stat(path)
			if os.IsNotExist(err) {
				continue
			}
			if os.IsPermission(err) {
				_, _ = fmt.Fprintf(os.Stderr, "%s is not writable: %v\n", path, err)
				continue
			}
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "path %s err = %v\n", path, err)
			} else {
				if !env.Bool(AmGoEnvNeverDelete, !(state.write || state.rm)) {
					err = os.Remove(path)
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "%s is not writable: %v\n", path, err)
						os.Exit(1)
					}
				}
				if !(state.write || state.rm) {
					fmt.Printf("The -write flag can be used to remove %s\n", path)
				}
			}
		}
		os.Exit(0)
	}

	Sanity(figs, state)
	triedWrite := false
retry:
	_, err = os.Lstat(state.Path)
	if os.IsNotExist(err) {
		if state.init && !triedWrite {
			if writeErr := os.WriteFile(state.Path, []byte{}, 0644); writeErr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "-init failed with: %v", writeErr)
				os.Exit(1)
			}
			os.Exit(0)
		} else if state.write && !triedWrite {
			var bb bytes.Buffer
			bb.WriteString(state.env)
			bb.WriteString("=")
			bb.WriteString(state.value)
			bb.WriteString("\n")
			if writeErr := os.WriteFile(state.Path, bb.Bytes(), 0644); writeErr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error writing %d bytes to %s due to %v", bb.Len(), state.Path, errors.Join(err, writeErr))
				os.Exit(1)
			}
			triedWrite = true
			goto retry
		}
		if !triedWrite {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if os.IsPermission(err) {
		_, _ = fmt.Fprintln(os.Stderr, "Error: permission denied")
		os.Exit(1)
	}

	var contents []byte
	contents, err = os.ReadFile(state.Path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "os.ReadFile(%s) returned err: %v", state.Path, err)
		os.Exit(1)
	}

	if size := len(contents); size == 0 && !(state.init || state.write || state.add) {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s %d bytes", state.Path, size)
		os.Exit(1)
	}

	envs := make(map[string]string)
	for _, line := range strings.Split(string(contents), "\n") {
		line = strings.TrimSpace(line)
		if len(line) < 3 {
			continue
		}
		parts := strings.SplitN(line, env.MapItemSeparator, env.MapSplitN)
		if len(parts) != 2 {
			continue
		}

		isThis := strings.EqualFold(strings.TrimSpace(parts[0]), strings.TrimSpace(state.env))
		if state.rm && isThis {
			continue
		}
		if state.has && isThis {
			code := 0
			if state.not {
				code = 1
			}
			if state.printer {
				if code == 1 {
					fmt.Println("YES")
				} else {
					fmt.Println("NO")
				}
			}
			os.Exit(code)
		}

		isThat := strings.EqualFold(strings.TrimSpace(parts[1]), strings.TrimSpace(state.value))
		if state.rm && isThat {
			continue
		}
		if state.is && isThat {
			code := 0
			if *figs.Bool(argNot) {
				code = 1
			}
			if state.printer {
				if code == 1 {
					fmt.Println("YES")
				} else {
					fmt.Println("NO")
				}
			}
			os.Exit(code)
		}

		envs[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	wrote := false
	if state.add {
		if _, exists := envs[strings.TrimSpace(state.env)]; !exists {
			envs[strings.TrimSpace(state.env)] = strings.TrimSpace(state.value)
			wrote = true
		}
	}

	if wrote && state.has {
		for k, v := range envs {
			if strings.EqualFold(strings.TrimSpace(state.env), strings.TrimSpace(k)) {
				code := 0
				if *figs.Bool(argNot) {
					code = 1
				}
				if state.printer {
					if code == 1 {
						fmt.Println(v)
					} else {
						fmt.Println("NO")
					}
				} else {
					os.Exit(code)
				}
			}
		}
	}

	if wrote && state.is {
		for _, v := range envs {
			if strings.EqualFold(strings.TrimSpace(state.value), strings.TrimSpace(v)) {
				code := 0
				if *figs.Bool(argNot) {
					code = 1
				}
				if state.printer {
					if code == 1 {
						fmt.Println(v)
					} else {
						fmt.Println("NO")
					}
				} else {
					os.Exit(code)
				}
			}
		}
	}

	Result(figs, envs, state)
}
