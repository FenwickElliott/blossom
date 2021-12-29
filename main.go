package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/logger"
)

const clusterID uint64 = 128

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
	nodeName = "node_" + nodeName

	log.Printf("starting as: %s", nodeName)

	dataDir := filepath.Join(".", "data", nodeName)

	node, err := dragonboat.NewNodeHost(config.NodeHostConfig{
		WALDir:         path.Join(dataDir, "wal"),
		NodeHostDir:    path.Join(dataDir, "lib"),
		RTTMillisecond: 20,
		RaftAddress:    seeds[int(nodeID)-1], // cheap hack for dev
	})
	fatal(err)

	// seedNodes := map[uint64]string{}
	seedNodes := map[uint64]string{1: seeds[0], 2: seeds[1], 3: seeds[2]}
	err = node.StartCluster(seedNodes, false, NewStateMachine, config.Config{
		ClusterID:          clusterID,
		NodeID:             nodeID,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	})
	fatal(err)

	// validate that the cluster is alive before continuing
	{
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // time to
		defer cancel()

		for range time.NewTicker(time.Second).C {
			_, err := node.SyncGetClusterMembership(ctx, clusterID)
			if err != nil {
				log.Println(err)
			} else {
				break
			}
		}
	}

	for t := range time.NewTicker(time.Second).C {
		func() {
			sess := node.GetNoOPSession(clusterID)
			ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
			defer cancel()

			res, err := node.SyncPropose(ctx, sess, []byte(fmt.Sprintf("%s - %s", nodeName, t)))
			if err != nil {
				log.Println(err)
				return
			}

			log.Println(res)
		}()
	}

	select {}
}

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	// // intercept logs

	logger.GetLogger("raft").SetLevel(logger.INFO)
	logger.GetLogger("rsm").SetLevel(logger.WARNING)
	logger.GetLogger("transport").SetLevel(logger.WARNING)
	logger.GetLogger("grpc").SetLevel(logger.WARNING)
}

func fatal(err error) {
	if err != nil {
		_, f, l, _ := runtime.Caller(1)
		log.Fatalf("fatal error thrown by: %s:%d, error: %s", f, l, err)
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
