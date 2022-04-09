package plouf

import "github.com/labstack/echo/v4"

func ValidateAndBind(c echo.Context, value interface{}) error {
	if err := c.Bind(value); err != nil {
		return err
	}

	if err := c.Validate(value); err != nil {
		return err
	}

	return nil
}
