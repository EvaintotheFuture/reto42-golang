package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/EvaintotheFuture/reto42-golang/models"
	"github.com/EvaintotheFuture/reto42-golang/responses"
	"github.com/EvaintotheFuture/reto42-golang/configs"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var asteroidsCollection *mongo.Collection = configs.GetCollection(configs.DB, "asteroides")
var validate = validator.New()

func CreateAsteroid(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var asteroid models.Asteroid
	defer cancel()

	if err := c.Bind(&asteroid); err != nil {
		return c.JSON(http.StatusBadRequest, responses.AsteroidResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&asteroid); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.AsteroidResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	newAsteroid := models.Asteroid {

		ID:				primitive.NewObjectID(),
		Name:			asteroid.Name,
		Diameter:		asteroid.Diameter,
		Discovery_date:	asteroid.Discovery_date,
		Observations:	asteroid.Observations,
		Distances:		asteroid.Distances,
	}

	result, err := asteroidsCollection.InsertOne(ctx, newAsteroid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.AsteroidResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, responses.AsteroidResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": result}})
}