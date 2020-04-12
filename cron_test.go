package cron

import (
	"errors"
	"fmt"
	"testing"
)

var (
	testLocales = map[LocaleType][]localeTestCase{
		Locale_cs: cs_TestCases(),
		Locale_da: da_TestCases(),
		Locale_de: de_TestCases(),
		Locale_en: en_TestCases(),
		Locale_es: es_TestCases(),
		// Locale_fa: fa_TestCases(), // TODO: Need help from Farsi translator
		Locale_fi: fi_TestCases(),
		Locale_fr: fr_TestCases(),
		// Locale_he: he_TestCases(),  // TODO: Need help from Hebrew translator
		Locale_it: it_TestCases(),
		Locale_ja: ja_TestCases(),
		Locale_ko: ko_TestCases(),
		Locale_nb: nb_TestCases(),
		Locale_nl: nl_TestCases(),
	}
)

func TestExpressionDescriptor_ToDescription(t *testing.T) {
	for loc, localeTestCases := range testLocales {

		t.Run(fmt.Sprintf("%s_%d_case(s)", loc, len(localeTestCases)), func(t *testing.T) {
			for i, tc := range localeTestCases {
				exprDesc, err := NewDescriptor(
					Verbose(tc.isVerbose),
					DayOfWeekStartsAtOne(tc.isDOWStartsAtOne),
					Use24HourTimeFormat(tc.is24HourTimeFormat),
					SetLocales(loc),
				)
				if err != nil {
					t.Errorf("failed to create expression descriptor: %s", err)
					return
				}

				gotDesc, err := exprDesc.ToDescription(tc.inExpr, loc)
				if tc.outErr != nil {
					if err == nil {
						t.Errorf("%d. %s: expected error, got nil", i, tc.name)
						return
					}
					if !errors.Is(err, tc.outErr) {
						t.Errorf("%d. %s: expected '%v' error, got '%v'", i, tc.name, tc.outErr, err)
						return
					}
					if gotDesc != tc.outDesc && gotDesc != "" {
						t.Errorf("%d. %s: expected return empty string when error, got '%v'", i, tc.name, gotDesc)
						return
					}
					return
				}

				if gotDesc != tc.outDesc {
					t.Errorf("%d. %s: expected '%v', got '%v'", i, tc.name, tc.outDesc, gotDesc)
					return
				}
			}
		})
	}
}

var _desc string

func BenchmarkExpressionDescriptor_ToDescription(b *testing.B) {
	b.StopTimer()
	exprDesc, err := NewDescriptor(SetLocales(Locale_en))
	if err != nil {
		b.Errorf("failed to init expression descriptor: %s", err)
		return
	}
	expr := "0/5 1,5,10,15 */2 L JAN-OCT 1-5/2 2000-2050/10"
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_desc, err = exprDesc.ToDescription(expr, Locale_en)
		if err != nil {
			b.Fatalf("expected nil, got error: %s", err)
		}
	}
}
