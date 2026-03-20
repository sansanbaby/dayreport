package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sansanbaby/dayreport/service"
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

var scheduleService *service.ScheduleService

func init() {
	// 初始化排班服务
	tokenService := service.NewDefaultTokenService()
	memberRepo := service.NewDefaultMemberRepository()
	scheduleService = service.NewScheduleService(tokenService, memberRepo)
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

	// 调用服务层执行排班操作
	result, err := scheduleService.SetSchedules(
		r.Context(),
		req.UserNames,
		req.Dates,
		req.ScheduleType,
	)
	if err != nil {
		sendResponse(w, &ScheduleResponse{
			Success: false,
			Message: err.Error(),
			ErrCode: 500,
		})
		return
	}

	sendResponse(w, &ScheduleResponse{
		Success:   true,
		Message:   fmt.Sprintf("成功为 %d 个用户在 %d 个日期设置 %s 排班", len(req.UserNames), len(req.Dates), req.ScheduleType),
		RequestID: result.RequestID,
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
