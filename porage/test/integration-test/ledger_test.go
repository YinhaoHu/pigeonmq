package integrationtest_test

import (
	"bytes"
	"fmt"
	"os"
	"path"
	entrylogger "porage/internal/entry_logger"
	"porage/internal/index"
	"porage/internal/journal"
	"porage/internal/ledger"
	"porage/internal/memtable"
	"porage/internal/pkg"
	"porage/internal/recovery"
	"porage/test/utilities"
	"sync"
	"testing"
	"time"

	"github.com/divan/num2words"
	"github.com/fatih/color"
)

// Test Sceanrio:
//  1. Write a large number of entries to the ledger.
//  2. Read the entries back from the ledger in a tailing read and catch up read manner.
//  3. Recovery.
//  4. Delete.
//  5. Journal Trim.

var (
	dataDir = "./_data"
)

func TestNewLedger(t *testing.T) {
	utilities.Logger.Logf("TestNewLedger: Start.")
	const ledgerID = uint64(1)
	const nIterations = 1_000_000
	var nIterationsWord = num2words.Convert(nIterations)

	setCleanEnvironment()
	config, err := pkg.ParseConfigFile("./config.toml")
	if err != nil {
		panic(err)
	}
	setup(config)
	defer clean()

	ledger, err := ledger.NewLedger(ledgerID)
	utilities.Logger.FatalIfErr(err, "Failed to create new ledger: %v", err)

	expectedDb := sync.Map{}

	// Check: Write Entry
	utilities.Logger.Logf("Testing write entries.")
	utilities.Logger.SetProgressLogFrequency(nIterations / 10)
	writeWaitGroup := sync.WaitGroup{}
	writeConsumedTime := time.Duration(0)
	writeTimeBegin := time.Now()
	writeWaitGroup.Add(nIterations)
	for i := 0; i < nIterations; i++ {
		utilities.Logger.ReportProgress(i+1, nIterations)
		entryID := i
		expectedEntryPayload := generatePayloadWithEntryID(entryID)

		go func() {
			entryID, err := ledger.PutEntry(expectedEntryPayload)
			expectedDb.Store(entryID, expectedEntryPayload)
			utilities.Logger.FatalIfErr(err, "Failed to put entry: %v", err)
			writeWaitGroup.Done()
		}()
	}
	utilities.Logger.Logf("Waiting for all writes to complete.")
	writeWaitGroup.Wait()
	writeConsumedTime = time.Since(writeTimeBegin)
	utilities.Logger.Logf("All writes completed.")

	// Check: Get Entry
	utilities.Logger.Logf("Testing read entries.")
	utilities.Logger.SetProgressLogFrequency(nIterations / 10)
	readConsumedTime := time.Duration(0)
	for i := 1; i <= nIterations; i++ {
		utilities.Logger.ReportProgress(i, nIterations)

		entryID := i - 1
		thisReadBeginTime := time.Now()
		entry, err := ledger.GetEntry(entryID)
		readConsumedTime += time.Since(thisReadBeginTime)
		utilities.Logger.FatalIfErr(err, "Failed to get entry: %v", err)

		expectedEntryPayload, ok := expectedDb.Load(entryID)
		if !ok {
			utilities.Logger.Logf("Entry %d not found.", entryID)
			panic("Entry not found.")
		}
		expectPayloadEq(t, expectedEntryPayload.([]byte), entry.Payload)
	}
	utilities.Logger.Logf("All reads completed.")

	// Check: Get non-existent entry
	utilities.Logger.Logf("Testing get non-existent entry.")
	nonExistedEntryID := nIterations + 1
	entry, err := ledger.GetEntry(nonExistedEntryID)
	if entry != nil || err != nil {
		utilities.Logger.Logf("Entry %d should not exist and no error. Err: %v", nonExistedEntryID, err)
		panic("Entry should not exist.")
	}

	// Check: Trim journal
	utilities.Logger.Logf("Testing journal trim.")
	time.Sleep(1 * time.Second)
	files, err := os.ReadDir(config.Journal.StoragePath)
	utilities.Logger.FatalIfErr(err, "Failed to read journal storage directory: %v", err)
	utilities.Logger.Logf("All journal segment files are trimmed.")
	if len(files) != 1 {
		t.Fatalf("Journal storage directory does not only have current segment file. Found %d files.", len(files))
	}

	// Check: Close ledger
	utilities.Logger.Logf("Testing close ledger.")
	err = ledger.Close()
	utilities.Logger.FatalIfErr(err, "Failed to close ledger: %v", err)

	utilities.Logger.Logf("TestLedger: %s | Write Time: %v, Read Time: %v, Number of Iterations: %v(%s)", color.HiGreenString("PASS"), writeConsumedTime,
		readConsumedTime, nIterations, nIterationsWord)
}

