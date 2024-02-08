/*
Model Registry REST API

REST API for Model Registry to create and manage ML model metadata

API version: v1alpha1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the BaseResourceList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BaseResourceList{}

// BaseResourceList struct for BaseResourceList
type BaseResourceList struct {
	// Token to use to retrieve next page of results.
	NextPageToken string `json:"nextPageToken"`
	// Maximum number of resources to return in the result.
	PageSize int32 `json:"pageSize"`
	// Number of items in result list.
	Size int32 `json:"size"`
}

// NewBaseResourceList instantiates a new BaseResourceList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBaseResourceList(nextPageToken string, pageSize int32, size int32) *BaseResourceList {
	this := BaseResourceList{}
	this.NextPageToken = nextPageToken
	this.PageSize = pageSize
	this.Size = size
	return &this
}

// NewBaseResourceListWithDefaults instantiates a new BaseResourceList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBaseResourceListWithDefaults() *BaseResourceList {
	this := BaseResourceList{}
	return &this
}

// GetNextPageToken returns the NextPageToken field value
func (o *BaseResourceList) GetNextPageToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value
// and a boolean to check if the value has been set.
func (o *BaseResourceList) GetNextPageTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NextPageToken, true
}

// SetNextPageToken sets field value
func (o *BaseResourceList) SetNextPageToken(v string) {
	o.NextPageToken = v
}

// GetPageSize returns the PageSize field value
func (o *BaseResourceList) GetPageSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.PageSize
}

// GetPageSizeOk returns a tuple with the PageSize field value
// and a boolean to check if the value has been set.
func (o *BaseResourceList) GetPageSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PageSize, true
}

// SetPageSize sets field value
func (o *BaseResourceList) SetPageSize(v int32) {
	o.PageSize = v
}

// GetSize returns the Size field value
func (o *BaseResourceList) GetSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Size
}

// GetSizeOk returns a tuple with the Size field value
// and a boolean to check if the value has been set.
func (o *BaseResourceList) GetSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Size, true
}

// SetSize sets field value
func (o *BaseResourceList) SetSize(v int32) {
	o.Size = v
}

func (o BaseResourceList) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BaseResourceList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["nextPageToken"] = o.NextPageToken
	toSerialize["pageSize"] = o.PageSize
	toSerialize["size"] = o.Size
	return toSerialize, nil
}

type NullableBaseResourceList struct {
	value *BaseResourceList
	isSet bool
}

func (v NullableBaseResourceList) Get() *BaseResourceList {
	return v.value
}

func (v *NullableBaseResourceList) Set(val *BaseResourceList) {
	v.value = val
	v.isSet = true
}

func (v NullableBaseResourceList) IsSet() bool {
	return v.isSet
}

func (v *NullableBaseResourceList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBaseResourceList(val *BaseResourceList) *NullableBaseResourceList {
	return &NullableBaseResourceList{value: val, isSet: true}
}

func (v NullableBaseResourceList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBaseResourceList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
