package docshandler

import "github.com/curioswitch/go-docs-handler/specification"

type jsonSchemaField struct {
	Description                string                     `json:"description"`
	Ref                        string                     `json:"$ref,omitempty"`
	Type                       string                     `json:"type,omitempty"`
	Enum                       []string                   `json:"enum,omitempty"`
	Properties                 map[string]jsonSchemaField `json:"properties,omitempty"`
	AdditionalPropertiesSchema *jsonSchemaField           `json:"additionalProperties,omitempty"`
	AdditionalPropertiesBool   *bool                      `json:"additionalProperties,omitempty"`
	Items                      *jsonSchemaField           `json:"items,omitempty"`
}

type jsonSchema struct {
	ID                   string                     `json:"$id"`
	Title                string                     `json:"title"`
	Description          string                     `json:"description"`
	Properties           map[string]jsonSchemaField `json:"properties"`
	AdditionalProperties bool                       `json:"additionalProperties"`
	Type                 string                     `json:"type"`
}

func generateJSONSchema(spec *specification.Specification) []jsonSchema {
	g := &jSONSchemaGenerator{
		typeNameToEnum:   map[string]specification.Enum{},
		typeNameToStruct: map[string]specification.Struct{},
	}

	for _, enum := range spec.Enums {
		g.typeNameToEnum[enum.Name] = enum
	}
	for _, s := range spec.Structs {
		g.typeNameToStruct[s.Name] = s
	}

	return g.generate(spec)
}

type jSONSchemaGenerator struct {
	typeNameToEnum   map[string]specification.Enum
	typeNameToStruct map[string]specification.Struct
}

func (g *jSONSchemaGenerator) generate(spec *specification.Specification) []jsonSchema {
	var schemas []jsonSchema

	for _, service := range spec.Services {
		for _, method := range service.Methods {
			schemas = append(schemas, g.generateMethodSchema(method))
		}
	}

	return schemas
}

func (g *jSONSchemaGenerator) generateMethodSchema(method specification.Method) jsonSchema {
	schema := jsonSchema{
		ID:                   method.ID,
		Title:                method.Name,
		Description:          method.DescriptionInfo.DocString,
		AdditionalProperties: false,
		Type:                 "object",
	}

	var methodFields []specification.Field
	visited := map[string]string{}
	currentPath := "#"

	if method.UseParameterAsRoot {
		sig := method.Parameters[0].TypeSignature
		if s, ok := g.typeNameToStruct[sig.Signature()]; ok {
			methodFields = s.Fields
		} else {
			// Couldn't resolve parameter, allow any parameters.
			schema.AdditionalProperties = true
		}
		visited[sig.Signature()] = currentPath
	}

	schema.Properties = g.generateProperties(methodFields, map[string]string{}, "#")

	return schema
}

func (g *jSONSchemaGenerator) generateField(field specification.Field, visited map[string]string, path string) jsonSchemaField {
	schema := jsonSchemaField{
		Description: field.DescriptionInfo.DocString,
	}

	if ref, ok := visited[field.TypeSignature.Signature()]; ok {
		schema.Ref = ref
		return schema
	}

	schema.Type = getSchemaType(field.TypeSignature)
	if field.TypeSignature.Type() == specification.TypeSignatureTypeEnum {
		schema.Enum = g.getEnumType(field.TypeSignature)
	}

	currentPath := path
	if len(field.Name) > 0 {
		currentPath += "/" + field.Name
	}

	// We can only have references to struct types, not primitives
	switch schema.Type {
	case "array", "object":
		visited[field.TypeSignature.Signature()] = currentPath
	}

	switch {
	case field.TypeSignature.Type() == specification.TypeSignatureTypeMap:
		aProps := g.generateMapAdditionalProperties(field, visited, currentPath)
		schema.AdditionalPropertiesSchema = &aProps
	case field.TypeSignature.Type() == specification.TypeSignatureTypeIterable:
		items := g.generateArrayItems(field, visited, currentPath)
		schema.Items = &items
	case schema.Type == "object":
		props, found := g.generateStructProperties(field, visited, currentPath)
		if found {
			allowAdditionalProperties := false
			schema.AdditionalPropertiesBool = &allowAdditionalProperties
			schema.Properties = props
		} else {
			// When we have a struct but the definition cannot not be found, we go ahead
			// and allow all additional properties since it may still be more useful than
			// being completely unusable.
			allowAdditionalProperties := true
			schema.AdditionalPropertiesBool = &allowAdditionalProperties
		}
	}

	return schema
}

func (g *jSONSchemaGenerator) generateProperties(fields []specification.Field, visited map[string]string, path string) map[string]jsonSchemaField {
	properties := make(map[string]jsonSchemaField)

	for _, field := range fields {
		switch field.Location {
		case specification.FieldLocationBody, specification.FieldLocationUnspecified:
			properties[field.Name] = g.generateField(field, visited, path+"/properties")
		}
	}

	return properties
}

func (g *jSONSchemaGenerator) generateMapAdditionalProperties(field specification.Field, visited map[string]string, path string) jsonSchemaField {
	type mt interface {
		KeyType() specification.TypeSignature
		ValueType() specification.TypeSignature
	}

	valueType := field.TypeSignature.(mt).ValueType()
	valueFieldInfo := specification.Field{
		Location:      specification.FieldLocationBody,
		TypeSignature: valueType,
	}

	return g.generateField(valueFieldInfo, visited, path+"/additionalProperties")
}

func (g *jSONSchemaGenerator) generateArrayItems(field specification.Field, visited map[string]string, path string) jsonSchemaField {
	type at interface {
		ItemType() specification.TypeSignature
	}

	itemType := field.TypeSignature.(at).ItemType()
	itemField := specification.Field{
		Location:      specification.FieldLocationBody,
		TypeSignature: itemType,
	}

	return g.generateField(itemField, visited, path+"/items")
}

func (g *jSONSchemaGenerator) generateStructProperties(
	field specification.Field, visited map[string]string, path string) (map[string]jsonSchemaField, bool) {
	if s, ok := g.typeNameToStruct[field.TypeSignature.Signature()]; ok {
		return g.generateProperties(s.Fields, visited, path), true
	}
	return map[string]jsonSchemaField{}, false
}

func (g *jSONSchemaGenerator) getEnumType(t specification.TypeSignature) []string {
	var res []string

	if e, ok := g.typeNameToEnum[t.Signature()]; ok {
		for _, value := range e.Values {
			res = append(res, value.Name)
		}
	}

	return res
}

func getSchemaType(t specification.TypeSignature) string {
	switch t.Type() {
	case specification.TypeSignatureTypeEnum:
		return "string"
	case specification.TypeSignatureTypeIterable:
		return "array"
	case specification.TypeSignatureTypeMap:
		return "object"
	case specification.TypeSignatureTypeBase:
		return getBaseType(t.Signature())
	default:
		return "object"
	}
}

func getBaseType(signature string) string {
	switch signature {
	case "bool", "boolean":
		return "boolean"
	case "short", "number", "float", "double":
		return "number"
	case "i", "i8", "i16", "i32", "i64", "integer", "int",
		"l32", "l64", "long", "long32", "long64", "int32", "int64",
		"uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64",
		"sfixed32", "sfixed64":
		return "integer"
	case "binary", "byte", "bytes", "string":
		return "string"
	default:
		return "object"
	}
}
