package inspection

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Status string

const (
	Open Status = "open"
	Done Status = "done"
)

var ErrStatusInvalid = errors.New("巡检状态无效")

type Definition struct {
	Status      Status `json:"status"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var definitions = []Definition{
	{Status: Open, Name: "待处理", Description: "巡检任务仍在处理流程中"},
	{Status: Done, Name: "已完成", Description: "巡检任务已经完成"},
}

func ParseStatus(text string) (Status, error) {
	status := Status(strings.ToLower(strings.TrimSpace(text)))
	if !status.Valid() {
		return "", fmt.Errorf("%w: %q", ErrStatusInvalid, text)
	}
	return status, nil
}

func (p Status) String() string {
	return string(p)
}

func (p Status) Valid() bool {
	switch p {
	case Open, Done:
		return true
	default:
		return false
	}
}

func (p Status) Definition() (Definition, bool) {
	for _, definition := range definitions {
		if definition.Status == p {
			return definition, true
		}
	}
	return Definition{}, false
}

func Statuses() []Status {
	ret := make([]Status, 0, len(definitions))
	for _, definition := range definitions {
		ret = append(ret, definition.Status)
	}
	return ret
}

func Definitions() []Definition {
	return slices.Clone(definitions)
}

func (p Status) MarshalText() ([]byte, error) {
	if !p.Valid() {
		return nil, fmt.Errorf("%w: %q", ErrStatusInvalid, p)
	}
	return []byte(p), nil
}

func (p *Status) UnmarshalText(text []byte) error {
	if p == nil {
		return ErrStatusInvalid
	}
	status, err := ParseStatus(string(text))
	if err != nil {
		return err
	}
	*p = status
	return nil
}
