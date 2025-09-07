package vesselentrypoint

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	vesselqueries "github.com/soulcodex/deus-cargo-tracker/internal/vessel/application/queries"
	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/bus"
	querybus "github.com/soulcodex/deus-cargo-tracker/pkg/bus/query"
	httpserver "github.com/soulcodex/deus-cargo-tracker/pkg/http-server"
	jsonapiresponse "github.com/soulcodex/deus-cargo-tracker/pkg/json-api/response"
)

type FetchVesselByIDResponse struct {
	ID        string    `jsonapi:"primary,vessel"`
	Name      string    `jsonapi:"attr,name"`
	Capacity  uint64    `jsonapi:"attr,capacity"`
	Latitude  float64   `jsonapi:"attr,latitude"`
	Longitude float64   `jsonapi:"attr,longitude"`
	CreatedAt time.Time `jsonapi:"attr,created_at"`
	UpdatedAt time.Time `jsonapi:"attr,updated_at"`
}

func newFetchVesselByIDResponse(resp vesselqueries.VesselResponse) *FetchVesselByIDResponse {
	return &FetchVesselByIDResponse{
		ID:        resp.ID,
		Name:      resp.Name,
		Capacity:  resp.Capacity,
		Latitude:  resp.Latitude,
		Longitude: resp.Longitude,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

func HandleGETFetchVesselByIDV1HTTP(
	queryBus querybus.Bus,
	middleware *httpserver.JSONAPIResponseMiddleware,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vesselID := mux.Vars(r)["vessel_id"]
		if vesselID == "" {
			res, statusCode := jsonapiresponse.NewBadRequest("vessel_id is required"), http.StatusBadRequest
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, vesseldomain.ErrInvalidVesselIDProvided)
			return
		}

		query := &vesselqueries.FetchVesselByIDQuery{ID: vesselID}

		result, err := bus.DispatchWithResponse[*vesselqueries.FetchVesselByIDQuery, vesselqueries.VesselResponse](
			queryBus,
		)(r.Context(), query)

		switch {
		case err == nil:
			middleware.WriteResponse(r.Context(), w, newFetchVesselByIDResponse(result), http.StatusOK)
		case vesseldomain.IsVesselNotExistsError(err):
			res, statusCode := jsonapiresponse.NewNotFound("vessel not found"), http.StatusNotFound
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		case errors.Is(err, vesseldomain.ErrInvalidVesselIDProvided):
			res, statusCode := jsonapiresponse.NewBadRequest("invalid vessel id provided"), http.StatusBadRequest
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		default:
			res, statusCode := jsonapiresponse.NewInternalServerError(), http.StatusInternalServerError
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		}
	}
}
