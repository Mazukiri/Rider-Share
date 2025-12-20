package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	TripsCollection     = "trips"
	RideFaresCollection = "ride_fares"
)

// MongoConfig holds MongoDB connection configuration
type MongoConfig struct {
	URI      string
	Database string
}
