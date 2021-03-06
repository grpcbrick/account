package provider

import (
	"context"
	"testing"

	"github.com/grpcbrick/account/standard"
)

func TestService_CreateLabel(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.CreateLabelRequest
		wantState standard.State
		wantErr   bool
	}{
		{"正常创建", &standard.CreateLabelRequest{Name: "TEST", Category: "Category", State: "State", Value: "Value"},
			standard.State_SUCCESS, false},
		{"正常创建", &standard.CreateLabelRequest{Name: "TEST2", Category: "Category", State: "State", Value: "Value"},
			standard.State_SUCCESS, false},
		{"正常创建", &standard.CreateLabelRequest{Name: "TEST3", Category: "Category", State: "State", Value: "Value"},
			standard.State_SUCCESS, false},
		{"重复的 Name", &standard.CreateLabelRequest{Name: "TEST", Category: "Category", State: "State", Value: "Value"},
			standard.State_SUCCESS, false}, // 标签允许重复
		{"空的 Name", &standard.CreateLabelRequest{Name: "", Category: "Category", State: "Nickname", Value: "Value"},
			standard.State_PARAMS_INVALID, false},
		{"空的 Category", &standard.CreateLabelRequest{Name: "TEST", Category: "", State: "Nickname", Value: "Value"},
			standard.State_PARAMS_INVALID, false},
		{"空的 State", &standard.CreateLabelRequest{Name: "TEST", Category: "Category", State: "", Value: "Value"},
			standard.State_PARAMS_INVALID, false},
		{"空的 Value", &standard.CreateLabelRequest{Name: "TEST", Category: "Category", State: "Nickname", Value: ""},
			standard.State_PARAMS_INVALID, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.CreateLabel(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateLabel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.CreateLabel() = %v, want %v", gotResp, tt.wantState)
				return
			}
			if gotResp.State == standard.State_SUCCESS {
				if gotResp.Data.Name != tt.args.Name {
					t.Errorf("Service.CreateLabel() = %v, want %v", gotResp, tt.wantState)
					return
				}
			}

		})
	}
}

func TestService_CreateLabelForUser(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.CreateLabelForUserRequest
		wantState standard.State
		wantErr   bool
	}{
		{"正常创建", &standard.CreateLabelForUserRequest{UserID: 1, Name: "TEST1", Category: "Category", State: "State", Value: "Value"},
			standard.State_SUCCESS, false},
		{"正常创建", &standard.CreateLabelForUserRequest{UserID: 1, Name: "TEST2", Category: "Category", State: "State", Value: "Value"},
			standard.State_SUCCESS, false},
		{"正常创建", &standard.CreateLabelForUserRequest{UserID: 1, Name: "TEST3", Category: "Category", State: "State", Value: "Value"},
			standard.State_SUCCESS, false},
		{"重复的 Name", &standard.CreateLabelForUserRequest{UserID: 1, Name: "TEST4", Category: "Category", State: "State", Value: "Value"},
			standard.State_SUCCESS, false}, // 标签允许重复
		{"空的 Name", &standard.CreateLabelForUserRequest{UserID: 1, Name: "5", Category: "Category", State: "Nickname", Value: "Value"},
			standard.State_PARAMS_INVALID, false},
		{"空的 Category", &standard.CreateLabelForUserRequest{UserID: 1, Name: "TEST6", Category: "", State: "Nickname", Value: "Value"},
			standard.State_PARAMS_INVALID, false},
		{"空的 State", &standard.CreateLabelForUserRequest{UserID: 1, Name: "TEST7", Category: "Category", State: "", Value: "Value"},
			standard.State_PARAMS_INVALID, false},
		{"空的 Value", &standard.CreateLabelForUserRequest{UserID: 1, Name: "TEST8", Category: "Category", State: "Nickname", Value: ""},
			standard.State_PARAMS_INVALID, false},
		{"空的 UserID", &standard.CreateLabelForUserRequest{Name: "TEST9", Category: "Category", State: "Nickname", Value: ""},
			standard.State_USER_NOT_EXIST, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.CreateLabelForUser(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateLabelForUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.CreateLabelForUser() = %v, want %v", gotResp, tt.wantState)
				return
			}

			if gotResp.State == standard.State_SUCCESS {
				if gotResp.Data.Name != tt.args.Name {
					t.Errorf("Service.CreateLabelForUser() = %v, want %v", gotResp, tt.wantState)
					return
				}
			}
		})
	}
}

