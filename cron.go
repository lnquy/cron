package cron

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	specialChars = []rune{'/', '-', ',', '*'}

	weekdaysNumberRegex = regexp.MustCompile(`(\d{1,2}w)|(w\d{1,2})`)
	lastDayOffsetRegex  = regexp.MustCompile(`L-(\d{1,2})`)
)

func containsAny(s string, matches []rune) bool {
	runes := []rune(s)
	for _, r := range runes {
		for _, c := range matches {
			if r == c {
				return true
			}
		}
	}

	return false
}

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

func NewDescriptor(options ...Option) (exprDesc *ExpressionDescriptor, err error) {
	exprDesc = &ExpressionDescriptor{}
	for _, option := range options {
		option(exprDesc)
	}

	// Init defaults
	if exprDesc.parser == nil {
		exprDesc.parser = &cronParser{
			isDOWStartsAtZero: exprDesc.isDOWStartsAtZero,
		}
	}

	// Always load EN locale so we can fallback to it
	if exprDesc.locales == nil {
		exprDesc.locales = make(map[LocaleType]Locale)
	}
	if _, ok := exprDesc.locales[Locale_en]; !ok {
		localeLoader, err := NewLocaleLoaders(Locale_en)
		if err != nil {
			return nil, fmt.Errorf("failed to init default locale EN: %w", err)
		}
		exprDesc.locales[Locale_en] = localeLoader[0]
	}

	return exprDesc, nil
}

