package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sansanbaby/dayreport/config"
	"github.com/sansanbaby/dayreport/emailsend"
	"github.com/sansanbaby/dayreport/handler"
	"github.com/sansanbaby/dayreport/members"
	"github.com/sansanbaby/dayreport/printattendance"
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

// 生成日报, 并生成 Excel 文件, 并发送邮件
func generateDailyReport() {
	fmt.Println("开始生成考勤报表...")

	token, err := config.GetAccessToken()
	if err != nil {
		fmt.Printf("获取 token 失败：%v\n", err)
		return
	}

	userIds, err := members.GetAttendanceGroupMembersId(token, config.Config.OpUserID, config.Config.GroupID)
	if err != nil {
		fmt.Printf("获取成员列表失败：%v\n", err)
		return
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	reportDir := getReportDir()
	if _, err := os.Stat(reportDir); os.IsNotExist(err) {
		if err := os.MkdirAll(reportDir, 0755); err != nil {
			fmt.Printf("创建目录失败：%v\n", err)
			return
		}
	}

	filename := filepath.Join(reportDir, fmt.Sprintf("考勤报表_%s.xlsx", yesterday))

	err = printattendance.ExportAttendanceToExcel(token, userIds, yesterday, filename)
	if err != nil {
		fmt.Printf("生成报表失败：%v\n", err)
		return
	}

	fmt.Printf("考勤报表生成成功：%s\n", filename)

	time.Sleep(2 * time.Second)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("错误：文件不存在：%s\n", filename)
		return
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("错误：无法获取文件信息：%v\n", err)
		return
	}
	fmt.Printf("文件大小：%d 字节\n", fileInfo.Size())

	if fileInfo.Size() == 0 {
		fmt.Println("错误：生成的 Excel 文件为空")
		return
	}

	subject := fmt.Sprintf("考勤日报 - %s", yesterday)
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
	`, yesterday)

	sender := emailsend.NewEmailSender()
	err = sender.SendWithAttachment(subject, body, filename)
	if err != nil {
		fmt.Printf("发送邮件失败：%v\n", err)
		return
	}

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

	// 启动 HTTP 服务器
	go func() {
		http.HandleFunc("/api/schedule", handler.HandleSchedule)

		port := "192.168.0.246:8080"
		//port := ":8080"
		fmt.Printf("HTTP API 服务器已启动在 http://%s\n", port)
		fmt.Println("API 端点：POST /api/schedule")
		fmt.Println("===========================================")

		if err := http.ListenAndServe(port, nil); err != nil {
			fmt.Printf("HTTP 服务启动失败：%v\n", err)
			os.Exit(1)
		}
	}()

	c := cron.New(cron.WithLocation(time.Local))

	_, err := c.AddFunc("30 8 * * *", func() {
		generateDailyReport()
	})
	if err != nil {
		fmt.Printf("创建定时任务失败：%v\n", err)
		return
	}

	fmt.Println("定时任务已设置：每天 8:30 执行")

	c.Start()

	select {}
}
