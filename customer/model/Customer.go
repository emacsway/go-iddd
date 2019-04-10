package model

import (
	"errors"
	"go-iddd/customer/model/commands"
	"go-iddd/customer/model/valueobjects"
	"go-iddd/shared"
)

type Customer interface {
	Apply(cmd shared.Command) error
}

type customer struct {
	id           valueobjects.ID
	emailAddress valueobjects.ConfirmableEmailAddress
	name         valueobjects.Name
	isRegistered bool
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
	case commands.Register:
		err = customer.register(command)
	case commands.ConfirmEmailAddress:
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
	case commands.Register:
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

		if customer.emailAddress == nil {
			return errors.New("customer - was registered but has no emailAddress")
		}

		if customer.name == nil {
			return errors.New("customer - was registered but has no name")
		}
	}

	return nil
}

func (customer *customer) register(register commands.Register) error {
	customer.id = register.ID()
	customer.emailAddress = register.ConfirmableEmailAddress()
	customer.name = register.Name()

	return nil
}

func (customer *customer) confirmEmailAddress(confirmEmailAddress commands.ConfirmEmailAddress) error {
	var err error

	if customer.emailAddress.IsConfirmed() {
		return nil
	}

	if !customer.emailAddress.Equals(confirmEmailAddress.EmailAddress()) {
		return errors.New("customer - emailAddress can not be confirmed because it has changed")
	}

	if customer.emailAddress, err = customer.emailAddress.Confirm(confirmEmailAddress.ConfirmationHash()); err != nil {
		return err
	}

	return nil
}