package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Fully-qualified proto type names that receive special handling during parsing.
const (
	fqnStringValueOrRef  = "org.openmcf.shared.foreignkey.v1.StringValueOrRef"
	fqnInt32ValueOrRef   = "org.openmcf.shared.foreignkey.v1.Int32ValueOrRef"
	fqnProtobufStruct    = "google.protobuf.Struct"
	fqnProtobufTimestamp = "google.protobuf.Timestamp"
)

// Parser orchestrates proto file parsing and schema extraction.
type Parser struct {
	importPaths []string
	kindEnum    *desc.EnumDescriptor
}

// NewParser creates a Parser with the given proto import paths.
func NewParser(importPaths []string) *Parser {
	return &Parser{importPaths: importPaths}
}

// Init parses the CloudResourceKind enum and caches it for kind name resolution.
// Must be called before ParseProvider.
func (p *Parser) Init() error {
	parser := &protoparse.Parser{
		ImportPaths:           p.importPaths,
		IncludeSourceCodeInfo: true,
	}

	fds, err := parser.ParseFiles("org/openmcf/shared/cloudresourcekind/cloud_resource_kind.proto")
	if err != nil {
		return fmt.Errorf("parsing CloudResourceKind proto: %w", err)
	}

	for _, fd := range fds {
		for _, ed := range fd.GetEnumTypes() {
			if ed.GetName() == "CloudResourceKind" {
				p.kindEnum = ed
				return nil
			}
		}
	}

	return fmt.Errorf("CloudResourceKind enum not found")
}

// ParseProvider parses a single provider's api.proto and spec.proto,
// returning a complete ProviderSchema.
func (p *Parser) ParseProvider(cloudProvider, resourceName string) (*ProviderSchema, error) {
	basePath := fmt.Sprintf("org/openmcf/provider/%s/%s/v1", cloudProvider, resourceName)
	apiProto := filepath.Join(basePath, "api.proto")
	specProto := filepath.Join(basePath, "spec.proto")

	parser := &protoparse.Parser{
		ImportPaths:           p.importPaths,
		IncludeSourceCodeInfo: true,
	}

	fds, err := parser.ParseFiles(apiProto, specProto)
	if err != nil {
		return nil, fmt.Errorf("parsing protos: %w", err)
	}

	fdByName := make(map[string]*desc.FileDescriptor, len(fds))
	for _, fd := range fds {
		fdByName[fd.GetName()] = fd
	}

	apiFD := fdByName[apiProto]
	if apiFD == nil {
		return nil, fmt.Errorf("api.proto not found in parsed descriptors")
	}

	kind, apiVersion, description, err := p.extractResourceInfo(apiFD)
	if err != nil {
		return nil, fmt.Errorf("extracting resource info from api.proto: %w", err)
	}

	specFD := fdByName[specProto]
	if specFD == nil {
		return nil, fmt.Errorf("spec.proto not found in parsed descriptors")
	}

	specMsg := findMessage(specFD, kind+"Spec")
	if specMsg == nil {
		return nil, fmt.Errorf("message %sSpec not found in spec.proto", kind)
	}

	nestedTypes := make(map[string]*TypeSchema)
	fields := p.extractFields(specMsg, nestedTypes)

	schema := &ProviderSchema{
		Name:          kind,
		Kind:          kind,
		CloudProvider: cloudProvider,
		APIVersion:    apiVersion,
		Description:   description,
		ProtoPackage:  apiFD.GetPackage(),
		ProtoFiles: ProtoFiles{
			API:  apiProto,
			Spec: specProto,
		},
		Spec: SpecSchema{
			Name:   specMsg.GetName(),
			Fields: fields,
		},
		NestedTypes: sortedNestedTypes(nestedTypes),
	}

	return schema, nil
}

