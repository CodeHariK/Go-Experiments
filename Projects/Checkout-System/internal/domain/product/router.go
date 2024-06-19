package product

import (
	"github.com/gofiber/fiber/v3"
)

func SetupRouter(r *fiber.App) fiber.Router {
	productRepository := NewProductRepository()
	productService := NewProductService(*productRepository)
	productRouter := NewProductRouter(*productService)

	productRouterGroup := r.Group("/product")

	productRouterGroup.Get("/list", productRouter.list)
	productRouterGroup.Post("/create", productRouter.create)
	productRouterGroup.Post("/update", productRouter.update)
	productRouterGroup.Post("/delete", productRouter.delete)

	return productRouterGroup
}

type ProductRouter struct {
	ProductService ProductService
}

func NewProductRouter(productService ProductService) *ProductRouter {
	return &ProductRouter{
		ProductService: productService,
	}
}

func (p *ProductRouter) list(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (p *ProductRouter) create(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (p *ProductRouter) update(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func (p *ProductRouter) delete(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
