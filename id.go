package souvenir

import (
	"crypto/rand"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// type

type Type interface {
	Prefix() string
}

type AnyType struct{}

func (a AnyType) Prefix() string {
	return ""
}

func TypePrefix[T Type]() string {
	var t T
	return t.Prefix()
}

// encoding

var zero = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// id

type ID[T Type] struct {
	data [16]byte
}

func ZeroID[T Type]() ID[T] {
	return ID[T]{zero}
}

func NewID[T Type](data [16]byte) ID[T] {
	return ID[T]{data}
}

func ParseID[T Type](data string) (ID[T], error) {
	parts := strings.SplitN(data, "_", 2)

	if parts[0] != TypePrefix[T]() {
		return ZeroID[T](), errors.New("prefix mismatch")
	}

	decoded, err := decode(parts[1])

	if err != nil {
		return ZeroID[T](), errors.New("parse error")
	}

	return NewID[T](decoded), nil
}

func ParseUUID[T Type](data uuid.UUID) ID[T] {
	parsed, err := data.MarshalBinary()

	if err != nil {
		panic("could not parse uuid into id")
	}

	return NewID[T](*(*[16]byte)(parsed))
}

func RandomID[T Type]() ID[T] {
	data := make([]byte, 16)
	_, err := rand.Read(data)

	if err != nil {
		panic("random id generation failed")
	}

	return NewID[T](*(*[16]byte)(data))
}

func (id ID[T]) Prefix() string {
	return TypePrefix[T]()
}

func (id ID[T]) String() string {
	return id.Prefix() + "_" + encode(id.data)
}

func (id ID[T]) Bytes() [16]byte {
	return id.data
}

func (id ID[T]) UUID() uuid.UUID {
	parsed, err := uuid.FromBytes(id.data[:])

	if err != nil {
		panic("could not convert id to uuid")
	}

	return parsed
}

func CastID[T Type, U Type](id ID[T]) ID[U] {
	return NewID[U](id.data)
}

