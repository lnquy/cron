package main

import (
	"log"
	"os"

	"github.com/lnquy/cron"
)

const expr = "* 9 * JAN 1-5"

func main() {
	exprDesc, err := cron.NewDescriptor(
		cron.DayOfWeekStartsAtZero(),
		cron.Use24HourTimeFormat(),
		cron.Verbose(),
		cron.SetLogger(log.New(os.Stdout, "cron: ", log.LstdFlags)),
		cron.SetLocales(cron.Locale_da, cron.Locale_de, cron.Locale_en, cron.Locale_es),
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
