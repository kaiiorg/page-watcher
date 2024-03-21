package models

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type Page struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name           string
	NormalizedText string
	RawText        string
	Diff           []byte
}

func (p *Page) EncodeDiff(diff []diffmatchpatch.Diff) error {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)

	err := encoder.Encode(diff)
	if err != nil {
		return err
	}

	p.Diff = buffer.Bytes()
	return nil
}

func (p *Page) DecodeDiff() ([]diffmatchpatch.Diff, error) {
	diff := []diffmatchpatch.Diff{}
	buffer := bytes.NewReader(p.Diff)
	decoder := gob.NewDecoder(buffer)

	err := decoder.Decode(&diff)
	if err != nil {
		return nil, err
	}
	return diff, nil
}
