package cron

func de_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "Jede Minute"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "Alle 5 Minuten, zwischen 03:00 PM und 03:59 PM, Montag bis Freitag"},
	}
}
