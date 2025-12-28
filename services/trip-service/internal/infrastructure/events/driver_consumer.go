package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"
	pbd "ride-sharing/shared/proto/driver"

	"github.com/rabbitmq/amqp091-go"
)

type driverConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  domain.TripService
}

func NewDriverConsumer(rabbitmq *messaging.RabbitMQ, service domain.TripService) *driverConsumer {
	return &driverConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (c *driverConsumer) Listen() error {
	return c.rabbitmq.ConsumeMessages(messaging.DriverTripResponseQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		var payload messaging.DriverTripResponseData
		if err := json.Unmarshal(message.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		log.Printf("driver response received message: %+v", payload)

		switch msg.RoutingKey {
		case contracts.DriverCmdTripAccept:
			if err := c.handleTripAccepted(ctx, payload.TripID, payload.Driver); err != nil {
				log.Printf("Failed to handle the trip accept: %v", err)
				return err
			}
		case contracts.DriverCmdTripDecline:
			if err := c.handleTripDeclined(ctx, payload.TripID, payload.RiderID, payload.Driver.Id); err != nil {
				log.Printf("Failed to handle the trip decline: %v", err)
				return err
			}
			return nil
		case contracts.DriverCmdTripComplete:
			if err := c.handleTripCompleted(ctx, payload.TripID, payload.RiderID, payload.Driver.Id); err != nil {
				log.Printf("Failed to handle the trip complete: %v", err)
				return err
			}
			return nil
		}
		log.Printf("unknown trip event: %+v", payload)

		return nil
	})
}

func (c *driverConsumer) handleTripDeclined(ctx context.Context, tripID, riderID, driverID string) error {
	// 1. Get current trip to see previous declines
	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	// 2. Add this driver to the declined list
	newDeclinedList := append(trip.CandidateDriverIDs, driverID)

	// 3. Update DB
	if err := c.service.AddCandidateDrivers(ctx, tripID, newDeclinedList); err != nil {
		return err
	}

	// 4. Publish "Driver Not Interested" event with the UPDATED list of excluded drivers
	newPayload := messaging.TripEventData{
		Trip:              trip.ToProto(),
		DeclinedDriverIDs: newDeclinedList,
	}

	marshalledPayload, err := json.Marshal(newPayload)
	if err != nil {
		return err
	}

	if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventDriverNotInterested,
		contracts.AmqpMessage{
			OwnerID: riderID,
			Data:    marshalledPayload,
		},
	); err != nil {
		return err
	}

	return nil
}

func (c *driverConsumer) handleTripAccepted(ctx context.Context, tripID string, driver *pbd.Driver) error {
	// 1. Fetch the first
	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	if trip == nil {
		return fmt.Errorf("Trip was not found %s", tripID)
	}

	// Fix Race Condition: Check if trip is already accepted
	if trip.Status == "accepted" || trip.Driver != nil {
		log.Printf("Trip %s already accepted by driver %s. Ignoring request from %s", tripID, trip.Driver.Id, driver.Id)
		return nil
	}

	// 2. Update the trip
	if err := c.service.UpdateTrip(ctx, tripID, "accepted", driver); err != nil {
		log.Printf("Failed to update the trip: %v", err)
		return err
	}

	trip, err = c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	// 3. Driver has been assigned -> publish this event to RB
	marshalledTrip, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	// Notify the rider that a driver has been assigned
	if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventDriverAssigned, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    marshalledTrip,
	}); err != nil {
		return err
	}

	// Removed premature PaymentCmdCreateSession. 
	// Payment should be triggered only after Trip Completion via DriverCmdTripComplete (handled in API Gateway).

	return nil
}

func (c *driverConsumer) handleTripCompleted(ctx context.Context, tripID, riderID, driverID string) error {
	// 1. Get current trip
	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	if trip == nil {
		return fmt.Errorf("Trip was not found %s", tripID)
	}

	// 2. Update the trip status
	if err := c.service.UpdateTrip(ctx, tripID, "completed", nil); err != nil {
		log.Printf("Failed to update the trip status to completed: %v", err)
		return err
	}

	log.Printf("Trip %s completed by driver %s", tripID, driverID)

	// 3. Trigger Payment Session Creation
	// We need to construct the payload expected by Payment Service
	trip.Status = "completed"
	marshalledTrip, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	if err := c.rabbitmq.PublishMessage(ctx, contracts.PaymentCmdCreateSession, contracts.AmqpMessage{
		OwnerID: riderID,
		Data:    marshalledTrip,
	}); err != nil {
		log.Printf("Failed to publish payment creation message: %v", err)
		return err
	}

	return nil
}
