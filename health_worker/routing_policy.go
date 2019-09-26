package health_worker

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/pir5/health-worker/model"
)

// getRoutingPolicys is getting routingpolicies.
// @Summary get routingpolicies
// @Description get routingpolicies
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id query int false "RoutingPolicy ID"
// @Param name query string false "Name"
// @Success 200 {array} model.RoutingPolicy
// @Failure 404 {object} health_worker.HTTPError
// @Failure 500 {object} health_worker.HTTPError
// @Router /routingpolicies [get]
func (h *RoutingPolicyHandler) getRoutingPolicies(c echo.Context) error {
	whereParams := map[string]interface{}{}
	for k, v := range c.QueryParams() {
		if k != "id" && k != "name" {
			return c.JSON(http.StatusForbidden, nil)
		}
		whereParams[k] = v
	}

	ds, err := h.RoutingPolicyModel.FindBy(whereParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if ds == nil || len(ds) == 0 {
		return c.JSON(http.StatusNotFound, "routingpolicies does not exists")
	}

	return c.JSON(http.StatusOK, ds)
}

// updateRoutingPolicy is update healthCheck.
// @Summary update healthCheck
// @Description update healthCheck
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id path string true "RoutingPolicy ID"
// @Param healthCheck body model.RoutingPolicy true "RoutingPolicy Object"
// @Success 200 {object} model.RoutingPolicy
// @Failure 403 {object} health_worker.HTTPError
// @Failure 404 {object} health_worker.HTTPError
// @Failure 500 {object} health_worker.HTTPError
// @Router /routingpolicies/{id} [put]
func (h *RoutingPolicyHandler) updateRoutingPolicy(c echo.Context) error {
	nd := &model.RoutingPolicy{}
	if err := c.Bind(nd); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	updated, err := h.RoutingPolicyModel.UpdateByID(c.Param("id"), nd)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !updated {
		return c.JSON(http.StatusNotFound, "routingpolicies does not exists")
	}
	return c.JSON(http.StatusOK, nil)
}

// deleteRoutingPolicy is delete healthCheck.
// @Summary delete healthCheck
// @Description delete healthCheck
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param id path string true "RoutingPolicy ID"
// @Success 204 {object} model.RoutingPolicy
// @Failure 403 {object} health_worker.HTTPError
// @Failure 404 {object} health_worker.HTTPError
// @Failure 500 {object} health_worker.HTTPError
// @Router /routingpolicies/{id} [delete]
func (h *RoutingPolicyHandler) deleteRoutingPolicy(c echo.Context) error {
	deleted, err := h.RoutingPolicyModel.DeleteByID(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if !deleted {
		return c.JSON(http.StatusNotFound, "routingpolicies does not exists")
	}

	return c.NoContent(http.StatusNoContent)
}

// createRoutingPolicy is create healthCheck.
// @Summary create healthCheck
// @Description create healthCheck
// @Security ID
// @Security Secret
// @Accept  json
// @Produce  json
// @Param healthCheck body model.RoutingPolicy true "RoutingPolicy Object"
// @Success 201 {object} model.RoutingPolicy
// @Failure 403 {object} health_worker.HTTPError
// @Failure 404 {object} health_worker.HTTPError
// @Failure 500 {object} health_worker.HTTPError
// @Router /routingpolicies [post]
func (h *RoutingPolicyHandler) createRoutingPolicy(c echo.Context) error {
	d := &model.RoutingPolicy{}
	if err := c.Bind(d); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.RoutingPolicyModel.Create(d); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, nil)
}

type RoutingPolicyHandler struct {
	RoutingPolicyModel model.RoutingPolicyModel
}

func NewRoutingPolicyHandler(d model.RoutingPolicyModel) *RoutingPolicyHandler {
	return &RoutingPolicyHandler{
		RoutingPolicyModel: d,
	}
}
func RoutingPolicyEndpoints(g *echo.Group, db *gorm.DB) {
	h := NewRoutingPolicyHandler(model.NewRoutingPolicyModel(db, nil))
	g.GET("/routingpolicies", h.getRoutingPolicies)
	g.PUT("/routingpolicies/:id", h.updateRoutingPolicy)
	g.DELETE("/routingpolicies/:id", h.deleteRoutingPolicy)
	g.POST("/routingpolicies", h.createRoutingPolicy)
}
