package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sansanbaby/dayreport/config"
	"github.com/sansanbaby/dayreport/members"
	"github.com/sansanbaby/dayreport/schedule"
	"github.com/sansanbaby/dayreport/tools"
)

// ScheduleRequest 排班请求结构
type ScheduleRequest struct {
	UserNames    []string `json:"user_names"`    // 用户名切片
	Dates        []string `json:"dates"`         // 多个日期字符串数组，格式："2026-3-17"
	ScheduleType string   `json:"schedule_type"` // 排班类型："rest"=休息，"clear"=清空，"common"=生产日常班次，"special"=生产特殊班次 1
}

// ScheduleResponse 排班响应结构
type ScheduleResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
	ErrCode   int    `json:"errcode,omitempty"`
}

// HandleSchedule 处理排班请求的 HTTP handler
func HandleSchedule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	// 只接受 POST 请求
	if r.Method != http.MethodPost {
		sendResponse(w, &ScheduleResponse{
			Success: false,
			Message: "只支持 POST 请求",
			ErrCode: 405,
		})
		return
	}

	// 解析 JSON 请求体
	var req ScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendResponse(w, &ScheduleResponse{
			Success: false,
			Message: "JSON 解析失败：" + err.Error(),
			ErrCode: 400,
		})
		return
	}

	fmt.Printf("DEBUG: 解析后的请求数据 - UserNames: %v, Dates: %v, ScheduleType: %s\n", req.UserNames, req.Dates, req.ScheduleType)

	// 验证必填参数
	if len(req.UserNames) == 0 {
		sendResponse(w, &ScheduleResponse{
			Success: false,
			Message: "用户名不能为空",
			ErrCode: 400,
		})
		return
	}

	if len(req.Dates) == 0 {
		sendResponse(w, &ScheduleResponse{
			Success: false,
			Message: "日期不能为空，请使用 'date' 字段传入日期数组",
			ErrCode: 400,
		})
		return
	}

	if req.ScheduleType == "" {
		sendResponse(w, &ScheduleResponse{
			Success: false,
			Message: "排班类型不能为空，请使用 '休息'、'清空'、'生产日常班次' 或 '生产特殊班次'",
			ErrCode: 400,
		})
		return
	}

	// 将中文排班类型转换为英文
	scheduleTypeEn := convertScheduleTypeToEnglish(req.ScheduleType)
	if scheduleTypeEn == "" {
		sendResponse(w, &ScheduleResponse{
			Success: false,
			Message: "未知的排班类型：" + req.ScheduleType + "，请使用 '休息'、'清空'、'生产日常班次' 或 '生产特殊班次'",
			ErrCode: 400,
		})
		return
	}

	// 获取 access token
	token, err := config.GetAccessToken()
	if err != nil {
		sendResponse(w, &ScheduleResponse{
			Success: false,
			Message: "获取访问令牌失败：" + err.Error(),
			ErrCode: 500,
		})
		return
	}

	// 根据用户名获取用户 ID
	userIDs := make([]string, 0, len(req.UserNames))
	for _, userName := range req.UserNames {
		userID, err := members.GetUserIDByName(token, userName)
		if err != nil {
			sendResponse(w, &ScheduleResponse{
				Success: false,
				Message: fmt.Sprintf("获取用户 %s 的 ID 失败：%v", userName, err),
				ErrCode: 404,
			})
			return
		}
		userIDs = append(userIDs, userID)
	}

	// 批量执行排班操作（支持多个日期）
	totalSuccessCount := 0
	for _, dateStr := range req.Dates {
		// 将日期转换为毫秒时间戳
		timestamp, err := tools.DateToMillisecondTimestamp(dateStr)
		if err != nil {
			sendResponse(w, &ScheduleResponse{
				Success: false,
				Message: fmt.Sprintf("日期 %s 转换失败：%v", dateStr, err),
				ErrCode: 400,
			})
			return
		}

		// 执行排班操作（批量设置同一日期的多个用户）
		var resp *schedule.SetScheduleResp
		switch scheduleTypeEn {
		case "rest":
			resp, err = schedule.SetRestSchedule(token, timestamp, userIDs)
		case "clear":
			resp, err = schedule.ClearSchedule(token, timestamp, userIDs)
		case "common":
			resp, err = schedule.UpdateSchedule(token, timestamp, userIDs, "common")
		case "special":
			resp, err = schedule.UpdateSchedule(token, timestamp, userIDs, "special")
		default:
			sendResponse(w, &ScheduleResponse{
				Success: false,
				Message: "未知的排班类型：" + scheduleTypeEn,
				ErrCode: 400,
			})
			return
		}

		if err != nil {
			sendResponse(w, &ScheduleResponse{
				Success: false,
				Message: fmt.Sprintf("日期 %s 排班操作失败：%v", dateStr, err),
				ErrCode: 500,
			})
			return
		}

		if !resp.Success {
			sendResponse(w, &ScheduleResponse{
				Success:   false,
				Message:   fmt.Sprintf("日期 %s 排班失败：%s", dateStr, resp.ErrMsg),
				ErrCode:   resp.ErrCode,
				RequestID: resp.RequestID,
			})
			return
		}

		totalSuccessCount += len(userIDs)
	}

	sendResponse(w, &ScheduleResponse{
		Success:   true,
		Message:   fmt.Sprintf("成功为 %d 个用户在 %d 个日期设置 %s 排班", len(userIDs), len(req.Dates), req.ScheduleType),
		RequestID: fmt.Sprintf("%d", totalSuccessCount),
	})
}

// convertScheduleTypeToEnglish 将中文排班类型转换为英文
// 支持的映射关系：
//   - "休息" -> "rest"
//   - "清空" -> "clear"
//   - "生产日常班次" -> "common"
//   - "生产特殊班次" -> "special"
func convertScheduleTypeToEnglish(scheduleType string) string {
	switch scheduleType {
	case "休息":
		return "rest"
	case "清空":
		return "clear"
	case "生产日常班次":
		return "common"
	case "生产特殊班次":
		return "special"
	default:
		return ""
	}
}

// sendResponse 发送 JSON 响应
func sendResponse(w http.ResponseWriter, resp *ScheduleResponse) {
	json.NewEncoder(w).Encode(resp)
}
