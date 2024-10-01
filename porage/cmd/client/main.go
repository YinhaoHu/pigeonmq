package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"porage/pkg"

	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

var (
	porageClient *pkg.PorageClient

	commandUsageMapping = map[string]string{
		"create-ledger": "create-ledger <ledger_id>",
		"append-entry":  "append-entry <ledger_id> <payload>",
		"get-entry":     "get-entry <ledger_id> <entry_id>",
		"close-ledger":  "close-ledger <ledger_id>",
		"list-ledgers":  "list-ledgers",
		"list-workers":  "list-workers",
		"ledger-len":    "ledger-len <ledger_id>",
		"help":          "show help information",
		"quit":          "exit the client",
	}
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "porage-client",
		Short: "A client to interact with Porage service",
		Long:  "A simple interactive client for Porage. Purpose is to demonstrate the usage of Porage service.",
		Run:   runInteractivePrompt,
	}

	var serverAddr string
	rootCmd.PersistentFlags().StringVarP(&serverAddr, "server", "s", "localhost:32901", "Address of the Porage server")

	cobra.OnInitialize(func() {
		var err error
		porageClient, err = pkg.NewPorageClient(serverAddr)
		if err != nil {
			fmt.Printf("Failed to connect to Porage server: %v\n", err)
			os.Exit(1)
		}
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runInteractivePrompt(cmd *cobra.Command, args []string) {
	defer porageClient.Close()

	fmt.Println("Welcome to Porage CLI. Type 'quit' to exit.")
	p := prompt.New(
		executeCommand,
		completer,
		prompt.OptionPrefix("-> "),
		prompt.OptionTitle("Porage CLI"),
	)
	p.Run()
}

func executeCommand(input string) {
	input = strings.TrimSpace(input)
	handleCommand(input)
}

func completer(d prompt.Document) []prompt.Suggest {
	commandSuggestions := []prompt.Suggest{}
	for cmd, usage := range commandUsageMapping {
		commandSuggestions = append(commandSuggestions, prompt.Suggest{
			Text:        cmd,
			Description: usage,
		})
	}
	return prompt.FilterHasPrefix(commandSuggestions, d.GetWordBeforeCursor(), true)
}

func handleCommand(input string) {
	if input == "" {
		return
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		fmt.Println("Invalid command")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch parts[0] {
	case "create-ledger":
		handleCreateLedger(parts, ctx)
	case "append-entry":
		handleAppendEntry(parts, ctx)
	case "get-entry":
		handleGetEntry(parts, ctx)
	case "close-ledger":
		handleCloseLedger(parts, ctx)
	case "list-ledgers":
		handleListLedgers(parts, ctx)
	case "list-workers":
		handleListWorkers(parts, ctx)
	case "ledger-len":
		handleLedgerLength(parts, ctx)
	case "help":
		showHelp(parts)
	case "quit":
		fmt.Println("Exiting Porage CLI")
		os.Exit(0)
	default:
		fmt.Printf("Unknown command: %s. Type 'help' for more information.\n", parts[0])
	}
}

func handleCreateLedger(parts []string, ctx context.Context) {
	if !isValidCommandUsageLen(parts, 2) {
		return
	}
	ledgerID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		fmt.Printf("Invalid ledger ID: %v\n", err)
		return
	}
	err = porageClient.CreateLedger(ctx, ledgerID)
	if err != nil {
		fmt.Printf("Failed to create ledger: %v\n", status.Convert(err).Message())
	} else {
		fmt.Println("Ledger created successfully")
	}
}

func handleAppendEntry(parts []string, ctx context.Context) {
	if !isValidCommandUsageLen(parts, 3) {
		return
	}
	ledgerID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		fmt.Printf("Invalid ledger ID: %v\n", err)
		return
	}
	payload := parts[2]
	entryID, err := porageClient.AppendEntryOnLedger(ctx, ledgerID, []byte(payload))
	if err != nil {
		fmt.Printf("Failed to append entry: %v\n", status.Convert(err).Message())
	} else {
		fmt.Printf("Entry appended successfully, Entry ID: %d\n", entryID)
	}
}

func handleGetEntry(parts []string, ctx context.Context) {
	if !isValidCommandUsageLen(parts, 3) {
		return
	}
	ledgerID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		fmt.Printf("Invalid ledger ID: %v\n", err)
		return
	}
	entryID, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Printf("Invalid entry ID: %v\n", err)
		return
	}
	payload, err := porageClient.GetEntryFromLedger(ctx, ledgerID, entryID)
	if err != nil {
		fmt.Printf("Failed to get entry: %v\n", status.Convert(err).Message())
	} else {
		fmt.Printf("Entry retrieved: %s\n", string(payload))
	}
}

