package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"log"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func main() {
	//r, err := reporter.NewGRPCReporter("localhost:11800", reporter.WithAuthentication("xxx"))
	r, err := reporter.NewGRPCReporter("localhost:11800", reporter.WithAuthentication("xxx"))
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer r.Close()
	tracer, err := go2sky.NewTracer("go-example", go2sky.WithReporter(r))
	if err != nil {
		log.Fatalf("crate tracer error: %v \n", err)
	}

	// This for test
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		log.Fatalf("create new local span error %v \n", err)
	}
	span.SetOperationName("invoke data")
	span.Tag("kind", "outer")
	time.Sleep(time.Second)
	subSpan, _, err := tracer.CreateLocalSpan(ctx)
	if err != nil {
		log.Fatalf("create new sub local span error %v \n", err)
	}
	subSpan.SetOperationName("invoke inner")
	subSpan.Log(time.Now(), "inner", "this is right")
	time.Sleep(time.Second)
	subSpan.End()
	time.Sleep(500 * time.Millisecond)
	span.End()
	time.Sleep(time.Second)

	// gin
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.Run("localhost:8088")
}
