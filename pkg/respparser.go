package respparser

import (
	"fmt"
	"strconv"
	"strings"
)

const crlfTerminator = "\r\n"

func readSimpleString(idx *int, str *string) (string, error) {
	out := ""
	lim := strings.Index((*str)[*idx:], crlfTerminator)
	if lim == -1 {
		return "", fmt.Errorf("error of terminator wasn't found")
	}
	lim = lim + (*idx)
	out = (*str)[*idx:lim]
	*idx = lim + 1
	return out, nil
}

func readErrorMSG(idx *int, str *string) (error, error) {
	out := ""
	lim := strings.Index((*str)[*idx:], crlfTerminator)
	if lim == -1 {
		return nil, fmt.Errorf("error of terminator wasn't found")
	}
	lim = lim + (*idx)
	out = (*str)[*idx:lim]
	*idx = lim + 1
	return fmt.Errorf(out), nil
}

func readInteger(idx *int, str *string) (int64, error) {
	out := ""
	lim := strings.Index((*str)[*idx:], crlfTerminator)
	if lim == -1 {
		return -1, fmt.Errorf("error of terminator wasn't found")
	}
	lim = lim + (*idx)
	out = (*str)[*idx:lim]
	*idx = lim + 1
	integer, err := strconv.ParseInt(out, 10, 64)
	if err != nil {
		return -1, err
	}
	return integer, nil
}

func readBulkString(idx *int, str *string) (interface{}, error) {
	out := ""
	length, err := readInteger(idx, str)
	if err != nil {
		return "", err
	}
	if length == -1 {
		return nil, nil
	}
	(*idx)++
	lim := (*idx) + int(length)
	if lim > len(*str) {
		return "", fmt.Errorf("error of string lenght isn't enough")
	}
	out = (*str)[*idx:lim]
	if out == "-1" {
		return "", nil
	}
	*idx = lim + 1
	return out, nil
}

func readBoolean(idx *int, str *string) (bool, error) {
	if (*idx) >= len(*str) {
		return false, fmt.Errorf("error of string lenght isn't enough")
	}
	if (*str)[(*idx)] == 't' {
		*idx += 2
		return true, nil
	}
	*idx += 2
	return false, nil
}

func readArray(idx *int, str *string) ([]interface{}, error) {
	length, err := readInteger(idx, str)
	if err != nil {
		return nil, err
	}

	if length == -1 {
		return nil, nil
	}

	var out []interface{}
	(*idx)++
	for i := 0; i < int(length); i++ {
		var curDatatype rune = rune((*str)[*idx])
		switch curDatatype {
		case rune('+'):
			(*idx)++
			simpleStr, err := readSimpleString(idx, str)
			if err != nil {
				return nil, err
			}
			out = append(out, simpleStr)

		case rune('-'):
			(*idx)++
			msg, err := readErrorMSG(idx, str)
			if err != nil {
				return nil, err
			}
			out = append(out, msg)

		case rune(':'):
			(*idx)++
			integer, err := readInteger(idx, str)
			if err != nil {
				return nil, err
			}
			out = append(out, integer)

		case rune('$'):
			(*idx)++
			bulkString, err := readBulkString(idx, str)
			if err != nil {
				return nil, err
			}
			out = append(out, bulkString)

		case rune('*'):
			(*idx)++
			array, err := readArray(idx, str)
			if err != nil {
				return nil, err
			}
			out = append(out, array)

		case rune('_'):
			(*idx)++
			out = append(out, nil)
			(*idx)++

		case rune('#'):
			(*idx)++
			boolean, err := readBoolean(idx, str)
			if err != nil {
				return nil, err
			}
			out = append(out, boolean)

		case rune(','):
			(*idx)++
			float, err := readFloat(idx, str)
			if err != nil {
				return nil, err
			}
			out = append(out, float)

		case rune('!'):
			(*idx)++
			bulkErrorMSG, err := readBulkError(idx, str)
			if err != nil {
				return nil, err
			}
			out = append(out, bulkErrorMSG)

		case rune('%'):
			(*idx)++
			respMap, err := readMap(idx, str)
			if err != nil {
				return nil, err
			}
			out = append(out, respMap)
		default:
			return nil, fmt.Errorf("error of unsupported datatype")
		}
		(*idx)++
	}
	return out, nil
}

func readFloat(idx *int, str *string) (float64, error) {
	out := ""
	lim := strings.Index((*str)[*idx:], crlfTerminator)
	if lim == -1 {
		return -1, fmt.Errorf("error of terminator wasn't found")
	}
	lim = lim + (*idx)
	out = (*str)[*idx:lim]
	*idx = lim + 1
	integer, err := strconv.ParseFloat(out, 64)
	if err != nil {
		return -1, err
	}
	return integer, nil
}

