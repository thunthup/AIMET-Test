package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thunthup/aimet-test/configs"
	"github.com/thunthup/aimet-test/models"
)

// Get an event by ID
func GetEventById(c *gin.Context) {
	db := configs.DB
	id := c.Param("id")
	var event models.Event

	if err := db.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	eventDate, _ := time.Parse("2020-05-28T00:00:00Z", event.EventDate)

	event.EventDate = eventDate.Format("2006-01-02")

	c.JSON(http.StatusOK, event)
}

// Create a new event
func CreateEvent(c *gin.Context) {
	db := configs.DB

	// Bind JSON request body to Event struct
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check that end time is after start time and is a valid time
	startTime, err := time.Parse("15:04:05-07", event.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
		return
	}
	endTime, err := time.Parse("15:04:05-07", event.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
		return
	}
	if endTime.Before(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
		return
	}

	// Check that event date is a valid date
	_, err = time.Parse("2006-01-02", event.EventDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event date format"})
		return
	}

	// Check if there is any event at the same day and overlapping time
	var count int64
	if err := db.Model(&models.Event{}).
		Where("event_date = ? AND start_time < ? AND end_time > ?", event.EventDate, event.EndTime, event.StartTime).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event time is overlapping with existing events"})
		return
	}

	// Create new event
	if err := db.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// Get events with filtering and searching
func ListEvents(c *gin.Context) {
	db := configs.DB
	var events []models.Event

	// Get query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	keyword := c.Query("keyword")
	sortOrder := c.DefaultQuery("sort_order", "asc")
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	// Parse date range parameters
	var startDate, endDate time.Time
	var err error
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
			return
		}
	}
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
			return
		}
	}

	//overide start and end date if year and month is provided

	if yearStr != "" && monthStr != "" {
		year, err := time.Parse("2006", yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
			return
		}
		month, err := time.Parse("01", monthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
			return
		}
		startDate = time.Date(year.Year(), month.Month(), 1, 0, 0, 0, 0, time.Now().Location())
		endDate = startDate.AddDate(0, 1, 0)
		endDate = endDate.Add(-time.Millisecond)

	}

	// overide if only the year is provided
	if yearStr != "" && monthStr == "" {
		year, err := time.Parse("2006", yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
			return
		}
		startDate = time.Date(year.Year(), 1, 1, 0, 0, 0, 0, time.Now().Location())
		endDate = startDate.AddDate(1, 0, 0)
		endDate = endDate.Add(-time.Millisecond)

	}

	// Build query
	query := db.Model(&models.Event{})
	if !startDate.IsZero() {
		query = query.Where("event_date >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("event_date <= ?", endDate)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	// Sort by event date and start time
	sortDirection := "ASC"
	if strings.ToLower(sortOrder) == "desc" {
		sortDirection = "DESC"
	}
	query = query.Order("event_date " + sortDirection + ", start_time " + sortDirection)

	// Execute query
	if err := query.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// for i := range events {
	// 	tempEventDate, _ := time.Parse("2020-05-28T00:00:00Z", events[i].EventDate)
	// 	events[i].EventDate = tempEventDate.Format("2006-01-02")
	// }

	c.JSON(http.StatusOK, events)
}

// Update an existing event
func UpdateEvent(c *gin.Context) {
	db := configs.DB

	// Get event ID from URL parameter
	eventID := c.Param("id")

	// Check if event exists
	var existingEvent models.Event
	if err := db.Where("id = ?", eventID).First(&existingEvent).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Bind JSON request body to Event struct
	var updatedEvent models.Event
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check that end time is after start time and is a valid time
	startTime, err := time.Parse("15:04:05-07", updatedEvent.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
		return
	}
	endTime, err := time.Parse("15:04:05-07", updatedEvent.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
		return
	}
	if endTime.Before(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
		return
	}

	// Check that event date is a valid date
	_, err = time.Parse("2006-01-02", updatedEvent.EventDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event date format"})
		return
	}

	// Check if there is any event at the same day and overlapping time
	var count int64
	if err := db.Model(&models.Event{}).
		Where("event_date = ? AND start_time < ? AND end_time > ? AND id <> ?", updatedEvent.EventDate, updatedEvent.EndTime, updatedEvent.StartTime, eventID).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event time is overlapping with existing events"})
		return
	}

	// Update existing event
	existingEvent.Title = updatedEvent.Title
	existingEvent.EventDate = updatedEvent.EventDate
	existingEvent.StartTime = updatedEvent.StartTime
	existingEvent.EndTime = updatedEvent.EndTime

	if err := db.Save(&existingEvent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, existingEvent)
}

// Delete an event by ID
func DeleteEvent(c *gin.Context) {
	db := configs.DB
	eventID := c.Param("id")

	// Check if event exists
	var event models.Event
	if err := db.Where("id = ?", eventID).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Delete event
	if err := db.Delete(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
