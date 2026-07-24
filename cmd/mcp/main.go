package main

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/to4kawa/repomapper-go/internal/mapper"
)

type RepoMapInput struct {
	Path         string `json:"path" jsonschema:"repository path to map"`
	Tokens       int    `json:"tokens,omitempty" jsonschema:"optional max tokens (0 = no limit)"`
	IncludeTests bool   `json:"include_tests,omitempty" jsonschema:"include test symbols (default false)"`
}

type RepoMapOutput struct {
	Map string `json:"map" jsonschema:"generated repository map"`
}

func main() {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "repomapper-go",
			Version: "0.1.0",
		},
		nil,
	)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "repo_map",
		Description: "Generate a concise symbol map of a git repository for LLM context",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in RepoMapInput) (*mcp.CallToolResult, RepoMapOutput, error) {
		if in.Path == "" {
			return nil, RepoMapOutput{}, fmt.Errorf("path is required")
		}

		text, err := mapper.Generate(mapper.Options{
			Path:         in.Path,
			MaxTokens:    in.Tokens,
			IncludeTests: in.IncludeTests,
		})
		if err != nil {
			return nil, RepoMapOutput{}, err
		}

		return nil, RepoMapOutput{Map: text}, nil
	})

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
