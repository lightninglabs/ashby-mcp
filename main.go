package main

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/lightninglabs/ashby-mcp/ashby"
	"github.com/lightninglabs/ashby-mcp/tools"
)

func main() {
	client, err := ashby.NewClientFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "ashby-mcp",
			Version: "0.1.0",
		},
		&mcp.ServerOptions{
			Instructions: "Ashby ATS tools for " +
				"Lightning Labs recruiting. " +
				"Query jobs, applications, " +
				"candidates, and screen " +
				"applicants against hiring " +
				"criteria.",
		},
	)

	handler := tools.NewHandler(client)
	tools.RegisterAll(server, handler)

	if err := server.Run(
		context.Background(), &mcp.StdioTransport{},
	); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
