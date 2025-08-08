# ENV

This go package is designed to give you easy access to interact with the ENV to request types back when working with
Linux operating systems (including macOS Silicon and Intel). 

## Installation

```shell
go get -u github.com/andreimerlescu/env
```

## Package Usage

```go
package main

import (
    "fmt"
    "strings"
    "time"
    
    "github.com/andreimerlescu/env"
)

func main() {
    hostname := env.String("HOSTNAME", "localhost")
    port := env.Int("PORT", 3306)
    user := env.String("USERNAME", "user")
    pass := env.String("PASSWORD", "")
    timeout := env.UnitDuration("TIMEOUT", time.Duration(3), time.Second)
    pi := env.Float64("PI", 3.14)
    pii := env.Float32("PI", 3.33)
    tags := env.List("TAGS", env.ZeroList)
    properties := env.Map("TAGS", env.ZeroMap)
    
    fmt.Println(
        hostname, 
        port, 
        user, 
        strings.Repeat("*", len(pass)), 
        pi, 
        pii, 
        timeout, 
        strings.Join(tags, ","), 
        properties)
}
```

## Function List

- **Can `panic()` with `export AM_GO_ENV_ALWAYS_ALLOW_PANIC=true`**
    - `func User() (u *user.User)`
    - `func MustExist(env string)`
- **Types**
    - `func Bool(env string, fallback bool) bool`
    - `func String(env string, fallback string) string`
    - `func Int(env string, fallback int) int`
    - `func Int64(env string, fallback int64) int64`
    - `func Float32(env string, fallback float32) float32`
    - `func Float64(env string, fallback float64) float64`
    - `func Duration(env string, fallback time.Duration) time.Duration`
    - `func UnitDuration(env string, fallback, unit time.Duration) time.Duration`
    - `func List(env string, fallback []string) []string`
    - `func Map(env string, fallback map[string]string) map[string]string`
- **Data Control**
    - `func Unset(env string)`
    - `func ListLength(env string, fallback []string) int`
    - `func Set(env string, value string) `
- **Conditionals**
    - `func Exists(env string) (ok bool)`
    - `func IsFalse(env string) (ok bool)`
    - `func IsTrue(env string) (ok bool)`
    - `func AreTrue(envs ...string) (ok bool)`
    - `func AreFalse(envs ...string) (ok bool)`
    - `func WasSetWasSet(env string, value string) (ok bool)`
    - `func WasUnset(env string) (ok bool)`
    - `func Int64LessThan(env string, fallback, lessThan int64) (ok bool)`
    - `func Int64GreaterThan(env string, fallback, greaterThan int64) (ok bool)`
    - `func Int64InRange(env string, fallback, min, max int64) (ok bool)`
    - `func IntLessThan(env string, fallback, lessThan int) (ok bool)`
    - `func IntGreaterThan(env string, fallback, greaterThan int) (ok bool)`
    - `func IntInRange(env string, fallback, min, max int) (ok bool)`
    - `func ListContains(env string, fallback []string, contains string) (ok bool)`
    - `func ListIsLength(env string, fallback []string, wantLength int) (ok bool)`
    - `func MapHasKeys(env string, fallback map[string]string, keys ...string) (ok bool)`
    - `func MapHasKey(env string, fallback map[string]string, key string) (ok bool)`

## Package Control Variables

This package can be manipulated using some control variables in your runtime when you're consuming the
package. For instance, you can define what makes a list from a string, ie. how you split the string.
The property `ListSeparator` defaults to `,` but allows you to change it. Additionally, there is 
`MapSeparator` (default `,`) and `MapItemSeparator` (default `=`) that allows you to define how a string
can become a `map[string]string` via `key1=value,key2=value` via `,` & `=`, but you can change it to anything.

```go
package main

import (
	"fmt"
    
    "github.com/andreimerlescu/env"
)

func main(){
    // Lists
    env.Set("TEST_LIST", "item1|item2|item3")
    env.ListSeparator = "|" 
    items := env.List("TEST_LIST", env.ZeroList)
    for i, item := range items {
		fmt.Printf("TEST_LIST#%d=%s\n", i, item)
    }
    
    // Maps
    env.Set("TEST_MAP", "key1~value|key2~value")
    env.MapSeparator = "|" 
    env.MapItemSeparator = "~"
	env.MapSplitN = 2 // default
    properties := env.Map("TEST_MAP", env.ZeroMap)
    for key, value := range properties {
        fmt.Printf("TEST_MAP[%s]=%s\n", key, value)
    }
}

```

