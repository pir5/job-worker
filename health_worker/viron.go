package health_worker

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// vironAuthType
// @Summary get auth type
// @Description get auth type
// @ID viron_authtype#get
// @Accept  json
// @Produce  json
// @Router /viron_authtype [get]
// @Tags viron
func vironAuthType(c echo.Context) error {
	encodedJSON := []byte(`
[
    {
      "type": "email",
      "provider": "viron-demo",
      "url": "/signin",
      "method": "POST",
    },
    {
      "type": "signout",
      "provider": "",
      "url": "",
      "method": "POST",
    },
]
`)
	return c.JSONBlob(http.StatusOK, encodedJSON)

}

//vironGlobalMenu
// @Summary get global menu
// @Description get global menu
// @ID viron#get
// @Accept json
// @Produce json
// @Router /viron [get]
// @Tags viron
func vironGlobalMenu(c echo.Context) error {
	encodedJSON := []byte(`{
  "theme": "standard",
  "color": "white",
  "name": "Viron example - local",
  "tags": [
    "healthchecks",
    "routingpolicies"
  ],
  "pages": [
    {
      "section": "manage",
      "id": "healthchecks",
      "name": "HealthChecks",
      "components": [
        {
          "api": {
            "method": "get",
            "path": "/"
          },
	  "query": [
	    { key: "id", type: "integer" },
	    { key: "name", type: "string" },
          ],
	  "primary": "id",
          "name": "HealthChecks",
	  "style": "table",
          "pagination": true,
	  "table_labels": [
	    "id",
            "name",
            "type",
            "check_interval",
            "threshould",
            "params"
	  ]
        }
      ]
    },
    {
      "section": "manage",
      "id": "routingpolicies",
      "name": "RoutingPolicies",
      "components": [
        {
          "api": {
            "method": "get",
            "path": "/routingpolicies"
          },
	  "query": [
	    { key: "id", type: "integer" },
	    { key: "record_id", type: "integer" },
	    { key: "health_check_id", type: "integer" }
          ],
          "name": "Record",
	  "style": "table",
	  "primary": "id",
	  "table_labels": [
	    "id",
	    "record_id",
            "health_check_id",
            "type"
	  ]
        }
      ]
    }
  ]
}`)
	return c.JSONBlob(http.StatusOK, encodedJSON)

}

func VironEndpoints(g *echo.Group) {
	g.GET("/viron", vironGlobalMenu)
	g.GET("/viron_authtype", vironAuthType)
}
