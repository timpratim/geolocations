package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/timpratim/gogeo/models"
	"github.com/timpratim/gogeo/repository"
)

var dbUrl string

const namespace = "surrealdb-golab"
const database = "geocoding"

func init() {
	ipAddress := os.Getenv("IP_ADDRESS") // Get the DB_IP from the environment variables
	if ipAddress == "" {
		ipAddress = "192.168.0.15"
	}
	dbUrl = "ws://" + ipAddress + ":8000/rpc"
}

type RestaurantGeorecoSVC struct {
	repository *repository.GeoRepository
	client     *GeoClient
}

func NewRestaurantGeorecoSVC(repository *repository.GeoRepository, client *GeoClient) *RestaurantGeorecoSVC {
	return &RestaurantGeorecoSVC{repository, client}

}

func (svc *RestaurantGeorecoSVC) SaveUser(user models.User) (interface{}, error) {

	log.Printf("saving user %s", user.Username)
	lat, long, err := svc.client.getGeolocation(context.Background(), user.Address)
	if err != nil {
		return nil, err
	}
	user.Location = [2]float64{lat, long}
	log.Printf("user location: %v", user.Location)
	return svc.repository.SaveUser(user)
}

func (svc *RestaurantGeorecoSVC) SaveRestaurants() ([]interface{}, error) {
	// Read the JSON file
	jsonFile, err := ioutil.ReadFile("restaurants.json")
	if err != nil {
		return nil, err
	}

	// Define a slice to hold our data
	var restaurants []models.Restaurant

	// Unmarshal the JSON data into the slice
	err = json.Unmarshal(jsonFile, &restaurants)
	if err != nil {
		return nil, err
	}
	var results []interface{}
	for _, restaurant := range restaurants {
		r, err := svc.SaveRestaurant(restaurant)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}
	return results, nil
}

func (svc *RestaurantGeorecoSVC) SaveRestaurant(restaurant models.Restaurant) (interface{}, error) {
	lat, long, err := svc.client.getGeolocation(context.Background(), restaurant.Address)
	if err != nil {
		return nil, err
	}
	restaurant.Location = [2]float64{lat, long}

	return svc.repository.SaveRestaurant(restaurant)
}

func (svc *RestaurantGeorecoSVC) GetRestaurantsByLocation(user models.User) (interface{}, error) {
	log.Printf("the closest restaurant to %s is: ", user.Username)

	return svc.repository.GetRestaurantsNearLocation(user.Location)
}

func main() {

	log.Print("connecting to database with dbUrl: ", dbUrl)

	repository, err := repository.NewGeoRepository(dbUrl, "root", "root", namespace, database)
	defer repository.Close()

	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	language := "it"
	region := "it"
	client, err := NewGeoClient(language, region)
	if err != nil {
		log.Fatalf("failed to create geo client: %s", err)
	}

	svc := NewRestaurantGeorecoSVC(repository, client)

	restaurants, err := svc.SaveRestaurants()
	if err != nil {
		log.Fatalf("failed to save restaurants: %s", err)
	}
	fmt.Printf("restaurants: %v", restaurants)

	user := models.User{
		Username: "pratim",
		Address:  "Lungarno del Tempio, 44, 50121 Firenze FI, Italy",
	}
	u, err := svc.SaveUser(user)
	if err != nil {
		log.Fatalf("failed to save user: %s", err)
	}
	fmt.Printf("user: %v", u)

	results, err := svc.GetRestaurantsByLocation(user)
	if err != nil {
		log.Fatalf("failed to get restaurants by location: %s", err)
	}
	fmt.Printf("results: %v", results)
}

func check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}
