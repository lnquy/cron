package cron

func pt_BR_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "A cada minuto"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "A cada 5 minutos, entre 03:00 PM e 03:59 PM, de segunda-feira a sexta-feira"},
	}
}