func handleCloseLedger(parts []string, ctx context.Context) {
	if !isValidCommandUsageLen(parts, 2) {
		return
	}
	ledgerID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		fmt.Printf("Invalid ledger ID: %v\n", err)
		return
	}
	err = porageClient.DeleteLedger(ctx, ledgerID)
	if err != nil {
		fmt.Printf("Failed to close ledger: %v\n", status.Convert(err).Message())
	} else {
		fmt.Println("Ledger closed successfully")
	}
}

func handleListLedgers(parts []string, ctx context.Context) {
	if !isValidCommandUsageLen(parts, 1) {
		return
	}
	ledgerIDs, err := porageClient.ListLedgers(ctx)
	if err != nil {
		fmt.Printf("Failed to list ledgers: %v\n", status.Convert(err).Message())
		return
	}

	if len(ledgerIDs) == 0 {
		fmt.Println("No ledgers found")
		return
	}

	tableContent := make([][]string, 0, len(ledgerIDs))
	for id, ledgerID := range ledgerIDs {
		tableContent = append(tableContent, []string{strconv.Itoa(id + 1), strconv.FormatUint(ledgerID, 10)})
	}
	renderTable([]string{"ID", "Ledger ID"}, tableContent)
}

func handleListWorkers(parts []string, ctx context.Context) {
	if !isValidCommandUsageLen(parts, 1) {
		return
	}
	workerDescriptions, err := porageClient.GetWorkerDescriptions(ctx)
	if err != nil {
		fmt.Printf("Failed to list workers: %v\n", status.Convert(err).Message())
		return
	}

	if len(workerDescriptions) == 0 {
		fmt.Println("No workers found")
		return
	}

	tableContent := make([][]string, 0, len(workerDescriptions))
	id := 1
	for workerName, description := range workerDescriptions {
		tableContent = append(tableContent, []string{strconv.Itoa(id), workerName, description})
		id++
	}
	renderTable([]string{"ID", "Worker Name", "Description"}, tableContent)
}

func handleLedgerLength(parts []string, ctx context.Context) {
	if !isValidCommandUsageLen(parts, 2) {
		return
	}
	ledgerID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		fmt.Printf("Invalid ledger ID: %v\n", err)
		return
	}
	length, err := porageClient.GetLedgerLength(ctx, ledgerID)
	if err != nil {
		fmt.Printf("Failed to get ledger length: %v\n", status.Convert(err).Message())
	} else {
		fmt.Printf("Ledger length: %d\n", length)
	}
}

func showHelp(parts []string) {
	if !isValidCommandUsageLen(parts, 1) {
		return
	}
	tableContent := make([][]string, 0, len(commandUsageMapping))
	commandId := 1
	for command, usage := range commandUsageMapping {
		tableContent = append(tableContent, []string{strconv.Itoa(commandId), command, usage})
		commandId++
	}
	renderTable([]string{"ID", "Command", "Usage"}, tableContent)
}

func isValidCommandUsageLen(parts []string, expectedLen int) bool {
	if len(parts) != expectedLen {
		fmt.Printf("Usage: %v\n", commandUsageMapping[parts[0]])
		return false
	}
	return true
}

func renderTable(header []string, content [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	for _, v := range content {
		table.Append(v)
	}
	table.Render()
}
