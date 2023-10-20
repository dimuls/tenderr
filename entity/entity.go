package entity

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

type Class struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Name  string    `json:"name" db:"name"`
	Rules []string  `json:"rules" db:"rules"`
}

func (c *Class) CompileRules() ([]*regexp.Regexp, error) {
	regexps := make([]*regexp.Regexp, 0, len(c.Rules))
	for i, rule := range c.Rules {
		rx, err := regexp.Compile(rule)
		if err != nil {
			return nil, fmt.Errorf("invalid rule %d: %w", i, err)
		}

		regexps = append(regexps, rx)
	}
	return regexps, nil
}
