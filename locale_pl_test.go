package cron

func pl_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "Co minutę"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "Co 5 minut, od 03:00 PM do 03:59 PM, od poniedziałek do piątek"},
	}
}
