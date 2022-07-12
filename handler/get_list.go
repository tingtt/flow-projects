package handler

import (
	"flow-projects/flags"
	"flow-projects/jwt"
	"flow-projects/project"
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type QueryParam struct {
	Name       *string `query:"name" validate:"omitempty"`
	ShowHidden bool    `query:"show_hidden" validate:"omitempty"`
	Embed      *string `query:"embed" validate:"omitempty,oneof=sub_projects"`
}

func GetList(c echo.Context) error {
	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
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
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	if q.Embed == nil {
		// Get projects
		projects, err := project.GetList(userId, q.ShowHidden, q.Name)
		if err != nil {
			// 500: Internal server error
			c.Logger().Error(err)
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
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	if projects == nil {
		return c.JSONPretty(http.StatusOK, []interface{}{}, "	")
	}
	return c.JSONPretty(http.StatusOK, projects, "	")
}
