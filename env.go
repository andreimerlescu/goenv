package env

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

// ** EXISTENCE **

// Exists checks if an environment variable is set.
//
// Example:
//     ok := env.Exists("HOSTNAME")
func Exists(env string) (ok bool) {
	_, ok = os.LookupEnv(env)
	if !ok && ShowVerbose {
		OutLogger.Printf("Exists(%s) = %v", env, ok)
	}
	return
}

// MustExist panics or exits if an environment variable is not set.
// The behavior is controlled by the AllowPanic and PrintErrors global variables.
//
// Example:
//     func init() {
//        env.MustExist("REQUIRED_VAR")
//     }
func MustExist(env string) {
	if _, ok := os.LookupEnv(env); !ok {
		msg := fmt.Sprintf("required environment variable '%s' is not set", env)
		if AllowPanic {
			panic(msg)
		}
		logError(env, errors.New(msg))
		os.Exit(1)
	}
	if ShowVerbose {
		OutLogger.Printf("MustExist() confirmed %s exists!", env)
	}
}

// ** TRUTHY **

// IsFalse returns true if an environment variable is "false" or is not set.
//
// Example:
//     if env.IsFalse("USE_JSON") {
//        fmt.Println("Not using JSON.")
//     }
func IsFalse(env string) (ok bool) {
	// If the value is "true", Bool(env, false) will be true, so we return false.
	// If the value is "false" or unset, Bool(env, false) will be false, so we return true.
	ok = !Bool(env, false)
	if ShowVerbose {
		OutLogger.Printf("IsFalse(%s) = %v", env, ok)
	}
	return
}

// IsTrue returns true only if an environment variable is explicitly "true".
//
// Example:
//     if env.IsTrue("USE_JSON") {
//        // intend to use json
//     }
func IsTrue(env string) (ok bool) {
	ok = Bool(env, false)
	if ShowVerbose {
		OutLogger.Printf("IsTrue(%s) = %v", env, ok)
	}
	return
}

// AreTrue returns true only if all specified environment variables are true.
//
// Example:
//     if env.AreTrue("ENABLED", "FEATURE_X_ON") {
//        // ...
//     }
func AreTrue(envs ...string) (ok bool) {
	for _, env := range envs {
		if !IsTrue(env) {
			if ShowVerbose {
				OutLogger.Printf("AreTrue(%s) = false (failed on %s)", strings.Join(envs, ", "), env)
			}
			return false
		}
	}
	return true
}

// AreFalse returns true only if all specified environment variables are false or unset.
//
// Example:
//     if env.AreFalse("DISABLED", "LEGACY_MODE") {
//        // ...
//     }
func AreFalse(envs ...string) (ok bool) {
	for _, env := range envs {
		if !IsFalse(env) {
			if ShowVerbose {
				OutLogger.Printf("AreFalse(%s) = false (failed on %s)", strings.Join(envs, ", "), env)
			}
			return false // Short-circuit
		}
	}
	return true
}

// Bool parses an environment variable as a boolean.
// It accepts "true", "t", "1" as true and "false", "f", "0" as false.
// If the variable is unset or fails to parse, it returns the fallback value.
//
// Example:
//     useJson := env.Bool("USE_JSON", false)
func Bool(env string, fallback bool) bool {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	vb, err := strconv.ParseBool(v)
	if err != nil {
		logError(env, err)
		return fallback
	}
	return vb
}

// ** TYPES **

// String gets a string value from an environment variable, or a fallback if not set.
//
// Example:
//     hostname := env.String("HOSTNAME", "localhost")
func String(env string, fallback string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	return v
}

// Int gets an integer value from an environment variable, or a fallback if not set or invalid.
//
// Example:
//     port := env.Int("PORT", 3306)
func Int(env string, fallback int) int {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	vint, err := strconv.Atoi(v)
	if err != nil {
		logError(env, err)
		return fallback
	}
	return vint
}

// Int64 gets an int64 value from an environment variable, or a fallback if not set or invalid.
//
// Example:
//     ns := env.Int64("NANOSECONDS", int64(1_000_000_000))
func Int64(env string, fallback int64) int64 {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	vint, err := strconv.ParseInt(v, Int64Base, Int64BitSize)
	if err != nil {
		logError(env, err)
		return fallback
	}
	return vint
}

