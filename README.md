# ENV

This go package is designed to give you easy access to interact with the ENV to request types back.

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
    tags := env.MapStringString("TAGS", map[string]string{"test": "test"})
    
    fmt.Println(hostname, port, user, strings.Repeat("*", len(pass)), pi, pii, timeout, tags)
}
```

## Function List

- **Types**
    - `func Bool(env string, fallback bool) bool`
    - `func String(env string, fallback string) string`
    - `func Int(env string, fallback int) int`
    - `func IntLessThan(env string, fallback int, lessThan int) bool`
    - `func IntGreaterThan(env string, fallback int, greaterThan int) bool`
    - `func IntInRange(env string, fallback int, min int, max int) bool`
    - `func Int64(env string, fallback int64) int64`
    - `func Float32(env string, fallback float32) float32`
    - `func Float64(env string, fallback float64) float64`
    - `func Duration(env string, fallback time.Duration) time.Duration`
    - `func UnitDuration(env string, fallback, unit time.Duration) time.Duration`
    - `func List(env string, fallback []string) []string`
    - `func MapStringString(env string, fallback map[string]string) map[string]string`
- **Conditionals**
    - `func Exists(env string) bool`
    - `func IsFalse(env string) bool`
    - `func IsTrue(env string) bool`
    - `func AreTrue(envs ...string) bool`
    - `func AreFalse(envs ...string) bool`

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

### `MapStringString`

```go
env.MapStringString("<ENV>", map[string]string{"<KEY>": "<VALUE>"})
```

## Environment Controls

You can use the following environment variables to define defaults used by the package across the runtime of any
go binary that uses this package. Useful if you're enforcing base 8 math on in64 operations.

| **Environment Variable**           | **Type** | **Default** | **Usage**                                                       |
|------------------------------------|----------|-------------|-----------------------------------------------------------------|
| `AM_GO_ENV_MAP_SEPARATOR`          | `String` | `,`         | How map items are separated.                                    |
| `AM_GO_ENV_MAP_ITEM_SEPARATOR`     | `String` | `=`         | How keys and values in map items are separated.                 |
| `AM_GO_ENV_MAP_SPLIT_N`            | `Int`    | `1`         | How many `AM_GO_ENV_MAP_ITEM_SEPARATOR` to `strings.SplitN` on. |
| `AM_GO_ENV_LIST_SEPARATOR`         | `String` | `,`         | How list items are separated.                                   |
| `AM_GO_ENV_UNIT_DURATION_BASE`     | `Int`    | `10`        | How `strconv.ParseInt` defines `base` by default.               |
| `AM_GO_ENV_UNIT_DURATION_BIT_SIZE` | `Int`    | `64`        | How `strconv.ParseInt` defines `bitSize` by default.            |
| `AM_GO_ENV_DURATION_BASE`          | `Int`    | `10`        | How `strconv.ParseInt` defines `base` by default.               |
| `AM_GO_ENV_DURATION_BIT_SIZE`      | `Int`    | `64`        | How `strconv.ParseInt` defaines `bitSize` by default.           |
| `AM_GO_ENV_FLOAT64_BIT_SIZE`       | `Int`    | `64`        | How `stconv.ParseFloat` defines `bitSize` by default.           |
| `AM_GO_ENV_FLOAT32_BIT_SIZE`       | `Int`    | `32`        | How `strconv.ParseFloat` defines `bitSize` by default.          |
| `AM_GO_ENV_INT64_BASE`             | `Int`    | `10`        | How `strconv.ParseInt` defines `base` by default.               |
| `AM_GO_ENV_INT64_BIT_SIZE`         | `Int`    | `64`        | How `strconv.ParseInt` defines `bitSize` by default.            | 

