package db

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"xorm.io/builder"

	"eywa/execution-tracker/types"
)

// GetEventLogs ...
func (c *Client) GetEventLogs(criteria types.EventLogsQuery, pageNumber, perPage int) ([]types.EventLog, int, error) {
	query := c.Builder().
		Select(`el.*`).
		From("event_logs el")
	query = applyEventLogFilter(query, "el", criteria)
	query = query.OrderBy("el.timestamp DESC")

	eventLogs := []types.EventLog{}
	total, err := c.SelectWithCount(&eventLogs, query, pageNumber, perPage)
	if err != nil {
		return nil, 0, err
	}

	return eventLogs, total, nil
}

// BulkInsertEventLogs ...
func (c *Client) BulkInsertEventLogs(records []types.EventLog) (int, error) {
	tx, err := c.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmnt, err := tx.Preparex(pq.CopyIn("event_logs",
		"request_id", "user_id", "type", "function_name", "function_id",
		"message", "is_error", "timestamp", "expires_at"))
	if err != nil {
		return 0, err
	}

	total := 0
	for _, record := range records {
		_, err = stmnt.Exec(record.RequestID, record.UserID,
			record.Type, record.FunctionName, record.FunctionID,
			record.Message, record.IsError, record.Timestamp,
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

func applyEventLogFilter(query *builder.Builder, name string, filter types.EventLogsQuery) *builder.Builder {

	if filter.UserID != "" {
		query = query.And(builder.Eq{name + ".user_id": filter.UserID})
	}
	if !filter.TimestampMax.IsZero() {
		query = query.And(builder.Lte{name + ".timestamp": filter.TimestampMax})
	}
	if filter.TimestampMin.IsZero() {
		query = query.And(builder.Gte{name + ".timestamp": filter.TimestampMin})
	}

	if filter.Query != "" {
		query = query.And(builder.Or(
			ILike{name + ".request_id::text", filter.Query},
			ILike{name + ".function_id::text", filter.Query},
			ILike{name + ".function_name::text", filter.Query},
			ILike{name + ".message", filter.Query},
		))
	}

	if filter.Level != "all" && filter.Level != "" {
		query = query.And(builder.Eq{name + ".type": filter.Level})
	}

	if filter.OnlyErrors {
		query = query.And(builder.Eq{name + ".is_error": true})
	}

	return query
}

func getStatusCriteria(status string) (gte, lt int) {
	res := strings.TrimRight(status, "x")
	num, err := strconv.Atoi(res)
	if err != nil { // should never happen
		panic(fmt.Sprintf("Failed to convert status to number %s", err))
	}

	switch len(res) {
	case 1:
		gte = num * 100
		lt = gte + 100
	case 3:
		gte = num
		lt = gte + 1
	default:
		// should never happen
		panic(fmt.Sprintf("Invalid status criteria %q", status))
	}

	return
}
