package health_worker

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/pir5/health-worker/model"
)

// getHealthChecks is getting healthchecks.
// @Summary get healthchecks
// @Description get healthchecks
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id query int false "HealthCheck ID"
// @Param name query string false "Name"
// @Success 200 {array} model.HealthCheck
// @Failure 404 {object} health_worker.HTTPError
// @Failure 500 {object} health_worker.HTTPError
// @Router /healthchecks [get]
func (h *HealthCheckHandler) getHealthChecks(c echo.Context) error {
	whereParams := map[string]interface{}{}
	for k, v := range c.QueryParams() {
		if k != "id" && k != "name" {
			return c.JSON(http.StatusForbidden, nil)
		}
		whereParams[k] = v
	}

	ds, err := h.HealthCheckModeler.FindBy(whereParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ds)
}

// updateHealthCheck is update healthCheck.
// @Summary update healthCheck
// @Description update healthCheck
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id path string true "HealthCheck ID"
// @Param healthCheck body model.HealthCheck true "HealthCheck Object"
// @Success 200 {object} model.HealthCheck
// @Failure 403 {object} health_worker.HTTPError
// @Failure 404 {object} health_worker.HTTPError
// @Failure 500 {object} health_worker.HTTPError
// @Router /healthchecks/{id} [put]
func (h *HealthCheckHandler) updateHealthCheck(c echo.Context) error {
	nd := &model.HealthCheck{}
	if err := c.Bind(nd); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	updated, err := h.HealthCheckModeler.UpdateByID(c.Param("id"), nd)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !updated {
		return c.JSON(http.StatusNotFound, "healthchecks does not exists")
	}
	return c.JSON(http.StatusOK, nil)
}

// deleteHealthCheck is delete healthCheck.
// @Summary delete healthCheck
// @Description delete healthCheck
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id path string true "HealthCheck ID"
// @Success 204 {object} model.HealthCheck
// @Failure 403 {object} health_worker.HTTPError
// @Failure 404 {object} health_worker.HTTPError
// @Failure 500 {object} health_worker.HTTPError
// @Router /healthchecks/{id} [delete]
func (h *HealthCheckHandler) deleteHealthCheck(c echo.Context) error {
	deleted, err := h.HealthCheckModeler.DeleteByID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !deleted {
		return c.JSON(http.StatusNotFound, "healthchecks does not exists")
	}

	return c.NoContent(http.StatusNoContent)
}

// createHealthCheck is create healthCheck.
// @Summary create healthCheck
// @Description create healthCheck
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param healthCheck body model.HealthCheck true "HealthCheck Object"
// @Success 201 {object} model.HealthCheck
// @Failure 403 {object} health_worker.HTTPError
// @Failure 404 {object} health_worker.HTTPError
// @Failure 500 {object} health_worker.HTTPError
// @Router /healthchecks [post]
func (h *HealthCheckHandler) createHealthCheck(c echo.Context) error {
	d := &model.HealthCheck{}
	if err := c.Bind(d); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.HealthCheckModeler.Create(d); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

type HealthCheckHandler struct {
	HealthCheckModeler model.HealthCheckModeler
}

func NewHealthCheckHandler(d model.HealthCheckModeler) *HealthCheckHandler {
	return &HealthCheckHandler{
		HealthCheckModeler: d,
	}
}
func HealthCheckEndpoints(g *echo.Group, db *gorm.DB) {
	h := NewHealthCheckHandler(model.NewHealthCheckModeler(db))
	g.GET("/healthchecks", h.getHealthChecks)
	g.PUT("/healthchecks/:id", h.updateHealthCheck)
	g.DELETE("/healthchecks/:id", h.deleteHealthCheck)
	g.POST("/healthchecks", h.createHealthCheck)
}
