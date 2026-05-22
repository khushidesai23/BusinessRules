package router

import (
	// "calculationengine/service"
	"calculationengine/models"
	"calculationengine/service/attribute"
	"calculationengine/service/category"
	"calculationengine/service/formulas"
	"calculationengine/service/product"
	"fmt"
	"github.com/gin-contrib/cors"

	// storage "calculationengine/store"
	// "fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Router *gin.Engine

func parseRequest[T any](c *gin.Context) T {
	var request T
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err)
		log.Fatalln("Failed to Parse Request")
	}
	return request
}

func Api() {
	Router = gin.Default()
	Router.Use(cors.Default())

	Router.POST("/v1/attribute/create", func(c *gin.Context) {
		request := parseRequest[models.CreateAttributeRequest](c)
		validate := validator.New()
		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(400, gin.H{"error": validationErr.Error()})
			return
		}
		ctx := c.Request.Context()
		result, err := attribute.CreateAttribute(ctx, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, result)
	})

	Router.POST("/v1/formula/create", func(c *gin.Context) {
		request := parseRequest[models.CreateFormulaRequest](c)
		validate := validator.New()
		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(400, gin.H{"error": validationErr.Error()})
			return
		}
		ctx := c.Request.Context()
		result, err := formulas.CreateFormula(ctx, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, result)
	})

	Router.POST("/v1/formula/get-all", func(c *gin.Context) {
		ctx := c.Request.Context()
		result, err := formulas.GetFormulasList(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, result)
	})

	Router.POST("/v1/category/get-all", func(c *gin.Context) {
		ctx := c.Request.Context()
		result, err := category.GetAllCategories(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, result)
	})

	Router.POST("/v1/category/create", func(c *gin.Context) {
		request := parseRequest[models.CreateCategoryRequest](c)
		validate := validator.New()
		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(400, gin.H{"error": validationErr.Error()})
			return
		}
		ctx := c.Request.Context()
		result, err := category.CreateCategory(ctx, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, result)
	})

	Router.POST("/v1/attributes/get-all", func(c *gin.Context) {
		ctx := c.Request.Context()
		result, err := attribute.GetAllAttributes(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, result)
	})

	Router.POST("/v1/assignment/get-category-wise-common-attributes", func(c *gin.Context) {
		request := parseRequest[models.GetCategoryWiseCommonAttributesRequest](c)
		validate := validator.New()
		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(400, gin.H{"error": validationErr.Error()})
			return
		}
		ctx := c.Request.Context()
		result, err := attribute.GetCategoryWiseCommonAttributes(ctx, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, result)
	})

	Router.POST("/v1/assignment/change", func(c *gin.Context) {
		request := parseRequest[models.ChangeCategoryAttributeAssignmentRequest](c)
		validate := validator.New()
		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(400, gin.H{"error": validationErr.Error()})
			return
		}
		ctx := c.Request.Context()
		result, err := attribute.ChangeCategoryAttributeAssignment(ctx, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, result)
	})

	Router.POST("/v1/product/upsert", func(c *gin.Context) {
		request := parseRequest[models.CreateProductRequest](c)
		validate := validator.New()
		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(400, gin.H{"error": validationErr.Error()})
			return
		}
		ctx := c.Request.Context()
		result, err := product.UpsertProduct(ctx, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, result)
	})

	Router.POST("/v1/product/get-all", func(c *gin.Context) {
		ctx := c.Request.Context()
		result, err := product.GetProductList(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, result)
	})

	Router.POST("/v1/product/get-by-id", func(c *gin.Context) {
		request := parseRequest[models.GetProductDataRequest](c)
		validate := validator.New()
		validationErr := validate.Struct(request)
		if validationErr != nil {
			c.JSON(400, gin.H{"error": validationErr.Error()})
			return
		}
		ctx := c.Request.Context()
		result, err := product.GetSingleProductData(ctx, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, result)
	})

	Router.GET("/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "healthy",
		})
	})
}
