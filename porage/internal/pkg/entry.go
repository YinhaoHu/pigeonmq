package pkg

import (
	"bytes"
	"encoding/binary"
)

type JournalEntryPayload struct {
	LedgerID uint64
	EntryID  int
	Payload  []byte
}

func noPayloadSize() int {
	return 8 + 8
}

func (j *JournalEntryPayload) Serialize() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, j.LedgerID)
	binary.Write(buf, binary.BigEndian, int64(j.EntryID))
	buf.Write(j.Payload)
	return buf.Bytes()
}

// Size returns the number of bytes of this whole entry.
func (j *JournalEntryPayload) Size() uint64 {
	return uint64(noPayloadSize() + len(j.Payload))
}

func DeserializeJournalEntry(data []byte) *JournalEntryPayload {
	entry := &JournalEntryPayload{}
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.BigEndian, &entry.LedgerID)

	var entryID int64
	binary.Read(buf, binary.BigEndian, &entryID)
	entry.EntryID = int(entryID)
	entry.Payload = make([]byte, len(data)-noPayloadSize())
	buf.Read(entry.Payload)
	return entry
}
