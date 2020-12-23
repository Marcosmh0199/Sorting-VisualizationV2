package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var(
	bsBarChart = charts.NewBar()
	qsBarChart = charts.NewBar()
	baseSlice = randomSlice(MAX_PRIME, 50)
)

func main() {
	/*
		fmt.Print("Indique el tamaño del slice: ")
		var first int
		fmt.Scanln(&first)

		var bsSlice = make([]float64, first)
		copy(bsSlice, slice)
		var qsSlice = make([]float64, first)
		copy(qsSlice, slice)
	*/
	initBsBarChart()
	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8081", nil)
	/*page := components.NewPage()
	page.AddCharts(bsBarChart)
	f, err := os.Create("sorts/html/sorts.html")
	if err != nil{
		panic(err)
	}
	page.Render(io.MultiWriter(f))
	fs := http.FileServer(http.Dir("sorts/html"))
	log.Println("running server at http://localhost:8089")
	log.Fatal(http.ListenAndServe("localhost:8089", logRequest(fs)))*/
}

func generateLabels() []string {
	var labels = make([]string, len(baseSlice))
	for i := range baseSlice{
		labels[i] = strconv.Itoa(i)
	}
	return labels
}

func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := range baseSlice {
		items = append(items, opts.BarData{Value: baseSlice[i]})
	}
	return items
}

func generateBarRandItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}

func initBsBarChart()  {
	bsBarChart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "BubbleSort",
	}))
	bsBarChart.SetXAxis(generateLabels()).
		AddSeries("Values", generateBarItems())
}

func initQsBarChart(slice []float64){

}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	for {
		bsBarChart.Render(w)
		time.Sleep(5 * time.Second)
		bsBarChart.AddSeries("Values", generateBarRandItems())
	}
}

