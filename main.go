package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vanilla/go-pagination/api/dto"

	"github.com/vanilla/go-pagination/api/services"

	"github.com/gin-gonic/gin"
	"github.com/vanilla/go-pagination/api/config"
	"github.com/vanilla/go-pagination/api/entity"
	"github.com/vanilla/go-pagination/api/payload"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.SetupDatabase()
)

func main() {
	defer config.CloseDatabase(db)

	r := gin.Default()

	contactRoutes := r.Group("api/contact")
	{
		contactRoutes.POST("/create", func(c *gin.Context) {
			var contact entity.Contact

			// Validate json
			err := c.ShouldBindJSON(&contact)

			// validation error
			if err != nil {
				res := payload.GenerateValidationResponse(err)
				c.JSON(http.StatusBadRequest, res)
				return
			}

			code := http.StatusOK

			res := services.CreateContact(&contact, db)

			if !res.Success {
				code = http.StatusBadRequest
			}

			c.JSON(code, res)
		})
		contactRoutes.GET("/", func(c *gin.Context) {
			code := http.StatusOK

			res := services.FindAllContacts(db)

			if !res.Success {
				code = http.StatusBadRequest
			}

			c.JSON(code, res)
		})
		contactRoutes.GET("/show/:id", func(c *gin.Context) {
			id := c.Param("id")

			code := http.StatusOK

			res := services.FindOneContactById(id, db)

			if !res.Success {
				code = http.StatusBadRequest
			}

			c.JSON(code, res)
		})
		contactRoutes.PUT("/update/:id", func(c *gin.Context) {
			id := c.Param("id")

			var contact entity.Contact

			err := c.ShouldBindJSON(&contact)

			// validation error
			if err != nil {
				res := payload.GenerateValidationResponse(err)
				c.JSON(http.StatusBadRequest, res)
				return
			}

			code := http.StatusOK

			res := services.UpdateContactById(id, &contact, db)

			if !res.Success {
				code = http.StatusBadRequest
			}

			c.JSON(code, res)
		})
		contactRoutes.DELETE("/delete/:id", func(c *gin.Context) {
			id := c.Param("id")

			code := http.StatusOK

			res := services.DeleteOneContactById(id, db)

			if !res.Success {
				code = http.StatusBadRequest
			}

			c.JSON(code, res)
		})
		contactRoutes.POST("/delete", func(c *gin.Context) {
			var multID dto.MultiID

			err := c.ShouldBindJSON(&multID)

			// validation error
			if err != nil {
				res := payload.GenerateValidationResponse(err)
				c.JSON(http.StatusBadRequest, res)
				return
			}

			if len(multID.Ids) == 0 {
				res := dto.Response{Success: false, Message: "IDs cannot be empty."}
				c.JSON(http.StatusBadRequest, res)
				return
			}

			code := http.StatusOK

			res := services.DeleteContactByIds(&multID, db)

			if !res.Success {
				code = http.StatusBadRequest
			}

			c.JSON(code, res)
		})
		contactRoutes.GET("/pagination", func(c *gin.Context) {
			code := http.StatusOK

			// Default limit, page & sort parameter
			limit := 10
			page := 1
			sort := "created_at desc"

			var searchs []dto.Search

			query := c.Request.URL.Query()

			for key, value := range query {
				qValue := value[len(value)-1]

				switch key {
				case "limit":
					limit, _ = strconv.Atoi(qValue)
					break
				case "page":
					page, _ = strconv.Atoi(qValue)
					break
				case "sort":
					sort = qValue
					break
				}

				if strings.Contains(key, ".") {
					// split query parameter key by dot
					searchKeys := strings.Split(key, ".")

					// create search object
					search := dto.Search{Column: searchKeys[0], Action: searchKeys[1], Query: qValue}

					// add search object
					searchs = append(searchs, search)
				}
			}

			pagination := &dto.Pagination{Limit: limit, Page: page, Sort: sort, Searchs: searchs}

			res := services.Pagination(c, pagination, db)

			if !res.Success {
				code = http.StatusBadRequest
			}

			c.JSON(code, res)
		})
	}

	r.Run(":8000")
}