// ParseMetadata parses the shared CloudResourceMetadata message.
func (p *Parser) ParseMetadata() (*MetadataSchema, error) {
	parser := &protoparse.Parser{
		ImportPaths:           p.importPaths,
		IncludeSourceCodeInfo: true,
	}

	fds, err := parser.ParseFiles("org/openmcf/shared/metadata.proto")
	if err != nil {
		return nil, fmt.Errorf("parsing metadata proto: %w", err)
	}

	for _, fd := range fds {
		msg := findMessage(fd, "CloudResourceMetadata")
		if msg == nil {
			continue
		}

		nestedTypes := make(map[string]*TypeSchema)
		fields := p.extractFields(msg, nestedTypes)

		return &MetadataSchema{
			Name:        msg.GetName(),
			Fields:      fields,
			NestedTypes: sortedNestedTypes(nestedTypes),
		}, nil
	}

	return nil, fmt.Errorf("CloudResourceMetadata message not found")
}

// extractResourceInfo finds the main resource message in api.proto and
// extracts the kind name, apiVersion, and description from buf.validate
// string const rules.
//
// Every OpenMCF api.proto has exactly one message with:
//
//	field 1 "api_version" with (buf.validate.field).string.const
//	field 2 "kind" with (buf.validate.field).string.const
func (p *Parser) extractResourceInfo(fd *desc.FileDescriptor) (kind, apiVersion, description string, err error) {
	for _, msg := range fd.GetMessageTypes() {
		avField := msg.FindFieldByNumber(1)
		kindField := msg.FindFieldByNumber(2)

		if avField == nil || kindField == nil {
			continue
		}
		if avField.GetName() != "api_version" || kindField.GetName() != "kind" {
			continue
		}

		apiVersion = extractStringConst(avField)
		kind = extractStringConst(kindField)

		if kind != "" && apiVersion != "" {
			description = extractComments(msg)
			return kind, apiVersion, description, nil
		}
	}

	return "", "", "", fmt.Errorf("no resource message with api_version/kind const rules found in %s", fd.GetName())
}

// extractFields extracts all field schemas from a message descriptor,
// collecting any nested message types into nestedTypes.
func (p *Parser) extractFields(msg *desc.MessageDescriptor, nestedTypes map[string]*TypeSchema) []*FieldSchema {
	var fields []*FieldSchema
	for _, field := range msg.GetFields() {
		fields = append(fields, p.extractField(field, nestedTypes))
	}
	return fields
}

// extractField extracts a single field's schema, including type information,
// validation rules, and OpenMCF custom options.
func (p *Parser) extractField(field *desc.FieldDescriptor, nestedTypes map[string]*TypeSchema) *FieldSchema {
	mcfOpts := extractOpenMCFFieldOptions(field.GetFieldOptions())
	validation := extractValidation(field)

	fs := &FieldSchema{
		Name:        toPascalCase(field.GetName()),
		JSONName:    field.GetJSONName(),
		ProtoField:  field.GetName(),
		Type:        p.extractTypeSpec(field, nestedTypes),
		Description: extractFieldComments(field),
	}

	if mcfOpts.DefaultKindValue != 0 && p.kindEnum != nil {
		fs.ReferenceKind = p.resolveKindName(mcfOpts.DefaultKindValue)
	}
	if mcfOpts.DefaultKindFieldPath != "" {
		fs.ReferenceFieldPath = mcfOpts.DefaultKindFieldPath
	}
	if mcfOpts.Default != "" {
		fs.Default = mcfOpts.Default
	}
	if mcfOpts.RecommendedDefault != "" {
		fs.RecommendedDefault = mcfOpts.RecommendedDefault
	}

	if validation != nil {
		fs.Validation = validation
		if validation.Required {
			fs.Required = true
		}
	}

	if oo := field.GetOneOf(); oo != nil {
		fs.OneofGroup = oo.GetName()
	}

	return fs
}

