package cron

func he_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "כל דקה"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "כל 5 דקות, 15:00 עד 15:59, יום שני עד יום שישי"},
	}
}
