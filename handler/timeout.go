package handler

import (
	"net/http"
	"time"

	"ginfra/log"

	"github.com/gin-gonic/gin"
)

//TimedHandler 示例
func TimedHandler(c *gin.Context) {

	// get the underlying request context
	ctx := c.Request.Context()

	// create the response data type to use as a channel type
	type responseData struct {
		status int
		body   map[string]interface{}
	}

	// create a done channel to tell the request it's done
	doneChan := make(chan responseData)

	// here you put the actual work needed for the request
	// and then send the doneChan with the status and body
	// to finish the request by writing the response
	go func(c *gin.Context) {
		select {

		// if the context is done it timed out or was cancelled
		// so don't return anything
		case <-ctx.Done():
			c.AbortWithStatus(http.StatusGatewayTimeout)
			log.WithGinContext(c).Info("timeout, terminate sub goroutine...")
			return

		// use timer to simulate I/O or logical opertion
		case <-time.After(time.Second * 30):
			doneChan <- responseData{
				status: 200,
				body:   gin.H{"hello": "world"},
			}
		}
	}(c)

	// non-blocking select on two channels see if the request
	// times out or finishes
	select {

	// if the context is done it timed out or was cancelled
	// so don't return anything
	case <-ctx.Done():
		c.AbortWithStatus(http.StatusGatewayTimeout)
		log.WithGinContext(c).Info("timeout")
		return

	// if the request finished then finish the request by
	// writing the response
	case res := <-doneChan:
		c.JSON(res.status, res.body)
	}
}
