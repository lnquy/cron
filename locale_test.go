package cron

type localeTestCase struct {
	name    string
	inExpr  string
	outErr  error
	outDesc string

	isDOWStartsAtOne   bool
	isVerbose          bool
	is12HourTimeFormat bool
}
