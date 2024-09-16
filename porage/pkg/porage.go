package pkg

import (
	"context"
	pb "porage/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PorageClient is the client struct for Porage.
type PorageClient struct {
	connection *grpc.ClientConn
	rpcClient  pb.PorageServiceClient
}

// NewPorageClient creates a new PorageClient.
func NewPorageClient(serverAddr string) (*PorageClient, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, err
	}
	rpcClient := pb.NewPorageServiceClient(conn)
	client := &PorageClient{
		connection: conn,
		rpcClient:  rpcClient,
	}
	return client, nil
}

// Close closes the connection to the server.
func (c *PorageClient) Close() error {
	return c.connection.Close()
}

// CreateLedger creates a new ledger.
func (c *PorageClient) CreateLedger(ctx context.Context, ledgerID uint64) error {
	_, err := c.rpcClient.CreateLedger(ctx, &pb.CreateLedgerRequest{LedgerId: ledgerID})
	return err
}

// AppendEntryOnLedger appends an entry to a ledger.
func (c *PorageClient) AppendEntryOnLedger(ctx context.Context, ledgerID uint64, payload []byte) (int, error) {
	response, err := c.rpcClient.AppendEntryOnLedger(ctx, &pb.AppendEntryOnLedgerRequest{LedgerId: ledgerID, Payload: payload})
	return int(response.GetEntryId()), err
}

// GetEntryFromLedger gets an entry from a ledger.
func (c *PorageClient) GetEntryFromLedger(ctx context.Context, ledgerID uint64, entryID int) ([]byte, error) {
	response, err := c.rpcClient.GetEntryFromLedger(ctx, &pb.GetEntryFromLedgerRequest{LedgerId: ledgerID, EntryId: int64(entryID)})
	return response.GetPayload(), err
}

// CloseLedger closes a ledger.
func (c *PorageClient) CloseLedger(ctx context.Context, ledgerID uint64) error {
	_, err := c.rpcClient.CloseLedger(ctx, &pb.CloseLedgerRequest{LedgerId: ledgerID})
	return err
}

// ListLedgers returns ID list of all ledgers.
func (c *PorageClient) ListLedgers(ctx context.Context) ([]uint64, error) {
	response, err := c.rpcClient.ListLedgers(ctx, &emptypb.Empty{})
	return response.GetLedgerIds(), err
}

// GetWorkerDescription gets the description of a worker.
func (c *PorageClient) GetWorkerDescriptions(ctx context.Context) (map[string]string, error) {
	response, err := c.rpcClient.ListWorkers(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	workerDescriptions := make(map[string]string)
	for workerName, description := range response.GetWorkers() {
		workerDescriptions[workerName] = description.GetDescription()
	}
	return workerDescriptions, nil
}

// GetLedgerLength gets the length of a ledger.
func (c *PorageClient) GetLedgerLength(ctx context.Context, ledgerID uint64) (int, error) {
	response, err := c.rpcClient.LedgerLength(ctx, &pb.LedgerLengthRequest{LedgerId: ledgerID})
	return int(response.GetLength()), err
}
