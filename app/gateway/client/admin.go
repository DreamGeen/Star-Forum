package client

import (
	"context"
	"star/proto/admin/adminPb"
)

func LoadCategoryList(ctx context.Context, req *adminPb.LoadCategoryListRequest) (*adminPb.LoadCategoryListResponse, error) {
	return adminService.LoadCategoryList(ctx, req)
}

func DelCategory(ctx context.Context, req *adminPb.DelCategoryRequest) (*adminPb.DelCategoryResponse, error) {
	return adminService.DelCategory(ctx, req)
}

func SaveCategory(ctx context.Context, req *adminPb.SaveCategoryRequest) (*adminPb.SaveCategoryResponse, error) {
	return adminService.SaveCategory(ctx, req)
}

func ChangeSort(ctx context.Context, req *adminPb.ChangeSortRequest) (*adminPb.ChangeSortResponse, error) {
	return adminService.ChangeSort(ctx, req)
}
