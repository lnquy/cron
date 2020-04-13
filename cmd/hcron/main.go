package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/lnquy/cron"
)

var (
	fLocale               string
	fInputFilePath        string
	fDayOfWeekStartsAtOne bool
	fUse24HourTimeFormat  bool
	fVerbose              bool

	acceptedCharsRegex = regexp.MustCompile(`^[wWlL /?,*#\-0-9]*$`)
)

func init() {
	flag.StringVar(&fLocale, "locale", "en", "Output description in which locale")
	flag.StringVar(&fInputFilePath, "file", "", "Path to crontab file")
	flag.BoolVar(&fDayOfWeekStartsAtOne, "dow-starts-at-one", false, "Is day of the week starts at 1 (Monday-Sunday: 1-7)")
	flag.BoolVar(&fUse24HourTimeFormat, "24-hour", false, "Output description in 24 hour time format")
	flag.BoolVar(&fVerbose, "verbose", false, "Output description in verbose format")
}

func main() {
	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Println(`
Examples:
  $ hcron "0 15 * * 1-5"
  $ hcron "0 */10 9 * * 1-5 2020"
  $ hcron -locale fr "0 */10 9 * * 1-5 2020"
  $ hcron -file /var/spool/cron/crontabs/mycronfile
  $ another-app | hcron 
  $ another-app | hcron --dow-starts-at-one --24-hour -locale es`)
	}
	flag.Parse()

	exprDesc, locale, err := getExpressionDescriptor()
	if err != nil {
		fmt.Printf("failed to init expression descriptor: %s", err)
		os.Exit(1)
	}

	// Read from stdin
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("failed to get stdin info: %s", err)
		os.Exit(1)
	}
	isPiped := (fi.Mode() & os.ModeCharDevice) == 0

	// Run in piped mode, read from the stdin until reaching EOF
	if isPiped {
		if err := stream(exprDesc, locale, bufio.NewReader(os.Stdin)); err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(1)
		}
		return
	}

	// Run in standalone mode
	// Read from crontab input file
	if strings.TrimSpace(fInputFilePath) != "" {
		f, err := os.OpenFile(fInputFilePath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Printf("failed to open file: %s", err)
			os.Exit(1)
		}
		if err := stream(exprDesc, locale, bufio.NewReader(f)); err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(1)
		}
		return
	}

	if len(os.Args) <= 1 {
		fmt.Println("cron expression must be specified")
		os.Exit(1)
	}

	// Get description for the last cmd parameter
	expr := os.Args[len(os.Args)-1]
	desc, err := exprDesc.ToDescription(expr, locale)
	if err != nil {
		fmt.Printf("invalid cron expression '%s': %s", expr, err)
		os.Exit(1)
	}
	fmt.Printf("%s: %s\n", expr, desc)
}

func getExpressionDescriptor() (exprDesc *cron.ExpressionDescriptor, locType cron.LocaleType, err error) {
	opts := []cron.Option{
		cron.Verbose(fVerbose),
		cron.Use24HourTimeFormat(fUse24HourTimeFormat),
		cron.DayOfWeekStartsAtOne(fDayOfWeekStartsAtOne),
	}

	loc, err := cron.ParseLocale(fLocale)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get locale: %w", err)
	}
	opts = append(opts, cron.SetLocales(loc))

	exprDesc, err = cron.NewDescriptor(opts...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to init cron expression descriptor: %s", err)
	}
	return exprDesc, loc, nil
}

func stream(exprDesc *cron.ExpressionDescriptor, locType cron.LocaleType, reader *bufio.Reader) error {
	for {
		line, _, err := reader.ReadLine()
		if err != nil && err == io.EOF { // Source stream closed
			return nil
		}

		expr, remaining := normalize(string(line))
		if expr == "" { // Not a parse-able cron expression
			continue
		}

		// fmt.Printf("expr: %s, remaining: %s, line: %s", expr, remaining, line)
		desc, err := exprDesc.ToDescription(expr, locType)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}

		if remaining != "" {
			fmt.Printf("%s: %s | %s\n", expr, desc, remaining)
			continue
		}
		fmt.Printf("%s: %s\n", expr, desc)
	}
}

func normalize(line string) (expr string, remainder string) {
	if strings.HasPrefix(line, "#") {
		return "", line
	}

	parts := strings.Fields(line)
	if len(parts) < 5 {
		return "", line
	}

	// Line contains invalid chars => Assume it's in crontab format
	// First 5 parts is the cron expression, the remaining is user and commands
	if !acceptedCharsRegex.MatchString(line) {
		return strings.Join(parts[:5], " "), strings.Join(parts[5:], " ")
	}

	// Only contains accepted cron characters => Assume valid cron expression
	return line, ""
}
