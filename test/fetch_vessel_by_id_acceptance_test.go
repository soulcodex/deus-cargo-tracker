package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/soulcodex/deus-cargo-tracker/cmd/di"
	vesseldomain "github.com/soulcodex/deus-cargo-tracker/internal/vessel/domain"
	testutils "github.com/soulcodex/deus-cargo-tracker/test/utils"
)

type FetchVesselByIDAcceptanceTestSuite struct {
	suite.Suite

	common       *di.CommonServices
	vesselModule *di.VesselModule

	vesselID vesseldomain.VesselID
}

func TestFindRocketByID(t *testing.T) {
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
}

func (suite *FetchVesselByIDAcceptanceTestSuite) TestFindRocketByID_RocketNotFound() {
	vesselID := suite.common.ULIDProvider.New().String()
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodGet, "/vessels/"+vesselID, nil)
	suite.Equal(http.StatusNotFound, response.Code, "Expected status code 404 Not Found")
}
