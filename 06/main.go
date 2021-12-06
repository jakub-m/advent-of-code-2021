package main

import (
	"advent"
)

const (
	renewAge = 6
	spawnAge = 8
)

func Calc(csvInput string, iterations int) (int, error) {
	initialAges, err := advent.ReadIntsFromStringSep(csvInput, ",")
	if err != nil {
		return 0, err
	}

	ageCountMap := make(map[int]int)
	for _, age := range initialAges {
		ageCountMap[age]++
	}

	for i := 0; i < iterations; i++ {

		updatedAgeCountMap := make(map[int]int)
		for age, count := range ageCountMap {
			advent.Assertf(age >= 0 && age <= 8, "invalid age: %d", age)
			advent.Assertf(count >= 0, "invalid count: %d", count)

			if age == 0 {
				updatedAgeCountMap[renewAge] += count
				updatedAgeCountMap[spawnAge] += count
			} else {
				updatedAgeCountMap[age-1] += count
			}

		}

		ageCountMap = updatedAgeCountMap
	}

	totalCount := 0
	for _, count := range ageCountMap {
		totalCount += count
	}
	return totalCount, nil
}
