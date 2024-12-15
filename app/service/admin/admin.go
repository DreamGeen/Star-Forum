package main

import (
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/models"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/app/utils/logging"
	"star/proto/admin/adminPb"
	"strconv"
	"strings"
)

type AdminSrv struct {
}

var adminIns = new(AdminSrv)

func (a *AdminSrv) LoadCategoryList(ctx context.Context, req *adminPb.LoadCategoryListRequest, resp *adminPb.LoadCategoryListResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "LoadCategoryListService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "AdminService.LoadCategoryList")

	categoryListJson, err := redis.Client.Get(ctx, str.Redis_Key_Category).Result()
	if err != nil {
		logger.Error("mysql get all adminPb error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrCategoryError
	}
	if err := json.Unmarshal([]byte(categoryListJson), &resp.CategoryList); err != nil {
		logger.Error("json unmarshal categoryList error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrCategoryError
	}
	return nil
}

func (a *AdminSrv) DelCategory(ctx context.Context, req *adminPb.DelCategoryRequest, resp *adminPb.DelCategoryResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "DelCategoryService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "AdminService.DelCategory")

	if err := mysql.CheckCategoryExist(req.CategoryId); err != nil {
		logging.SetSpanError(span, err)
		if !errors.Is(err, str.ErrCategoryNotExists) {
			logger.Error("check adminPb is exist error",
				zap.Error(err))
			return str.ErrCategoryError
		}
		return err
	}
	if err := mysql.DelCategory(req.CategoryId); err != nil {
		logger.Error("del adminPb error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrCategoryError
	}
	updateCategoryCache(ctx, logger)
	return nil
}

func categoryList2Tree(categoryList []*models.Category) []*adminPb.Category {

	categoryListMap := make(map[int64]*adminPb.Category)
	for _, category := range categoryList {
		categoryProto := &adminPb.Category{
			CategoryId:   category.CategoryId,
			CategoryName: category.CategoryName,
			CategoryCode: category.CategoryCode,
			PCategoryId:  category.PCategoryId,
			Icon:         category.Icon,
			Background:   category.Background,
			Sort:         category.Sort,
		}
		categoryListMap[category.CategoryId] = categoryProto
	}

	var categoryTree []*adminPb.Category
	for _, category := range categoryList {
		if category.PCategoryId == 0 {
			categoryTree = append(categoryTree, categoryListMap[category.CategoryId])
		} else if pCategory, ok := categoryListMap[category.PCategoryId]; ok {
			pCategory.Children = append(pCategory.Children, categoryListMap[category.CategoryId])
		}
	}
	return categoryTree
}

func (a *AdminSrv) SaveCategory(ctx context.Context, req *adminPb.SaveCategoryRequest, resp *adminPb.SaveCategoryResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "SaveCategoryService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "AdminService.SaveCategory")

	category, err := mysql.QueryCategoryByCode(req.CategoryCode)
	if err != nil && !errors.Is(err, str.ErrCategoryNotExists) {
		logger.Error("mysql query adminPb by name  error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrCategoryError
	}

	if req.CategoryId == 0 && category != nil ||
		req.CategoryId != 0 && category != nil && req.CategoryId != category.CategoryId {
		return str.ErrCategoryIdExists
	}
	saveCategory := &models.Category{
		CategoryId:   req.CategoryId,
		CategoryCode: req.CategoryCode,
		CategoryName: req.CategoryName,
		PCategoryId:  req.PCategoryId,
		Icon:         req.Icon,
		Background:   req.Background,
	}
	logger.Debug("adminPb", zap.Any("x", saveCategory))
	if category == nil {
		sort, err := mysql.QueryMaxCategorySort(saveCategory.PCategoryId)
		if err != nil {
			logger.Error("mysql query max adminPb sort error",
				zap.Error(err))
			logging.SetSpanError(span, err)
			return str.ErrCategoryError
		}
		saveCategory.Sort = sort + 1
		err = mysql.InsertCategory(saveCategory)
		if err != nil {
			logger.Error("mysql insert adminPb error",
				zap.Error(err))
			logging.SetSpanError(span, err)
			return str.ErrCategoryError
		}
	} else {
		err := mysql.UpdateCategory(saveCategory)
		if err != nil {
			logger.Error("mysql update adminPb error",
				zap.Error(err))
			logging.SetSpanError(span, err)
			return str.ErrCategoryError
		}
	}
	updateCategoryCache(ctx, logger)
	return nil
}

func (a *AdminSrv) ChangeSort(ctx context.Context, req *adminPb.ChangeSortRequest, resp *adminPb.ChangeSortResponse) error {
	ctx, span := tracing.Tracer.Start(ctx, "ChangeSortService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "AdminService.ChangeSort")

	categoryIdsStr := strings.Split(req.CategoryIdsStr, ",")
	categorys := make([]*models.Category, len(categoryIdsStr))
	sort := uint32(1)
	for i, idStr := range categoryIdsStr {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		category := &models.Category{
			CategoryId:  id,
			PCategoryId: req.PCategoryId,
			Sort:        sort,
		}
		categorys[i] = category
		sort++
	}
	if err := mysql.BatchUpdateSort(categorys); err != nil {
		logger.Error("mysql batch update sort error",
			zap.Error(err))
		logging.SetSpanError(span, err)
		return str.ErrCategoryError
	}
	err := updateCategoryCache(ctx, logger)
	if err != nil {
		logger.Warn("update cache error",
			zap.Error(err))
	}
	return nil
}

func updateCategoryCache(ctx context.Context, logger *zap.Logger) error {
	categorys, err := mysql.GetAllCategory()
	if err != nil {
		logger.Error("mysql get all category error",
			zap.Error(err))
		return str.ErrCategoryError
	}
	categoryList := categoryList2Tree(categorys)
	categoryListJson, err := json.Marshal(categoryList)
	if err != nil {
		return err
	}
	return redis.Client.Set(ctx, str.Redis_Key_Category, categoryListJson, 0).Err()
}
