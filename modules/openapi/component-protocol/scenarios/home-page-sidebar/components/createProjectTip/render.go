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

package createProjectTip

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
	protocol "github.com/erda-project/erda/modules/openapi/component-protocol"
)

type CreateProjectTip struct {
	ctxBdl     protocol.ContextBundle
	Type       string               `json:"type"`
	Props      Props                `json:"props"`
	Operations map[string]Operation `json:"operations"`
	State      State                `json:"state"`
}

type Props struct {
	Visible    bool                   `json:"visible"`
	RenderType string                 `json:"renderType"`
	Value      map[string]interface{} `json:"value"`
}

type State struct {
	ProsNum int `json:"prosNum"`
}

type Command struct {
	Key     string `json:"key"`
	Target  string `json:"target"`
	JumpOut bool   `json:"jumpOut"`
	Visible bool   `json:"visible"`
}

type Operation struct {
	Command Command `json:"command"`
	Key     string  `json:"key"`
	Reload  bool    `json:"reload"`
	Show    bool    `json:"show"`
}

func (this *CreateProjectTip) SetCtxBundle(ctx context.Context) error {
	bdl := ctx.Value(protocol.GlobalInnerKeyCtxBundle.String()).(protocol.ContextBundle)
	if bdl.Bdl == nil || bdl.I18nPrinter == nil {
		return fmt.Errorf("invalid context bundle")
	}
	logrus.Infof("inParams:%+v, identity:%+v", bdl.InParams, bdl.Identity)
	this.ctxBdl = bdl
	return nil
}

func (this *CreateProjectTip) GenComponentState(c *apistructs.Component) error {
	if c == nil || c.State == nil {
		return nil
	}
	var state State
	cont, err := json.Marshal(c.State)
	if err != nil {
		logrus.Errorf("marshal component state failed, content:%v, err:%v", c.State, err)
		return err
	}
	err = json.Unmarshal(cont, &state)
	if err != nil {
		logrus.Errorf("unmarshal component state failed, content:%v, err:%v", cont, err)
		return err
	}
	this.State = state
	return nil
}

func (this *CreateProjectTip) getProjectsNum(orgID string) (int, error) {
	orgIntId, err := strconv.Atoi(orgID)
	if err != nil {
		return 0, err
	}
	req := apistructs.ProjectListRequest{
		OrgID:    uint64(orgIntId),
		PageNo:   1,
		PageSize: 1,
	}

	projectDTO, err := this.ctxBdl.Bdl.ListMyProject(this.ctxBdl.Identity.UserID, req)
	if err != nil {
		return 0, err
	}
	if projectDTO == nil {
		return 0, nil
	}
	return projectDTO.Total, nil
}

func (p *CreateProjectTip) Render(ctx context.Context, c *apistructs.Component, scenario apistructs.ComponentProtocolScenario, event apistructs.ComponentEvent, gs *apistructs.GlobalStateData) error {
	if err := p.SetCtxBundle(ctx); err != nil {
		return err
	}
	if err := p.GenComponentState(c); err != nil {
		return err
	}
	var visible bool
	if p.ctxBdl.Identity.OrgID != "" && p.State.ProsNum == 0 {
		visible = true
	}

	p.Type = "Text"
	p.Props.Visible = visible
	p.Props.RenderType = "linkText"
	p.Props.Value = map[string]interface{}{
		"text": []interface{}{map[string]interface{}{
			"text":         "如何创建项目",
			"operationKey": "createProjectDoc",
		}, " 或 ", map[string]interface{}{
			"text":         "通过公开组织浏览公开项目信息",
			"operationKey": "toPublicOrgPage",
		}},
	}
	p.Operations = map[string]Operation{
		"createProjectDoc": {
			Command: Command{
				Key:     "goto",
				Target:  "https://docs.erda.cloud/",
				JumpOut: true,
				Visible: visible,
			},
			Key:    "click",
			Reload: false,
			Show:   false,
		},
		"toPublicOrgPage": {
			Command: Command{
				Key:     "goto",
				Target:  "orgList",
				JumpOut: true,
				Visible: visible,
			},
			Key:    "click",
			Reload: false,
			Show:   false,
		},
	}
	return nil
}

func RenderCreator() protocol.CompRender {
	return &CreateProjectTip{}
}
