package services

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vanilla/go-pagination/api/dto"
	"github.com/vanilla/go-pagination/api/entity"
	"github.com/vanilla/go-pagination/api/repository"
	"gorm.io/gorm"
)

func CreateContact(contact *entity.Contact, db *gorm.DB) dto.Response {
	uuid, err := uuid.NewRandom()

	if err != nil {
		log.Fatalln(err)
	}

	contact.ID = uuid.String()
	repo := repository.NewContactRepository(db)

	operation, err := repo.Save(contact)

	if err != nil {
		return dto.Response{Success: false, Message: err.Error()}
	}

	return dto.Response{Success: true, Data: operation}
}

func FindAllContacts(db *gorm.DB) dto.Response {
	repo := repository.NewContactRepository(db)

	operation, err := repo.FindAll()

	if err != nil {
		return dto.Response{Success: false, Message: err.Error()}
	}

	return dto.Response{Success: true, Data: operation}
}

func FindOneContactById(id string, db *gorm.DB) dto.Response {
	repo := repository.NewContactRepository(db)

	operation, err := repo.FindOneById(id)

	if err != nil {
		return dto.Response{Success: false, Message: err.Error()}
	}

	return dto.Response{Success: true, Data: operation}
}

func UpdateContactById(id string, contact *entity.Contact, db *gorm.DB) dto.Response {
	existing := FindOneContactById(id, db)

	if !existing.Success {
		return existing
	}

	fmt.Println(existing)

	existingContact := existing.Data.(entity.Contact)

	existingContact.Name = contact.Name
	existingContact.Email = contact.Email
	existingContact.Address = contact.Address

	repo := repository.NewContactRepository(db)
	operation, err := repo.Save(&existingContact)

	if err != nil {
		return dto.Response{Success: false, Message: err.Error()}
	}

	return dto.Response{Success: true, Data: operation}
}

func DeleteOneContactById(id string, db *gorm.DB) dto.Response {
	repo := repository.NewContactRepository(db)
	operation, err := repo.DeleteOneById(id)

	if err != nil {
		return dto.Response{Success: false, Message: err.Error()}
	}

	return dto.Response{Success: true, Data: operation}
}

func DeleteContactByIds(multiId *dto.MultiID, db *gorm.DB) dto.Response {
	repo := repository.NewContactRepository(db)
	operation, err := repo.DeleteByIds(&multiId.Ids)

	if err != nil {
		return dto.Response{Success: false, Message: err.Error()}
	}

	return dto.Response{Success: true, Data: operation}
}

func Pagination(c *gin.Context, pagination *dto.Pagination, db *gorm.DB) dto.Response {
	repo := repository.NewContactRepository(db)
	operation, err, totalPages := repo.Pagination(pagination)

	if err != nil {
		return dto.Response{Success: false, Message: err.Error()}
	}

	var data = operation.(*dto.Pagination)

	// get current url path
	urlPath := c.Request.URL.Path

	// search query params
	searchQueryParams := ""

	for _, search := range pagination.Searchs {
		searchQueryParams += fmt.Sprintf("&%s.%s=%s", search.Column, search.Action, search.Query)
	}

	// set first & last page pagination response
	data.FirstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, 0, pagination.Sort) + searchQueryParams
	data.LastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, totalPages, pagination.Sort) + searchQueryParams

	if data.Page > 0 {
		// set previous page pagination response
		data.PreviousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, data.Page-1, pagination.Sort) + searchQueryParams
	}

	if data.Page < totalPages {
		// set next page pagination response
		data.NextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s", urlPath, pagination.Limit, data.Page+1, pagination.Sort) + searchQueryParams
	}

	if data.Page > totalPages {
		// reset previous page
		data.PreviousPage = ""
	}

	return dto.Response{Success: true, Data: operation}
}
