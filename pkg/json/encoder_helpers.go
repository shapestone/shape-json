package json

import (
	"errors"
	"strconv"
	"time"
)

// errNeedReflect is a sentinel error returned by appendInterface when
// the concrete type is not handled by the fast type-switch path and
// the caller should fall back to reflect-based encoding.
var errNeedReflect = errors.New("need reflect")

// appendInterface encodes a Go interface value to JSON using a type-switch
// over common concrete types. This avoids reflect entirely for the most
// frequent types seen in JSON data (primitives, []interface{},
// map[string]interface{}, time.Time, time.Duration, and Marshaler).
//
// Returns errNeedReflect for types not covered by the switch so the
// caller can fall back to the compiled encoder cache.
func appendInterface(buf []byte, v interface{}) ([]byte, error) {
	switch val := v.(type) {
	case nil:
		return append(buf, "null"...), nil
	case bool:
		if val {
			return append(buf, "true"...), nil
		}
		return append(buf, "false"...), nil
	case string:
		buf = append(buf, '"')
		buf = appendEscapedString(buf, val)
		buf = append(buf, '"')
		return buf, nil
	case int:
		return strconv.AppendInt(buf, int64(val), 10), nil
	case int8:
		return strconv.AppendInt(buf, int64(val), 10), nil
	case int16:
		return strconv.AppendInt(buf, int64(val), 10), nil
	case int32:
		return strconv.AppendInt(buf, int64(val), 10), nil
	case int64:
		return strconv.AppendInt(buf, val, 10), nil
	case uint:
		return strconv.AppendUint(buf, uint64(val), 10), nil
	case uint8:
		return strconv.AppendUint(buf, uint64(val), 10), nil
	case uint16:
		return strconv.AppendUint(buf, uint64(val), 10), nil
	case uint32:
		return strconv.AppendUint(buf, uint64(val), 10), nil
	case uint64:
		return strconv.AppendUint(buf, val, 10), nil
	case float32:
		return strconv.AppendFloat(buf, float64(val), 'g', -1, 32), nil
	case float64:
		return strconv.AppendFloat(buf, float64(val), 'g', -1, 64), nil
	case time.Time:
		buf = append(buf, '"')
		buf = val.AppendFormat(buf, time.RFC3339Nano)
		buf = append(buf, '"')
		return buf, nil
	case time.Duration:
		buf = append(buf, '"')
		buf = appendISO8601Duration(buf, val)
		buf = append(buf, '"')
		return buf, nil
	case []interface{}:
		buf = append(buf, '[')
		for i, elem := range val {
			if i > 0 {
				buf = append(buf, ',')
			}
			var err error
			buf, err = appendInterface(buf, elem)
			if err != nil {
				return buf, err
			}
		}
		buf = append(buf, ']')
		return buf, nil
	case map[string]interface{}:
		buf = append(buf, '{')
		// For consistent output, sort the keys.
		// Use a simple insertion sort for small maps (common case).
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sortStrings(keys)
		for i, k := range keys {
			if i > 0 {
				buf = append(buf, ',')
			}
			buf = append(buf, '"')
			buf = appendEscapedString(buf, k)
			buf = append(buf, '"', ':')
			var err error
			buf, err = appendInterface(buf, val[k])
			if err != nil {
				return buf, err
			}
		}
		buf = append(buf, '}')
		return buf, nil
	case Marshaler:
		b, err := val.MarshalJSON()
		if err != nil {
			return buf, err
		}
		return append(buf, b...), nil
	default:
		return buf, errNeedReflect
	}
}

// sortStrings sorts a string slice in-place using insertion sort.
// For the small key counts typical in JSON maps (< 20 keys) this is
// faster than sort.Strings because it avoids the interface overhead.
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		key := s[i]
		j := i - 1
		for j >= 0 && s[j] > key {
			s[j+1] = s[j]
			j--
		}
		s[j+1] = key
	}
}

// appendISO8601Duration formats a time.Duration as an ISO 8601 duration
// string (e.g. "PT1H30M5.5S") and appends it to buf without allocating.
func appendISO8601Duration(buf []byte, d time.Duration) []byte {
	if d == 0 {
		return append(buf, "PT0S"...)
	}

	buf = append(buf, 'P')

	if d < 0 {
		// ISO 8601 doesn't have negative durations in the standard,
		// but we prefix with '-' for round-trip compatibility.
		buf = append(buf, '-')
		d = -d
	}

	hours := int64(d / time.Hour)
	d -= time.Duration(hours) * time.Hour
	minutes := int64(d / time.Minute)
	d -= time.Duration(minutes) * time.Minute

	// Always emit T prefix before time components
	buf = append(buf, 'T')

	if hours > 0 {
		buf = strconv.AppendInt(buf, hours, 10)
		buf = append(buf, 'H')
	}
	if minutes > 0 {
		buf = strconv.AppendInt(buf, minutes, 10)
		buf = append(buf, 'M')
	}

	// Handle seconds + fractional
	secs := d.Seconds()
	if secs > 0 || (hours == 0 && minutes == 0) {
		// Use integer path when there's no fractional part
		if d%time.Second == 0 {
			buf = strconv.AppendInt(buf, int64(d/time.Second), 10)
		} else {
			buf = strconv.AppendFloat(buf, secs, 'f', -1, 64)
		}
		buf = append(buf, 'S')
	}

	return buf
}
