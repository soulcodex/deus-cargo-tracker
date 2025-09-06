package cargoentrypoint

import (
	"errors"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"

	cargocommands "github.com/soulcodex/deus-cargo-tracker/internal/cargo/application/commands"
	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/bus"
	commandbus "github.com/soulcodex/deus-cargo-tracker/pkg/bus/command"
	distributedsync "github.com/soulcodex/deus-cargo-tracker/pkg/distributed-sync"
	httpserver "github.com/soulcodex/deus-cargo-tracker/pkg/http-server"
	jsonapiresponse "github.com/soulcodex/deus-cargo-tracker/pkg/json-api/response"
)

type UpdateCargoStatusRequest struct {
	NewStatus string `jsonapi:"attr,new_status"`
}

func HandlePATCHUpdateCargoStatusV1HTTP(
	commandBus commandbus.Bus,
	mutex distributedsync.MutexService,
	middleware *httpserver.JSONAPIResponseMiddleware,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cargoID := mux.Vars(r)["cargo_id"]
		if cargoID == "" {
			res, statusCode := jsonapiresponse.NewBadRequest("cargo ID is required"), http.StatusBadRequest
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, cargodomain.ErrInvalidCargoIDProvided)
			return
		}

		var req UpdateCargoStatusRequest
		if err := jsonapi.UnmarshalPayload(r.Body, &req); err != nil {
			res, statusCode := jsonapiresponse.NewBadRequest("invalid received request"), http.StatusBadRequest
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
			return
		}

		cmd := &cargocommands.UpdateCargoStatusCommand{
			ID:        cargoID,
			NewStatus: req.NewStatus,
		}

		err := bus.DispatchBlocking(commandBus, mutex)(r.Context(), cmd)

		switch {
		case err == nil:
			middleware.WriteResponse(r.Context(), w, nil, http.StatusNoContent)
		case cargodomain.IsCargoNotExistsError(err):
			res, statusCode := jsonapiresponse.NewNotFound("cargo not found"), http.StatusNotFound
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		case errors.Is(err, cargodomain.ErrInvalidStatusProvided):
			res, statusCode := jsonapiresponse.NewBadRequest("invalid cargo status provided"), http.StatusBadRequest
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		case errors.Is(err, cargodomain.ErrStatusTransitionNotAllowed):
			res, statusCode := jsonapiresponse.NewBadRequest("status transition not allowed"), http.StatusConflict
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		default:
			res, statusCode := jsonapiresponse.NewInternalServerError(), http.StatusInternalServerError
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		}
	}
}
