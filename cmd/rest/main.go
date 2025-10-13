package main

// go run cmd/rest/main.go

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/luzmareto/go-grpc-ecommerce-be/internal/handler"
)

func handlerGetFileName(c *fiber.Ctx) error {
	fileNameParam := c.Params("filename")
	filePath := path.Join("storage", "product", fileNameParam)
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return c.Status(http.StatusNotFound).SendString("Not found")
		}
		log.Println(err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	ext := path.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)

	c.Set("Content-Type", mimeType)

	return c.SendStream(file)
}

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/storage/products/:filename", handlerGetFileName) // Untuk List Product
	app.Get("/storage/product/:filename", handlerGetFileName)  // Untuk Detail/Edit Product

	app.Post("/product/upload", handler.UploadProductImageHandler)
	app.Post("/products/upload", handler.UploadProductImageHandler)

	app.Listen(":3000")
}
