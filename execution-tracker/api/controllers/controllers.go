package controllers

import (
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"eywa/execution-tracker/db"
	"eywa/execution-tracker/types"
	"eywa/go-libs/auth"
)

// GetInvocationList returns a list of all users invocations
func GetInvocationList(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	db := c.Get("db").(*db.Client)
	perPage := c.Get("per_page").(int)
	pageNumber := c.Get("page_number").(int)
	query := c.QueryParam("query")

	timelines, total, err := db.GetTimelines(auth.UserID, query, perPage, pageNumber)
	if err != nil {
		log.Errorf("Error getting timeline logs: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	mappedTimelines := make(map[string][]types.TimelineLog)
	for _, t := range timelines {
		mappedTimelines[t.RequestID] = append(mappedTimelines[t.RequestID], t)
	}

	briefs := []types.TimelineBrief{}
	for rID, tl := range mappedTimelines {
		brief := types.TimelineBrief{
			RequestID:    rID,
			FunctionName: tl[0].EventName,
			FunctionID:   tl[0].FunctionID,
		}

		startedAt := tl[0].Timestamp
		lastEventAt := tl[len(tl)-1].Timestamp

		brief.Duration = lastEventAt.Sub(startedAt).Milliseconds()
		brief.Age = startedAt
		brief.Status = tl[len(tl)-1].Response

		for _, t := range tl {
			if isErrorResponse(t.Response) {
				brief.IsError = true
				break
			}
		}

		briefs = append(briefs, brief)
	}

	sort.Slice(briefs, func(i, j int) bool {
		return briefs[i].Age.After(briefs[j].Age)
	})

	return c.JSON(http.StatusOK, types.TimelineLogsResponse{
		Page:       pageNumber,
		PerPage:    perPage,
		TotalCount: total,
		Objects:    briefs,
	})
}

// GetInvocation returns detailed information about a specific invocation
func GetInvocation(c echo.Context) error {
	auth := c.Get("auth").(*auth.Auth)
	db := c.Get("db").(*db.Client)
	requestID := c.Param("request_id")

	timelines, err := db.GetTimeline(auth.UserID, requestID)
	if err != nil {
		log.Errorf("Failed to get timeline for request %q: %s", requestID, err)
		return err
	}

	if len(timelines) == 0 {
		return c.JSON(http.StatusNotFound, "Timeline Not Found")
	}

	tld := types.TimelineDetails{
		RequestID: timelines[0].RequestID,
		Method:    timelines[0].Method,
		Response:  timelines[len(timelines)-1].Response,
		Age:       timelines[0].Timestamp,
	}

	initEvent := types.EventDetails{
		Name:      timelines[0].EventName,
		Response:  timelines[0].Response,
		Duration:  timelines[0].Duration,
		IsError:   isErrorResponse(timelines[0].Response),
		Timestamp: timelines[0].Timestamp,
	}
	tld.Events = []types.EventDetails{initEvent}

	initCreatedAt := timelines[0].Timestamp
	// Async execution
	if timelines[0].EventType == types.TimelineEventTypeQueued {
		queueDetails := types.EventDetails{
			Name: "Dwell Time",
		}
		// Only queued message
		if len(timelines) == 1 {
			queueDetails = types.EventDetails{
				Response: timelines[0].Response,
				Duration: time.Since(initCreatedAt).Milliseconds(),
			}
		} else {
			// Dequeue message is present
			queueDetails.Response = timelines[1].Response
			queueDetails.Duration = timelines[1].Duration
		}

		queueDetails.Timestamp = timelines[0].Timestamp
		queueDetails.IsError = isErrorResponse(queueDetails.Response)
		tld.Events = append(tld.Events, queueDetails)

		queueEndedAt := initCreatedAt.Add(time.Duration(queueDetails.Duration))

		mappedAttempts := make(map[string][]types.TimelineLog)
		for _, t := range timelines[2:] {
			mappedAttempts[t.EventName] = append(mappedAttempts[t.EventName], t)
		}

		for attempt, events := range mappedAttempts {
			// Either running or outright failed
			if len(events) == 1 {
				duration := time.Since(queueEndedAt).Milliseconds()
				// Failed events have duration set
				if events[0].EventType == types.TimelineEventTypeFailed {
					duration = events[0].Duration
				}
				tld.Events = append(tld.Events, types.EventDetails{
					Name:     attempt,
					Response: events[0].Response,
					Duration: duration,
					IsError:  isErrorResponse(events[0].Response),
				})
			} else if len(events) == 2 {
				// Could be failed or completed
				tld.Events = append(tld.Events, types.EventDetails{
					Name:      attempt,
					Response:  events[1].Response,
					Duration:  events[1].Duration,
					IsError:   isErrorResponse(events[1].Response),
					Timestamp: events[0].Timestamp,
				})
			}
		}
	}

	tld.Duration = time.Since(tld.Events[len(tld.Events)-1].Timestamp).Milliseconds()

	return c.JSON(http.StatusOK, tld)
}

func isErrorResponse(response int) bool {
	return response < 200 || response >= 400
}

// GetEventLogs returns event logs
func GetEventLogs(c echo.Context) error {
	db := c.Get("db").(*db.Client)
	perPage := c.Get("per_page").(int)
	pageNumber := c.Get("page_number").(int)

	var query types.EventLogsQuery
	err := c.Bind(&query)
	if err != nil {
		return err
	}

	top, bottom, err := checkTimeRange(c, query.TimestampMax, query.TimestampMin)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Validation error",
			"details": map[string]string{
				"timestamp": err.Error(),
			},
		})
	}

	query.TimestampMax = top
	query.TimestampMin = bottom

	alogs, total, err := db.GetEventLogs(query, pageNumber, perPage)
	if err != nil {
		log.Errorf("Error getting request logs: %s", err)
		return c.JSON(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, types.EventLogsResponse{
		Page:       pageNumber,
		PerPage:    perPage,
		Objects:    alogs,
		TotalCount: total,
	})
}

func checkTimeRange(c echo.Context, top, bottom time.Time) (time.Time, time.Time, error) {
	if bottom.IsZero() {
		if top.IsZero() {
			t := time.Now().UTC()
			top = t
			delta := top.Add(time.Hour * 1)
			bottom = delta
			return top, bottom, nil
		}

		return top, bottom, errors.New("Minimum time range value is required")
	}

	if top.IsZero() {
		t := time.Now().UTC()
		top = t
	}

	if bottom.After(top) {
		return top, bottom, errors.New("Invalid time range")
	}

	return top, bottom, nil
}
