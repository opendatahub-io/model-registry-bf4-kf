package core

import (
	"github.com/opendatahub-io/model-registry/internal/ml_metadata/proto"
)

func NewListResult[T any](results []T, listOptions ListOptions, nextToken *string) ListResult {
	size := int32(len(results))
	listResult := ListResult{
		PageSize:      listOptions.PageSize,
		Size:          &size,
		NextPageToken: nextToken,
	}
	return listResult
}

func BuildListOperationOptions(listOptions ListOptions) (*proto.ListOperationOptions, error) {

	result := proto.ListOperationOptions{}
	if listOptions.PageSize != nil {
		result.MaxResultSize = listOptions.PageSize
	}
	if listOptions.NextPageToken != nil {
		result.NextPageToken = listOptions.NextPageToken
	}
	if listOptions.OrderBy != nil {
		so := listOptions.SortOrder

		// default is DESC
		isAsc := false
		if so != nil && *so == "ASC" {
			isAsc = true
		}

		var orderByField proto.ListOperationOptions_OrderByField_Field
		if val, ok := proto.ListOperationOptions_OrderByField_Field_value[*listOptions.OrderBy]; ok {
			orderByField = proto.ListOperationOptions_OrderByField_Field(val)
		}

		result.OrderByField = &proto.ListOperationOptions_OrderByField{
			Field: &orderByField,
			IsAsc: &isAsc,
		}
	}
	return &result, nil
}
