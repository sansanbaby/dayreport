package service

import (
	"context"
	"fmt"

	"github.com/sansanbaby/dayreport/members"
	"github.com/sansanbaby/dayreport/schedule"
	"github.com/sansanbaby/dayreport/tools"
)

// ScheduleService 排班服务
type ScheduleService struct {
	tokenProvider    TokenProvider
	memberRepository MemberRepository
}

// NewScheduleService 创建排班服务
func NewScheduleService(tokenProvider TokenProvider, memberRepository MemberRepository) *ScheduleService {
	return &ScheduleService{
		tokenProvider:    tokenProvider,
		memberRepository: memberRepository,
	}
}

// ScheduleResult 排班结果
type ScheduleResult struct {
	SuccessCount int
	RequestID    string
}

// SetSchedules 批量设置排班
func (s *ScheduleService) SetSchedules(
	ctx context.Context,
	userNames []string,
	dates []string,
	scheduleType string,
) (*ScheduleResult, error) {
	// 1. 获取 Token
	token, err := s.tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取 Token 失败：%w", err)
	}

	// 2. 将中文排班类型转换为英文
	scheduleTypeEn := convertScheduleTypeToEnglish(scheduleType)
	if scheduleTypeEn == "" {
		return nil, fmt.Errorf("未知的排班类型：%s", scheduleType)
	}

	// 3. 根据用户名获取用户 ID
	userIDs := make([]string, 0, len(userNames))
	for _, userName := range userNames {
		userID, err := members.GetUserIDByName(token, userName)
		if err != nil {
			return nil, fmt.Errorf("获取用户 %s 的 ID 失败：%w", userName, err)
		}
		userIDs = append(userIDs, userID)
	}

	// 4. 批量执行排班操作
	totalSuccessCount := 0
	var lastRequestID string

	for _, dateStr := range dates {
		// 将日期转换为毫秒时间戳
		timestamp, err := tools.DateToMillisecondTimestamp(dateStr)
		if err != nil {
			return nil, fmt.Errorf("日期 %s 转换失败：%w", dateStr, err)
		}

		// 执行排班操作
		resp, err := s.executeSchedule(token, timestamp, userIDs, scheduleTypeEn)
		if err != nil {
			return nil, fmt.Errorf("日期 %s 排班操作失败：%w", dateStr, err)
		}

		if !resp.Success {
			return nil, fmt.Errorf("日期 %s 排班失败：%s", dateStr, resp.ErrMsg)
		}

		totalSuccessCount += len(userIDs)
		lastRequestID = resp.RequestID
	}

	return &ScheduleResult{
		SuccessCount: totalSuccessCount,
		RequestID:    lastRequestID,
	}, nil
}

// executeSchedule 执行排班操作
func (s *ScheduleService) executeSchedule(token string, timestamp int64, userIDs []string, scheduleType string) (*schedule.SetScheduleResp, error) {
	switch scheduleType {
	case "rest":
		return schedule.SetRestSchedule(token, timestamp, userIDs)
	case "clear":
		return schedule.ClearSchedule(token, timestamp, userIDs)
	case "common":
		return schedule.UpdateSchedule(token, timestamp, userIDs, "common")
	case "special":
		return schedule.UpdateSchedule(token, timestamp, userIDs, "special")
	default:
		return nil, fmt.Errorf("未知的排班类型：%s", scheduleType)
	}
}

// convertScheduleTypeToEnglish 将中文排班类型转换为英文
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
