package main

import (
	"github.com/labstack/echo"
	"github.com/michaelbui/AWSSG-1710/controllers"
	"github.com/michaelbui/AWSSG-1710/types"
	"net/http"
)

func main() {
	configs := &types.AppConfigs{
		DB: &types.DBConfig{
			Type:      "sqlite3",
			Dsn:       "./db.sqlite3",
			Initiated: false,
		},
		Cloud: &types.CloudConfig{
			Type: "aws",
			Configs: types.Configs{
				"s3": types.AwsS3Config{
					Bucket: "awssg-1710",
				},
				"et": types.AwsETConfig{
					Pipeline: "1508403132974-cj4bdz",
					Preset:   "1351620000001-000061",
				},
			},
		},
	}

	e := echo.New()
	defineRoutes(e, configs)
	e.Logger.Fatal(e.Start(":1323"))
}

func defineRoutes(e *echo.Echo, configs *types.AppConfigs) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/files", func(c echo.Context) error {
		files, err := controllers.NewFileController(configs).List()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, files)
	})

	e.POST("/files", func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		id, err := controllers.NewFileController(configs).Post(file)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, struct {
			Id uint
		}{
			Id: id,
		})
	})
}
