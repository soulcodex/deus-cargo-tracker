package utils

import (
	"github.com/google/uuid"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/errutil"
)

type UUID string

func (u UUID) String() string {
	return string(u)
}

type UUIDProvider interface {
	New() UUID
}

type RandomUUIDProvider struct{}

func NewRandomUUIDProvider() *RandomUUIDProvider {
	return &RandomUUIDProvider{}
}

func (up RandomUUIDProvider) New() UUID {
	return UUID(uuid.New().String())
}

type FixedUUIDProvider struct {
	uuid string
}

func NewFixedUUIDProvider() *FixedUUIDProvider {
	return &FixedUUIDProvider{
		uuid: "",
	}
}

func (up *FixedUUIDProvider) New() UUID {
	if up.uuid == "" {
		up.uuid = uuid.New().String()
	}

	return UUID(up.uuid)
}

func GuardUUID(raw string) error {
	_, err := uuid.Parse(raw)
	if err != nil {
		return errutil.NewError("invalid UUID string provided").Wrap(err)
	}

	return nil
}

func NewUUID() UUID {
	return UUID(uuid.New().String())
}
