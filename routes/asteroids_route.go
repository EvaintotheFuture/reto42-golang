package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/EvaintotheFuture/reto42-golang/controllers"
)

func AsteroidRoute (e *echo.Echo) {

	api := e.Group("/api/v1", serverHeader)

	api.POST("/asteroides", controllers.CreateAsteroid)
	api.GET("/asteroides/:asteroidID", controllers.GetAsteroid)
	api.GET("/asteroides", controllers.GetAllAsteroids)
	api.PATCH("/asteroides/:asteroidID", controllers.EditAsteroid)
	api.DELETE("/asteroides/:asteroidID", controllers.DeleteAsteroid)

}

func serverHeader (next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("x-version", "Test/v1.0")
		return next(c)
	}
}