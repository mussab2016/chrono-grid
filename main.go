package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"unisched/server/bootsrap"
	"unisched/server/db"
	"unisched/server/http"
)

// Embed Angular build output into the final EXE binary
//
//go:embed client/dist
var angularAssets embed.FS

func main() {
	// 1. Start the embedded database
	dbEngine, err := db.InitDB(5433)
	if err != nil {
		log.Fatalf("Failed to start database: %v", err)
	}
	defer dbEngine.Close()

	// 1.5 Seed initial university data (Idempotent - runs safely once)
	//ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//defer cancel()
	//if err := seed.Run(ctx, dbEngine.Client); err != nil {
	//	log.Fatalf("Failed to seed database: %v", err)
	//}
	//dbEngine.GenerateData()

	// 2. Set up the embedded file system for the Angular UI
	distFS, err := fs.Sub(angularAssets, "client/dist/client/browser")
	if err != nil {
		log.Fatalf("Failed to read Angular files: %v", err)
	}
	serverHost := "localhost"
	serverPort := "8080"
	targetURL := fmt.Sprintf("http://localhost:%s", serverPort)

	// 3. Starting services
	services, err := bootsrap.StartServices(dbEngine)
	if err != nil {
		log.Fatalf("Failed to start services: %v", err)
	}

	go func() {
		log.Println("Starting server")
		if err = http.StartHTTPServer(serverHost, serverPort, distFS, services); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
		log.Printf("Server started at http://%s:%s", serverHost, serverPort)
	}()

	// Automatically open the standalone browser window as soon as the server is ready.
	go func() {
		// Wait half a second to confirm the HTTP server has started.
		time.Sleep(500 * time.Millisecond)
		openNativeWindow(targetURL)
	}()

	log.Printf("UniSched App is running natively at: %s\n", targetURL)

	// 5. Listen for termination signals (Ctrl+C) to gracefully shut down Postgres and prevent data corruption
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down the application and database...")
}

// openNativeWindow opens the browser in standalone window mode (App Mode)
func openNativeWindow(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Attempts to open MS Edge or Chrome in standalone app mode without the browser bar
		cmd = exec.Command("cmd", "/c", "start", "msedge", fmt.Sprintf("--app=%s", url))
		if err := cmd.Run(); err != nil {
			// If Edge is not found, opens the default browser
			exec.Command("cmd", "/c", "start", url).Run()
		}
	case "darwin": // macOS
		exec.Command("open", "-a", "Google Chrome", fmt.Sprintf("--app=%s", url)).Run()
	case "linux":
		exec.Command("google-chrome", fmt.Sprintf("--app=%s", url)).Run()
	}
}
