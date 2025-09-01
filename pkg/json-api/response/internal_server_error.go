package jsonapiresponse

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"

	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

const (
	internalServerErrorTitle          = "Internal Server Error"
	internalServerErrorCode           = "internal_server_error"
	internalServerDefaultErrorMessage = "Internal Server Error"
)

func NewInternalServerError() []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewULID().String(),
		Code:   internalServerErrorCode,
		Title:  internalServerErrorTitle,
		Detail: internalServerDefaultErrorMessage,
		Status: strconv.Itoa(http.StatusInternalServerError),
	}}
}

func NewInternalServerErrorWithDetails(detail string, items ...MetadataItem) []*jsonapi.ErrorObject {
	metadata := NewMetadata(items...).MetadataMap()

	return []*jsonapi.ErrorObject{{
		ID:     utils.NewULID().String(),
		Code:   internalServerErrorCode,
		Title:  internalServerErrorTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusInternalServerError),
		Meta:   &metadata,
	}}
}
