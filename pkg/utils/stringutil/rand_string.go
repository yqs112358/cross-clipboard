package stringutil

import (
	"math/rand"
)

const lowerCaseString string = "abcdefghijklmnopqrstuvwxyz"
const lowerCaseStringLen int = len(lowerCaseString)

const upperCaseString string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const upperCaseStringLen int = len(upperCaseString)

const letterString string = lowerCaseString + upperCaseString
const letterStringLen int = len(letterString)

const numberString string = "0123456789"
const numberStringLen int = len(numberString)

// RandStringNumber random string number by giving length
func RandStringNumber(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = numberString[rand.Int63()%int64(numberStringLen)]
	}
	return string(b)
}

// RandStringLetter random string upper letter by giving length
func RandStringLetter(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterString[rand.Int63()%int64(letterStringLen)]
	}
	return string(b)
}

// RandStringUpperLetter random string upper letter by giving length
func RandStringUpperLetter(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = upperCaseString[rand.Int63()%int64(upperCaseStringLen)]
	}
	return string(b)
}

// RandStringLowerLetter random string upper letter by giving length
func RandStringLowerLetter(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = lowerCaseString[rand.Int63()%int64(lowerCaseStringLen)]
	}
	return string(b)
}
