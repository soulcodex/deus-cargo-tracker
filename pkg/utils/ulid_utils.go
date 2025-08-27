package utils

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
)

type ULID string

func (u ULID) String() string {
	return string(u)
}

type ULIDProvider interface {
	New() ULID
}

type RandomULIDProvider struct {
}

func NewRandomULIDProvider() *RandomULIDProvider {
	return &RandomULIDProvider{}
}

func (up RandomULIDProvider) New() ULID {
	return NewULID()
}

type FixedULIDProvider struct {
	ulid ULID
	lock sync.Mutex
}

func NewFixedULIDProvider() *FixedULIDProvider {
	return &FixedULIDProvider{
		ulid: "",
		lock: sync.Mutex{},
	}
}

func (up *FixedULIDProvider) New() ULID {
	defer up.lock.Unlock()
	up.lock.Lock()

	if up.ulid == "" {
		up.ulid = NewULID()
	}

	return up.ulid
}

func GuardULID(rawUlid string) error {
	_, err := ulid.Parse(rawUlid)
	if err != nil {
		return errutil.NewError("invalid ULID string provided").Wrap(err)
	}

	return nil
}

func NewULID() ULID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	return ULID(ulid.MustNew(ulid.Timestamp(t), entropy).String())
}
