package models

type Category struct {
	CategoryCode string `db:"category_code" json:"category_code,omitempty"` //分类编码
	CategoryName string `db:"category_name" json:"category_name,omitempty"` //分类名称
	Icon         string `db:"icon" json:"icon,omitempty"`                   //图标
	Background   string `db:"background" json:"background,omitempty"`       //背景
	CategoryId   int64  `db:"category_id" json:"category_id,omitempty"`
	PCategoryId  int64  `db:"p_category_id" json:"p_category_id,omitempty"` //父分类id
	Sort         uint32 `db:"sort" json:"sort,omitempty"`                   //排序号
}
