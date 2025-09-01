package jsonapiresponse

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"

	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

const (
	notFoundCode  = "not_found"
	notFoundTitle = "Not Found"
)

func NewNotFound(detail string) []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewULID().String(),
		Code:   notFoundCode,
		Title:  notFoundTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusNotFound),
	}}
}

func NewNotFoundErrorWithDetails(detail string, items ...MetadataItem) []*jsonapi.ErrorObject {
	metadata := NewMetadata(items...).MetadataMap()

	return []*jsonapi.ErrorObject{{
		ID:     utils.NewULID().String(),
		Code:   notFoundCode,
		Title:  notFoundTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusNotFound),
		Meta:   &metadata,
	}}
}
