package cron

func uk_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "Щохвилини"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "Кожні 5 хвилин, між 15:00 та 15:59, понеділок по п'ятниця"},
	}
}
