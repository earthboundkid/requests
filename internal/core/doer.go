package core

import (
	"net/http"
)

type ResponseHandler = func(*http.Response) error

type DoerResponse int

const (
	DoerOK DoerResponse = iota
	DoerConnect
	DoerValidate
	DoerHandle
)

func Do(cl *http.Client, req *http.Request, validators []ResponseHandler, h ResponseHandler) (DoerResponse, error) {
	res, err := cl.Do(req)
	if err != nil {
		return DoerConnect, err
	}
	defer res.Body.Close()

	for _, v := range validators {
		if v == nil {
			continue
		}
		if err = v(res); err != nil {
			return DoerValidate, err
		}
	}
	if err = h(res); err != nil {
		return DoerHandle, err
	}

	return DoerOK, nil
}
