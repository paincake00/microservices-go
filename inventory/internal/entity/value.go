package entity

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ValueKind int32

const (
	StringValue ValueKind = iota
	Int64Value
	Float64Value
	BoolValue
)

type Value struct {
	kind ValueKind
	data any
}

// Допустимые методы (с допустимыми значениями)

func (v *Value) GetKind() ValueKind {
	return v.kind
}

func (v *Value) GetStringValue() string {
	return v.data.(string)
}

func (v *Value) GetInt64Value() int64 {
	return v.data.(int64)
}

func (v *Value) GetFloat64Value() float64 {
	return v.data.(float64)
}

func (v *Value) GetBoolValue() bool {
	return v.data.(bool)
}

func NewStringValue(value string) *Value {
	return &Value{
		kind: StringValue,
		data: value,
	}
}

func NewInt64Value(value int64) *Value {
	return &Value{
		kind: Int64Value,
		data: value,
	}
}

func NewFloat64Value(value float64) *Value {
	return &Value{
		kind: Float64Value,
		data: value,
	}
}

func NewBoolValue(value bool) *Value {
	return &Value{
		kind: BoolValue,
		data: value,
	}
}

// MarshalBSON - имплементация интерфейса bson.Marshaler
func (v *Value) MarshalBSON() ([]byte, error) {
	return bson.Marshal(
		struct {
			Kind ValueKind `bson:"kind"`
			Data any       `bson:"data"`
		}{
			v.kind,
			v.data,
		},
	)
}

// UnmarshalBSON - имплементация интерфейса bson.Unmarshaler
func (v *Value) UnmarshalBSON(data []byte) error {
	var raw struct {
		Kind ValueKind `bson:"kind"`
		// игнорирует Data при парсинге, кладет сырое значение из BSON
		Data bson.RawValue `bson:"data"`
	}

	err := bson.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	v.kind = raw.Kind

	switch v.GetKind() {
	case StringValue:
		var value string
		if err = raw.Data.Unmarshal(&value); err != nil {
			return err
		}
		v.data = value
	case Int64Value:
		var value int64
		if err = raw.Data.Unmarshal(&value); err != nil {
			return err
		}
		v.data = value
	case Float64Value:
		var value float64
		if err = raw.Data.Unmarshal(&value); err != nil {
			return err
		}
		v.data = value
	case BoolValue:
		var value bool
		if err = raw.Data.Unmarshal(&value); err != nil {
			return err
		}
		v.data = value
	default:
		return ErrValueKind
	}

	return nil
}
