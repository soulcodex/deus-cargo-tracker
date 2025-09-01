package jsonapiresponse

import (
	"net/http"
	"strconv"

	"github.com/google/jsonapi"

	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

const (
	serviceUnavailableCode  = "service_unavailable"
	serviceUnavailableTitle = "Service Unavailable"
)

func NewUnavailable(detail string) []*jsonapi.ErrorObject {
	return []*jsonapi.ErrorObject{{
		ID:     utils.NewULID().String(),
		Code:   serviceUnavailableCode,
		Title:  serviceUnavailableTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusServiceUnavailable),
	}}
}

func NewUnavailableWithDetails(detail string, items ...MetadataItem) []*jsonapi.ErrorObject {
	metadata := NewMetadata(items...).MetadataMap()

	return []*jsonapi.ErrorObject{{
		ID:     utils.NewULID().String(),
		Code:   serviceUnavailableCode,
		Title:  serviceUnavailableTitle,
		Detail: detail,
		Status: strconv.Itoa(http.StatusServiceUnavailable),
		Meta:   &metadata,
	}}
}
