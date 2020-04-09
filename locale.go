package cron

const (
	LocaleAll LocaleType = "all"

	Locale_cs    LocaleType = "cs"
	Locale_da    LocaleType = "da"
	Locale_de    LocaleType = "de"
	Locale_en    LocaleType = "en"
	Locale_es    LocaleType = "es"
	Locale_fi    LocaleType = "fi"
	Locale_fr    LocaleType = "fr"
	Locale_he    LocaleType = "he"
	Locale_it    LocaleType = "it"
	Locale_ja    LocaleType = "ja"
	Locale_ko    LocaleType = "ko"
	Locale_nb    LocaleType = "nb"
	Locale_pt_BR LocaleType = "pt_BR"
	Locale_ro    LocaleType = "ro"
	Locale_ru    LocaleType = "ru"
	Locale_sk    LocaleType = "sk"
	Locale_sl    LocaleType = "sl"
	Locale_sv    LocaleType = "sv"
	Locale_sw    LocaleType = "sw"
	Locale_tr    LocaleType = "tr"
	Locale_uk    LocaleType = "uk"
	Locale_zh_CN LocaleType = "zh_CN"
	Locale_zh_TW LocaleType = "zh_TW"
)

var (
	allLocales = []LocaleType{
		Locale_cs,
		Locale_da,
		Locale_de,
		Locale_en,
		Locale_es,
		Locale_fi,
		Locale_fr,
		Locale_he,
		Locale_it,
		Locale_ja,
		Locale_ko,
		Locale_nb,
		Locale_pt_BR,
		Locale_ro,
		Locale_ru,
		Locale_sk,
		Locale_sl,
		Locale_sv,
		Locale_sw,
		Locale_tr,
		Locale_uk,
		Locale_zh_CN,
		Locale_zh_TW,
	}
)

type (
	LocaleType string

	Locale interface {
	}
)
