package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/logger"	
	"github.com/bonarizki-dat/boilerplate-gin-dat/pkg/types"
)

// Default HTTP status messages mapping
var defaultHTTPMessages = map[int]string{
	http.StatusOK:                  "OK",
	http.StatusCreated:             "Created",
	http.StatusNoContent:           "No Content",
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusConflict:            "Conflict",
	http.StatusUnprocessableEntity: "Unprocessable Entity",
	http.StatusTooManyRequests:     "Too Many Requests",
	http.StatusInternalServerError: "Internal Server Error",
	http.StatusBadGateway:          "Bad Gateway",
}

// FormatValidationErrors converts validator.ValidationErrors into a map of field -> error message.
// Provides human-readable validation messages for each tag.
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field()
			var msg string

			switch e.Tag() {
			case "required":
				msg = field + " is required"
			case "omitempty":
				msg = field + " is optional"
			case "email":
				msg = field + " must be a valid email address"
			case "url":
				msg = field + " must be a valid URL"
			case "uuid":
				msg = field + " must be a valid UUID"
			case "len":
				msg = field + " must be exactly " + e.Param() + " characters long"
			case "min":
				msg = field + " must be at least " + e.Param()
			case "max":
				msg = field + " must be at most " + e.Param()
			case "lt":
				msg = field + " must be less than " + e.Param()
			case "lte":
				msg = field + " must be less than or equal to " + e.Param()
			case "gt":
				msg = field + " must be greater than " + e.Param()
			case "gte":
				msg = field + " must be greater than or equal to " + e.Param()
			case "eq":
				msg = field + " must be equal to " + e.Param()
			case "ne":
				msg = field + " must not be equal to " + e.Param()
			case "oneof":
				msg = field + " must be one of [" + e.Param() + "]"
			case "datetime":
				msg = field + " must be in format " + e.Param()
			case "numeric":
				msg = field + " must be a numeric value"
			case "alpha":
				msg = field + " must contain only letters"
			case "alphanum":
				msg = field + " must contain only letters and numbers"
			case "boolean":
				msg = field + " must be a boolean value"
			case "ip":
				msg = field + " must be a valid IP address"
			case "ipv4":
				msg = field + " must be a valid IPv4 address"
			case "ipv6":
				msg = field + " must be a valid IPv6 address"
			case "cidr":
				msg = field + " must be a valid CIDR notation"
			default:
				msg = "Invalid value for " + field
			}

			errors[field] = msg

			// Log details for easier tracing
			logger.Errorf("[VALIDATION] field=%s tag=%s value=%v msg=%s",
				field, e.Tag(), e.Value(), msg)
		}
	} else if err != nil {
		logger.Errorf("[VALIDATION] %v", err)
	}

	return errors
}

// HandleErrors sends an error response in a standard format.
// If the error is a validation error, it will return detailed field errors.
func HandleErrors(c *gin.Context, code int, err error, message string) {
	if message == "" {
		if msg, ok := defaultHTTPMessages[code]; ok {
			message = msg
		} else {
			message = "Error"
		}
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		validationErrors := FormatValidationErrors(errs)
		c.JSON(code, types.ErrorResponse{
			Success: false,
			Message: message,
			Data:    nil,
			Errors:  validationErrors,
		})
		return
	}

	errMsg := message
	if err != nil {
		errMsg = err.Error()
	}

	c.JSON(code, types.ErrorResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  gin.H{"error": errMsg},
	})
}

// HandleSuccess sends a success response in a standard format.
// If message is empty, it will fall back to default HTTP messages.
func HandleSuccess(c *gin.Context, code int, data interface{}, message string) {
	if message == "" {
		if msg, ok := defaultHTTPMessages[code]; ok {
			message = msg
		} else {
			message = "Success"
		}
	}

	c.JSON(code, types.SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
	})
}

// Ok is a shortcut to send a 200 OK response
func Ok(c *gin.Context, data interface{}, message string) {
	HandleSuccess(c, http.StatusOK, data, message)
}

// Created is a shortcut to send a 201 Created response
func Created(c *gin.Context, data interface{}, message string) {
	HandleSuccess(c, http.StatusCreated, data, message)
}

// NoContent is a shortcut to send a 204 No Content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest is a shortcut to send a 400 Bad Request response
func BadRequest(c *gin.Context, err error, message string) {
	HandleErrors(c, http.StatusBadRequest, err, message)
}

// Unauthorized is a shortcut to send a 401 Unauthorized response
func Unauthorized(c *gin.Context, err error, message string) {
	HandleErrors(c, http.StatusUnauthorized, err, message)
}

// Forbidden is a shortcut to send a 403 Forbidden response
func Forbidden(c *gin.Context, err error, message string) {
	HandleErrors(c, http.StatusForbidden, err, message)
}

// NotFound is a shortcut to send a 404 Not Found response
func NotFound(c *gin.Context, err error, message string) {
	HandleErrors(c, http.StatusNotFound, err, message)
}

// Conflict is a shortcut to send a 409 Conflict response
func Conflict(c *gin.Context, err error, message string) {
	HandleErrors(c, http.StatusConflict, err, message)
}

// UnprocessableEntity is a shortcut to send a 422 Unprocessable Entity response
func UnprocessableEntity(c *gin.Context, err error, message string) {
	HandleErrors(c, http.StatusUnprocessableEntity, err, message)
}

// TooManyRequests is a shortcut to send a 429 Too Many Requests response
func TooManyRequests(c *gin.Context, err error, message string) {
	HandleErrors(c, http.StatusTooManyRequests, err, message)
}

// InternalServerError is a shortcut to send a 500 Internal Server Error response
func InternalServerError(c *gin.Context, err error, message string) {
	HandleErrors(c, http.StatusInternalServerError, err, message)
}

// BadGateway is a shortcut to send a 502 Bad Gateway response
func BadGateway(c *gin.Context, err error, message string) {
	HandleErrors(c, http.StatusBadGateway, err, message)
}

// RespondErrorIntrajasa is a helper to send standard Intrajasa error response
func RespondErrorIntrajasa(c *gin.Context, code string, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"responseData": gin.H{
			"responseCode": code,
			"responseMsg":  msg,
		},
	})
}