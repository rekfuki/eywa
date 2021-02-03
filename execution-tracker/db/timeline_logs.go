package db

import (
	"github.com/lib/pq"
	"xorm.io/builder"

	"eywa/execution-tracker/types"
)

// GetTimelines returns all timelines from db
func (c *Client) GetTimelines(userID, filter string, perPage, pageNumber int) ([]types.TimelineLog, int, error) {

	subQuery := builder.
		Select("tl.request_id").
		From("timeline_logs", "tl").
		Where(builder.Eq{"tl.user_id": userID})

	subQuery = applyTimelineLogFilter(subQuery, "tl", filter)

	query := c.Builder().Select("tl.*").
		From("timeline_logs", "tl").
		Where(builder.In("tl.request_id", subQuery)).
		OrderBy("tl.timestamp")

	timelineLogs := []types.TimelineLog{}
	total, err := c.SelectWithCount(&timelineLogs, query, pageNumber, perPage)
	if err != nil {
		return nil, 0, err
	}

	return timelineLogs, total, nil
}

// GetTimeline returns timeline related to request and user ids
func (c *Client) GetTimeline(userID, requestID string) ([]types.TimelineLog, error) {
	query := c.Builder().Select("tl.*").
		From("timeline_logs", "tl").
		Where(builder.Eq{
			"tl.user_id":    userID,
			"tl.request_id": requestID,
		}).OrderBy("tl.timestamp")

	timelineLogs := []types.TimelineLog{}
	err := c.Select(&timelineLogs, query)
	if err != nil {
		return nil, err
	}

	return timelineLogs, nil
}

// BulkInsertTimelineLogs ...
func (c *Client) BulkInsertTimelineLogs(records []types.TimelineLog) (int, error) {
	tx, err := c.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmnt, err := tx.Preparex(pq.CopyIn("timeline_logs",
		"request_id", "user_id", "function_id",
		"event_name", "event_type", "response",
		"method", "duration", "timestamp",
		"expires_at"))
	if err != nil {
		return 0, err
	}

	total := 0
	for _, record := range records {
		_, err = stmnt.Exec(record.RequestID, record.UserID,
			record.FunctionID, record.EventName, record.EventType,
			record.Response, record.Method, record.Duration, record.Timestamp,
			record.ExpiresAt)
		if err != nil {
			return 0, err
		}
		total++
	}

	_, err = stmnt.Exec()
	if err != nil {
		return 0, err
	}

	err = stmnt.Close()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return total, nil
}

func applyTimelineLogFilter(query *builder.Builder, name string, filter string) *builder.Builder {
	if filter != "" {
		query = query.And(builder.Or(
			ILike{name + ".request_id::text", filter},
			ILike{name + ".function_id::text", filter},
			ILike{name + ".event_name", filter},
			ILike{name + ".response::text", filter},
		))
	}

	return query
}
