package main

import (
	"io"

	"github.com/lni/dragonboat/v3/statemachine"
)

type State struct {
	NodeID    uint64
	ClusterID uint64
	Count     uint64
}

func NewStateMachine(clusterID uint64, nodeID uint64) statemachine.IStateMachine {
	return &State{
		ClusterID: clusterID,
		NodeID:    nodeID,
		Count:     0,
	}
}

func (s *State) Update([]byte) (statemachine.Result, error) {
	return statemachine.Result{}, nil
}
func (s *State) Lookup(interface{}) (interface{}, error) {
	return nil, nil
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
