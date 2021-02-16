package utils

import (
	"github.com/bwmarrin/discordgo"
)

// QuickSort -
func QuickSort(arr []*discordgo.Channel) []*discordgo.Channel {
	newArr := make([]*discordgo.Channel, len(arr))

	for i, v := range arr {
		newArr[i] = v
	}

	sort(newArr, 0, len(arr)-1)

	return newArr
}

// Sort for quicksort
func sort(arr []*discordgo.Channel, start, end int) {
	if (end - start) < 1 {
		return
	}

	pivot := arr[end]
	splitIndex := start

	for i := start; i < end; i++ {
		if arr[i].Position < pivot.Position {
			temp := arr[splitIndex]

			arr[splitIndex] = arr[i]
			arr[i] = temp

			splitIndex++
		}
	}

	arr[end] = arr[splitIndex]
	arr[splitIndex] = pivot

	sort(arr, start, splitIndex-1)
	sort(arr, splitIndex+1, end)
}

// StringInSlice - Check if a slice contains a string
func StringInSlice(str string, sl []string) bool {
	for _, s := range sl {
		if s == str {
			return true
		}
	}
	return false
}
