package main

import (
	"sync"
	"testing"
)

func TestFindAvailableDrivers(t *testing.T) {
	svc := NewService()

	// Register mixed drivers
	_, _ = svc.RegisterDriver("driver-1", "suv")
	_, _ = svc.RegisterDriver("driver-2", "sedan")
	_, _ = svc.RegisterDriver("driver-3", "suv")

	// Test SUV matching
	suvs := svc.FindAvailableDrivers("suv")
	if len(suvs) != 2 {
		t.Errorf("expected 2 suvs, got %d", len(suvs))
	}

	// Test Sedan matching
	sedans := svc.FindAvailableDrivers("sedan")
	if len(sedans) != 1 {
		t.Errorf("expected 1 sedan, got %d", len(sedans))
	}
	if sedans[0] != "driver-2" {
		t.Errorf("expected driver-2, got %s", sedans[0])
	}

	// Test non-existent package
	vans := svc.FindAvailableDrivers("van")
	if len(vans) != 0 {
		t.Errorf("expected 0 vans, got %d", len(vans))
	}
}

func TestRegisterDriverConcurrency(t *testing.T) {
	svc := NewService()
	var wg sync.WaitGroup
	count := 100

	// Concurrent registration
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := svc.RegisterDriver("driver-concurrent", "suv")
			if err != nil {
				t.Errorf("failed to register: %v", err)
			}
		}(i)
	}

	wg.Wait()

	// Verify count (thread safety check)
	drivers := svc.FindAvailableDrivers("suv")
	if len(drivers) != count {
		t.Errorf("Race condition detected! Expected %d drivers, got %d", count, len(drivers))
	}
}