Additionally, you can modify the **default** behavior in the way `strconv.ParseInt` and `strconv.ParseFloat`
customize the `base` and `bitSize` argument values when used in conjunction with `env`. 

| Variable              | Type  | Default | Usage                                     |
|-----------------------|-------|---------|-------------------------------------------|
| `Int64Base`           | `int` | `10`    | The `base` value of the number.           |
| `Int64BitSize`        | `int` | `64`    | The `bitSize` value of the number.        |
| `Float32BitSize`      | `int` | `32`    | The `bitSize` value of the number.        |
| `Float64BitSize`      | `int` | `64`    | The `bitSize` of the value of the number. |
| `DurationBase`        | `int` | `10`    | The `base` value of number.               |
| `DurationBitSize`     | `int` | `10`    | The `bitSize` value of the number.        |
| `UnitDurationBase`    | `int` | `10`    | The `base` value of the number.           |
| `UnitDurationBitSize` | `int` | `64`    | The `bitSize` value of the number.        |

Additionally, you can pre-define these values in your environment without manually changing the code of
any package that you submit. If your local workstation requires a difference here, for whatever reason,
this package allows you to assert the **Default** values from the Environment itself. 

| ENV                                | Expected Type | Variable                   | Type     |
|------------------------------------|---------------|----------------------------|----------|
| `AM_GO_ENV_MAP_SEPARATOR`          | `String`      | `env.MapSeparator`         | `String` |
| `AM_GO_ENV_MAP_ITEM_SEPARATOR`     | `String`      | `env.MapItemSeparator`     | `String` |
| `AM_GO_ENV_MAP_SPLIT_N`            | `Int`         | `env.MapSplitN`            | `Int`    |
| `AM_GO_ENV_LIST_SEPARATOR`         | `String`      | `env.ListSeparator`        | `String` |
| `AM_GO_ENV_UNIT_DURATION_BASE`     | `Int`         | `env.UnitDurationBase`     | `Int`    |
| `AM_GO_ENV_UNIT_DURATION_BIT_SIZE` | `Int`         | `env.UnitDurationBitSize`  | `Int`    |
| `AM_GO_ENV_DURATION_BASE`          | `Int`         | `env.DurationBase`         | `Int`    |
| `AM_GO_ENV_DURATION_BIT_SIZE`      | `Int`         | `env.DurationBitSize`      | `Int`    | 
| `AM_GO_ENV_FLOAT64_BIT_SIZE`       | `Int`         | `env.Float64BitSize`       | `Int`    |
| `AM_GO_ENV_FLOAT32_BIT_SIZE`       | `Int`         | `env.Float32BitSize`       | `Int`    |
| `AM_GO_ENV_INT64_BASE`             | `Int`         | `env.Int64Base`            | `Int`    |
| `AM_GO_ENV_INT64_BIT_SIZE`         | `Int`         | `env.Int64BitSize`         | `Int`    |
| `AM_GO_ENV_ALWAYS_ALLOW_PANIC`     | `String`      | `env.AllowPanic`           | `Bool`   |
| `AM_GO_ENV_ALWAYS_PRINT_ERRORS`    | `String`      | `env.PrintErrors`          | `Bool`   |
| `AM_GO_ENV_ALWAYS_USE_LOGGER`      | `String`      | `env.UseLogger`            | `Bool`   |
| `AM_GO_ENV_ENABLE_VERBOSE_LOGGING` | `String`      | `env.EnableVerboseLogging` | `Bool`   |
| `AM_GO_ENV_PANIC_NO_USER`          | `String`      | `env.PanicNoUser`          | `Bool`   |


> *NOTE* When **Type** `Bool` has an **Expected Type** of `String`, _true_ is `"true"` and _false_ is `"false"`.

## Types 

### `String`

```go
env.String("<ENV>", "<DEFAULT>")
```

### `Int`

