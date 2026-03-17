package schedule

import (
	"testing"

	"github.com/sansanbaby/dayreport/config"
	"github.com/sansanbaby/dayreport/tools"
)

func TestSchedule(t *testing.T) {
	var scheduleResp *setScheduleResp
	token, _ := config.GetAccessToken()
	timestamp, _ := tools.DateToMillisecondTimestamp("2026-3-9")
	scheduleResp, _ = UpdateSchedule(token, timestamp, "196109540729105228", "common")

	t.Logf("排班设置结果：errcode=%d, success=%v, errmsg=%s, request_id=%s",
		scheduleResp.ErrCode, scheduleResp.Success, scheduleResp.ErrMsg, scheduleResp.RequestID)
}
