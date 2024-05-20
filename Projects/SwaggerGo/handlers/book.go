package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// ResponseHTTP represents response body of this API
type ResponseHTTP struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// GetAllBooks is a function to get all books data from database
//
//	@Summary		Get all books
//	@Description	Get all books
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ResponseHTTP{data=[]models.Book}
//	@Failure		503	{object}	ResponseHTTP{}
//	@Router			/v1/books [get]
func GetAllBooks(c *fiber.Ctx) error {
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get all books.",
	})
}

// GetBookByID is a function to get a book by ID
//
//	@Summary		Get book by ID
//	@Description	Get book by ID
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Book ID"
//	@Success		200	{object}	ResponseHTTP{data=[]models.Book}
//	@Failure		404	{object}	ResponseHTTP{}
//	@Failure		503	{object}	ResponseHTTP{}
//	@Router			/v1/books/{id} [get]
func GetBookByID(c *fiber.Ctx) error {
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success get book by ID.",
	})
}

// RegisterBook registers a new book data
//
//	@Summary		Register a new book
//	@Description	Register book
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			book	body		models.Book	true	"Register book"
//	@Success		200		{object}	ResponseHTTP{data=models.Book}
//	@Failure		400		{object}	ResponseHTTP{}
//	@Router			/v1/books [post]
func RegisterBook(c *fiber.Ctx) error {
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success register a book.",
	})
}

// DeleteBook function removes a book by ID
//
//	@Summary		Remove book by ID
//	@Description	Remove book by ID
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Book ID"
//	@Success		200	{object}	ResponseHTTP{}
//	@Failure		404	{object}	ResponseHTTP{}
//	@Failure		503	{object}	ResponseHTTP{}
//	@Router			/v1/books/{id} [delete]
func DeleteBook(c *fiber.Ctx) error {
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Success delete book.",
		Data:    nil,
	})
}
