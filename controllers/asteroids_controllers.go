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
	"go.mongodb.org/mongo-driver/bson"
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
		DiscoveryDate:	asteroid.DiscoveryDate,
		Observations:	asteroid.Observations,
		Distances:		asteroid.Distances,
	}

	result, err := asteroidsCollection.InsertOne(ctx, newAsteroid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.AsteroidResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, responses.AsteroidResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": result}})
}

func GetAsteroid(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	asteroidID := c.Param("asteroidID")
	var asteroid models.Asteroid
	defer cancel()

	ObjId, _ := primitive.ObjectIDFromHex(asteroidID)

	err := asteroidsCollection.FindOne(ctx, bson.M{"_id": ObjId}).Decode(&asteroid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.AsteroidResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
}

	return c.JSON(http.StatusOK, responses.AsteroidResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": asteroid}})
}

func EditAsteroid(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	asteroidID := c.Param("asteroidID")
	var asteroid models.Asteroid
	defer cancel()

	ObjId, _ := primitive.ObjectIDFromHex(asteroidID)

	if err := c.Bind(&asteroid); err != nil {
		return c.JSON(http.StatusBadRequest, responses.AsteroidResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if validationErr := validate.Struct(&asteroid); validationErr != nil {
		return c.JSON(http.StatusBadRequest, responses.AsteroidResponse{Status: http.StatusBadRequest, Message: "error", Data: &echo.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"name": asteroid.Name,"diameter": asteroid.Diameter,  "discovery_date": asteroid.DiscoveryDate, "observations": asteroid.Observations, "distances": asteroid.Distances}

	result, err := asteroidsCollection.UpdateOne(ctx, bson.M{"_id": ObjId}, bson.M{"$set": update})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.AsteroidResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	var updateAsteroid models.Asteroid
	if result.MatchedCount == 1 {
		err := asteroidsCollection.FindOne(ctx, bson.M{"_id": ObjId}).Decode(&updateAsteroid)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.AsteroidResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}
	}
	return c.JSON(http.StatusOK, responses.AsteroidResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updateAsteroid}})
}

func DeleteAsteroid(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	asteroidID := c.Param("asteroidID")
	defer cancel()

	ObjId, _ := primitive.ObjectIDFromHex(asteroidID)

	result, err := asteroidsCollection.DeleteOne(ctx, bson.M{"_id": ObjId})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.AsteroidResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.JSON(http.StatusNotFound, responses.AsteroidResponse{Status: http.StatusNotFound, Message: "success", Data: &echo.Map{"data": "Asteroid with specified ID not found"}})
	}
	return c.JSON(http.StatusOK, responses.AsteroidResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": "Asteroid successfully deleted"}})
}

func GetAllAsteroids(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var asteroids []models.Asteroid
	defer cancel()

	results, err := asteroidsCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.AsteroidResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleAsteroid models.Asteroid
		if err = results.Decode(&singleAsteroid); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.AsteroidResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	asteroids = append(asteroids, singleAsteroid)
	}

	return c.JSON(http.StatusOK, responses.AsteroidResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": asteroids}})

}