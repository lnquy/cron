package cron

const (
	Locale_all LocaleType = "all"
	Locale_en  LocaleType = "en"
)

var (
	allLocales = []LocaleType{Locale_en}
)

type (
	LocaleType string

	Locale     interface {
	}
)
