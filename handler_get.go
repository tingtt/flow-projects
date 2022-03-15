package main

import (
	"flow-projects/jwt"
	"flow-projects/project"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type QueryParam struct {
	ShowHidden bool    `query:"show_hidden" validate:"omitempty"`
	Embed      *string `query:"embed" validate:"omitempty,oneof=sub_projects"`
}

func get(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*jwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": err.Error()}, "	")
	}

	// Bind query
	q := new(QueryParam)
	if err = c.Bind(q); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate query
	if err = c.Validate(q); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	if q.Embed == nil {
		// Get projects
		projects, err := project.GetList(userId, q.ShowHidden)
		if err != nil {
			// 500: Internal server error
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		if projects == nil {
			return c.JSONPretty(http.StatusOK, []interface{}{}, "	")
		}
		return c.JSONPretty(http.StatusOK, projects, "	")
	}

	// Get projects with sub
	projects, err := project.GetListEmbed(userId, q.ShowHidden)
	if err != nil {
		// 500: Internal server error
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	if projects == nil {
		return c.JSONPretty(http.StatusOK, []interface{}{}, "	")
	}
	return c.JSONPretty(http.StatusOK, projects, "	")
}
