package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
)

func dev_visual(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	// get all the policy data in memory held under each namespace name
	policyInMemory := make(map[string][]*policy.Attribute)
	attrList, err := h.ListAttributes(common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE)
	if err != nil {
		cli.ExitWithError("Failed to list attributes", err)
	}
	for _, attr := range attrList {
		ns := attr.GetNamespace().GetName()
		if _, ok := policyInMemory[ns]; !ok {
			policyInMemory[ns] = []*policy.Attribute{}
		}
		policyInMemory[ns] = append(policyInMemory[ns], attr)
	}

	// Create a HTTP server to handle any API calls from the browser
	srv := &http.Server{Addr: ":3000"}
	stop := make(chan os.Signal, 1)

	// Start a web server and serve whatever visualization you want to show
	http.HandleFunc("/visual", func(w http.ResponseWriter, r *http.Request) {
		// Get the hello from the query parameters.
		// code := r.URL.Query().Get("hello")

		respBody := `{"data":"Hello world from CLI server"}`
		// Write the user info to the response.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(respBody)

		// Send a value to the stop channel to simulate the SIGINT signal.
		stop <- syscall.SIGINT
	})
	url := "http://localhost:3000/visual?hello=world"
	fmt.Println(url)
	openBrowser(url)

	// Start the HTTP server in a separate goroutine.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for a SIGINT or SIGTERM signal to shutdown the server.
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down HTTP server...")
	// TODO: resolve the panic that seems to occur every time here
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		cli.ExitWithError("Failed to shutdown HTTP server gracefully: %v", err)
	}
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return fmt.Errorf("failed to open browser: %v", err)
	}

	return nil
}

var devVisualCmd = man.Docs.GetCommand("dev/visual",
	man.WithRun(dev_visual),
)

func init() {
	doc := man.Docs.GetCommand("dev/visual")

	dev_visualCmd := &doc.Command
	devCmd.AddCommand(dev_visualCmd)
}
