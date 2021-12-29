package main

import (
	"context"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/client"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/logger"
)

const clusterID uint64 = 128

var (
	nodeName string
	nodeID   uint64
	node     *dragonboat.NodeHost
	sess     *client.Session

	seeds = []string{
		"localhost:3001",
		"localhost:3002",
		"localhost:3003",
	}
	externalPorts = []string{
		":2001",
		":2002",
		":2003",
	}
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("requires ecatly two arguments, received: %d", len(os.Args))
	}

	nodeName = os.Args[1]
	var err error
	nodeID, err = strconv.ParseUint(nodeName, 10, 64)
	fatal(err)
	nodeName = "node_" + nodeName

	log.Printf("starting as: %s", nodeName)

	dataDir := filepath.Join(".", "data", nodeName)

	node, err = dragonboat.NewNodeHost(config.NodeHostConfig{
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

	sess = node.GetNoOPSession(clusterID)

	log.Fatal(ServeHTTP())
}

func init() {
	logger.GetLogger("raft").SetLevel(logger.WARNING)
	logger.GetLogger("rsm").SetLevel(logger.WARNING)
	logger.GetLogger("transport").SetLevel(logger.WARNING)
	logger.GetLogger("grpc").SetLevel(logger.WARNING)
	gin.SetMode(gin.ReleaseMode)
}

func fatal(err error) {
	if err != nil {
		_, f, l, _ := runtime.Caller(1)
		log.Fatalf("fatal error thrown by: %s:%d, error: %s", f, l, err)
	}
}
