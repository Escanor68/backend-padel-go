package utils

import (
	"errors"
	"reflect"
)

// ValidatePointer verifica que un puntero no sea nil
func ValidatePointer(ptr interface{}, fieldName string) error {
	if ptr == nil {
		return errors.New(fieldName + " cannot be nil")
	}

	// Usar reflexión para verificar si es un puntero
	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Ptr {
		return errors.New(fieldName + " must be a pointer")
	}

	// Verificar si el puntero apunta a nil
	if val.IsNil() {
		return errors.New(fieldName + " cannot be nil")
	}

	return nil
}

// ValidateStringPointer verifica que un puntero a string no sea nil y no esté vacío
func ValidateStringPointer(ptr *string, fieldName string) error {
	if err := ValidatePointer(ptr, fieldName); err != nil {
		return err
	}

	if *ptr == "" {
		return errors.New(fieldName + " cannot be empty")
	}

	return nil
}

// ValidateNumberPointer verifica que un puntero a número no sea nil y sea válido
func ValidateNumberPointer(ptr *float64, fieldName string) error {
	if err := ValidatePointer(ptr, fieldName); err != nil {
		return err
	}

	if *ptr < 0 {
		return errors.New(fieldName + " must be positive")
	}

	return nil
}

// ValidateBoolPointer verifica que un puntero a bool no sea nil
func ValidateBoolPointer(ptr *bool, fieldName string) error {
	return ValidatePointer(ptr, fieldName)
}

// ValidateIntPointer verifica que un puntero a int no sea nil y sea válido
func ValidateIntPointer(ptr *int, fieldName string) error {
	if err := ValidatePointer(ptr, fieldName); err != nil {
		return err
	}

	if *ptr < 0 {
		return errors.New(fieldName + " must be positive")
	}

	return nil
}
