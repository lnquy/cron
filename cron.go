package cron

type (
	ExpressionDescriptor struct {
		isVerbose          bool
		isDOWStartsAtZero  bool
		is24HourTimeFormat bool

		logger  Logger
		parser  Parser
		locales map[LocaleType]Locale
	}

	Logger interface {
		Printf(format string, v ...interface{})
	}

	Option func(exprDesc *ExpressionDescriptor)
)

func NewDescriptor(options ...Option) *ExpressionDescriptor {
	var exprDesc ExpressionDescriptor
	for _, option := range options {
		option(&exprDesc)
	}

	// Init defaults
	if exprDesc.parser == nil {
		// TODO
	}
	if len(exprDesc.locales) == 0 {
		// TODO
	}
	// Always load EN locale so we can fallback to it
	if _, ok := exprDesc.locales[Locale_en]; !ok {
		// TODO
	}

	return &exprDesc
}

func (e *ExpressionDescriptor) ToDescription(expr string, locale LocaleType) (desc string, err error) {
	// TODO
	return "", nil
}

func (e *ExpressionDescriptor) log(format string, v ...interface{}) {
	if e.logger == nil {
		return
	}
	e.logger.Printf(format, v...)
}

func (e *ExpressionDescriptor) verbose(format string, v ...interface{}) {
	if !e.isVerbose || e.logger == nil {
		return
	}
	e.logger.Printf(format, v...)
}
