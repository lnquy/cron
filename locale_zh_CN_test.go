package cron

func zh_CN_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "每分钟"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "每隔 5 分钟, 在 下午 03:00 和 下午 03:59 之间, 星期一至星期五"},
	}
}
