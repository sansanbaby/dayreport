package schedule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sansanbaby/dayreport/config"
)

type setScheduleReq struct {
	OpUserID  string     `json:"op_user_id"`
	GroupID   int        `json:"group_id"`
	Schedules []Schedule `json:"schedules"`
}

type Schedule struct {
	ShiftID  int64  `json:"shift_id"`
	WorkDate int64  `json:"work_date"`
	IsRest   bool   `json:"is_rest"`
	UserID   string `json:"userid"`
}

type setScheduleResp struct {
	ErrCode   int    `json:"errcode"`
	Success   bool   `json:"success"`
	ErrMsg    string `json:"errmsg"`
	RequestID string `json:"request_id"`
}

// httpPostJSON 发送 POST JSON 请求
func httpPostJSON(url string, body interface{}) (*setScheduleResp, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data setScheduleResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// SetRestSchedule 设置休息排班
// 参数:
//   - accessToken: 访问令牌
//   - workDate: 工作日期（毫秒时间戳）
//   - userID: 用户 ID
//
// 返回: 钉钉 API 响应和错误
func SetRestSchedule(accessToken string, workDate int64, userID string) (*setScheduleResp, error) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/attendance/group/schedule/async?access_token=%s", accessToken)

	reqBody := &setScheduleReq{
		OpUserID: config.Config.OpUserID,
		GroupID:  config.Config.GroupID,
		Schedules: []Schedule{
			{
				ShiftID:  1,
				WorkDate: workDate,
				IsRest:   true,
				UserID:   userID,
			},
		},
	}

	return httpPostJSON(url, reqBody)
}

// ClearSchedule 清空排班
// 参数:
//   - accessToken: 访问令牌
//   - workDate: 工作日期（毫秒时间戳）
//   - userID: 用户 ID
//
// 返回: 钉钉 API 响应和错误
func ClearSchedule(accessToken string, workDate int64, userID string) (*setScheduleResp, error) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/attendance/group/schedule/async?access_token=%s", accessToken)

	reqBody := &setScheduleReq{
		OpUserID: config.Config.OpUserID,
		GroupID:  config.Config.GroupID,
		Schedules: []Schedule{
			{
				ShiftID:  -2,
				WorkDate: workDate,
				IsRest:   false,
				UserID:   userID,
			},
		},
	}

	return httpPostJSON(url, reqBody)
}

// UpdateSchedule 调整排班
// 参数:
//   - accessToken: 访问令牌
//   - workDate: 工作日期（毫秒时间戳）
//   - userID: 用户 ID
//   - scheduleType: 排班类型 ("common"=生产日常班次，"special"=生产特殊班次 1)
//
// 返回: 钉钉 API 响应和错误
func UpdateSchedule(accessToken string, workDate int64, userID string, scheduleType string) (*setScheduleResp, error) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/attendance/group/schedule/async?access_token=%s", accessToken)

	var shiftID int
	switch scheduleType {
	case "common":
		shiftID = config.Config.CommonScheduleID1
	case "special":
		shiftID = config.Config.SpecialScheduleID2
	default:
		return nil, fmt.Errorf("未知的排班类型：%s，请使用 'common' 或 'special'", scheduleType)
	}

	reqBody := &setScheduleReq{
		OpUserID: config.Config.OpUserID,
		GroupID:  config.Config.GroupID,
		Schedules: []Schedule{
			{
				ShiftID:  int64(shiftID),
				WorkDate: workDate,
				IsRest:   false,
				UserID:   userID,
			},
		},
	}

	return httpPostJSON(url, reqBody)
}
