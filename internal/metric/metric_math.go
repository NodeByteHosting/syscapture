package metric

import (
	"crypto/rand"
	"encoding/binary"
	"math"
	"math/big"
)

// RoundFloatPtr rounds a float64 to the specified precision and returns a pointer to the result.
func RoundFloatPtr(val float64, precision uint) *float64 {
	ratio := math.Pow(10, float64(precision))
	prc := math.Round(val*ratio) / ratio
	return &prc
}

// RandomIntPtr generates a random integer up to the specified max value and returns a pointer to the result.
func RandomIntPtr(max int64) (*int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return nil, err
	}
	result := int(n.Int64())
	return &result, nil
}

// RandomUInt64Ptr generates a random uint64 and returns a pointer to the result.
func RandomUInt64Ptr() (*uint64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return nil, err
	}
	result := binary.BigEndian.Uint64(b[:])
	return &result, nil
}

// RandomFloatPtr generates a random float64 between 0 and 1 and returns a pointer to the result.
func RandomFloatPtr() (*float64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return nil, err
	}
	randomUint64 := binary.BigEndian.Uint64(b[:])
	result := float64(randomUint64) / float64(math.MaxUint64)
	return &result, nil
}
