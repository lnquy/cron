package cron

func nl_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "Elke minuut"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "Elke 5 minuten, tussen 03:00 PM en 03:59 PM, maandag t/m vrijdag"},
	}
}
