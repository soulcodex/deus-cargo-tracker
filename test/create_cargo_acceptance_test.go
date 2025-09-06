package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/soulcodex/deus-cargo-tracker/cmd/di"
	cargodomain "github.com/soulcodex/deus-cargo-tracker/internal/cargo/domain"
	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb/postgres"
	testarrangers "github.com/soulcodex/deus-cargo-tracker/test/arrangers"
	cargotest "github.com/soulcodex/deus-cargo-tracker/test/cargo"
	testutils "github.com/soulcodex/deus-cargo-tracker/test/utils"
	vesseltest "github.com/soulcodex/deus-cargo-tracker/test/vessel"
)

type CreateCargoAcceptanceTestSuite struct {
	suite.Suite

	common       *di.CommonServices
	vesselModule *di.VesselModule
	cargoModule  *di.CargoModule

	dbArranger *testarrangers.PostgresSQLArranger

	vesselID vesseldomain.VesselID
	cargoID  cargodomain.CargoID
}

func TestCreateCargo(t *testing.T) {
	suite.Run(t, new(CreateCargoAcceptanceTestSuite))
}

func (suite *CreateCargoAcceptanceTestSuite) SetupSuite() {
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

func (suite *CreateCargoAcceptanceTestSuite) SetupTest() {
	suite.dbArranger.MustArrange(suite.T().Context())

	vesselID := vesseltest.WithVesselID(suite.vesselID.String())
	vessel := vesseltest.NewVesselMother(vesselID).Build(suite.T())
	err := suite.vesselModule.Repository.Save(suite.T().Context(), vessel)
	suite.Require().NoError(err, "failed to save vessel for suite setup")

	cargo := cargotest.NewCargoMother(cargotest.WithVesselID(suite.vesselID.String())).Build(suite.T())
	saveCargoErr := suite.cargoModule.Repository.Save(suite.T().Context(), cargo)
	suite.Require().NoError(saveCargoErr, "failed to save cargo for suite setup")
	suite.cargoID = cargo.ID()
}

func (suite *CreateCargoAcceptanceTestSuite) TestCreateCargo_Success() {
	body := []byte(fmt.Sprintf(`
		{
			"data": {
				"id": "01K4BBCBY7MQCC5CVGKMRHBBTM",
				"type": "cargo",
				"attributes": {
					"vessel_id": "%s",
					"items": [
						{"name": "Item 1", "weight": 1000},
						{"name": "Item 2", "weight": 2000},
						{"name": "Item 3", "weight": 3000}
					]
				}
			}
		}
	`, suite.vesselID.String()))
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodPost, "/cargoes", body)
	suite.Equal(http.StatusNoContent, response.Code, "Expected status code 204 No Content")
}

func (suite *CreateCargoAcceptanceTestSuite) TestCreateCargo_FailIfAlreadyExists() {
	body := []byte(fmt.Sprintf(`
		{
			"data": {
				"id": "%s",
				"type": "cargo",
				"attributes": {
					"vessel_id": "%s",
					"items": [
						{"name": "Item 1", "weight": 1000},
						{"name": "Item 2", "weight": 2000},
						{"name": "Item 3", "weight": 3000}
					]
				}
			}
		}
	`, suite.cargoID.String(), suite.vesselID.String()))
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodPost, "/cargoes", body)
	suite.Equal(http.StatusConflict, response.Code, "Expected status code 409 Conflict")
}

func (suite *CreateCargoAcceptanceTestSuite) TestCreateCargo_FailIfVesselNotFound() {
	body := []byte(fmt.Sprintf(`
		{
			"data": {
				"id": "%s",
				"type": "cargo",
				"attributes": {
					"vessel_id": "%s",
					"items": [
						{"name": "Item 1", "weight": 1000},
						{"name": "Item 2", "weight": 2000},
						{"name": "Item 3", "weight": 3000}
					]
				}
			}
		}
	`, suite.cargoID.String(), suite.common.ULIDProvider.New().String()))
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodPost, "/cargoes", body)
	suite.Equal(http.StatusNotFound, response.Code, "Expected status code 404 Not Found")
}

func (suite *CreateCargoAcceptanceTestSuite) TestCreateCargo_FailIfCargoItemsExceedLimit() {
	body := []byte(fmt.Sprintf(`
		{
			"data": {
				"id": "%s",
				"type": "cargo",
				"attributes": {
					"vessel_id": "%s",
					"items": [
						{"name": "Item 1", "weight": 1000},
						{"name": "Item 2", "weight": 2000},
						{"name": "Item 3", "weight": 3000},
						{"name": "Item 4", "weight": 3000},
						{"name": "Item 5", "weight": 3000},
						{"name": "Item 6", "weight": 3000},
						{"name": "Item 7", "weight": 3000},
						{"name": "Item 8", "weight": 3000},
						{"name": "Item 9", "weight": 3000},
						{"name": "Item 10", "weight": 3000},
						{"name": "Item 11", "weight": 3000}
					]
				}
			}
		}
	`, suite.common.ULIDProvider.New().String(), suite.vesselID.String()))
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodPost, "/cargoes", body)
	suite.Equal(http.StatusBadRequest, response.Code, "Expected status code 400 Bad Request")
}
