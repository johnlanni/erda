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

package endpoint_api

import (
	"context"

	"github.com/erda-project/erda/modules/hepa/common"
	"github.com/erda-project/erda/modules/hepa/gateway/dto"
	"github.com/erda-project/erda/modules/hepa/repository/orm"
	"github.com/erda-project/erda/modules/hepa/repository/service"
	"github.com/erda-project/erda/modules/hepa/services/runtime_service"
)

var Service GatewayOpenapiService

type PackageApiInfo struct {
	*orm.GatewayPackageApi
	Hosts               []string
	ProjectId           string
	Env                 string
	Az                  string
	InjectRuntimeDomain bool
}

type GatewayOpenapiService interface {
	Clone(context.Context) GatewayOpenapiService
	// 清理所有流量入口内和对应 runtime service 关联的 api
	ClearRuntimeRoute(id string) error
	// 创建或更新 erda.yml 中指定的 endpoint 和 endpoint 的 api
	SetRuntimeEndpoint(runtime_service.RuntimeEndpointInfo) error
	// 创建或更新 APIM 里访问管理条目对应的根路径转发
	TouchPackageRootApi(packageId string, reqDto *dto.OpenapiDto) (bool, error)
	// 在 runtimeservice 的关联发生变化时（服务删除了、不依赖网关了）尝试清理对应的流量入口
	TryClearRuntimePackage(*orm.GatewayRuntimeService, *service.SessionHelper, ...bool) error
	// 创建或更新和 runtimeservice 关联的流量入口的元信息
	TouchRuntimePackageMeta(*orm.GatewayRuntimeService, *service.SessionHelper) (string, bool, error)
	// 更新和 runtimeservice 关联的流量入口下的所有 API 相关的 ingress
	RefreshRuntimePackage(string, *orm.GatewayRuntimeService, *service.SessionHelper) error
	// 为了 api policy 中适配统一域名入口，创建对应的 zone
	CreateUnityPackageZone(string, *service.SessionHelper) (*orm.GatewayZone, error)
	// 用于在创建项目租户时，创建对应的统一域名入口
	CreateTenantPackage(string, *service.SessionHelper) error

	// 创建流量入口
	CreatePackage(*dto.DiceArgsDto, *dto.PackageDto) (*dto.PackageInfoDto, string, error)
	// 获取流量入口列表
	GetPackages(*dto.GetPackagesDto) (common.NewPageQuery, error)
	// 获取流量入口详情
	GetPackage(string) (*dto.PackageInfoDto, error)
	// 获取租户下流量入口的名称
	GetPackagesName(*dto.GetPackagesDto) ([]dto.PackageInfoDto, error)
	// 编辑流量入口
	UpdatePackage(string, string, *dto.PackageDto) (*dto.PackageInfoDto, error)
	// 删除流量入口
	DeletePackage(string) (bool, error)
	// 创建流量入口下的 API
	CreatePackageApi(string, *dto.OpenapiDto) (string, bool, error)
	// 获取流量入口下的 API 列表
	GetPackageApis(string, *dto.GetOpenapiDto) (common.NewPageQuery, error)
	// 更新流量入口内的 API
	UpdatePackageApi(string, string, *dto.OpenapiDto) (*dto.OpenapiInfoDto, bool, error)
	// 删除流量入口内的 API
	DeletePackageApi(string, string) (bool, error)
	// 创建或修改流量入口 API 对应的 zone，同时创建对应的 ingress
	TouchPackageApiZone(info PackageApiInfo, session ...*service.SessionHelper) (string, error)
}
