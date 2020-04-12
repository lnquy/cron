package cron

func fr_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "Toutes les minutes"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "Toutes les 5 minutes, de 03:00 PM à 03:59 PM, de lundi à vendredi"},
	}
}
