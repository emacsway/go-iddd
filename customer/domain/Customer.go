package domain

import (
	"errors"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

//go:generate mockery -name Customer -output ../application/mocks -outpkg mocks -note "Regenerate by running `go generate` in domain/Customer"

type Customer interface {
	Apply(cmd shared.Command) error

	shared.Aggregate
}

type customer struct {
	id                      *values.ID
	confirmableEmailAddress *values.ConfirmableEmailAddress
	personName              *values.PersonName
	isRegistered            bool
}

func NewUnregisteredCustomer() *customer {
	return &customer{}
}

func (customer *customer) Apply(command shared.Command) error {
	var err error

	if err := customer.assertCustomerIsInValidState(command); err != nil {
		return err
	}

	switch command := command.(type) {
	case *commands.Register:
		customer.register(command)
	case *commands.ConfirmEmailAddress:
		err = customer.confirmEmailAddress(command)
	case nil:
		err = errors.New("customer - nil command applied")
	default:
		err = errors.New("customer - unknown command applied")
	}

	return err
}

func (customer *customer) assertCustomerIsInValidState(command shared.Command) error {
	switch command.(type) {
	case *commands.Register:
		if customer.isRegistered {
			return errors.New("customer - was already registered")
		}
	default:
		if !customer.isRegistered {
			return errors.New("customer - was not registered yet")
		}

		if customer.id == nil {
			return errors.New("customer - was registered but has no id")
		}

		if customer.confirmableEmailAddress == nil {
			return errors.New("customer - was registered but has no emailAddress")
		}

		if customer.personName == nil {
			return errors.New("customer - was registered but has no personName")
		}
	}

	return nil
}

func (customer *customer) register(given *commands.Register) {
	customer.id = given.ID()
	customer.confirmableEmailAddress = given.EmailAddress().ToConfirmable()
	customer.personName = given.PersonName()
}

func (customer *customer) confirmEmailAddress(given *commands.ConfirmEmailAddress) error {
	var err error

	if customer.confirmableEmailAddress.IsConfirmed() {
		return nil
	}

	if !customer.confirmableEmailAddress.Equals(given.EmailAddress()) {
		return errors.New("customer - emailAddress can not be confirmed because it has changed")
	}

	if customer.confirmableEmailAddress, err = customer.confirmableEmailAddress.Confirm(given.ConfirmationHash()); err != nil {
		return err
	}

	return nil
}

func (customer *customer) AggregateIdentifier() shared.AggregateIdentifier {
	return customer.id
}

func (customer *customer) AggregateName() string {
	return shared.BuildAggregateNameFor(customer)
}
