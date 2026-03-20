package attendance

import (
	"fmt"
	"testing"

	"github.com/sansanbaby/dayreport/config"
)

func Test_PersonalAttendance(t *testing.T) {
	token, _ := config.GetAccessToken()
	fmt.Println("token:", token)
	AttendanceDe, _ := GetPersonalAttendance(token, "026916505820853671", "2026-03-07")
	fmt.Println("AttendanceDe:", AttendanceDe)
	//userinfo, _ := members.GetAttendanceGroupMembersId(token, config.Config.OpUserID, config.Config.GroupID)
	//AttendanceDe, _ := BatchGetPersonalAttendance(token, userinfo, "2026-03-07")
	//fmt.Println("AttendanceDe:", AttendanceDe)
}
