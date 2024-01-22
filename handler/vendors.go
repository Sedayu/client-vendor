package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Sedayu/client-vendor/entity"
	"github.com/Sedayu/client-vendor/service"
	"github.com/labstack/echo/v4"
)

type Vendors struct {
	vendorFinderService  service.VendorFinderServiceInterface
	vendorCreatorService service.VendorCreatorServiceInterface
	vendorUpdaterService service.VendorUpdaterServiceInterface
}

func NewVendors(vendorFinderService service.VendorFinderServiceInterface,
	vendorCreatorService service.VendorCreatorServiceInterface,
	vendorUpdaterService service.VendorUpdaterServiceInterface) *Vendors {
	return &Vendors{
		vendorFinderService:  vendorFinderService,
		vendorCreatorService: vendorCreatorService,
		vendorUpdaterService: vendorUpdaterService,
	}
}

func (s *Vendors) GetVendors(c echo.Context) error {
	var limit int
	var offset int
	var err error

	limit, err = strconv.Atoi(c.Request().URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	offset, err = strconv.Atoi(c.Request().URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	data, err := s.vendorFinderService.GetVendors(c.Request().Context(), limit, offset)
	if err != nil {
		log.Default().Printf("Error getting vendors data: %v", err)
		if errors.Is(err, entity.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"errors": []map[string]interface{}{
					{
						"error": "vendors not found",
					},
				},
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "internal server error",
				},
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
	})
}

func (s *Vendors) CreateVendor(c echo.Context) error {
	var vendor entity.Vendor
	if err := c.Bind(&vendor); err != nil {
		log.Default().Printf("Error binding vendor data: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "bad request",
				},
			},
		})
	}

	vendorID, err := s.vendorCreatorService.CreateVendor(c.Request().Context(), vendor)
	if err != nil {
		log.Default().Printf("Error creating vendor: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "internal server error",
				},
			},
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"vendor_id": vendorID,
	})
}

func (s *Vendors) UpdateVendor(c echo.Context) error {
	var vendor entity.Vendor
	if err := c.Bind(&vendor); err != nil {
		log.Default().Printf("Error binding vendor data: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "bad request",
				},
			},
		})
	}

	vendorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Default().Printf("Error parsing vendor ID: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "invalid vendor ID",
				},
			},
		})
	}

	vendor.ID = int(vendorID)

	err = s.vendorUpdaterService.UpdateVendor(c.Request().Context(), vendor)
	if err != nil {
		log.Default().Printf("Error updating vendor: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "internal server error",
				},
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "vendor updated successfully",
	})
}

func (s *Vendors) GetVendorByID(c echo.Context) error {
	vendorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Default().Printf("Error parsing vendor ID: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "invalid vendor ID",
				},
			},
		})
	}

	vendor, err := s.vendorFinderService.GetVendorByID(c.Request().Context(), vendorID)
	if err != nil {
		log.Default().Printf("Error getting vendor by ID: %v", err)
		if errors.Is(err, entity.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"errors": []map[string]interface{}{
					{
						"error": "vendor not found",
					},
				},
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "internal server error",
				},
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": vendor,
	})
}

func (s *Vendors) UpdateVendorEmail(c echo.Context) error {
	vendorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Default().Printf("Error parsing vendor ID: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "invalid vendor ID",
				},
			},
		})
	}

	// Extract email from the request body
	var updateEmailRequest struct {
		Email string `json:"email"`
	}

	if err := c.Bind(&updateEmailRequest); err != nil {
		log.Default().Printf("Error binding update email request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "bad request",
				},
			},
		})
	}

	// Validate email presence
	if updateEmailRequest.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "email is required in the request body",
				},
			},
		})
	}

	// Call the service to update the email
	err = s.vendorUpdaterService.UpdateVendorEmail(c.Request().Context(), int(vendorID), updateEmailRequest.Email)
	if err != nil {
		log.Default().Printf("Error updating vendor email: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"error": "internal server error",
				},
			},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "vendor email updated successfully",
	})
}
