package service

import (
	"context"
	"errors"
	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	pbd "ride-sharing/shared/proto/driver"
	"testing"
)

// MockRepository implements domain.TripRepository for testing
type MockRepository struct {
	createTripFunc func(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error)
	saveRideFareFunc func(ctx context.Context, f *domain.RideFareModel) error
}

func (m *MockRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	if m.createTripFunc != nil {
		return m.createTripFunc(ctx, trip)
	}
	return trip, nil
}

func (m *MockRepository) SaveRideFare(ctx context.Context, f *domain.RideFareModel) error {
	if m.saveRideFareFunc != nil {
		return m.saveRideFareFunc(ctx, f)
	}
	return nil
}

// Stubs for other methods to satisfy interface
func (m *MockRepository) GetRideFareByID(ctx context.Context, id string) (*domain.RideFareModel, error) { return nil, nil }
func (m *MockRepository) GetTripByID(ctx context.Context, id string) (*domain.TripModel, error) { return nil, nil }
func (m *MockRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *pbd.Driver) error { return nil }
func (m *MockRepository) RemoveCandidateDriver(ctx context.Context, tripID, driverID string) error { return nil }
func (m *MockRepository) AddCandidateDrivers(ctx context.Context, tripID string, driverIDs []string) error { return nil }
func (m *MockRepository) SaveOutboxEvent(ctx context.Context, eventType string, payload interface{}) error { return nil }

func TestCreateTrip(t *testing.T) {
	tests := []struct {
		name    string
		fare    *domain.RideFareModel
		mockErr error
		wantErr bool
	}{
		{
			name: "Success",
			fare: &domain.RideFareModel{
				UserID: "user-123",
				TotalPriceInCents: 1000,
			},
			wantErr: false,
		},
		{
			name: "Repository Error",
			fare: &domain.RideFareModel{
				UserID: "user-123",
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepository{
				createTripFunc: func(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					// Validation
					if trip.Status != "pending" {
						t.Errorf("expected status pending, got %s", trip.Status)
					}
					return trip, nil
				},
			}

			svc := NewService(repo)
			_, err := svc.CreateTrip(context.Background(), tt.fare)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTrip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateTripFares(t *testing.T) {
	userID := "user-high-value"
	route := &tripTypes.OsrmApiResponse{}
	baseFares := []*domain.RideFareModel{
		{TotalPriceInCents: 200, PackageSlug: "suv"},
		{TotalPriceInCents: 350, PackageSlug: "luxury"},
	}

	repo := &MockRepository{
		saveRideFareFunc: func(ctx context.Context, f *domain.RideFareModel) error {
			if f.UserID != userID {
				t.Errorf("expected user ID %s, got %s", userID, f.UserID)
			}
			return nil
		},
	}

	svc := NewService(repo)
	fares, err := svc.GenerateTripFares(context.Background(), baseFares, userID, route)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fares) != 2 {
		t.Errorf("expected 2 fares, got %d", len(fares))
	}
}
