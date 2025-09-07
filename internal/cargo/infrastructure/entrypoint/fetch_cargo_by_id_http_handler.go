package cargoentrypoint

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	cargoqueries "github.com/soulcodex/deus-cargo-tracker/internal/cargo/application/queries"
	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/bus"
	querybus "github.com/soulcodex/deus-cargo-tracker/pkg/bus/query"
	httpserver "github.com/soulcodex/deus-cargo-tracker/pkg/http-server"
	jsonapiresponse "github.com/soulcodex/deus-cargo-tracker/pkg/json-api/response"
)

type CargoTracking struct {
	ID           string    `jsonapi:"primary,tracking"`
	EntryType    string    `jsonapi:"attr,entry_type"`
	StatusBefore *string   `jsonapi:"attr,status_before"`
	StatusAfter  *string   `jsonapi:"attr,status_after"`
	CreatedAt    time.Time `jsonapi:"attr,created_at,rfc3339"`
}

func newCargoTracking(resp cargoqueries.CargoResponse) []*CargoTracking {
	trackingItems := make([]*CargoTracking, 0, len(resp.Tracking))
	for _, t := range resp.Tracking {
		trackingItems = append(trackingItems, &CargoTracking{
			ID:           t.ID,
			EntryType:    t.EntryType,
			StatusBefore: t.StatusBefore,
			StatusAfter:  t.StatusAfter,
			CreatedAt:    t.CreatedAt,
		})
	}

	return trackingItems
}

type FetchCargoByIDResponse struct {
	ID       string `jsonapi:"primary,cargo"`
	VesselID string `jsonapi:"attr,vessel_id"`
	Items    []struct {
		Name   string `json:"name"`
		Weight uint64 `json:"weight"`
	} `jsonapi:"attr,items"`
	Tracking  []*CargoTracking `jsonapi:"relation,tracking,omitempty"`
	Status    string           `jsonapi:"attr,status"`
	Weight    uint64           `jsonapi:"attr,weight"`
	CreatedAt time.Time        `jsonapi:"attr,created_at,rfc3339"`
	UpdatedAt time.Time        `jsonapi:"attr,updated_at,rfc3339"`
}

func newFetchCargoByIDResponse(resp cargoqueries.CargoResponse) *FetchCargoByIDResponse {
	items := make([]struct {
		Name   string `json:"name"`
		Weight uint64 `json:"weight"`
	}, len(resp.Items))
	for i, item := range resp.Items {
		items[i] = struct {
			Name   string `json:"name"`
			Weight uint64 `json:"weight"`
		}{
			Name:   item.Name,
			Weight: item.Weight,
		}
	}

	return &FetchCargoByIDResponse{
		ID:        resp.ID,
		VesselID:  resp.VesselID,
		Items:     items,
		Tracking:  newCargoTracking(resp),
		Status:    resp.Status,
		Weight:    resp.Weight,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

const (
	trackingQueryParam = "tracking"
)

func HandleGETFetchCargoByIDV1HTTP(
	queryBus querybus.Bus,
	middleware *httpserver.JSONAPIResponseMiddleware,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cargoID := mux.Vars(r)["cargo_id"]
		if cargoID == "" {
			res, statusCode := jsonapiresponse.NewBadRequest("cargo ID is required"), http.StatusBadRequest
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, cargodomain.ErrInvalidCargoIDProvided)
			return
		}

		// This must be by made using an include query param to follow JSON API standard
		// but for simplicity I'm going to use a custom query param.
		tracking := httpserver.FetchBoolQueryParamValue(r.URL.Query(), trackingQueryParam, false)

		query := &cargoqueries.FetchCargoByID{ID: cargoID, Tracking: tracking}

		result, err := bus.DispatchWithResponse[*cargoqueries.FetchCargoByID, cargoqueries.CargoResponse](queryBus)(
			r.Context(),
			query,
		)

		switch {
		case err == nil:
			middleware.WriteResponse(r.Context(), w, newFetchCargoByIDResponse(result), http.StatusOK)
		case cargodomain.IsCargoNotExistsError(err):
			res, statusCode := jsonapiresponse.NewNotFound("cargo not found"), http.StatusNotFound
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		case errors.Is(err, cargodomain.ErrInvalidCargoIDProvided):
			res, statusCode := jsonapiresponse.NewBadRequest("invalid cargo ID provided"), http.StatusBadRequest
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		default:
			res, statusCode := jsonapiresponse.NewInternalServerError(), http.StatusInternalServerError
			middleware.WriteErrorResponse(r.Context(), w, res, statusCode, err)
		}
	}
}
