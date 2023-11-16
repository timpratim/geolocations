package repository

import (
	"fmt"

	surreal "github.com/surrealdb/surrealdb.go"
	"github.com/timpratim/gogeo/models"
)

type GeoRepository struct {
	db *surreal.DB
}

func NewGeoRepository(address, user, password, namespace, database string) (*GeoRepository, error) {

	db, err := surreal.New(address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", err)
	}
	_, err = db.Signin(map[string]interface{}{
		"user": user,
		"pass": password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sign in: %w", err)
	}

	_, err = db.Use(namespace, database)
	if err != nil {
		return nil, err
	}

	return &GeoRepository{db}, nil
}

func (r GeoRepository) Close() {
	r.db.Close()
}

// save user
func (r GeoRepository) SaveUser(user models.User) (interface{}, error) {
	fmt.Println("saving user %s", user.Username)
	return r.db.Create("users", map[string]interface{}{
		"username": user.Username,
		"address":  user.Address,
		"location": user.Location,
	})
}

// get user location
func (r GeoRepository) GetUserLocation(username string) (interface{}, error) {
	fmt.Println("getting user location %s", username)
	return r.db.Query("SELECT location FROM users WHERE username = $username LIMIT 1", map[string]interface{}{
		"username": username,
	})
}

func (r GeoRepository) SaveRestaurant(restaurant models.Restaurant) (interface{}, error) {
	fmt.Println("saving restaurant %s", restaurant.Name)
	return r.db.Create("restaurants", map[string]interface{}{
		"name":     restaurant.Name,
		"address":  restaurant.Address,
		"location": restaurant.Location,
	})
}

func (r GeoRepository) GetRestaurantsNearLocation(userLocation [2]float64) (interface{}, error) {

	return r.db.Query("SELECT name, geo::distance(type::point($userLocation), type::point(location)) as distance FROM restaurants ORDER BY distance", map[string]interface{}{
		"userLocation": userLocation,
	})
}
