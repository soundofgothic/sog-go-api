package mcp

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	Version = "0.1.0"
)

type MCPServer struct {
	*server.MCPServer
}

type ResourceHandler interface {
	Handler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)
}

type Resource interface {
	ResourceHandler
	Resource() mcp.Resource
}

type TemplateResource interface {
	ResourceHandler
	Resource() mcp.ResourceTemplate
}

type Tool interface {
	Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	Tool() mcp.Tool
}

type config struct {
	resources         []Resource
	templateResources []TemplateResource
	tools             []Tool
}

type Option func(*config)

func WithResources(resources ...Resource) Option {
	return func(o *config) {
		o.resources = append(o.resources, resources...)
	}
}

func WithTemplateResources(templateResources ...TemplateResource) Option {
	return func(o *config) {
		o.templateResources = append(o.templateResources, templateResources...)
	}
}

func WithTools(tools ...Tool) Option {
	return func(o *config) {
		o.tools = append(o.tools, tools...)
	}
}

func NewMCPServer(options ...Option) *MCPServer {
	config := &config{}

	for _, opt := range options {
		opt(config)
	}

	hooks := &server.Hooks{}
	hooks.AddAfterCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
		slog.Info("MCP CallTool", slog.String("tool", message.Params.Name), slog.Any("req", message.Params.Arguments))
	})

	s := server.NewMCPServer(
		"Sound of Gothic",
		Version,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
		server.WithHooks(hooks),
	)

	for _, resource := range config.resources {
		s.AddResource(resource.Resource(), resource.Handler)
	}

	for _, templateResource := range config.templateResources {
		s.AddResourceTemplate(templateResource.Resource(), templateResource.Handler)
	}

	for _, tool := range config.tools {
		toolDefinition := tool.Tool()
		slog.Info("MCP Tool Registration", slog.String("tool", toolDefinition.Name))
		s.AddTool(toolDefinition, tool.Handler)
	}

	return &MCPServer{s}
}