func TestService_QueryLabelByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.QueryLabelByIDRequest
		wantState standard.State
		wantName  string
		wantErr   bool
	}{
		{"正常查询", &standard.QueryLabelByIDRequest{ID: 1},
			standard.State_SUCCESS, "TEST", false},
		{"空的 ID", &standard.QueryLabelByIDRequest{},
			standard.State_PARAMS_INVALID, "ignore", false},
		{"不存在的 ID", &standard.QueryLabelByIDRequest{ID: 9999999},
			standard.State_LABEL_NOT_EXIST, "ignore", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.QueryLabelByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.QueryLabelByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.QueryLabelByID() = %v, want %v", gotResp, tt.wantState)
				return
			}

			if tt.wantName != "ignore" {
				if gotResp.Data.Name != tt.wantName {
					t.Errorf("Service.QueryLabelByID() = %v, want %v", gotResp.Data.Name, tt.wantName)
				}
			}
		})
	}
}

func TestService_DeleteLabelByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.DeleteLabelByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的 ID", &standard.DeleteLabelByIDRequest{},
			standard.State_PARAMS_INVALID, false},
		{"正常删除", &standard.DeleteLabelByIDRequest{ID: 1},
			standard.State_SUCCESS, false},
		{"不存在的 ID", &standard.DeleteLabelByIDRequest{ID: 9999999},
			standard.State_LABEL_NOT_EXIST, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.DeleteLabelByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.DeleteLabelByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.DeleteLabelByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_UpdateLabelNameByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.UpdateLabelNameByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的 ID", &standard.UpdateLabelNameByIDRequest{Name: "Update1"},
			standard.State_PARAMS_INVALID, false},
		{"正常更新", &standard.UpdateLabelNameByIDRequest{ID: 2, Name: "Update1"},
			standard.State_SUCCESS, false},
		{"不存在的 ID", &standard.UpdateLabelNameByIDRequest{ID: 99999, Name: "Update1"},
			standard.State_LABEL_NOT_EXIST, false},
		{"空的 Name", &standard.UpdateLabelNameByIDRequest{ID: 2},
			standard.State_PARAMS_INVALID, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.UpdateLabelNameByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateLabelNameByIDRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.UpdateLabelNameByIDRequest() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_UpdateLabelCategoryByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.UpdateLabelCategoryByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的 ID", &standard.UpdateLabelCategoryByIDRequest{Category: "Update1"},
			standard.State_PARAMS_INVALID, false},
		{"正常更新", &standard.UpdateLabelCategoryByIDRequest{ID: 2, Category: "Update1"},
			standard.State_SUCCESS, false},
		{"不存在的 ID", &standard.UpdateLabelCategoryByIDRequest{ID: 99999, Category: "Update1"},
			standard.State_LABEL_NOT_EXIST, false},
		{"空的 Name", &standard.UpdateLabelCategoryByIDRequest{ID: 2},
			standard.State_PARAMS_INVALID, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.UpdateLabelCategoryByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateLabelCategoryByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.UpdateLabelCategoryByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_UpdateLabelStateByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.UpdateLabelStateByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的 ID", &standard.UpdateLabelStateByIDRequest{State: "Update1"},
			standard.State_PARAMS_INVALID, false},
		{"正常更新", &standard.UpdateLabelStateByIDRequest{ID: 2, State: "Update1"},
			standard.State_SUCCESS, false},
		{"不存在的 ID", &standard.UpdateLabelStateByIDRequest{ID: 99999, State: "Update1"},
			standard.State_LABEL_NOT_EXIST, false},
		{"空的 Name", &standard.UpdateLabelStateByIDRequest{ID: 2},
			standard.State_PARAMS_INVALID, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.UpdateLabelStateByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateLabelStateByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.UpdateLabelStateByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_UpdateLabelValueByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.UpdateLabelValueByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的 ID", &standard.UpdateLabelValueByIDRequest{Value: "Update1"},
			standard.State_PARAMS_INVALID, false},
		{"正常更新", &standard.UpdateLabelValueByIDRequest{ID: 2, Value: "Update1"},
			standard.State_SUCCESS, false},
		{"不存在的 ID", &standard.UpdateLabelValueByIDRequest{ID: 99999, Value: "Update1"},
			standard.State_LABEL_NOT_EXIST, false},
		{"空的 Name", &standard.UpdateLabelValueByIDRequest{ID: 2},
			standard.State_PARAMS_INVALID, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.UpdateLabelValueByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UpdateLabelValueByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.UpdateLabelValueByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_QueryLabels(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name            string
		args            *standard.QueryLabelsRequest
		wantState       standard.State
		wantDataSize    int64
		wantCurrentPage int64
		wantTotalPage   int64
		wantErr         bool
	}{
		{"正常查询", &standard.QueryLabelsRequest{Page: 1, Limit: 90},
			standard.State_SUCCESS, 8, 1, 1, false},
		{"只查一条", &standard.QueryLabelsRequest{Page: 1, Limit: 1},
			standard.State_SUCCESS, 1, 1, 8, false},
		{"第二页", &standard.QueryLabelsRequest{Page: 1, Limit: 2},
			standard.State_SUCCESS, 2, 1, 4, false},
		{"空的 ID", &standard.QueryLabelsRequest{Page: 1, Limit: 0},
			standard.State_PARAMS_INVALID, 0, 1, 0, false},
		{"不存在的 ID", &standard.QueryLabelsRequest{},
			standard.State_PARAMS_INVALID, 0, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.QueryLabels(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.QueryLabels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.QueryLabels() = %v, want %v", gotResp, tt.wantState)
				return
			}

			if tt.wantState == standard.State_SUCCESS {
				if gotResp.TotalPage != tt.wantTotalPage {
					t.Errorf("Service.QueryLabels() = %v, want %d", gotResp, tt.wantTotalPage)
					return
				}
				if int64(len(gotResp.Data)) != tt.wantDataSize {
					t.Errorf("Service.QueryLabels() = %v, want %d", gotResp.Data, tt.wantDataSize)
					return
				}
				if gotResp.CurrentPage != tt.wantCurrentPage {
					t.Errorf("Service.QueryLabels() = %v, want %d", gotResp, tt.wantCurrentPage)
					return
				}
			}
		})
	}
}

