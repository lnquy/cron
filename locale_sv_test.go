package cron

func sv_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "Varje minut"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "Var 5 minut, mellan 03:00 PM och 03:59 PM, måndag till fredag"},
		{inExpr: "0 12 * * *", outErr: nil, outDesc: "Kl 12:00 PM"},
		{inExpr: "0 15 10 ? * 6#3", outErr: nil, outDesc: "Kl 10:15 AM, den tredje lördagen av månaden"},
		{inExpr: "0 0 15 ? * MON *", outErr: nil, outDesc: "Kl 03:00 PM, varje måndag"},
	}
}
