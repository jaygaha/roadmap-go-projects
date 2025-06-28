package handler

import "fmt"

// PrintUsage prints the command-line usage
func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("  broadcast-server start [-port PORT]")
	fmt.Println("  broadcast-server connect [-host HOST] [-port PORT] [-username USERNAME]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  start    Start the broadcast server")
	fmt.Println("  connect  Connect to the broadcast server as a client")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  broadcast-server start")
	fmt.Println("  broadcast-server start -port 6000")
	fmt.Println("  broadcast-server connect")
	fmt.Println("  broadcast-server connect -host localhost -port 6000 -username jay")
}
