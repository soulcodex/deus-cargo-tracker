package jsonapiresponse

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"

	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

const (
	unauthorizedDefaultTitle = "Unauthorized"
	unauthorizedDefaultCode  = "unauthorized"
)

func NewUnauthorized(detail string) []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewULID().String(),
		Code:   unauthorizedDefaultCode,
		Title:  unauthorizedDefaultTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusUnauthorized),
	}}
}

func NewUnauthorizedWithDetails(detail string, items ...MetadataItem) []*jsonapi.ErrorObject {
	metadata := NewMetadata(items...).MetadataMap()

	return []*jsonapi.ErrorObject{{
		ID:     utils.NewULID().String(),
		Code:   unauthorizedDefaultCode,
		Title:  unauthorizedDefaultTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusUnauthorized),
		Meta:   &metadata,
	}}
}
