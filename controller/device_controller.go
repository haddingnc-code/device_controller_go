package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"devices-api-go/model"
	"devices-api-go/service"
)

// DeviceController handles incoming REST API endpoints for device lifecycle management.
type DeviceController struct {
	svc *service.DeviceService
}

// NewDeviceController creates a new instance of DeviceController.
func NewDeviceController(svc *service.DeviceService) *DeviceController {
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
	var dto model.DeviceDTO

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

// GetAll handles HTTP GET requests to fetch a paginated collection block of devices.
// @Summary      Fetch all devices with efficient pagination
// @Description  Retrieves a paginated chunk of devices using page and size query parameters.
// @Tags         Devices API
// @Accept       json
// @Produce      json
// @Param        page  query     int  false  "Page number (default 0)"
// @Param        size  query     int  false  "Page size (default 20)"
// @Success      200   {array}   model.Device
// @Failure      500   {object}  model.ApiError
// @Router       /devices [get]
func (ctrl *DeviceController) GetAll(c *gin.Context) {
	limit, offset := ctrl.getPaginationParams(c)

	devices, err := ctrl.svc.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, devices)
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

// GetByBrand handles HTTP GET requests to filter devices by brand with pagination.
// @Summary      Fetch devices by brand with pagination
// @Description  Filters the device collection case-insensitively by brand matching text.
// @Tags         Devices API
// @Accept       json
// @Produce      json
// @Param        brand  query     string  true   "Name of the device brand"
// @Param        page   query     int     false  "Page number"
// @Param        size   query     int     false  "Page size"
// @Success      200    {array}   model.Device
// @Failure      400    {object}  model.ApiError
// @Failure      500    {object}  model.ApiError
// @Router       /devices/search/brand [get]
func (ctrl *DeviceController) GetByBrand(c *gin.Context) {
	brand := c.Query("brand")
	if brand == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Brand query parameter is required"})
		return
	}

	limit, offset := ctrl.getPaginationParams(c)
	devices, err := ctrl.svc.GetByBrand(c.Request.Context(), brand, limit, offset)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, devices)
}

// GetByState handles HTTP GET requests to filter devices by operational status with pagination.
// @Summary      Fetch devices by state with pagination
// @Description  Filters the device collection matching the target state enum value.
// @Tags         Devices API
// @Accept       json
// @Produce      json
// @Param        state  query     string  true   "Target state (AVAILABLE, IN_USE, INACTIVE)"
// @Param        page   query     int     false  "Page number"
// @Param        size   query     int     false  "Page size"
// @Success      200    {array}   model.Device
// @Failure      400    {object}  model.ApiError
// @Failure      500    {object}  model.ApiError
// @Router       /devices/search/state [get]
func (ctrl *DeviceController) GetByState(c *gin.Context) {
	stateStr := c.Query("state")
	if stateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State query parameter is required"})
		return
	}

	limit, offset := ctrl.getPaginationParams(c)
	devices, err := ctrl.svc.GetByState(c.Request.Context(), model.DeviceState(stateStr), limit, offset)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, devices)
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

	var dto model.DeviceDTO
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

// getPaginationParams parses standard query variables 'page' and 'size' to compute limit/offset parameters.
func (ctrl *DeviceController) getPaginationParams(c *gin.Context) (int, int) {
	pageStr := c.DefaultQuery("page", "0")
	sizeStr := c.DefaultQuery("size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 0
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 20
	}

	limit := size
	offset := page * size

	return limit, offset
}
