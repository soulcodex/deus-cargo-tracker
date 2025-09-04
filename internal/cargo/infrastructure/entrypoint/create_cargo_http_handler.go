package cargoentrypoint

import (
	"errors"
	"net/http"

	"github.com/google/jsonapi"

	cargocommands "github.com/soulcodex/deus-cargo-tracker/internal/cargo/application/commands"
	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	cargoinfra "github.com/soulcodex/deus-cargo-tracker/internal/cargo/infrastructure"
	"github.com/soulcodex/deus-cargo-tracker/pkg/bus"
	commandbus "github.com/soulcodex/deus-cargo-tracker/pkg/bus/command"
	httpserver "github.com/soulcodex/deus-cargo-tracker/pkg/http-server"
	jsonapiresponse "github.com/soulcodex/deus-cargo-tracker/pkg/json-api/response"
)

type CreateCargoRequest struct {
	ID       string `jsonapi:"primary,cargo"`
	VesselID string `jsonapi:"attr,vessel_id"`
	Items    []struct {
		Name   string `jsonapi:"attr,name"`
		Weight uint64 `jsonapi:"attr,weight"`
	} `jsonapi:"attr,items"`
}

func HandlePOSTCreateCargoV1HTTP(
	commandBus commandbus.Bus,
	middleware *httpserver.JSONAPIResponseMiddleware,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateCargoRequest
		if err := jsonapi.UnmarshalPayload(r.Body, &req); err != nil {
			res, statusCode := jsonapiresponse.NewBadRequest("invalid received request"), http.StatusBadRequest
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
			return
		}

		cmd := &cargocommands.CreateCargoCommand{
			ID:       req.ID,
			VesselID: req.VesselID,
			Items: []struct {
				Name   string `json:"name"`
				Weight uint64 `json:"weight"`
			}(req.Items),
		}

		err := bus.Dispatch(commandBus)(r.Context(), cmd)

		switch {
		case err == nil:
			middleware.WriteResponse(r.Context(), w, nil, http.StatusNoContent)
		case errors.Is(err, cargoinfra.ErrVesselNotFound):
			res, statusCode := jsonapiresponse.NewNotFound("cargo vessel not found"), http.StatusNotFound
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		case cargodomain.IsCargoAlreadyExistsError(err):
			res, statusCode := jsonapiresponse.NewBadRequest("cargo already exists"), http.StatusConflict
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		case errors.Is(err, cargodomain.ErrInvalidVesselIDProvided):
			res, statusCode := jsonapiresponse.NewBadRequest("invalid vessel ID provided"), http.StatusBadRequest
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		case errors.Is(err, cargodomain.ErrInvalidCargoIDProvided):
			res, statusCode := jsonapiresponse.NewBadRequest("invalid cargo ID provided"), http.StatusBadRequest
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		case errors.Is(err, cargodomain.ErrInvalidItemsProvided):
			res, statusCode := jsonapiresponse.NewBadRequest("invalid cargo items provided"), http.StatusBadRequest
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		default:
			res, statusCode := jsonapiresponse.NewInternalServerError(), http.StatusInternalServerError
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		}
	}
}
