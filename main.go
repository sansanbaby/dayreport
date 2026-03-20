package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sansanbaby/dayreport/handler"
	"github.com/sansanbaby/dayreport/service"
	"github.com/sansanbaby/dayreport/tools"
)

// 获取输出目录
func getReportDir() string {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("获取程序路径失败：%v\n", err)
		return "report"
	}

	exeDir := filepath.Dir(exePath)
	reportDir := filepath.Join(exeDir, "report")

	return reportDir
}

// 生成日报，并生成 Excel 文件，并发送邮件
func generateDailyReport(reportService *service.ReportService) {
	fmt.Println("开始生成考勤报表...")

	result, err := reportService.GenerateDailyReport(context.Background())
	if err != nil {
		tools.LogErrorf("生成报表失败：%v", err)
		return
	}

	fmt.Printf("考勤报表生成成功：%s\n", result.Filename)
	fmt.Printf("文件大小：%d 字节\n", result.FileSize)
	fmt.Printf("共导出 %d 条记录\n", result.Count)
	fmt.Println("邮件发送成功！")
}

// 主函数，程序入口，设置定时任务每天 8 点 30 分执行 generateDailyReport 函数
func main() {
	fmt.Println("===========================================")
	fmt.Println("考勤报表自动生成服务已启动...")
	fmt.Printf("当前操作系统：%s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("每天上午 8 点 30 分自动生成前一天的考勤报表")

	reportDir := getReportDir()
	fmt.Printf("输出目录：%s\n", reportDir)

	// 初始化服务层
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

	// 启动 HTTP 服务器
	go func() {
		http.HandleFunc("/api/schedule", handler.HandleSchedule)

		port := "192.168.0.246:8080"
		//port := ":8080"
		fmt.Printf("HTTP API 服务器已启动在 http://%s\n", port)
		fmt.Println("API 端点：POST /api/schedule")
		fmt.Println("===========================================")

		if err := http.ListenAndServe(port, nil); err != nil {
			tools.LogErrorf("HTTP 服务启动失败：%v", err)
			os.Exit(1)
		}
	}()

	c := cron.New(cron.WithLocation(time.Local))

	_, err := c.AddFunc("46 8 * * *", func() {
		generateDailyReport(reportService)
	})
	if err != nil {
		tools.LogErrorf("创建定时任务失败：%v", err)
		return
	}

	fmt.Println("定时任务已设置：每天 8:30 执行")

	c.Start()

	select {}
}
