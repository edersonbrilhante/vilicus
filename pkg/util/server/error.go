package server

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

var validationErrors = map[string]string{
	"required": " is required, but was not received",
}

func getVldErrorMsg(s string) string {
	if v, ok := validationErrors[s]; ok {
		return v
	}
	return " failed on " + s + " validation"
}

type customErrHandler struct {
	e *echo.Echo
}

func (ce *customErrHandler) handler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	type resp struct {
		Message interface{} `json:"message"`
	}

	if ce.e.Debug {
		msg = err.Error()
		switch err.(type) {
		case *echo.HTTPError:
			code = err.(*echo.HTTPError).Code
		case validator.ValidationErrors:
			code = http.StatusBadRequest
		}
	} else {
		switch e := err.(type) {
		case *echo.HTTPError:
			code = e.Code
			msg = e.Message
			if e.Internal != nil {
				msg = fmt.Sprintf("%v, %v", err, e.Internal)
			}
		case validator.ValidationErrors:
			var errMsg []string
			for _, v := range e {
				errMsg = append(errMsg, fmt.Sprintf("%s%s", v.Field(), getVldErrorMsg(v.ActualTag())))
			}
			msg = resp{Message: errMsg}
			code = http.StatusBadRequest
		default:
			msg = http.StatusText(code)
		}
		if _, ok := msg.(string); ok {
			msg = resp{Message: msg}
		}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == "HEAD" {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			ce.e.Logger.Error(err)
		}
	}
}
