package entity

import (
	"encoding/json"
	"fmt"

	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/pedramktb/go-base-lib/types"
)

type SortDirection string

const (
	SortAsc  SortDirection = "+"
	SortDesc SortDirection = "-"
)

// Cursor is here because its meaningless without Sort.
type Sort[E Entity] struct {
	Field      string
	Direction  SortDirection
	CursorPart any
}

type Sorts[E Entity] []Sort[E]

func (s *Sorts[E]) FromDTO(sorts []string, cursor types.Nillable[string]) error {
	var cur []string
	if cursor.NotNil {
		// Split the cursor into parts
		cur = strings.Split(cursor.Val, ",")

		// Either null cursor or the same number of parts as sorters
		if len(cur) != len(sorts) {
			return ErrInvalidCursor
		}
	}

	ss := make(Sorts[E], len(sorts))
	for i := range sorts {
		if len(sorts[i]) == 0 {
			continue
		}
		var direction SortDirection
		if sorts[i][0] == '+' {
			direction = SortAsc
		} else if sorts[i][0] == '-' {
			direction = SortDesc
		}
		field := sorts[i][1:]
		// Create cursor with the same type as the field
		cursorPart, ok := (*new(E)).Fields()[field]
		if !ok {
			return ErrUnknownField.Wrap(fmt.Errorf("unknown field %q in sort parameters", field))
		}
		// If there is a cursor part, unmarshal it
		if len(cur) != 0 {
			if err := json.Unmarshal([]byte(cur[i]), &cursorPart); err != nil {
				return ErrInvalidCursor
			}
		}
		ss[i] = Sort[E]{
			Field:      field,
			Direction:  direction,
			CursorPart: cursorPart,
		}
	}

	*s = ss

	return nil
}

// NextCursor returns the cursor for the next list based on the next item
func (s Sorts[E]) NextCursor(nextItem types.Nillable[E]) Cursor {
	if nextItem.NotNil {
		cursor := make(Cursor, len(s))
		for i := range s {
			cursor[i] = nextItem.Val.Fields()[s[i].Field]
		}
		return cursor
	}
	return nil
}

func (s Sorts[E]) ToSQLQuery(query squirrel.SelectBuilder) squirrel.SelectBuilder {
	var cursorCondition squirrel.Sqlizer
	for i := range s {
		if s[i].CursorPart == nil {
			break
		}
		if i == 0 {
			if s[i].Direction == SortAsc {
				cursorCondition = squirrel.Gt{string(s[i].Field): s[i].CursorPart}
			} else {
				cursorCondition = squirrel.Lt{string(s[i].Field): s[i].CursorPart}
			}

		} else {
			var subCondition squirrel.Sqlizer
			for j := 0; j < i; j++ {
				if j == 0 {
					subCondition = squirrel.Eq{string(s[j].Field): s[j].CursorPart}
				} else {
					subCondition = squirrel.And{subCondition, squirrel.Eq{string(s[j].Field): s[i].CursorPart}}
				}
			}
			if s[i].Direction == SortAsc {
				subCondition = squirrel.And{subCondition, squirrel.GtOrEq{string(s[i].Field): s[i].CursorPart}}
			} else {
				subCondition = squirrel.And{subCondition, squirrel.LtOrEq{string(s[i].Field): s[i].CursorPart}}
			}
			cursorCondition = squirrel.Or{cursorCondition, subCondition}
		}
	}
	if cursorCondition != nil {
		query = query.Where(cursorCondition)
	}
	for _, s := range s {
		query = query.OrderBy(string(s.Field) + " " + string(s.Direction))
	}
	return query
}
