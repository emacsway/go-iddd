package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

func ChangeName(eventStream es.EventStream, command domain.ChangeCustomerName) (es.RecordedEvents, error) {
	customer := buildCurrentStateFrom(eventStream)

	if err := assertNotDeleted(customer); err != nil {
		return nil, errors.Wrap(err, "changeCustomerName")
	}

	if customer.personName.Equals(command.PersonName()) {
		return nil, nil
	}

	event := domain.BuildCustomerNameChanged(
		customer.id,
		command.PersonName(),
		customer.currentStreamVersion+1,
	)

	return es.RecordedEvents{event}, nil
}
