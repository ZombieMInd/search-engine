package sqlstore

const (
	InsertLog   string = `INSERT INTO log (uuid, ip, user_uuid, timestamp) VALUES ($1, $2, $3, TO_TIMESTAMP($4));`
	InsertEvent string = `INSERT INTO event (uuid, event_type, event_txt, log_uuid) VALUES ($1, $2, $3, $4);`
)
