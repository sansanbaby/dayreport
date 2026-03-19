package schedule

import (
	"testing"

	"github.com/sansanbaby/dayreport/config"
	"github.com/sansanbaby/dayreport/tools"
)

func TestSchedule(t *testing.T) {
	var scheduleResp *setScheduleResp
	token, _ := config.GetAccessToken()
	//timestamp, _ := tools.DateToMillisecondTimestamp("2026-3-7")
	timestamp, _ := tools.DatesToMillisecondTimestamps([]string{"2026-3-7", "2026-3-8"})
	//userids := []string{"196109540729105228", "01344549142634894610"}

	/*scheduleResp, _ = UpdateSchedule(token, timestamp, userids, "special")*/
	/*scheduleResp, _ = ClearSchedule(token, timestamp, userids)*/

	//scheduleResp, _ = UpdateScheduleByDates(token, timestamp, "196109540729105228", "common")
	scheduleResp, _ = ClearScheduleByDates(token, timestamp, "196109540729105228")
	t.Logf("排班设置结果：errcode=%d, success=%v, errmsg=%s, request_id=%s",
		scheduleResp.ErrCode, scheduleResp.Success, scheduleResp.ErrMsg, scheduleResp.RequestID)
}
