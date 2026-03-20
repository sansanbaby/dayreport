package printattendance

import (
	"fmt"
	"testing"

	"github.com/sansanbaby/dayreport/config"
	"github.com/sansanbaby/dayreport/members"
	"github.com/sansanbaby/dayreport/schedule"
)

func Test_PrintAttendanceReport(t *testing.T) {
	token, _ := config.GetAccessToken()
	userinfo, _ := members.GetAttendanceGroupMembersId(token, config.Config.OpUserID, config.Config.GroupID)
	fmt.Println(userinfo)
	userlist2, _ := schedule.GetScheduleInfo(token, config.Config.OpUserID, userinfo, 1773849600000, 1773849600000)
	//err1 := PrintAttendanceReport(token, userinfo, "2026-03-07")
	//if err1 != nil {
	//	t.Error("PrintAttendanceReport error")
	//}
	err2 := ExportAttendanceToExcel(token, userlist2, "2026-03-19", "2026-03-19.xlsx")
	if err2 != nil {
		t.Error("ExportAttendanceToExcel error")
	}

}
