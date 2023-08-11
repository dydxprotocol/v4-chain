package types

import "github.com/dydxprotocol/v4/lib"

var (
	_ RandomBooler = &RealRandomBooler{}
	_ RandomBooler = &AlwaysTrueRandomBooler{}
	_ RandomBooler = &AlwaysFalseRandomBooler{}
)

type RandomBooler interface {
	RandomBool() bool
}

type RealRandomBooler struct{}

func (r *RealRandomBooler) RandomBool() bool {
	return lib.RandomBool()
}

type AlwaysTrueRandomBooler struct{}

func (r *AlwaysTrueRandomBooler) RandomBool() bool {
	return true
}

type AlwaysFalseRandomBooler struct{}

func (r *AlwaysFalseRandomBooler) RandomBool() bool {
	return false
}
