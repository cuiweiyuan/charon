// Copyright © 2022-2023 Obol Labs Inc. Licensed under the terms of a Business Source License 1.1

package testdata

// Code generated by genssz. DO NOT EDIT.

import (
	ssz "github.com/ferranbt/fastssz"

	"github.com/obolnetwork/charon/app/errors"
	"github.com/obolnetwork/charon/app/z"
)

// HashTreeRootWith ssz hashes the Foo object with a hasher
func (f Foo) HashTreeRootWith(hw ssz.HashWalker) (err error) {
	indx := hw.Index()

	// Field 0: 'ByteList' ssz:"ByteList[32]"
	err = putByteList(hw, []byte(f.ByteList[:]), 32, "ByteList")
	if err != nil {
		return err
	}

	// Field 1: 'Number' ssz:"uint64"
	hw.PutUint64(uint64(f.Number))

	// Field 2: 'Bytes4' ssz:"Bytes4"
	err = putBytesN(hw, []byte(f.Bytes4[:]), 4)
	if err != nil {
		return err
	}

	// Field 3: 'Bytes2' ssz:"Bytes2"
	err = putBytesN(hw, []byte(f.Bytes2[:]), 2)
	if err != nil {
		return err
	}

	// Field 4: 'Bar' ssz:"Composite"
	err = f.Bar.HashTreeRootWith(hw)
	if err != nil {
		return err
	}

	// Field 5: 'Quxes' ssz:"CompositeList[256]"
	{
		listIdx := hw.Index()
		for _, item := range f.Quxes {
			err = item.HashTreeRootWith(hw)
			if err != nil {
				return err
			}
		}

		hw.MerkleizeWithMixin(listIdx, uint64(len(f.Quxes)), uint64(256))
	}

	// Field 6: 'UnixTime' ssz:"uint64"
	hw.PutUint64(uint64(f.UnixTime.Unix()))

	hw.Merkleize(indx)

	return nil
}

// HashTreeRootWith ssz hashes the Bar object with a hasher
func (b Bar) HashTreeRootWith(hw ssz.HashWalker) (err error) {
	indx := hw.Index()

	// Field 0: 'Name' ssz:"ByteList[32]"
	err = putByteList(hw, []byte(b.Name[:]), 32, "Name")
	if err != nil {
		return err
	}

	hw.Merkleize(indx)

	return nil
}

// HashTreeRootWith ssz hashes the Qux object with a hasher
func (q Qux) HashTreeRootWith(hw ssz.HashWalker) (err error) {
	indx := hw.Index()

	// Field 0: 'Number' ssz:"uint64"
	hw.PutUint64(uint64(q.Number))

	hw.Merkleize(indx)

	return nil
}

// putByteList appends a ssz byte list.
// See reference: github.com/attestantio/go-eth2-client/spec/bellatrix/executionpayload_encoding.go:277-284.
func putByteList(h ssz.HashWalker, b []byte, limit int, field string) error {
	elemIndx := h.Index()
	byteLen := len(b)
	if byteLen > limit {
		return errors.Wrap(ssz.ErrIncorrectListSize, "put byte list", z.Str("field", field))
	}
	h.AppendBytes32(b)
	h.MerkleizeWithMixin(elemIndx, uint64(byteLen), uint64(limit+31)/32)

	return nil
}

// putByteList appends b as a ssz fixed size byte array of length n.
func putBytesN(h ssz.HashWalker, b []byte, n int) error {
	if len(b) > n {
		return errors.New("bytes too long", z.Int("n", n), z.Int("l", len(b)))
	}

	h.PutBytes(leftPad(b, n))

	return nil
}

// leftPad returns the byte slice left padded with zero to ensure a length of at least l.
func leftPad(b []byte, l int) []byte {
	for len(b) < l {
		b = append([]byte{0x00}, b...)
	}

	return b
}