func TestRecovery(t *testing.T) {
	utilities.Logger.Logf("TestRecovery: Start.")
	const ledgerID = uint64(2)
	const nIterations = 500_000
	const nJournalSegments = 2

	setCleanEnvironment()
	config, err := pkg.ParseConfigFile("./config.toml")
	if err != nil {
		panic(err)
	}

	// Generate a journal file and ledger file manually.
	ledgerFileName := fmt.Sprintf("ledger_%d", ledgerID)
	ledgerFilePath := path.Join(config.Ledger.StoragePath, ledgerFileName)
	if err := os.MkdirAll(config.Ledger.StoragePath, 0755); err != nil {
		panic(err)
	}
	if _, err := os.Create(ledgerFilePath); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(config.Journal.StoragePath, 0755); err != nil {
		panic(err)
	}
	for segmentID := 0; segmentID < nJournalSegments; segmentID++ {
		// Create journal file.
		utilities.Logger.Logf("Generating journal file %v.", segmentID)
		journalFileName := fmt.Sprintf("%v.journal", time.Now().UnixNano())
		journalFilePath := fmt.Sprintf("%s/%s", config.Journal.StoragePath, journalFileName)
		journalFile, err := os.Create(journalFilePath)
		if err != nil {
			panic(err)
		}

		// Write journal entries.
		utilities.Logger.Logf("Writing journal entries.")
		for i := 0; i < nIterations; i++ {
			entryID := i + nIterations*segmentID
			entryPayload := generatePayloadWithEntryID(entryID)
			journalEntry := journal.NewJournalEntry(&pkg.JournalEntryPayload{
				LedgerID: ledgerID,
				EntryID:  entryID,
				Payload:  entryPayload,
			})
			err := journalEntry.WriteTo(journalFile)
			if err != nil {
				panic(err)
			}
		}
		if err := journalFile.Sync(); err != nil {
			panic(err)
		}
		if err := journalFile.Close(); err != nil {
			panic(err)
		}
	}

	setup(config)
	defer clean()

	// Recover.
	utilities.Logger.Logf("Testing recovery.")
	ledgers, err := recovery.Recover()
	utilities.Logger.FatalIfErr(err, "Failed to recover ledgers: %v", err)

	// Check.
	utilities.Logger.Logf("Testing recovered ledgers.")
	ledger := ledgers[0]
	for i := 0; i < nIterations*nJournalSegments; i++ {
		entryID := i
		entryPayload := generatePayloadWithEntryID(entryID)
		entry, err := ledger.GetEntry(entryID)
		utilities.Logger.FatalIfErr(err, "Failed to get entry: %v", err)
		expectPayloadEq(t, entryPayload, entry.Payload)
	}

	utilities.Logger.Logf("TestRecovery: %s", color.HiGreenString("PASS"))
}

func generatePayloadWithEntryID(entryID int) []byte {
	return []byte(fmt.Sprintf("%24s%d", "xxxxxxxxxxxxxxxxxxxxxxxx", entryID))
}

func expectPayloadEq(t *testing.T, expected, actual []byte) {
	if !bytes.Equal(expected, actual) {
		t.Fatalf("Payloads do not match. Expected(len=%d): %v, got(len=%d): %v", len(expected), expected, len(actual), actual)
	}
}

func setCleanEnvironment() {
	os.RemoveAll(dataDir)
}

func setup(config *pkg.Config) {
	memtable.Startup(&config.Memtable)
	index.Startup(&config.IndexFile)
	entrylogger.Startup(&config.EntryLogger)
	journal.Startup(&config.Journal)
	ledger.Startup(config)
}

func clean() {
	journal.Stop()
	ledger.Stop()
}
