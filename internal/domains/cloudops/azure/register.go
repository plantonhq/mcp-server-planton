package azure

import "github.com/modelcontextprotocol/go-sdk/mcp"

func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ListVirtualMachinesTool(), ListVirtualMachinesHandler(serverAddress))
	mcp.AddTool(srv, ListBlobContainersTool(), ListBlobContainersHandler(serverAddress))
}
