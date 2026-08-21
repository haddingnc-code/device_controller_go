package service

import (
	"context"
	"devices-api-go/internal/domain" // Certifica-se de que mapeia para o teu pacote domain atual
	"errors"
)

// JoinPoint represents the execution hook for the targeted core business service logic.
type JoinPoint func(ctx context.Context) (interface{}, error)

// DeviceServiceAspect acts as the AOP interceptor for the concrete DeviceService.
type DeviceServiceAspect struct {
	next *DeviceService
	repo domain.DeviceRepository // Changed from *repository.DeviceRepository to the interface
}

func NewDeviceServiceAspect(next *DeviceService, repo domain.DeviceRepository) *DeviceServiceAspect {
	return &DeviceServiceAspect{
		next: next,
		repo: repo,
	}
}

// BeforeDeviceStateCheck runs as a generic @Before aspect advice triggered prior to mutations.
func (a *DeviceServiceAspect) BeforeDeviceStateCheck(ctx context.Context, id int64, isDelete bool, proceed JoinPoint) (interface{}, error) {
	// 1. Advice Execution: Fetch current resource state snapshot to validate constraints
	current, err := a.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Aspect Rule 1: In use devices cannot be deleted
	if isDelete && current.State == domain.InUse {
		return nil, errors.New(domain.ErrDeviceInUseDelete)
	}

	// 2. Proceed to target: Execute the encapsulated JoinPoint method
	return proceed(ctx)
}

// --- ADVISED INTERCEPTED METHODS (AOP) ---

func (a *DeviceServiceAspect) FullUpdate(ctx context.Context, id int64, dto domain.DeviceDTO) (*domain.Device, error) {
	res, err := a.BeforeDeviceStateCheck(ctx, id, false, func(c context.Context) (interface{}, error) {
		current, fetchErr := a.repo.FindByID(c, id)
		if fetchErr != nil {
			return nil, fetchErr
		}

		// Aspect Rule 2: Name and brand properties cannot be updated if the device is in use
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
	res, err := a.BeforeDeviceStateCheck(ctx, id, false, func(c context.Context) (interface{}, error) {
		current, fetchErr := a.repo.FindByID(c, id)
		if fetchErr != nil {
			return nil, fetchErr
		}

		// Aspect Rule 3: Name and brand properties cannot be updated if the device is in use (Partial check)
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

// --- PASS-THROUGH METHODS (Bypassing state validation aspect interceptors) ---

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