// Float32 gets a float32 value from an environment variable, or a fallback if not set or invalid.
//
// Example:
//     pi := env.Float32("PI", float32(3.14))
func Float32(env string, fallback float32) float32 {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	vf, err := strconv.ParseFloat(v, Float32BitSize)
	if err != nil {
		logError(env, err)
		return fallback
	}
	return float32(vf)
}

// Float64 gets a float64 value from an environment variable, or a fallback if not set or invalid.
//
// Example:
//     pi := env.Float64("PI", float64(3.14))
func Float64(env string, fallback float64) float64 {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	vf, err := strconv.ParseFloat(v, Float64BitSize)
	if err != nil {
		logError(env, err)
		return fallback
	}
	return vf
}

// Duration parses an environment variable as a time.Duration. It can parse duration
// strings like "300ms", "1.5h" or an integer number of nanoseconds.
//
// Example:
//     timeout := env.Duration("TIMEOUT", 5*time.Second)
func Duration(env string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	// First, try to parse as a duration string like "10s"
	if d, err := time.ParseDuration(v); err == nil {
		return d
	}
	// If that fails, try to parse as an integer (nanoseconds)
	vf, err := strconv.ParseInt(v, DurationBase, DurationBitSize)
	if err != nil {
		logError(env, err)
		return fallback
	}
	return time.Duration(vf)
}

// UnitDuration parses an environment variable as a number and multiplies it by the given unit.
// It can also parse full duration strings like "1h30m", in which case the unit is ignored.
// NOTE: The value of the env **AND** the fallback will be multiplied by the unit variable.
//
// Example:
//     # TIMEOUT=10 will result in 10 * time.Second
//     timeout := env.UnitDuration("TIMEOUT", 5, time.Second)
func UnitDuration(env string, fallback, unit time.Duration) time.Duration {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback * unit
	}
	// First, try to parse as a full duration string (e.g., "1h30m"). If so, the unit is ignored.
	if d, err := time.ParseDuration(v); err == nil {
		return d * unit
	}
	// If not, parse as a number and apply the unit.
	vf, err := strconv.ParseInt(v, UnitDurationBase, UnitDurationBitSize)
	if err != nil {
		logError(env, err)
		return fallback * unit
	}
	return time.Duration(vf) * unit
}

// List parses a delimited string from an environment variable into a []string.
// The delimiter is configured by the global `ListSeparator` variable.
//
// Example:
//     # With ENV TAGS="go,docker,linux"
//     tags := env.List("TAGS", []string{})
func List(env string, fallback []string) []string {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	parts := strings.Split(v, ListSeparator)
	list := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			list = append(list, part)
		}
	}
	if len(list) == 0 && len(parts) > 0 && ShowVerbose {
		OutLogger.Printf("List(%s) parsed %d parts but found 0 non-empty items", env, len(parts))
	}
	return list
}

// Map parses a delimited string from an env var into a map[string]string.
// Delimiters are configured by `MapSeparator` and `MapItemSeparator`.
//
// Example:
//     # With ENV DATA="key1=val1,key2=val2"
//     data := env.Map("DATA", env.ZeroMap)
func Map(env string, fallback map[string]string) map[string]string {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	parts := strings.Split(v, MapSeparator)
	result := make(map[string]string, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		pieces := strings.SplitN(part, MapItemSeparator, MapSplitN)
		if len(pieces) == 2 {
			key := strings.TrimSpace(pieces[0])
			value := strings.TrimSpace(pieces[1])
			if key != "" {
				result[key] = value
			}
		}
	}
	if len(result) == 0 && len(parts) > 0 && ShowVerbose {
		OutLogger.Printf("Map(%s) parsed %d parts but found 0 valid key-value pairs", env, len(parts))
	}
	return result
}

// ** DATA VALIDATION **

// ListContains checks if a list parsed from an env var contains a specific string (case-insensitive).
//
// Example:
//     if env.ListContains("FEATURES", env.ZeroList, "beta") {
//        // ...
//     }
func ListContains(env string, fallback []string, contains string) (ok bool) {
	L := List(env, fallback)
	for _, l := range L {
		if strings.EqualFold(l, contains) {
			ok = true
			break
		}
	}
	if !ok && ShowVerbose {
		OutLogger.Printf("ListContains(%s): '%s' not found in list. = %v", env, contains, ok)
	}
	return
}

