package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrLen    = errors.New("length does not match")
	ErrRegexp = errors.New("value does not match regexp")
	ErrInt    = errors.New("value not in allowed list")
	ErrMin    = errors.New("value less than minimum")
	ErrMax    = errors.New("value more than maximum")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, e := range v {
		sb.WriteString(fmt.Sprintf("Field: %s, Error:%s; ", e.Field, e.Err))
	}
	return sb.String()
}

type CheckList struct {
	Len    *int
	Regexp *regexp.Regexp
	In     []string
	Min    *int
	Max    *int
}

func GetCheckListFromStructTag(structTag string) (*CheckList, error) {
	res := CheckList{}
	for _, rule := range strings.Split(structTag, "|") {
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("unknown validation key: %s", parts[0])
		}
		key, val := parts[0], parts[1]
		switch key {
		case "len":
			n, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			res.Len = &n

		case "min":
			n, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			res.Min = &n

		case "max":
			n, err := strconv.Atoi(val)
			if err != nil {
				return nil, err
			}
			res.Max = &n

		case "in":
			res.In = strings.Split(val, ",")

		case "regexp":
			re, err := regexp.Compile(val)
			if err != nil {
				return nil, err
			}
			res.Regexp = re

		default:
			return nil, fmt.Errorf("unknown validation key: %s", key)
		}
	}
	return &res, nil
}

func Validate(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return errors.New("expected a struct")
	}

	rt := rv.Type()
	var allErrs ValidationErrors

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if !f.IsExported() {
			continue
		}
		tag := f.Tag.Get("validate")
		if tag == "" {
			continue
		}

		cl, err := GetCheckListFromStructTag(tag)
		if err != nil {
			return err
		}

		val := rv.Field(i)
		name := f.Name

		//nolint:exhaustive
		switch f.Type.Kind() {
		case reflect.String:
			allErrs = append(allErrs, validateString(name, val.String(), cl)...)

		case reflect.Int:
			allErrs = append(allErrs, validateInt(name, int(val.Int()), cl)...)

		case reflect.Slice:
			for j := 0; j < val.Len(); j++ {
				elem := val.Index(j)
				elemName := fmt.Sprintf("%s[%d]", name, j)
				//nolint:exhaustive
				switch elem.Kind() {
				case reflect.String:
					allErrs = append(allErrs, validateString(elemName, elem.String(), cl)...)
				case reflect.Int:
					allErrs = append(allErrs, validateInt(elemName, int(elem.Int()), cl)...)
				default:
					return fmt.Errorf("invalid slice elem kind: %s", elem.Kind())
				}
			}
		default:
			return fmt.Errorf("invalid kind: %s", f.Type.Kind())
		}
	}

	if len(allErrs) > 0 {
		return allErrs
	}
	return nil
}

func validateString(field, s string, cl *CheckList) ValidationErrors {
	var errs ValidationErrors
	if cl.Len != nil && len(s) != *cl.Len {
		errs = append(errs, ValidationError{Field: field, Err: ErrLen})
	}
	if cl.Regexp != nil && !cl.Regexp.MatchString(s) {
		errs = append(errs, ValidationError{Field: field, Err: ErrRegexp})
	}
	if len(cl.In) > 0 {
		ok := false
		for _, v := range cl.In {
			if s == v {
				ok = true
				break
			}
		}
		if !ok {
			errs = append(errs, ValidationError{Field: field, Err: ErrInt})
		}
	}
	return errs
}

func validateInt(field string, x int, cl *CheckList) ValidationErrors {
	var errs ValidationErrors
	if cl.Min != nil && x < *cl.Min {
		errs = append(errs, ValidationError{Field: field, Err: ErrMin})
	}
	if cl.Max != nil && x > *cl.Max {
		errs = append(errs, ValidationError{Field: field, Err: ErrMax})
	}
	if len(cl.In) > 0 {
		ok := false
		for _, v := range cl.In {
			if n, err := strconv.Atoi(v); err == nil && n == x {
				ok = true
				break
			}
		}
		if !ok {
			errs = append(errs, ValidationError{Field: field, Err: ErrInt})
		}
	}
	return errs
}
