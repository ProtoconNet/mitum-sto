package sto

import (
	"regexp"

	"github.com/ProtoconNet/mitum2/util"
)

var (
	MinLengthPartition = 3
	MaxLengthPartition = 10
	ReValidPartition   = regexp.MustCompile(`^[A-Z0-9][A-Z0-9_\.\!\$\*\@]*[A-Z0-9]$`)
)

type Partition string

func (p Partition) Bytes() []byte {
	return []byte(p)
}

func (p Partition) String() string {
	return string(p)
}

func (p Partition) IsValid([]byte) error {
	if l := len(p); l < MinLengthPartition || l > MaxLengthPartition {
		return util.ErrInvalid.Errorf(
			"invalid length of partition, %d <= %d <= %d", MinLengthPartition, l, MaxLengthPartition)
	} else if !ReValidPartition.Match([]byte(p)) {
		return util.ErrInvalid.Errorf("wrong partition, %q", p)
	}

	return nil
}
