package eventstore

import (
	"database/sql"
	"math"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

const streamPrefix = "customer"

type CustomerEventStore struct {
	eventStore           es.EventStore
	uniqueEmailAddresses ForCheckingUniqueEmailAddresses
	db                   *sql.DB
}

func NewCustomerEventStore(
	eventStore es.EventStore,
	uniqueEmailAddresses ForCheckingUniqueEmailAddresses,
	db *sql.DB,
) *CustomerEventStore {

	return &CustomerEventStore{
		eventStore:           eventStore,
		uniqueEmailAddresses: uniqueEmailAddresses,
		db:                   db,
	}
}

func (store *CustomerEventStore) EventStreamFor(id values.CustomerID) (es.DomainEvents, error) {
	wrapWithMsg := "customerEventStore.EventStreamFor"

	eventStream, err := store.eventStore.LoadEventStream(store.streamID(id), 0, math.MaxUint32)
	if err != nil {
		return nil, errors.Wrap(err, wrapWithMsg)
	}

	if len(eventStream) == 0 {
		err := errors.New("customer not found")
		return nil, lib.MarkAndWrapError(err, lib.ErrNotFound, wrapWithMsg)
	}

	return eventStream, nil
}

func (store *CustomerEventStore) CreateStreamFrom(recordedEvents es.DomainEvents, id values.CustomerID) error {
	var err error
	wrapWithMsg := "customerEventStore.CreateStreamFrom"

	tx, err := store.db.Begin()
	if err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	for _, event := range recordedEvents {
		switch actualEvent := event.(type) {
		case events.CustomerRegistered:
			if err = store.uniqueEmailAddresses.AddUniqueEmailAddress(actualEvent.EmailAddress(), tx); err != nil {
				_ = tx.Rollback()

				return errors.Wrap(err, wrapWithMsg)
			}
		}
	}

	if err = store.eventStore.AppendEventsToStream(store.streamID(id), recordedEvents, tx); err != nil {
		_ = tx.Rollback()

		if errors.Is(err, lib.ErrConcurrencyConflict) {
			return lib.MarkAndWrapError(errors.New("found duplicate customer"), lib.ErrDuplicate, wrapWithMsg)
		}

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (store *CustomerEventStore) Add(recordedEvents es.DomainEvents, id values.CustomerID) error {
	var err error
	wrapWithMsg := "customerEventStore.Add"

	tx, err := store.db.Begin()
	if err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	for _, event := range recordedEvents {
		switch actualEvent := event.(type) {
		case events.CustomerEmailAddressChanged:
			if err = store.uniqueEmailAddresses.ReplaceUniqueEmailAddress(actualEvent.PreviousEmailAddress(), actualEvent.EmailAddress(), tx); err != nil {
				_ = tx.Rollback()

				return errors.Wrap(err, wrapWithMsg)
			}
		case events.CustomerDeleted:
			if err = store.uniqueEmailAddresses.RemoveUniqueEmailAddress(actualEvent.EmailAddress(), tx); err != nil {
				_ = tx.Rollback()

				return errors.Wrap(err, wrapWithMsg)
			}
		}
	}

	if err = store.eventStore.AppendEventsToStream(store.streamID(id), recordedEvents, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (store *CustomerEventStore) Delete(id values.CustomerID) error {
	if err := store.eventStore.PurgeEventStream(store.streamID(id)); err != nil {
		return errors.Wrap(err, "customerEventStore.Delete")
	}

	return nil
}

func (store *CustomerEventStore) streamID(id values.CustomerID) es.StreamID {
	return es.NewStreamID(streamPrefix + "-" + id.ID())
}
