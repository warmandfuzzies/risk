package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestNewRiskHandler(t *testing.T) {

	posturl := "http://localhost:8080/v1/risks"

	body := []byte(`{
		"state": "open",
		"title": "I am a title",
		"description": "I am a description"
	}`)

	req, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newRiskHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var risk Risk

	derr := json.NewDecoder(rr.Body).Decode(&risk)
	if derr != nil {
		t.Errorf("Unable to decode risk")
	}

	if risk.State != "open" {
		t.Errorf("handler returned unexpected state: got %v want %v",
			risk.State, "open")
	}

	if risk.Title != "I am a title" {
		t.Errorf("handler returned unexpected title: got %v want %v",
			risk.Title, "I am a title")
	}

	if risk.Description != "I am a description" {
		t.Errorf("handler returned unexpected description: got %v want %v",
			risk.Title, "I am a description")
	}

	delete(risks, risk.Id)
}

func TestGetRiskHandler(t *testing.T) {

	geturl := "http://localhost:8080/v1/risks/00112233-4455-6677-8899-aabbccddeeff"

	testUUID, err := uuid.Parse("00112233-4455-6677-8899-aabbccddeeff")
	if err != nil {
		panic(err)
	}

	risks[testUUID] = Risk{
		Id:          testUUID,
		State:       "open",
		Title:       "I am a title",
		Description: "I am a description",
	}

	req, err := http.NewRequest("GET", geturl, nil)
	if err != nil {
		panic(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "00112233-4455-6677-8899-aabbccddeeff"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getRiskHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var risk Risk

	derr := json.NewDecoder(rr.Body).Decode(&risk)
	if derr != nil {
		t.Errorf("Unable to decode risk")
	}

	if string(risk.Id.String()) != "00112233-4455-6677-8899-aabbccddeeff" {
		t.Errorf("handler returned unexpected Id: got %v want %v",
			risk.Id, "00112233-4455-6677-8899-aabbccddeeff")
	}

	if risk.State != "open" {
		t.Errorf("handler returned unexpected state: got %v want %v",
			risk.State, "open")
	}

	if risk.Title != "I am a title" {
		t.Errorf("handler returned unexpected title: got %v want %v",
			risk.Title, "I am a title")
	}

	if risk.Description != "I am a description" {
		t.Errorf("handler returned unexpected description: got %v want %v",
			risk.Title, "I am a description")
	}

	delete(risks, testUUID)

}

func TestListRisksHandler(t *testing.T) {

	geturl := "http://localhost:8080/v1/risks"

	testUUID, err := uuid.Parse("00112233-4455-6677-8899-aabbccddeeff")
	if err != nil {
		panic(err)
	}

	risks[testUUID] = Risk{
		Id:          testUUID,
		State:       "open",
		Title:       "I am a title",
		Description: "I am a description",
	}

	testUUID2, err := uuid.Parse("acde070d-8c4c-4f0d-9d8a-162843c10333")
	if err != nil {
		panic(err)
	}

	risks[testUUID2] = Risk{
		Id:          testUUID2,
		State:       "closed",
		Title:       "test title",
		Description: "test description",
	}

	req, err := http.NewRequest("GET", geturl, nil)
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(listRisksHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var riskValues []Risk

	derr := json.NewDecoder(rr.Body).Decode(&riskValues)
	if derr != nil {
		t.Errorf("Unable to decode risks")
	}

	if len(riskValues) != 2 {
		t.Errorf("less than 2 risks were returned: got %v want %v",
			len(riskValues), 2)
	}

	if riskValues[0].Id.String() != "00112233-4455-6677-8899-aabbccddeeff" &&
		riskValues[0].Id.String() != "acde070d-8c4c-4f0d-9d8a-162843c10333" {
		t.Errorf("handler returned unexpected Id: got %v want 00112233-4455-6677-8899-aabbccddeeff or acde070d-8c4c-4f0d-9d8a-162843c10333",
			riskValues[0].Id.String())
	}

	testuuidIndex := 0
	testuuid2Index := 0
	if riskValues[0].Id.String() == "00112233-4455-6677-8899-aabbccddeeff" {
		testuuidIndex = 0
		testuuid2Index = 1
	} else {
		testuuidIndex = 1
		testuuid2Index = 0
	}

	if riskValues[testuuidIndex].State != "open" {
		t.Errorf("handler returned unexpected state: got %v want %v",
			riskValues[testuuidIndex].State, "open")
	}

	if riskValues[testuuidIndex].Title != "I am a title" {
		t.Errorf("handler returned unexpected title: got %v want %v",
			riskValues[testuuidIndex].Title, "I am a title")
	}

	if riskValues[testuuidIndex].Description != "I am a description" {
		t.Errorf("handler returned unexpected description: got %v want %v",
			riskValues[testuuidIndex].Description, "I am a description")
	}

	if riskValues[testuuid2Index].State != "closed" {
		t.Errorf("handler returned unexpected state: got %v want %v",
			riskValues[testuuid2Index].State, "closed")
	}

	if riskValues[testuuid2Index].Title != "test title" {
		t.Errorf("handler returned unexpected title: got %v want %v",
			riskValues[testuuid2Index].Title, "test title")
	}

	if riskValues[testuuid2Index].Description != "test description" {
		t.Errorf("handler returned unexpected description: got %v want %v",
			riskValues[testuuid2Index].Description, "testdescription")
	}

	delete(risks, testUUID)
	delete(risks, testUUID2)

}