```go
env.Int("<ENV>", <DEFAULT>)
```

### `Int64`

```go
env.Int64("<ENV>", <DEFAULT>)
```

### `Float32`

```go
env.Float32("<ENV>", <DEFAULT>)
```

### `Float64`

```go
env.Float64("<ENV>", <DEFAULT>)
```

### `Duration`

```go
env.Duration("<ENV>", <DEFAULT>)
```

### `UnitDuration`

```go
env.UnitDuration("<ENV>", <DEFAULT>, time.<UNIT>)
```

### `Map`

```go
env.Map("<ENV>", map[string]string{"<KEY>": "<VALUE>"})
```

## Environment Controls

You can use the following environment variables to define defaults used by the package across the runtime of any
go binary that uses this package. Useful if you're enforcing base 8 math on in64 operations.

| **Environment Variable**           | **Type** | **Default** | **Usage**                                                                  |
|------------------------------------|----------|-------------|----------------------------------------------------------------------------|
| `AM_GO_ENV_MAP_SEPARATOR`          | `String` | `,`         | How map items are separated.                                               |
| `AM_GO_ENV_MAP_ITEM_SEPARATOR`     | `String` | `=`         | How keys and values in map items are separated.                            |
| `AM_GO_ENV_MAP_SPLIT_N`            | `Int`    | `1`         | How many `AM_GO_ENV_MAP_ITEM_SEPARATOR` to `strings.SplitN` on.            |
| `AM_GO_ENV_LIST_SEPARATOR`         | `String` | `,`         | How list items are separated.                                              |
| `AM_GO_ENV_UNIT_DURATION_BASE`     | `Int`    | `10`        | How `strconv.ParseInt` defines `base` by default.                          |
| `AM_GO_ENV_UNIT_DURATION_BIT_SIZE` | `Int`    | `64`        | How `strconv.ParseInt` defines `bitSize` by default.                       |
| `AM_GO_ENV_DURATION_BASE`          | `Int`    | `10`        | How `strconv.ParseInt` defines `base` by default.                          |
| `AM_GO_ENV_DURATION_BIT_SIZE`      | `Int`    | `64`        | How `strconv.ParseInt` defines `bitSize` by default.                       |
| `AM_GO_ENV_FLOAT64_BIT_SIZE`       | `Int`    | `64`        | How `strconv.ParseFloat` defines `bitSize` by default.                     |
| `AM_GO_ENV_FLOAT32_BIT_SIZE`       | `Int`    | `32`        | How `strconv.ParseFloat` defines `bitSize` by default.                     |
| `AM_GO_ENV_INT64_BASE`             | `Int`    | `10`        | How `strconv.ParseInt` defines `base` by default.                          |
| `AM_GO_ENV_INT64_BIT_SIZE`         | `Int`    | `64`        | How `strconv.ParseInt` defines `bitSize` by default.                       | 
| `AM_GO_ENV_ALWAYS_ALLOW_PANIC`     | `Bool`   | `false`     | Permit the use of `panic()` within the package.                            |
| `AM_GO_ENV_ALWAYS_PRINT_ERRORS`    | `Bool`   | `false`     | Permit writing to STDERR to be returned.                                   |
| `AM_GO_ENV_ALWAYS_USE_LOGGER`      | `Bool`   | `false`     | Permit the STDOUT, STDERR logger to be used.                               |
| `AM_GO_ENV_ENABLE_VERBOSE_LOGGING` | `Bool`   | `false`     | Permit the STDOUT logger to report when something wasn't fully successful. |
| `AM_GO_ENV_PANIC_NO_USER`          | `Bool`   | `false`     | Use `panic()` over `os.Exit(1)` when `user.Current()` returns an error.    |

## License

This package is Open Source and is licensed with the Apache 2.0 License. 

## Why?

I really like clean code, and I find that when working with the environment in an application, it can get messy. This
package is my solution that. I wanted a way to quickly interact with the `os.Environ()` without having to do so much
repetitive work across each project. Thus, why this package was born. I just merged that work and forked it into its own
package so it can be used externally and enhanced to a better purpose. 

This package is called `env` and it is **NOT** made for Windows, but if you're on Windows, feel free to try it out. 
Enjoy using it!
