package cron

func ro_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "În fiecare minut"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "La fiecare 5 minute, între 03:00 PM și 03:59 PM, de luni până vineri"},
	}
}
