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
	AppKey:             "xx",
	AppSecret:          "xx",
	OpUserID:           "xx1",
	GroupID:            11,
	AppAgentID:         11,
	CommonScheduleID1:  11,
	SpecialScheduleID2: 11,
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
	Username:   "11.com",
	Password:   "11",
	From:       "11com",
	To:         []string{"11"},
}
