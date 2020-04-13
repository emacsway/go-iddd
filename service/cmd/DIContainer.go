package cmd

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/adapter/secondary/postgres"

	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"
	"github.com/AntonStoeckl/go-iddd/service/customer/application/query"
	customergrpc "github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/adapter/primary/grpc"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

const (
	eventStoreTableName           = "eventstore"
	uniqueEmailAddressesTableName = "unique_email_addresses"
)

type DIContainer struct {
	postgresDBConn         *sql.DB
	marshalCustomerEvent   es.MarshalDomainEvent
	unmarshalCustomerEvent es.UnmarshalDomainEvent
	customerEventStore     *postgres.CustomerEventStore
	customerCommandHandler *command.CustomerCommandHandler
	customerQueryHandler   *query.CustomerQueryHandler
	customerGRPCServer     customergrpc.CustomerServer
}

func NewDIContainer(
	postgresDBConn *sql.DB,
	marshalCustomerEvent es.MarshalDomainEvent,
	unmarshalCustomerEvent es.UnmarshalDomainEvent,
) (*DIContainer, error) {

	if postgresDBConn == nil {
		return nil, errors.Mark(errors.New("newDIContainer: postgres DB connection must not be nil"), lib.ErrTechnical)
	}

	container := &DIContainer{
		postgresDBConn:         postgresDBConn,
		marshalCustomerEvent:   marshalCustomerEvent,
		unmarshalCustomerEvent: unmarshalCustomerEvent,
	}

	container.init()

	return container, nil
}

func (container DIContainer) init() {
	container.GetCustomerEventStore()
	container.GetCustomerCommandHandler()
	container.GetCustomerQueryHandler()
	container.GetCustomerGRPCServer()
}

func (container DIContainer) GetPostgresDBConn() *sql.DB {
	return container.postgresDBConn
}

func (container DIContainer) GetCustomerEventStore() *postgres.CustomerEventStore {
	if container.customerEventStore == nil {
		container.customerEventStore = postgres.NewCustomerEventStore(
			container.postgresDBConn,
			eventStoreTableName,
			uniqueEmailAddressesTableName,
			container.marshalCustomerEvent,
			container.unmarshalCustomerEvent,
		)
	}

	return container.customerEventStore
}

func (container DIContainer) GetCustomerCommandHandler() *command.CustomerCommandHandler {
	if container.customerCommandHandler == nil {
		container.customerCommandHandler = command.NewCustomerCommandHandler(
			container.GetCustomerEventStore().RetrieveCustomerEventStream,
			container.GetCustomerEventStore().RegisterCustomer,
			container.GetCustomerEventStore().AppendToCustomerEventStream,
			lib.RetryOnConcurrencyConflict,
		)
	}

	return container.customerCommandHandler
}

func (container DIContainer) GetCustomerQueryHandler() *query.CustomerQueryHandler {
	if container.customerQueryHandler == nil {
		container.customerQueryHandler = query.NewCustomerQueryHandler(
			container.GetCustomerEventStore().RetrieveCustomerEventStream,
		)
	}

	return container.customerQueryHandler
}

func (container DIContainer) GetCustomerGRPCServer() customergrpc.CustomerServer {
	if container.customerGRPCServer == nil {
		container.customerGRPCServer = customergrpc.NewCustomerServer(
			container.GetCustomerCommandHandler().RegisterCustomer,
			container.GetCustomerCommandHandler().ConfirmCustomerEmailAddress,
			container.GetCustomerCommandHandler().ChangeCustomerEmailAddress,
			container.GetCustomerCommandHandler().ChangeCustomerName,
			container.GetCustomerCommandHandler().DeleteCustomer,
			container.GetCustomerQueryHandler().CustomerViewByID,
		)
	}

	return container.customerGRPCServer
}
