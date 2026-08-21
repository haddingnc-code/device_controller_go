package controller

import (
	"devices-api-go/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DeviceController handles incoming REST API endpoints for device lifecycle management.
type DeviceController struct {
	svc domain.DeviceService // Changed from *service.DeviceServiceAspect to the interface
}

// NewDeviceController creates a new instance of DeviceController.
func NewDeviceController(svc domain.DeviceService) *DeviceController {
	return &DeviceController{svc: svc}
}

// Create handles HTTP POST requests to register a new device resource.
// @Summary      Create a new device
// @Description  Persists a device into the database with an automatic creation time.
// @Tags         Devices API
// @Accept       json
// @Produce      json
// @Param        device  body      model.DeviceDTO  true  "Device Input Payload Data"
// @Success      201     {object}  model.Device
// @Failure      400     {object}  model.ApiError
// @Failure      500     {object}  model.ApiError
// @Router       /devices [post]
func (ctrl *DeviceController) Create(c *gin.Context) {
	var dto domain.DeviceDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	device, err := ctrl.svc.Create(c.Request.Context(), dto)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, device)
}

// GetAll handles HTTP GET requests to fetch a paginated collection block using cursors.
// @Summary      Fetch all devices with cursor pagination
// @Description  Retrieves a high-performance chunk of devices starting after a target cursor ID.
// @Tags         Devices API
// @Accept       json
// @Produce      json
// @Param        limit   query     int    false  "Max items to return (default 20)"
// @Param        cursor  query     int    false  "ID of the last item from the previous page (default 0)"
// @Success      200     {object}  model.CursorPageResult
// @Failure      500     {object}  model.ApiError
// @Router       /devices [get]
func (ctrl *DeviceController) GetAll(c *gin.Context) {
	limit, cursor := ctrl.getCursorParams(c)

	devices, err := ctrl.svc.GetAll(c.Request.Context(), limit, cursor)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	var nextCursor int64 = 0
	hasNext := false
	if len(devices) > 0 {
		// The next cursor identifier is always the ID of the very last element in the slice
		nextCursor = devices[len(devices)-1].ID
		hasNext = len(devices) == limit
	}

	c.JSON(http.StatusOK, domain.CursorPageResult{
		Content:    devices,
		NextCursor: nextCursor,
		HasNext:    hasNext,
	})
}

// getCursorParams parses query strings 'limit' and 'cursor' safely into numeric values.
func (ctrl *DeviceController) getCursorParams(c *gin.Context) (int, int64) {
	limitStr := c.DefaultQuery("limit", "20")
	cursorStr := c.DefaultQuery("cursor", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	cursor, err := strconv.ParseInt(cursorStr, 10, 64)
	if err != nil || cursor < 0 {
		cursor = 0
	}

	return limit, cursor
}

// GetByID handles HTTP GET requests to locate a single device by its ID key.
// @Summary      Fetch a single device by ID
// @Description  Retrieves a specific device resource matching the unique path variable identifier.
// @Tags         Devices API
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID of the target device to find"
// @Success      200  {object}  model.Device
// @Failure      400  {object}  model.ApiError
// @Failure      404  {object}  model.ApiError
// @Router       /devices/{id} [get]
func (ctrl *DeviceController) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	device, err := ctrl.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, device)
}

// FullUpdate handles HTTP PUT requests to perform complete resource modifications.
// @Summary      Fully update an existing device
// @Description  Replaces all fields of an existing device record. Protected by active state guards.
// @Tags         Devices API
// @Accept       json
// @Produce      json
// @Param        id      path      int              true  "ID of the target device"
// @Param        device  body      model.DeviceDTO  true  "Complete new state data payload"
// @Success      200     {object}  model.Device
// @Failure      400     {object}  model.ApiError
// @Failure      404     {object}  model.ApiError
// @Failure      500     {object}  model.ApiError
// @Router       /devices/{id} [put]
func (ctrl *DeviceController) FullUpdate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var dto domain.DeviceDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	device, err := ctrl.svc.FullUpdate(c.Request.Context(), id, dto)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, device)
}

