package cron

func es_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "Cada minuto"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "Cada 5 minutos, entre las 03:00 PM y las 03:59 PM, de lunes a viernes"},
	}
}
