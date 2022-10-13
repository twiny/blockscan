package api

import (
	"fmt"
	"strconv"
	"strings"
)

// parseRange
func parseRange(s string) (start int64, end int64, err error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		err = fmt.Errorf("range must be in format start:end")
		return
	}

	start, err = strconv.ParseInt(parts[0], 10, 0)
	if err != nil {
		return
	}

	if start < 0 {
		start = start * -1
	}

	end, err = strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		return
	}

	if end < 0 {
		end = end * -1
	}

	return
}
