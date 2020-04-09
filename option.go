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
		// TODO
	}
}
