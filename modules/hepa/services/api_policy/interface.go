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

package api_policy

import (
	"context"

	"github.com/erda-project/erda/modules/hepa/apipolicy"
	"github.com/erda-project/erda/modules/hepa/repository/orm"
	"github.com/erda-project/erda/modules/hepa/repository/service"
)

var Service GatewayApiPolicyService

type GatewayApiPolicyService interface {
	Clone(context.Context) GatewayApiPolicyService
	// 1. 流量入口下的全局策略基于此生效
	// 2. 提供 hepa 内部调用，为一个流量入口及其下的所有 API 的 zone 生效指定的策略
	SetPackageDefaultPolicyConfig(category, packageId string, az *orm.GatewayAzInfo, config []byte, helper ...*service.SessionHelper) (string, error)
	GetPolicyConfig(category, packageId, packageApiId string) (interface{}, error)
	SetPolicyConfig(category, packageId, packageApiId string, config []byte) (interface{}, error)
	// 更新 zone 的 ingress 上的 Annotation
	RefreshZoneIngress(zone orm.GatewayZone, az orm.GatewayAzInfo) error
	// 为 zone 设置生效指定策略
	SetZonePolicyConfig(zone *orm.GatewayZone, category string, config []byte, helper *service.SessionHelper, needDeployTag ...bool) (apipolicy.PolicyDto, string, error)
	// zone 初建时，获取所属流量入口的全局策略进行生效
	SetZoneDefaultPolicyConfig(packageId string, zone *orm.GatewayZone, az *orm.GatewayAzInfo, session ...*service.SessionHelper) (map[string]*string, *string, *service.SessionHelper, error)
}
