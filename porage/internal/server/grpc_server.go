package server

import (
	"context"
	"fmt"
	"net"
	"porage/internal/control"
	"porage/internal/pkg"
	porage "porage/pkg"
	pb "porage/proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PorageRPCServiceServer struct {
	pb.UnimplementedPorageServiceServer
	grpcServer *grpc.Server

	ledgerControl *control.LedgerControl
	workerControl *control.WorkerControl
}

// newPorageRPCServiceServer creates a new PorageServiceServer.
func newPorageRPCServiceServer(ledgerControl *control.LedgerControl, workerControl *control.WorkerControl) *PorageRPCServiceServer {
	grpcServer := grpc.NewServer()
	return &PorageRPCServiceServer{
		ledgerControl: ledgerControl,
		grpcServer:    grpcServer,
		workerControl: workerControl,
	}
}

// stop stops the gRPC server gracefully.
func (s *PorageRPCServiceServer) stop() {
	s.grpcServer.GracefulStop()
	pkg.Logger.Infof("gRPC server stopped")
}

// start starts the gRPC server.
//
// This function blocks.
func (s *PorageRPCServiceServer) start(config *pkg.Config) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Server.GRPCPort))
	if err != nil {
		return err
	}
	pb.RegisterPorageServiceServer(s.grpcServer, s)

	return s.grpcServer.Serve(listener)
}

// CreateLedger creates a new ledger. Error would be returned if the ledger already exists.
func (s *PorageRPCServiceServer) CreateLedger(ctx context.Context, in *pb.CreateLedgerRequest) (*emptypb.Empty, error) {
	err := s.ledgerControl.CreateLedger(in.LedgerId)
	return nil, err
}

// AppendEntryOnLedger puts an entry on a ledger.
func (s *PorageRPCServiceServer) AppendEntryOnLedger(ctx context.Context, in *pb.AppendEntryOnLedgerRequest) (*pb.AppendEntryOnLedgerResponse, error) {
	ledger := s.ledgerControl.GetLedger(in.LedgerId)
	if ledger == nil {
		return nil, porage.ErrLedgerNotFound
	}
	entryID, err := ledger.PutEntry(in.Payload)
	if err != nil {
		return nil, err
	}
	response := &pb.AppendEntryOnLedgerResponse{
		EntryId: int64(entryID),
	}
	return response, nil
}

// GetEntryFromLedger gets an entry from a ledger.
func (s *PorageRPCServiceServer) GetEntryFromLedger(ctx context.Context, in *pb.GetEntryFromLedgerRequest) (*pb.GetEntryFromLedgerResponse, error) {
	ledger := s.ledgerControl.GetLedger(in.LedgerId)
	if ledger == nil {
		return nil, porage.ErrLedgerNotFound
	}
	entry, err := ledger.GetEntry(int(in.EntryId))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, porage.ErrEntryNotFound
	}
	response := &pb.GetEntryFromLedgerResponse{
		Payload: entry.Payload,
	}
	return response, nil
}

// DeleteLedger closes a ledger.
func (s *PorageRPCServiceServer) DeleteLedger(ctx context.Context, in *pb.DeleteLedgerRequest) (*emptypb.Empty, error) {
	err := s.ledgerControl.RemoveLedger(in.LedgerId)
	return nil, err
}

// LedgerLength returns the length of a ledger.
func (s *PorageRPCServiceServer) LedgerLength(ctx context.Context, in *pb.LedgerLengthRequest) (*pb.LedgerLengthResponse, error) {
	ledger := s.ledgerControl.GetLedger(in.LedgerId)
	if ledger == nil {
		return nil, porage.ErrLedgerNotFound
	}
	length, err := ledger.Length()
	if err != nil {
		return nil, err
	}
	response := &pb.LedgerLengthResponse{
		Length: int64(length),
	}
	return response, nil
}

// ListLedgers returns a list of ledgerIDs in the ledger control.
func (s *PorageRPCServiceServer) ListLedgers(ctx context.Context, in *emptypb.Empty) (*pb.ListLedgersResponse, error) {
	ledgerIDs := s.ledgerControl.ListLedgers()
	response := &pb.ListLedgersResponse{
		LedgerIds: ledgerIDs,
	}
	return response, nil
}

// ListWorkers returns a list of workers in the ledger control.
func (s *PorageRPCServiceServer) ListWorkers(ctx context.Context, in *emptypb.Empty) (*pb.ListWorkersResponse, error) {
	workerDescriptions := s.workerControl.List()
	response := pb.ListWorkersResponse{
		Workers: make(map[string]*pb.WorkerDescription),
	}
	for workerName, description := range workerDescriptions {
		response.Workers[workerName] = description.ToPb()
	}
	return &response, nil
}
