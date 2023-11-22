/*
Model Registry REST API

REST API for Model Registry to create and manage ML model metadata

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
)

// ModelVersionState - LIVE: A state indicating that the `ModelVersion` exists - ARCHIVED: A state indicating that the `ModelVersion` has been archived.
type ModelVersionState string

// List of ModelVersionState
const (
	MODELVERSIONSTATE_LIVE     ModelVersionState = "LIVE"
	MODELVERSIONSTATE_ARCHIVED ModelVersionState = "ARCHIVED"
)

// All allowed values of ModelVersionState enum
var AllowedModelVersionStateEnumValues = []ModelVersionState{
	"LIVE",
	"ARCHIVED",
}

func (v *ModelVersionState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ModelVersionState(value)
	for _, existing := range AllowedModelVersionStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ModelVersionState", value)
}

// NewModelVersionStateFromValue returns a pointer to a valid ModelVersionState
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewModelVersionStateFromValue(v string) (*ModelVersionState, error) {
	ev := ModelVersionState(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ModelVersionState: valid values are %v", v, AllowedModelVersionStateEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ModelVersionState) IsValid() bool {
	for _, existing := range AllowedModelVersionStateEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ModelVersionState value
func (v ModelVersionState) Ptr() *ModelVersionState {
	return &v
}

type NullableModelVersionState struct {
	value *ModelVersionState
	isSet bool
}

func (v NullableModelVersionState) Get() *ModelVersionState {
	return v.value
}

func (v *NullableModelVersionState) Set(val *ModelVersionState) {
	v.value = val
	v.isSet = true
}

func (v NullableModelVersionState) IsSet() bool {
	return v.isSet
}

func (v *NullableModelVersionState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableModelVersionState(val *ModelVersionState) *NullableModelVersionState {
	return &NullableModelVersionState{value: val, isSet: true}
}

func (v NullableModelVersionState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableModelVersionState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
