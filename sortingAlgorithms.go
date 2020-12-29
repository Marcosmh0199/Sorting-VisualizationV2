/*
Notas:
1. Para una visualización correcta del gráfico de barras, por favor ejecutar el programa
en una terminal a pantalla completa.
2. Para ejecutar de nuevo, presionar una tecla distinta de 0 y enter luego de que ambos algoritmos hayan acabado.
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
	"regexp"
	"strconv"
	"sync"
	"time"
)

const (
	BAR_WIDTH = 3
	FONT_WIDTH = 8
	FONT_HEIGHT = 16
	MAX_NUMBER_SIZE = 32
)

var(
	width int = int(win.GetSystemMetrics(win.SM_CXSCREEN) / FONT_WIDTH)
	height int = int(win.GetSystemMetrics(win.SM_CYSCREEN) / (FONT_HEIGHT*2))
	bsChart widgets.BarChart
	qsChart widgets.BarChart
	m sync.Mutex

	bsTime time.Duration
	bsSwaps = 0
	bsComparisons = 0
	bsIterations = 0

	qsTime time.Duration
	qsSwaps = 0
	qsComparisons = 0
	qsIterations = 0
)


func main() {
	barNumber := width / BAR_WIDTH - 1
	fmt.Print("Indique la cantidad de numeros(Se recomienda " + strconv.Itoa(barNumber) +" maximo para una visualizacion correcta): ")
	var size int
	fmt.Scanln(&size)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	baseSlice := randomSlice(size)
	initBsChart(baseSlice)
	initQsChart(baseSlice)
	ui.Render(&bsChart)
	ui.Render(&qsChart)
	go bsChartDrawer(baseSlice)
	qsChartDrawer(baseSlice)
	fmt.Scanln() //end until any key is pressed
	ui.Close()
}

/*
Generates a 3 digit number between 0-559 from the system hour
*/
func generateSeed() int{
	currentTime := time.Now()
	seed := currentTime.Format(":05.0")
	var regex = regexp.MustCompile("[0-9]+")
	bytes := regex.FindAll([]byte(seed),2)
	num, _ := strconv.Atoi(string(bytes[0])+string(bytes[1]))
	return num
}

/*
creates a N size slice with random numbers based on the linear congruential method
output: slice with N random integers
*/
func randomSlice(size int) []float64{
	var slice = make([]float64, size)
	var m = int(math.Pow(2, MAX_NUMBER_SIZE))  // modulus
	var a = 22695477                    // multiplier
	var c = 1                           //increment
	seed := generateSeed()
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
	initTime := time.Now()
	n := len(slice) - 1
	for true {
		swapped := false
		for i := 0; i < n; i++{
			bsComparisons++
			if slice[i] > slice[i+1]{
				pair <- []int{i, i+1}
				bsSwaps++
				swapped = true
			}
		}
		bsIterations++
		if !swapped{
			break
		}
		n--
	}
	close(pair)
	bsTime = time.Since(initTime)
}

func partition(slice []float64, start int, end int, pair chan []int) int {
	pivot := slice[end]
	index := start
	for i := start; i < end; i++{
		qsComparisons++
		if slice[i] <= pivot{
			pair <- []int{i,index}
			qsSwaps++
			index++
		}
	}
	pair <- []int{index,end}
	qsSwaps++
	return index
}

/*
Iterative quicksort
 */
func quickSort(slice []float64, size int, pair chan []int) {
	initTime := time.Now()
	var stack [][]int
	start := 0
	end := size
	stack = append(stack, []int{start,end})
	for len(stack) > 0{
		start, end = stack[0][0],stack[0][1]
		stack = remove(stack, 0)
		pivot := partition(slice, start, end, pair)
		qsComparisons++
		if pivot-1 > start {
			stack = append(stack, []int{start,pivot-1})
		}
		qsComparisons++
		if pivot+1 < end {
			stack = append(stack, []int{pivot+1,end})
		}
		qsIterations++
	}
	close(pair)
	qsTime = time.Since(initTime)
}

/*
Bubblesort graphic drawer
 */
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
	bsChart.Title = "BubbleSort-Finalizado-" +
		"Tiempo:"+strconv.FormatInt(bsTime.Milliseconds(),10)+"ms-" +
		"Swaps:"+strconv.Itoa(bsSwaps)+"-" +
		"Comparaciones:"+strconv.Itoa(bsComparisons)+"-"+
		"Iteraciones:"+strconv.Itoa(bsIterations)
	m.Lock()
	ui.Render(&bsChart)
	m.Unlock()
}

/*
Quicksort graphic drawer
*/
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
	qsChart.Title = "QuickSort-Finalizado-" +
		"Tiempo:"+strconv.FormatInt(qsTime.Milliseconds(),10)+"ms-" +
		"Swaps:"+strconv.Itoa(qsSwaps)+"-" +
		"Comparaciones:"+strconv.Itoa(qsComparisons)+"-"+
		"Iteraciones:"+strconv.Itoa(qsIterations)
	m.Lock()
	ui.Render(&qsChart)
	m.Unlock()
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
	qsChart.BarWidth = BAR_WIDTH
	qsChart.BarGap = 0
	qsChart.Labels = generateLabels(slice)
	qsChart.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
	qsChart.BarColors = []ui.Color{ui.ColorRed}
	qsChart.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack)}
}