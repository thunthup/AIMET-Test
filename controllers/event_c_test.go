package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thunthup/aimet-test/configs"
	"github.com/thunthup/aimet-test/models"
	"gotest.tools/v3/assert"
)

func init() {
	var path = "../.env"
	configs.LoadEnvVar(&path)
	configs.ConnectPostgresDB()
	gin.SetMode(gin.TestMode)
}

// Unmarshal response body
type EventResp struct {
	Title     string `json:"title"`
	EventDate string `json:"event_date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func TestCreateEvent(t *testing.T) {
	// Setup
	r := gin.Default()
	r.POST("/events", CreateEvent)
	db := configs.DB
	// Test case 1: valid input
	requestBody := []byte(`{"title": "Test Event 9835-5dc547a01713", "event_date": "4000-05-15", "start_time": "15:00:00+07", "end_time": "16:00:00+07"}`)
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var event EventResp
	err := json.Unmarshal(resp.Body.Bytes(), &event)
	assert.NilError(t, err)

	// Assert response body
	expectedEvent := EventResp{
		Title:     "Test Event 9835-5dc547a01713",
		EventDate: "4000-05-15",
		StartTime: "15:00:00+07",
		EndTime:   "16:00:00+07",
	}
	assert.DeepEqual(t, expectedEvent, event)

	db.Delete(&models.Event{}, "title = ?", "Test Event 9835-5dc547a01713")

	// Test case 2: invalid input (end time before start time)
	requestBody = []byte(`{"title": "Test Event 9835-5dc547a01713", "event_date": "4000-05-15", "start_time": "16:00:00+07", "end_time": "15:00:00+07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"End time must be after start time"}`, resp.Body.String())

	// Test case 3: invalid input (invalid date format)
	requestBody = []byte(`{"title": "Test Event 9835-5dc547a01713", "event_date": "4000-05-35", "start_time": "15:00:00-07", "end_time": "16:00:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid event date format"}`, resp.Body.String())

	// Test case 4: invalid input (invalid start time format)
	requestBody = []byte(`{"title": "Test Event 9835-5dc547a01713", "event_date": "4000-05-35", "start_time": "1a5:00:00-07", "end_time": "16:00:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid start time format"}`, resp.Body.String())

	// Test case 5: invalid input (invalid end time format)
	requestBody = []byte(`{"title": "Test Event 9835-5dc547a01713", "event_date": "4000-05-35", "start_time": "15:00:00-07", "end_time": "1s6:00:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	assert.Equal(t, `{"error":"Invalid end time format"}`, resp.Body.String())

	// Test case 6: invalid input (missing title)
	requestBody = []byte(`{ "event_date": "4000-05-12", "start_time": "15:00:00-07", "end_time": "1s6:00:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	assert.Equal(t, `{"error":"Key: 'Event.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`, resp.Body.String())

	// Test case 7: overlapping time
	requestBody = []byte(`{ "title": "Test Event 9835-5dc547a01713", "event_date": "4000-05-12", "start_time": "15:00:00-07", "end_time": "16:00:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	requestBody = []byte(`{ "title": "Test Event overlapped 9835-5dc547a01713", "event_date": "4000-05-12", "start_time": "13:00:00-07", "end_time": "15:30:00-07"}`)
	req, _ = http.NewRequest("POST", "/events", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	db.Delete(&models.Event{}, "title = ?", "Test Event 9835-5dc547a01713")
	db.Delete(&models.Event{}, "title = ?", "Test Event overlapped 9835-5dc547a01713")
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Event time is overlapping with existing events"}`, resp.Body.String())

}

func TestGetEventById(t *testing.T) {
	// Setup
	r := gin.Default()
	r.GET("/events/:id", GetEventById)
	db := configs.DB

	// Create test event
	event := models.Event{
		Title:     "Test Event 9835-5dc547a01713",
		EventDate: "2022-05-15",
		StartTime: "15:00:00+07",
		EndTime:   "16:00:00+07",
	}
	db.Create(&event)

	// Test case 1: valid event ID
	req, _ := http.NewRequest("GET", "/events/"+strconv.Itoa(int(event.ID)), nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var eventResp EventResp
	err := json.Unmarshal(resp.Body.Bytes(), &eventResp)
	assert.NilError(t, err)

	// Assert response body
	expectedEvent := EventResp{
		Title:     "Test Event 9835-5dc547a01713",
		EventDate: "2022-05-15",
		StartTime: "15:00:00+07",
		EndTime:   "16:00:00+07",
	}
	assert.DeepEqual(t, expectedEvent, eventResp)
	invalidID := strconv.Itoa(int(event.ID))

	// Cleanup
	db.Delete(&event)

	// Test case 2: invalid event ID
	req, _ = http.NewRequest("GET", "/events/"+invalidID, nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, `{"error":"Event not found"}`, resp.Body.String())

}

func TestUpdateEvent(t *testing.T) {
	// Setup
	r := gin.Default()
	r.PUT("/events/:id", UpdateEvent)
	db := configs.DB

	// Create an event to update
	event := models.Event{
		Title:     "Test Event 9835-5dc547a01713",
		EventDate: "9999-06-15",
		StartTime: "15:00:00+07",
		EndTime:   "16:00:00+07",
	}
	overlappedEvent := models.Event{
		Title:     "Test Overlapped Event 9835-5dc547a01713",
		EventDate: "9999-06-17",
		StartTime: "15:00:00+07",
		EndTime:   "16:00:00+07",
	}
	nonExistEvent := models.Event{
		Title:     "Test Event 9835-5dc547a01713",
		EventDate: "9999-12-12",
		StartTime: "15:00:00+07",
		EndTime:   "16:00:00+07",
	}
	db.Create(&event)
	db.Create(&overlappedEvent)
	db.Create(&nonExistEvent)
	nonExistEventID := nonExistEvent.ID
	db.Delete(&nonExistEvent)
	//tear down
	defer db.Delete(&event)
	defer db.Delete(&overlappedEvent)

	// Test case 1: valid input
	requestBody := []byte(`{"title": "Updated Event 9835-5dc547a01713", "event_date": "9999-05-16", "start_time": "17:00:00+07", "end_time": "18:00:00+07"}`)
	req, _ := http.NewRequest("PUT", "/events/"+fmt.Sprintf("%d", event.ID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var updatedEvent models.Event
	db.Find(&updatedEvent, event.ID)
	tempEventDate, _ := time.Parse(time.RFC3339, updatedEvent.EventDate)

	updatedEvent.EventDate = tempEventDate.Format("2006-01-02")

	// Assert updated event
	expectedEvent := models.Event{
		ID:        event.ID,
		Title:     "Updated Event 9835-5dc547a01713",
		EventDate: "9999-05-16",
		StartTime: "17:00:00+07",
		EndTime:   "18:00:00+07",
		CreatedAt: updatedEvent.CreatedAt,
		UpdatedAt: updatedEvent.UpdatedAt,
		DeletedAt: updatedEvent.DeletedAt,
	}
	assert.DeepEqual(t, expectedEvent, updatedEvent)

	// Test case 2: invalid input (end time before start time)
	requestBody = []byte(`{"title": "Invalid Event 9835-5dc547a01713", "event_date": "9999-05-17", "start_time": "18:00:00+07", "end_time": "17:00:00+07"}`)
	req, _ = http.NewRequest("PUT", "/events/"+fmt.Sprintf("%d", event.ID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"End time must be after start time"}`, resp.Body.String())

	// Test case 3: invalid input (invalid date format)
	requestBody = []byte(`{"title": "Invalid Event 9835-5dc547a01713", "event_date": "9999-05-35", "start_time": "17:00:00+07", "end_time": "18:00:00+07"}`)
	req, _ = http.NewRequest("PUT", "/events/"+fmt.Sprintf("%d", event.ID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid event date format"}`, resp.Body.String())

	// Test case 4: invalid input (invalid start time format)
	requestBody = []byte(`{"title": "Test Event 9835-5dc547a01713", "event_date": "9999-05-35", "start_time": "1a5:00:00-07", "end_time": "16:00:00-07"}`)
	req, _ = http.NewRequest("PUT", "/events/"+fmt.Sprintf("%d", event.ID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid start time format"}`, resp.Body.String())

	// Test case 5: invalid input (invalid end time format)
	requestBody = []byte(`{"title": "Test Event 9835-5dc547a01713", "event_date": "9999-05-35", "start_time": "15:00:00-07", "end_time": "1s6:00:00-07"}`)
	req, _ = http.NewRequest("PUT", "/events/"+fmt.Sprintf("%d", event.ID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid end time format"}`, resp.Body.String())

	// Test case 6: invalid input (missing title)
	requestBody = []byte(`{ "event_date": "9999-05-12", "start_time": "15:00:00-07", "end_time": "1s6:00:00-07"}`)
	req, _ = http.NewRequest("PUT", "/events/"+fmt.Sprintf("%d", event.ID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Key: 'Event.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`, resp.Body.String())

	// Test case 7: overlapping time

	requestBody = []byte(`{ "title": "Overlapped Updated Event 9835-5dc547a01713", "event_date": "9999-06-17", "start_time": "13:00:00+07", "end_time": "16:30:00+07"}`)
	req, _ = http.NewRequest("PUT", "/events/"+fmt.Sprintf("%d", event.ID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Event time is overlapping with existing events"}`, resp.Body.String())

	// Test case 8: invalid input (non-exist event)

	requestBody = []byte(`{"title": "Updated Event 9835-5dc547a01713", "event_date": "9999-05-16", "start_time": "17:00:00+07", "end_time": "18:00:00+07"}`)
	req, _ = http.NewRequest("PUT", "/events/"+fmt.Sprintf("%d", nonExistEventID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

}

func TestListEvents(t *testing.T) {
	// Setup
	r := gin.Default()
	r.GET("/events", ListEvents)
	db := configs.DB

	// Create test events
	event1 := models.Event{
		Title:     "Test Ant Event 1 9835-5dc547a01713",
		EventDate: "9998-04-07",
		StartTime: "15:00:00+07",
		EndTime:   "16:00:00+07",
	}
	event2 := models.Event{
		Title:     "Test Antler Event 2 9835-5dc547a01713",
		EventDate: "9998-04-08",
		StartTime: "18:00:00+07",
		EndTime:   "19:00:00+07",
	}
	event3 := models.Event{
		Title:     "Test Bird Event 3 9835-5dc547a01713",
		EventDate: "9998-06-15",
		StartTime: "18:00:00+07",
		EndTime:   "19:00:00+07",
	}
	event4 := models.Event{
		Title:     "Test Cat Event 4 9835-5dc547a01713",
		EventDate: "9999-06-20",
		StartTime: "12:00:00+07",
		EndTime:   "14:00:00+07",
	}
	event5 := models.Event{
		Title:     "Test Dog Event 5 9835-5dc547a01713",
		EventDate: "9999-06-20",
		StartTime: "18:00:00+07",
		EndTime:   "19:00:00+07",
	}
	db.Create(&event1)
	db.Create(&event2)
	db.Create(&event3)
	db.Create(&event4)
	db.Create(&event5)

	//tear down
	defer db.Delete(&event1)
	defer db.Delete(&event2)
	defer db.Delete(&event3)
	defer db.Delete(&event4)
	defer db.Delete(&event5)

	// Test case 1: list all events with keyword, date range, sort order = desc
	req, _ := http.NewRequest("GET", "/events?keyword=9835-5dc547a01713&start_date=9998-04-08&end_date=9999-06-17&sort_order=desc", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Unmarshal response body
	var eventsResp []models.Event
	err := json.Unmarshal(resp.Body.Bytes(), &eventsResp)
	assert.NilError(t, err)

	// Assert response body
	expectedEvents := []models.Event{event3, event2}
	assert.Equal(t, len(expectedEvents), len(eventsResp))
	for i := range eventsResp {
		eventsResp[i].CreatedAt = expectedEvents[i].CreatedAt
		eventsResp[i].UpdatedAt = expectedEvents[i].UpdatedAt
	}
	assert.DeepEqual(t, expectedEvents, eventsResp)

	// Test case 2: invalid input (invalid start date format)
	req, _ = http.NewRequest("GET", "/events?keyword=9835-5dc547a01713&start_date=9998-04-0a8&end_date=9999-06-17", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid start date"}`, resp.Body.String())

	// Test case 3: invalid input (invalid end date format)
	req, _ = http.NewRequest("GET", "/events?keyword=9835-5dc547a01713&start_date=9998-04-08&end_date=9999-06-1a7", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid end date"}`, resp.Body.String())

	// Test case 4: filter by year only
	req, _ = http.NewRequest("GET", "/events?year=9998&keyword=9835-5dc547a01713", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Unmarshal response body
	err = json.Unmarshal(resp.Body.Bytes(), &eventsResp)
	assert.NilError(t, err)

	// Assert response body
	expectedEvents = []models.Event{event1, event2, event3}
	assert.Equal(t, len(expectedEvents), len(eventsResp))
	for i := range eventsResp {
		eventsResp[i].CreatedAt = expectedEvents[i].CreatedAt
		eventsResp[i].UpdatedAt = expectedEvents[i].UpdatedAt
	}
	assert.DeepEqual(t, expectedEvents, eventsResp)

	// Test case 5: filter by year and month
	req, _ = http.NewRequest("GET", "/events?year=9998&month=04&keyword=9835-5dc547a01713", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Unmarshal response body
	err = json.Unmarshal(resp.Body.Bytes(), &eventsResp)
	assert.NilError(t, err)

	// Assert response body
	expectedEvents = []models.Event{event1, event2}
	assert.Equal(t, len(expectedEvents), len(eventsResp))
	for i := range eventsResp {
		eventsResp[i].CreatedAt = expectedEvents[i].CreatedAt
		eventsResp[i].UpdatedAt = expectedEvents[i].UpdatedAt
	}
	assert.DeepEqual(t, expectedEvents, eventsResp)

	// Test case 6: filter by year with invalid year
	req, _ = http.NewRequest("GET", "/events?year=99988", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid year"}`, resp.Body.String())

	// Test case 7: filter by year and month with invalid year and valid month
	req, _ = http.NewRequest("GET", "/events?year=99988&month=12", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid year"}`, resp.Body.String())

	// Test case 8: filter by year and month with valid year and invalid month
	req, _ = http.NewRequest("GET", "/events?year=9998&month=13", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, `{"error":"Invalid month"}`, resp.Body.String())

}

func TestDeleteEvent(t *testing.T) {
	// Setup
	r := gin.Default()
	r.DELETE("/events/:id", DeleteEvent)
	db := configs.DB

	// Create an event to delete
	event := models.Event{
		Title:     "Test Event 9835-5dc547a01713",
		EventDate: "9999-06-15",
		StartTime: "15:00:00+07",
		EndTime:   "16:00:00+07",
	}
	db.Create(&event)
	eventID := event.ID

	// Test case 1: valid input
	req, _ := http.NewRequest("DELETE", "/events/"+fmt.Sprintf("%d", eventID), nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `{"message":"Event deleted successfully"}`, resp.Body.String())

	// Verify that the event has been deleted
	var deletedEvent models.Event
	if err := db.Where("id = ?", eventID).First(&deletedEvent).Error; err == nil {
		t.Errorf("Expected event with ID %d to be deleted, but found %+v", eventID, deletedEvent)
	}

	// Test case 2: event not found
	req, _ = http.NewRequest("DELETE", "/events/"+fmt.Sprintf("%d", eventID), nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, `{"error":"Event not found"}`, resp.Body.String())
}
