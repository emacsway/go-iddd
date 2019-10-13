package domain

import (
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"

	"github.com/cockroachdb/errors"
)

type Customer struct {
	id                      *values.ID
	confirmableEmailAddress *values.ConfirmableEmailAddress
	personName              *values.PersonName
	currentStreamVersion    uint
	recordedEvents          shared.DomainEvents
}

func blankCustomer() *Customer {
	return &Customer{}
}

func (customer *Customer) ID() shared.IdentifiesAggregates {
	return customer.id
}

func ReconstituteCustomerFrom(eventStream shared.DomainEvents) (*Customer, error) {
	newCustomer := blankCustomer()

	if err := eventStream.FirstEventShouldBeOfSameTypeAs(&events.Registered{}); err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrInvalidEventStream), "shared.ErrInvalidEventStream")
	}

	newCustomer.apply(eventStream...)

	return newCustomer, nil
}

func (customer *Customer) recordThat(event shared.DomainEvent) {
	customer.recordedEvents = append(customer.recordedEvents, event)
	customer.apply(event)
}

func (customer *Customer) apply(eventStream ...shared.DomainEvent) {
	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case *events.Registered:
			customer.whenItWasRegistered(actualEvent)
		case *events.EmailAddressConfirmed:
			customer.whenEmailAddressWasConfirmed(actualEvent)
		case *events.EmailAddressChanged:
			customer.whenEmailAddressWasChanged(actualEvent)
		}

		customer.currentStreamVersion = event.StreamVersion()
	}
}

func (customer *Customer) RecordedEvents() shared.DomainEvents {
	return customer.recordedEvents
}

func (customer *Customer) PurgeRecordedEvents() {
	customer.recordedEvents = nil
}

func (customer *Customer) StreamVersion() uint {
	return customer.currentStreamVersion
}
