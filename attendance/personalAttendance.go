package attendance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type getUpdateDataReq struct {
	UserID   string `json:"userid"`
	WorkDate string `json:"work_date"`
}

type getUpdateDataResp struct {
	ErrCode int          `json:"errcode"`
	ErrMsg  string       `json:"errmsg"`
	Success bool         `json:"success"`
	Result  UpdateResult `json:"result"`
}

type UpdateResult struct {
	ClassSettingInfo     ClassSettingInfo   `json:"class_setting_info"`
	ApproveList          []interface{}      `json:"approve_list"`
	AttendanceResultList []AttendanceResult `json:"attendance_result_list"`
	CorpID               string             `json:"corpId"`
	WorkDate             string             `json:"work_date"`
	UserID               string             `json:"userid"`
	CheckRecordList      []CheckRecord      `json:"check_record_list"`
}

type ClassSettingInfo struct {
	RestTimeVoList []RestTimeVo `json:"rest_time_vo_list"`
}

type RestTimeVo struct {
	RestBeginTime int64 `json:"rest_begin_time"`
	RestEndTime   int64 `json:"rest_end_time"`
}

type AttendanceResult struct {
	LocationMethod string `json:"location_method"`
	RecordID       int64  `json:"record_id"`
	GroupID        int    `json:"group_id"`
	LocationResult string `json:"location_result"`
	ClassID        int    `json:"class_id"`
	TimeResult     string `json:"time_result"`
	UserAddress    string `json:"user_address"`
	UserCheckTime  string `json:"user_check_time"`
	PlanCheckTime  string `json:"plan_check_time"`
	CheckType      string `json:"check_type"`
	SourceType     string `json:"source_type"`
	PlanID         int64  `json:"plan_id"`
}

type CheckRecord struct {
	RecordID      int64  `json:"record_id"`
	UserCheckTime string `json:"user_check_time"`
	ValidMatched  bool   `json:"valid_matched"`
	SourceType    string `json:"source_type"`
}

type AttendanceDetail struct {
	UserID        string
	UserCheckTime string
	//PlanCheckTime  string
	CheckType  string
	TimeResult string
	//LocationResult string
	//UserAddress    string
}

func GetPersonalAttendance(accessToken, userID, workDate string) ([]AttendanceResult, error) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/attendance/getupdatedata?access_token=%s", accessToken)

	reqBody := getUpdateDataReq{
		UserID:   userID,
		WorkDate: workDate,
	}

	client := &http.Client{Timeout: 30 * time.Second}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data getUpdateDataResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if data.ErrCode != 0 {
		return nil, fmt.Errorf("getupdatedata error: %d %s", data.ErrCode, data.ErrMsg)
	}

	return data.Result.AttendanceResultList, nil
}

func BatchGetPersonalAttendance(accessToken string, userIdList []string, workDate string) ([]AttendanceDetail, error) {
	var allDetails []AttendanceDetail

	for _, userID := range userIdList {
		results, err := GetPersonalAttendance(accessToken, userID, workDate)
		if err != nil {
			fmt.Printf("获取用户 %s 的考勤数据失败：%v\n", userID, err)
			continue
		}

		for _, result := range results {
			detail := AttendanceDetail{
				UserID:        userID,
				UserCheckTime: result.UserCheckTime,
				//PlanCheckTime:  result.PlanCheckTime,
				CheckType:  result.CheckType,
				TimeResult: result.TimeResult,
				//LocationResult: result.LocationResult,
				//UserAddress:    result.UserAddress,
			}
			allDetails = append(allDetails, detail)
		}

		time.Sleep(100 * time.Millisecond)
	}

	return allDetails, nil
}
