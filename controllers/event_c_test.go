package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/thunthup/aimet-test/configs"
	"github.com/thunthup/aimet-test/models"
	"gotest.tools/v3/assert"
)

func init() {
	configs.LoadEnvVar()
	configs.ConnectPostgresDB()
}

func TestCreateEvent(t *testing.T) {
	// Setup
	r := gin.Default()
	r.POST("/events", CreateEvent)
	db := configs.DB
	// Test case 1: valid input
	requestBody := []byte(`{"title": "Test Event", "event_date": "4000-05-15", "start_time": "15:00:00+07", "end_time": "16:00:00+07"}`)
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	// Unmarshal response body
	type EventResp struct {
		Title     string `json:"title"`
		EventDate string `json:"event_date"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	var event EventResp
	err := json.Unmarshal(resp.Body.Bytes(), &event)
	assert.NilError(t, err)

	// Assert response body
	expectedEvent := EventResp{
		Title:     "Test Event",
		EventDate: "4000-05-15",
		StartTime: "15:00:00+07",
		EndTime:   "16:00:00+07",
	}
	assert.DeepEqual(t, expectedEvent, event)

	db.Delete(&models.Event{}, "title = ?", "Test Event")

	// Test case 2: invalid input (end time before start time)
	requestBody = []byte(`{"title": "Test Event", "event_date": "4000-05-15", "start_time": "16:00:00+07", "end_time": "15:00:00+07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"End time must be after start time"}`, resp.Body.String())

	// Test case 3: invalid input (invalid date format)
	requestBody = []byte(`{"title": "Test Event", "event_date": "4000-05-35", "start_time": "15:00:00-07", "end_time": "16:00:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid event date format"}`, resp.Body.String())

	// Test case 4: invalid input (invalid start time format)
	requestBody = []byte(`{"title": "Test Event", "event_date": "4000-05-35", "start_time": "1a5:00:00-07", "end_time": "16:00:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid start time format"}`, resp.Body.String())

	// Test case 5: invalid input (invalid end time format)
	requestBody = []byte(`{"title": "Test Event", "event_date": "4000-05-35", "start_time": "15:00:00-07", "end_time": "1s6:00:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	fmt.Println(resp.Body.String())
	assert.Equal(t, `{"error":"Invalid end time format"}`, resp.Body.String())

	// Test case 6: invalid input (missing title)
	requestBody = []byte(`{ "event_date": "4000-05-12", "start_time": "15:00:00-07", "end_time": "1s6:00:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	// fmt.Println(resp.Body.String())
	assert.Equal(t, `{"error":"Key: 'Event.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`, resp.Body.String())

	// Test case 7: overlapping time
	requestBody = []byte(`{ "title": "Test Event", "event_date": "4000-05-12", "start_time": "15:00:00-07", "end_time": "16:00:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	requestBody = []byte(`{ "title": "Test Event overlapped", "event_date": "4000-05-12", "start_time": "13:00:00-07", "end_time": "15:30:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	db.Delete(&models.Event{}, "title = ?", "Test Event")
	db.Delete(&models.Event{}, "title = ?", "Test Event overlapped")
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Event time is overlapping with existing events"}`, resp.Body.String())

}
