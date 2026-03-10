package config

type DingDingConfig struct {
	AppKey     string
	AppSecret  string
	OpUserID   string
	GroupID    int
	AppAgentID int
}

var Config = DingDingConfig{
	AppKey:     "ding8vdd94lrv6norsyx",
	AppSecret:  "RjgCO3uTzuq9kdGujy6Esn8rsjW7esOylgrBhH-taIeYBlxo-kvOZf6LdFKH_TTu",
	OpUserID:   "011041333440857971",
	GroupID:    1358900177,
	AppAgentID: 4314663907,
}

type EmailConfig struct {
	SMTPServer string
	SMTPPort   int
	Username   string
	Password   string
	From       string
	To         []string
	ReportDir  string
}

var Email = EmailConfig{
	SMTPServer: "smtphz.qiye.163.com",
	SMTPPort:   465,
	Username:   "lisuo@cxic.com",
	Password:   "RY5hpYMghyD$Newp",
	From:       "lisuo@cxic.com",
	To:         []string{"2074307487@qq.com", "l2074307487@gmail.com"},
}
