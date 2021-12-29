package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/client"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/logger"
	"github.com/spf13/viper"
)

var (
	clusterID uint64
	nodeName string
	nodeID uint64
)

var (
	node     *dragonboat.NodeHost
	sess     *client.Session
)

func main() {
	log.Printf("clusterID: %d, nodeName: %s, nodeID: %d", clusterID, nodeName, nodeID)

	go log.Fatal(ServeHTTP())
	// any other master routines
	select{}
}

func init() {
	err := loadConfig()
	fatal(err)

	clusterID = viper.GetUint64("cluster_id")
	nodeName = viper.GetString("node_name")
	nodeID = viper.GetUint64("node_id")

	err = bootStrapRaft()
	fatal(err)
}

func loadConfig() error {
	// parse a config flag when not being envoked by test
	if !strings.HasSuffix(os.Args[0], ".test") {
		configFile := flag.String("c", "", "config file")
		flag.Parse()
		if *configFile != "" {
			viper.SetConfigFile(*configFile)
		} else {
			viper.AddConfigPath(path.Join(".", "config"))
			if env := os.Getenv("ENV"); env != "" {
				viper.SetConfigName(env)
			} else {
				viper.SetConfigName("dev")
			}
		}
	}
	return viper.ReadInConfig()
}

func bootStrapRaft() error {
	logger.GetLogger("raft").SetLevel(logger.WARNING)
	logger.GetLogger("rsm").SetLevel(logger.WARNING)
	logger.GetLogger("transport").SetLevel(logger.WARNING)
	logger.GetLogger("grpc").SetLevel(logger.WARNING)
	
	dataDir := filepath.Join(".", "data", nodeName)
	seeds := viper.GetStringSlice("seed_hosts")

	var err error
	node, err = dragonboat.NewNodeHost(config.NodeHostConfig{
		WALDir:         path.Join(dataDir, "wal"),
		NodeHostDir:    path.Join(dataDir, "lib"),
		RTTMillisecond: 100,
		RaftAddress:    seeds[int(nodeID)-1], // cheap hack for dev
	})
	if err != nil {
		return err
	}

	seedMap := map[uint64]string{}
	for i, v := range seeds {
		seedMap[uint64(i+1)] = v
	}

	err = node.StartCluster(seedMap, false, NewStateMachine, config.Config{
		ClusterID:          clusterID,
		NodeID:             nodeID,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	})
	if err != nil {
		return err
	}

	sess = node.GetNoOPSession(clusterID)

	return nil
}

func fatal(err error) {
	if err != nil {
		_, f, l, _ := runtime.Caller(1)
		log.Fatalf("fatal error thrown by: %s:%d, error: %s", f, l, err)
	}
}
