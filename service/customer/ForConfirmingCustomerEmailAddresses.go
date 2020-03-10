package customer

import (
	"go-iddd/service/customer/application/domain/customer/commands"
)

type ForConfirmingCustomerEmailAddresses func(command commands.ConfirmCustomerEmailAddress) error
