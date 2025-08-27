package distributedsync

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/retry"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/logger"
)

const mutexName = "distributed-sync-mutex"

type MutexService interface {
	Mutex(ctx context.Context, key string, fn MutexCallback) (interface{}, error)
}

type MutexCallback func() (interface{}, error)

type RedisMutexService struct {
	options *MutexServiceOptions
	sync    *redsync.Redsync
	logger  logger.ZerologLogger
}

const (
	redsyncDefaultExpiry        = 30 * time.Second
	redsyncDefaultRetryDelay    = 250 * time.Millisecond
	redsyncDefaultTimeoutFactor = 0.05
)

func NewRedisMutexService(redisClient *redis.Client, logger logger.ZerologLogger, options ...MutexServiceOptFunc) *RedisMutexService {
	pool := goredis.NewPool(redisClient)

	return &RedisMutexService{sync: redsync.New(pool), logger: logger, options: NewMutexServiceOptions(options...)}
}

func (rm *RedisMutexService) Mutex(ctx context.Context, key string, fn MutexCallback) (interface{}, error) {
	mutex := rm.sync.NewMutex(
		rm.lockingKey(key),
		redsync.WithExpiry(redsyncDefaultExpiry),
		redsync.WithRetryDelay(redsyncDefaultRetryDelay),
		redsync.WithTimeoutFactor(redsyncDefaultTimeoutFactor),
	)

	if _, lockingErr := retry.Do(func() (interface{}, error) {
		return nil, rm.acquireLock(ctx, mutex)
	}, int(rm.options.Retries)); lockingErr != nil {
		rm.logger.Error().
			Ctx(ctx).
			Err(lockingErr).
			Str("mutex_key", mutex.Name()).
			Msg("error locking mutex sync")
		return nil, NewMutexLockingError(key).Wrap(lockingErr)
	}

	result, err := fn()

	if _, unlockingErr := retry.Do(func() (interface{}, error) {
		return nil, rm.releaseLock(ctx, mutex)
	}, int(rm.options.Retries)); unlockingErr != nil {
		rm.logger.Error().
			Ctx(ctx).
			Err(unlockingErr).
			Str("mutex_key", mutex.Name()).
			Msg("error unlocking mutex sync")
		return nil, NewMutexUnlockingError(key).Wrap(unlockingErr)
	}

	return result, err
}

func (rm *RedisMutexService) lockingKey(key string) string {
	if rm.options.ServicePrefix != nil {
		return mutexName + ":" + *rm.options.ServicePrefix + ":" + key
	}

	return mutexName + ":" + key
}

func (rm *RedisMutexService) releaseLock(ctx context.Context, mutex *redsync.Mutex) error {
	if ok, err := mutex.UnlockContext(ctx); !ok || err != nil {
		rm.logger.Warn().
			Ctx(ctx).
			Err(err).
			Str("db.operation.parameter.mutex_key", mutex.Name()).
			Msg("error unlocking mutex")
		if err != nil {
			return fmt.Errorf("error while unlocking mutex: %w", err)
		}

		return errors.New("redis mutex invalid status when unlocking")
	}

	return nil
}

func (rm *RedisMutexService) acquireLock(ctx context.Context, mutex *redsync.Mutex) error {
	if err := mutex.LockContext(ctx); err != nil {
		rm.logger.Warn().
			Ctx(ctx).
			Err(err).
			Str("mutex_key", mutex.Name()).
			Msg("error locking mutex sync - retrying")
		return fmt.Errorf("error while releasing mutex: %w", err)
	}

	return nil
}
