package env

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

func init() {
	UseLogger = Bool(AmGoEnvAlwaysUseLogger, false)
	AllowPanic = Bool(AmGoEnvAlwaysAllowPanic, true)
	PanicNoUser = Bool(AmGoEnvPanicNoUser, AllowPanic)
	PrintErrors = Bool(AmGoEnvAlwaysPrintErrors, UseLogger)
	EnableVerboseLogging = Bool(AmGoEnvEnableVerboseLogging, false)

	OutLogger = log.New(os.Stdout, "INFO: ${env.git} ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrLogger = log.New(os.Stderr, " ERR: ${env.git} ", log.Ldate|log.Ltime|log.Lshortfile)

	ShowVerbose = UseLogger && EnableVerboseLogging

	Float64BitSize = Int(AmGoEnvFloat64BitSize, 64)
	Float32BitSize = Int(AmGoEnvFloat32BitSize, 32)
	ListSeparator = String(AmGoEnvListSeparator, ",")
	DurationBase = Int(AmGoEnvDurationBase, 10)
	DurationBitSize = Int(AmGoEnvDurationBitSize, 64)
	MapSeparator = String(AmGoEnvMapSeparator, ",")
	MapItemSeparator = String(AmGoEnvMapItemSeparator, "=")
	MapSplitN = Int(AmGoEnvMapSplitN, 2)
	UnitDurationBase = Int(AmGoEnvUnitDurationBase, 10)
	UnitDurationBitSize = Int(AmGoEnvUnitDurationBitSize, 64)
	Int64Base = Int(AmGoEnvInt64Base, 10)
	Int64BitSize = Int(AmGoEnvInt64BitSize, 64)

	ZeroList = make([]string, 0)
	ZeroMap = make(map[string]string)
}

// ** EXISTENCE **

// Exists allows you to use an ENV such as "HOSTNAME=localhost" and check if it exists
//
// Example:
// 		ok := env.Exists("HOSTNAME")
func Exists(env string) (ok bool) {
	_, ok = os.LookupEnv(env)
	if !ok && ShowVerbose {
		OutLogger.Printf("Exits(%s) = %v", env, ok)
	}
	return
}

// MustExist throws a panic if the env is not found
//
// Example:
// 		// os.Environ[1] => "HOSTNAME="
// 		func init() {
//			env.MustExists("HOSTNAME")
// 		}
// 		func main() {
// 			hostname := env.String("HOSTNAME", "")
// 			if hostname == "" {
// 				// hostname is set, but its an empty string
// 			}
// 		}
func MustExist(env string) {
	_, ok := os.LookupEnv(env)
	if !ok {
		if AllowPanic {
			panic(fmt.Sprintf("MustExist() failed because %s was not found", env))
		}
		if PrintErrors {
			msg := fmt.Sprintf("MustExist() failed because environment variable %s was not found", env)
			if UseLogger {
				ErrLogger.Fatal(msg)
			} else {
				_, _ = fmt.Fprintln(os.Stderr, msg)
				os.Exit(1)
			}
		}
	}
	if EnableVerboseLogging {
		OutLogger.Printf("MustExist() has found %s!", env)
	}
}

// ** TRUTHY **

// IsFalse verifies that Bool is false
//
// Example:
// 		if env.IsFalse("USE_JSON") {
// 			fmt.Println("Not using JSON.")
// 		} else {
//			// intend to use json
// 		}
func IsFalse(env string) (ok bool) {
	ok = !Bool(env, false) == false
	if !ok && ShowVerbose {
		OutLogger.Printf("IsFalse(%s) = %v", env, ok)
	}
	return
}

// IsTrue verifies that Bool is true
//
// Example:
// 		if env.IsTrue("USE_JSON") {
//			// intend to use json
// 		} else {
// 			fmt.Println("Not using JSON.")
// 		}
func IsTrue(env string) (ok bool) {
	ok = Bool(env, false) == true
	if !ok && ShowVerbose {
		OutLogger.Printf("IsTrue(%s) = %v", env, ok)
	}
	return
}

// AreTrue verifies a list of ENV variables and runs IsTrue to verify all are true
//
// Example:
// 		if env.AreTrue("ALWAYS_PRINT", "USE_JSON") {
// 			// intend to print to STDOUT in JSON format
// 		}
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

// AreFalse verifies a list of ENV variables and runs IsFalse to verify all are false
//
// Example:
// 		if env.AreFalse("NEVER_SAVE", "NEVER_OVERWRITE") {
// 			// act accordingly
// 		}
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

// Bool allows you to use an ENV to store a string "true" and "false" to parse as a bool
//
// Example:
// 		useJson := env.Bool("USE_JSON", false)
// 		# when no USE_JSON set, useJson is false
func Bool(env string, fallback bool) bool {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	vb, err := strconv.ParseBool(v)
	if err != nil {
		if PrintErrors {
			if UseLogger {
				ErrLogger.Printf("Bool(%s) strconv.ParseBool(%s) threw err: %v", env, v, err)
			} else {
				fmt.Printf("Bool(%s) strconv.ParseBool(%s) threw err: %v", env, v, err)
			}
		}
		return fallback
	}
	return vb
}

// ** TYPES **

// String retries a string value from an ENV name with a fallback value when not defined
//
// Example:
// 		hostname := env.String("HOSTNAME", "localhost")
// 		if strings.EqualFold(hostname, "localhost") {
// 			// localhost
// 		} else {
// 			// not localhost
// 		}
func String(env string, fallback string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	return v
}

// Int retrieves an int value from an ENV name with a fallback int value when not defined
//
// Example:
// 		port := env.Int("PORT", 3306)
// 		if port < 1 || port > 65534 {
// 			panic("port out of range")
// 		}
func Int(env string, fallback int) int {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	vint, err := strconv.Atoi(v)
	if err != nil {
		if PrintErrors {
			if UseLogger {
				ErrLogger.Printf("Int(%s) failed due to: %s", env, err.Error())
			} else {
				fmt.Printf("Int(%s) failed due to: %s", env, err.Error())
			}
		}
		return fallback
	}
	return vint
}

// Int64 retrieves an int64 value from an ENV name with a fallback value if not defined or has an error
//
// Example
// 		ns := env.Int64("NANOSECONDS", int64(1_000_000_000))
func Int64(env string, fallback int64) int64 {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	vint, err := strconv.ParseInt(v, Int64Base, Int64BitSize)
	if err != nil {
		if PrintErrors {
			if UseLogger {
				OutLogger.Printf("Int64(%s) failed due to: %s", env, err.Error())
			} else {
				fmt.Printf("Int64(%s) failed due to: %s", env, err.Error())
			}
		}
		return fallback
	}
	return int64(vint)
}

// Float32 retrieves a float32 value from an ENV name with a fallback value if undefined or has an error
//
// Example:
// 		pi := env.Float32("PI", float32(3.14))
func Float32(env string, fallback float32) float32 {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	vf, err := strconv.ParseFloat(v, Float32BitSize)
	if err != nil {
		if PrintErrors {
			if UseLogger {
				ErrLogger.Printf("Float32(%s) failed due to %s", env, err.Error())
			} else {
				fmt.Printf("Float32(%s) failed due to %s", env, err.Error())
			}
		}
		return fallback
	}
	return float32(vf)
}

// Float64 retrieves a float64 value from an ENV name with a fallback value if undefined or has an error
//
// Example:
// 		pi := env.Float64("PI", float64(3.14))
func Float64(env string, fallback float64) float64 {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	vf, err := strconv.ParseFloat(v, Float64BitSize)
	if err != nil {
		if PrintErrors {
			if UseLogger {
				ErrLogger.Printf("Float64(%s) failed due to: %s", env, err.Error())
			} else {
				fmt.Printf("Float64(%s) failed due to: %s", env, err.Error())
			}
		}
		return fallback
	}
	return vf
}

// Duration interacts with time.Duration to allow an ENV var to be set to a value of such or return a fallback
//
// Example:
// 		ns := env.Duration("NANOSECONDS", 3*time.Nanosecond)
// 		if ns > time.Hour {
// 			panic("cant exceed an hour")
// 		}
func Duration(env string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err == nil {
		return d
	}
	vf, err := strconv.ParseInt(v, DurationBase, DurationBitSize)
	if err != nil {
		if PrintErrors {
			if UseLogger {
				ErrLogger.Printf("Duration(%s) failed due to %s", env, err.Error())
			} else {
				fmt.Printf("Duration(%s) failed due to %s", env, err.Error())
			}
		}
		return fallback
	}
	return time.Duration(vf)
}

// UnitDuration takes Duration to a scalar of unit time.Duration
//
// Example:
// 		seconds := env.UnitDuration("SECONDS", 3*time.Second, time.Second)
// 		# setting ENV of SECONDS=3 will return time.Duration(3*time.Second) result
// 		if seconds > time.Hour {
// 			panic("cannot exceed an hour")
// 		}
func UnitDuration(env string, fallback, unit time.Duration) time.Duration {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback * unit
	}
	d, err := time.ParseDuration(v)
	if err == nil {
		return d * unit
	}
	vf, err := strconv.ParseInt(v, UnitDurationBase, UnitDurationBitSize)
	if err != nil {
		if PrintErrors {
			if UseLogger {
				ErrLogger.Printf("UnitDuration(%s) failed due to %s", env, err.Error())
			} else {
				fmt.Printf("UnitDuration(%s) failed due to %s", env, err.Error())
			}
		}
		return fallback * unit
	}
	return time.Duration(vf) * unit
}

