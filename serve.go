package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

func ServeHTTP() error {

	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
	r.Use(gin.Logger())

	// var stuff string
	// r.GET("/stuff", func(c *gin.Context) {
	// 	c.String(200, stuff)
	// })
	// r.PUT("/stuff/:msg", func(c *gin.Context) {
	// 	stuff = c.Param("msg")
	// 	c.JSON(200, gin.H{"message": fmt.Sprintf("set stuff to: %s", stuff)})
	// })

	r.GET("/stuff", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, time.Second)
		defer cancel()


		res, err := node.SyncRead(ctx, clusterID, nil)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"message": "failed to lookup message", "error": err.Error()})
			return
		}

		msg, ok := res.(string)
		if !ok {
				c.AbortWithStatusJSON(500, gin.H{"message": fmt.Sprintf("typecasting error, expected string, got %s", reflect.TypeOf(res))})
				return
		}

		// log.Println(res)
		c.String(200, msg)
	})

	r.PUT("/stuff/:msg", func(c *gin.Context) {

		msg := c.Param("msg")

		ctx, cancel := context.WithTimeout(c, time.Second)
		defer cancel()

		_, err := node.SyncPropose(ctx, sess, []byte(msg))
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"message": "failed to write message", "error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": fmt.Sprintf("set stuff to: %s", msg)})
	})

	port := externalPorts[int(nodeID)-1] // cheap hack for dev
	log.Println("blossom-external listening on", port)
	return r.Run(port)
}
