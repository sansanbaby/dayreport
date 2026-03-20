package main

import (
	"testing"

	"github.com/sansanbaby/dayreport/service"
)

func Test_Main(t *testing.T) {
	reportDir := getReportDir()
	tokenService := service.NewDefaultTokenService()
	memberRepo := service.NewDefaultMemberRepository()
	reportGen := service.NewDefaultReportGenerator()
	emailSender := service.NewDefaultEmailSender()

	// 创建报表服务
	reportService := service.NewReportService(
		tokenService,
		memberRepo,
		reportGen,
		emailSender,
		reportDir,
	)

	generateDailyReport(reportService)
}
