package main

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Proto field numbers for OpenMCF custom extensions.
//
// These extensions are defined in openmcf protos but are not compiled into
// our Go binary, so they appear as unknown fields in FieldOptions.
// We extract them using raw protowire decoding.
const (
	// org/openmcf/shared/foreignkey/v1/foreign_key.proto
	fieldNumDefaultKind          protowire.Number = 200001
	fieldNumDefaultKindFieldPath protowire.Number = 200002

	// org/openmcf/shared/options/options.proto
	fieldNumDefault            protowire.Number = 60001
	fieldNumRecommendedDefault protowire.Number = 60002
)

// OpenMCFFieldOptions holds custom option values extracted from a proto field.
type OpenMCFFieldOptions struct {
	DefaultKindValue     int32  // Raw CloudResourceKind enum value (varint)
	DefaultKindFieldPath string // Output field path for cross-resource references
	Default              string // Default value for the field
	RecommendedDefault   string // Recommended default value
}

// extractOpenMCFFieldOptions reads OpenMCF custom extension values from
// proto field options using raw protowire decoding.
//
// OpenMCF extensions are not in our Go proto registry (we don't compile
// the openmcf Go stubs), so they appear as unknown wire bytes in
// descriptorpb.FieldOptions. This function walks the unknown bytes
// and extracts values by their well-known field numbers.
func extractOpenMCFFieldOptions(opts *descriptorpb.FieldOptions) OpenMCFFieldOptions {
	if opts == nil {
		return OpenMCFFieldOptions{}
	}

	raw := opts.ProtoReflect().GetUnknown()
	return decodeOpenMCFOptions(raw)
}

// decodeOpenMCFOptions walks raw protobuf wire bytes and extracts
// OpenMCF extension values by field number.
func decodeOpenMCFOptions(raw []byte) OpenMCFFieldOptions {
	var result OpenMCFFieldOptions

	for len(raw) > 0 {
		num, typ, tagLen := protowire.ConsumeTag(raw)
		if tagLen < 0 {
			break
		}
		raw = raw[tagLen:]

		switch typ {
		case protowire.VarintType:
			v, n := protowire.ConsumeVarint(raw)
			if n < 0 {
				return result
			}
			if num == fieldNumDefaultKind {
				result.DefaultKindValue = int32(v)
			}
			raw = raw[n:]

		case protowire.Fixed32Type:
			if len(raw) < 4 {
				return result
			}
			raw = raw[4:]

		case protowire.Fixed64Type:
			if len(raw) < 8 {
				return result
			}
			raw = raw[8:]

		case protowire.BytesType:
			v, n := protowire.ConsumeBytes(raw)
			if n < 0 {
				return result
			}
			switch num {
			case fieldNumDefaultKindFieldPath:
				result.DefaultKindFieldPath = string(v)
			case fieldNumDefault:
				result.Default = string(v)
			case fieldNumRecommendedDefault:
				result.RecommendedDefault = string(v)
			}
			raw = raw[n:]

		default:
			// Unknown wire type; cannot safely continue parsing.
			return result
		}
	}

	return result
}
