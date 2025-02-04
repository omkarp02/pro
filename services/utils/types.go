package utils

type UPDATE_RESULT struct {
	MatchedCount  int
	ModifiedCount int
}

type MODEL_STATUS string

const (
	ACTIVE_STATUS   MODEL_STATUS = "active"
	INACTIVE_STATUS MODEL_STATUS = "inActive"
	ARCHIVE_STATUS  MODEL_STATUS = "archive"
)
