package barang

import (
	"net/http"

	"RTLS_API/pkg/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Get(c *gin.Context) {
	deviceID := c.Query("device_id")

	data, err := h.Service.GetBarang(deviceID)
	if err != nil {
		c.JSON(500, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(200, data)
}

func (h *Handler) Create(c *gin.Context) {
	var input models.InputTransaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	result, err := h.Service.CreateBarang(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("device_id")

	var payload models.UpdateTransaction
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"detail": err.Error()})
		return
	}

	update, err := h.Service.UpdateBarang(id, payload)
	if err != nil {
		c.JSON(400, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":        "Barang berhasil diupdate",
		"device_id":      id,
		"updated_fields": update,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("device_id")

	if err := h.Service.DeleteBarang(id); err != nil {
		c.JSON(404, gin.H{"detail": "Barang tidak ditemukan"})
		return
	}

	c.JSON(200, gin.H{
		"message":   "Barang berhasil dihapus",
		"device_id": id,
	})
}

func (h *Handler) ResetSystem(c *gin.Context) {
	if err := h.Service.ResetSystem(); err != nil {
		c.JSON(404, gin.H{
			"detail": "system sudah kosong / ter-reset",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "system berhasil di-reset (barang + counter dihapus)",
	})
}
