package test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/jsonapi"
	"github.com/stretchr/testify/suite"

	"github.com/soulcodex/deus-cargo-tracker/cmd/di"
	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
	vesselentrypoint "github.com/soulcodex/deus-cargo-tracker/internal/vessel/infrastructure/entrypoint"
	"github.com/soulcodex/deus-cargo-tracker/pkg/sqldb/postgres"
	testarrangers "github.com/soulcodex/deus-cargo-tracker/test/arrangers"
	testutils "github.com/soulcodex/deus-cargo-tracker/test/utils"
	vesseltest "github.com/soulcodex/deus-cargo-tracker/test/vessel"
)

type FetchVesselByIDAcceptanceTestSuite struct {
	suite.Suite

	common       *di.CommonServices
	vesselModule *di.VesselModule

	dbArranger *testarrangers.PostgresSQLArranger

	vesselID vesseldomain.VesselID
}

func TestFetchVesselByID(t *testing.T) {
	suite.Run(t, new(FetchVesselByIDAcceptanceTestSuite))
}

func (suite *FetchVesselByIDAcceptanceTestSuite) SetupSuite() {
	suite.common = di.MustInitCommonServicesWithEnvFiles(
		suite.T().Context(),
		"../.env",
		".test.env",
	)
	suite.vesselModule = di.NewVesselModule(suite.T().Context(), suite.common)
	suite.common.RedisClient.FlushAll(suite.T().Context())
	suite.vesselID = vesseldomain.VesselID(suite.common.ULIDProvider.New().String())

	vesselID := vesseltest.WithVesselID(suite.vesselID.String())
	vessel := vesseltest.NewVesselMother(vesselID).Build(suite.T())
	err := suite.vesselModule.Repository.Save(suite.T().Context(), vessel)
	suite.Require().NoError(err, "failed to save vessel for suite setup")

	dbPool, match := suite.common.DBPool.(*postgres.ConnectionPool)
	suite.Require().True(match, "expected *postgres.ConnectionPool, got different type")

	suite.dbArranger = testarrangers.NewPostgresSQLArranger(suite.common.Config.PostgresSchema, dbPool)
}

func (suite *FetchVesselByIDAcceptanceTestSuite) SetupTest() {
	suite.dbArranger.MustArrange(suite.T().Context())
}

func (suite *FetchVesselByIDAcceptanceTestSuite) TestFetchVesselByID_VesselNotFound() {
	vesselID := suite.common.ULIDProvider.New().String()
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodGet, "/vessels/"+vesselID, nil)
	suite.Equal(http.StatusNotFound, response.Code, "Expected status code 404 Not Found")
}

func (suite *FetchVesselByIDAcceptanceTestSuite) TestFetchVesselByID_InvalidVesselID() {
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodGet, "/vessels/1", nil)
	suite.Equal(http.StatusBadRequest, response.Code, "Expected status code 400 Bad Request")
}

func (suite *FetchVesselByIDAcceptanceTestSuite) TestFetchVesselByID_Success() {
	vesselID := suite.common.ULIDProvider.New().String()
	vessel := vesseltest.NewVesselMother(vesseltest.WithVesselID(vesselID))
	err := suite.vesselModule.Repository.Save(suite.T().Context(), vessel.Build(suite.T()))
	suite.Require().NoError(err, "failed to soft delete vessel for suite setup")

	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodGet, "/vessels/"+vesselID, nil)
	suite.Equal(http.StatusOK, response.Code, "Expected status code 200 OK")

	responseBody := suite.parseResponseBody(response.Body)
	suite.Equal(vesselID, responseBody.ID, "vessel name mismatch")
	suite.assertVessel(responseBody)
}

func (suite *FetchVesselByIDAcceptanceTestSuite) TestFetchVesselByID_SoftDeletedSuccess() {
	vesselID := suite.vesselID.String()
	vessel := vesseltest.NewVesselMother(vesseltest.WithVesselID(vesselID), vesseltest.WithSoftDeletion(time.Now()))
	err := suite.vesselModule.Repository.Save(suite.T().Context(), vessel.Build(suite.T()))
	suite.Require().NoError(err, "failed to soft delete vessel for suite setup")

	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodGet, "/vessels/"+vesselID, nil)
	suite.Equal(http.StatusNotFound, response.Code, "Expected status code 404 Not Found")
}

func (suite *FetchVesselByIDAcceptanceTestSuite) assertVessel(response *vesselentrypoint.FetchVesselByIDResponse) {
	suite.T().Helper()

	suite.Equal("Falcon 9", response.Name, "vessel name mismatch")
	suite.Equal(uint64(5000), response.Capacity, "vessel capacity mismatch")
	suite.Equal(37.7749, response.Latitude, "vessel latitude mismatch")
	suite.Equal(-122.4194, response.Longitude, "vessel longitude mismatch")
}

func (suite *FetchVesselByIDAcceptanceTestSuite) parseResponseBody(body io.Reader) *vesselentrypoint.FetchVesselByIDResponse {
	suite.T().Helper()

	var response vesselentrypoint.FetchVesselByIDResponse
	err := jsonapi.UnmarshalPayload(body, &response)
	suite.NoError(err, "failed to unmarshal vessel response")

	return &response
}
