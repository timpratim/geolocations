package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
)

type GeoClient struct {
	client   *maps.Client
	language string
	region   string
}

func NewGeoClient(language string, region string) (*GeoClient, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the value of the environment variable
	GEO_API_KEY := os.Getenv("GEO_API_KEY")
	client, err := maps.NewClient(maps.WithAPIKey(GEO_API_KEY))
	if err != nil {
		return nil, err
	}
	return &GeoClient{client, language, region}, nil
}

func (c *GeoClient) getGeolocation(ctx context.Context, address string) (float64, float64, error) {
	r := &maps.GeocodingRequest{
		Address:  address,
		Language: c.language,
		Region:   c.region,
	}
	resp, err := c.client.Geocode(ctx, r)
	if err != nil {
		return 0, 0, err
	}
	lat := resp[0].Geometry.Location.Lat
	long := resp[0].Geometry.Location.Lng
	return lat, long, nil
}
