package httpdoc

import (
	"net/http"

	"github.com/Aiya594/doctor-service/internal/repository"
	usecase "github.com/Aiya594/doctor-service/internal/use-case"
)

func parseError(err error) int {

	switch err {
	case repository.ErrNotFound:
		return http.StatusNotFound
	case usecase.ErrAlreadyExists:
		return http.StatusConflict
	case usecase.ErrInvalidFields:
		return http.StatusBadRequest

	default:
		return http.StatusInternalServerError

	}

}
