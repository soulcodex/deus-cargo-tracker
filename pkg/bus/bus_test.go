package bus_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/soulcodex/deus-cargo-tracker/pkg/bus"
	dsync "github.com/soulcodex/deus-cargo-tracker/pkg/distributed-sync"
	distributedsyncmock "github.com/soulcodex/deus-cargo-tracker/pkg/distributed-sync/mock"
)

type FakeResponse struct {
	FakeID string
}

func newFakeResponse() *FakeResponse {
	return &FakeResponse{FakeID: "fake_dto_id"}
}

type FakeDto struct {
	FakeID string
}

func (fkd *FakeDto) Type() string {
	return "fake_dto"
}

func (fkd *FakeDto) BlockingKey() string {
	return fkd.FakeID
}

func newFakeDto() *FakeDto {
	return &FakeDto{FakeID: "fake_dto_id"}
}

type FakeUnregisteredDto struct{}

func (fu *FakeUnregisteredDto) Type() string {
	return "fake_unregistered_dto"
}

func (fu *FakeUnregisteredDto) BlockingKey() string {
	return "fake_unregistered_dto_id"
}

func newFakeUnregisteredDto() *FakeUnregisteredDto {
	return &FakeUnregisteredDto{}
}

func createBus() *bus.SyncBus {
	return bus.InitSyncBus()
}

var _ bus.Handler[any, *FakeDto] = (*FakeHandler)(nil)

type FakeHandler struct {
	response interface{}
}

func (fh *FakeHandler) Handle(_ context.Context, _ *FakeDto) (interface{}, error) {
	return fh.response, nil
}

func (fh *FakeHandler) SetResponse(response interface{}) {
	fh.response = response
}

func dispatchScenarios() []struct {
	name             string
	input            func() bus.Dto
	bus              func(t *testing.T) *bus.SyncBus
	expectedResponse func() interface{}
	err              func(i bus.Dto, err error) error
} {
	return []struct {
		name             string
		input            func() bus.Dto
		bus              func(t *testing.T) *bus.SyncBus
		expectedResponse func() interface{}
		err              func(i bus.Dto, err error) error
	}{
		{
			name: "should return an error on unexpected output",
			input: func() bus.Dto {
				return newFakeDto()
			},
			bus: func(t *testing.T) *bus.SyncBus {
				syncBus := createBus()
				handler := &FakeHandler{}
				handler.SetResponse(nil)

				err := syncBus.Register(&FakeDto{}, bus.WrapAsAnyHandler(handler))
				assert.NoError(t, err)

				return syncBus
			},
			expectedResponse: func() interface{} {
				return nil
			},
			err: func(_ bus.Dto, _ error) error {
				return &bus.InvalidOutputReceivedError{}
			},
		},
		{
			name: "should return response",
			input: func() bus.Dto {
				return newFakeDto()
			},
			bus: func(t *testing.T) *bus.SyncBus {
				syncBus := createBus()
				handler := &FakeHandler{}
				handler.SetResponse(newFakeResponse())

				err := syncBus.Register(&FakeDto{}, bus.WrapAsAnyHandler(handler))
				assert.NoError(t, err)

				return syncBus
			},
			expectedResponse: func() interface{} {
				return newFakeResponse()
			},
			err: func(_ bus.Dto, _ error) error {
				return nil
			},
		},
		{
			name: "should return an error on unregistered dto",
			input: func() bus.Dto {
				return newFakeUnregisteredDto()
			},
			bus: func(t *testing.T) *bus.SyncBus {
				syncBus := createBus()
				handler := &FakeHandler{}
				handler.SetResponse(newFakeResponse())

				err := syncBus.Register(&FakeDto{}, bus.WrapAsAnyHandler(handler))
				assert.NoError(t, err)

				return syncBus
			},
			expectedResponse: func() interface{} {
				return nil
			},
			err: func(i bus.Dto, err error) error {
				return bus.ErrNoHandlerForInput(i, err)
			},
		},
	}
}

func Test_SyncBus_DispatchWithResponse(t *testing.T) {
	for _, scenario := range dispatchScenarios() {
		t.Run(scenario.name, func(t *testing.T) {
			syncBus, dto := scenario.bus(t), scenario.input()

			response, err := bus.DispatchWithResponse[bus.Dto, *FakeResponse](syncBus)(
				context.Background(),
				dto,
			)

			if scenarioErr := scenario.err(dto, err); scenarioErr != nil {
				require.Error(t, err)
				require.IsType(t, scenarioErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, scenario.expectedResponse(), response)
			}
		})
	}
}

func Test_SyncBus_HandlerAlreadyRegistered(t *testing.T) {
	t.Run("should return error on already registered handler", func(t *testing.T) {
		syncBus := createBus()
		handler := &FakeHandler{}

		err := syncBus.Register(&FakeDto{}, bus.WrapAsAnyHandler(handler))
		require.NoError(t, err)

		err = syncBus.Register(&FakeDto{}, bus.WrapAsAnyHandler(handler))
		require.IsType(t, &bus.HandlerAlreadyRegisteredError{}, err)
	})
}

func Test_SyncBus_DispatchBlocking(t *testing.T) {
	for _, scenario := range dispatchScenarios()[1:] {
		t.Run(scenario.name, func(t *testing.T) {
			syncBus, dto := scenario.bus(t), scenario.input()
			mutex := &distributedsyncmock.MutexServiceMock{}
			mutex.MutexFunc = func(_ context.Context, _ string, fn dsync.MutexCallback) (interface{}, error) {
				return fn()
			}

			err := bus.DispatchBlocking(syncBus, mutex)(
				context.Background(),
				dto.(bus.BlockingDto),
			)

			if scenarioErr := scenario.err(dto, err); scenarioErr != nil {
				require.Error(t, err)
				require.IsType(t, scenarioErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_SyncBus_Dispatch(t *testing.T) {
	for _, scenario := range dispatchScenarios()[1:] {
		t.Run(scenario.name, func(t *testing.T) {
			syncBus, dto := scenario.bus(t), scenario.input()

			err := bus.Dispatch(syncBus)(
				context.Background(),
				dto,
			)

			if scenarioErr := scenario.err(dto, err); scenarioErr != nil {
				require.Error(t, err)
				require.IsType(t, scenarioErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
