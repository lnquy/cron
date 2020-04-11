package cron

import (
	"fmt"
	"strings"
)

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
		exprDesc.parser = &cronParser{
			isDOWStartsAtZero: exprDesc.isDOWStartsAtZero,
		}
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
	var exprParts []string
	if exprParts, err = e.parser.Parse(expr); err != nil {
		return "", fmt.Errorf("failed to parse CRON expression: %w", err)
	}

	e.log("parsed: %v", strings.Join(exprParts, " | "))
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
