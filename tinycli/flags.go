package tinycli

import (
	"fmt"
	"strconv"
	"time"
)

// GetTimeDuration get time from string
func GetTimeDuration(v string) (time.Duration, error) {
	if v == "" || v == "0" {
		return 0, nil
	}
	var ret time.Duration
	hasSet := false
	if len(v) > 2 {
		hasSet = true
		vv, err := strconv.ParseInt(v[:len(v)-2], 10, 64)
		switch v[len(v)-2:] {
		case "ns":
			ret = time.Duration(vv)
		case "us":
			ret = time.Duration(vv) * time.Microsecond
		case "ms":
			ret = time.Duration(vv) * time.Millisecond
		default:
			hasSet = false
		}
		if hasSet && err != nil {
			return 0, fmt.Errorf("bad forma of time value: %s", v)
		}
	}
	if !hasSet {
		hasSet = true
		vv, err := strconv.ParseInt(v[:len(v)-1], 10, 64)
		switch v[len(v)-1:] {
		case "s":
			ret = time.Duration(vv) * time.Second
		case "m":
			ret = time.Duration(vv) * time.Minute
		case "h":
			ret = time.Duration(vv) * time.Hour
		case "d":
			ret = time.Duration(vv) * time.Hour * 24
		default:
			return 0, fmt.Errorf("bad forma of time value: %s", v)
		}
		if hasSet && err != nil {
			return 0, fmt.Errorf("bad forma of time value: %s", v)
		}
	}

	return ret, nil
}

func setFlagInt(envExists bool, ev string, val *int, defaultVal int, usage string) (err error) {
	*val = defaultVal
	if envExists {
		if *val, err = strconv.Atoi(ev); err != nil {
			return err
		}
	}
	return nil
}

func setFlagString(envExists bool, ev string, val *string, defaultVal string, usage string) error {
	*val = defaultVal
	if envExists {
		*val = ev
	}
	return nil
}

func setFlagDuration(envExists bool, ev string, val *time.Duration, defaultVal string,
	usage string) (err error) {
	if *val, err = GetTimeDuration(defaultVal); err != nil {
		return err
	}
	if envExists {
		if *val, err = GetTimeDuration(ev); err != nil {
			return err
		}
	}
	return nil
}
