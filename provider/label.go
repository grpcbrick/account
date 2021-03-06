package provider

import (
	"context"

	"github.com/grpcbrick/account/dao"
	"github.com/grpcbrick/account/standard"
	validators "github.com/grpcbrick/account/validators"
)

// CreateLabel 创建标签
func (srv *Service) CreateLabel(ctx context.Context, req *standard.CreateLabelRequest) (resp *standard.CreateLabelResponse, err error) {
	resp = new(standard.CreateLabelResponse)
	if ok, msg := validators.LabelName(req.Name); ok != true {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = msg
		return resp, nil
	}

	if ok, msg := validators.LabelCategory(req.Category); ok != true {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = msg
		return resp, nil
	}

	if ok, msg := validators.LabelState(req.State); ok != true {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = msg
		return resp, nil
	}

	if ok, msg := validators.LabelValue(req.Value); ok != true {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = msg
		return resp, nil
	}

	id, err := dao.CreateLabel(req.Name, req.Category, req.State, req.Value)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	// 查询数据
	queryResult, err := srv.QueryLabelByID(ctx, &standard.QueryLabelByIDRequest{ID: id})
	if err != nil {
		resp.State = standard.State_SERVICE_ERROR
		resp.Message = err.Error()
		return resp, nil
	}

	// 查询失败了
	if queryResult.State != standard.State_SUCCESS {
		resp.State = queryResult.State
		resp.Message = queryResult.Message
		return resp, nil
	}

	resp.State = queryResult.State
	resp.Data = queryResult.Data
	resp.Message = "创建成功"
	return resp, nil
}

// QueryLabels 查询标签
func (srv *Service) QueryLabels(ctx context.Context, req *standard.QueryLabelsRequest) (resp *standard.QueryLabelsResponse, err error) {
	resp = new(standard.QueryLabelsResponse)

	if req.Page <= 0 || req.Limit <= 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的参数"
		return resp, nil
	}

	totalPage, currentPage, labels, err := dao.QueryLabels(req.Page, req.Limit)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	data := []*standard.Label{}
	for _, label := range labels {
		data = append(data, label.OutProtoStruct())
	}

	resp.State = standard.State_SUCCESS
	resp.CurrentPage = currentPage
	resp.TotalPage = totalPage
	resp.Message = "查询成功"
	resp.Data = data
	return resp, nil
}

