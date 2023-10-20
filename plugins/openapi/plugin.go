package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	docshandler "github.com/curioswitch/go-docs-handler"
	"github.com/curioswitch/go-docs-handler/specification"
)

var (
	errOpenAPIV3Only       = errors.New("only OpenAPI v3+ is supported")
	errOperationIdRequired = errors.New("operationId is required")
)

func NewPlugin(spec []byte, opts ...Option) docshandler.Plugin {
	c := newConfig(spec)
	for _, opt := range opts {
		opt(c)
	}
	return &plugin{
		config: c,
	}
}

type plugin struct {
	config *config
}

func (p *plugin) GenerateSpecification() (*specification.Specification, error) {
	loader, err := libopenapi.NewDocumentWithConfiguration(p.config.spec, &datamodel.DocumentConfiguration{
		AvoidIndexBuild: true,
	})
	if err != nil {
		return nil, fmt.Errorf("openapi: parsing spec metadata : %w", err)
	}
	if loader.GetSpecInfo().SpecFormat != datamodel.OAS3 {
		return nil, errOpenAPIV3Only
	}

	doc, errs := loader.BuildV3Model()
	if len(errs) > 0 {
		return nil, fmt.Errorf("openapi: parsing spec: %w", errors.Join(errs...))
	}

	// Because Armeria specification separates structs and enums, we need to look for them first
	// to later resolve refs. We go ahead and also convert the ref to a name to simplify code later.
	schemas := schemas{
		structSchemas:   make(map[string]*base.Schema),
		enumSchemas:     make(map[string]*base.Schema),
		arraySchemas:    make(map[string]*base.Schema),
		exampleRequests: p.config.exampleRequests,
	}

	if c := doc.Model.Components; c != nil {
		for name, schemaP := range c.Schemas {
			schema := schemaP.Schema()
			if schema == nil {
				return nil, fmt.Errorf("openapi: parsing schema: %w", schemaP.GetBuildError())
			}

			switch {
			case len(schema.Type) == 0 || schema.Type[0] == "object":
				if len(schema.Enum) > 0 {
					schemas.enumSchemas[name] = schema
				} else {
					schemas.structSchemas[name] = schema
				}
			case schema.Type[0] == "array":
				schemas.arraySchemas[name] = schema
			}
		}
	}

	if p := doc.Model.Paths; p != nil {
		schemas.pathSchemas = p.PathItems
	}

	if err := schemas.convert(); err != nil {
		return nil, err
	}

	// Assume one service
	svc := specification.Service{
		Name:    "Service",
		Methods: schemas.methods,
	}

	return &specification.Specification{
		Structs:  schemas.structs,
		Enums:    schemas.enums,
		Services: []specification.Service{svc},
	}, nil
}

type schemas struct {
	structSchemas map[string]*base.Schema
	enumSchemas   map[string]*base.Schema
	pathSchemas   map[string]*v3.PathItem
	arraySchemas  map[string]*base.Schema

	exampleRequests map[string][]any

	methods []specification.Method
	structs []specification.Struct
	enums   []specification.Enum
}

func (s *schemas) convert() error {
	for name, schema := range s.structSchemas {
		res, err := s.convertStruct(name, schema)
		if err != nil {
			return err
		}
		s.structs = append(s.structs, res)
	}

	for name, enum := range s.enumSchemas {
		res, err := s.convertEnum(name, enum)
		if err != nil {
			return err
		}
		s.enums = append(s.enums, res)
	}

	for path, schema := range s.pathSchemas {
		res, err := s.convertPath(path, schema)
		if err != nil {
			return err
		}
		s.methods = append(s.methods, res...)
	}

	return nil
}

