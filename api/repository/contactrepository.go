package repository

import (
	"fmt"
	"math"
	"strings"

	"github.com/vanilla/go-pagination/api/dto"
	"github.com/vanilla/go-pagination/api/entity"
	"gorm.io/gorm"
)

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) Save(contact *entity.Contact) (interface{}, error) {
	err := r.db.Save(contact).Error

	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (r *ContactRepository) FindAll() (interface{}, error) {
	var contacts entity.Contacts

	err := r.db.Find(&contacts).Error

	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (r *ContactRepository) FindOneById(id string) (interface{}, error) {
	var contact entity.Contact

	err := r.db.Where(&entity.Contact{ID: id}).Take(&contact).Error

	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (r *ContactRepository) DeleteOneById(id string) (interface{}, error) {
	err := r.db.Delete(&entity.Contact{ID: id}).Error

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *ContactRepository) DeleteByIds(ids *[]string) (interface{}, error) {
	err := r.db.Where("id IN(?)", *ids).Delete(&entity.Contact{}).Error

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *ContactRepository) Pagination(pagination *dto.Pagination) (interface{}, error, int) {
	var contacts entity.Contacts

	totalRows, totalPages, fromRow, toRow := 0, 0, 0, 0

	offset := pagination.Page * pagination.Limit

	// get data with limit, offset & order
	find := r.db.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)

	// generate where query
	searchs := pagination.Searchs

	if searchs != nil {
		for _, value := range searchs {
			column := value.Column
			action := value.Action
			query := value.Query

			switch action {
			case "equals":
				whereQuery := fmt.Sprintf("%s = ?", column)
				find = find.Where(whereQuery, query)
				break
			case "contains":
				whereQuery := fmt.Sprintf("%s LIKE ?", column)
				find = find.Where(whereQuery, "%"+query+"%")
				break
			case "in":
				whereQuery := fmt.Sprintf("%s IN (?)", column)
				queryArray := strings.Split(query, ",")
				find = find.Where(whereQuery, queryArray)
				break

			}
		}
	}

	find = find.Find(&contacts)

	// has error find data
	errFind := find.Error

	if errFind != nil {
		return nil, errFind, totalPages
	}

	pagination.Rows = contacts

	counting := int64(totalRows)
	// count all data
	errCount := r.db.Model(&entity.Contact{}).Count(&counting).Error

	if errCount != nil {
		return nil, errCount, totalPages
	}

	totalRows = int(counting)

	pagination.TotalRows = totalRows

	// calculate total pages
	totalPages = int(math.Ceil(float64(totalRows)/float64(pagination.Limit))) - 1

	if pagination.Page == 0 {
		// set from & to row on first page
		fromRow = 1
		toRow = pagination.Limit
	} else {
		if pagination.Page <= totalPages {
			// calculate from & to row
			fromRow = pagination.Page*pagination.Limit + 1
			toRow = (pagination.Page + 1) * pagination.Limit
		}
	}

	if toRow > totalRows {
		// set to row with total rows
		toRow = totalRows
	}

	pagination.FromRow = fromRow
	pagination.ToRow = toRow

	return pagination, nil, totalPages
}
