package cron

func tr_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "Her dakika"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "Her 5 dakikada bir, 03:00 PM ile 03:59 PM arasında, Pazartesi ile Cuma arasında"},
	}
}
