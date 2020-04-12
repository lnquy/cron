package cron

func nb_TestCases() []localeTestCase {
	return []localeTestCase{
		// TODO: Need help
		{inExpr: "* * * * *", outErr: nil, outDesc: "Hvert minutt"},
		{inExpr: "*/5 15 * * MON-FRI", outErr: nil, outDesc: "Hvert 5 minutt, mellom 03:00 PM og 03:59 PM, mandag til og med fredag"},
		{inExpr: "0 5 1/1 * *", outErr: nil, outDesc: "Kl.05:00 AM"},
		{inExpr: "15 11 * 1/1 MON#1", outErr: nil, outDesc: "Kl.11:15 AM, på første mandag i måneden"},
		{inExpr: "15 11 * 1/5 MON#1", outErr: nil, outDesc: "Kl.11:15 AM, på første mandag i måneden, hver 5 måned"},
		{inExpr: "0 7 * * MON,TUE,THU,FRI,SUN", outErr: nil, outDesc: "Kl.07:00 AM, på mandag, tirsdag, torsdag, fredag, og søndag"},
	}
}
