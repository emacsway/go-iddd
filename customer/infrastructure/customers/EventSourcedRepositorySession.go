package customers

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"math"

	"golang.org/x/xerrors"
)

type EventSourcedRepositorySession struct {
	eventStoreSession shared.EventStore
	customerFactory   func(eventStream shared.DomainEvents) (*domain.Customer, error)
}

/***** Implement domain.Customers *****/

func (session *EventSourcedRepositorySession) Register(customer *domain.Customer) error {
	streamID := shared.NewStreamID(streamPrefix + "-" + customer.ID().String())

	recordedEvents := customer.RecordedEvents()
	customer.PurgeRecordedEvents()

	if err := session.eventStoreSession.AppendEventsToStream(streamID, recordedEvents); err != nil {
		if xerrors.Is(err, shared.ErrConcurrencyConflict) {
			return xerrors.Errorf("eventSourcedRepositorySession.Register: %s: %w", err, shared.ErrDuplicate)
		}

		return xerrors.Errorf("eventSourcedRepositorySession.Register: %w", err)
	}

	return nil
}

func (session *EventSourcedRepositorySession) Of(id *values.ID) (*domain.Customer, error) {
	streamID := shared.NewStreamID(streamPrefix + "-" + id.String())

	eventStream, err := session.eventStoreSession.LoadEventStream(streamID, 0, math.MaxUint32)
	if err != nil {
		return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: %w", err)
	}

	if len(eventStream) == 0 {
		return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: event stream is empty: %w", shared.ErrNotFound)
	}

	customer, err := session.customerFactory(eventStream)
	if err != nil {
		return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: %w", err)
	}

	return customer, nil
}

/***** Implement application.PersistsCustomers *****/

func (session *EventSourcedRepositorySession) Persist(customer *domain.Customer) error {
	streamID := shared.NewStreamID(streamPrefix + "-" + customer.ID().String())

	recordedEvents := customer.RecordedEvents()
	customer.PurgeRecordedEvents()

	if err := session.eventStoreSession.AppendEventsToStream(streamID, recordedEvents); err != nil {
		return xerrors.Errorf("eventSourcedRepositorySession.Persist: %w", err)
	}

	return nil
}
