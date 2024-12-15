package mysql

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"star/app/constant/str"
	"star/app/models"
)

const (
	getAllCategorySQL       = "select  category_id, category_code, category_name, p_category_id, icon, background, sort from category_info order by sort "
	checkCategoryExistSQL   = "select  category_id from category_info  where category_id=?"
	delCategorySQL          = "delete  from  category_info where category_id=? or p_category_id=?"
	queryCategoryByCodeSQL  = "select  category_id, category_code, category_name, p_category_id, icon, background, sort from category_info where category_code=?"
	insertCategorySQL       = "insert into category_info(category_code, category_name, p_category_id, icon, background, sort) values(?,?,?,?,?,?)"
	queryMaxCategorySortSQL = "select ifnull(max(sort),0) from category_info where p_category_id=? "
	updateCategorySQL       = "update category_info  set category_name=?,category_code=?,icon=?,background=? where category_id=? "
	updateSortSQL           = "update category_info  set sort=? where category_id=? and  p_category_id=?"
)

func GetAllCategory() ([]*models.Category, error) {
	var categoryList []*models.Category
	if err := Client.Select(&categoryList, getAllCategorySQL); err != nil {
		return nil, err
	}
	return categoryList, nil
}

func CheckCategoryExist(categoryId int64) error {
	var id int
	if err := Client.Get(&id, checkCategoryExistSQL, categoryId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return str.ErrCategoryNotExists
		}
		return err
	}
	return nil
}

func DelCategory(categoryId int64) error {
	if _, err := Client.Exec(delCategorySQL, categoryId, categoryId); err != nil {
		return err
	}
	return nil
}

func QueryCategoryByCode(categoryCode string) (*models.Category, error) {
	category := new(models.Category)
	if err := Client.Get(category, queryCategoryByCodeSQL, categoryCode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, str.ErrCategoryNotExists
		}
		return nil, err
	}
	return category, nil
}

func QueryMaxCategorySort(pCategoryId int64) (uint32, error) {
	var sort uint32
	if err := Client.Get(&sort, queryMaxCategorySortSQL, pCategoryId); err != nil {
		return 0, err
	}
	return sort, nil
}

func InsertCategory(category *models.Category) error {
	if _, err := Client.Exec(insertCategorySQL, category.CategoryCode,
		category.CategoryName, category.PCategoryId, category.Icon, category.Background,
		category.Sort); err != nil {
		return err
	}
	return nil
}

func UpdateCategory(category *models.Category) error {
	if _, err := Client.Exec(updateCategorySQL, category.CategoryName,
		category.CategoryCode, category.Icon, category.Background,
		category.CategoryId); err != nil {
		return err
	}
	return nil
}

func BatchUpdateSort(categorys []*models.Category) (err error) {
	var tx *sqlx.Tx
	tx, err = Client.Beginx()
	if err != nil {
		return str.ErrMessageError
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		}
	}()
	for _, category := range categorys {
		_, err = tx.Exec(updateSortSQL, category.Sort, category.CategoryId, category.PCategoryId)
		if err != nil {
			return
		}
	}
	err = tx.Commit()
	return
}
