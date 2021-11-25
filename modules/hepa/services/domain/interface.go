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

package domain

import (
	"context"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/hepa/common"
	"github.com/erda-project/erda/modules/hepa/endpoint"
	"github.com/erda-project/erda/modules/hepa/gateway/dto"
	"github.com/erda-project/erda/modules/hepa/repository/orm"
	"github.com/erda-project/erda/modules/hepa/repository/service"
	"github.com/erda-project/erda/pkg/parser/diceyml"
)

var Service GatewayDomainService

type GatewayDomainService interface {
	// 云管平台应用资源中的域名列表分页接口
	GetOrgDomainInfo(*dto.ManageDomainReq) (common.NewPageQuery, error)
	// 获取项目环境租户下的所有域名
	GetTenantDomains(projectId, env string) ([]string, error)
	// 在部署详情页获取对应服务绑定的域名
	GetRuntimeDomains(runtimeId string) (dto.RuntimeDomainsDto, error)
	// 在部署详情页更新对应服务绑定的域名
	UpdateRuntimeServiceDomain(orgId, runtimeId, serviceName string, reqDto *dto.ServiceDomainReqDto) (bool, string, error)
	CreateOrUpdateComponentIngress(apistructs.ComponentIngressUpdateRequest) (bool, error)
	Clone(context.Context) GatewayDomainService
	FindDomains(domain, projectId, workspace string, matchType orm.OptionType, domainType ...string) ([]orm.GatewayDomain, error)
	// 获取 erda.yml 中的 service 端口，用于更新 ingress
	UpdateRuntimeServicePort(runtimeService *orm.GatewayRuntimeService, releaseInfo *diceyml.Object) error
	// 基于runtimeservice 信息生成最新的 material，去更新 ingress
	RefreshRuntimeDomain(runtimeService *orm.GatewayRuntimeService, session *service.SessionHelper) error
	// 将绑定在服务上的域名同步给流量入口，用于 runtimeservice 接收应用部署成功后的处理
	GiveRuntimeDomainToPackage(runtimeService *orm.GatewayRuntimeService, session *service.SessionHelper) (bool, error)
	// 更新 runtime 以及 runtime 关联的流量入口的 ingress
	TouchRuntimeDomain(orgId string, runtimeService *orm.GatewayRuntimeService, material endpoint.EndpointMaterial, domains []dto.EndpointDomainDto, audits *[]apistructs.Audit, session *service.SessionHelper) (string, error)
	// 更新指定流量入口关联的域名，只更新数据库，不更新 ingress
	TouchPackageDomain(orgId, packageId, clusterName string, domains []string, session *service.SessionHelper) ([]string, error)
	// 获取流量入口下的所有域名
	GetPackageDomains(packageId string, session ...*service.SessionHelper) ([]string, error)
	// 判断流量入口的域名是否发生了变化
	IsPackageDomainsDiff(packageId, clusterName string, domains []string, session *service.SessionHelper) (bool, error)
}
