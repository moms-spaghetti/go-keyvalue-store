package protocols

import (
	"encoding/json"
	"errors"
	"net/http"
	"task1/internal/logger"
	"task1/internal/store"
)

var ErrRouteForbidden = errors.New("method forbidden")

func statusFromError(err error) int {
	switch {
	case errors.Is(err, nil):
		return http.StatusOK
	case errors.Is(err, store.ErrStoreKeyNotFound):
		return http.StatusNotFound
	case errors.Is(err, store.ErrKeyEmpty):
		return http.StatusBadRequest
	case errors.Is(err, ErrRouteForbidden):
		return http.StatusMethodNotAllowed
	default:
		return http.StatusInternalServerError
	}
}

func BuildJsonResponse(err error, data interface{}, logger *logger.Logger) (int, []byte) {
	res := jsonResponse{
		Err:    "",
		Status: statusFromError(err),
		Data:   data,
	}

	if err != nil {
		logger.Log("buildJsonResponseErr:" + err.Error())
		res.Err = err.Error()
		res.Data = nil
	}

	out, err1 := json.Marshal(res)
	if err1 != nil {
		logger.Log("jsonError encoding error")
		res.Err = err1.Error()
		res.Status = http.StatusInternalServerError
	}

	return res.Status, out
}