// PartialUpdate handles HTTP PATCH requests to execute partial field changes.
// @Summary      Partially update an existing device
// @Description  Modifies specific payload key-value fields. Protected by active state guards.
// @Tags         Devices API
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "ID of the target device"
// @Param        updates  body      map[string]interface{} true  "Key-value properties map to update"
// @Success      200      {object}  model.Device
// @Failure      400      {object}  model.ApiError
// @Failure      404      {object}  model.ApiError
// @Failure      500      {object}  model.ApiError
// @Router       /devices/{id} [patch]
func (ctrl *DeviceController) PartialUpdate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	device, err := ctrl.svc.PartialUpdate(c.Request.Context(), id, updates)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, device)
}

// Delete handles HTTP DELETE requests to remove a device resource.
// @Summary      Delete a single device
// @Description  Permanently removes a specific device matching the path identifier key. Protected by active state guards.
// @Tags         Devices API
// @Param        id   path      int  true  "ID of the target device to delete"
// @Success      24   {string}  string "No Content"
// @Failure      400  {object}  model.ApiError
// @Failure      404  {object}  model.ApiError
// @Failure      500  {object}  model.ApiError
// @Router       /devices/{id} [delete]
func (ctrl *DeviceController) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := ctrl.svc.Delete(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusNoContent)
}

// GetByBrand handles HTTP GET requests to filter devices by brand with cursor pagination.
// @Summary      Fetch devices by brand with cursor pagination
// @Description  Filters the device collection case-insensitively by brand starting after a target cursor ID.
// @Tags         Devices API
// @Accept       json
// @Produce      json
// @Param        brand   query     string  true   "Name of the device brand"
// @Param        limit   query     int     false  "Max items to return (default 20)"
// @Param        cursor  query     int     false  "ID of the last item from the previous page (default 0)"
// @Success      200     {object}  model.CursorPageResult
// @Failure      400     {object}  model.ApiError
// @Failure      500     {object}  model.ApiError
// @Router       /devices/search/brand [get]
func (ctrl *DeviceController) GetByBrand(c *gin.Context) {
	brand := c.Query("brand")
	if brand == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand query parameter is required"})
		return
	}

	limit, cursor := ctrl.getCursorParams(c)
	devices, err := ctrl.svc.GetByBrand(c.Request.Context(), brand, limit, cursor)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	var nextCursor int64 = 0
	hasNext := false
	if len(devices) > 0 {
		nextCursor = devices[len(devices)-1].ID
		hasNext = len(devices) == limit
	}

	c.JSON(http.StatusOK, domain.CursorPageResult{
		Content:    devices,
		NextCursor: nextCursor,
		HasNext:    hasNext,
	})
}

// GetByState handles HTTP GET requests to filter devices by operational status with cursor pagination.
// @Summary      Fetch devices by state with cursor pagination
// @Description  Filters the device collection matching the target state enum starting after a target cursor ID.
// @Tags         Devices API
// @Accept       json
// @Produce      json
// @Param        state   query     string  true   "Target state (AVAILABLE, IN_USE, INACTIVE)"
// @Param        limit   query     int     false  "Max items to return (default 20)"
// @Param        cursor  query     int     false  "ID of the last item from the previous page (default 0)"
// @Success      200     {object}  model.CursorPageResult
// @Failure      400     {object}  model.ApiError
// @Failure      500     {object}  model.ApiError
// @Router       /devices/search/state [get]
func (ctrl *DeviceController) GetByState(c *gin.Context) {
	stateStr := c.Query("state")
	if stateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State query parameter is required"})
		return
	}

	limit, cursor := ctrl.getCursorParams(c)
	devices, err := ctrl.svc.GetByState(c.Request.Context(), domain.DeviceState(stateStr), limit, cursor)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	var nextCursor int64 = 0
	hasNext := false
	if len(devices) > 0 {
		nextCursor = devices[len(devices)-1].ID
		hasNext = len(devices) == limit
	}

	c.JSON(http.StatusOK, domain.CursorPageResult{
		Content:    devices,
		NextCursor: nextCursor,
		HasNext:    hasNext,
	})
}
