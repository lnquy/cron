package cron

func it_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "Ogni minuto"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "Ogni 5 minuti, tra le 03:00 PM e le 03:59 PM, lunedì al venerdì"},
	}
}
