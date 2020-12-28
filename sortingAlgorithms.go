/*
Notas:

*/

package main

import (
	"math"
	"math/rand"
)

const (
	FONT_WIDTH = 8
	FONT_HEIGHT = 16
	MAX_NUMBER_SIZE = 32
	MAX_PRIME = 101
	MILI_SECONDS = 10
)

/*
Generates a prime number from the system hour
*/
func generateSeed() int{
	return 111
}

/*
creates a N size slice with random numbers based on the linear congruential method using only the seed as param
output: slice with N random integers
 */
func randomSlice(seed int, size int) []int{
	var slice = make([]int, size)
	var m = int(math.Pow(2, MAX_NUMBER_SIZE)) // modulus
	var a = 22695477                          // multiplier
	var c = 1                                 //increment
	for i := 0; i < size; i++ {
		seed = (a * seed + c) % m
		slice[i] = seed % MAX_NUMBER_SIZE //cast to make the slice values compatible with the barChart
	}
	return slice
}

/*
Aux function to swap two numbers
*/
func swap (a *float64, b *float64){
	temp := *a
	*a = *b
	*b = temp
}

func bubbleSort(slice []float64, pair chan []int){
	n := len(slice) - 1
	for true {
		swapped := false
		for i := 0; i < n; i++{
			if slice[i] > slice[i+1]{
				pair <- []int{i, i+1}
				swapped = true
			}
		}
		if !swapped{
			break
		}
		n--
	}
	close(pair)
}

func quickSort(slice []float64, pair chan []int){
	if len(slice) < 2 {
		return
	}
	low := 0
	high := len(slice)-1

	pivot := rand.Int() % len(slice)

	pair <- []int{pivot, high}
	for i := range slice {
		if slice[i] < slice[high] {
			pair <- []int{i, low}
			low++
		}
	}
	pair <- []int{low, high}
	quickSort(slice[:low], pair)
	quickSort(slice[low+1:], pair)
	//close(pair)
}

