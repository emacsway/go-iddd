package customercli

import (
	"fmt"
	"go-iddd/customer/application"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"

	"github.com/urfave/cli"
)

type CustomerApp struct {
	commandHandler *application.CommandHandler
}

func NewCustomerApp(commandHandler *application.CommandHandler) *CustomerApp {
	return &CustomerApp{commandHandler: commandHandler}
}

func (app *CustomerApp) Commands() []cli.Command {
	return []cli.Command{
		{
			Name:      "RegisterCustomer",
			Aliases:   []string{"rc"},
			Usage:     "Register a Customer",
			Action:    app.RegisterCustomer,
			ArgsUsage: "emailAddress givenName familyName",
		},
		{
			Name:      "ConfirmCustomerEmailAddress",
			Aliases:   []string{"cocea"},
			Usage:     "Confirm a Customer's emailAddress",
			Action:    app.ConfirmCustomerEmailAddress,
			ArgsUsage: "id emailAddress confirmationHash",
		},
		{
			Name:      "ChangeCustomerEmailAddress",
			Aliases:   []string{"chcea"},
			Usage:     "Change a Customer's emailAddress",
			Action:    app.ChangeCustomerEmailAddress,
			ArgsUsage: "id emailAddress",
		},
	}
}

func (app *CustomerApp) RegisterCustomer(ctx *cli.Context) error {
	emailAddress := ctx.Args().Get(0)
	givenName := ctx.Args().Get(1)
	familyName := ctx.Args().Get(2)
	id := values.GenerateCustomerID()

	command, err := commands.NewRegister(id.ID(), emailAddress, givenName, familyName)
	if err != nil {
		return err
	}

	if err := app.commandHandler.Register(command); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(
		ctx.App.Writer,
		"Customer registered with id '%s'\n",
		id.ID(),
	)

	return nil
}

func (app *CustomerApp) ConfirmCustomerEmailAddress(ctx *cli.Context) error {
	id := ctx.Args().Get(0)
	emailAddress := ctx.Args().Get(1)
	confirmationHash := ctx.Args().Get(2)

	command, err := commands.NewConfirmEmailAddress(id, emailAddress, confirmationHash)
	if err != nil {
		return err
	}

	if err := app.commandHandler.ConfirmEmailAddress(command); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(
		ctx.App.Writer,
		"successfully confirmed the emailAddress of Customer with id '%s'\n",
		command.CustomerID().ID(),
	)

	return nil
}

func (app *CustomerApp) ChangeCustomerEmailAddress(ctx *cli.Context) error {
	id := ctx.Args().Get(0)
	emailAddress := ctx.Args().Get(1)

	command, err := commands.NewChangeEmailAddress(id, emailAddress)
	if err != nil {
		return err
	}

	if err := app.commandHandler.ChangeEmailAddress(command); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(
		ctx.App.Writer,
		"successfully changed the emailAddress of Customer with id '%s' to '%s\n",
		command.CustomerID().ID(),
		command.EmailAddress().EmailAddress(),
	)

	return nil
}
