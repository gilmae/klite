package object

import (
	"bytes"
	"fmt"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ   = "INTEGER"
	STRING_OBJ    = "STRING"
	RECORDSET_OBJ = "RECORDSET"
	ERROR_OBJ     = "ERROR"
	NULL_OBJ      = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return fmt.Sprint(s.Value) }
func (s *String) Type() ObjectType { return STRING_OBJ }

type RecordSet struct {
	Value   [][]Object
	Columns []String
}

func (r *RecordSet) Inspect() string {
	var out bytes.Buffer

	for _, p := range r.Columns {
		out.WriteString(fmt.Sprintf("%-32v", p.Inspect()))
	}

	out.WriteString("\n")

	for _, row := range r.Value {
		for _, col := range row {
			out.WriteString(fmt.Sprintf("%-32v", strings.Trim(string(col.Inspect()), string([]byte{'\x00'}))))
		}
		out.WriteString("\n")
	}

	return out.String()
}
func (r *RecordSet) Type() ObjectType { return RECORDSET_OBJ }

type Error struct {
	Message string
	//TODO we're going to want stack traces
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