// CreateLabelForUser 给指定用户创建标签
func (srv *Service) CreateLabelForUser(ctx context.Context, req *standard.CreateLabelForUserRequest) (resp *standard.CreateLabelForUserResponse, err error) {
	resp = new(standard.CreateLabelForUserResponse)

	// 先判断条件
	userCount, err := dao.CountUserByID(req.UserID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if userCount <= 0 { // 没有找到用户
		resp.State = standard.State_USER_NOT_EXIST
		resp.Message = "该用户不存在"
		return resp, nil
	}

	createResult, err := srv.CreateLabel(ctx, &standard.CreateLabelRequest{Name: req.Name, Category: req.Category, State: req.State, Value: req.Value})
	if err != nil {
		resp.State = standard.State_SERVICE_ERROR
		resp.Message = err.Error()
		return resp, nil
	}

	if createResult.State != standard.State_SUCCESS {
		resp.State = createResult.State
		resp.Message = createResult.Message
		return resp, nil
	}

	addResult, err := srv.AddLabelToUserByID(ctx, &standard.AddLabelToUserByIDRequest{ID: createResult.Data.ID, UserID: req.UserID})
	if err != nil {
		resp.State = standard.State_SERVICE_ERROR
		resp.Message = err.Error()
		return resp, nil
	}

	// 添加失败了
	if addResult.State != standard.State_SUCCESS {
		resp.State = addResult.State
		resp.Message = addResult.Message
		return resp, nil
	}

	resp.State = addResult.State
	resp.Data = createResult.Data
	resp.Message = "创建成功"
	return resp, nil
}

// QueryLabelByID 通过 ID 查询
func (srv *Service) QueryLabelByID(ctx context.Context, req *standard.QueryLabelByIDRequest) (resp *standard.QueryLabelByIDResponse, err error) {
	resp = new(standard.QueryLabelByIDResponse)
	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	count, err := dao.CountLabelByID(req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 { // 没有找到
		resp.State = standard.State_LABEL_NOT_EXIST
		resp.Message = "该标签不存在"
		return resp, nil
	}

	label, err := dao.QueryLabelByID(req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Data = label.OutProtoStruct()
	resp.Message = "查询成功"
	return resp, nil
}

// DeleteLabelByID 通过 ID 删除
func (srv *Service) DeleteLabelByID(ctx context.Context, req *standard.DeleteLabelByIDRequest) (resp *standard.DeleteLabelByIDResponse, err error) {
	resp = new(standard.DeleteLabelByIDResponse)
	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	count, err := dao.CountLabelByID(req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 { // 没有找到
		resp.State = standard.State_LABEL_NOT_EXIST
		resp.Message = "该标签不存在"
		return resp, nil
	}

	err = dao.DeleteLabelByID(req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "删除成功"
	return resp, nil
}

// UpdateLabelNameByID 通过 ID 更新分类
func (srv *Service) UpdateLabelNameByID(ctx context.Context, req *standard.UpdateLabelNameByIDRequest) (resp *standard.UpdateLabelNameByIDResponse, err error) {
	resp = new(standard.UpdateLabelNameByIDResponse)
	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	if ok, msg := validators.LabelName(req.Name); ok != true {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = msg
		return resp, nil
	}

	count, err := dao.CountLabelByID(req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 { // 没有找到
		resp.State = standard.State_LABEL_NOT_EXIST
		resp.Message = "该标签不存在"
		return resp, nil
	}

	err = dao.UpdateLabelNameByID(req.ID, req.Name)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "更新成功"
	return resp, nil
}

// UpdateLabelCategoryByID 通过 ID 更新分类
func (srv *Service) UpdateLabelCategoryByID(ctx context.Context, req *standard.UpdateLabelCategoryByIDRequest) (resp *standard.UpdateLabelCategoryByIDResponse, err error) {
	resp = new(standard.UpdateLabelCategoryByIDResponse)
	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	if ok, msg := validators.LabelCategory(req.Category); ok != true {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = msg
		return resp, nil
	}

	count, err := dao.CountLabelByID(req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 { // 没有找到
		resp.State = standard.State_LABEL_NOT_EXIST
		resp.Message = "该标签不存在"
		return resp, nil
	}

	err = dao.UpdateLabelCategoryByID(req.ID, req.Category)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "更新成功"
	return resp, nil
}

// UpdateLabelStateByID 通过 ID 更新状态
func (srv *Service) UpdateLabelStateByID(ctx context.Context, req *standard.UpdateLabelStateByIDRequest) (resp *standard.UpdateLabelStateByIDResponse, err error) {
	resp = new(standard.UpdateLabelStateByIDResponse)
	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	if ok, msg := validators.LabelState(req.State); ok != true {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = msg
		return resp, nil
	}

	count, err := dao.CountLabelByID(req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 { // 没有找到
		resp.State = standard.State_LABEL_NOT_EXIST
		resp.Message = "该标签不存在"
		return resp, nil
	}

	err = dao.UpdateLabelStateByID(req.ID, req.State)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "更新成功"
	return resp, nil
}

// UpdateLabelValueByID 通过 ID 更新值
func (srv *Service) UpdateLabelValueByID(ctx context.Context, req *standard.UpdateLabelValueByIDRequest) (resp *standard.UpdateLabelValueByIDResponse, err error) {
	resp = new(standard.UpdateLabelValueByIDResponse)
	if req.ID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	if ok, msg := validators.LabelValue(req.Value); ok != true {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = msg
		return resp, nil
	}

	count, err := dao.CountLabelByID(req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if count <= 0 { // 没有找到
		resp.State = standard.State_LABEL_NOT_EXIST
		resp.Message = "该标签不存在"
		return resp, nil
	}

	err = dao.UpdateLabelValueByID(req.ID, req.Value)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "更新成功"
	return resp, nil
}

// AddLabelToUserByID 添加标签到用户
func (srv *Service) AddLabelToUserByID(ctx context.Context, req *standard.AddLabelToUserByIDRequest) (resp *standard.AddLabelToUserByIDResponse, err error) {
	resp = new(standard.AddLabelToUserByIDResponse)
	if req.ID == 0 || req.UserID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	labelCount, err := dao.CountLabelByID(req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if labelCount <= 0 { // 没有找到
		resp.State = standard.State_LABEL_NOT_EXIST
		resp.Message = "该标签不存在"
		return resp, nil
	}

	userCount, err := dao.CountUserByID(req.UserID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if userCount <= 0 { // 没有找到用户
		resp.State = standard.State_USER_NOT_EXIST
		resp.Message = "该用户不存在"
		return resp, nil
	}

	already, err := dao.IsAlreadyInGroup(req.ID, req.UserID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	// 已经在组里了
	if already == true {
		resp.State = standard.State_SUCCESS
		resp.Message = "添加成功"
		return resp, nil
	}

	err = dao.AddLabelToUserByID(req.ID, req.UserID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "添加成功"
	return resp, nil
}

// RemoveLabelFromUserByID 从用户身上撕下标签
func (srv *Service) RemoveLabelFromUserByID(ctx context.Context, req *standard.RemoveLabelFromUserByIDRequest) (resp *standard.RemoveLabelFromUserByIDResponse, err error) {
	resp = new(standard.RemoveLabelFromUserByIDResponse)
	if req.ID == 0 || req.LabelID == 0 {
		resp.State = standard.State_PARAMS_INVALID
		resp.Message = "无效的 ID"
		return resp, nil
	}

	labelCount, err := dao.CountLabelByID(req.LabelID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if labelCount <= 0 { // 没有找到
		resp.State = standard.State_LABEL_NOT_EXIST
		resp.Message = "该标签不存在"
		return resp, nil
	}

	userCount, err := dao.CountUserByID(req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	if userCount <= 0 { // 没有找到用户
		resp.State = standard.State_USER_NOT_EXIST
		resp.Message = "该用户不存在"
		return resp, nil
	}

	already, err := dao.IsAlreadyOwnLabel(req.LabelID, req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	// 不存在
	if already == false {
		resp.State = standard.State_SUCCESS
		resp.Message = "移除成功"
		return resp, nil
	}

	err = dao.RemoveLabelFromUserByID(req.LabelID, req.ID)
	if err != nil {
		resp.State = standard.State_DB_OPERATION_FATLURE
		resp.Message = err.Error()
		return resp, nil
	}

	resp.State = standard.State_SUCCESS
	resp.Message = "移除成功"
	return resp, nil
}
