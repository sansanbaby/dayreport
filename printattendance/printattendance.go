package printattendance

import (
	"fmt"
	"sort"
	"time"

	"github.com/sansanbaby/dayreport/attendance"
	"github.com/sansanbaby/dayreport/members"
	"github.com/sansanbaby/dayreport/tools"
	"github.com/xuri/excelize/v2"
)

//打印考勤报表到控制台
//func PrintAttendanceReport(accessToken string, userIdList []string, workDate string) error {
//	userInfos, err := members.GetUserRosterInfo(accessToken, userIdList)
//	if err != nil {
//		return fmt.Errorf("获取员工信息失败：%v", err)
//	}
//
//	userInfoMap := make(map[string]members.UserInfo)
//	for _, info := range userInfos {
//		userInfoMap[info.UserID] = info
//	}
//
//	details, err := attendance.BatchGetPersonalAttendance(accessToken, userIdList, workDate)
//	if err != nil {
//		return fmt.Errorf("获取考勤数据失败：%v", err)
//	}
//
//	for i := 0; i < len(details); i += 2 {
//		if i+1 >= len(details) {
//			break
//		}
//
//		record1 := details[i]
//		record2 := details[i+1]
//
//		var onDuty, offDuty attendance.AttendanceDetail
//
//		if record1.CheckType == "OnDuty" {
//			onDuty = record1
//			offDuty = record2
//		} else {
//			onDuty = record2
//			offDuty = record1
//		}
//
//		userInfo := userInfoMap[onDuty.UserID]
//
//		fmt.Printf("\n员工 ID: %s\n", onDuty.UserID)
//		fmt.Printf("姓名：%s\n", userInfo.Name)
//		fmt.Printf("部门：%s\n", userInfo.Dept)
//		fmt.Println("----------------------------------------")
//
//		switch onDuty.TimeResult {
//		case "NotSigned":
//			fmt.Println("上班时间：未打卡")
//		case "Normal":
//			fmt.Printf("上班时间：%s (正常)\n", onDuty.UserCheckTime)
//		case "Late":
//			fmt.Printf("上班时间：%s (迟到)\n", onDuty.UserCheckTime)
//		case "SeriousLate":
//			fmt.Printf("上班时间：%s (严重迟到)\n", onDuty.UserCheckTime)
//		case "Absenteeism":
//			fmt.Printf("上班时间：%s (旷工迟到)\n", onDuty.UserCheckTime)
//		case "Early":
//			fmt.Printf("上班时间：%s (早退)\n", onDuty.UserCheckTime)
//		default:
//			fmt.Printf("上班时间：%s\n", onDuty.UserCheckTime)
//		}
//
//		switch offDuty.TimeResult {
//		case "NotSigned":
//			fmt.Println("下班时间：未打卡")
//		case "Normal":
//			fmt.Printf("下班时间：%s (正常)\n", offDuty.UserCheckTime)
//		case "Late":
//			fmt.Printf("下班时间：%s (迟到)\n", offDuty.UserCheckTime)
//		case "SeriousLate":
//			fmt.Printf("下班时间：%s (严重迟到)\n", offDuty.UserCheckTime)
//		case "Absenteeism":
//			fmt.Printf("下班时间：%s (旷工迟到)\n", offDuty.UserCheckTime)
//		case "Early":
//			fmt.Printf("下班时间：%s (早退)\n", offDuty.UserCheckTime)
//		default:
//			fmt.Printf("下班时间：%s\n", offDuty.UserCheckTime)
//		}
//	}
//
//	return nil
//}

type AttendanceRecord struct {
	UserID        string
	Name          string
	Dept          string
	OnDutyTime    string
	OnDutyStatus  string
	OffDutyTime   string
	OffDutyStatus string
}

// 考勤状态判断英转中
func getStatusText(status string) string {
	switch status {
	case "Normal":
		return "正常"
	case "Late":
		return "迟到"
	case "SeriousLate":
		return "严重迟到"
	case "Absenteeism":
		return "旷工迟到"
	case "Early":
		return "早退"
	case "NotSigned":
		return "未打卡"
	default:
		return status
	}
}

