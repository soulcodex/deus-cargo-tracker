package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/soulcodex/deus-cargo-tracker/cmd/di"
	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	cargotrackingdomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain/tracking"
	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb/postgres"
	testarrangers "github.com/soulcodex/deus-cargo-tracker/test/arrangers"
	cargotest "github.com/soulcodex/deus-cargo-tracker/test/cargo"
	testutils "github.com/soulcodex/deus-cargo-tracker/test/utils"
	vesseltest "github.com/soulcodex/deus-cargo-tracker/test/vessel"
)

type FetchCargoByIDAcceptanceTestSuite struct {
	suite.Suite

	common       *di.CommonServices
	vesselModule *di.VesselModule
	cargoModule  *di.CargoModule

	dbArranger *testarrangers.PostgresSQLArranger

	vesselID vesseldomain.VesselID
	cargoID  cargodomain.CargoID
	cargo    *cargodomain.Cargo
}

func TestFetchCargoByID(t *testing.T) {
	suite.Run(t, new(FetchCargoByIDAcceptanceTestSuite))
}

func (suite *FetchCargoByIDAcceptanceTestSuite) SetupSuite() {
	suite.common = di.MustInitCommonServicesWithEnvFiles(
		suite.T().Context(),
		"../.env",
		".test.env",
	)
	suite.vesselModule = di.NewVesselModule(suite.T().Context(), suite.common)
	suite.cargoModule = di.NewCargoModule(suite.T().Context(), suite.common)
	suite.common.RedisClient.FlushAll(suite.T().Context())
	suite.vesselID = vesseldomain.VesselID(suite.common.ULIDProvider.New().String())

	dbPool, match := suite.common.DBPool.(*postgres.ConnectionPool)
	suite.Require().True(match, "expected *postgres.ConnectionPool, got different type")

	suite.dbArranger = testarrangers.NewPostgresSQLArranger(suite.common.Config.PostgresSchema, dbPool)
}

func (suite *FetchCargoByIDAcceptanceTestSuite) SetupTest() {
	suite.dbArranger.MustArrange(suite.T().Context())

	vesselID := vesseltest.WithVesselID(suite.vesselID.String())
	vessel := vesseltest.NewVesselMother(vesselID).Build(suite.T())
	err := suite.vesselModule.Repository.Save(suite.T().Context(), vessel)
	suite.Require().NoError(err, "failed to save vessel for suite setup")

	cargo := cargotest.NewCargoMother(cargotest.WithVesselID(suite.vesselID.String())).Build(suite.T())
	saveCargoErr := suite.cargoModule.Repository.Save(suite.T().Context(), cargo)
	suite.Require().NoError(saveCargoErr, "failed to save cargo for suite setup")
	suite.cargoID = cargo.ID()
	suite.cargo = cargo
}

func (suite *FetchCargoByIDAcceptanceTestSuite) TestFetchCargoByID_Success() {
	params := []any{
		suite.cargoID.String(),
		suite.vesselID.String(),
		suite.cargo.Primitives().CreatedAt.Format(time.RFC3339),
		suite.cargo.Primitives().UpdatedAt.Format(time.RFC3339),
	}
	body := []byte(fmt.Sprintf(`
		{
			"data": {
				"id": "%s",
				"type": "cargo",
				"attributes": {
					"vessel_id": "%s",
					"weight": 3500,
					"status": "pending",
					"items": [
						{"name": "Electronics", "weight": 1500},
						{"name": "Clothing", "weight": 2000}
					],
					"created_at": "%s",
					"updated_at": "%s"
				}
			}
		}
	`, params...))
	response := testutils.ExecuteJSONRequest(
		suite.T(),
		suite.common.Router,
		http.MethodGet,
		"/cargoes/"+suite.cargo.ID().String(),
		body,
	)
	testutils.CheckResponse(suite.T(), http.StatusOK, string(body), response)
}

