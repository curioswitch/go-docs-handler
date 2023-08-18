package specification

import (
	"encoding/json"
	"fmt"
)

type TypeSignatureType int

const (
	TypeSignatureTypeBase TypeSignatureType = iota
	TypeSignatureTypeStruct
	TypeSignatureTypeEnum
	TypeSignatureTypeIterable
	TypeSignatureTypeMap
)

type TypeSignature interface {
	Type() TypeSignatureType
	Signature() string
	MarshalJSON() ([]byte, error)
}

func NewStructTypeSignature(name string) TypeSignature {
	return descriptiveTypeSignature{
		type_:     TypeSignatureTypeStruct,
		signature: name,
	}
}

func NewEnumTypeSignature(name string) TypeSignature {
	return descriptiveTypeSignature{
		type_:     TypeSignatureTypeEnum,
		signature: name,
	}
}

type descriptiveTypeSignature struct {
	type_     TypeSignatureType
	signature string
}

func (s descriptiveTypeSignature) Type() TypeSignatureType {
	return s.type_
}

func (s descriptiveTypeSignature) Signature() string {
	return s.signature
}

func (s descriptiveTypeSignature) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Signature())
}

func NewIterableTypeSignature(name string, itemType TypeSignature) TypeSignature {
	return iterableTypeSignature{
		name:     name,
		itemType: itemType,
	}
}

type iterableTypeSignature struct {
	name     string
	itemType TypeSignature
}

func (s iterableTypeSignature) Type() TypeSignatureType {
	return TypeSignatureTypeIterable
}

func (s iterableTypeSignature) Signature() string {
	return fmt.Sprintf("%s<%s>", s.name, s.itemType.Signature())
}

func (s iterableTypeSignature) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Signature())
}

func (s iterableTypeSignature) ItemType() TypeSignature {
	return s.itemType
}

type mapTypeSignature struct {
	keyType   TypeSignature
	valueType TypeSignature
}

func NewMapTypeSignature(keyType TypeSignature, valueType TypeSignature) TypeSignature {
	return mapTypeSignature{
		keyType:   keyType,
		valueType: valueType,
	}
}

func (m mapTypeSignature) Type() TypeSignatureType {
	return TypeSignatureTypeMap
}

func (m mapTypeSignature) Signature() string {
	return fmt.Sprintf("map<%s, %s>", m.keyType.Signature(), m.valueType.Signature())
}

func (m mapTypeSignature) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Signature())
}

func (m mapTypeSignature) KeyType() TypeSignature {
	return m.keyType
}

func (m mapTypeSignature) ValueType() TypeSignature {
	return m.valueType
}
