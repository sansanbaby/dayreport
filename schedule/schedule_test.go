package schedule

import (
	"fmt"
	"testing"

	"github.com/sansanbaby/dayreport/config"
)

func TestSchedule(t *testing.T) {
	//var scheduleResp *SetScheduleResp
	token, _ := config.GetAccessToken()

	//timestamp, _ := tools.DateToMillisecondTimestamp("2026-3-7")
	//timestamp, _ := tools.DatesToMillisecondTimestamps([]string{"2026-3-7", "2026-3-8"})
	////userids := []string{"196109540729105228", "01344549142634894610"}
	//
	///*scheduleResp, _ = UpdateSchedule(token, timestamp, userids, "special")*/
	///*scheduleResp, _ = ClearSchedule(token, timestamp, userids)*/
	//
	////scheduleResp, _ = UpdateScheduleByDates(token, timestamp, "196109540729105228", "common")
	//scheduleResp, _ = ClearScheduleByDates(token, timestamp, "196109540729105228")
	//t.Logf("排班设置结果：errcode=%d, success=%v, errmsg=%s, request_id=%s",
	//	scheduleResp.ErrCode, scheduleResp.Success, scheduleResp.ErrMsg, scheduleResp.RequestID)
	ss, _ := GetScheduleInfo(token, config.Config.OpUserID, []string{"196109540729105228", "323064103223254547", "03412742364626279380"}, 1773849600000, 1773849600000)
	fmt.Println(ss)

}
