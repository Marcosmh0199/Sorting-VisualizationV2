/*
Notas: 
1. Para una visualización correcta del gráfico de barras, por favor ejecutar el programa
en una terminal a pantalla completa.
2. Si se quiere cambiar el tamaño del slice, modificar la constante SLICE_SIZE.
*/

package main

import (
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
const SLICE_SIZE = 50
const MAX_PRIME = 101
const MILI_SECONDS = 10

//global variables
var width int = int(win.GetSystemMetrics(win.SM_CXSCREEN) / FONT_WIDTH)
var height int = int(win.GetSystemMetrics(win.SM_CYSCREEN) / (FONT_HEIGHT*2))
var bsChart widgets.BarChart
var qsChart widgets.BarChart

func main() {
	slice := randomSlice(MAX_PRIME)
	barChartDriver(slice)
}

func remove(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}

/*
the sieve of eratosthenes algorithm, used to generate prime number in certain range
*/
func generateSeed() int{
	return 0
}

/*
creates a N size slice with random numbers based on the linear congruential method using only the seed as param
output: slice with N random integers
 */
func randomSlice(seed int) []float64{
	var slice = make([]float64, SLICE_SIZE)
	var m = int(math.Pow(2, MAX_NUMBER_SIZE))  // modulus
	var a = 22695477                    // multiplier
	var c = 1                           //increment
	for i := 0; i < SLICE_SIZE; i++ {
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

func bubbleSort(slice []float64) []float64{
	n := len(slice) - 1
	for true {
		swapped := false
		for i := 0; i < n; i++{
			if slice[i] > slice[i+1]{
				swap(&slice[i], &slice[i+1])
				swapped = true
				updateChart() //update the barChart, not related to Bubblesort
			}
		}
		if !swapped{
			break
		}
		n--
	}
	return slice
}

func quickSort(slice []float64) []float64 {
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
			updateChart() //update the barChart, not related to Quicksort
			low++
		}
	}
	swap(&slice[low], &slice[high])
	quickSort(slice[:low])
	quickSort(slice[low+1:])
	return slice
}

func displayHelp(){
	for i := int(win.GetSystemMetrics(win.SM_CYSCREEN) / 24); i > 0; i-- {
		println("")
	}
	println("______________________________________________ ")
	println("|Comandos:                                   | ")
	println("|1: Ejecutar quicksort sobre el slice actual | ")
	println("|2: Ejecutar bubblesort sobre el slice actual| ")
	println("|3: Regresar el slice a su estado original   | ")
	println("|4: Crear un slice nuevo                     | ")
	println("|5: Salir                                    | ")
	println("---------------------------------------------- ")
}

func updateChart(){
	ui.Render(&bsChart)
	time.Sleep(MILI_SECONDS * time.Millisecond)
}

func barChartDriver(slice []float64) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	initBsChart(slice)
	initQsChart(slice)
	ui.Render(&bsChart)
	ui.Render(&qsChart)
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
				bsChart.Title = "BubbleSort"
				bsChart.Data = make([]float64, len(slice))
				copy(bsChart.Data, slice)
				bubbleSort(bsChart.Data)
			case "3":
				bsChart.Title = "Sort"
				bsChart.Data = make([]float64, len(slice))
				copy(bsChart.Data, slice)
				ui.Render(&bsChart)
			case "4":
				slice = randomSlice(111)
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
