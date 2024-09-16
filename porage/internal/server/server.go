package server

import (
	"porage/internal/control"
	entrylogger "porage/internal/entry_logger"
	"porage/internal/index"
	"porage/internal/journal"
	"porage/internal/ledger"
	"porage/internal/memtable"
	"porage/internal/pkg"
	"porage/internal/recovery"
)

// PoraServer is the main server struct.
type PoraServer struct {
	config *pkg.Config

	workerControl *control.WorkerControl
	ledgerControl *control.LedgerControl
	grpcServer    *PorageRPCServiceServer
}

// NewPorageServer creates a new PoraServer with the given config.
func NewPorageServer(config *pkg.Config) *PoraServer {
	return &PoraServer{
		config: config,
	}
}

// Start starts the PoraServer.
//
// This function blocks.
func (ps *PoraServer) Start() {
	ps.startLog()
	ps.startWorkerControl()
	ps.startLedgerControl()
	ps.startLocalStorage()
	ps.startRecovery()
	ps.startRPCServer()
}

// Stop stops the PoraServer gracefully.
func (ps *PoraServer) Stop() {
	ps.grpcServer.stop()
	journal.Stop()
	ledger.Stop()
	pkg.Logger.Infof("Porage server stopped")
}

// startLog starts the log configuration.
func (ps *PoraServer) startLog() {
	logLevel := pkg.Logger.LogLevelFromString(ps.config.Log.Level)
	pkg.Logger.SetLevel(logLevel)
	pkg.Logger.SetOutput(ps.config.Log.Output)
	pkg.Logger.SetWithColor(ps.config.Log.WithColor)

	pkg.Logger.Infof("Starting Porage server: accomplished log configuration")
}

func (ps *PoraServer) startWorkerControl() {
	ps.workerControl = control.NewWorkerControl()
	pkg.Logger.Infof("Starting Porage server: accomplished worker control initialization")
}

// startLedgerControl starts the ledger control.
func (ps *PoraServer) startLedgerControl() {
	ps.ledgerControl = control.NewLedgerControl()
	pkg.Logger.Infof("Starting Porage server: accomplished ledger control initialization")
}

// startLocalStorage starts the local storage.
func (ps *PoraServer) startLocalStorage() {
	memtable.Startup(&ps.config.Memtable)
	index.Startup(&ps.config.IndexFile)
	entrylogger.Startup(&ps.config.EntryLogger)
	journal.Startup(&ps.config.Journal)
	ledger.Startup(ps.config)
	pkg.Logger.Infof("Starting Porage server: accomplished local storage initialization")
}

// startRecovery starts the recovery process.
func (ps *PoraServer) startRecovery() {
	recoveredLedgers, err := recovery.Recover()
	if err != nil {
		pkg.Logger.Fatalf("Failed to recover ledgers: %v", err)
	}
	for _, ledger := range recoveredLedgers {
		ps.ledgerControl.AddLedger(ledger)
	}
	pkg.Logger.Infof("Starting Porage server: accomplished recovery process")
}

// startRPCServer starts the RPC.
//
// This function blocks.
func (ps *PoraServer) startRPCServer() {
	ps.grpcServer = newPorageRPCServiceServer(ps.ledgerControl, ps.workerControl)
	pkg.Logger.Infof("Starting Porage server: accomplished RPC server initialization")
	err := ps.grpcServer.start(ps.config)
	if err != nil {
		pkg.Logger.Errorf("Failed to run gRPC server: %v", err)
	}
}