// ListIsLength checks if a list parsed from an env var has a specific length.
//
// Example:
//     if !env.ListIsLength("HOSTS", env.ZeroList, 3) {
//        panic("env HOSTS must have exactly 3 items")
//     }
func ListIsLength(env string, fallback []string, wantLength int) (ok bool) {
	length := ListLength(env, fallback)
	ok = length == wantLength
	if !ok && ShowVerbose {
		OutLogger.Printf("ListIsLength(%s) got = %d; want %d", env, length, wantLength)
	}
	return
}

// ListLength returns the number of elements in a list parsed from an env var.
//
// Example:
//     if env.ListLength("HOSTS", env.ZeroList) > 3 {
//        fmt.Println("Warning: too many hosts specified.")
//     }
func ListLength(env string, fallback []string) int {
	return len(List(env, fallback))
}

// MapHasKey checks if a map parsed from an env var contains a specific key.
//
// Example:
//     if !env.MapHasKey("CONFIG", env.ZeroMap, "domain") {
//        panic("missing required key 'domain' from CONFIG env value")
//     }
func MapHasKey(env string, fallback map[string]string, key string) (ok bool) {
	M := Map(env, fallback)
	_, ok = M[key]
	if !ok && ShowVerbose {
		OutLogger.Printf("MapHasKey(%s)[%s] = %v", env, key, ok)
	}
	return
}

// MapHasKeys checks if a map parsed from an env var contains all specified keys.
//
// Example:
//     required := []string{"env", "author", "customer"}
//     if env.MapHasKeys("TAGS", env.ZeroMap, required...) {
//        fmt.Println("TAGS contain required keys!")
//     }
func MapHasKeys(env string, fallback map[string]string, keys ...string) (ok bool) {
	parsedMap := Map(env, fallback)
	foundCount := 0
	for _, key := range keys {
		if _, exists := parsedMap[key]; exists {
			foundCount++
		}
	}
	ok = foundCount == len(keys)
	if !ok && ShowVerbose {
		OutLogger.Printf("MapHasKeys(%s)[%s] = %v", env, strings.Join(keys, ","), ok)
	}
	return ok
}

// Int64LessThan checks if an int64 from an env var is less than a given value.
//
// Example:
//     if env.Int64LessThan("HITS", 0, 100) {
//        fmt.Println("Still under threshold.")
//     }
func Int64LessThan(env string, fallback, lessThan int64) (ok bool) {
	i := Int64(env, fallback)
	ok = i < lessThan
	if !ok && ShowVerbose {
		OutLogger.Printf("Int64LessThan(%s): %d < %d = %v", env, i, lessThan, ok)
	}
	return ok
}

// Int64GreaterThan checks if an int64 from an env var is greater than a given value.
//
// Example:
//     if env.Int64GreaterThan("HITS", 0, 1000) {
//        log.Fatal("Maximum hits received.")
//     }
func Int64GreaterThan(env string, fallback, greaterThan int64) (ok bool) {
	i := Int64(env, fallback)
	ok = i > greaterThan
	if !ok && ShowVerbose {
		OutLogger.Printf("Int64GreaterThan(%s): %d > %d = %v", env, i, greaterThan, ok)
	}
	return
}

// Int64InRange checks if an int64 from an env var is within a given range (inclusive).
//
// Example:
//     if env.Int64InRange("YEAR", 2020, 2000, 2025) {
//        fmt.Println("Valid year!")
//     }
func Int64InRange(env string, fallback, min, max int64) (ok bool) {
	i := Int64(env, fallback)
	ok = i >= min && i <= max
	if !ok && ShowVerbose {
		OutLogger.Printf("Int64InRange(%s): %d is not in range [%d, %d]", env, i, min, max)
	}
	return
}

// IntLessThan checks if an int from an env var is less than a given value.
//
// Example:
//     if env.IntLessThan("RETRIES", 0, 3) {
//        fmt.Println("Still have retries left.")
//     }
func IntLessThan(env string, fallback, lessThan int) (ok bool) {
	i := Int(env, fallback)
	ok = i < lessThan
	if !ok && ShowVerbose {
		OutLogger.Printf("IntLessThan(%s): %d < %d = %v", env, i, lessThan, ok)
	}
	return
}

// IntGreaterThan checks if an int from an env var is greater than a given value.
//
// Example:
//     if env.IntGreaterThan("CONNECTIONS", 0, 100) {
//        log.Fatal("Connection limit exceeded.")
//     }
func IntGreaterThan(env string, fallback, greaterThan int) bool {
	i := Int(env, fallback)
	ok := i > greaterThan
	if !ok && ShowVerbose {
		OutLogger.Printf("IntGreaterThan(%s): %d > %d = %v", env, i, greaterThan, ok)
	}
	return ok
}

