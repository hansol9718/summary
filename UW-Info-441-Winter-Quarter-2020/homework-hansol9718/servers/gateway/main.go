package main
import (
	"fmt"
	"os"
	"log"
	"net/http"
	"github.com/UW-Info-441-Winter-Quarter-2020/homework-hansol9718/servers/gateway/handlers"
)

//main is the main entry point for the server
func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}
	tlsKeyPath := os.Getenv("TLSKEY")
    tlsCertPath := os.Getenv("TLSCERT")

	if len(tlsKeyPath) == 0 {
		fmt.Print("empty KeyPath")
		os.Exit(3)
	}
	if len(tlsCertPath) == 0 {
		fmt.Print("empty CertPath")
		os.Exit(2)
	}
	  /*
	- Create a new mux for the web server.
	*/
	mux := http.NewServeMux()
	/*
	- Tell the mux to call your handlers.SummaryHandler function
	  when the "/v1/summary" URL path is requested.
	  */
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)

	  /*
	- Start a web server listening on the address you read from
	  the environment variable, using the mux you created as
	  the root handler. Use log.Fatal() to report any errors
	  that occur when trying to start the web server.
	*/
	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, mux))
}
