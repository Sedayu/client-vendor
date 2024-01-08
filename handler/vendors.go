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
}

func NewVendors(vendorFinderService service.VendorFinderServiceInterface, vendorCreatorService service.VendorCreatorServiceInterface) *Vendors {
	return &Vendors{
		vendorFinderService:  vendorFinderService,
		vendorCreatorService: vendorCreatorService,
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
