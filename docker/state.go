package main

import (
	"io"

	"github.com/lni/dragonboat/v3/statemachine"
)

type State struct {
	NodeID    uint64
	ClusterID uint64

	Stuff string
}

func NewStateMachine(clusterID uint64, nodeID uint64) statemachine.IStateMachine {
	return &State{ClusterID: clusterID,	NodeID: nodeID }
}

func (s *State) Update(data []byte) (statemachine.Result, error) {
	s.Stuff = string(data)
	return statemachine.Result{}, nil
}
func (s *State) Lookup(_ interface{}) (interface{}, error) {
	return s.Stuff, nil
}
func (s *State) SaveSnapshot(io.Writer, statemachine.ISnapshotFileCollection, <-chan struct{}) error {
	return nil
}
func (s *State) RecoverFromSnapshot(io.Reader, []statemachine.SnapshotFile, <-chan struct{}) error {
	return nil
}
func (s *State) Close() error {
	return nil
}
