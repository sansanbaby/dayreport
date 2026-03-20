package service

import (
	"context"

	"github.com/sansanbaby/dayreport/config"
	"github.com/sansanbaby/dayreport/emailsend"
	"github.com/sansanbaby/dayreport/members"
	"github.com/sansanbaby/dayreport/printattendance"
	"github.com/sansanbaby/dayreport/schedule"
	"github.com/sansanbaby/dayreport/tools"
)

// DefaultTokenService 默认的 Token 服务
type DefaultTokenService struct{}

// NewDefaultTokenService 创建默认 Token 服务
func NewDefaultTokenService() *DefaultTokenService {
	return &DefaultTokenService{}
}

// GetToken 获取访问令牌
func (s *DefaultTokenService) GetToken(ctx context.Context) (string, error) {
	return config.GetAccessToken()
}

// DefaultMemberRepository 默认成员仓库
type DefaultMemberRepository struct{}

// NewDefaultMemberRepository 创建默认成员仓库
func NewDefaultMemberRepository() *DefaultMemberRepository {
	return &DefaultMemberRepository{}
}

// GetAttendanceUserIDs 获取考勤组成员 ID（仅返回指定日期有排班的成员）
func (r *DefaultMemberRepository) GetAttendanceUserIDs(ctx context.Context, token string, workDate string) ([]string, error) {
	memberIDs, err := members.GetAttendanceGroupMembersId(token, config.Config.OpUserID, config.Config.GroupID)
	if err != nil {
		return nil, tools.LogError(err)
	}

	timestamp, err := tools.DateToMillisecondTimestamp(workDate)
	if err != nil {
		return nil, tools.LogErrorf("日期转换失败：%w", err)
	}

	toDateTime := timestamp

	scheduleUserIDs, err := schedule.GetScheduleInfo(token, config.Config.OpUserID, memberIDs, timestamp, toDateTime)
	if err != nil {
		return nil, tools.LogError(err)
	}

	return scheduleUserIDs, nil
}

// DefaultReportGenerator 默认报表生成器
type DefaultReportGenerator struct{}

// NewDefaultReportGenerator 创建默认报表生成器
func NewDefaultReportGenerator() *DefaultReportGenerator {
	return &DefaultReportGenerator{}
}

// Generate 生成考勤报表
func (g *DefaultReportGenerator) Generate(ctx context.Context, accessToken string, userIds []string, workDate string, filename string) error {
	return printattendance.ExportAttendanceToExcel(accessToken, userIds, workDate, filename)
}

// DefaultEmailSender 默认邮件发送器
type DefaultEmailSender struct {
	sender *emailsend.EmailSender
}

// NewDefaultEmailSender 创建默认邮件发送器
func NewDefaultEmailSender() *DefaultEmailSender {
	return &DefaultEmailSender{
		sender: emailsend.NewEmailSender(),
	}
}

// SendWithAttachment 发送带附件的邮件
func (s *DefaultEmailSender) SendWithAttachment(subject, body, attachmentPath string) error {
	return s.sender.SendWithAttachment(subject, body, attachmentPath)
}
