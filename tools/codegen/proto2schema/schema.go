package main

// ProviderSchema is the top-level schema for an OpenMCF provider kind.
// One schema is generated per provider (e.g., AwsAlb, GcpGkeCluster).
// Consumed by the Stage 2 generator and by MCP resource template handlers.
type ProviderSchema struct {
	Name          string       `json:"name"`
	Kind          string       `json:"kind"`
	CloudProvider string       `json:"cloudProvider"`
	APIVersion    string       `json:"apiVersion"`
	Description   string       `json:"description,omitempty"`
	ProtoPackage  string       `json:"protoPackage"`
	ProtoFiles    ProtoFiles   `json:"protoFiles"`
	Spec          SpecSchema   `json:"spec"`
	NestedTypes   []TypeSchema `json:"nestedTypes,omitempty"`
}

// ProtoFiles records which proto source files define this provider.
type ProtoFiles struct {
	API  string `json:"api"`
	Spec string `json:"spec"`
}

// SpecSchema describes the {Kind}Spec message for a provider.
type SpecSchema struct {
	Name   string         `json:"name"`
	Fields []*FieldSchema `json:"fields"`
}

// TypeSchema describes a nested message type referenced by a spec.
type TypeSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	ProtoType   string         `json:"protoType,omitempty"`
	Fields      []*FieldSchema `json:"fields"`
}

// FieldSchema describes a single field in a spec or nested type.
//
// StringValueOrRef fields are simplified to "string" type with referenceKind
// metadata preserved. This keeps the MCP tool schema small while retaining
// cross-resource relationship information for documentation and future use.
type FieldSchema struct {
	Name               string      `json:"name"`
	JSONName           string      `json:"jsonName"`
	ProtoField         string      `json:"protoField"`
	Type               TypeSpec    `json:"type"`
	Description        string      `json:"description,omitempty"`
	Required           bool        `json:"required"`
	Validation         *Validation `json:"validation,omitempty"`
	ReferenceKind      string      `json:"referenceKind,omitempty"`
	ReferenceFieldPath string      `json:"referenceFieldPath,omitempty"`
	Default            string      `json:"default,omitempty"`
	RecommendedDefault string      `json:"recommendedDefault,omitempty"`
	OneofGroup         string      `json:"oneofGroup,omitempty"`
}

// TypeSpec describes the type of a field.
type TypeSpec struct {
	Kind        string    `json:"kind"`
	KeyType     *TypeSpec `json:"keyType,omitempty"`
	ValueType   *TypeSpec `json:"valueType,omitempty"`
	ElementType *TypeSpec `json:"elementType,omitempty"`
	MessageType string    `json:"messageType,omitempty"`
	EnumType    string    `json:"enumType,omitempty"`
	EnumValues  []string  `json:"enumValues,omitempty"`
}

// Validation holds buf.validate rules extracted from proto annotations.
type Validation struct {
	Required  bool     `json:"required,omitempty"`
	MinLength int      `json:"minLength,omitempty"`
	MaxLength int      `json:"maxLength,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
	Min       int      `json:"min,omitempty"`
	Max       int      `json:"max,omitempty"`
	MinItems  int      `json:"minItems,omitempty"`
	MaxItems  int      `json:"maxItems,omitempty"`
	Enum      []string `json:"enum,omitempty"`
	Unique    bool     `json:"unique,omitempty"`
	Const     string   `json:"const,omitempty"`
}

// MetadataSchema describes CloudResourceMetadata fields shared across all providers.
type MetadataSchema struct {
	Name        string       `json:"name"`
	Fields      []*FieldSchema `json:"fields"`
	NestedTypes []TypeSchema `json:"nestedTypes,omitempty"`
}

// Registry indexes all generated provider schemas for quick lookup.
// Used by the Stage 2 generator and MCP resource template handlers.
type Registry struct {
	Providers map[string]RegistryEntry `json:"providers"`
}

// RegistryEntry maps a provider kind to its schema file location.
type RegistryEntry struct {
	CloudProvider string `json:"cloudProvider"`
	APIVersion    string `json:"apiVersion"`
	SchemaFile    string `json:"schemaFile"`
}
