// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"net/http"

	moovhttp "github.com/moov-io/base/http"
	"github.com/moov-io/ofac"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

type customerWatchResponse struct {
	WatchID string `json:"watchId"`
}

func addCustomerRoutes(logger log.Logger, r *mux.Router) {
	r.Methods("GET").Path("/customers/{id}").HandlerFunc(getCustomer(logger))
	r.Methods("PUT").Path("/customers/{id}").HandlerFunc(updateCustomer(logger))

	r.Methods("POST").Path("/customers/{id}/watch").HandlerFunc(addCustomerWatch(logger))
	r.Methods("DELETE").Path("/customers/{id}/watch").HandlerFunc(removeCustomerWatch(logger))

	r.Methods("POST").Path("/customers/watch").HandlerFunc(addCustomerNameWatch(logger))
	r.Methods("DELETE").Path("/customers/watch/{watchId}").HandlerFunc(removeCustomerNameWatch(logger))
}

func getCustomer(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w, err := wrapResponseWriter(logger, w, r)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)

		customer := ofac.Customer{
			ID: "13ou1fohfkajfah", // "random"
			SDN: &ofac.SDN{
				EntityID: "306",
				SDNName:  "BANCO NACIONAL DE CUBA",
				SDNType:  "individual",
				Program:  "CUBA",
				Title:    "",
				Remarks:  "a.k.a. 'BNC'.",
			},
			Addresses: []*ofac.Address{
				{
					EntityID:                    "306",
					AddressID:                   "201",
					Address:                     "Dai-Ichi Bldg. 6th Floor, 10-2 Nihombashi, 2-chome, Chuo-ku",
					CityStateProvincePostalCode: "Tokyo 103",
					Country:                     "Japan",
				},
			},
			Alts: []*ofac.AlternateIdentity{
				{
					EntityID:      "306",
					AlternateID:   "220",
					AlternateType: "aka",
					AlternateName: "NATIONAL BANK OF CUBA",
				},
			},
		}
		if err := json.NewEncoder(w).Encode(customer); err != nil {
			moovhttp.Problem(w, err) // TODO(adam): replace with wrapped moovhttp.InternalError
			return
		}
	}
}

func addCustomerNameWatch(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w, err := wrapResponseWriter(logger, w, r)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)

		// TODO: read ?name=... param

		if err := json.NewEncoder(w).Encode(customerWatchResponse{"cust-name-watch"}); err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}

func addCustomerWatch(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w, err := wrapResponseWriter(logger, w, r)
		if err != nil {
			return
		}

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(customerWatchResponse{"cust-watch"}); err != nil {
			moovhttp.Problem(w, err)
			return
		}
	}
}

func updateCustomer(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w, err := wrapResponseWriter(logger, w, r)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusOK) // TODO
	}
}

func removeCustomerWatch(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w, err := wrapResponseWriter(logger, w, r)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusOK) // TODO
	}
}

func removeCustomerNameWatch(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w, err := wrapResponseWriter(logger, w, r)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusOK) // TODO
	}
}
