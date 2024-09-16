package pkg

import (
	"os"

	"github.com/pelletier/go-toml"
)

type JournalConfig struct {
	// StoragePath is the path to the directory where the journal files are stored
	StoragePath string `toml:"storage_path"`
	// SegmentSoftThreshold is the threshold at which the journal storage will create a new segment
	// for journal.
	SegmentSoftThreshold uint64 `toml:"segment_soft_threshold"`

	MessageBufferSize          uint64 `toml:"message_buffer_size"`
	MessageBufferBusyThreshold uint64 `toml:"message_buffer_busy_threshold"`
	// GroupCommitThreasold is the number of messages threshold at which the journal storage will commit the entries.
	GroupCommitThreasold uint64 `toml:"group_commit_threshold"`
	// GroupCommitInterval is the time interval (in millsecond) at which the journal storage will commit the entries when the group commit
	// threshold is not reached.
	GroupCommitInterval uint64 `toml:"group_commit_interval"`
	// TrimInterval is the time interval (in second) at which the journal storage will trim the entries.
	TrimInterval uint64 `toml:"trim_interval"`
}

type MemtableConfig struct {
	// TrimThreshold is the threshold at which the ledger will trim the entries in memtable
	TrimThreshold int `toml:"trim_threshold"`
}

type EntryLoggerConfig struct {
	// StoragePath is the path to the directory where the entry logger files are stored
	StoragePath string `toml:"storage_path"`
	// MessageBufferSize is the size of the message buffer
	MessageBufferSize uint64 `toml:"message_buffer_size"`
	// MessageBufferBusyThreshold is the threshold at which the message buffer is considered busy. When the threshold
	// is reached, the entry logger will stop accepting new messages and error will be returned.
	MessageBufferBusyThreshold uint64 `toml:"message_buffer_busy_threshold"`
	// FlushRate is the rate at which the entry logger will flush the message buffer. After the number of processed
	// entries reaches the flush rate, the entry logger will flush the file system.
	FlushRate uint64 `toml:"flush_rate"`
	// FlushInterval is the time interval (in second) at which the entry logger will flush the message buffer. After the time
	// interval reaches the flush interval, the entry logger will flush the file system.
	FlushInterval uint64 `toml:"flush_interval"`
}

type IndexFileConfig struct {
	StoragePath  string `toml:"storage_path"`
	MemtableSize int64  `toml:"memtable_size"`
}

type LedgerConfig struct {
	StoragePath string `toml:"storage_path"`
}

type LogConfig struct {
	Level     string `toml:"level"`
	Output    string `toml:"output"`
	WithColor bool   `toml:"with_color"`
}

type ServerConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	GRPCPort int    `toml:"grpc_port"`
}

type Config struct {
	Ledger      LedgerConfig      `toml:"Ledger"`
	Journal     JournalConfig     `toml:"Journal"`
	Memtable    MemtableConfig    `toml:"Memtable"`
	EntryLogger EntryLoggerConfig `toml:"EntryLogger"`
	IndexFile   IndexFileConfig   `toml:"IndexFile"`
	Log         LogConfig         `toml:"Log"`
	Server      ServerConfig      `toml:"Server"`
}

func ParseConfigFile(filePath string) (*Config, error) {
	// Open the TOML file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the TOML file into a Config struct
	var config Config
	err = toml.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