func (s *schemas) convertPath(path string, schema *v3.PathItem) ([]specification.Method, error) {
	var res []specification.Method

	if o := schema.Get; o != nil {
		m, err := s.convertOperation(path, "GET", schema, o)
		if err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	if o := schema.Post; o != nil {
		m, err := s.convertOperation(path, "POST", schema, o)
		if err != nil {
			return nil, err
		}
		res = append(res, m)
	}

	return res, nil
}

func (s *schemas) convertOperation(path string, httpMethod string, schema *v3.PathItem, operation *v3.Operation) (specification.Method, error) {
	oid := operation.OperationId
	if oid == "" {
		return specification.Method{}, errOperationIdRequired
	}
	var pathParts []string
	regex := false
	for _, p := range strings.Split(path, "/") {
		if p == "" || p[0] != '{' {
			pathParts = append(pathParts, p)
			continue
		}
		regex = true
		pathParts = append(pathParts, "[^/]+")
	}
	path = strings.Join(pathParts, "/")
	if regex {
		path = fmt.Sprintf("regex:%s", path)
	} else {
		path = fmt.Sprintf("exact:%s", path)
	}
	ep := specification.Endpoint{
		PathMapping: path,
		// Always add mime type which triggers Armeria debug form.
		AvailableMimeTypes: []string{"application/json; charset=utf-8"},
	}
	var parameters []specification.Field
	if operation.RequestBody != nil {
		setDefault := false
		for ct, reqSchema := range operation.RequestBody.Content {
			ep.AvailableMimeTypes = append(ep.AvailableMimeTypes, ct)
			if !setDefault {
				ep.DefaultMimeType = ct
				setDefault = true
			}
			if !reqSchema.Schema.IsReference() {
				return specification.Method{}, fmt.Errorf("openapi: non-$ref request body in operation %s", operation.OperationId)
			}
			f, err := s.convertField("body", reqSchema.Schema)
			if err != nil {
				return specification.Method{}, err
			}
			parameters = append(parameters, f)
		}
	}
	for _, param := range operation.Parameters {
		f, err := s.convertField(param.Name, param.Schema)
		if err != nil {
			return specification.Method{}, err
		}
		switch param.In {
		case "path":
			f.Location = specification.FieldLocationPath
		case "query":
			f.Location = specification.FieldLocationQuery
		case "header":
			f.Location = specification.FieldLocationHeader
		}
		parameters = append(parameters, f)
	}

	res := specification.Method{
		Name:               oid,
		ID:                 fmt.Sprintf("Service/%s/%s", oid, httpMethod),
		HTTPMethod:         httpMethod,
		Endpoints:          []specification.Endpoint{ep},
		Parameters:         parameters,
		UseParameterAsRoot: operation.RequestBody != nil,
	}

	if reqs, ok := s.exampleRequests[oid]; ok {
		var serialized []string
		for _, req := range reqs {
			j, err := json.MarshalIndent(req, "", "  ")
			if err != nil {
				return specification.Method{}, fmt.Errorf("openapi: serializing example request: %w", err)
			}
			serialized = append(serialized, string(j))
		}
		res.ExampleRequests = serialized
	}

	if operation.Responses != nil {
		// Armeria only supports a single response type, so we just pick the first one.
		for _, response := range operation.Responses.Codes {
			if response.Content != nil {
				for _, respSchema := range response.Content {
					typ, err := s.fieldTypeSignature(respSchema.Schema)
					if err != nil {
						return specification.Method{}, err
					}
					res.ReturnTypeSignature = typ
				}
			}
			break
		}
	}

	return res, nil
}

func (s *schemas) convertStruct(name string, schema *base.Schema) (specification.Struct, error) {
	res := specification.Struct{
		Name: name,
	}

	if schema.Description != "" {
		res.DescriptionInfo = description(schema.Description)
	}

	props := map[string]*base.SchemaProxy{}
	if len(schema.AllOf) > 0 {
		for _, allOf := range schema.AllOf {
			if allOf.IsReference() {
				refName, ok := strings.CutPrefix(allOf.GetReference(), "#/components/schemas/")
				if !ok {
					return specification.Struct{}, fmt.Errorf("openapi: only component/schemas property refs are currently supported, got %s", refName)
				}
				if _, ok := s.structSchemas[refName]; ok {
					for k, v := range s.structSchemas[refName].Properties {
						props[k] = v
					}
				}
				if _, ok := s.enumSchemas[refName]; ok {
					for k, v := range s.enumSchemas[refName].Properties {
						props[k] = v
					}
				}
			} else {
				for k, v := range allOf.Schema().Properties {
					props[k] = v
				}
			}
		}
	} else {
		for name, prop := range schema.Properties {
			props[name] = prop
		}
	}

	for name, prop := range props {
		field, err := s.convertField(name, prop)
		if err != nil {
			return specification.Struct{}, err
		}
		if slices.Contains(schema.Required, name) {
			field.Requirement = "REQUIRED"
		} else {
			field.Requirement = "OPTIONAL"
		}
		res.Fields = append(res.Fields, field)
	}

	return res, nil
}

func (s *schemas) convertEnum(name string, schema *base.Schema) (specification.Enum, error) {
	res := specification.Enum{
		Name: name,
	}

	if schema.Description != "" {
		res.DescriptionInfo = description(schema.Description)
	}

	for _, v := range schema.Enum {
		vStr, ok := v.(string)
		if !ok {
			// Best effort string conversion, consider changing to an error if it causes problems.
			vStr = fmt.Sprintf("%v", v)
		}
		res.Values = append(res.Values, specification.Value{
			Name: vStr,
		})
	}

	return res, nil
}

func (s *schemas) convertField(name string, schemaP *base.SchemaProxy) (specification.Field, error) {
	res := specification.Field{
		Name: name,
	}

	t, err := s.fieldTypeSignature(schemaP)
	if err != nil {
		return specification.Field{}, err
	}
	res.TypeSignature = t

	if d := schemaP.Schema().Description; d != "" {
		res.DescriptionInfo = description(d)
	}

	return res, nil
}

func (s *schemas) fieldTypeSignature(schemaP *base.SchemaProxy) (specification.TypeSignature, error) {
	if schemaP.IsReference() {
		refName, ok := strings.CutPrefix(schemaP.GetReference(), "#/components/schemas/")
		if !ok {
			return nil, fmt.Errorf("openapi: only component/schemas property refs are currently supported, got %s", refName)
		}
		if _, ok := s.structSchemas[refName]; ok {
			return specification.NewStructTypeSignature(refName), nil
		}
		if _, ok := s.enumSchemas[refName]; ok {
			return specification.NewEnumTypeSignature(refName), nil
		}
		if schema, ok := s.arraySchemas[refName]; ok {
			return s.typeSignature(schema)
		}
		// This shouldn't happen in practice since it is an invalid schema.
		return nil, fmt.Errorf("openapi: unknown ref %s", refName)
	}

	schema := schemaP.Schema()
	if schema == nil {
		return nil, fmt.Errorf("openapi: parsing property schema: %w", schemaP.GetBuildError())
	}
	return s.typeSignature(schema)
}

func (s *schemas) typeSignature(schema *base.Schema) (specification.TypeSignature, error) {
	switch t := schema.Type[0]; t {
	case "boolean":
		return jsonTypeSignatureBoolean, nil
	case "integer":
		unsigned := false
		if m := schema.Minimum; m != nil && *m == 0.0 {
			unsigned = true
		}
		switch schema.Format {
		case "int32":
			if unsigned {
				return jsonTypeSignatureUint32, nil
			} else {
				return jsonTypeSignatureInt32, nil
			}
		case "int64":
			if unsigned {
				return jsonTypeSignatureUint64, nil
			} else {
				return jsonTypeSignatureInt64, nil
			}
		default:
			if unsigned {
				return jsonTypeSignatureUint, nil
			} else {
				return jsonTypeSignatureInt, nil
			}
		}
	case "number":
		switch schema.Format {
		case "float":
			return jsonTypeSignatureFloat, nil
		case "double":
			return jsonTypeSignatureDouble, nil
		default:
			return jsonTypeSignatureNumber, nil
		}
	case "string":
		return jsonTypeSignatureString, nil
	case "object":
		// TODO: Will need to convert to an auto-struct.
		panic("implement me")
	case "array":
		if schema.Items.IsB() {
			// Array of arbitrary objects
			return specification.NewIterableTypeSignature("items",
				specification.NewMapTypeSignature(jsonTypeSignatureUnknown, jsonTypeSignatureUnknown)), nil
		} else {
			item, err := s.fieldTypeSignature(schema.Items.A)
			if err != nil {
				return nil, fmt.Errorf("openapi: parsing array item type: %w", err)
			}
			return specification.NewIterableTypeSignature("items", item), nil
		}
	default:
		return jsonTypeSignatureUnknown, nil
	}
}

func description(d string) specification.DescriptionInfo {
	return specification.DescriptionInfo{
		DocString: d,
		Markup:    "MARKDOWN",
	}
}
