DROP INDEX IF EXISTS verification_events_request_id_idx;
DROP INDEX IF EXISTS verification_sessions_request_id_idx;
DROP INDEX IF EXISTS verification_requests_status_idx;

DROP TABLE IF EXISTS verification_events;
DROP TABLE IF EXISTS verification_sessions;
DROP TABLE IF EXISTS verification_requests;
