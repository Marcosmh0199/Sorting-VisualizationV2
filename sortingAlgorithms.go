/*
Notas:
1. Para una visualización correcta del gráfico de barras, por favor ejecutar el programa
en una terminal a pantalla completa.
2. Si se quiere cambiar el tamaño del slice, modificar la constante SLICE_SIZE.
https://www.geeksforgeeks.org/iterative-quick-sort/
*/

package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/lxn/win"
	"log"
	"math"
	"strconv"
	"sync"
)

//Constants
const (
	BAR_WIDTH = 3
)
const FONT_WIDTH = 8
const FONT_HEIGHT = 16
const MAX_NUMBER_SIZE = 32
const MAX_PRIME = 101
const MILI_SECONDS = 10

//global variables
var width int = int(win.GetSystemMetrics(win.SM_CXSCREEN) / FONT_WIDTH)
var height int = int(win.GetSystemMetrics(win.SM_CYSCREEN) / (FONT_HEIGHT*2))
var bsChart widgets.BarChart
var qsChart widgets.BarChart
var m sync.Mutex

func main() {
	barNumber := width / BAR_WIDTH - 1
	fmt.Print("Indique el tamaño del slice(Se recomienda que contenga " + strconv.Itoa(barNumber) +" elementos para una visualizacion correcta): ")
	var first int
	fmt.Scanln(&first)
	slice := randomSlice(MAX_PRIME, first)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	initBsChart(slice)
	initQsChart(slice)
	ui.Render(&bsChart)
	ui.Render(&qsChart)
	go bsChartDrawer(slice)
	go qsChartDrawer(slice)
	fmt.Scanln()
}

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
func randomSlice(seed int, size int) []float64{
	var slice = make([]float64, size)
	var m = int(math.Pow(2, MAX_NUMBER_SIZE))  // modulus
	var a = 22695477                    // multiplier
	var c = 1                           //increment
	for i := 0; i < size; i++ {
		seed = (a * seed + c) % m
		slice[i] = float64(seed % MAX_NUMBER_SIZE) //cast to make the slice values compatible with the barChart
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

/*
Aux function to remove an element from a slice
*/
func remove(slice [][]int, index int) [][]int {
	return append(slice[:index], slice[index+1:]...)
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

func partition(slice []float64, start int, end int, pair chan []int) int {
	pivot := slice[end]
	index := start
	for i := start; i < end; i++{
		if slice[i] <= pivot{
			pair <- []int{i,index}
			index++
		}
	}
	pair <- []int{index,end}
	return index
}

func quickSort(slice []float64, size int, pair chan []int) {
	var stack [][]int
	start := 0
	end := size
	stack = append(stack, []int{start,end})
	for len(stack) > 0{
		start, end = stack[0][0],stack[0][1]
		stack = remove(stack, 0)
		pivot := partition(slice, start, end, pair)

		if pivot-1 > start {
			stack = append(stack, []int{start,pivot-1})
		}
		if pivot+1 < end {
			stack = append(stack, []int{pivot+1,end})
		}
	}
	close(pair)
}

func qsChartDrawer(slice []float64){
	qsChart.Data = make([]float64, len(slice))
	copy(qsChart.Data, slice)
	pairsChannel := make(chan []int)
	go quickSort(qsChart.Data, len(slice)-1, pairsChannel)
	for pair := range pairsChannel{
		swap(&qsChart.Data[pair[0]], &qsChart.Data[pair[1]])
		m.Lock()
		ui.Render(&qsChart)
		m.Unlock()
	}
}

func bsChartDrawer(slice []float64)  {
	bsChart.Data = make([]float64, len(slice))
	copy(bsChart.Data, slice)
	pairsChannel := make(chan []int)
	go bubbleSort(bsChart.Data, pairsChannel)
	for pair := range pairsChannel{
		swap(&bsChart.Data[pair[0]], &bsChart.Data[pair[1]])
		m.Lock()
		ui.Render(&bsChart)
		m.Unlock()
	}
}

func generateLabels(slice []float64) []string {
	var labels = make([]string, len(slice))
	for i := range slice{
		labels[i] = strconv.Itoa(i)
	}
	return labels
}

func initBsChart(slice []float64)  {
	bsChart = *widgets.NewBarChart()
	bsChart.Data = slice
	bsChart.Title = "BubbleSort"
	bsChart.SetRect(0, 0, width, height - 2)
	bsChart.BarWidth = BAR_WIDTH
	bsChart.BarGap = 0
	bsChart.Labels = generateLabels(slice)
	bsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	bsChart.BorderBottom = false
	bsChart.BarColors = []ui.Color{ui.ColorRed}
	bsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
}

func initQsChart(slice []float64){
	qsChart = *widgets.NewBarChart()
	qsChart.Data = slice
	qsChart.Title = "QuickSort"
	qsChart.SetRect(0, height-2, width, height*2 - 3)
	qsChart.BarWidth = 3
	qsChart.BarGap = 0
	qsChart.Labels = generateLabels(slice)
	qsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	qsChart.BarColors = []ui.Color{ui.ColorRed}
	qsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
}