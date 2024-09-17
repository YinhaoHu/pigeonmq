package e2e_test

import (
	"context"
	"fmt"
	"os"
	"porage/internal/pkg"
	"porage/internal/server"
	porage "porage/pkg"
	"porage/test/utilities"
	"sync"
	"testing"
	"time"
)

const (
	serverConfigFilePath              = "./config.toml"
	ledgerID                          = uint64(0)
	nIterations                       = 100000
	rwLogFrequency                    = nIterations / 10
	dataDir                           = "./_data"
	getterGoroutineAssignedEntryCount = 100
)

var (
	serverAddr   string
	porageClient *porage.PorageClient = nil
	expectedDB                        = sync.Map{}
)

func TestE2E(t *testing.T) {
	var err error
	err = os.RemoveAll("./data")
	utilities.Logger.FatalIfErr(err, "Failed to remove data directory")

	// Start a Porage server
	utilities.Logger.Logf("Starting Porage server")
	err = startPorageServerInBackground()
	utilities.Logger.FatalIfErr(err, "Failed to start Porage server")

	// Start a Porage client
	utilities.Logger.Logf("Starting Porage client")
	porageClient, err = porage.NewPorageClient(serverAddr)
	utilities.Logger.FatalIfErr(err, "Failed to start Porage client")

	// Test logic
	ctx := context.Background()
	testCreateLedger(ctx)
	testAppendEntry(ctx)
	testGetEntry(ctx)
	testGetLedgerLength(ctx)
	testListWorkers(ctx)
	testCloseEntry(ctx)
}

func startPorageServerInBackground() error {
	config, err := pkg.ParseConfigFile(serverConfigFilePath)
	if err != nil {
		return err
	}
	server := server.NewPorageServer(config)
	serverAddr = fmt.Sprintf("localhost:%d", config.Server.GRPCPort)
	go server.Start()
	time.Sleep(3 * time.Second)
	return nil
}

func testCreateLedger(ctx context.Context) {
	utilities.Logger.Logf("Testing CreateLedger")
	err := porageClient.CreateLedger(ctx, ledgerID)
	utilities.Logger.FatalIfErr(err, "Failed to create ledger")

	ledgerIDList, err := porageClient.ListLedgers(ctx)
	utilities.Logger.FatalIfErr(err, "Failed to list ledgers")
	if len(ledgerIDList) != 1 && ledgerIDList[0] != ledgerID {
		panic("Failed to list ledgers")
	}
}

func testAppendEntry(ctx context.Context) {
	utilities.Logger.Logf("Testing AppendEntry")
	wg := sync.WaitGroup{}
	wg.Add(nIterations)
	startTime := time.Now()
	utilities.Logger.SetProgressLogFrequency(nIterations / 10)
	for i := 0; i < nIterations; i++ {
		utilities.Logger.ReportProgress(i+1, nIterations)
		payload := generatePayloadWithEntryID(i)
		go func() {
			entryID, err := porageClient.AppendEntryOnLedger(ctx, ledgerID, payload)
			utilities.Logger.FatalIfErr(err, "Failed to append entry")
			expectedDB.Store(entryID, payload)
			wg.Done()
		}()
	}
	wg.Wait()
	utilities.Logger.Logf("AppendEntry done, elapsed time: %v", time.Since(startTime))
}

func getterGoroutine(start int, wg *sync.WaitGroup, ctx context.Context) {
	for i := start; i < start+getterGoroutineAssignedEntryCount; i++ {
		payload, err := porageClient.GetEntryFromLedger(ctx, ledgerID, i)
		utilities.Logger.FatalIfErr(err, "Failed to get entry")
		expectedPayload, ok := expectedDB.Load(i)
		if !ok {
			panic("Failed to get entry")
		}
		if string(payload) != string(expectedPayload.([]byte)) {
			panic("Failed to get entry")
		}
	}
	wg.Done()
}

func testGetEntry(ctx context.Context) {
	utilities.Logger.Logf("Testing GetEntry")
	startTime := time.Now()
	utilities.Logger.SetProgressLogFrequency(nIterations / 10)
	wg := sync.WaitGroup{}
	for i := 0; i < nIterations; i += getterGoroutineAssignedEntryCount {
		utilities.Logger.ReportProgress(i, nIterations)
		wg.Add(1)
		go getterGoroutine(i, &wg, ctx)
	}
	wg.Wait()
	utilities.Logger.Logf("GetEntry done, elapsed time: %v(Each goroutine gets 100 entries)", time.Since(startTime))
}

func testGetLedgerLength(ctx context.Context) {
	utilities.Logger.Logf("Testing GetLedgerLength")
	ledgerLength, err := porageClient.GetLedgerLength(ctx, ledgerID)
	utilities.Logger.FatalIfErr(err, "Failed to get ledger length")
	if ledgerLength != nIterations {
		panic("Failed to get ledger length")
	}
}

func testListWorkers(ctx context.Context) {
	utilities.Logger.Logf("Testing ListWorkers")
	workerDescriptions, err := porageClient.GetWorkerDescriptions(ctx)
	utilities.Logger.FatalIfErr(err, "Failed to list workers")
	if len(workerDescriptions) != 3 {
		// Expected: 1. ledger persistence worker; 2. journal worker; 3. journal trim worker
		panic("Failed to list workers. Expected 3 workers. If there is any missing update in e2e, please update this number.")
	}
}

func testCloseEntry(ctx context.Context) {
	utilities.Logger.Logf("Testing CloseEntry")
	err := porageClient.CloseLedger(ctx, ledgerID)
	utilities.Logger.FatalIfErr(err, "Failed to close ledger")

	ledgerIDList, err := porageClient.ListLedgers(ctx)
	utilities.Logger.FatalIfErr(err, "Failed to list ledgers")
	if len(ledgerIDList) != 0 {
		panic("Failed to list ledgers")
	}
}

func generatePayloadWithEntryID(entryID int) []byte {
	return []byte(fmt.Sprintf("Entry ID: %16d", entryID))
}
