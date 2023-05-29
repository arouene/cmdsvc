package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"cmdsvc/config"
	"cmdsvc/executor"
)

var dispatcher executor.Executor

func init() {
	dispatcher = executor.Default()
}

func main() {
	if !config.General.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// handle CORS, allow all origins
	router.Use(cors.Default())

	// api routing
	api := router.Group("/api")
	{
		// simple health check
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// returns a list of all the jobs actually running
		api.GET("/list", func(c *gin.Context) {
			jobs := dispatcher.GetJob("")
			c.JSON(http.StatusOK, gin.H{
				"jobs": jobs,
			})
		})

		// set custom routes
		for _, service := range config.Services {
			service := service
			log.Printf("Set route: /api%s for service %s", service.Route, service.Name)
			api.POST(service.Route, func(c *gin.Context) {
				execute(c, service.Name)
			})
		}
	}

	log.Println("Listening on", config.General.Address)
	if len(config.General.CrtFile) > 0 && len(config.General.KeyFile) > 0 {
		log.Fatal(router.RunTLS(config.General.Address, config.General.CrtFile, config.General.KeyFile))
	} else {
		log.Fatal(router.Run(config.General.Address))
	}
}

// execute dispatch a command to the executor and then
// returns a status code to the http client
func execute(c *gin.Context, svcName string) {
	// Get service svcName
	if svc, ok := config.Services.Get(svcName); ok {
		// early returns if command is not runnable
		if runnable := isRunnable(c, svc); !runnable {
			return
		}

		// Dispatch execution
		err := dispatcher.Execute(svcName, svc.Command, svc.WorkDir, svc.Environ)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Cannot start command",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "OK",
			})
		}
	}
}

// isRunnable checks if an exclusive command is already running
// an exclusive command cannot run concurrently
func isRunnable(c *gin.Context, svc config.Service) bool {
	if svc.Exclusive {
		jobs := dispatcher.GetJob(svc.Name)
		if len(jobs) > 0 {
			c.AbortWithStatusJSON(http.StatusAccepted, gin.H{
				"message": "Command already running",
			})
			return false
		}
	}
	return true
}
