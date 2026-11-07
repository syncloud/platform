package btrfs

import (
	"strconv"
	"strings"
)

type Counter string

func (c *Counter) UnmarshalJSON(data []byte) error {
	*c = Counter(strings.Trim(string(data), `"`))
	return nil
}

func (c Counter) Int64() int64 {
	value, err := strconv.ParseInt(string(c), 10, 64)
	if err != nil {
		return 0
	}
	return value
}
