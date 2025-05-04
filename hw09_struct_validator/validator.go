package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {

	// TODO: мб стрингс билдер
	var res string
	for _, ve := range v {
		res += fmt.Sprintf("Field: %s, Error:%s", ve.Field, ve.Err)
	}
	panic("implement me")
}

type CheckList struct {
	Len    *int
	Regexp *string
	In     []string
	Min    *int
	Max    *int
}

const NotSet = 0

func GetCheckListFromStructTag(structTag string) (*CheckList, error) {

	var res CheckList
	rules := strings.Split(structTag, "|")
	for _, rule := range rules {
		keyVal := strings.SplitN(rule, ":", 2)
		if len(keyVal) != 2 {

			// TODO:
			return nil, fmt.Errorf("not enough args: %v", keyVal)
		}
		switch keyVal[0] {
		case "len":
			num, err := strconv.Atoi(keyVal[1])
			if err != nil {
				return nil, err
			}
			res.Len = &num
		case "min":
			num, err := strconv.Atoi(keyVal[1])
			if err != nil {
				return nil, err
			}
			res.Min = &num
		case "max":
			num, err := strconv.Atoi(keyVal[1])
			if err != nil {
				return nil, err
			}
			res.Max = &num
		case "in":
			res.In = strings.Split(keyVal[1], ",")

		case "regexp":
			// TODO:
			// ParseRegexp(keyVal[1])
			res.Regexp = &keyVal[1]
		}
	}

	return &res, nil
}

// TODO:
func ParseRegexp(str string) {

}

// Подсказки:
// reflect.StructTag
// regexp.Compile

// Функиця валидирует ПУБЛИЧНЫЕ поля на основе структурного тега validate.
// Функция может возвращать или програмную ошибку, или ValidationErrors произошедшую во время валидации.
func Validate(v any) error {

	//как проверить публичное поле или нет?

	var resCustomErrors ValidationErrors

	// 	check if not struct
	// 	Проверка что входной интерфейс структура
	structVal, err := GetStruct(v)
	if err != nil {
		return err
	}

	structType := structVal.Type()

	for i := 0; i < structVal.NumField(); i++ {
		// 	ф-я игнорирует поля без структурных тегов/тегов validate
		structFieldInfo := structType.Field(i)
		structFieldVal := structVal.Field(i)

		if !structFieldInfo.IsExported() {
			continue
		}
		structTag := structFieldInfo.Tag.Get("validate")
		if structTag == "" {
			continue
		}

		checkList, err := GetCheckListFromStructTag(structTag)
		if err != nil {
			// TODO: нужно понять как обрабатывать
			// ошибки невеных стракт тегов в поле in: итд
			// fmt.Errorf("not enough args: может вернуть например
			return err
		}

		// var cmpTmpErr ValidationErrors
		// скорее всего для каждой страктфилд одна ошибка и не нужно ебаться
		// err = ProcessStructFields(sf, structTag)
		err = checkList.ProcessStructFields(structFieldInfo, structFieldVal)
		if err != nil {
			var customValidateErrors ValidationErrors
			if errors.As(err, &customValidateErrors) {
				resCustomErrors = append(resCustomErrors, customValidateErrors...)
			} else {
				return err
			}
		}

	}

	// 	Parse struct tag:
	// 	min max regexp
	// 	in, len

	// TODO: ДОП: Ошибки валидации вынести в отдельные переменные и заврапить %w

	return resCustomErrors
}

// Обработка поля структуры и запуск валидности разные свич кейсы или типо того.
// Возвращает либо ошибку и программа завершается, либо результат проверки.
// int, []int;
// string, []string.

func (cl *CheckList) ProcessStructFields(info reflect.StructField, val reflect.Value) error {
	switch info.Type.Kind() {
	case reflect.Int:
		return cl.ProcessInt(val.Int())
	}
	return nil
}

func (cl *CheckList) ProcessInt(val int64) error {
	if cl.Min != nil && val < int64(*cl.Min){
		return ValidationErrors{ValidationError{
			Field: ,
		}}
	}
	return nil
}

func ProcessStructFields(sf reflect.StructField, tag string) error {

	switch sf.Type.Kind() {
	case reflect.Int:

	}
	return nil
}

func ProcessInts(structTag string) {

}
func GetStruct(v any) (reflect.Value, error) {
	value := reflect.ValueOf(v)

	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return value, errors.New("value is nil pointer")
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return value, errors.New("value is not a struct")
	}

	return value, nil

}

func CheckStructTags(field any) bool {
	// a := reflect.StructOf()
	return false
}

func ValidateLen(value reflect.Value) error {
	return nil
}