// extractTypeSpec determines the full type specification for a field,
// handling maps, repeated fields, and scalar/message types.
func (p *Parser) extractTypeSpec(field *desc.FieldDescriptor, nestedTypes map[string]*TypeSchema) TypeSpec {
	if field.IsMap() {
		keyType := p.extractScalarType(field.GetMapKeyType(), nestedTypes)
		valueType := p.extractScalarType(field.GetMapValueType(), nestedTypes)
		return TypeSpec{
			Kind:      "map",
			KeyType:   &keyType,
			ValueType: &valueType,
		}
	}

	if field.IsRepeated() {
		elemType := p.extractScalarType(field, nestedTypes)
		return TypeSpec{
			Kind:        "array",
			ElementType: &elemType,
		}
	}

	return p.extractScalarType(field, nestedTypes)
}

// extractScalarType extracts the base type for a single field
// (not considering repeated/map wrappers).
func (p *Parser) extractScalarType(field *desc.FieldDescriptor, nestedTypes map[string]*TypeSchema) TypeSpec {
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return TypeSpec{Kind: "string"}
	case descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		return TypeSpec{Kind: "int32"}
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		return TypeSpec{Kind: "uint32"}
	case descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		return TypeSpec{Kind: "int64"}
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		return TypeSpec{Kind: "uint64"}
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return TypeSpec{Kind: "bool"}
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return TypeSpec{Kind: "float"}
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return TypeSpec{Kind: "double"}
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return TypeSpec{Kind: "bytes"}
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		return p.extractMessageType(field, nestedTypes)
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return extractEnumType(field)
	default:
		return TypeSpec{Kind: "string"}
	}
}

// extractMessageType handles message-typed fields with special treatment for
// ValueOrRef types (simplified to their value type), well-known google.protobuf
// types, and provider-specific nested messages.
func (p *Parser) extractMessageType(field *desc.FieldDescriptor, nestedTypes map[string]*TypeSchema) TypeSpec {
	msgType := field.GetMessageType()
	fqn := msgType.GetFullyQualifiedName()

	switch fqn {
	case fqnStringValueOrRef:
		return TypeSpec{Kind: "string"}
	case fqnInt32ValueOrRef:
		return TypeSpec{Kind: "int32"}
	case fqnProtobufStruct:
		return TypeSpec{Kind: "object"}
	case fqnProtobufTimestamp:
		return TypeSpec{Kind: "timestamp"}
	}

	if strings.HasPrefix(fqn, "google.protobuf.") {
		return TypeSpec{Kind: strings.ToLower(strings.TrimPrefix(fqn, "google.protobuf."))}
	}

	typeName := msgType.GetName()
	if _, exists := nestedTypes[typeName]; !exists {
		nestedTypes[typeName] = &TypeSchema{
			Name:        typeName,
			Description: extractComments(msgType),
			ProtoType:   fqn,
			Fields:      p.extractFields(msgType, nestedTypes),
		}
	}

	return TypeSpec{
		Kind:        "message",
		MessageType: typeName,
	}
}

// resolveKindName maps a CloudResourceKind enum integer to its string name.
func (p *Parser) resolveKindName(value int32) string {
	if p.kindEnum == nil {
		return ""
	}
	for _, v := range p.kindEnum.GetValues() {
		if v.GetNumber() == value {
			return v.GetName()
		}
	}
	return ""
}

// --- Validation extraction (buf.validate rules) ---

