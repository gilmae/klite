package parser

type StatementType string

const (
	STATEMENT_INSERT = "insert"
	STATEMENT_SELECT = "select"
)

const (
	COLUMN_USERNAME_SIZE = 255
	COLUMN_EMAIL_SIZE    = 255
)

// Hard coded for now while we work through the tutorial
type Row struct {
	Id       uint32
	Username [COLUMN_USERNAME_SIZE]byte
	Email    [COLUMN_EMAIL_SIZE]byte
}
