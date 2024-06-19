package docs

import (
	"embed"

	"github.com/gofiber/fiber/v3"
)

//go:embed swagger.json
var swaggerjson embed.FS

func SwaggerHandler(app *fiber.App) {
	app.Get("/docs", func(c fiber.Ctx) error {
		c.Type("html")
		return c.SendString(`<!DOCTYPE html>
<html>

<head>
    <title>Redoc</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">

    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>

<body>
    <redoc spec-url='./docs/swagger.json' required-props-first=true></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
</body>

</html>`)
	},
	)
}
