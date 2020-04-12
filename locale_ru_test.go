package cron

func ru_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "Каждую минуту"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "Каждые 5 минут, с 03:00 PM по 03:59 PM, понедельник по пятница"},
	}
}
