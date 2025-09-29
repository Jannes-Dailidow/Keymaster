// Copyright (c) 2025 ToeiRei
// Keymaster - SSH key management system
// This source code is licensed under the MIT license found in the LICENSE file.

package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/bobg/go-generics/v4/slices"
)

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func toPascalCase(s string) string {
	words := strings.Split(capitalize(s), "_")

	return slices.Accum(words, func(result, word string) string {
		if len(word) > 0 {
			return result + capitalize(word)
		}
		return result
	})
}

func getFieldPointers(s interface{}, keys []string) ([]any, error) {
	v := reflect.ValueOf(s).Elem()

	return slices.Mapx(keys, func(_ int, key string) (any, error) {
		field := v.FieldByName(key)
		if !field.IsValid() {
			return nil, fmt.Errorf("Field %s not found", key)
		}

		if field.CanAddr() {
			return field.Addr().Interface(), nil
		} else {
			return nil, fmt.Errorf("Field %s is not addressable", key)
		}
	})
}

func QueryModels[T any](ctx context.Context, db *sql.DB, query string, args ...any) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var model_keys []string
	var records []T

	if column_names, err := rows.Columns(); err != nil {
		return nil, err
	} else {
		model_keys = slices.Map(column_names, toPascalCase)
	}

	for rows.Next() {
		var record T
		fields, err := getFieldPointers(&record, model_keys)
		if err != nil {
			return nil, err
		}

		if err = rows.Scan(fields...); err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func QueryModel[T any](ctx context.Context, db *sql.DB, query string, args ...any) (*T, error) {
	records, err := QueryModels[T](ctx, db, query, args...)
	if err != nil {
		return nil, err
	}
	if len(records) != 1 {
		return nil, fmt.Errorf("Expected 1 row, got %v", len(records))
	}
	return &records[0], nil
}
