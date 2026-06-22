package utils

import (
	"errors"
	"strconv"
	"strings"
)

type ArrayList interface {
	int | int32 | int64 | uint | string | float32 | float64
}

func Map[T any, R any](collection []T, iteratee func(T) R) []R {
	result := make([]R, len(collection))
	for i, item := range collection {
		result[i] = iteratee(item)
	}
	return result
}

func Filter[T any](slice []T, f func(T) bool) []T {
	var filteredSlice []T
	for _, value := range slice {
		conditionFullfilled := f(value)
		if conditionFullfilled {
			filteredSlice = append(filteredSlice, value)
		}
	}
	return filteredSlice
}

func ArrayDifference[T ArrayList](array1, array2 []T) []T {
	array2Map := make(map[T]struct{})
	for _, elem := range array2 {
		array2Map[elem] = struct{}{}
	}

	var difference []T
	for _, elem := range array1 {
		if _, ok := array2Map[elem]; !ok {
			difference = append(difference, elem)
		}
	}
	return difference
}

func RemoveArrayDuplicates[T ArrayList](array []T) []T {
	var result []T
	hashMap := make(map[T]int)
	for _, value := range array {
		if _, ok := hashMap[value]; !ok {
			hashMap[value] = 1
		} else {
			hashMap[value]++
		}
	}
	for key := range hashMap {
		result = append(result, key)
	}
	return result
}

func StringToBoolean(s string) (bool, error) {
	lowerInput := strings.ToLower(s)
	if lowerInput == "true" {
		return true, nil
	}
	if lowerInput == "false" {
		return false, nil
	}
	val, err := strconv.Atoi(lowerInput)
	if err == nil {
		if val < 0 {
			return false, errors.New("invalid input: negative integers not supported")
		}

		if val > 0 {
			return true, nil
		}
		return false, nil
	}
	return false, errors.New("invalid boolean representation: " + lowerInput)
}
