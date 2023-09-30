package main

import (
	"net/http"

	"github.com/elmehdiabdi-src/exclusive"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {

	engine := gin.Default()

	engine.StaticFile("/docs", "./public/index.html")

	engine.GET("/swagger.json", func(c *gin.Context) {

		schema := exclusive.Swag(engine, c, &exclusive.Configure{
			Responses: map[int]any{
				429: &ErrorResponse{},
				400: &ErrorResponse{},
			},
		})

		c.String(http.StatusOK, schema)

	})

	v1 := engine.Group("api/v1")

	{

		v1.POST("/login", func(c *gin.Context) {

			type LoginRequest struct {
				Type     string `query:"type" binding:"required"`
				Email    string `json:"email" binding:"required"`
				Password string `json:"password" binding:"required"`
			}

			type LoginResponse struct {
				Token string `json:"token"`
			}

			c.Set("doc", exclusive.Doc{
				IsDeprecated: false,
				Tags:         "Auth",
				ID:           "Login",
				Request:      new(LoginRequest),
				Response:     new(LoginResponse),
				Description:  "Endpoint to generate any token.",
			})

		})
	}

	engine.Run()
}
