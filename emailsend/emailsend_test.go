package emailsend

import "testing"

func Test_Emailsend(t *testing.T) {
	emailsender := NewEmailSender()
	emailsender.SendWithAttachment("测试", "测试", "C:\\Users\\lishuo\\Desktop\\dayreport\\report\\考勤报表_2026-03-09.xlsx")
}
