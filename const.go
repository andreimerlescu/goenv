package env

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

	// AmGoEnvAlwaysAllowPanic allows you to control whether panic() is permitted in MustExist
	AmGoEnvAlwaysAllowPanic string = "AM_GO_ENV_ALWAYS_ALLOW_PANIC"

	// AmGoEnvAlwaysPrintErrors allows you to control whether OutLogger and ErrLogger is leveraged
	AmGoEnvAlwaysPrintErrors string = "AM_GO_ENV_ALWAYS_PRINT_ERRORS"

	// AmGoEnvAlwaysUseLogger allows you to always use the OutLogger and ErrLogger
	AmGoEnvAlwaysUseLogger string = "AM_GO_ENV_ALWAYS_USE_LOGGER"

	// AmGoEnvEnableVerboseLogging allows you to enable verbose logging that writs to OutLogger
	AmGoEnvEnableVerboseLogging string = "AM_GO_ENV_ENABLE_VERBOSE_LOGGING"

	// AmGoEnvPanicNoUser allows you to panic() when user.Current() returns an error
	AmGoEnvPanicNoUser string = "AM_GO_ENV_PANIC_NO_USER"
)
