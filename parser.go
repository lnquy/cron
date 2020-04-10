package cron

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	EmptyExprError   = errors.New("expression is empty")
	InvalidExprError = errors.New("invalid expression")

	InvalidExprSecondError     = errors.New("invalid expression, second part")
	InvalidExprMinuteError     = errors.New("invalid expression, minute part")
	InvalidExprHourError       = errors.New("invalid expression, hour part")
	InvalidExprDayOfMonthError = errors.New("invalid expression, day of month part")
	InvalidExprMonthError      = errors.New("invalid expression, month part")
	InvalidExprDayOfWeekError  = errors.New("invalid expression, day of week part")
	InvalidExprYearError       = errors.New("invalid expression, year part")
)

var (
	yearRegex = regexp.MustCompile(`^\d{4}$`)
)

var (
	zeroRune  int32 = 48
	oneRune   int32 = 49
	twoRune   int32 = 50
	threeRune int32 = 51
	fourRune  int32 = 52
	fiveRune  int32 = 53
	sixRune   int32 = 54
	sevenRune int32 = 55
)

type (
	cronParser struct {
		isDOWStartsAtZero bool
	}

	Parser interface {
		Parse(expr string) (exprParts []string, err error)
	}
)

func (p *cronParser) Parse(expr string) (exprParts []string, err error) {
	exprParts, err = p.extractExprParts(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to extract expression parts: %w", err)
	}

	// TODO
	return nil, nil
}

func (p *cronParser) extractExprParts(expr string) (exprParts []string, err error) {
	if expr == "" {
		return nil, EmptyExprError
	}

	exprParts = make([]string, 7, 7)
	parts := strings.Fields(expr)

	switch {
	case len(parts) < 5:
		return nil, fmt.Errorf("expression has only %d part(s), at least 5 parts required: %w", len(parts), InvalidExprError)
	case len(parts) == 5:
		// Expression has 5 parts (standard POSIX cron)
		// => Prepend 1 and append 1 empty part at the beginning and the end of exprParts
		copy(exprParts[1:], append(parts, ""))
	case len(parts) == 6:
		// Has year (last part) or second (first part)
		if yearRegex.MatchString(parts[5]) {
			// Year provided => Prepend 1 empty part at the beginning for second
			copy(exprParts[1:], parts)
			break
		}
		// Second provided => Last parts (year) is empty
		copy(exprParts, parts)
	case len(parts) > 7:
		return nil, fmt.Errorf("expression has %d parts, at most 7 parts allowed: %w", len(parts), InvalidExprError)
	default: // Expression has 7 parts
		exprParts = parts
	}

	return exprParts, nil
}

func (p *cronParser) normalize(exprParts []string) (err error) {
	second := exprParts[0]
	minute := exprParts[1]
	hour := exprParts[2]
	dayOfMonth := exprParts[3]
	month := exprParts[4]
	dayOfWeek := exprParts[5]
	year := exprParts[6]

	// Convert ? to * for DOM and DOW
	dayOfMonth = strings.Replace(dayOfMonth, "?", "*", 1)
	dayOfWeek = strings.Replace(dayOfWeek, "?", "*", 1)
	// Convert ? to * for hour. ? isn't valid for hour position but we can work around it
	hour = strings.Replace(hour, "?", "*", 1)

	// Convert 0/, 1/ to */
	if strings.Index(second, "0/") == 0 {
		second = strings.Replace(second, "0/", "*/", 1)
	}
	if strings.Index(minute, "0/") == 0 {
		minute = strings.Replace(minute, "0/", "*/", 1)
	}
	if strings.Index(hour, "0/") == 0 {
		hour = strings.Replace(hour, "0/", "*/", 1)
	}
	if strings.Index(dayOfMonth, "1/") == 0 {
		dayOfMonth = strings.Replace(dayOfMonth, "1/", "*/", 1)
	}
	if strings.Index(month, "1/") == 0 {
		month = strings.Replace(month, "1/", "*/", 1)
	}
	if strings.Index(dayOfWeek, "1/") == 0 {
		dayOfWeek = strings.Replace(dayOfWeek, "1/", "*/", 1)
	}
	if strings.Index(year, "1/") == 0 {
		year = strings.Replace(year, "1/", "*/", 1)
	}

	// Adjust DOW based on isDOWStartsAtZero option
	// Normalized DOW: 0=Sunday/6=Saturday
	dowRunes := []rune(dayOfWeek)
	for i, c := range dowRunes {
		if c < zeroRune || c > sevenRune {
			continue
		}

		if p.isDOWStartsAtZero {
			if c != sevenRune {
				continue
			}
			c = oneRune // Accept 7 means Sunday too
		} else {
			if c == zeroRune {
				return fmt.Errorf("day of week starts at 1, must be from 1 to 7: %w", InvalidExprDayOfWeekError)
			}
			c -= 1 // Day of week start at 1 (Monday), so shift it 1
		}

		// Replace adjusted day of week
		dowRunes[i] = c
	}
	dayOfWeek = string(dowRunes)

	// Convert DOW 'L' to '6' (Saturday)
	if dayOfWeek == "L" {
		dayOfWeek = "6"
	}

	if strings.Index(dayOfMonth, "W") > -1 &&
		(strings.Index(dayOfMonth, ",") > -1 || strings.Index(dayOfMonth, "-") > -1) {
		return fmt.Errorf("the 'W' character can be specified only when the day-of-month is a single day, not a range or list of days: %w", InvalidExprDayOfMonthError)
	}

	// TODO: WIP

	return nil
}

func isNumberChar(c int32) bool {
	return c >= 48 && c <= 57
}
