package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/platform48-functions/metric_test"
)

func main() {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	host := ""
	if os.Getenv("LOCAL_ONLY") == "true" {
		host = "127.0.0.1"
	}
	hostport := host + ":" + port
	fmt.Printf("DEBUG: hostport = %q\n", hostport)
	if err := funcframework.StartHostPort(host, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}
