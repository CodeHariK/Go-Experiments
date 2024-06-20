package docs

import (
	"embed"

	"github.com/gofiber/fiber/v3"
)

//go:embed swagger.json
var swaggerjson embed.FS

func SwaggerHandler(app *fiber.App) {
	// Serve the embedded swagger.json file at /docs/swagger.json
	app.Get("/docs/swagger.json", func(c fiber.Ctx) error {
		swaggerBytes, err := swaggerjson.ReadFile("swagger.json")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to read swagger.json")
		}
		c.Type("json")
		return c.Send(swaggerBytes)
	})

	// Serve the Redoc HTML
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
		.react-tabs__tab-panel {
			border-radius: 20px;
		}
		input {
			padding: 15px !important;
			margin: 6px !important;
			border: 1px solid white !important;
			border-radius: 5px !important;
			font-size: 16px !important;
		}
		.topbar-wrapper img {
			content:url('https://www.shutterstock.com/shutterstock/photos/1992498650/display_1500/stock-vector-colors-font-alphabet-letters-modern-logo-typography-color-creative-art-typographic-design-1992498650.jpg');
		}
    </style>
</head>

<body>
	<div id="redoc-container"></div>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
	<script>
        Redoc.init('/docs/swagger.json', {
            theme: {
				sidebar: {					
					backgroundColor: '#43395d',
					textColor: '#ffffff'
				},
				rightPanel: {					
					backgroundColor: '#43395d',
				},
                colors: {
                    primary: {
                        main: '#5a5a5a'
                    }
                },
                typography: {
                    fontSize: '16px',
                    fontFamily: 'Roboto, sans-serif',
                    headings: {
                        fontFamily: 'Montserrat, sans-serif'
                    }
                }
            }
        }, document.getElementById('redoc-container'))


		setTimeout(() => {
			var searchDiv = document.querySelector('div[role="search"]');
			var img = document.createElement('img');
			img.src = 'https://cdn.pixabay.com/photo/2017/03/16/21/18/logo-2150297_640.png';
			img.style.width = "100%"
			img.style.borderRadius = "20px"
			img.style.padding = "15px"
			searchDiv.appendChild(img);
		}, 2000);
    </script>
</body>

</html>`)
	})
}
