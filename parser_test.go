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
			outErr:   InvalidExprError,
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
			inExpr:   "3 * * 1/5 1/5 1/2 1/10",
			outExprs: []string{"3", "*", "*", "*/5", "*/5", "*/2", "*/10"},
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
			inExpr:   "5 * * 3/5 2/5 2/2 2000/10",
			outExprs: []string{"5", "*", "*", "3/5", "2-12/5", "2-6/2", "2000-2099/10"},
			outErr:   nil,
		},

		// validate
		{
			name:     "should error if second is invalid 2",
			inExpr:   "60 * * * * * *",
			outExprs: nil,
			outErr:   InvalidExprSecondError,
		}, {
			name:     "should error if second is invalid 3",
			inExpr:   "9223372036854775808 * * * * * *",
			outExprs: nil,
			outErr:   InvalidExprSecondError,
		}, {
			name:     "should error if minute is invalid 2",
			inExpr:   "* 60 * * * * *",
			outExprs: nil,
			outErr:   InvalidExprMinuteError,
		}, {
			name:     "should error if hour is invalid 2",
			inExpr:   "* * 24 * * * *",
			outExprs: nil,
			outErr:   InvalidExprHourError,
		}, {
			name:     "should error if DOM is invalid 2",
			inExpr:   "* * * 32 * * *",
			outExprs: nil,
			outErr:   InvalidExprDayOfMonthError,
		}, {
			name:     "should error if DOM contains invalid characters",
			inExpr:   "* * LX * *",
			outExprs: nil,
			outErr:   InvalidExprDayOfMonthError,
		}, {
			name:     "should error if month is invalid 2",
			inExpr:   "* * * * 13 * *",
			outExprs: nil,
			outErr:   InvalidExprMonthError,
		}, {
			name:     "should error if DOW is invalid 2",
			inExpr:   "* * * * * 8 *",
			outExprs: nil,
			outErr:   InvalidExprDayOfWeekError,
		}, {
			name:     "should error if DOW contains invalid characters",
			inExpr:   "* * * * MO",
			outExprs: nil,
			outErr:   InvalidExprDayOfWeekError,
		}, {
			name:     "should error if year is invalid 1",
			inExpr:   "* * * * * * 0",
			outExprs: nil,
			outErr:   InvalidExprYearError,
		}, {
			name:     "should error if year is invalid 2",
			inExpr:   "* * * * * * 2100",
			outExprs: nil,
			outErr:   InvalidExprYearError,
		},
		// Cannot test due to no reliable way to detect negative number in cron expression
		// {
		// 	name:     "should error if second is invalid 1",
		// 	inExpr:   "-1 * * * * * *",
		// 	outExprs: nil,
		// 	outErr:   InvalidExprSecondError,
		// }, {
		// 	name:     "should error if minute is invalid 1",
		// 	inExpr:   "* -1 * * * * *",
		// 	outExprs: nil,
		// 	outErr:   InvalidExprMinuteError,
		// },{
		// 	name:     "should error if hour is invalid 1",
		// 	inExpr:   "* * -1 * * * *",
		// 	outExprs: nil,
		// 	outErr:   InvalidExprHourError,
		// },{
		// 	name:     "should error if DOM is invalid 1",
		// 	inExpr:   "* * * -1 * * *",
		// 	outExprs: nil,
		// 	outErr:   InvalidExprDayOfMonthError,
		// },{
		// 	name:     "should error if month is invalid 1",
		// 	inExpr:   "* * * * 0 * *",
		// 	outExprs: nil,
		// 	outErr:   InvalidExprMonthError,
		// }, {
		// 	name:     "should error if DOW is invalid 1",
		// 	inExpr:   "* * * * * -1 *",
		// 	outExprs: nil,
		// 	outErr:   InvalidExprDayOfWeekError,
		// },
	}

	parser := cronParser{}

	for i, tc := range tcs {
		parser.isDOWStartsAtOne = tc.inTestDOWStartsAtOne

		parsed, err := parser.Parse(tc.inExpr)
		if tc.outErr != nil {
			if err == nil {
				t.Errorf("%d. %s: expected error, got nil", i, tc.name)
				return
			}
			if !errors.Is(err, tc.outErr) {
				t.Errorf("%d. %s: expected '%v' error, got '%v'", i, tc.name, tc.outErr, err)
				return
			}
			if !reflect.DeepEqual(parsed, tc.outExprs) {
				t.Errorf("%d. %s: expected return nil when error, got '%v'", i, tc.name, parsed)
				return
			}
			continue
		}

		if !reflect.DeepEqual(parsed, tc.outExprs) {
			t.Errorf("%d. %s: expected '%v', got '%v'", i, tc.name, tc.outExprs, parsed)
			return
		}
	}
}

var _parsed []string

func BenchmarkCronParser_Parse(b *testing.B) {
	b.StopTimer()
	parser := &cronParser{}
	expr := "0/5 1,5,10,15 */2 L JAN-OCT 1-5/2 2000-2050/10"
	var err error
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_parsed, err = parser.Parse(expr)
		if err != nil {
			b.Fatalf("expected nil, got error: %s", err)
		}
	}
}
