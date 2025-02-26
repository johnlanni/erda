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

// Package iteration 封装 迭代 相关操作
package iteration

import (
	"github.com/jinzhu/gorm"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/modules/cmdb/dao"
	"github.com/erda-project/erda/modules/cmdb/services/apierrors"
	"github.com/erda-project/erda/modules/cmdb/services/issue"
)

// Iteration 迭代操作封装
type Iteration struct {
	db    *dao.DBClient
	issue *issue.Issue
}

// Option 定义 Iteration 配置选项
type Option func(*Iteration)

// New 新建 Iteration 实例
func New(options ...Option) *Iteration {
	itr := &Iteration{}
	for _, op := range options {
		op(itr)
	}
	return itr
}

// WithDBClient 配置 Iteration 数据库选项
func WithDBClient(db *dao.DBClient) Option {
	return func(itr *Iteration) {
		itr.db = db
	}
}

// WithIssue 配置 issue service
func WithIssue(is *issue.Issue) Option {
	return func(itr *Iteration) {
		itr.issue = is
	}
}

// Create 创建迭代
func (itr *Iteration) Create(req *apistructs.IterationCreateRequest) (*dao.Iteration, error) {
	// 请求校验
	if req.ProjectID == 0 {
		return nil, apierrors.ErrCreateIteration.MissingParameter("projectID")
	}
	if req.Title == "" {
		return nil, apierrors.ErrCreateIteration.MissingParameter("title")
	}

	// 创建 iteration
	iteration := dao.Iteration{
		StartedAt:  req.StartedAt,
		FinishedAt: req.FinishedAt,
		ProjectID:  req.ProjectID,
		Title:      req.Title,
		Content:    req.Content,
		Creator:    req.UserID,
		State:      apistructs.IterationStateUnfiled,
	}
	if err := itr.db.CreateIteration(&iteration); err != nil {
		return nil, apierrors.ErrCreateIteration.InternalError(err)
	}
	return &iteration, nil
}

// Update 更新 iteration
func (itr *Iteration) Update(id uint64, req apistructs.IterationUpdateRequest) error {
	iteration, err := itr.db.GetIteration(id)
	if err != nil {
		return err
	}

	// 更新 iteration
	iteration.Title = req.Title
	iteration.Content = req.Content
	iteration.StartedAt = req.StartedAt
	iteration.FinishedAt = req.FinishedAt
	iteration.State = req.State
	if err := itr.db.UpdateIteration(iteration); err != nil {
		return err
	}
	return nil
}

// Get 获取 iteration 详情
func (itr *Iteration) Get(id uint64) (*dao.Iteration, error) {
	iteration, err := itr.db.GetIteration(id)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, apierrors.ErrGetIteration.NotFound()
		}
		return nil, apierrors.ErrGetIteration.InternalError(err)
	}
	return iteration, nil
}

// GetByTitle 根据 title 获取 iteration 详情
func (itr *Iteration) GetByTitle(projectID uint64, title string) (*dao.Iteration, error) {
	iteration, err := itr.db.GetIterationByTitle(projectID, title)
	if err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, apierrors.ErrGetIteration.InternalError(err)
		}
	}
	return iteration, nil
}

// Delete 删除 iteration
func (itr *Iteration) Delete(id uint64) error {
	iteration, err := itr.db.GetIteration(id)
	if err != nil {
		return err
	}

	// 检查 iteration 下是否有需求/任务/bug; 若有，不可删除
	issueReq := apistructs.IssuePagingRequest{
		IssueListRequest: apistructs.IssueListRequest{
			ProjectID:   iteration.ProjectID,
			IterationID: int64(id),
			External:    true,
		},
		PageNo:   1,
		PageSize: 1,
	}
	issues, _, err := itr.issue.Paging(issueReq)
	if err != nil {
		return err
	}
	if len(issues) > 0 {
		return apierrors.ErrDeleteIteration.InvalidParameter("该迭代下存在事件，请先删除事件后再删除迭代")
	}

	if err := itr.db.DeleteIteration(id); err != nil {
		return apierrors.ErrDeleteIteration.InternalError(err)
	}

	return nil
}

func (itr *Iteration) Paging(req apistructs.IterationPagingRequest) ([]dao.Iteration, uint64, error) {
	// 请求校验
	if req.ProjectID == 0 {
		return nil, 0, apierrors.ErrPagingIterations.InvalidParameter("missing projectID")
	}
	if req.PageNo == 0 {
		req.PageNo = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	// 创建
	iterations, total, err := itr.db.PagingIterations(req)
	if err != nil {
		return nil, 0, apierrors.ErrPagingIterations.InternalError(err)
	}

	return iterations, total, nil
}

func (itr *Iteration) GetIssueSummary(iterationID int64, projectID uint64) (apistructs.ISummary, error) {
	var taskState, bugState, requirementState []int64
	states, err := itr.db.GetIssuesStatesByProjectID(projectID, apistructs.IssueTypeTask)
	if err != nil {
		return apistructs.ISummary{}, err
	}
	// 获取每个类型的已完成状态ID
	for _, v := range states {
		if v.Belong == apistructs.IssueStateBelongDone {
			taskState = append(taskState, v.ID)
		}
	}
	states, err = itr.db.GetIssuesStatesByProjectID(projectID, apistructs.IssueTypeBug)
	if err != nil {
		return apistructs.ISummary{}, err
	}
	for _, v := range states {
		if v.Belong == apistructs.IssueStateBelongClosed {
			bugState = append(bugState, v.ID)
		}
	}
	states, err = itr.db.GetIssuesStatesByProjectID(projectID, apistructs.IssueTypeRequirement)
	if err != nil {
		return apistructs.ISummary{}, err
	}
	for _, v := range states {
		if v.Belong == apistructs.IssueStateBelongDone {
			requirementState = append(requirementState, v.ID)
		}
	}

	return itr.db.GetIssueSummary(iterationID, taskState, bugState, requirementState), nil
}
