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

package openapi_consumer

import (
	"context"

	"github.com/erda-project/erda/modules/hepa/common"
	"github.com/erda-project/erda/modules/hepa/gateway/dto"
	"github.com/erda-project/erda/modules/hepa/repository/orm"
)

var Service GatewayOpenapiConsumerService

type GatewayOpenapiConsumerService interface {
	Clone(context.Context) GatewayOpenapiConsumerService
	GrantPackageToConsumer(consumerId, packageId string) error
	RevokePackageFromConsumer(consumerId, packageId string) error
	CreateClientConsumer(clientName, clientId, clientSecret, clusterName string) (*orm.GatewayConsumer, error)
	CreateConsumer(*dto.DiceArgsDto, *dto.OpenConsumerDto) (string, bool, error)
	GetConsumers(*dto.GetOpenConsumersDto) (common.NewPageQuery, error)
	GetConsumersName(*dto.GetOpenConsumersDto) ([]dto.OpenConsumerInfoDto, error)
	UpdateConsumer(string, *dto.OpenConsumerDto) (*dto.OpenConsumerInfoDto, error)
	DeleteConsumer(string) (bool, error)
	// 获取调用方的所有认证凭证信息
	GetConsumerCredentials(string) (dto.ConsumerCredentialsDto, error)
	// 更新调用方的认证凭证信息
	UpdateConsumerCredentials(string, *dto.ConsumerCredentialsDto) (dto.ConsumerCredentialsDto, string, error)
	// 获取调用方能调用哪些流量入口
	GetConsumerAcls(string) ([]dto.ConsumerAclInfoDto, error)
	// 更新调用方能调用哪些流量入口
	UpdateConsumerAcls(string, *dto.ConsumerAclsDto) (bool, error)
	// 获取流量入口授权给了哪些调用方，提供比如更新流量入口的授权时使用
	GetConsumersOfPackage(string) ([]orm.GatewayConsumer, error)
	// 对应 kong 里 consumer 的 customId
	GetKongConsumerName(consumer *orm.GatewayConsumer) string
	// 获取流量入口能被哪些调用方调用
	GetPackageAcls(string) ([]dto.PackageAclInfoDto, error)
	// 更新流量入口能被哪些调用方调用
	UpdatePackageAcls(string, *dto.PackageAclsDto) (bool, error)
	// 获取流量入口下指定的 API 能被哪些调用方调用
	GetPackageApiAcls(string, string) ([]dto.PackageAclInfoDto, error)
	// 更新流量入口下指定的 API 能被哪些调用方调用
	UpdatePackageApiAcls(string, string, *dto.PackageAclsDto) (bool, error)
}
