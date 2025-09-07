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

type UpdateCargoStatusAcceptanceTestSuite struct {
	suite.Suite

	common       *di.CommonServices
	vesselModule *di.VesselModule
	cargoModule  *di.CargoModule

	dbArranger *testarrangers.PostgresSQLArranger

	vesselID vesseldomain.VesselID
	cargoID  cargodomain.CargoID
}

func TestUpdateCargoStatus(t *testing.T) {
	suite.Run(t, new(UpdateCargoStatusAcceptanceTestSuite))
}

func (suite *UpdateCargoStatusAcceptanceTestSuite) SetupSuite() {
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

func (suite *UpdateCargoStatusAcceptanceTestSuite) SetupTest() {
	suite.dbArranger.MustArrange(suite.T().Context())
	suite.common.RedisClient.FlushAll(suite.T().Context())

	vesselID := vesseltest.WithVesselID(suite.vesselID.String())
	vessel := vesseltest.NewVesselMother(vesselID).Build(suite.T())
	err := suite.vesselModule.Repository.Save(suite.T().Context(), vessel)
	suite.Require().NoError(err, "failed to save vessel for suite setup")

	cargo := cargotest.NewCargoMother(cargotest.WithVesselID(suite.vesselID.String())).Build(suite.T())
	saveCargoErr := suite.cargoModule.Repository.Save(suite.T().Context(), cargo)
	suite.Require().NoError(saveCargoErr, "failed to save cargo for suite setup")
	suite.cargoID = cargo.ID()
}

func (suite *UpdateCargoStatusAcceptanceTestSuite) TestUpdateCargoStatus_Success() {
	body := []byte(fmt.Sprintf(`
		{
			"data": {
				"type": "cargo",
				"attributes": {
					"new_status": "%s"
				}
			}
		}
	`, cargodomain.StatusInTransit))
	route := fmt.Sprintf("/cargoes/%s/update-status", suite.cargoID.String())
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodPatch, route, body)
	suite.Equal(http.StatusNoContent, response.Code, "Expected status code 200 OK")
}

func (suite *UpdateCargoStatusAcceptanceTestSuite) TestUpdateCargoStatus_SuccessWithSameStatus() {
	body := []byte(fmt.Sprintf(`
		{
			"data": {
				"type": "cargo",
				"attributes": {
					"new_status": "%s"
				}
			}
		}
	`, "pending"))
	route := fmt.Sprintf("/cargoes/%s/update-status", suite.cargoID.String())
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodPatch, route, body)
	suite.Equal(http.StatusNoContent, response.Code, "Expected status code 200 OK")
}

func (suite *UpdateCargoStatusAcceptanceTestSuite) TestUpdateCargoStatus_FailIfInvalidCargoStatusProvided() {
	body := []byte(`
		{
			"data": {
				"type": "cargo",
				"attributes": {
					"new_status": "invalid_status"
				}
			}
		}
	`)
	route := fmt.Sprintf("/cargoes/%s/update-status", suite.cargoID.String())
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodPatch, route, body)
	suite.Equal(http.StatusBadRequest, response.Code, "Expected status code 400 Bad Request")
}

func (suite *UpdateCargoStatusAcceptanceTestSuite) TestUpdateCargoStatus_FailIfInvalidStatusTransition() {
	body := []byte(`
		{
			"data": {
				"type": "cargo",
				"attributes": {
					"new_status": "delivered"
				}
			}
		}
	`)
	route := fmt.Sprintf("/cargoes/%s/update-status", suite.cargoID.String())
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodPatch, route, body)
	suite.Equal(http.StatusConflict, response.Code, "Expected status code 409 Conflict")
}

func (suite *UpdateCargoStatusAcceptanceTestSuite) TestUpdateCargoStatus_FailIfCargoNotExists() {
	body := []byte(fmt.Sprintf(`
		{
			"data": {
				"type": "cargo",
				"attributes": {
					"new_status": "%s"
				}
			}
		}
	`, cargodomain.StatusInTransit))
	route := fmt.Sprintf("/cargoes/%s/update-status", suite.common.ULIDProvider.New().String())
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodPatch, route, body)
	suite.Equal(http.StatusNotFound, response.Code, "Expected status code 404 Not Found")
}
