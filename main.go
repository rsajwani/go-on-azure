package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/fsnotify.v1"
	"github.com/yusufpapurcu/wmi"
)

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Hello from Go and Gin running on Azure App Service",
			"link":  "/json",
		})
	})

	router.GET("/json", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"foo": "bar",
		})
	})

	router.GET("/json", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"foo": "bar",
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

type Win32_DiskDrive struct {
	DeviceID          string
	MediaType         string
	Model             string
	InterfaceType     string
	Partitions        uint32
	Size              uint64
	TotalCylinders    uint64
	TotalHeads        uint64
	TotalSectors      uint64
	TotalTracks       uint64
	TracksPerCylinder uint32
	SectorsPerTrack   uint32
	BytesPerSector    uint32
}

func getDrives() []Win32_DiskDrive {
	var ssdList []Win32_DiskDrive
	//query := "SELECT DeviceID, MediaType, Model, InterfaceType FROM Win32_DiskDrive WHERE MediaType='SSD'"
	query := "SELECT DeviceID, MediaType, Model, InterfaceType FROM Win32_DiskDrive"
	err := Query(query, &ssdList)
	if err != nil {
		fmt.Printf("Query failed: %s/n", err)
		return []Win32_DiskDrive{}
	}

	/*fmt.Println("List of SSDs:")
	for _, ssd := range ssdList {
		fmt.Printf("Device ID: %s\n", ssd.DeviceID)
		fmt.Printf("Model: %s\n", ssd.Model)
		fmt.Printf("Interface Type: %s\n", ssd.InterfaceType)
		fmt.Printf("Media Type: %s\n", ssd.MediaType)
		fmt.Println("----------------------------------------------------")
	}*/
	return ssdList
}

func toJson(ssdList []Win32_DiskDrive) string {
	if len(ssdList) == 0 {
		return "no data"
	}
	var jsonResult string
	for _, ssd := range ssdList {
		fmt.Printf("Device ID: %s\n", ssd.DeviceID)
		fmt.Printf("Model: %s\n", ssd.Model)
		fmt.Printf("Interface Type: %s\n", ssd.InterfaceType)
		fmt.Printf("Media Type: %s\n", ssd.MediaType)
		fmt.Println("----------------------------------------------------")
	}
	return jsonResult
}