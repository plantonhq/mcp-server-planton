package graph

import (
	graphv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/graph/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

var (
	nodeTypeResolver         = domains.NewEnumResolver[graphv1.GraphNode_Type](graphv1.GraphNode_Type_value, "node type", "type_unspecified")
	relationshipTypeResolver = domains.NewEnumResolver[graphv1.GraphRelationship_Type](graphv1.GraphRelationship_Type_value, "relationship type", "type_unspecified")
	changeTypeResolver       = domains.NewEnumResolver[graphv1.GetImpactAnalysisInput_ChangeType](graphv1.GetImpactAnalysisInput_ChangeType_value, "change type", "change_type_unspecified")
)
