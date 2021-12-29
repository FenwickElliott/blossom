package main

import (
	"encoding/binary"
	"io"
	"log"

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

func (s *State) Update(data []byte) (statemachine.Result, error) {
	// return statemachine.Result{}, nil
	s.Count++
	log.Printf("from ExampleStateMachine.Update(), msg: %s, count:%d\n", string(data), s.Count)
	return statemachine.Result{Value: uint64(len(data))}, nil
}
func (s *State) Lookup(query interface{}) (interface{}, error) {
	// return nil, nil
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, s.Count)
	return result, nil
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
