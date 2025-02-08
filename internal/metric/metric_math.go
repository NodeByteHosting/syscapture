package metric

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"math/big"
)

// RoundFloatPtr rounds a float to a given precision and returns the pointer of the result.
func RoundFloatPtr(val float64, precision uint) *float64 {
	r := RoundFloat(val, precision)
	return &r
}

// RoundFloat rounds a float to a given precision and returns the result.
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	prc := math.Round(val*ratio) / ratio
	return prc
}

// RandomIntPtr generates a random integer up to a maximum value and returns the pointer of the result.
func RandomIntPtr(maximum int64) *int {
	n, err := rand.Int(rand.Reader, big.NewInt(maximum))
	if err != nil {
		panic(err) // handle error appropriately in production
	}
	result := int(n.Int64())
	return &result
}

// RandomUInt64Ptr generates a random uint64 and returns the pointer of the result.
func RandomUInt64Ptr() *uint64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err) // handle error appropriately in production
	}
	result := binary.BigEndian.Uint64(b[:])
	return &result
}

// RandomFloatPtr generates a random float64 and returns the pointer of the result.
func RandomFloatPtr() *float64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err) // handle error appropriately in production
	}
	randomUint64 := binary.BigEndian.Uint64(b[:])
	result := float64(randomUint64) / float64(math.MaxUint64)
	return &result
}
