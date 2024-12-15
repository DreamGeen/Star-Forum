package httpHandler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"star/app/constant/settings"
	"star/app/constant/str"
	"star/app/extra/tracing"
	"star/app/gateway/client"
	"star/app/gateway/models"
	models2 "star/app/models"
	"star/app/utils/jwt"
	"star/app/utils/logging"
	"star/proto/admin/adminPb"
)

func LoginAdminHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "LoginAdminHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.LoginAdmin")

	//参数校验
	u := new(models.LoginPassword)
	if err := c.ShouldBindJSON(u); err != nil {
		logger.Error("adminPb login error invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}

	//图形验证码效验
	if !store.Verify(u.CheckCodeKey, u.CheckCode, true) {
		logger.Warn("img captcha error",
			zap.String("user", u.User))
		str.Response(c, str.ErrInvalidImgCaptcha, nil)
		return
	}
	if u.User != settings.Conf.Admin.Username || u.Password != settings.Conf.Admin.Password {
		logger.Error("adminPb login error because invalid password")
		str.Response(c, str.ErrInvalidPassword, nil)
		return
	}
	token, _, err := jwt.GetToken(&models2.User{UserId: settings.Conf.Admin.Id})
	if err != nil {
		logger.Error("adminPb get token error",
			zap.Error(err))
		str.Response(c, err, nil)
	}
	str.Response(c, nil, map[string]interface{}{
		"token": token,
	})
}

func LoadCategoryListHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "LoadCategoryListHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.LoadCategoryList")

	resp, err := client.LoadCategoryList(c.Request.Context(), &adminPb.LoadCategoryListRequest{})
	if err != nil {
		logger.Error("load adminPb list error",
			zap.Error(err))
		return
	}
	str.Response(c, nil, map[string]interface{}{
		"data": resp.CategoryList,
	})
}

func DelCategoryHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "DelCategoryHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.DelCategory")

	var body struct {
		Id int64 `form:"categoryId"`
	}
	if err := c.ShouldBind(&body); err != nil {
		logger.Error("del adminPb error because  invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	_, err := client.DelCategory(c.Request.Context(), &adminPb.DelCategoryRequest{
		CategoryId: body.Id,
	})
	if err != nil {
		logger.Error("del adminPb error",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, nil)
}

func SaveCategoryHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "SaveCategoryHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.SaveCategory")

	category := new(models.Category)
	if err := c.ShouldBind(category); err != nil {
		logger.Error("save adminPb error because invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	if _, err := client.SaveCategory(c.Request.Context(), &adminPb.SaveCategoryRequest{
		CategoryId:   category.CategoryId,
		CategoryName: category.CategoryName,
		CategoryCode: category.CategoryCode,
		Icon:         category.Icon,
		PCategoryId:  category.PCategoryId,
		Background:   category.Background,
	}); err != nil {
		logger.Error("save adminPb error",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}
	str.Response(c, nil, nil)
}

func ChangeSortHandler(c *gin.Context) {
	_, span := tracing.Tracer.Start(c.Request.Context(), "ChangeSortHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogServiceWithTrace(span, "GateWay.ChangeSort")

	change := new(models.ChangeSort)
	if err := c.ShouldBind(change); err != nil {
		logger.Error("change adminPb sort error because  invalid param",
			zap.Error(err))
		str.Response(c, str.ErrInvalidParam, nil)
		return
	}
	_, err := client.ChangeSort(c.Request.Context(), &adminPb.ChangeSortRequest{
		PCategoryId:    change.PCategoryId,
		CategoryIdsStr: change.CategoryIds,
	})
	if err != nil {
		logger.Error("change sort  error",
			zap.Error(err))
		str.Response(c, err, nil)
		return
	}

	str.Response(c, nil, nil)
}