func (suite *FetchCargoByIDAcceptanceTestSuite) TestFetchCargoByID_SuccessWithTracking() {
	trackingItem := cargotrackingdomain.NewTrackingOnCargoCreated(
		cargotrackingdomain.TrackingID(suite.common.ULIDProvider.New().String()),
		cargodomain.StatusPending.String(),
		suite.common.TimeProvider.Now(),
	)

	cargoID, vesselID := suite.common.ULIDProvider.New().String(), suite.vesselID.String()
	tracking := cargotest.WithTracking(cargotrackingdomain.NewTrackingItemPrimitives(cargoID, trackingItem))

	cargo := cargotest.NewCargoMother(cargotest.WithID(cargoID), cargotest.WithVesselID(vesselID), tracking).Build(suite.T())
	saveCargoErr := suite.cargoModule.Repository.Save(suite.T().Context(), cargo)
	suite.Require().NoError(saveCargoErr, "failed to save cargo for suite setup")
	primitives := cargo.Primitives()

	params := []any{
		primitives.ID,
		primitives.VesselID,
		primitives.CreatedAt.Format(time.RFC3339),
		primitives.UpdatedAt.Format(time.RFC3339),
		primitives.Tracking[0].ID,
		primitives.Tracking[0].ID,
		primitives.Tracking[0].CreatedAt.Format(time.RFC3339),
	}

	body := []byte(fmt.Sprintf(`
		{
			"data": {
				"id": "%s",
				"type": "cargo",
				"attributes": {
					"vessel_id": "%s",
					"weight": 3500,
					"status": "pending",
					"items": [
						{"name": "Electronics", "weight": 1500},
						{"name": "Clothing", "weight": 2000}
					],
					"created_at": "%s",
					"updated_at": "%s"
				},
				"relationships" : {
				  "tracking" : {
					"data" : [ {
					  "type" : "tracking",
					  "id" : "%s"
					} ]
				  }
				}
			},
			"included": [
				{
    				"type" : "tracking",
    				"id" : "%s",
					"attributes" : {
					  "created_at" : "%s",
					  "entry_type" : "cargo.created",
					  "status_after" : "pending",
					  "status_before" : "pending"
					}
				}
			]
		}
	`, params...))
	response := testutils.ExecuteJSONRequest(
		suite.T(),
		suite.common.Router,
		http.MethodGet,
		"/cargoes/"+primitives.ID+"?tracking=true",
		body,
	)
	testutils.CheckResponse(suite.T(), http.StatusOK, string(body), response)
}

func (suite *FetchCargoByIDAcceptanceTestSuite) TestFetchCargoByID_FailIfInvalidCargoIDIsProvided() {
	body := []byte(`
		{
		  "errors" : [ {
			"id" : "<<PRESENCE>>",
			"title" : "Bad Request",
			"detail" : "invalid cargo ID provided",
			"status" : "400",
			"code" : "bad_request"
		  } ]
		}
	`)
	response := testutils.ExecuteJSONRequest(
		suite.T(),
		suite.common.Router,
		http.MethodGet,
		"/cargoes/1",
		body,
	)
	testutils.CheckResponse(suite.T(), http.StatusBadRequest, string(body), response)
}

func (suite *FetchCargoByIDAcceptanceTestSuite) TestFetchCargoByID_FailIfCargoNotExists() {
	body := []byte(`
		{
		  "errors" : [ {
			"id" : "<<PRESENCE>>",
			"title" : "Not Found",
			"detail" : "cargo not found",
			"status" : "404",
			"code" : "not_found"
		  } ]
		}
	`)
	response := testutils.ExecuteJSONRequest(
		suite.T(),
		suite.common.Router,
		http.MethodGet,
		"/cargoes/"+suite.common.ULIDProvider.New().String(),
		body,
	)
	testutils.CheckResponse(suite.T(), http.StatusNotFound, string(body), response)
}
