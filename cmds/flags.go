package cmds

import (
	stotypes "github.com/ProtoconNet/mitum-sto/types/sto"
)

type PartitionFlag struct {
	Partition stotypes.Partition
}

func (v *PartitionFlag) UnmarshalText(b []byte) error {
	p := stotypes.Partition(string(b))
	if err := p.IsValid(nil); err != nil {
		return err
	}
	v.Partition = p

	return nil
}

func (v *PartitionFlag) String() string {
	return v.Partition.String()
}
