package metric

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
)

var (
	ErrRandomGeneration = errors.New("failed to generate random number")
)

// RoundFloatPtr rounds a float to a given precision and returns the pointer of the result.
func RoundFloatPtr(val float64, precision uint) *float64 {
	r := RoundFloat(val, precision)
	return &r
}

// RoundFloat rounds a float to a given precision and returns the result.
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// RandomIntPtr generates a random integer up to a maximum value and returns the pointer of the result.
func RandomIntPtr(maximum int64) (*int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(maximum))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRandomGeneration, err)
	}
	result := int(n.Int64())
	return &result, nil
}

// RandomUInt64Ptr generates a random uint64 and returns the pointer of the result.
func RandomUInt64Ptr() (*uint64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRandomGeneration, err)
	}
	result := binary.BigEndian.Uint64(b[:])
	return &result, nil
}

// RandomFloatPtr generates a random float64 between 0 and 1 and returns the pointer of the result.
func RandomFloatPtr() (*float64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRandomGeneration, err)
	}
	randomUint64 := binary.BigEndian.Uint64(b[:])
	result := float64(randomUint64) / float64(math.MaxUint64)
	return &result, nil
}

// RandomFloatRangePtr generates a random float64 between min and max and returns the pointer of the result.
func RandomFloatRangePtr(min, max float64) (*float64, error) {
	if min >= max {
		return nil, errors.New("min must be less than max")
	}

	randFloat, err := RandomFloatPtr()
	if err != nil {
		return nil, err
	}

	result := min + (*randFloat * (max - min))
	return &result, nil
}

// BytesToGB converts bytes to gigabytes and returns the result.
func BytesToGB(bytes uint64) float64 {
	return RoundFloat(float64(bytes)/float64(1<<30), 2)
}

// BytesToMB converts bytes to megabytes and returns the result.
func BytesToMB(bytes uint64) float64 {
	return RoundFloat(float64(bytes)/float64(1<<20), 2)
}

// PercentagePtr calculates the percentage of part in total and returns the pointer of the result.
func PercentagePtr(part, total uint64) *float64 {
	if total == 0 {
		return RoundFloatPtr(0, 4)
	}
	result := (float64(part) / float64(total)) * 100
	return RoundFloatPtr(result, 4)
}