// extractValidation extracts buf.validate rules from field options.
func extractValidation(field *desc.FieldDescriptor) *Validation {
	opts := field.GetFieldOptions()
	if opts == nil {
		return nil
	}

	ext := proto.GetExtension(opts, validate.E_Field)
	if ext == nil {
		return nil
	}

	fc, ok := ext.(*validate.FieldRules)
	if !ok || fc == nil {
		return nil
	}

	v := &Validation{}
	hasRules := false

	if fc.GetRequired() {
		v.Required = true
		hasRules = true
	}

	if sr := fc.GetString(); sr != nil {
		if c := sr.GetConst(); c != "" {
			v.Const = c
			hasRules = true
		}
		if sr.GetMinLen() > 0 {
			v.MinLength = int(sr.GetMinLen())
			hasRules = true
		}
		if sr.GetMaxLen() > 0 {
			v.MaxLength = int(sr.GetMaxLen())
			hasRules = true
		}
		if sr.GetPattern() != "" {
			v.Pattern = sr.GetPattern()
			hasRules = true
		}
		if len(sr.GetIn()) > 0 {
			v.Enum = sr.GetIn()
			hasRules = true
		}
	}

	if ir := fc.GetInt32(); ir != nil {
		if gte := ir.GetGte(); gte != 0 {
			v.Min = int(gte)
			hasRules = true
		}
		if lte := ir.GetLte(); lte != 0 {
			v.Max = int(lte)
			hasRules = true
		}
		if gt := ir.GetGt(); gt != 0 {
			v.Min = int(gt) + 1
			hasRules = true
		}
		if lt := ir.GetLt(); lt != 0 {
			v.Max = int(lt) - 1
			hasRules = true
		}
	}

	if ir := fc.GetInt64(); ir != nil {
		if gte := ir.GetGte(); gte != 0 {
			v.Min = int(gte)
			hasRules = true
		}
		if lte := ir.GetLte(); lte != 0 {
			v.Max = int(lte)
			hasRules = true
		}
	}

	if rr := fc.GetRepeated(); rr != nil {
		if rr.GetMinItems() > 0 {
			v.MinItems = int(rr.GetMinItems())
			hasRules = true
		}
		if rr.GetMaxItems() > 0 {
			v.MaxItems = int(rr.GetMaxItems())
			hasRules = true
		}
		if rr.GetUnique() {
			v.Unique = true
			hasRules = true
		}
	}

	if mr := fc.GetMap(); mr != nil {
		if mr.GetMinPairs() > 0 {
			v.MinItems = int(mr.GetMinPairs())
			hasRules = true
		}
		if mr.GetMaxPairs() > 0 {
			v.MaxItems = int(mr.GetMaxPairs())
			hasRules = true
		}
	}

	if !hasRules {
		return nil
	}
	return v
}

// extractStringConst extracts a buf.validate string const value from a field.
// Used to read api_version and kind constants from api.proto.
func extractStringConst(field *desc.FieldDescriptor) string {
	opts := field.GetFieldOptions()
	if opts == nil {
		return ""
	}

	ext := proto.GetExtension(opts, validate.E_Field)
	if ext == nil {
		return ""
	}

	fc, ok := ext.(*validate.FieldRules)
	if !ok || fc == nil {
		return ""
	}

	if sr := fc.GetString(); sr != nil {
		return sr.GetConst()
	}
	return ""
}

// --- Enum type extraction ---

func extractEnumType(field *desc.FieldDescriptor) TypeSpec {
	enumDesc := field.GetEnumType()
	fqn := fmt.Sprintf("%s.%s", enumDesc.GetFile().GetPackage(), enumDesc.GetName())

	var values []string
	for _, v := range enumDesc.GetValues() {
		if v.GetNumber() == 0 {
			continue
		}
		values = append(values, v.GetName())
	}

	return TypeSpec{
		Kind:       "enum",
		EnumType:   fqn,
		EnumValues: values,
	}
}

// --- Comment extraction ---

func extractComments(msg *desc.MessageDescriptor) string {
	si := msg.GetSourceInfo()
	if si == nil {
		return ""
	}
	return strings.TrimSpace(si.GetLeadingComments())
}

func extractFieldComments(field *desc.FieldDescriptor) string {
	si := field.GetSourceInfo()
	if si == nil {
		return ""
	}
	return strings.TrimSpace(si.GetLeadingComments())
}

// --- Helpers ---

func findMessage(fd *desc.FileDescriptor, name string) *desc.MessageDescriptor {
	for _, msg := range fd.GetMessageTypes() {
		if msg.GetName() == name {
			return msg
		}
	}
	return nil
}

func sortedNestedTypes(m map[string]*TypeSchema) []TypeSchema {
	if len(m) == 0 {
		return nil
	}
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)

	result := make([]TypeSchema, 0, len(m))
	for _, name := range names {
		result = append(result, *m[name])
	}
	return result
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}