func TestService_AddLabelToUserByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.AddLabelToUserByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的用户 ID", &standard.AddLabelToUserByIDRequest{ID: 2},
			standard.State_PARAMS_INVALID, false},
		{"不存在的用户 ID", &standard.AddLabelToUserByIDRequest{ID: 2, UserID: 999999},
			standard.State_USER_NOT_EXIST, false},
		{"空的标签 ID", &standard.AddLabelToUserByIDRequest{UserID: 999999},
			standard.State_PARAMS_INVALID, false},
		{"不存在的标签 ID", &standard.AddLabelToUserByIDRequest{ID: 999999, UserID: 1},
			standard.State_LABEL_NOT_EXIST, false},
		{"正常添加", &standard.AddLabelToUserByIDRequest{ID: 1, UserID: 1},
			standard.State_SUCCESS, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.AddLabelToUserByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddLabelToUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.AddLabelToUserByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}

func TestService_RemoveLabelFromUserByID(t *testing.T) {
	srv := NewService()
	tests := []struct {
		name      string
		args      *standard.RemoveLabelFromUserByIDRequest
		wantState standard.State
		wantErr   bool
	}{
		{"空的用户 ID", &standard.RemoveLabelFromUserByIDRequest{LabelID: 2},
			standard.State_PARAMS_INVALID, false},
		{"不存在的用户 ID", &standard.RemoveLabelFromUserByIDRequest{LabelID: 2, ID: 2},
			standard.State_USER_NOT_EXIST, false},
		{"空的标签 ID", &standard.RemoveLabelFromUserByIDRequest{ID: 2},
			standard.State_PARAMS_INVALID, false},
		{"不存在的标签 ID", &standard.RemoveLabelFromUserByIDRequest{LabelID: 9999, ID: 2},
			standard.State_LABEL_NOT_EXIST, false},
		{"正常添加", &standard.RemoveLabelFromUserByIDRequest{LabelID: 1, ID: 1},
			standard.State_SUCCESS, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := srv.RemoveLabelFromUserByID(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RemoveLabelFromUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotResp.State.String() != tt.wantState.String() {
				t.Errorf("Service.RemoveLabelFromUserByID() = %v, want %v", gotResp, tt.wantState)
				return
			}
		})
	}
}
