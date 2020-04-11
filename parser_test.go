package cron

import (
	"errors"
	"reflect"
	"testing"
)

func TestCronParser_Parse(t *testing.T) {
	type testCase struct {
		name                 string
		inTestDOWStartsAtOne bool
		inExpr               string
		outExprs             []string
		outErr               error
	}

	tcs := []testCase{
		// extractExprParts
		{
			name:     "should failed on empty",
			inExpr:   "   ",
			outExprs: nil,
			outErr:   EmptyExprError,
		}, {
			name:     "should parse cron with multiple spaces between parts",
			inExpr:   "30  2  *    *  *",
			outExprs: []string{"", "30", "2", "*", "*", "*", ""},
			outErr:   nil,
		}, {
			name:     "should parse 5 part cron",
			inExpr:   "* * * * *",
			outExprs: []string{"", "*", "*", "*", "*", "*", ""},
			outErr:   nil,
		}, {
			name:     "should parse 6 part cron with year",
			inExpr:   "* * * * * 2020",
			outExprs: []string{"", "*", "*", "*", "*", "*", "2020"},
			outErr:   nil,
		}, {
			name:     "should parse 6 part cron with second",
			inExpr:   "5 * * * * *",
			outExprs: []string{"5", "*", "*", "*", "*", "*", ""},
			outErr:   nil,
		}, {
			name:     "should parse 7 part cron",
			inExpr:   "  5 * * *    * *  2020  ",
			outExprs: []string{"5", "*", "*", "*", "*", "*", "2020"},
			outErr:   nil,
		}, {
			name:     "should error if more than 7 part",
			inExpr:   "5 * * * * *  2020 *",
			outExprs: nil,
			outErr:   InvalidExprError,
		}, {
			name:     "should error if expression is not a cron schedule",
			inExpr:   "sdlksCRAPdlkskl- dds",
			outExprs: nil,
			outErr:   InvalidExprError,
		},

		// normalize
		{
			name:     "should parse cron in uppercase",
			inExpr:   "* * L JAN-NOV MON-FRI",
			outExprs: []string{"", "*", "*", "l", "1-11", "1-5", ""},
			outErr:   nil,
		}, {
			name:     "should parse cron in lowercase",
			inExpr:   "* * l jan-NOV MON-fri",
			outExprs: []string{"", "*", "*", "l", "1-11", "1-5", ""},
			outErr:   nil,
		}, {
			name:     "should convert 0 second to empty",
			inExpr:   "0 * * * * * *",
			outExprs: []string{"", "*", "*", "*", "*", "*", "*"},
			outErr:   nil,
		}, {
			name:     "should convert ? to *",
			inExpr:   "* ? ? * ?",
			outExprs: []string{"", "*", "*", "*", "*", "*", ""},
			outErr:   nil,
		}, {
			name:     "should convert 0/ to */",
			inExpr:   "0/5 0/5 0/5 * * *",
			outExprs: []string{"*/5", "*/5", "*/5", "*", "*", "*", ""},
			outErr:   nil,
		}, {
			name:     "should convert 1/ to */",
			inExpr:   "3 * * 1/5 1/5 1/2 1/1000",
			outExprs: []string{"3", "*", "*", "*/5", "*/5", "*/2", "*/1000"},
			outErr:   nil,
		}, {
			name:     "should use 7 as Sunday when DOWStartsAtZero is true",
			inExpr:   "* * * * 1-7",
			outExprs: []string{"", "*", "*", "*", "*", "1-0", ""},
			outErr:   nil,
		}, {
			name:                 "should failed when DOWStartsAtZero is false and have 0",
			inExpr:               "* * * * 0-7",
			inTestDOWStartsAtOne: true,
			outExprs:             nil,
			outErr:               InvalidExprDayOfWeekError,
		}, {
			name:                 "should parse when DOWStartsAtZero is false",
			inExpr:               "* * * * 1-7",
			inTestDOWStartsAtOne: true,
			outExprs:             []string{"", "*", "*", "*", "*", "0-6", ""},
			outErr:               nil,
		}, {
			name:     "should convert L DOW to 6",
			inExpr:   "* * * * L",
			outExprs: []string{"", "*", "*", "*", "*", "6", ""},
			outErr:   nil,
		}, {
			name:     "should failed if DOM 'W' comes with list of dates 1",
			inExpr:   "* * 1-3W * *",
			outExprs: nil,
			outErr:   InvalidExprDayOfMonthError,
		}, {
			name:     "should failed if DOM 'W' comes with list of dates 2",
			inExpr:   "* * 1,2,5W * *",
			outExprs: nil,
			outErr:   InvalidExprDayOfMonthError,
		}, {
			name:     "should convert hour to range 1",
			inExpr:   "0-20/3 9 * * *",
			outExprs: []string{"", "0-20/3", "9-9", "*", "*", "*", ""},
			outErr:   nil,
		}, {
			name:     "should convert hour to range 2",
			inExpr:   "*/5 3 * * *",
			outExprs: []string{"", "*/5", "3-3", "*", "*", "*", ""},
			outErr:   nil,
		}, {
			name:     "should normalize */1 to *",
			inExpr:   "*/1 3 * * *",
			outExprs: []string{"", "*", "3-3", "*", "*", "*", ""},
			outErr:   nil,
		}, {
			name:     "should normalize to range",
			inExpr:   "5 * * 3/5 2/5 2/2 2000/1000",
			outExprs: []string{"5", "*", "*", "3/5", "2-12/5", "2-6/2", "2000-9999/1000"},
			outErr:   nil,
		},

		// validate
		{
			name:     "should error if DOW part is not valid",
			inExpr:   "* * * * MO",
			outExprs: nil,
			outErr:   InvalidExprDayOfWeekError,
		}, {
			name:     "should error if DOM part is not valid",
			inExpr:   "* * LX * *",
			outExprs: nil,
			outErr:   InvalidExprDayOfMonthError,
		},
	}

	parser := cronParser{isDOWStartsAtZero: true}

	for i, tc := range tcs {
		testFunc := func() {
			parsed, err := parser.Parse(tc.inExpr)
			if tc.outErr != nil {
				if err == nil {
					t.Errorf("%d. %s: expected error, got nil", i, tc.name)
					return
				}
				if !errors.Is(err, tc.outErr) {
					t.Errorf("%d. %s: expected %v error, got %v", i, tc.name, tc.outErr, err)
					return
				}
				if !reflect.DeepEqual(parsed, tc.outExprs) {
					t.Errorf("%d. %s: expected return nil when error, got %v", i, tc.name, parsed)
					return
				}
				return
			}

			if !reflect.DeepEqual(parsed, tc.outExprs) {
				t.Errorf("%d. %s: expected %v, got %v", i, tc.name, tc.outExprs, parsed)
				return
			}
		}

		// Test the DOWStartsAtZero is false
		if tc.inTestDOWStartsAtOne {
			parser.isDOWStartsAtZero = false
			testFunc()
			parser.isDOWStartsAtZero = true // Set it back to default
			continue
		}

		testFunc() // Otherwise test with default isDOWStartsAtZero
	}
}