// List allows you to get a []string from a comma separated list stored in an ENV
//
// Example:
// 		aliases := env.List("ALIASES", []string{})
// 		# store as ALIASES as "dog,canine,beast,creature"
// 		if len(aliases) > 0 {
// 			for _, a := range aliases {
//				fmt.Println(a)
// 			}
// 		}
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
	ll, lp := len(list), len(parts)
	if ll == 0 && lp > 0 && ShowVerbose {
		OutLogger.Printf("List(%s) has %d items, needs %d items", env, ll, lp)
	}
	return list
}

// Map allows you to store a comma separated list of key=value pairs in an ENV with a fallback
//
// Example:
// 		tags := env.Map("TAGS",
//			map[string]string{
//				"Environment", env.String("ENVIRONMENT", "dev",
//			)})
func Map(env string, fallback map[string]string) map[string]string {
	defer func() {
		if r := recover(); r != nil {
			if AllowPanic {
				panic(r)
			}
			if PrintErrors {
				if UseLogger {
					ErrLogger.Printf("Map(%s) failed due to %v", env, r)
				} else {
					_, _ = fmt.Fprintf(os.Stderr, "Map(%s) failed due to %v", env, r)
				}
			}
		}
	}()
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	parts := strings.Split(v, MapSeparator)
	result := make(map[string]string, len(parts))
	for _, part := range parts {
		pieces := strings.SplitN(strings.TrimSpace(part), MapItemSeparator, MapSplitN)
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

// ** CONTROL DATA **

// ListContains uses strings.EqualFold on the elements in the list to check contains
//
// Example:
// 		if env.ListContains("HOSTS", []string{}, "localhost:27017") {
// 			fmt.Println("MongoDB")
// 		}
func ListContains(env string, fallback []string, contains string) (ok bool) {
	L := List(env, fallback)
	defer func() {
		if !ok && ShowVerbose {
			OutLogger.Printf("ListContains(%s) '%s' = %v", env, contains, ok)
		}
	}()
	for _, l := range L {
		if strings.EqualFold(l, contains) {
			ok = true
			return
		}
	}
	ok = false
	if !ok && ShowVerbose {
		OutLogger.Printf("ListContains(%s) '%s' = %v", env, contains, ok)
	}
	return
}

// ListIsLength uses <= on the length against the length of the list by env name
//
// Example:
// 		if !env.ListIsLength("HOSTS", []string{}, 3) {
// 			panic("env HOSTS must have 3 items inside it")
// 		}
func ListIsLength(env string, fallback []string, wantLength int) (ok bool) {
	i := List(env, fallback)
	length := len(i)
	ok = length == wantLength
	if !ok && ShowVerbose {
		OutLogger.Printf("ListIsLength(%s) got = %d; want %d", env, length, wantLength)
	}
	return
}

// ListLength returns the number of elements in the List
//
// Example:
// 		if env.ListLength("HOSTS", []string{}) > 3 {
// 			fmt.Println("Thank you for following directions and listening to the guide.")
// 		}
func ListLength(env string, fallback []string) int {
	L := List(env, fallback)
	return len(L)
}

// MapHasKey allows you to interact with an ENV string storing a map to verify if it has a key or not
//
// Example:
// 		if !env.MapHasKey("TAGS", map[string]string{}, "domain") {
//			panic("missing required key 'domain' from TAGS env value")
// 		}
func MapHasKey(env string, fallback map[string]string, key string) (ok bool) {
	M := Map(env, fallback)
	_, ok = M[key]
	if !ok && ShowVerbose {
		OutLogger.Printf("MapHasKey(%s)[%s] = %v", env, key, ok)
	}
	return
}

// MapHasKeys allows you to interact with an ENV string storing a map to verify if it contains multiple keys
//
// Example:
// 		if env.MapHasKeys("TAGS", map[string]string{}, []string{"environment", "author", "customer"}) {
// 			fmt.Println("TAGS contain required keys!")
// 		}
func MapHasKeys(env string, fallback map[string]string, keys ...string) (ok bool) {
	J := Map(env, fallback)
	R := 0
	for _, F := range keys {
		if _, K := J[F]; K {
			R++
		}
	}
	ok = R == len(keys)
	if !ok && ShowVerbose {
		OutLogger.Printf("MapHasKeys(%s)[%s] = %v", env, strings.Join(keys, ","), ok)
	}
	return ok
}

// Int64LessThan uses Int64 to check if the lessThan is less than
//
// Example:
// 		if env.IntLessThan("HITS", 0, 33) {
//			fmt.Printf("Congrats, hits are at: %d\n", env.Int("HITS"))
//		}
func Int64LessThan(env string, fallback, lessThan int64) (ok bool) {
	i := Int64(env, fallback)
	ok = i < lessThan
	if !ok && ShowVerbose {
		OutLogger.Printf("Int64(%s) %d < %d = %v", env, i, lessThan, ok)
	}
	return ok
}

// Int64GreaterThan uses Int64 to allow you to check if the greaterThan is greater than
//
// Example:
// 		if env.IntGreaterThan("HITS", 0, 1000) {
// 			log.Fatal("Maximum hits received.")
// 		}
func Int64GreaterThan(env string, fallback, greaterThan int64) (ok bool) {
	i := Int64(env, fallback)
	ok = i > greaterThan
	if !ok && ShowVerbose {
		OutLogger.Printf("Int64(%s) %d > %d = %v", env, i, greaterThan, ok)
	}
	return
}

// Int64InRange uses Int64 to allow you to check for a min and max value as < and > exclusive
//
// Example:
// 		if env.IntInRange("HITS", 0, 1,999) {
// 			fmt.Println("Valid season!")
// 		}
func Int64InRange(env string, fallback, min, max int64) (ok bool) {
	i := Int64(env, fallback)
	ok = i >= min && i <= max
	if !ok && ShowVerbose {
		OutLogger.Printf("Int64InRange(%s): %d is not in range [%d, %d]", env, i, min, max)
	}
	return
}

// IntLessThan uses Int to check if the lessThan is less than
//
// Example:
// 		if env.IntLessThan("HITS", 0, 33) {
//			fmt.Printf("Congrats, hits are at: %d\n", env.Int("HITS"))
//		}
func IntLessThan(env string, fallback, lessThan int) (ok bool) {
	i := Int(env, fallback)
	ok = i < lessThan
	if !ok && ShowVerbose {
		OutLogger.Printf("Int(%s) %d < %d = %v", env, i, lessThan, ok)
	}
	return
}

// IntGreaterThan uses Int to allow you to check if the greaterThan is greater than
//
// Example:
// 		if env.IntGreaterThan("HITS", 0, 1000) {
// 			log.Fatal("Maximum hits received.")
// 		}
func IntGreaterThan(env string, fallback, greaterThan int) bool {
	i := Int(env, fallback)
	ok := i > greaterThan
	if !ok && ShowVerbose {
		OutLogger.Printf("Int(%d) > %d == %v", i, greaterThan, ok)
	}
	return ok
}

// IntInRange uses Int to allow you to check for a min and max value as < and > exclusive
//
// Example:
// 		if env.IntInRange("HITS", 0, 1,999) {
// 			fmt.Println("Valid season!")
// 		}
func IntInRange(env string, fallback, min, max int) bool {
	i := Int(env, fallback)
	ok := i >= min && i <= max
	if !ok && ShowVerbose {
		OutLogger.Printf("IntInRange(%s): %d is not in range [%d, %d]", env, i, min, max)
	}
	return ok
}

// ** SYSTEM CALLS **

// User returns a *user.User where if an error occurs, an empty user.User{} is returned pointing to os.TempDir
//
// Example:
// 		fmt.Printf("Your Home Directory: %s", env.User().HomeDir)
func User() (u *user.User) {
	var err error
	u, err = user.Current()
	defer func() {
		if ShowVerbose {
			var sb strings.Builder
			sb.WriteString("User()\n")
			sb.WriteString("  Username: " + u.Username + "\n")
			sb.WriteString("  Name: " + u.Name + "\n")
			sb.WriteString("  Uid: " + u.Uid + "\n")
			sb.WriteString("  Gid: " + u.Gid + "\n")
			sb.WriteString("  HomeDir: " + u.HomeDir + "\n")
			OutLogger.Println(sb.String())
		}
	}()
	if err != nil {
		willPanic := AllowPanic && PanicNoUser
		if PrintErrors && !willPanic {
			ErrLogger.Printf("user.Current() aka env.User() failed due to err: %s", err.Error())
		}
		if willPanic {
			panic(fmt.Errorf("user.Current() aka env.User() failed due to err: %s", err.Error()))
		}
		u = &user.User{
			HomeDir:  os.TempDir(),
			Uid:      "-1",
			Gid:      "-1",
			Username: "unknown",
			Name:     "Unknown",
		}
	}
	return
}

// ** CHANGES **

// Set wraps os.Setenv() and uses the PrintErrors, EnableVerboseLogging, UseLogger, ErrLogger, and OutLogger
//
// Example:
// 		env.Set("BACKUPS_DIR", filepath.Join(env.User().HomeDir, "backups"))
func Set(env string, value string) {
	err := os.Setenv(env, value)
	if err != nil {
		if PrintErrors {
			msg := fmt.Sprintf("Set(%s) failed: %s", env, err.Error())
			if UseLogger {
				ErrLogger.Println(msg)
			} else {
				_, _ = fmt.Fprintln(os.Stderr, msg)
			}
		}
		return
	}
	if ShowVerbose {
		OutLogger.Printf("Set(%s) to value '%s'", env, value)
	}
}

// Unset wraps os.Unsetenv and uses the PrintErrors EnableVerboseLogging UseLogger ErrLogger and OutLogger
//
// Example:
// 		env.Unset("BACKUPS_DIR")
func Unset(env string) {
	err := os.Unsetenv(env)
	if err != nil {
		if PrintErrors {
			msg := fmt.Sprintf("Unset(%s) failed: %s", env, err.Error())
			if UseLogger {
				ErrLogger.Println(msg)
			} else {
				_, _ = fmt.Fprintln(os.Stderr, msg)
			}
		}
		return
	}
	if ShowVerbose {
		OutLogger.Printf("Unset(%s) successful", env)
	}
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
		if PrintErrors {
			msg := fmt.Sprintf("WasUnset(%s) failed: %s", env, err.Error())
			if UseLogger {
				ErrLogger.Println(msg)
			} else {
				_, _ = fmt.Fprintln(os.Stderr, msg)
			}
		}
		return false
	}
	if ShowVerbose {
		OutLogger.Printf("WasUnset(%s) successful", env)
	}
	return true
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
		if PrintErrors {
			msg := fmt.Sprintf("WasSet(%s) failed: %s", env, err.Error())
			if UseLogger {
				ErrLogger.Println(msg)
			} else {
				_, _ = fmt.Fprintln(os.Stderr, msg)
			}
		}
		return false
	}
	// Verify it was set correctly
	v, ok := os.LookupEnv(env)
	if !ok || v != value {
		if PrintErrors {
			msg := fmt.Sprintf("WasSet(%s) failed verification after setting", env)
			if UseLogger {
				ErrLogger.Println(msg)
			} else {
				_, _ = fmt.Fprintln(os.Stderr, msg)
			}
		}
		return false
	}
	if ShowVerbose {
		OutLogger.Printf("WasSet(%s) successful", env)
	}
	return true
}
