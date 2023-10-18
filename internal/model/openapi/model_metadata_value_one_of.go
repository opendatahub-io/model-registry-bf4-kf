/*
Model Registry REST API

REST API for Model Registry to create and manage ML model metadata

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the MetadataValueOneOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &MetadataValueOneOf{}

// MetadataValueOneOf struct for MetadataValueOneOf
type MetadataValueOneOf struct {
	IntValue *string `json:"int_value,omitempty"`
}

// NewMetadataValueOneOf instantiates a new MetadataValueOneOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMetadataValueOneOf() *MetadataValueOneOf {
	this := MetadataValueOneOf{}
	return &this
}

// NewMetadataValueOneOfWithDefaults instantiates a new MetadataValueOneOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMetadataValueOneOfWithDefaults() *MetadataValueOneOf {
	this := MetadataValueOneOf{}
	return &this
}

// GetIntValue returns the IntValue field value if set, zero value otherwise.
func (o *MetadataValueOneOf) GetIntValue() string {
	if o == nil || IsNil(o.IntValue) {
		var ret string
		return ret
	}
	return *o.IntValue
}

// GetIntValueOk returns a tuple with the IntValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetadataValueOneOf) GetIntValueOk() (*string, bool) {
	if o == nil || IsNil(o.IntValue) {
		return nil, false
	}
	return o.IntValue, true
}

// HasIntValue returns a boolean if a field has been set.
func (o *MetadataValueOneOf) HasIntValue() bool {
	if o != nil && !IsNil(o.IntValue) {
		return true
	}

	return false
}

// SetIntValue gets a reference to the given string and assigns it to the IntValue field.
func (o *MetadataValueOneOf) SetIntValue(v string) {
	o.IntValue = &v
}

func (o MetadataValueOneOf) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o MetadataValueOneOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.IntValue) {
		toSerialize["int_value"] = o.IntValue
	}
	return toSerialize, nil
}

type NullableMetadataValueOneOf struct {
	value *MetadataValueOneOf
	isSet bool
}

func (v NullableMetadataValueOneOf) Get() *MetadataValueOneOf {
	return v.value
}

func (v *NullableMetadataValueOneOf) Set(val *MetadataValueOneOf) {
	v.value = val
	v.isSet = true
}

func (v NullableMetadataValueOneOf) IsSet() bool {
	return v.isSet
}

func (v *NullableMetadataValueOneOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMetadataValueOneOf(val *MetadataValueOneOf) *NullableMetadataValueOneOf {
	return &NullableMetadataValueOneOf{value: val, isSet: true}
}

func (v NullableMetadataValueOneOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMetadataValueOneOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}