package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saufiroja/fin-ai/internal/contracts/requests"
	"github.com/saufiroja/fin-ai/internal/contracts/responses"
	"github.com/saufiroja/fin-ai/internal/domains/categories"
)

type categoryController struct {
	categoryService categories.CategoryManager
}

func NewCategoryController(categoryService categories.CategoryManager) categories.CategoryController {
	return &categoryController{
		categoryService: categoryService,
	}
}

func (cc *categoryController) CreateCategory(c *fiber.Ctx) error {
	req := &requests.CategoryRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	err := cc.categoryService.CreateCategory(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(responses.Response{
		Status:  fiber.StatusCreated,
		Message: "Category created successfully",
		Data:    nil,
	})
}

func (cc *categoryController) GetAllCategories(c *fiber.Ctx) error {
	req := &requests.GetAllCategoryQuery{
		Limit:  10, // Default limit
		Offset: 1,  // Default offset
		Search: "", // Default search term
	}
	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid query parameters",
		})
	}

	categories, err := cc.categoryService.FindAllCategories(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to retrieve categories",
		})
	}

	return c.JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Categories retrieved successfully",
		Data:    categories.Categories,
		Pagination: &responses.Pagination{
			CurrentPage: categories.CurrentPage,
			TotalPages:  categories.TotalPages,
			Total:       categories.Total,
		},
	})
}

func (cc *categoryController) UpdateCategoryById(c *fiber.Ctx) error {
	categoryId := c.Params("category_id")
	if categoryId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Category ID is required",
		})
	}

	req := &requests.UpdateCategoryRequest{}
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	err := cc.categoryService.UpdateCategoryById(categoryId, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to update category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Category updated successfully",
		Data:    nil,
	})
}

func (cc *categoryController) DeleteCategoryById(c *fiber.Ctx) error {
	categoryId := c.Params("category_id")
	if categoryId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Category ID is required",
		})
	}

	err := cc.categoryService.DeleteCategoryById(categoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.Response{
			Status:  fiber.StatusInternalServerError,
			Message: "Failed to delete category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.Response{
		Status:  fiber.StatusOK,
		Message: "Category deleted successfully",
		Data:    nil,
	})
}
