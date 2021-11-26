// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package micro_api

import (
	"context"

	"github.com/erda-project/erda/modules/hepa/common"
	"github.com/erda-project/erda/modules/hepa/common/vars"
	"github.com/erda-project/erda/modules/hepa/gateway/dto"
	"github.com/erda-project/erda/modules/hepa/repository/orm"
	"github.com/erda-project/erda/modules/hepa/repository/service"
)

var Service GatewayApiService

type GatewayApiService interface {
	Clone(context.Context) GatewayApiService
	// 对应微服务 API 管理中，获取指定 runtimeservice 的 API 列表
	GetRuntimeApis(runtimeServiceId string, registerType ...string) ([]dto.ApiDto, error)
	// 在 kong 中创建 runtimeservice API 的内部接口
	CreateRuntimeApi(dto *dto.ApiDto, session ...*service.SessionHelper) (string, vars.StandardErrorCode, error)
	// 对接 kong 的 API 创建
	CreateApi(*dto.ApiReqDto) (string, error)
	// legacy openapi 接口，拿所有注册到 kong 里的 runtimeservice 以及 sdk 自动注册的 API
	GetApiInfos(*dto.GetApisDto) (*common.PageQuery, error)
	// 对接 kong 的 API 删除
	DeleteApi(string) error
	// 对接 kong 的 API 更新
	UpdateApi(string, *dto.ApiReqDto) (*dto.ApiInfoDto, error)
	// 创建 upstream api
	CreateUpstreamBindApi(*orm.GatewayConsumer, string, string, string, *orm.GatewayUpstreamApi, string) (string, error)
	UpdateUpstreamBindApi(*orm.GatewayConsumer, string, string, *orm.GatewayUpstreamApi, string) error
	DeleteUpstreamBindApi(*orm.GatewayUpstreamApi) error
	// 在 kong 中创建 runtimeservice 的 API
	TouchRuntimeApi(*orm.GatewayRuntimeService, *service.SessionHelper, bool) error
	// 对应清理 runtimeservice 的 API
	ClearRuntimeApi(*orm.GatewayRuntimeService) error
}
