package config

type DingDingConfig struct {
	AppKey             string
	AppSecret          string
	OpUserID           string
	GroupID            int
	AppAgentID         int
	CommonScheduleID1  int
	SpecialScheduleID2 int
}

var Config = DingDingConfig{
	AppKey:             "ding8vdd94lrv6norsyx",
	AppSecret:          "RjgCO3uTzuq9kdGujy6Esn8rsjW7esOylgrBhH-taIeYBlxo-kvOZf6LdFKH_TTu",
	OpUserID:           "011041333440857971",
	GroupID:            1358900177,
	AppAgentID:         4314663907,
	CommonScheduleID1:  1508085204,
	SpecialScheduleID2: 1436380221,
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
	To:         []string{"fuwei@cxic.com", "heyifei@cxic.com", "gongxianzhong@cxic.com", "qiyuansheng@cxic.com", "xym@cxic.com", "yujue@cxic.com", "qianjun@cxic.com", "xym@cxic.com", "yujue@cxic.com", "qianjun@cxic.com", "lisuo@cxic.com"},
}
