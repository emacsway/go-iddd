package customergrpc

import (
	"context"
	"go-iddd/customer/application"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"

	"github.com/golang/protobuf/ptypes/empty"
)

type customerServer struct {
	commandHandler *application.CommandHandler
}

func NewCustomerServer(commandHandler *application.CommandHandler) *customerServer {
	return &customerServer{commandHandler: commandHandler}
}

func (server *customerServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	id := values.GenerateCustomerID()

	command, err := commands.NewRegister(id.ID(), req.EmailAddress, req.GivenName, req.FamilyName)
	if err != nil {
		return nil, err
	}

	if err := server.commandHandler.Register(command); err != nil {
		return nil, err
	}

	return &RegisterResponse{Id: id.ID()}, nil
}

func (server *customerServer) ConfirmEmailAddress(ctx context.Context, req *ConfirmEmailAddressRequest) (*empty.Empty, error) {
	command, err := commands.NewConfirmEmailAddress(req.Id, req.EmailAddress, req.ConfirmationHash)
	if err != nil {
		return nil, err
	}

	if err := server.commandHandler.ConfirmEmailAddress(command); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (server *customerServer) ChangeEmailAddress(ctx context.Context, req *ChangeEmailAddressRequest) (*empty.Empty, error) {
	command, err := commands.NewChangeEmailAddress(req.Id, req.EmailAddress)
	if err != nil {
		return nil, err
	}

	if err := server.commandHandler.ChangeEmailAddress(command); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
