package index

import (
	"encoding/binary"
	"errors"
	"os"
	"porage/internal/pkg"
	"strconv"

	"github.com/dgraph-io/badger/v3"
)

type Index struct {
	ledgerID uint64

	db *badger.DB
}

func NewIndex(ledgerID uint64) (*Index, error) {
	// Open the Badger database
	dbStoragePath := makeStoragePathByLedgerID(ledgerID)
	option := badger.DefaultOptions(dbStoragePath)
	option.MemTableSize = myConfig.MemtableSize
	option.Logger = pkg.Logger
	db, err := badger.Open(option)

	if err != nil {
		return nil, err
	}
	return &Index{
		ledgerID: ledgerID,
		db:       db,
	}, nil
}

// Put writes the index value to the badger database.
func (i *Index) Put(entryID int, value *IndexValue) error {
	// Write entryID as key and offset as value to myBadger.
	err := i.db.Update(func(txn *badger.Txn) error {
		key := []byte(strconv.Itoa(entryID))
		value := value.serialize()
		err := txn.Set(key, value)
		return err
	})
	return err
}

// Get the index of the entryID in the ledgerID. If the entryID does not exist, return nil.
func (i *Index) Get(entryID int) (*IndexValue, error) {
	var value *IndexValue = nil
	err := i.db.View(func(txn *badger.Txn) error {
		key := []byte(strconv.Itoa(entryID))
		item, err := txn.Get(key)
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		value = &IndexValue{}
		err = item.Value(func(val []byte) error {
			value.deserialize(val)
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}

// LastItem returns the last entryID and index value in the ledger. If there is no item, indexValue
// will be nil.
//
// Expected to be called in recovery.
func (i *Index) LastItem() (entryID int, indexValue *IndexValue, err error) {
	err = i.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()
		it.Rewind()
		if !it.Valid() {
			return nil
		}
		item := it.Item()
		entryID, err = strconv.Atoi(string(item.Key()))
		if err != nil {
			return err
		}
		indexValue = &IndexValue{}
		err = item.Value(func(val []byte) error {
			indexValue.deserialize(val)
			return nil
		})
		return err
	})
	return
}

// Delete deletes the index file of the ledger.
func (i *Index) Delete() error {
	storagePath := makeStoragePathByLedgerID(i.ledgerID)
	err := os.RemoveAll(storagePath)
	return err
}

type IndexValue struct {
	Offset int
	Size   int
}

func (iv *IndexValue) deserialize(data []byte) {
	iv.Offset = int(binary.BigEndian.Uint64(data[:8]))
	iv.Size = int(binary.BigEndian.Uint64(data[8:]))
}

func (iv *IndexValue) serialize() []byte {
	data := make([]byte, 16)
	binary.BigEndian.PutUint64(data[:8], uint64(iv.Offset))
	binary.BigEndian.PutUint64(data[8:], uint64(iv.Size))
	return data
}
