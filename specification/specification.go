package specification

import "encoding/json"

type Array[T any] []T

func (a Array[T]) MarshalJSON() ([]byte, error) {
	if a == nil {
		a = make([]T, 0)
	}
	return json.Marshal([]T(a))
}

type FieldLocation int

const (
	FieldLocationUnspecified FieldLocation = iota
	FieldLocationPath
	FieldLocationHeader
	FieldLocationQuery
	FieldLocationBody
)

func (l FieldLocation) MarshalJSON() ([]byte, error) {
	var loc string
	switch l {
	case FieldLocationUnspecified:
		loc = "UNSPECIFIED"
	case FieldLocationPath:
		loc = "PATH"
	case FieldLocationHeader:
		loc = "HEADER"
	case FieldLocationQuery:
		loc = "QUERY"
	case FieldLocationBody:
		loc = "BODY"

	}
	return json.Marshal(loc)
}

type DescriptionInfo struct {
	DocString string `json:"docString"`
	Markup    string `json:"markup"`
}

type Endpoint struct {
	HostnamePattern    string        `json:"hostnamePattern"`
	PathMapping        string        `json:"pathMapping"`
	DefaultMimeType    string        `json:"defaultMimeType,omitempty"`
	AvailableMimeTypes Array[string] `json:"availableMimeTypes"`
	RegexPathPrefix    string        `json:"regexPathPrefix,omitempty"`
	Fragment           string        `json:"fragment,omitempty"`
}

type Field struct {
	Name            string          `json:"name"`
	Location        FieldLocation   `json:"location"`
	Requirement     string          `json:"requirement"`
	TypeSignature   TypeSignature   `json:"typeSignature"`
	DescriptionInfo DescriptionInfo `json:"descriptionInfo"`
}

type Struct struct {
	Name            string          `json:"name"`
	Alias           string          `json:"alias,omitempty"`
	Fields          Array[Field]    `json:"fields"`
	DescriptionInfo DescriptionInfo `json:"descriptionInfo"`
}

type Method struct {
	Name                    string                   `json:"name"`
	ID                      string                   `json:"id"`
	ReturnTypeSignature     TypeSignature            `json:"returnTypeSignature"`
	Parameters              Array[Field]             `json:"parameters"`
	UseParameterAsRoot      bool                     `json:"-"`
	ExceptionTypeSignatures Array[TypeSignature]     `json:"exceptionTypeSignatures"`
	Endpoints               Array[Endpoint]          `json:"endpoints"`
	ExampleHeaders          Array[map[string]string] `json:"exampleHeaders"`
	ExampleRequests         Array[string]            `json:"exampleRequests"`
	ExamplePaths            Array[string]            `json:"examplePaths"`
	ExampleQueries          Array[string]            `json:"exampleQueries"`
	HTTPMethod              string                   `json:"httpMethod"`
	DescriptionInfo         DescriptionInfo          `json:"descriptionInfo"`
}

type Service struct {
	Name            string                   `json:"name"`
	Methods         Array[Method]            `json:"methods"`
	ExampleHeaders  Array[map[string]string] `json:"exampleHeaders"`
	DescriptionInfo DescriptionInfo          `json:"descriptionInfo"`
}

type Value struct {
	Name            string          `json:"name"`
	IntValue        *int            `json:"intValue,omitempty"`
	DescriptionInfo DescriptionInfo `json:"descriptionInfo"`
}

type Enum struct {
	Name            string          `json:"name"`
	Values          Array[Value]    `json:"values"`
	DescriptionInfo DescriptionInfo `json:"descriptionInfo"`
}

type Specification struct {
	Services       Array[Service]           `json:"services"`
	Enums          Array[Enum]              `json:"enums"`
	Structs        Array[Struct]            `json:"structs"`
	Exceptions     Array[Struct]            `json:"exceptions"`
	ExampleHeaders Array[map[string]string] `json:"exampleHeaders"`
}