func (e *ExpressionDescriptor) ToDescription(expr string, locale LocaleType) (desc string, err error) {
	var exprParts []string
	if exprParts, err = e.parser.Parse(expr); err != nil {
		return "", fmt.Errorf("failed to parse CRON expression: %w", err)
	}

	e.log("parsed: %v", strings.Join(exprParts, " | "))

	var timeSegment = e.getTimeOfDayDescription(exprParts, locale)
	var dayOfMonthDesc = e.getDayOfMonthDescription(exprParts, locale)
	var monthDesc = e.getMonthDescription(exprParts)
	var dayOfWeekDesc = e.getDayOfWeekDescription(exprParts)
	var yearDesc = e.getYearDescription(exprParts)

	desc = timeSegment + dayOfMonthDesc + dayOfWeekDesc + monthDesc + yearDesc
	// TODO: desc = transformVerbosity(desc, e.isVerbose)

	return desc, nil
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

func (e *ExpressionDescriptor) getTimeOfDayDescription(exprParts []string, loc LocaleType) string {
	second := exprParts[0]
	minute := exprParts[1]
	hour := exprParts[2]
	var desc string
	locale := e.getLocale(loc)

	if !containsAny(second, specialChars) && !containsAny(minute, specialChars) && !containsAny(hour, specialChars) {
		// specific time of day (i.e. 10:14:00)
		desc += locale.GetString(atSpace) + formatTime(hour, minute, second, locale, e.is24HourTimeFormat)
	} else if second == "" &&
		strings.Index(minute, "-") > -1 &&
		!(strings.Index(minute, ",") > -1) &&
		!(strings.Index(minute, "/") > -1) &&
		!containsAny(hour, specialChars) {
		// minute range in single hour (i.e. 0-10 11)
		idx := strings.Index(minute, "-")
		desc += fmt.Sprintf(locale.GetString(everyMinuteBetweenX0AndX1), minute[:idx], minute[idx:])
	} else if second == "" &&
		strings.Index(hour, ",") > 1 &&
		strings.Index(hour, "-") == -1 &&
		strings.Index(hour, "/") == -1 &&
		!containsAny(minute, specialChars) {
		// hours list with single minute (i.e. 30 6,14,16)
		hourParts := strings.Split(hour, ",")
		desc += locale.GetString(at)
		for i, p := range hourParts {
			desc += " "
			desc += formatTime(p, minute, "", locale, e.is24HourTimeFormat)
			if i < len(hourParts)-2 {
				desc += ", "
			}
			if i == len(hourParts)-2 {
				desc += locale.GetString(spaceAnd)
			}
		}
	} else {
		// default time description
		secondDesc := e.getSecondDescription(exprParts, locale)
		minuteDesc := e.getMinuteDescription(exprParts, locale)
		hourDesc := e.getHourDescription(exprParts, locale)

		desc += secondDesc
		if desc != "" && minuteDesc != "" {
			desc += ", "
		}
		desc += minuteDesc

		if desc != "" && hourDesc != "" {
			desc += ", "
		}
		desc += hourDesc
	}

	return desc
}

func (e *ExpressionDescriptor) getSecondDescription(exprParts []string, locale Locale) string {
	desc := getSegmentDescription(
		exprParts[0],
		locale.GetString(everySecond),
		func(s string) string {
			return s
		},
		func(s string) string {
			return fmt.Sprintf(locale.GetString(everyX0Seconds), s)
		},
		func(s string) string {
			return locale.GetString(secondsX0ThroughX1PastTheMinute)
		},
		func(s string) string {
			if s == "" {
				return ""
			}
			sInt, _ := strconv.Atoi(s)
			if sInt < 20 {
				return locale.GetString(atX0SecondsPastTheMinute)
			}
			if msg := locale.GetString(atX0SecondsPastTheMinuteGt20); msg != "" {
				return msg
			}
			return locale.GetString(atX0SecondsPastTheMinute)
		},
		locale,
	)

	return desc
}

func (e *ExpressionDescriptor) getDayOfMonthDescription(exprParts []string, loc LocaleType) string {
	desc := ""
	dom := exprParts[3]
	locale := e.getLocale(loc)

	switch dom {
	case "l":
		desc = locale.GetString(commaOnTheLastDayOfTheMonth)
	case "wl":
		fallthrough
	case "lw":
		desc = locale.GetString(commaOnTheLastWeekdayOfTheMonth)
	default:
		weekdaysNumberMatches := weekdaysNumberRegex.FindAllString(dom, -1)
		if len(weekdaysNumberMatches) > 0 {
			dayNumber, _ := strconv.Atoi(strings.Replace(weekdaysNumberMatches[0], "w", "", -1))
			dayStr := ""
			if dayNumber == 1 {
				dayStr = locale.GetString(firstWeekday)
			} else {
				dayStr = fmt.Sprintf(locale.GetString(weekdayNearestDayX0), strconv.Itoa(dayNumber))
			}
			desc = fmt.Sprintf(locale.GetString(commaOnTheX0OfTheMonth), dayStr)
			break
		}

		// Handle "last day offset" (i.e. L-5:  "5 days before the last day of the month")
		lastDayOffsetMatches := lastDayOffsetRegex.FindAllString(dom, -1)
		if len(lastDayOffsetMatches) > 0 {
			desc = fmt.Sprintf(locale.GetString(commaDaysBeforeTheLastDayOfTheMonth), lastDayOffsetMatches[0]) // TODO: cronstrue used 1
			break
		}
		// * dayOfMonth and dayOfWeek specified so use dayOfWeek verbiage instead
		if dom == "*" && exprParts[5] != "*" {
			return ""
		}
		desc = getSegmentDescription(
			dom,
			locale.GetString(commaEveryDay),
			func(s string) string {
				if s == "l" {
					return locale.GetString(lastDay)
				}
				return fmt.Sprintf(locale.GetString(dayX0), s) // TODO
			},
			func(s string) string {
				if s == "1" {
					return locale.GetString(commaEveryDay)
				}
				return locale.GetString(commaEveryX0Days)
			},
			func(s string) string {
				return locale.GetString(commaBetweenDayX0AndX1OfTheMonth)
			},
			func(s string) string {
				return locale.GetString(commaOnDayX0OfTheMonth)
			},
			locale,
		)
		break
	}

	return desc
}

func (e *ExpressionDescriptor) getMonthDescription(exprParts []string) string {
	return "*"
}

func (e *ExpressionDescriptor) getDayOfWeekDescription(exprParts []string) string {
	return "*"
}

func (e *ExpressionDescriptor) getYearDescription(exprParts []string) string {
	return "*"
}

func (e *ExpressionDescriptor) getLocale(loc LocaleType) Locale {
	v, ok := e.locales[loc]
	if !ok {
		return e.locales[Locale_en] // Fall back to default
	}
	return v
}

func formatTime(hour, minute, second string, locale Locale, isUse24HourTimeFormat bool) string {
	hourInt, _ := strconv.Atoi(hour)
	minuteInt, _ := strconv.Atoi(minute)
	secondInt, _ := strconv.Atoi(second)
	period := ""

	if !isUse24HourTimeFormat {
		period = getPeriod(hourInt, locale)
		if hourInt > 12 {
			hourInt -= 12
		}
		if hourInt == 0 {
			hourInt = 12
		}
	}

	hour = fmt.Sprintf("%02d", hourInt)
	minute = fmt.Sprintf("%02d", minuteInt)
	second = fmt.Sprintf("%02d", secondInt)
	return fmt.Sprintf("%s:%s:%s %s", hour, minute, second, period)
}

func getPeriod(hour int, locale Locale) string {
	if hour >= 12 {
		period := locale.GetString(pm)
		if period == "" {
			return "PM"
		}
		return period
	}

	period := locale.GetString(am)
	if period == "" {
		return "AM"
	}
	return period
}

type getStringFunc func(string) string

func getSegmentDescription(expr, allDesc string,
	getSingleItemDescription,
	getIntervalDescriptionFormat,
	getBetweenDescriptionFormat,
	getDescriptionFormat getStringFunc,
	locale Locale) string {
	desc := ""
	if expr == "" {
		desc = ""
	} else if expr == "*" {
		desc = allDesc
	} else if !containsAny(expr, []rune{'/', '-', ','}) {
		desc = fmt.Sprintf(getDescriptionFormat(expr), getSingleItemDescription(expr))
	} else if strings.Index(expr, "/") > -1 {
		segments := strings.Split(expr, "/")
		desc = fmt.Sprintf(getIntervalDescriptionFormat(segments[1]), segments[1])

		// interval contains 'between' piece (i.e. 2-59/3 )
		if strings.Index(segments[0], "-") > -1 {
			betweenDesc := generateBetweenSegmentDescription(segments[0], getBetweenDescriptionFormat, getSingleItemDescription)
			if strings.Index(betweenDesc, ", ") != 0 {
				desc += ", "
			}
			desc += betweenDesc
		} else if !containsAny(segments[0], []rune{'*', ','}) {
			rangeDesc := fmt.Sprintf(getDescriptionFormat(segments[0]), getSingleItemDescription(segments[0]))
			rangeDesc = strings.Replace(rangeDesc, ", ", "", 1)
			desc += fmt.Sprintf(locale.GetString(commaStartingX0), rangeDesc)
		}
	} else if strings.Index(expr, ",") > -1 {
		segments := strings.Split(expr, ",")
		contentDesc := ""
		for i, seg := range segments {
			if i > 0 && len(segments) > 2 {
				contentDesc += ","
				if i < len(segments)-1 {
					contentDesc += " "
				}
			}

			if i > 0 && len(segments) > 1 && (i == len(segments)-1 || len(segments) == 2) {
				contentDesc += locale.GetString(spaceAnd) + " "
			}

			getBetweenFmtFunc := func(s string) string { return locale.GetString(commaX0ThroughX1) }
			if strings.Index(seg, "-") > -1 {
				betweenDesc := generateBetweenSegmentDescription(
					seg,
					getBetweenFmtFunc,
					getSingleItemDescription,
				)
				betweenDesc = strings.Replace(betweenDesc, ", ", "", 1)
				contentDesc += betweenDesc
			} else {
				contentDesc += getSingleItemDescription(seg)
			}
		}

		desc += fmt.Sprintf(getDescriptionFormat(expr), contentDesc)
	} else if strings.Index(expr, "-") > -1 {
		desc = generateBetweenSegmentDescription(
			expr,
			getBetweenDescriptionFormat,
			getSingleItemDescription,
		)
	}

	return desc
}

func generateBetweenSegmentDescription(betweenDesc string, getBetweenDescriptionFormat, getSingleItemDescription getStringFunc) string {
	desc := ""
	betweenSegments := strings.Split(betweenDesc, "-")
	seg1 := getSingleItemDescription(betweenSegments[0])
	seg2 := getSingleItemDescription(betweenSegments[1])
	seg2 = strings.Replace(seg2, ":00", ":59", 1)
	desc += fmt.Sprintf(getBetweenDescriptionFormat(betweenDesc), seg1, seg2)
	return desc
}

func (e *ExpressionDescriptor) getMinuteDescription(exprParts []string, locale Locale) string {
	second := exprParts[0]
	hour := exprParts[2]

	desc := getSegmentDescription(
		exprParts[1],
		locale.GetString(everyMinute),
		func(s string) string {
			return s
		},
		func(s string) string {
			return fmt.Sprintf(locale.GetString(everyX0Minutes), s)
		},
		func(s string) string {
			return locale.GetString(minutesX0ThroughX1PastTheHour)
		},
		func(s string) string {
			if s == "0" && strings.Index(hour, "/") == -1 && second == "" {
				return locale.GetString(everyHour)
			}
			sInt, _ := strconv.Atoi(s)
			if sInt < 20 {
				return locale.GetString(atX0MinutesPastTheHour)
			}
			if msg := locale.GetString(atX0MinutesPastTheHourGt20); msg != "" {
				return msg
			}
			return locale.GetString(atX0MinutesPastTheHour)
		},
		locale)

	return desc
}

func (e *ExpressionDescriptor) getHourDescription(exprParts []string, locale Locale) string {
	desc := getSegmentDescription(
		exprParts[2],
		locale.GetString(everyHour),
		func(s string) string {
			return formatTime(s, "0", "", locale, e.is24HourTimeFormat)
		},
		func(s string) string {
			return fmt.Sprintf(locale.GetString(everyX0Hours), s)
		},
		func(s string) string {
			return locale.GetString(betweenX0AndX1)
		},
		func(s string) string {
			return locale.GetString(atX0)
		},
		locale,
	)

	return desc
}
