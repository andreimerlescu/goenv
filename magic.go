package env

// Magic will read the environment of the runtime where this package is imported and will attempt to reassign
// the variables used throughout the application like UseLogger, AllowPanic, etc. where these values are
// retrieved using their related Bool(), Int(), String() functions accordingly. Naming pattern of these are
// AM_GO_ENV_<UPPER_UNDERLINE_CONVERSION> so UseLogger would be AM_GO_ENV_USE_LOGGER
func Magic() {
	MapSplitN = Int(AmGoEnvMapSplitN, 2)
	Int64Base = Int(AmGoEnvInt64Base, 10)
	PanicNoUser = Bool(AmGoEnvPanicNoUser, AllowPanic)
	Int64BitSize = Int(AmGoEnvInt64BitSize, 64)
	DurationBase = Int(AmGoEnvDurationBase, 10)
	PrintErrors = Bool(AmGoEnvAlwaysPrintErrors, UseLogger)
	Float64BitSize = Int(AmGoEnvFloat64BitSize, 64)
	Float32BitSize = Int(AmGoEnvFloat32BitSize, 32)
	MapSeparator = String(AmGoEnvMapSeparator, ",")
	UseLogger = Bool(AmGoEnvAlwaysUseLogger, false)
	AllowPanic = Bool(AmGoEnvAlwaysAllowPanic, true)
	DurationBitSize = Int(AmGoEnvDurationBitSize, 64)
	ListSeparator = String(AmGoEnvListSeparator, ",")
	UnitDurationBase = Int(AmGoEnvUnitDurationBase, 10)
	MapItemSeparator = String(AmGoEnvMapItemSeparator, "=")
	UnitDurationBitSize = Int(AmGoEnvUnitDurationBitSize, 64)
	EnableVerboseLogging = Bool(AmGoEnvEnableVerboseLogging, false)
}
