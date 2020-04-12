package cron

import (
	"errors"
	"testing"
)

var (
	testLocales = map[LocaleType][]localeTestCase{
		Locale_en: en_TestCases(),
	}
)

func TestExpressionDescriptor_ToDescription(t *testing.T) {
	for loc, localeTestCases := range testLocales {
		t.Logf("=== Test '%s' locale, %d test case(s) ===", loc, len(localeTestCases))

		for i, tc := range localeTestCases {
			exprDesc, err := NewDescriptor(
				Verbose(tc.isVerbose),
				DayOfWeekStartsAtOne(tc.isDOWStartsAtOne),
				Use12HourTimeFormat(tc.is12HourTimeFormat),
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
					t.Errorf("%d. %s: expected %v error, got %v", i, tc.name, tc.outErr, err)
					return
				}
				if gotDesc != tc.outDesc && gotDesc != "" {
					t.Errorf("%d. %s: expected return empty string when error, got %v", i, tc.name, gotDesc)
					return
				}
				return
			}

			if gotDesc != tc.outDesc {
				t.Errorf("%d. %s: expected %v, got %v", i, tc.name, tc.outDesc, gotDesc)
				return
			}
		}
	}
}
