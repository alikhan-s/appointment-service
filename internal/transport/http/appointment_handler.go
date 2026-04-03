package http

import (
	"net/http"

	"github.com/alikhan-s/appointment-s/internal/model"
	"github.com/alikhan-s/appointment-s/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	usecase usecase.AppointmentUseCase
}

func NewAppointmentHandler(r *gin.Engine, u usecase.AppointmentUseCase) {
	handler := &AppointmentHandler{usecase: u}

	r.POST("/appointments", handler.Create)
	r.GET("/appointments/:id", handler.GetByID)
	r.GET("/appointments", handler.GetAll)
	r.PATCH("/appointments/:id/status", handler.UpdateStatus)
}

func (h *AppointmentHandler) Create(c *gin.Context) {
	var appt model.Appointment
	if err := c.ShouldBindJSON(&appt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.usecase.Create(c.Request.Context(), &appt); err != nil {
		if err == model.ErrInvalidTitle || err == model.ErrInvalidDoctorID || err == model.ErrDoctorNotExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err == model.ErrDoctorServiceDown {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}

	c.JSON(http.StatusOK, appt)
}

func (h *AppointmentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	appt, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == model.ErrAppointmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "there is no appointment like this"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, appt)
}

func (h *AppointmentHandler) GetAll(c *gin.Context) {
	appointments, err := h.usecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}

	if appointments == nil {
		appointments = []*model.Appointment{}
	}

	c.JSON(http.StatusOK, appointments)
}

type updateStatusRequest struct {
	Status model.Status `json:"status"`
}

func (h *AppointmentHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req updateStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := h.usecase.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		if err == model.ErrAppointmentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "appointment not found"})
			return
		}
		if err == model.ErrInvalidStatus || err == model.ErrInvalidTransition {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}
