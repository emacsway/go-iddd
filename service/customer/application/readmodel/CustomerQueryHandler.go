package readmodel

import (
	"go-iddd/service/customer/application/readmodel/domain/customer"

	"github.com/cockroachdb/errors"
)

type CustomerQueryHandler struct {
	customerEvents ForReadingCustomerEventStreams
}

func NewCustomerQueryHandler(customerEvents ForReadingCustomerEventStreams) *CustomerQueryHandler {
	return &CustomerQueryHandler{
		customerEvents: customerEvents,
	}
}

func (h *CustomerQueryHandler) CustomerViewByID(customerID customer.ID) (customer.View, error) {
	eventStream, err := h.customerEvents.EventStreamFor(customerID)
	if err != nil {
		return customer.View{}, errors.Wrap(err, "customerQueryHandler.CustomerViewByID")
	}

	customerView := customer.BuildViewFrom(eventStream)

	return customerView, nil
}