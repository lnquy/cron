package cron

func SetLogger(logger Logger) Option {
	return func(exprDesc *ExpressionDescriptor) {
		exprDesc.logger = logger
	}
}

func Verbose() Option {
	return func(exprDesc *ExpressionDescriptor) {
		exprDesc.isVerbose = true
	}
}

func DayOfWeekStartsAtZero() Option {
	return func(exprDesc *ExpressionDescriptor) {
		exprDesc.isDOWStartsAtZero = true
	}
}

func Use24HourTimeFormat() Option {
	return func(exprDesc *ExpressionDescriptor) {
		exprDesc.isDOWStartsAtZero = true
	}
}

func SetLocales(locales ...LocaleType) Option {
	return func(exprDesc *ExpressionDescriptor) {
		loaders, err := NewLocalLoaders(locales...)
		if err != nil {
			exprDesc.log("failed to init locale loaders: %s", err)
		}

		if exprDesc.locales == nil {
			exprDesc.locales = make(map[LocaleType]Locale)
		}
		for _, loader := range loaders {
			exprDesc.locales[loader.Type] = loader
		}
	}
}
