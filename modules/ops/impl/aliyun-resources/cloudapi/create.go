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

package cloudapi

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	aliyun_slb "github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	aliyun_errors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"

	"github.com/erda-project/erda/apistructs"
	aliyun_resources "github.com/erda-project/erda/modules/ops/impl/aliyun-resources"
	"github.com/erda-project/erda/modules/ops/impl/aliyun-resources/slb"
	"github.com/erda-project/erda/pkg/uuid"
)

// create ApiGateway VPC grant access
func CreateVpcGrant(ctx aliyun_resources.Context, req *apistructs.ApiGatewayVpcGrantRequest) (string, error) {
	client, err := cloudapi.NewClientWithAccessKey(ctx.Region, ctx.AccessKeyID, ctx.AccessSecret)
	if err != nil {
		return "", errors.Wrap(err, "create cloudapi client error")
	}
	client.SetReadTimeout(3 * time.Second)
	if req.Slb.Port == 0 {
		req.Slb.Port = 80
	}
	if req.Slb.ID == "" {
		sreq := aliyun_slb.CreateCreateLoadBalancerRequest()
		sreq.Scheme = "https"
		sreq.LoadBalancerName = req.Slb.Name
		sreq.VpcId = req.VpcID
		sreq.VSwitchId = req.VSwitchID
		sreq.LoadBalancerSpec = req.Slb.Spec
		if strings.ToLower(req.Slb.ChargeType) == aliyun_resources.ChargeTypePrepaid {
			sreq.PayType = "PrePay"
			sreq.PricingCycle = "month"
			sreq.Duration = requests.Integer(req.Slb.ChargePeriod)
			sreq.AutoPay = requests.NewBoolean(req.Slb.AutoRenew)
		} else if strings.ToLower(req.Slb.ChargeType) == "postpaid" {
			sreq.PayType = "PayOnDemand"
		}
		lreq := aliyun_slb.CreateCreateLoadBalancerTCPListenerRequest()
		lreq.Scheme = "https"
		lreq.ListenerPort = requests.NewInteger(req.Slb.Port)
		lreq.Bandwidth = requests.NewInteger(5120)
		lbid, err := slb.CreateIntranetSLBTCPListener(ctx, req.ClusterName, sreq, lreq)
		if err != nil {
			return "", errors.Wrap(err, "create slb failed")
		}
		req.Slb.ID = lbid
	}

	request := cloudapi.CreateDescribeVpcAccessesRequest()
	request.Scheme = "https"
	request.SecurityToken = uuid.UUID()
	grantName := fmt.Sprintf("%s_port%d", req.Slb.Name, req.Slb.Port)
	request.Name = grantName
	response, err := client.DescribeVpcAccesses(request)
	if err != nil {
		return "", errors.Wrap(err, "describe vpc access failed")
	}
	logrus.Infof("vpc access list: %+v", response.VpcAccessAttributes.VpcAccessAttribute)
	if len(response.VpcAccessAttributes.VpcAccessAttribute) > 0 {
		return grantName, nil
	}
	vpcCreateReq := cloudapi.CreateSetVpcAccessRequest()
	vpcCreateReq.Scheme = "https"
	vpcCreateReq.InstanceId = req.Slb.ID
	vpcCreateReq.SecurityToken = uuid.UUID()
	vpcCreateReq.Port = requests.NewInteger(req.Slb.Port)
	vpcCreateReq.VpcId = req.VpcID
	vpcCreateReq.Name = grantName
	_, err = client.SetVpcAccess(vpcCreateReq)
	logrus.Infof("set vpc access req: %+v", vpcCreateReq)
	if err != nil {
		if serr, ok := err.(*aliyun_errors.ServerError); ok && serr.ErrorCode() == "vpcAccessExists" {
			return vpcCreateReq.Name, nil
		}
		return "", errors.Wrap(err, "set vpc access failed")
	}
	logrus.Info("set vpc access success")
	return vpcCreateReq.Name, nil
}

func CreateAPIGateway(ctx aliyun_resources.Context, req *cloudapi.CreateInstanceRequest) (string, error) {
	client, err := cloudapi.NewClientWithAccessKey(ctx.Region, ctx.AccessKeyID, ctx.AccessSecret)
	if err != nil {
		return "", err
	}
	zoneReq := cloudapi.CreateDescribeZonesRequest()
	zoneReq.Scheme = "https"
	zoneReq.SecurityToken = uuid.UUID()
	zoneResp, err := client.DescribeZones(zoneReq)
	if err != nil {
		return "", err
	}
	if len(zoneResp.Zones.Zone) == 0 {
		return "", errors.Errorf("no available zone in this region: %s", ctx.Region)
	}
	req.ZoneId = zoneResp.Zones.Zone[0].ZoneId
	req.Scheme = "https"
	response, err := client.CreateInstance(req)
	if err != nil {
		return "", err
	}
	return response.InstanceId, nil
}
