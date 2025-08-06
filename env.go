package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// AmGoEnvMapSeparator allows you to define the default separator used for map items. Default is ","
	AmGoEnvMapSeparator string = "AM_GO_ENV_MAP_SEPARATOR"

	// AmGoEnvMapItemSeparator allows you to define how key and values in map items are split. Default is "="
	AmGoEnvMapItemSeparator string = "AM_GO_ENV_MAP_ITEM_SEPARATOR"

	// AmGoEnvMapSplitN allows you to define the strings.SplitN n value. Default is "1"
	AmGoEnvMapSplitN string = "AM_GO_ENV_MAP_SPLIT_N"

	// AmGoEnvListSeparator allows you to define how list items are separated. Default is ","
	AmGoEnvListSeparator string = "AM_GO_ENV_LIST_SEPARATOR"

	// AmGoEnvUnitDurationBase allows you to define strconv.ParseInt base. Default is "10"
	AmGoEnvUnitDurationBase string = "AM_GO_ENV_UNIT_DURATION_BASE"

	// AmGoEnvUnitDurationBitSize allows you to define strconv.ParseInt bitSize. Default is "64"
	AmGoEnvUnitDurationBitSize string = "AM_GO_ENV_UNIT_DURATION_BIT_SIZE"

	// AmGoEnvDurationBase allows you to define strconv.ParseInt base. Default is "10"
	AmGoEnvDurationBase string = "AM_GO_ENV_DURATION_BASE"

	// AmGoEnvDurationBitSize allows you to define strconv.ParseInt bitSize. Default is "64"
	AmGoEnvDurationBitSize string = "AM_GO_ENV_DURATION_BIT_SIZE"

	// AmGoEnvFloat64BitSize allows you to define strconv.ParseFloat bitSize. Default is "64"
	AmGoEnvFloat64BitSize string = "AM_GO_ENV_FLOAT64_BIT_SIZE"

	// AmGoEnvFloat32BitSize allows you to define strconv.ParseFloat bitSize. Default is "32"
	AmGoEnvFloat32BitSize string = "AM_GO_ENV_FLOAT32_BIT_SIZE"

	// AmGoEnvInt64Base allows you to define strconv.ParseInt base. Default is "10"
	AmGoEnvInt64Base string = "AM_GO_ENV_INT64_BASE"

	// AmGoEnvInt64BitSize allows you to define strconv.ParseInt bitSize. Default is "64"
	AmGoEnvInt64BitSize string = "AM_GO_ENV_INT64_BIT_SIZE"
)

// Exists allows you to use an ENV such as "HOSTNAME=localhost" and check if it exists
//
// Example:
// 		ok := env.Exists("HOSTNAME")
func Exists(env string) bool {
	_, ok := os.LookupEnv(env)
	if !ok {
		return false
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
		return fallback
	}
	return vb
}

// IsFalse verifies that Bool is false
//
// Example:
// 		if env.IsFalse("USE_JSON") {
// 			fmt.Println("Not using JSON.")
// 		} else {
//			// intend to use json
// 		}
func IsFalse(env string) bool {
	return Bool(env, false) == false
}

// IsTrue verifies that Bool is true
//
// Example:
// 		if env.IsTrue("USE_JSON") {
//			// intend to use json
// 		} else {
// 			fmt.Println("Not using JSON.")
// 		}
func IsTrue(env string) bool {
	return Bool(env, false) == true
}

// AreTrue verifies a list of ENV variables and runs IsTrue to verify all are true
//
// Example:
// 		if env.AreTrue("ALWAYS_PRINT", "USE_JSON") {
// 			// intend to print to STDOUT in JSON format
// 		}
func AreTrue(envs ...string) bool {
	o := make(map[string]bool, len(envs))
	for _, env := range envs {
		o[env] = IsTrue(env)
	}
	O := 0
	for _, g := range o {
		if g {
			O++
		}
	}
	return O == len(envs)
}

// AreFalse verifies a list of ENV variables and runs IsFalse to verify all are false
//
// Example:
// 		if env.AreFalse("NEVER_SAVE", "NEVER_OVERWRITE") {
// 			// act accordingly
// 		}
func AreFalse(envs ...string) bool {
	o := make(map[string]bool, len(envs))
	for _, env := range envs {
		o[env] = IsFalse(env)
	}
	O := 0
	for _, g := range o {
		if g {
			O++
		}
	}
	return O == len(envs)
}

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
		return fallback
	}
	return vint
}

