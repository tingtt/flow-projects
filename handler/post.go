package handler

import (
	"flow-projects/flags"
	"flow-projects/jwt"
	"flow-projects/project"
	"fmt"
	"net/http"
	"strings"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func Post(c echo.Context) error {
	// Check `Content-Type`
	if !strings.Contains(c.Request().Header.Get("Content-Type"), "application/json") {
		// 415: Invalid `Content-Type`
		return c.JSONPretty(http.StatusUnsupportedMediaType, map[string]string{"message": "unsupported media type"}, "	")
	}

	// Check token
	u := c.Get("user").(*jwtGo.Token)
	userId, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// Bind request body
	post := new(project.PostBody)
	if err = c.Bind(post); err != nil {
		// 400: Bad request
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// Validate request body
	if err = c.Validate(post); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	p, usedName, parentNotFound, parentHasParent, err := project.Post(userId, *post)
	if err != nil {
		// 500: Internal server error
		c.Logger().Error(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if usedName {
		// 400: Bad Request
		c.Logger().Debugf("project name: %s already used", post.Name)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("project name: %s already used", post.Name)}, "	")
	}
	if parentNotFound && post.ParentId != nil {
		// 400: Bad Request
		c.Logger().Debugf("project id: %d does not exists", *post.ParentId)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": fmt.Sprintf("project id: %d does not exists", *post.ParentId)}, "	")
	}
	if parentHasParent && post.ParentId != nil {
		// 400: Bad Request
		c.Logger().Debug("cannot create child's child project")
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": "cannot create child's child project"}, "	")
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, p, "	")
}
