package cron

func SetLogger(logger Logger) Option {
	return func(exprDesc *ExpressionDescriptor) {
		exprDesc.logger = logger
	}
}

func Verbose(v bool) Option {
	return func(exprDesc *ExpressionDescriptor) {
		exprDesc.isVerbose = v
	}
}

func DayOfWeekStartsAtOne(v bool) Option {
	return func(exprDesc *ExpressionDescriptor) {
		exprDesc.isDOWStartsAtOne = v
	}
}

func Use24HourTimeFormat(v bool) Option {
	return func(exprDesc *ExpressionDescriptor) {
		exprDesc.is24HourTimeFormat = v
	}
}

func SetLocales(locales ...LocaleType) Option {
	return func(exprDesc *ExpressionDescriptor) {
		loaders, err := NewLocaleLoaders(locales...)
		if err != nil {
			exprDesc.log("failed to init locale loaders: %s", err)
		}

		if exprDesc.locales == nil {
			exprDesc.locales = make(map[LocaleType]Locale)
		}
		for _, loader := range loaders {
			exprDesc.locales[loader.GetLocaleType()] = loader
		}
	}
}
