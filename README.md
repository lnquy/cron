# cron
<p align="left">
  <a href="https://godoc.org/github.com/lnquy/cron" title="GoDoc Reference" rel="nofollow"><img src="https://img.shields.io/badge/go-documentation-blue.svg?style=flat" alt="GoDoc Reference"></a>
  <a href="https://github.com/github.com/lnquy/cron/releases/tag/v0.0.1" title="0.0.1 Release" rel="nofollow"><img src="https://img.shields.io/badge/version-0.0.1-blue.svg?style=flat" alt="0.0.1 release"></a>
  <a href="https://goreportcard.com/report/github.com/lnquy/cron"><img src="https://goreportcard.com/badge/github.com/lnquy/cron" alt="Code Status" /></a>
  <a href="https://travis-ci.org/lnquy/cron"><img src="https://travis-ci.org/lnquy/cron.svg?branch=master" alt="Build Status" /></a>
  <a href='https://coveralls.io/github/lnquy/cron?branch=master'><img src='https://coveralls.io/repos/github/lnquy/cron/badge.svg?branch=master' alt='Coverage Status' /></a>
  <br />
</p>

cron is a Go library that parses a cron expression and outputs a human readable description of the cron schedule.  
For example, given the expression "*/5 * * * *" it will output "Every 5 minutes".  

Translated to Go from [cron-expression-descriptor](https://github.com/bradymholt/cron-expression-descriptor) (C#) via [cRonstrue](https://github.com/bradymholt/cRonstrue) (Javascript).  
Original Author & Credit: Brady Holt (http://www.geekytidbits.com).

## Features
- Zero dependencies
- Supports all cron expression special characters including * / , - ? L W, #
- Supports 5, 6 (w/ seconds or year), or 7 (w/ seconds and year) part cron expressions
- Supports [Quartz Job Scheduler](http://www.quartz-scheduler.org/) cron expressions
- i18n support with 25 languages

## Installation
cron module can be used with Go >= 1.13.
```
go get github.com/lnquy/cron
```

## Usage

```go
// Init with default EN locale
exprDesc, _ := cron.NewDescriptor()

desc, _ := exprDesc.ToDescription("* * * * *", cron.Locale_en)
// "Every minute" 

desc, _ := exprDesc.ToDescription("0 23 ? * MON-FRI", cron.Locale_en)
// "At 11:00 PM, Monday through Friday" 

desc, _ := exprDesc.ToDescription("23 14 * * SUN#2", cron.Locale_en)
// "At 02:23 PM, on the second Sunday of the month"

// Init with custom configs
exprDesc, _ := cron.NewDescriptor(
    cron.Use24HourTimeFormat(true),
    cron.DayOfWeekStartsAtOne(true),
    cron.Verbose(true),
    cron.SetLogger(log.New(os.Stdout, "cron: ", 0)),
    cron.SetLocales(cron.Locale_en, cron.Locale_fr),
)
```

For more usage examples, including a demonstration of how cRonstrue can handle some very complex cron expressions, you can reference [the unit tests](https://github.com/lnquy/cron/blob/develop/locale_en_test.go) or [the example codes](https://github.com/lnquy/cron/tree/develop/examples).

## Demo

[...]

## i18n

To use the i18n support, you must configure the locales when create a new `ExpressionDescriptor` via `SetLocales()` option.
```go
exprDesc, _ := cron.NewDescriptor(
    cron.SetLocales(cron.Locale_en, cron.Locale_es, cron.Locale_fr),
)
// or load all cron.LocaleAll
exprDesc, _ := cron.NewDescriptor(cron.SetLocales(cron.LocaleAll))

desc, _ := exprDesc.ToDescription("* * * * *", cron.Locale_fr)
// Toutes les minutes
```

By default, `ExpressionDescriptor` always load the `Locale_en`. If you passes an unregistered locale into `ToDescription()` function, the result will be returned in English.


### Supported Locales

- en - English
- cs - Czech
- es - Spanish
- da - Danish
- de - German
- fi - Finnish
- fr - French
- fa - Farsi
- he - Hebrew
- it - Italian
- ja - Japanese
- ko - Korean
- nb - Norwegian
- nl - Dutch
- pl - Polish
- pt_BR - Portuguese (Brazil)
- ro - Romanian
- ru - Russian
- sk - Slovakian
- sl - Slovenian
- sw - Swahili
- sv - Swedish
- tr - Turkish
- uk - Ukrainian
- zh_CN - Chinese (Simplified)
- zh_TW - Chinese (Traditional)  
[TODO]


## Project status
- [x] Port 1-1 code from cRonstrue Javascript
- [X] Port and pass all test cases from cRonstrue
- [X] i18n for 25 languages
- [X] Test cases i18n
- [x] Fix i18n issues of FA, HE, RO, RU, UK, ZH_CN and ZH_TW
- [x] hcron CLI tool
- [ ] Performance improvement
- [ ] Release v1.0

## License

This project is under the MIT License. See the [LICENSE](https://github.com/lnquy/cron/blob/master/LICENSE) file for the full license text.
