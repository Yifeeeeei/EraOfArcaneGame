package game

import (
	"fmt"
)

func requiredBoardPosition(data map[string]any, colKey string, rowKey string) (Position, error) {
	col, err := requiredBoardCoordinate(data, colKey)
	if err != nil {
		return Position{}, err
	}
	row, err := requiredBoardCoordinate(data, rowKey)
	if err != nil {
		return Position{}, err
	}
	pos := Position{Col: col, Row: row}
	if !pos.Valid() {
		return Position{}, fmt.Errorf("invalid position")
	}
	return pos, nil
}

func optionalBoardPosition(data map[string]any, colKey string, rowKey string) (Position, bool, error) {
	_, hasCol := data[colKey]
	_, hasRow := data[rowKey]
	if !hasCol && !hasRow {
		return Position{}, false, nil
	}
	if hasCol != hasRow {
		return Position{}, false, fmt.Errorf("%s and %s must be provided together", colKey, rowKey)
	}
	pos, err := requiredBoardPosition(data, colKey, rowKey)
	if err != nil {
		return Position{}, false, err
	}
	return pos, true, nil
}

func requiredBoardCoordinate(data map[string]any, key string) (int, error) {
	value, ok := data[key]
	if !ok {
		return 0, fmt.Errorf("missing %s", key)
	}
	switch typed := value.(type) {
	case float64:
		coord := int(typed)
		if typed != float64(coord) {
			return 0, fmt.Errorf("%s must be an integer", key)
		}
		if coord < 0 || coord > 2 {
			return 0, fmt.Errorf("%s must be between 0 and 2", key)
		}
		return coord, nil
	case int:
		if typed < 0 || typed > 2 {
			return 0, fmt.Errorf("%s must be between 0 and 2", key)
		}
		return typed, nil
	default:
		return 0, fmt.Errorf("%s must be a number", key)
	}
}
