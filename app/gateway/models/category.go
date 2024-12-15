package models

type Category struct {
	CategoryId   int64  `json:"categoryId,omitempty"  form:"categoryId"`
	PCategoryId  int64  `json:"pCategoryId" form:"pCategoryId"`
	CategoryCode string `json:"categoryCode" binding:"required" form:"categoryCode"`
	CategoryName string `json:"categoryName" binding:"required" form:"categoryName"`
	Icon         string `json:"icon,omitempty" form:"icon"`
	Background   string `json:"background,omitempty" form:"background"`
}

type ChangeSort struct {
	PCategoryId int64  `form:"pCategoryId"`
	CategoryIds string `form:"categoryIds"`
}
