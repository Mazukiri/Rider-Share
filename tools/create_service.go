package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	serviceName := flag.String("name", "", "Name of the service (e.g., user, payment)")
	flag.Parse()

	if *serviceName == "" {
		fmt.Println("Please provide a service name using -name flag")
		os.Exit(1)
	}

	// Create service directory structure
	basePath := filepath.Join("services", *serviceName+"-service")
	dirs := []string{
		"cmd",
		"internal/domain",
		"internal/service",
		"internal/infrastructure/events",
		"internal/infrastructure/grpc",
		"internal/infrastructure/repository",
		"pkg/types",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(basePath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// Create an empty README.md
	readmePath := filepath.Join(basePath, "README.md")
}