package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
   // "encoding/json"
   "io"
   "strconv"
   log "github.com/sirupsen/logrus"


	"github.com/gin-gonic/gin"
	"gopkg.in/fsnotify.v1"
)

type Content struct {
    test string
  
}

func main() {
	var request_nr int
	f, err := os.OpenFile("log.txt", os.O_WRONLY | os.O_CREATE, 0755)
if err != nil {
}
log.SetOutput(f)
	router := gin.Default()


	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Hello from Go and Gin running on Azure App Service Demo",
			"link":  "/json",
		})
	})

	router.POST("/test", func(ctx *gin.Context) {

		body, _ := io.ReadAll(ctx.Request.Body)

		log.Println("request no",strconv.Itoa(request_nr))
		request_nr=request_nr + 1
		var content  string
		content=ctx.Query("validationToken")

		log.Println("body",string(body))

	   if len(content)>81 {

		log.Println("query",string(content))
	   ctx.Data(http.StatusOK,"text/plain; charset=utf-8", []byte(content))
	   } else {
		
		ctx.JSON(http.StatusOK, "no content length")

	   }
	})

	router.GET("/json", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"foo": "baaa",
		})
	})

	router.Static("/public", "./public")

	// creates a new file watcher for App_offline.htm
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	// watch for App_offline.htm and exit the program if present
	// This allows continuous deployment on App Service as the .exe will not be
	// terminated otherwise
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if strings.HasSuffix(event.Name, "app_offline.htm") {
					fmt.Println("Exiting due to app_offline.htm being present")
					os.Exit(0)
				}
			}
		}
	}()

	// get the current working directory and watch it
	currentDir, err := os.Getwd()
	if err := watcher.Add(currentDir); err != nil {
		fmt.Println("ERROR", err)
	}

	// Azure App Service sets the port as an Environment Variable
	// This can be random, so needs to be loaded at startup
	port := os.Getenv("HTTP_PLATFORM_PORT")

	// default back to 8080 for local dev
	if port == "" {
		port = "8080"
	}

	router.Run("127.0.0.1:" + port)
}
