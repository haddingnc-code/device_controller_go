package service

import (
	"context"
	"devices-api-go/internal/domain"
	"errors"
)

type JoinPoint func(ctx context.Context) (interface{}, error)

type DeviceServiceAspect struct {
	next domain.DeviceService
	repo domain.DeviceRepository
}

func NewDeviceServiceAspect(next domain.DeviceService, repo domain.DeviceRepository) *DeviceServiceAspect {
	return &DeviceServiceAspect{
		next: next,
		repo: repo,
	}
}

func (a *DeviceServiceAspect) BeforeDeviceStateCheck(ctx context.Context, id int64, isDelete bool, proceed JoinPoint) (interface{}, error) {
	current, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Rule: In use devices cannot be deleted
	if isDelete && current.State == domain.InUse {
		return nil, errors.New(domain.ErrDeviceInUseDelete)
	}

	return proceed(ctx)
}

func (a *DeviceServiceAspect) FullUpdate(ctx context.Context, id int64, dto domain.DeviceDTO) (*domain.Device, error) {
	res, err := a.BeforeDeviceStateCheck(ctx, id, false, func(c context.Context) (interface{}, error) {
		current, fetchErr := a.repo.FindByID(c, id)
		if fetchErr != nil {
			return nil, fetchErr
		}

		// Rule: Name and brand properties cannot be updated if the device is in use
		if current.State == domain.InUse && (dto.Name != current.Name || dto.Brand != current.Brand) {
			return nil, errors.New(domain.ErrDeviceInUseLocked)
		}

		return a.next.FullUpdate(c, id, dto)
	})
	if err != nil {
		return nil, err
	}
	return res.(*domain.Device), nil
}

func (a *DeviceServiceAspect) PartialUpdate(ctx context.Context, id int64, dto map[string]interface{}) (*domain.Device, error) {
	// Rule: Creation time cannot be updated (Intercept before database load)
	if _, creationTimeAttempt := dto["creationTime"]; creationTimeAttempt {
		return nil, errors.New(domain.ErrCreationTimeImmutable)
	}

	res, err := a.BeforeDeviceStateCheck(ctx, id, false, func(c context.Context) (interface{}, error) {
		current, fetchErr := a.repo.FindByID(c, id)
		if fetchErr != nil {
			return nil, fetchErr
		}

		// Rule: Name and brand properties cannot be updated if the device is in use
		if current.State == domain.InUse {
			_, nameExists := dto["name"]
			_, brandExists := dto["brand"]
			if nameExists || brandExists {
				return nil, errors.New(domain.ErrDeviceInUseLocked)
			}
		}

		return a.next.PartialUpdate(c, id, dto)
	})
	if err != nil {
		return nil, err
	}
	return res.(*domain.Device), nil
}

func (a *DeviceServiceAspect) Delete(ctx context.Context, id int64) error {
	_, err := a.BeforeDeviceStateCheck(ctx, id, true, func(c context.Context) (interface{}, error) {
		return nil, a.next.Delete(c, id)
	})
	return err
}

// Pass-through methods
func (a *DeviceServiceAspect) Create(ctx context.Context, dto domain.DeviceDTO) (*domain.Device, error) {
	return a.next.Create(ctx, dto)
}

func (a *DeviceServiceAspect) GetAll(ctx context.Context, limit int, cursor int64) ([]domain.Device, error) {
	return a.next.GetAll(ctx, limit, cursor)
}

func (a *DeviceServiceAspect) GetByID(ctx context.Context, id int64) (*domain.Device, error) {
	return a.next.GetByID(ctx, id)
}

func (a *DeviceServiceAspect) GetByBrand(ctx context.Context, brand string, limit int, cursor int64) ([]domain.Device, error) {
	return a.next.GetByBrand(ctx, brand, limit, cursor)
}

func (a *DeviceServiceAspect) GetByState(ctx context.Context, state domain.DeviceState, limit int, cursor int64) ([]domain.Device, error) {
	return a.next.GetByState(ctx, state, limit, cursor)
}
