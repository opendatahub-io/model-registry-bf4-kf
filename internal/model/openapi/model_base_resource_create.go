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

// checks if the BaseResourceCreate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BaseResourceCreate{}

// BaseResourceCreate struct for BaseResourceCreate
type BaseResourceCreate struct {
	// User provided custom properties which are not defined by its type.
	CustomProperties *map[string]MetadataValue `json:"customProperties,omitempty"`
	// The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance.
	ExternalID *string `json:"externalID,omitempty"`
	// The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set.
	Name *string `json:"name,omitempty"`
}

// NewBaseResourceCreate instantiates a new BaseResourceCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBaseResourceCreate() *BaseResourceCreate {
	this := BaseResourceCreate{}
	return &this
}

// NewBaseResourceCreateWithDefaults instantiates a new BaseResourceCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBaseResourceCreateWithDefaults() *BaseResourceCreate {
	this := BaseResourceCreate{}
	return &this
}

// GetCustomProperties returns the CustomProperties field value if set, zero value otherwise.
func (o *BaseResourceCreate) GetCustomProperties() map[string]MetadataValue {
	if o == nil || IsNil(o.CustomProperties) {
		var ret map[string]MetadataValue
		return ret
	}
	return *o.CustomProperties
}

// GetCustomPropertiesOk returns a tuple with the CustomProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseResourceCreate) GetCustomPropertiesOk() (*map[string]MetadataValue, bool) {
	if o == nil || IsNil(o.CustomProperties) {
		return nil, false
	}
	return o.CustomProperties, true
}

// HasCustomProperties returns a boolean if a field has been set.
func (o *BaseResourceCreate) HasCustomProperties() bool {
	if o != nil && !IsNil(o.CustomProperties) {
		return true
	}

	return false
}

// SetCustomProperties gets a reference to the given map[string]MetadataValue and assigns it to the CustomProperties field.
func (o *BaseResourceCreate) SetCustomProperties(v map[string]MetadataValue) {
	o.CustomProperties = &v
}

// GetExternalID returns the ExternalID field value if set, zero value otherwise.
func (o *BaseResourceCreate) GetExternalID() string {
	if o == nil || IsNil(o.ExternalID) {
		var ret string
		return ret
	}
	return *o.ExternalID
}

// GetExternalIDOk returns a tuple with the ExternalID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseResourceCreate) GetExternalIDOk() (*string, bool) {
	if o == nil || IsNil(o.ExternalID) {
		return nil, false
	}
	return o.ExternalID, true
}

// HasExternalID returns a boolean if a field has been set.
func (o *BaseResourceCreate) HasExternalID() bool {
	if o != nil && !IsNil(o.ExternalID) {
		return true
	}

	return false
}

// SetExternalID gets a reference to the given string and assigns it to the ExternalID field.
func (o *BaseResourceCreate) SetExternalID(v string) {
	o.ExternalID = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *BaseResourceCreate) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BaseResourceCreate) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *BaseResourceCreate) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *BaseResourceCreate) SetName(v string) {
	o.Name = &v
}

func (o BaseResourceCreate) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BaseResourceCreate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CustomProperties) {
		toSerialize["customProperties"] = o.CustomProperties
	}
	if !IsNil(o.ExternalID) {
		toSerialize["externalID"] = o.ExternalID
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	return toSerialize, nil
}

type NullableBaseResourceCreate struct {
	value *BaseResourceCreate
	isSet bool
}

func (v NullableBaseResourceCreate) Get() *BaseResourceCreate {
	return v.value
}

func (v *NullableBaseResourceCreate) Set(val *BaseResourceCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableBaseResourceCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableBaseResourceCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBaseResourceCreate(val *BaseResourceCreate) *NullableBaseResourceCreate {
	return &NullableBaseResourceCreate{value: val, isSet: true}
}

func (v NullableBaseResourceCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBaseResourceCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
