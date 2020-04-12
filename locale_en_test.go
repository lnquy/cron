package cron

func en_TestCases() []localeTestCase {
	return []localeTestCase{
		{inExpr: "* * * * * *", outErr: nil, outDesc: "Every second"},
		{inExpr: "* * * * *", outErr: nil, outDesc: "Every minute"},
		{inExpr: "* * * * *", isVerbose: true, outErr: nil, outDesc: "Every minute, every hour, every day"},
		{inExpr: "*/1 * * * *", outErr: nil, outDesc: "Every minute"},
		{inExpr: "*/5 * * * *", outErr: nil, outDesc: "Every 5 minutes"},
		{inExpr: "0 0/1 * * * ?", outErr: nil, outDesc: "Every minute"},
		{inExpr: "0 0 * * * ?", outErr: nil, outDesc: "Every hour"},
		{inExpr: "0 0 0/1 * * ?", outErr: nil, outDesc: "Every hour"},
		{inExpr: "* * * 3 *", outErr: nil, outDesc: "Every minute, only in March"},
		{inExpr: "* * * 3,6 *", outErr: nil, outDesc: "Every minute, only in March and June"},
		{inExpr: "* * * * * * 2013", outErr: nil, outDesc: "Every second, only in 2013"},
		{inExpr: "* * * * * 2013", outErr: nil, outDesc: "Every minute, only in 2013"},
		{inExpr: "* * * * * 2013,2014", outErr: nil, outDesc: "Every minute, only in 2013 and 2014"},
	}
}
