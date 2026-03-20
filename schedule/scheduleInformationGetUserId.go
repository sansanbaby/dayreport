package schedule

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sansanbaby/dayreport/tools"
)

// getScheduleReq 排班查询请求体
type getScheduleReq struct {
	OpUserID     string `json:"op_user_id"`
	UserIDs      string `json:"userids"`
	FromDateTime int64  `json:"from_date_time"`
	ToDateTime   int64  `json:"to_date_time"`
}

// getScheduleResp 排班查询响应体
type getScheduleResp struct {
	ErrCode   int              `json:"errcode"`
	ErrMsg    string           `json:"errmsg"`
	Success   bool             `json:"success"`
	Result    []ScheduleResult `json:"result"`
	RequestID string           `json:"request_id"`
}

// ScheduleResult 排班结果
type ScheduleResult struct {
	WorkDate      string   `json:"work_date"`
	GroupID       int      `json:"group_id"`
	ShiftVersions []int    `json:"shift_versions"`
	CorpID        string   `json:"corp_id"`
	ShiftNames    []string `json:"shift_names"`
	UserID        string   `json:"userid"`
	ShiftIDs      []int    `json:"shift_ids"`
}

// GetScheduleInfo 获取用户排班信息并提取有排班的用户 ID 列表
// accessToken: 访问令牌
// opUserID: 操作人用户 ID
// userIDs: 待查询的用户 ID 列表
// fromDateTime: 开始时间戳（毫秒）
// toDateTime: 结束时间戳（毫秒）
// 返回：有排班的用户 ID 列表
func GetScheduleInfo(accessToken, opUserID string, userIDs []string, fromDateTime, toDateTime int64) ([]string, error) {
	var allScheduleUserIDs []string
	batchSize := 20

	for i := 0; i < len(userIDs); i += batchSize {
		end := i + batchSize
		if end > len(userIDs) {
			end = len(userIDs)
		}

		batchIDs := userIDs[i:end]
		userIDsStr := ""
		for j, id := range batchIDs {
			if j == 0 {
				userIDsStr = id
			} else {
				userIDsStr += "," + id
			}
		}

		url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/attendance/schedule/shift/listbydays?access_token=%s", accessToken)

		reqBody := getScheduleReq{
			OpUserID:     opUserID,
			UserIDs:      userIDsStr,
			FromDateTime: fromDateTime,
			ToDateTime:   toDateTime,
		}

		// fmt.Printf("\n=== 批次 %d: 查询 %d 个用户 ===\n", i/batchSize+1, len(batchIDs))
		// fmt.Printf("userids: %s\n", userIDsStr)

		client := &http.Client{Timeout: 30 * time.Second}
		b, err := json.Marshal(reqBody)
		if err != nil {
			return nil, tools.LogErrorf("序列化请求失败：%w", err)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
		if err != nil {
			return nil, tools.LogErrorf("创建请求失败：%w", err)
		}
		req.Header.Set("Content-Type", "application/json;charset=utf-8")

		resp, err := client.Do(req)
		if err != nil {
			return nil, tools.LogErrorf("发送请求失败：%w", err)
		}
		defer resp.Body.Close()

		var data getScheduleResp
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, tools.LogErrorf("解析响应失败：%w", err)
		}

		if data.ErrCode != 0 {
			return nil, tools.LogErrorf("API 调用失败：%d - %s", data.ErrCode, data.ErrMsg)
		}

		for _, result := range data.Result {
			allScheduleUserIDs = append(allScheduleUserIDs, result.UserID)
		}

		// fmt.Printf("本批次找到 %d 个有排班的用户\n", len(data.Result))

		time.Sleep(100 * time.Millisecond)
	}

	// fmt.Printf("\n=== 汇总 ===\n")
	// fmt.Printf("总共查询到 %d 个有排班的用户\n", len(allScheduleUserIDs))

	return allScheduleUserIDs, nil
}
