package cron

func zh_TW_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "每分鐘"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "每 5 分鐘, 在 03:00 PM 和 03:59 PM 之間, 星期一 到 星期五"},
	}
}