// IntInRange checks if an int from an env var is within a given range (inclusive).
//
// Example:
//     if env.IntInRange("PORT", 8080, 1024, 49151) {
//        fmt.Println("Port is in the user range.")
//     }
func IntInRange(env string, fallback, min, max int) bool {
	i := Int(env, fallback)
	ok := i >= min && i <= max
	if !ok && ShowVerbose {
		OutLogger.Printf("IntInRange(%s): %d is not in range [%d, %d]", env, i, min, max)
	}
	return ok
}

// ** SYSTEM & OS **

// User returns the current user. If an error occurs, it may panic or return a
// default user, depending on global configuration (PanicNoUser).
//
// Example:
//     fmt.Printf("Your Home Directory: %s", env.User().HomeDir)
func User() (u *user.User) {
	var err error
	u, err = user.Current()

	if err != nil {
		willPanic := AllowPanic && PanicNoUser
		if !willPanic {
			logError("user.Current()", err)
		}
		if willPanic {
			panic(fmt.Sprintf("user.Current() failed: %v", err))
		}
		// Return a sensible default if not panicking
		return &user.User{
			HomeDir:  os.TempDir(),
			Uid:      "-1",
			Gid:      "-1",
			Username: "unknown",
			Name:     "Unknown",
		}
	}

	if ShowVerbose {
		var sb strings.Builder
		sb.WriteString("User() Details:\n")
		sb.WriteString(fmt.Sprintf("  Username: %s\n", u.Username))
		sb.WriteString(fmt.Sprintf("  Name:     %s\n", u.Name))
		sb.WriteString(fmt.Sprintf("  Uid:      %s\n", u.Uid))
		sb.WriteString(fmt.Sprintf("  Gid:      %s\n", u.Gid))
		sb.WriteString(fmt.Sprintf("  HomeDir:  %s\n", u.HomeDir))
		OutLogger.Print(sb.String())
	}
	return
}

// Set wraps os.Setenv, returning an error if it fails.
//
// Example:
//     err := env.Set("BACKUPS_DIR", "/tmp/backups")
//     if err != nil {
//        log.Fatal(err)
//     }
func Set(env string, value string) error {
	err := os.Setenv(env, value)
	if err != nil {
		logError(env, err)
		return err
	}
	if ShowVerbose {
		OutLogger.Printf("Set(%s) to value '%s'", env, value)
	}
	return nil
}

// Unset wraps os.Unsetenv, returning an error if it fails.
//
// Example:
//     err := env.Unset("OLD_VAR")
func Unset(env string) error {
	err := os.Unsetenv(env)
	if err != nil {
		logError(env, err)
		return err
	}
	if ShowVerbose {
		OutLogger.Printf("Unset(%s) successful", env)
	}
	return nil
}

// WasSet allows you to conditionally check if an ENV was defined for the runtime
//
// Example:
// 		u, _ := user.
// 		if !env.WasSet("BACKUPS_DIR", filepath.Join(env.User().HomeDir, "backups")) {
// 			// failed to set ENV "BACKUPS_DIR"
// 		}
func WasSet(env string, value string) bool {
	err := os.Setenv(env, value)
	if err != nil {
		logError(env, err)
		return false
	}
	// Verify it was set correctly
	v, ok := os.LookupEnv(env)
	if !ok || v != value {
		msg := fmt.Sprintf("WasSet(%s) failed verification after setting", env)
		logError(env, errors.New(msg))
		return false
	}
	if ShowVerbose {
		OutLogger.Printf("WasSet(%s) successful", env)
	}
	return true
}

// WasUnset allows you to conditionally check if an ENV was unset for the runtime
//
// Example:
// 		u, _ := user.
// 		if !env.WasUnset("BACKUPS_DIR") {
// 			// failed to unset ENV "BACKUPS_DIR"
// 		}
func WasUnset(env string) bool {
	err := os.Unsetenv(env)
	if err != nil {
		msg := fmt.Sprintf("WasUnset(%s) failed: %s", env, err.Error())
		logError(env, errors.New(msg))
		return false
	}
	if ShowVerbose {
		OutLogger.Printf("WasUnset(%s) successful", env)
	}
	return true
}
