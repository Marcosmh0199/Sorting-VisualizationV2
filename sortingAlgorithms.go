/*
Notas: 
1. Para una visualización correcta del gráfico de barras, por favor ejecutar el programa
en una terminal a pantalla completa.
2. Si se quiere cambiar el tamaño del slice, modificar la constante SLICE_SIZE.
*/

package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/lxn/win"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

//Constants
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

func main() {
	fmt.Print("Indique el tamaño del slice: ")
	var first int
	fmt.Scanln(&first)
	slice := randomSlice(MAX_PRIME, first)
	//barChartDriver(slice)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	initBsChart(slice)
	initQsChart(slice)
	ui.Render(&bsChart)
	ui.Render(&qsChart)
	bsChartDrawer(slice)
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

func quickSort(slice []float64, pair chan []int) []float64 {
	if len(slice) < 2 {
		return slice
	}
	low := 0
	high := len(slice)-1

	pivot := rand.Int() % len(slice)

	swap(&slice[pivot], &slice[high])
	for i := range slice {
		if slice[i] < slice[high] {
			swap(&slice[i], &slice[low])
			low++
		}
	}
	swap(&slice[low], &slice[high])
	quickSort(slice[:low])
	quickSort(slice[low+1:])
	return slice
}

func qsChartDrawer(slice []float64){
	qsChart.Data = make([]float64, len(slice))
	copy(qsChart.Data, slice)
	pairsChannel := make(chan []int)
	go quickSort(qsChart.Data, pairsChannel)
	for pair := range pairsChannel{
		swap(&qsChart.Data[pair[0]], &qsChart.Data[pair[1]])
		ui.Render(&qsChart)
		time.Sleep(MILI_SECONDS * time.Millisecond)
	}
}

func bsChartDrawer(slice []float64)  {
	bsChart.Data = make([]float64, len(slice))
	copy(bsChart.Data, slice)
	pairsChannel := make(chan []int)
	go bubbleSort(bsChart.Data, pairsChannel)
	for pair := range pairsChannel{
		swap(&bsChart.Data[pair[0]], &bsChart.Data[pair[1]])
		ui.Render(&bsChart)
		time.Sleep(MILI_SECONDS * time.Millisecond)
	}
}

func barChartDriver(slice []float64) {
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "1":
				bsChart.Title = "QuickSort"
				bsChart.Data = make([]float64, len(slice))
				copy(bsChart.Data, slice)
				quickSort(bsChart.Data)
			case "2":
				bsChartDrawer(slice)
			case "3":
				bsChart.Title = "Sort"
				bsChart.Data = make([]float64, len(slice))
				copy(bsChart.Data, slice)
				ui.Render(&bsChart)
			case "4":
				slice = randomSlice(generateSeed(), 5)
				bsChart.Data = make([]float64, len(slice))
				copy(bsChart.Data, slice)
				ui.Render(&bsChart)
			case "5":
				return
			}
		}
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
	bsChart.BarWidth = 3
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
	qsChart.SetRect(0, height-2, width, height*2 - 2)
	qsChart.BarWidth = 3
	qsChart.BarGap = 0
	qsChart.Labels = generateLabels(slice)
	qsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	qsChart.BarColors = []ui.Color{ui.ColorRed}
	qsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
}
