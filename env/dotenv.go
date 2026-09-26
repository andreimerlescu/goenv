package env

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// LoadFile reads the .env file at path (default EnvFile when path is empty) into the runtime environment.
// Variables that are already defined in the environment are NOT overwritten. A missing file is not an error.
//
// Example:
//     if err := env.LoadFile("./.env"); err != nil {
//        log.Fatal(err)
//     }
func LoadFile(path string) error {
	return loadFile(path, false)
}

// OverloadFile reads the .env file at path (default EnvFile when path is empty) into the runtime environment,
// overwriting any variables that are already defined. A missing file is not an error.
//
// Example:
//     err := env.OverloadFile(".env.local")
func OverloadFile(path string) error {
	return loadFile(path, true)
}

// ReadFile parses the .env file at path (default EnvFile when path is empty) and returns its key/value pairs
// without modifying the runtime environment.
//
// Example:
//     vars, err := env.ReadFile(".env")
func ReadFile(path string) (map[string]string, error) {
	if len(path) == 0 {
		path = EnvFile
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return Parse(f)
}

// Parse reads .env formatted content from r and returns its key/value pairs.
//
// Supported syntax:
//     # comments and blank lines are ignored
//     KEY=value
//     export KEY=value
//     KEY="double quoted value with \n escapes"
//     KEY='single quoted literal value'
//     KEY=value # inline comments on unquoted values
func Parse(r io.Reader) (map[string]string, error) {
	envs := make(map[string]string)
	scanner := bufio.NewScanner(r)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if lineNo == 1 {
			line = strings.TrimPrefix(line, "\uFEFF")
		}
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, found := strings.Cut(line, "=")
		if !found {
			return envs, fmt.Errorf("line %d: missing '=' in %q", lineNo, line)
		}
		key = strings.TrimSpace(key)
		if len(key) == 0 {
			return envs, fmt.Errorf("line %d: empty key", lineNo)
		}
		envs[key] = parseValue(strings.TrimSpace(value))
	}
	return envs, scanner.Err()
}

// parseValue unquotes value and strips inline comments from unquoted values
func parseValue(value string) string {
	if len(value) == 0 {
		return value
	}
	switch q := value[0]; q {
	case '"', '\'':
		if end := strings.IndexByte(value[1:], q); end >= 0 {
			inner := value[1 : end+1]
			if q == '"' {
				inner = strings.NewReplacer(`\n`, "\n", `\r`, "\r", `\t`, "\t", `\"`, `"`, `\\`, `\`).Replace(inner)
			}
			return inner
		}
		return value
	}
	if i := strings.Index(value, " #"); i >= 0 {
		value = strings.TrimSpace(value[:i])
	}
	return value
}

// loadFile is the shared implementation of LoadFile and OverloadFile
func loadFile(path string, overwrite bool) error {
	if len(path) == 0 {
		path = EnvFile
	}
	envs, err := ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		if ShowVerbose {
			OutLogger.Printf("LoadFile(%s) skipped: file does not exist", path)
		}
		return nil
	}
	if err != nil {
		logError(path, err)
		return err
	}
	for key, value := range envs {
		if _, exists := os.LookupEnv(key); exists && !overwrite {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			logError(key, err)
			return err
		}
	}
	if ShowVerbose {
		OutLogger.Printf("LoadFile(%s) loaded %d variables", path, len(envs))
	}
	return nil
}
