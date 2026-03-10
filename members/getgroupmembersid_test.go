package members

import (
	"fmt"
	"testing"

	"github.com/sansanbaby/dayreport/config"
)

func TestGetAttendanceGroupMembersId(t *testing.T) {
	token, _ := config.GetAccessToken()
	fmt.Println("token:", token)
	members, err := GetAttendanceGroupMembersId(token, config.Config.OpUserID, config.Config.GroupID)
	if err != nil {
		panic(err)
	}
	//for _, member := range members {
	//	println(member)
	//}
	UserInfom, err2 := GetUserRosterInfo(token, members)
	if err2 != nil {
		panic(err2)
	}
	for _, member := range UserInfom {
		println(member.UserID, member.Name, member.Dept)
	}
}