func readBulkError(idx *int, str *string) (error, error) {
	out := ""
	length, err := readInteger(idx, str)
	if err != nil {
		return nil, err
	}
	(*idx)++
	lim := (*idx) + int(length)
	if lim > len(*str) {
		return nil, fmt.Errorf("error of string lenght isn't enough")
	}
	out = (*str)[*idx:lim]
	if out == "-1" {
		return nil, nil
	}
	*idx = lim + 1

	return fmt.Errorf(out), nil
}

func readMap(idx *int, str *string) (map[interface{}]interface{}, error) {

	out := make(map[interface{}]interface{})
	length, err := readInteger(idx, str)
	if err != nil {
		return nil, err
	}

	if length == -1 {
		return nil, nil
	}

	(*idx)++
	var prev interface{}
	for i := 0; i < int(length)*2; i++ {
		var curDatatype rune = rune((*str)[*idx])

		var val interface{}
		switch curDatatype {
		case rune('+'):
			(*idx)++
			simpleStr, err := readSimpleString(idx, str)
			if err != nil {
				return nil, err
			}
			val = simpleStr

		case rune('-'):
			(*idx)++
			msg, err := readErrorMSG(idx, str)
			if err != nil {
				return nil, err
			}
			val = msg

		case rune(':'):
			(*idx)++
			integer, err := readInteger(idx, str)
			if err != nil {
				return nil, err
			}
			val = integer

		case rune('$'):
			(*idx)++
			bulkString, err := readBulkString(idx, str)
			if err != nil {
				return nil, err
			}
			val = bulkString

		case rune('*'):
			(*idx)++
			array, err := readArray(idx, str)
			if err != nil {
				return nil, err
			}
			val = array

		case rune('_'):
			(*idx)++
			val = nil
			(*idx)++

		case rune('#'):
			(*idx)++
			boolean, err := readBoolean(idx, str)
			if err != nil {
				return nil, err
			}
			val = boolean

		case rune(','):
			(*idx)++
			float, err := readFloat(idx, str)
			if err != nil {
				return nil, err
			}
			val = float

		case rune('!'):
			(*idx)++
			bulkErrorMSG, err := readBulkError(idx, str)
			if err != nil {
				return nil, err
			}
			val = bulkErrorMSG

		case rune('%'):
			(*idx)++
			respMap, err := readMap(idx, str)
			if err != nil {
				return nil, err
			}
			val = respMap
		default:
			return nil, fmt.Errorf("error of unsupported datatype")
		}
		(*idx)++
		if i%2 == 1 {
			out[prev] = val
		}

		prev = val
	}
	return out, nil

}

// Parser reads an resp string, returning an array of interfaces of items
func Parser(str string) ([]interface{}, error) {
	var out []interface{}
	var idx int = 0
	for ; idx < len(str); idx++ {
		var curDatatype rune = rune(str[idx])
		switch curDatatype {
		case rune('+'):
			idx++
			simpleStr, err := readSimpleString(&idx, &str)
			if err != nil {
				return nil, err
			}
			out = append(out, simpleStr)

		case rune('-'):
			idx++
			msg, err := readErrorMSG(&idx, &str)
			if err != nil {
				return nil, err
			}
			out = append(out, msg)

		case rune(':'):
			idx++
			integer, err := readInteger(&idx, &str)
			if err != nil {
				return nil, err
			}
			out = append(out, integer)

		case rune('$'):
			idx++
			bulkString, err := readBulkString(&idx, &str)
			if err != nil {
				return nil, err
			}
			out = append(out, bulkString)

		case rune('*'):
			idx++
			array, err := readArray(&idx, &str)
			if err != nil {
				return nil, err
			}
			out = append(out, array)

		case rune('_'):
			idx++
			out = append(out, nil)
			idx++

		case rune('#'):
			idx++
			boolean, err := readBoolean(&idx, &str)
			if err != nil {
				return nil, err
			}
			out = append(out, boolean)

		case rune(','):
			idx++
			float, err := readFloat(&idx, &str)
			if err != nil {
				return nil, err
			}
			out = append(out, float)

		case rune('!'):
			idx++
			bulkErrorMSG, err := readBulkError(&idx, &str)
			if err != nil {
				return nil, err
			}
			out = append(out, bulkErrorMSG)

		case rune('%'):
			idx++
			respMap, err := readMap(&idx, &str)
			if err != nil {
				return nil, err
			}
			out = append(out, respMap)
		default:
			return nil, fmt.Errorf("error of unsupported datatype")
		}
	}

	return out, nil
}