// IntLessThan uses Int to check if the lessThan is less than
//
// Example:
// 		if env.IntLessThan("HITS", 0, 33) {
//			fmt.Printf("Congrats, hits are at: %d\n", env.Int("HITS"))
//		}
func IntLessThan(env string, fallback int, lessThan int) bool {
	return Int(env, fallback) < lessThan
}

// IntGreaterThan uses Int to allow you to check if the greaterThan is greater than
//
// Example:
// 		if env.IntGreaterThan("HITS", 0, 1000) {
// 			log.Fatal("Maximum hits received.")
// 		}
func IntGreaterThan(env string, fallback int, greaterThan int) bool {
	return Int(env, fallback) > greaterThan
}

// IntInRange uses Int to allow you to check for a min and max value as < and > exclusive
//
// Example:
// 		if env.IntInRange("HITS", 0, 1,999) {
// 			fmt.Println("Valid season!")
// 		}
func IntInRange(env string, fallback int, min int, max int) bool {
	i := Int(env, fallback)
	return i < min || i > max
}

var (
	// Int64Base allows you to define the default behavior of strconv.ParseInt inside Int64
	Int64Base int = Int(AmGoEnvInt64Base, 10)

	// Int64BitSize allows you to define the default behavior of strconv.ParseInt inside Int64
	Int64BitSize int = Int(AmGoEnvInt64BitSize, 64)
)

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
		return fallback
	}
	return int64(vint)
}

// Float32BitSize allows you to define the default behavior of the strconv.ParseFloat inside Float32
var Float32BitSize int = Int(AmGoEnvFloat32BitSize, 32)

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
		return fallback
	}
	return float32(vf)
}

// Float64BitSize allows you to define the default behavior of the strconv.ParseFloat inside Float64
var Float64BitSize int = Int(AmGoEnvFloat64BitSize, 64)

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
		return fallback
	}
	return vf
}

var (
	// DurationBase allows you to define the default behavior of the strconv.ParseInt inside Duration
	DurationBase int = Int(AmGoEnvDurationBase, 10)

	// DurationBitSize allows you to define the default behavior of the strconv.ParseInt inside Duration
	DurationBitSize int = Int(AmGoEnvDurationBitSize, 64)
)

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
	vf, err := strconv.ParseInt(v, DurationBase, DurationBitSize)
	if err != nil {
		return fallback
	}
	return time.Duration(vf)
}

var (
	// UnitDurationBase allows you to define the default behavior of the strconv.ParseInt inside UnitDuration
	UnitDurationBase int = Int(AmGoEnvUnitDurationBase, 10)
	// UnitDurationBitSize allows you to define the default behavior of the strconv.ParseInt inside UnitDuration
	UnitDurationBitSize int = Int(AmGoEnvUnitDurationBitSize, 64)
)

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
	vf, err := strconv.ParseInt(v, UnitDurationBase, UnitDurationBitSize)
	if err != nil {
		return fallback * unit
	}
	return time.Duration(vf) * unit
}

// ListSeparator defines how a list is split, CSV uses "," and TSV uses "|", you can customize this to anything
var ListSeparator string = String(AmGoEnvListSeparator, ",")

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
	return list
}

var (
	// MapSeparator defines how key=value pairs in a map are split, csv uses ",", tsv uses "|" inside MapStringString
	MapSeparator string = String(AmGoEnvMapSeparator, ",")

	// MapItemSeparator defines how keys and values are split in the string inside MapStringString
	MapItemSeparator = String(AmGoEnvMapItemSeparator, "=")

	// MapSplitN defines the default behavior on strings.SplitN inside MapStringString
	MapSplitN int = Int(AmGoEnvMapSplitN, 1)
)

// MapStringString allows you to store a comma separated list of key=value pairs in an ENV with a fallback
//
// Example:
// 		tags := env.MapStringString("TAGS",
//			map[string]string{
//				"Environment", env.String("ENVIRONMENT", "dev",
//			)})
func MapStringString(env string, fallback map[string]string) map[string]string {
	v, ok := os.LookupEnv(env)
	if !ok {
		return fallback
	}
	parts := strings.Split(v, MapSeparator)
	result := make(map[string]string, len(parts))
	for _, part := range parts {
		pieces := strings.SplitN(strings.TrimSpace(part), MapItemSeparator, MapSplitN)
		if len(pieces) != 2 {
			result[pieces[0]] = pieces[1]
		}
	}
	return result
}
