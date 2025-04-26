package entity

import "github.com/pedramktb/go-base-lib/types"

const (
	PaginationDefaultLimit = 10
)

type PaginationLimit uint8

type PaginationMeta struct {
	Total uint64
	Next  Cursor
}

type Paginated[E Entity] struct {
	Items []E
	Meta  PaginationMeta
}

func (l *PaginationLimit) FromDTO(limit types.Nillable[uint8]) {
	*l = PaginationDefaultLimit
	if limit.NotNil {
		*l = PaginationLimit(limit.Val)
	}
}

type paginationMetaDTO struct {
	Total uint64                 `json:"total"`
	Next  types.Nillable[string] `json:"next"`
}

func (m *PaginationMeta) ToDTO() (*paginationMetaDTO, error) {
	next, err := m.Next.ToDTO()
	if err != nil {
		return nil, err
	}
	return &paginationMetaDTO{
		Total: m.Total,
		Next:  next,
	}, nil
}
