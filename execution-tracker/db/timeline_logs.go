package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"xorm.io/builder"

	"eywa/execution-tracker/types"
)

// GetTimelines returns all timelines from db
func (c *Client) GetTimelines(userID, functionID, filter string, perPage, pageNumber int) ([]types.TimelineLog, int, error) {
	subQuery := builder.
		Select("tl.request_id").
		From("timeline_logs", "tl").
		Where(builder.Eq{"tl.user_id": userID})

	subQuery = applyTimelineLogFilter(subQuery, "tl", functionID, filter)

	query := c.Builder().
		Select("tl.*").
		From("timeline_logs", "tl").
		Where(builder.In("tl.request_id", subQuery)).
		OrderBy("tl.timestamp")

	// copy the struct so we can replace the select fields
	countQuery := &builder.Builder{}
	*countQuery = *query

	countQuery.Select("tl.request_id, max(tl.timestamp) as maxtime").
		OrderBy("maxtime desc").
		GroupBy("tl.request_id")
	countQuery = c.Builder().
		Select("array_agg(request_id::text)::text[] ").
		From(countQuery, "tmp")

	sql, args, err := countQuery.ToSQL()
	if err != nil {
		return nil, 0, err
	}
	sql, err = builder.ConvertPlaceholder(sql, "$")
	if err != nil {
		return nil, 0, err
	}

	ids := []string{}
	err = sqlx.Get(c.ex, pq.Array(&ids), sql, args...)
	if err != nil {
		log.Debugf("ERROR Getting Totals: %s", err)
		return nil, 0, err
	}

	paginatedIds := paginate(ids, pageNumber, perPage)

	query = c.Builder().
		Select("tl.*").
		From(query, "tl").
		Where(builder.In("request_id", paginatedIds)).OrderBy("tl.timestamp")

	sql, args, err = query.ToSQL()
	if err != nil {
		log.Debugf("ERROR Getting Rows: %s", err)
		return nil, 0, err
	}
	sql, err = builder.ConvertPlaceholder(sql, "$")
	if err != nil {
		return nil, 0, err
	}

	timelineLogs := []types.TimelineLog{}
	err = sqlx.Select(c.ex, &timelineLogs, sql, args...)
	if err != nil {
		return nil, 0, err
	}

	return timelineLogs, len(ids), nil
}

func paginate(x []string, page int, perPage int) []string {
	skip := perPage * (page - 1)
	if skip > len(x) {
		page = len(x)
	}

	end := skip + perPage
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
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

func applyTimelineLogFilter(query *builder.Builder, name string, functionID, filter string) *builder.Builder {
	if functionID != "" {
		query = query.And(builder.Eq{name + ".function_id": functionID})
	}

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
