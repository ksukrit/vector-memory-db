package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

// Simple API server to interface with the vector memory database

type NearestNeighborArgs struct {
	Query  Vector `json:"query"`
	K      int    `json:"k"`
	OpCode int    `json:"op_code"`
}

type VectorArgs struct {
	Vector Vector `json:"vector"`
	Id     string `json:"id"`
}

func Cleanup(db *InMemVectorDB) {
	fmt.Println("Cleaning up and exiting...")
	db.DumpFile("db.json")
}

func StartServer(db *InMemVectorDB) {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/get_embedding_size", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"embedding_size": db.GetEmbeddingSize(),
		})
	})

	router.GET("/get_size", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"size": db.Size(),
		})
	})

	router.GET("/get_vector", func(c *gin.Context) {
		key := c.Query("id")

		if key == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "missing  vector id",
			})
			return
		}

		vector := db.Get(key)

		c.JSON(http.StatusOK, gin.H{
			"value": vector,
		})
	})

	router.POST("/insert_vector", func(c *gin.Context) {
		var vector Vector

		var args VectorArgs
		c.BindJSON(&args)
		vector = args.Vector

		if len(vector) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing vector",
			})
			return
		}

		if len(vector) != db.GetEmbeddingSize() {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "vector size does not match the embedding size",
			})
			return
		}

		key := db.Insert(vector)
		c.JSON(http.StatusOK, gin.H{
			"id": key,
		})
	})

	router.POST("/insert_vector_with_id", func(c *gin.Context) {
		var vector Vector
		var args VectorArgs
		c.BindJSON(&args)
		vector = args.Vector

		if len(vector) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing vector",
			})
			return
		}

		key := args.Id

		if key == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "missing  vector id",
			})
			return
		}

		if db.Exists(key) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "vector id already exists",
			})
			return
		}

		if len(vector) != db.GetEmbeddingSize() {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "vector size does not match the embedding size",
			})
			return
		}

		key = db.InsertKey(key, vector)
		c.JSON(http.StatusOK, gin.H{
			"id": key,
		})
	})

	router.POST("/get_nearest_neighbors", func(c *gin.Context) {
		var args NearestNeighborArgs
		c.BindJSON(&args)

		if len(args.Query) == 0 || len(args.Query) != db.GetEmbeddingSize() || args.K <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "missing query vector or invalid K or op_code value",
			})
			return
		}

		neighbors := db.KNearestNeighbor(args.Query, args.K, args.OpCode)
		c.JSON(http.StatusOK, gin.H{
			"neighbors": neighbors,
		})

	})

	router.POST("/update_embedding_size", func(c *gin.Context) {
		var args struct {
			EmbeddingSize int `json:"embedding_size"`
		}
		c.BindJSON(&args)
		if args.EmbeddingSize <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid embedding size",
			})
			return
		}

		db.IncreaseEmbeddingSize(args.EmbeddingSize)
		c.JSON(http.StatusOK, gin.H{
			"message": "updated embedding size",
		})

	})

	router.GET("/dump", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": db.GetDBData(),
		})
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		Cleanup(db)
		os.Exit(1)
	}()

	router.Run(":8080")
}
