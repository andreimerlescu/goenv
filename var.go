package env

import (
	"log"
	"time"
)

// RUNTIME CONTROL
var (
	UseMagic bool = true

	// UseLogger controls whether or not the env package will use OutLogger or ErrLogger to report problems
	UseLogger bool

	// AllowPanic determines whether or not panic() should be called on MustExist
	AllowPanic bool

	// PanicNoUser determines on User() whether to throw a panic if the user.Current() returns an error
	PanicNoUser bool

	// PrintErrors determines whether or not
	PrintErrors bool

	// EnableVerboseLogging allows you to use the OutLogger to receive helpful debugging information
	EnableVerboseLogging bool

	// OutLogger defaults to STDOUT but requires UseLogger or AmGoEnvAlwaysUseLogger to be defined as 'true'
	OutLogger *log.Logger

	// ErrLogger defaults to STDERR but requires UseLogger or AmGoEnvAlwaysUseLogger to be defined as 'true'
	ErrLogger *log.Logger

	// ShowVerbose is enabled when UseLogger and EnableVerboseLogging are enabled
	ShowVerbose bool
)

// DATA CONTROL

var (
	// ListSeparator defines how a list is split, CSV uses "," and TSV uses "|", you can customize this to anything
	ListSeparator string

	// MapSeparator defines how key=value pairs in a map are split, csv uses ",", tsv uses "|" inside Map
	MapSeparator string

	// MapItemSeparator defines how keys and values are split in the string inside Map
	MapItemSeparator string

	// MapSplitN defines the default behavior on strings.SplitN inside Map
	MapSplitN int

	// Float64BitSize allows you to define the default behavior of the strconv.ParseFloat inside Float64
	Float64BitSize int

	// Float32BitSize allows you to define the default behavior of the strconv.ParseFloat inside Float32
	Float32BitSize int

	// DurationBase allows you to define the default behavior of the strconv.ParseInt inside Duration
	DurationBase int

	// DurationBitSize allows you to define the default behavior of the strconv.ParseInt inside Duration
	DurationBitSize int
	// UnitDurationBase allows you to define the default behavior of the strconv.ParseInt inside UnitDuration
	UnitDurationBase int

	// UnitDurationBitSize allows you to define the default behavior of the strconv.ParseInt inside UnitDuration
	UnitDurationBitSize int

	// Int64Base allows you to define the default behavior of strconv.ParseInt inside Int64
	Int64Base int

	// Int64BitSize allows you to define the default behavior of strconv.ParseInt inside Int64
	Int64BitSize int
)

// EMPTY TYPES

var (
	ZeroInt          int           = 0
	ZeroInt64        int64         = 0
	ZeroFloat64      float64       = 0
	ZeroFloat32      float32       = 0
	ZeroBool         bool          = false
	ZeroString       string        = ""
	ZeroDuration     time.Duration = 0
	ZeroUnitDuration time.Duration = 0
	ZeroList         []string
	ZeroMap          map[string]string
)
