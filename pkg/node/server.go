package node

import (
	"github.com/gin-gonic/gin"
	"github.com/rusik69/govnocloud/pkg/types"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
)

// Serve serves the node.
func Serve() {
	r := gin.New()
	r.POST("/api/v1/vms", CreateVMHandler)
	r.GET("/api/v1/vm/:id", VMInfoHandler)
	r.DELETE("/api/v1/vm/:id", DeleteVMHandler)
	r.GET("/api/v1/vms", ListVMHandler)
	r.GET("api/v1/vm/start/:id", StartVMHandler)
	r.GET("api/v1/vm/stop/:id", StopVMHandler)
	r.GET("/api/v1/container/:id", GetContainerHandler)
	r.POST("/api/v1/container/create", CreateContainerHandler)
	r.DELETE("/api/v1/container/:id", DeleteContainerHandler)
	r.GET("/api/v1/containers", ListContainersHandler)
	r.POST("/api/v1/files", PostFileHandler)
	r.DELETE("/api/v1/file/:name", DeleteFileHandler)
	r.GET("/api/v1/files", ListFilesHandler)
	r.GET("/api/v1/file/:name", GetFileHandler)
	r.GET("/api/v1/node/stats", StatsHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	logrus.Println("Node is listening on port " + string(types.NodeEnvInstance.Port))
	r.Run(":" + string(types.NodeEnvInstance.Port))
}

// StatsHandler handles the stats request.
func StatsHandler(c *gin.Context) {
	nodeStats, err := GetStats()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, nodeStats)
}

// Get gets the stats.
func GetStats() (types.Stats, error) {
	numCPUs, err := cpu.Counts(true)
	if err != nil {
		return types.Stats{}, err
	}
	mem, err := mem.VirtualMemory()
	if err != nil {
		return types.Stats{}, err
	}
	disk, err := disk.Usage("/")
	if err != nil {
		return types.Stats{}, err
	}
	return types.Stats{
		CPUs: numCPUs,
		MEM:  int64(mem.Total),
		DISK: int64(disk.Total),
	}, nil
}