// 打印考勤报表到Excel文件
func ExportAttendanceToExcel(accessToken string, userIdList []string, workDate string, filename string) error {
	userInfos, err := members.GetUserRosterInfo(accessToken, userIdList)
	if err != nil {
		return tools.LogErrorf("获取员工信息失败：%v", err)
	}

	userInfoMap := make(map[string]members.UserInfo)
	for _, info := range userInfos {
		userInfoMap[info.UserID] = info
	}

	details, err := attendance.BatchGetPersonalAttendance(accessToken, userIdList, workDate)
	if err != nil {
		return tools.LogErrorf("获取考勤数据失败：%v", err)
	}

	userRecords := make(map[string][]attendance.AttendanceDetail)
	for _, detail := range details {
		userRecords[detail.UserID] = append(userRecords[detail.UserID], detail)
	}

	var recordList []*AttendanceRecord
	for userID, userDetailList := range userRecords {
		userInfo := userInfoMap[userID]

		record := &AttendanceRecord{
			UserID: userID,
			Name:   userInfo.Name,
			Dept:   userInfo.Dept,
		}

		for _, detail := range userDetailList {
			if detail.CheckType == "OnDuty" {
				record.OnDutyTime = detail.UserCheckTime
				record.OnDutyStatus = detail.TimeResult
				if detail.TimeResult == "NotSigned" {
					record.OnDutyTime = "未打卡"
				}
			} else if detail.CheckType == "OffDuty" {
				record.OffDutyTime = detail.UserCheckTime
				record.OffDutyStatus = detail.TimeResult
				if detail.TimeResult == "NotSigned" {
					record.OffDutyTime = "未打卡"
				}
			}
		}

		recordList = append(recordList, record)
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			tools.LogError(err)
		}
	}()

	sheetName := "考勤报表"
	f.SetSheetName("Sheet1", sheetName)

	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	timestampStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 11},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#CCCCCC"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	lateStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFFF00"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	seriousLateStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFA500"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	absenteeismStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FF0000"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	earlyStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FFC0CB"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	notSignedStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#808080"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	nameWarningStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#FF0000"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})

	getStyle := func(status string) int {
		switch status {
		case "Normal":
			return 0
		case "Late":
			return lateStyle
		case "SeriousLate":
			return seriousLateStyle
		case "Absenteeism":
			return absenteeismStyle
		case "Early":
			return earlyStyle
		case "NotSigned":
			return notSignedStyle
		default:
			return 0
		}
	}

	f.MergeCell(sheetName, "A1", "G1")
	f.SetCellValue(sheetName, "A1", fmt.Sprintf("考勤报表 - %s", workDate))
	f.SetCellStyle(sheetName, "A1", "A1", titleStyle)

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	f.MergeCell(sheetName, "A2", "G2")
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("生成时间：%s", currentTime))
	f.SetCellStyle(sheetName, "A2", "A2", timestampStyle)

	headers := []string{"员工 ID", "员工姓名", "员工部门", "上班时间", "上班状态", "下班时间", "下班状态"}
	for col, header := range headers {
		cell := fmt.Sprintf("%c3", 'A'+col)
		f.SetCellValue(sheetName, cell, header)
	}
	for col := 0; col < len(headers); col++ {
		cell := fmt.Sprintf("%c3", 'A'+col)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	sort.Slice(recordList, func(i, j int) bool {
		if recordList[i].Dept != recordList[j].Dept {
			return recordList[i].Dept < recordList[j].Dept
		}
		return recordList[i].Name < recordList[j].Name
	})

	rowNum := 4
	for _, record := range recordList {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), record.UserID)

		hasAbnormal := record.OnDutyStatus != "Normal" || record.OffDutyStatus != "Normal"
		if hasAbnormal {
			f.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowNum), fmt.Sprintf("B%d", rowNum), nameWarningStyle)
		}

		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), record.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowNum), record.Dept)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowNum), record.OnDutyTime)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowNum), getStatusText(record.OnDutyStatus))
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", rowNum), record.OffDutyTime)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", rowNum), getStatusText(record.OffDutyStatus))

		onDutyStyleID := getStyle(record.OnDutyStatus)
		cellE := fmt.Sprintf("E%d", rowNum)
		f.SetCellStyle(sheetName, cellE, cellE, onDutyStyleID)

		offDutyStyleID := getStyle(record.OffDutyStatus)
		cellG := fmt.Sprintf("G%d", rowNum)
		f.SetCellStyle(sheetName, cellG, cellG, offDutyStyleID)

		rowNum++
	}

	columnWidths := []float64{15, 20, 25, 20, 15, 20, 15}
	for i, width := range columnWidths {
		colName := fmt.Sprintf("%c", 'A'+i)
		f.SetColWidth(sheetName, colName, colName, width)
	}

	f.SetRowHeight(sheetName, 1, 25)
	f.SetRowHeight(sheetName, 2, 20)
	f.SetRowHeight(sheetName, 3, 25)

	if filename == "" {
		filename = fmt.Sprintf("考勤报表_%s.xlsx", workDate)
	}
	if err := f.SaveAs(filename); err != nil {
		return tools.LogErrorf("保存 Excel 文件失败：%v", err)
	}

	fmt.Printf("考勤报表已导出到：%s\n", filename)
	fmt.Printf("共导出 %d 条记录\n", len(recordList))

	return nil
}
