package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
)

const clusterID uint64 = 1

var (
	nodeName string

	seeds = []string{
		"localhost:1111",
		"localhost:2222",
		"localhost:3333",
	}
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("requires ecatly two arguments, received: %d", len(os.Args))
	}

	nodeName = os.Args[1]
	nodeID, err := strconv.ParseUint(nodeName, 10, 64)
	fatal(err)
	nodeName = "node_" +nodeName
	
	log.Printf("starting as: %s", nodeName)

	dataDir := filepath.Join(".", "data", nodeName)

	node, err := dragonboat.NewNodeHost(config.NodeHostConfig{
		// NodeID: nodeID,
		// ClusterID: clusterID,
		WALDir:         path.Join(dataDir, "wal"),
		NodeHostDir:    path.Join(dataDir, "lib"),
		RTTMillisecond: 10,
		RaftAddress:    seeds[int(nodeID)-1], // cheap hack for dev
	})
	fatal(err)

	seedNodes := map[uint64]string{}
	err = node.StartCluster(seedNodes, true, NewStateMachine, config.Config{
		ClusterID:          clusterID,
		NodeID:             nodeID,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	})
	fatal(err)

	select {}
}

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	// // intercept logs
	// logger.GetLogger("raft").SetLevel(logger.ERROR)
	// logger.GetLogger("rsm").SetLevel(logger.WARNING)
	// logger.GetLogger("transport").SetLevel(logger.WARNING)
	// logger.GetLogger("grpc").SetLevel(logger.WARNING)
}

func fatal(err error) {
	if err != nil {
		_, f, l, _ := runtime.Caller(1)
		log.Fatalf("fataal error thrown by: %s:%d, error: %s", f, l, err)
	}
}

func notes() {
	if false {
		// https://github.com/golang/go/issues/17393
		if runtime.GOOS == "darwin" {
			signal.Ignore(syscall.Signal(0xd))
		}
	}
}
