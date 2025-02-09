package json2image

import (
	"encoding/json"
	"strconv"
	"strings"
)

func JsonCrop(input interface{}, rules []string) ([]byte, error) {
	output := make(map[string]interface{})
	for _, rule := range rules {
		steps := parseRule(rule)
		processStep(input, steps, output, nil)
	}

	return json.Marshal(output)
}

type PathStep struct {
	Key     string
	Indices []int // 将单个 Index 改为 Indices 数组
}

func parseRule(rule string) []PathStep {
	parts := strings.Split(rule, ".")
	var steps []PathStep

	for _, part := range parts {
		if idxStart := strings.Index(part, "["); idxStart != -1 && strings.HasSuffix(part, "]") {
			key := part[:idxStart]
			indexStr := part[idxStart+1 : len(part)-1]

			if indexStr == "*" {
				steps = append(steps, PathStep{Key: key, Indices: nil})
			} else {
				// 处理多个索引值，用逗号分隔
				indices := []int{}
				for _, idx := range strings.Split(indexStr, ",") {
					if index, err := strconv.Atoi(strings.TrimSpace(idx)); err == nil {
						indices = append(indices, index)
					}
				}
				steps = append(steps, PathStep{Key: key, Indices: indices})
			}
		} else {
			steps = append(steps, PathStep{Key: part})
		}
	}
	return steps
}

func processStep(input interface{}, steps []PathStep, output map[string]interface{}, pathSoFar []PathStep) {
	if len(steps) == 0 {
		return
	}

	currentStep := steps[0]
	remainingSteps := steps[1:]

	switch currentStep.Key {
	case "*":
		switch input := input.(type) {
		case map[string]interface{}:
			for key := range input {
				newStep := PathStep{Key: key}
				newPath := append(pathSoFar, newStep)
				processStep(input[key], remainingSteps, output, newPath)
			}
		case []interface{}:
			for i := range input {
				newStep := PathStep{Key: currentStep.Key, Indices: []int{i}} // 修复：使用切片包含单个索引
				newPath := append(pathSoFar, newStep)
				processStep(input[i], remainingSteps, output, newPath)
			}
		}
	default:
		switch input := input.(type) {
		case map[string]interface{}:
			if child, exists := input[currentStep.Key]; exists {
				newPath := append(pathSoFar, currentStep)
				if len(currentStep.Indices) > 0 {
					if slice, ok := child.([]interface{}); ok {
						// 处理多个索引
						for _, index := range currentStep.Indices {
							if index < len(slice) {
								newStep := PathStep{Key: currentStep.Key, Indices: []int{index}}
								newPath := append(pathSoFar, newStep)
								processStep(slice[index], remainingSteps, output, newPath)
							}
						}
					}
				} else {
					if slice, ok := child.([]interface{}); ok {
						// 处理 [*] 的情况
						for i := range slice {
							newStep := PathStep{Key: currentStep.Key, Indices: []int{i}}
							newPath := append(pathSoFar, newStep)
							processStep(slice[i], remainingSteps, output, newPath)
						}
					} else {
						processStep(child, remainingSteps, output, newPath)
					}
				}
			}
		case []interface{}:
			if len(currentStep.Indices) > 0 {
				for _, index := range currentStep.Indices {
					if index < len(input) {
						newStep := PathStep{Key: currentStep.Key, Indices: []int{index}}
						newPath := append(pathSoFar, newStep)
						processStep(input[index], remainingSteps, output, newPath)
					}
				}
			}
		}
	}

	if len(remainingSteps) == 0 {
		var value interface{}
		switch input := input.(type) {
		case map[string]interface{}:
			if len(currentStep.Indices) > 0 {
				if slice, ok := input[currentStep.Key].([]interface{}); ok {
					for _, index := range currentStep.Indices {
						if index < len(slice) {
							value = slice[index]
							fullPath := append(pathSoFar, PathStep{Key: currentStep.Key, Indices: []int{index}})
							setValue(output, fullPath, value)
						}
					}
					return
				}
			} else {
				value = input[currentStep.Key]
			}
		case []interface{}:
			if len(currentStep.Indices) > 0 {
				for _, index := range currentStep.Indices {
					if index < len(input) {
						value = input[index]
						fullPath := append(pathSoFar, PathStep{Key: currentStep.Key, Indices: []int{index}})
						setValue(output, fullPath, value)
					}
				}
				return
			}
		}

		if value != nil {
			fullPath := append(pathSoFar, currentStep)
			setValue(output, fullPath, value)
		}
	}
}

func setValue(output map[string]interface{}, path []PathStep, value interface{}) {
	current := output
	for i, step := range path {
		isLast := i == len(path)-1

		if len(step.Indices) > 0 {
			key := step.Key
			index := step.Indices[0] // 使用第一个索引

			var slice []interface{}
			if existing, ok := current[key].([]interface{}); ok {
				slice = existing
				if index >= len(slice) {
					newSlice := make([]interface{}, index+1)
					copy(newSlice, slice)
					for j := len(slice); j <= index; j++ {
						newSlice[j] = make(map[string]interface{})
					}
					slice = newSlice
					current[key] = slice
				}
			} else {
				slice = make([]interface{}, index+1)
				for j := 0; j <= index; j++ {
					slice[j] = make(map[string]interface{})
				}
				current[key] = slice
			}

			if isLast {
				slice[index] = value
				return
			} else {
				if elem, ok := slice[index].(map[string]interface{}); ok {
					current = elem
				} else {
					elem := make(map[string]interface{})
					slice[index] = elem
					current = elem
				}
			}
		} else {
			key := step.Key
			if isLast {
				current[key] = value
			} else {
				if existing, ok := current[key].(map[string]interface{}); ok {
					current = existing
				} else {
					newMap := make(map[string]interface{})
					current[key] = newMap
					current = newMap
				}
			}
		}
	}
}
