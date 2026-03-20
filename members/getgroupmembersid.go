package members

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sansanbaby/dayreport/config"
	"github.com/sansanbaby/dayreport/tools"
)

type getGroupMembersResp struct {
	ErrCode int        `json:"errcode"`
	ErrMsg  string     `json:"errmsg"`
	Success bool       `json:"success"`
	Result  PageResult `json:"result"`
}

type PageResult struct {
	Cursor  int      `json:"cursor"`
	Result  []string `json:"result"`
	HasMore bool     `json:"has_more"`
}

type getUserRosterReq struct {
	UserIDList         []string `json:"userIdList"`
	FieldFilterList    []string `json:"fieldFilterList"`
	AppAgentID         int64    `json:"appAgentId"`
	Text2SelectConvert bool     `json:"text2SelectConvert,omitempty"`
}

type getUserRosterResp struct {
	Result []RosterInfo `json:"result"`
}

type getUserRosterErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RosterInfo struct {
	CorpID        string      `json:"corpId"`
	UserID        string      `json:"userId"`
	FieldDataList []FieldData `json:"fieldDataList"`
}

type FieldData struct {
	FieldName      string       `json:"fieldName"`
	FieldCode      string       `json:"fieldCode"`
	GroupID        string       `json:"groupId"`
	FieldValueList []FieldValue `json:"fieldValueList"`
}

type FieldValue struct {
	ItemIndex int    `json:"itemIndex"`
	Label     string `json:"label"`
	Value     string `json:"value"`
}

type UserInfo struct {
	UserID string
	Name   string
	Dept   string
}

// httpPostJSON 是一个简单的工具函数，用于发送 POST 请求并将请求体编码为 JSON 格式。
func httpPostJSON(url string, body interface{}) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, tools.LogError(err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, tools.LogError(err)
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	return client.Do(req)
}

// GetAttendanceGroupMembersId 获取考勤组成员ID
func GetAttendanceGroupMembersId(accessToken, opUserId string, groupId int) ([]string, error) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/attendance/group/memberusers/list?access_token=%s", accessToken)

	reqBody := map[string]interface{}{
		"op_user_id": opUserId,
		"group_id":   groupId,
	}

	var allMembers []string
	cursor := 0

	for {
		reqBody["cursor"] = cursor
		resp, err := httpPostJSON(url, reqBody)
		if err != nil {
			return nil, tools.LogError(err)
		}
		defer resp.Body.Close()

		var data getGroupMembersResp
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, tools.LogError(err)
		}
		if data.ErrCode != 0 {
			return nil, tools.LogErrorf("getgroupmembers error: %d %s", data.ErrCode, data.ErrMsg)
		}

		allMembers = append(allMembers, data.Result.Result...)

		if !data.Result.HasMore {
			break
		}
		cursor = data.Result.Cursor
	}

	return allMembers, nil
}

// GetUserIDByName 根据人员名字获取对应的用户 ID
func GetUserIDByName(accessToken, name string) (string, error) {
	// 首先获取考勤组所有成员的 ID
	memberIDs, err := GetAttendanceGroupMembersId(accessToken, config.Config.OpUserID, config.Config.GroupID)
	if err != nil {
		return "", tools.LogErrorf("获取考勤组成员失败：%w", err)
	}

	// 批量获取用户信息（每次最多 100 个）
	var allUserInfos []UserInfo
	batchSize := 100

	for i := 0; i < len(memberIDs); i += batchSize {
		end := i + batchSize
		if end > len(memberIDs) {
			end = len(memberIDs)
		}

		batchIDs := memberIDs[i:end]
		userInfos, err := GetUserRosterInfo(accessToken, batchIDs)
		if err != nil {
			return "", tools.LogErrorf("获取用户信息失败：%w", err)
		}

		allUserInfos = append(allUserInfos, userInfos...)
	}

	// 查找匹配名字的用户 ID
	for _, userInfo := range allUserInfos {
		if userInfo.Name == name {
			return userInfo.UserID, nil
		}
	}

	return "", tools.LogErrorf("未找到姓名为 %s 的用户", name)
}

// GetUserRosterInfo 获取用户信息姓名和部门
func GetUserRosterInfo(accessToken string, userIdList []string) ([]UserInfo, error) {
	var allUserInfos []UserInfo
	batchSize := 100

	for i := 0; i < len(userIdList); i += batchSize {
		end := i + batchSize
		if end > len(userIdList) {
			end = len(userIdList)
		}

		batchIDs := userIdList[i:end]
		infos, err := getUserRosterInfoBatch(accessToken, batchIDs)
		if err != nil {
			return nil, tools.LogError(err)
		}

		allUserInfos = append(allUserInfos, infos...)
	}

	return allUserInfos, nil
}

// getUserRosterInfoBatch 批量获取用户信息（单次请求，最多 100 个用户）
func getUserRosterInfoBatch(accessToken string, userIdList []string) ([]UserInfo, error) {
	url := "https://api.dingtalk.com/v1.0/hrm/rosters/lists/query"

	reqBody := getUserRosterReq{
		UserIDList:         userIdList,
		FieldFilterList:    []string{"sys00-name", "sys00-dept"},
		AppAgentID:         int64(config.Config.AppAgentID),
		Text2SelectConvert: false,
	}

	client := &http.Client{Timeout: 30 * time.Second}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, tools.LogError(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, tools.LogError(err)
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("x-acs-dingtalk-access-token", accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, tools.LogError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp getUserRosterErrorResp
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, tools.LogErrorf("request failed: status=%d, message=%s", resp.StatusCode, errResp.Message)
	}

	var data getUserRosterResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, tools.LogError(err)
	}

	fmt.Printf("获取到 %d 条用户信息\n", len(data.Result))

	var userInfos []UserInfo
	for _, roster := range data.Result {
		userInfo := UserInfo{
			UserID: roster.UserID,
		}

		for _, fieldData := range roster.FieldDataList {
			if len(fieldData.FieldValueList) > 0 {
				value := fieldData.FieldValueList[0].Value
				if fieldData.FieldCode == "sys00-name" {
					userInfo.Name = value
				} else if fieldData.FieldCode == "sys00-dept" {
					userInfo.Dept = value
				}
			}
		}

		userInfos = append(userInfos, userInfo)
	}

	return userInfos, nil
}
