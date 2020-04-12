package main

import (
	"log"
	"os"

	"github.com/lnquy/cron"
)

const expr = "0 * 9 LW JAN-OCT 1-5 2000/10"

// const expr = "00 11,16 * * *"

func main() {
	exprDesc, err := cron.NewDescriptor(
		cron.Use24HourTimeFormat(true),
		cron.DayOfWeekStartsAtOne(true),
		cron.Verbose(true),
		cron.SetLogger(log.New(os.Stdout, "cron: ", 0)),
		cron.SetLocales(cron.Locale_en),
	)
	if err != nil {
		log.Panicf("failed to create CRON expression descriptor: %s", err)
	}

	desc, err := exprDesc.ToDescription(expr, cron.Locale_en)
	if err != nil {
		log.Panicf("failed to convert CRON expression to human readable description: %s", err)
	}
	log.Printf("Expression: %s\nDescription: %s", expr, desc)
}
