package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sansanbaby/dayreport/tools"
)

// ReportService 报表服务
type ReportService struct {
	tokenProvider    TokenProvider
	memberRepository MemberRepository
	reportGenerator  ReportGenerator
	emailSender      EmailSender
	reportDir        string
}

// TokenProvider Token 提供者接口
type TokenProvider interface {
	GetToken(ctx context.Context) (string, error)
}

// MemberRepository 成员仓库接口
type MemberRepository interface {
	GetAttendanceUserIDs(ctx context.Context, token string, workDate string) ([]string, error)
}

// ReportGenerator 报表生成器接口
type ReportGenerator interface {
	Generate(ctx context.Context, accessToken string, userIds []string, workDate string, filename string) error
}

// EmailSender 邮件发送器接口
type EmailSender interface {
	SendWithAttachment(subject, body, attachmentPath string) error
}

// NewReportService 创建报表服务
func NewReportService(
	tokenProvider TokenProvider,
	memberRepository MemberRepository,
	reportGenerator ReportGenerator,
	emailSender EmailSender,
	reportDir string,
) *ReportService {
	return &ReportService{
		tokenProvider:    tokenProvider,
		memberRepository: memberRepository,
		reportGenerator:  reportGenerator,
		emailSender:      emailSender,
		reportDir:        reportDir,
	}
}

// DailyReportResult 日报生成结果
type DailyReportResult struct {
	Filename string
	FileSize int64
	Count    int
}

// GenerateDailyReport 生成日报
func (s *ReportService) GenerateDailyReport(ctx context.Context) (*DailyReportResult, error) {
	// 1. 获取 Token
	token, err := s.tokenProvider.GetToken(ctx)
	if err != nil {
		return nil, tools.LogErrorf("获取 Token 失败：%w", err)
	}
	// 3. 计算日期（昨天）
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// 2. 获取考勤组成员 ID 昨天已经排班的成员
	userIds, err := s.memberRepository.GetAttendanceUserIDs(ctx, token, yesterday)
	if err != nil {
		return nil, tools.LogErrorf("获取成员列表失败：%w", err)
	}

	// 4. 确保输出目录存在
	if err := s.ensureReportDir(); err != nil {
		return nil, tools.LogErrorf("创建输出目录失败：%w", err)
	}

	// 5. 生成文件名
	filename := filepath.Join(s.reportDir, fmt.Sprintf("考勤报表_%s.xlsx", yesterday))

	// 6. 生成 Excel 报表
	if err := s.reportGenerator.Generate(ctx, token, userIds, yesterday, filename); err != nil {
		return nil, tools.LogErrorf("生成报表失败：%w", err)
	}

	// 7. 验证文件
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, tools.LogErrorf("文件验证失败：%w", err)
	}

	if fileInfo.Size() == 0 {
		return nil, tools.LogErrorf("生成的 Excel 文件为空")
	}

	// 8. 发送邮件
	if err := s.sendReportEmail(filename, yesterday); err != nil {
		return nil, tools.LogErrorf("发送邮件失败：%w", err)
	}

	return &DailyReportResult{
		Filename: filename,
		FileSize: fileInfo.Size(),
		Count:    len(userIds),
	}, nil
}

// ensureReportDir 确保输出目录存在
func (s *ReportService) ensureReportDir() error {
	if _, err := os.Stat(s.reportDir); os.IsNotExist(err) {
		return os.MkdirAll(s.reportDir, 0755)
	}
	return nil
}

// sendReportEmail 发送报表邮件
func (s *ReportService) sendReportEmail(filename, date string) error {
	subject := fmt.Sprintf("考勤日报 - %s", date)
	body := fmt.Sprintf(`
<html>
<body>
	<h2>尊敬的领导：</h2>
	<p>您好！</p>
	<p>附件为 <strong>%s</strong> 的考勤日报，请查收。</p>
	<p>此邮件由系统自动发送</p>
	<br/>
</body>
</html>
`, date)

	return s.emailSender.SendWithAttachment(subject, body, filename)
}

// GetReportDir 获取报告目录
func (s *ReportService) GetReportDir() string {
	return s.reportDir
}
