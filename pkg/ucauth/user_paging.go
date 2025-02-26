// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package ucauth

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/discover"
	"github.com/erda-project/erda/pkg/httpclient"
	"github.com/erda-project/erda/pkg/strutil"
	"github.com/pkg/errors"
)

func HandlePagingUsers(req *apistructs.UserPagingRequest, token OAuthToken) (*userPaging, error) {
	if token.TokenType == OryCompatibleClientId {
		users, err := getUserPage(token.AccessToken, req.PageNo, req.PageSize)
		if err != nil {
			return nil, err
		}
		var p userPaging
		p.Total = 1000
		for _, u := range users {
			p.Data = append(p.Data, userToUserInPaging(u))
		}
		return &p, nil
	}
	v := httpclient.New().Get(discover.UC()).Path("/api/user/admin/paging").
		Header("Authorization", strutil.Concat("Bearer ", token.AccessToken))
	if req.Name != "" {
		v.Param("username", req.Name)
	}
	if req.Nick != "" {
		v.Param("nickname", req.Nick)
	}
	if req.Phone != "" {
		v.Param("mobile", req.Phone)
	}
	if req.Email != "" {
		v.Param("email", req.Email)
	}
	if req.Locked != nil {
		v.Param("locked", strconv.Itoa(*req.Locked))
	}
	if req.Source != "" {
		v.Param("source", req.Source)
	}
	if req.PageNo > 0 {
		v.Param("pageNo", strconv.Itoa(req.PageNo))
	}
	if req.PageSize > 0 {
		v.Param("pageSize", strconv.Itoa(req.PageSize))
	}
	// 批量查询用户
	var resp struct {
		Success bool        `json:"success"`
		Result  *userPaging `json:"result"`
		Error   string      `json:"error"`
	}
	r, err := v.Do().JSON(&resp)
	if err != nil {
		return nil, err
	}
	if !r.IsOK() {
		return nil, fmt.Errorf("internal status code: %v", r.StatusCode())
	}
	if !resp.Success {
		return nil, errors.New(resp.Error)
	}
	return resp.Result, nil
}

func ConvertToUserInfoExt(user *userPaging) *apistructs.UserPagingData {
	var ret apistructs.UserPagingData
	ret.Total = user.Total
	ret.List = make([]apistructs.UserInfoExt, 0)
	for _, u := range user.Data {
		ret.List = append(ret.List, apistructs.UserInfoExt{
			UserInfo: apistructs.UserInfo{
				ID:          u.Id,
				Name:        u.Username,
				Nick:        u.Nickname,
				Avatar:      u.Avatar,
				Phone:       u.Mobile,
				Email:       u.Email,
				LastLoginAt: time.Time(u.LastLoginAt).Format("2006-01-02 15:04:05"),
				PwdExpireAt: time.Time(u.PwdExpireAt).Format("2006-01-02 15:04:05"),
				Source:      u.Source,
			},
			Locked: u.Locked,
		})
	}
	return &ret
}

// userInPaging 用户中心分页用户数据结构
type userInPaging struct {
	Id            string      `json:"id"`            // 主键
	Avatar        string      `json:"avatar"`        // 头像
	Username      string      `json:"username"`      // 用户名
	Nickname      string      `json:"nickname"`      // 昵称
	Mobile        string      `json:"mobile"`        // 手机号
	Email         string      `json:"email"`         // 邮箱
	Enabled       bool        `json:"enabled"`       // 是否启用
	UserDetail    interface{} `json:"userDetail"`    // 用户详细信息
	Locked        bool        `json:"locked"`        // 冻结FLAG(0:NOT,1:YES)
	PasswordExist bool        `json:"passwordExist"` // 密码是否存在
	PwdExpireAt   timestamp   `json:"pwdExpireAt"`   // 过期时间
	Extra         interface{} `json:"extra"`         // 扩展字段
	Source        string      `json:"source"`        // 用户来源
	SourceType    string      `json:"sourceType"`    // 来源类型
	Tag           string      `json:"tag"`           // 标签
	Channel       string      `json:"channel"`       // 注册渠道
	ChannelType   string      `json:"channelType"`   // 渠道类型
	TenantId      int         `json:"tenantId"`      // 租户ID
	CreatedAt     timestamp   `json:"createdAt"`     // 创建时间
	UpdatedAt     timestamp   `json:"updatedAt"`     // 更新时间
	LastLoginAt   timestamp   `json:"lastLoginAt"`   // 最后登录时间
}

// millisecond epoch
type timestamp time.Time

func (t timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t *timestamp) UnmarshalJSON(s []byte) (err error) {
	r := strings.Replace(string(s), `"`, ``, -1)
	if r == "null" {
		return
	}

	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q/1000, 0)
	return
}

type userPaging struct {
	Data  []userInPaging `json:"data"`
	Total int            `json:"total"`
}
