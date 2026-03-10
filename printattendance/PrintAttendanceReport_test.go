package printattendance

import (
	"testing"

	"github.com/sansanbaby/dayreport/config"
	"github.com/sansanbaby/dayreport/members"
)

func Test_PrintAttendanceReport(t *testing.T) {
	token, _ := config.GetAccessToken()
	userinfo, _ := members.GetAttendanceGroupMembersId(token, config.Config.OpUserID, config.Config.GroupID)
	//err1 := PrintAttendanceReport(token, userinfo, "2026-03-07")
	//if err1 != nil {
	//	t.Error("PrintAttendanceReport error")
	//}
	err2 := ExportAttendanceToExcel(token, userinfo, "2026-03-07", "2026-03-07.xlsx")
	if err2 != nil {
		t.Error("ExportAttendanceToExcel error")
	}

}
