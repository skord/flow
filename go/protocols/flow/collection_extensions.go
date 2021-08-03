package flow

import (
	pb "go.gazette.dev/core/broker/protocol"
)

// GetProjectionByField finds the projection with the given field name, or nil if one does not exist
func GetProjectionByField(field string, projections []Projection) *Projection {
	for p := range projections {
		if projections[p].Field == field {
			return &projections[p]
		}
	}
	return nil
}

// GetProjection finds the projection with the given field name, or nil if one does not exist
func (m *CollectionSpec) GetProjection(field string) *Projection {
	return GetProjectionByField(field, m.Projections)
}

// Validate returns an error if the CollectionSpec is invalid.
func (m *CollectionSpec) Validate() error {
	if err := m.Collection.Validate(); err != nil {
		return pb.ExtendContext(err, "Collection")
	}

	var keyPointers = make(map[string]struct{})

	for i, proj := range m.Projections {
		var err error
		if proj.Field == "" {
			err = pb.NewValidationError("missing field")
		}
		if err != nil {
			return pb.ExtendContext(err, "Projections[%d]", i)
		}

		if proj.IsPrimaryKey {
			keyPointers[proj.Ptr] = struct{}{}
		}
	}

	if m.SchemaUri == "" {
		return pb.NewValidationError("missing schema URI")
	}
	if len(m.KeyPtrs) == 0 {
		return pb.NewValidationError("key pointers are empty")
	}
	for _, p := range m.KeyPtrs {
		if _, ok := keyPointers[p]; !ok {
			return pb.NewValidationError("no keyed projection for key pointer %q", p)
		}
	}

	return nil
}

// Validate returns an error if the DerivationSpec is invalid.
func (m *DerivationSpec) Validate() error {
	if err := m.Collection.Validate(); err != nil {
		return pb.ExtendContext(err, "Collection")
	}
	if m.RegisterSchemaUri == "" {
		return pb.NewValidationError("missing RegisterSchemaUri")
	}
	if len(m.RegisterInitialJson) == 0 {
		return pb.NewValidationError("missing RegisterInitialJson")
	}
	for i, tf := range m.Transforms {
		if err := tf.Validate(); err != nil {
			return pb.ExtendContext(err, "Transform[%d]", i)
		}
	}

	return nil
}

// IsRootDocumentProjection returns true only if this is a projection of the entire document,
// meaning that the json pointer is the empty string.
func (projection *Projection) IsRootDocumentProjection() bool {
	return len(projection.Ptr) == 0
}

// IsSingleType returns true if this inference may only hold a single type besides null For
// example, if the types are ["string", "null"] or just ["string"], then this would return true.
func (i *Inference) IsSingleType() bool {
	var nTypes = 0
	for _, ty := range i.Types {
		if ty != "null" {
			nTypes++
		}
	}
	return nTypes == 1
}

// IsSingleScalarType returns true if this inference may hold a single scalar type besides null.
func (i *Inference) IsSingleScalarType() bool {
	var isScalar = false
	var nTypes = 0
	for _, ty := range i.Types {
		switch ty {
		case "null":
		case "integer", "number", "boolean", "string":
			isScalar = true
			nTypes++
		default:
			nTypes++
		}
	}
	return isScalar && nTypes == 1
}

// Validate returns an error if the Inference is invalid.
func (i *Inference) Validate() error {
	return nil
}
