package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapTo_SimpleStruct(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	params := map[string]any{
		"name": "Alice",
		"age":  30,
	}
	result, err := MapTo[TestStruct](params)
	require.NoError(t, err)
	assert.Equal(t, "Alice", result.Name)
	assert.Equal(t, 30, result.Age)
}

func TestMapTo_MissingTag(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		City string `json:"city"`
	}
	params := map[string]any{
		"name": "Bob",
		"age":  25,
		"city": "Paris",
	}
	result, err := MapTo[TestStruct](params)
	require.NoError(t, err)
	assert.Equal(t, "Bob", result.Name)
	assert.Equal(t, 25, result.Age)
	assert.Equal(t, "Paris", result.City)
}

func TestMapTo_TypeMismatch(t *testing.T) {
	type TestStruct struct {
		Age int `json:"age"`
	}
	params := map[string]any{
		"age": "not-an-int",
	}
	_, err := MapTo[TestStruct](params)
	assert.Error(t, err)
}

func TestMapTo_ConvertibleType(t *testing.T) {
	type TestStruct struct {
		Score float64 `json:"score"`
	}
	params := map[string]any{
		"score": float32(42.5),
	}
	result, err := MapTo[TestStruct](params)
	require.NoError(t, err)
	assert.Equal(t, 42.5, result.Score)
}

func TestMapTo_EmptyParams(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
	}
	params := map[string]any{}
	result, err := MapTo[TestStruct](params)
	require.NoError(t, err)
	assert.Equal(t, "", result.Name)
}

func TestMapTo_ValidationTags(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"min=18"`
	}
	params := map[string]any{
		"name": "Charlie",
		"age":  20,
	}
	result, err := MapTo[TestStruct](params)
	require.NoError(t, err)
	assert.Equal(t, "Charlie", result.Name)
	assert.Equal(t, 20, result.Age)
}

func TestMapTo_ValidationFails(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"min=18"`
	}
	params := map[string]any{
		"name": "",
		"age":  15,
	}
	_, err := MapTo[TestStruct](params)
	require.Error(t, err)
	expectedErr := "Name: required, Age: min=18"
	assert.Equal(t, expectedErr, err.Error())
}

func TestMapTo_SliceElementConversion(t *testing.T) {
	type TestStruct struct {
		IDs []int64 `json:"ids"`
	}
	params := map[string]any{
		"ids": []any{1.0, 2.0, 3.0},
	}
	result, err := MapTo[TestStruct](params)
	require.NoError(t, err)
	assert.Equal(t, []int64{1, 2, 3}, result.IDs)
}

func TestMapTo_RecursiveConversion(t *testing.T) {
	type InnerStruct struct {
		Value string `json:"value"`
	}
	type TestStruct struct {
		Inner InnerStruct `json:"inner"`
	}
	params := map[string]any{
		"inner": map[string]any{
			"value": "test",
		},
	}
	result, err := MapTo[TestStruct](params)
	require.NoError(t, err)
	assert.Equal(t, "test", result.Inner.Value)
}
