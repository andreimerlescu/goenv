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
		is:    *figs.Bool(argIs),
		has:   *figs.Bool(argHas),
		add:   *figs.Bool(argAdd),
		rm:    *figs.Bool(argRm),
		write: *figs.Bool(argWrite),
		prod:  *figs.Bool(argProd),

		toIni:  *figs.Bool(argIni),
		toYaml: *figs.Bool(argYaml),
		toXml:  *figs.Bool(argXml),
		toJson: *figs.Bool(argJson),
		toToml: *figs.Bool(argToml),
	}

	if len(state.Path) == 0 && state.write {
		// when no path is provided
		if state.prod {
			state.Path = ".env.production"
		} else {
			state.Path = ".env"
		}
	}
	d, err := os.Stat(state.Path)
	if os.IsNotExist(err) && !state.write {
		_, _ = fmt.Fprintf(os.Stderr, "%s does not exists, use -write to create\n", state.Path)
		os.Exit(1)
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

	Sanity(figs, state)

retry:
	_, err = os.Lstat(state.Path)
	if os.IsNotExist(err) {
		_, _ = fmt.Fprintln(os.Stderr, err)
		if state.write {
			var bb bytes.Buffer
			bb.WriteString(state.env)
			bb.WriteString("=")
			bb.WriteString(state.value)
			bb.WriteString("\n")
			if writeErr := os.WriteFile(state.Path, bb.Bytes(), 0644); writeErr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error writing %d bytes to %s due to %v", bb.Len(), state.Path, errors.Join(err, writeErr))
				os.Exit(1)
			}
			goto retry
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

	if size := len(contents); size == 0 {
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
