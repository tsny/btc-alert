package eps

import (
	"fmt"

	chart "github.com/wcharczuk/go-chart"
)

// Creates a chart based on a CandleQueue
// See chart-example.png for an example
func QueueToGraph(q CandleQueue) chart.Chart {
	ser := chart.TimeSeries{Style: chart.StyleShow()}
	for _, c := range q.inner {
		ser.XValues = append(ser.XValues, c.Begin)
		ser.YValues = append(ser.YValues, c.Close)
	}
	graph := chart.Chart{
		Series: []chart.Series{ser},
		YAxis:  chart.YAxis{Style: chart.StyleShow()},
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeMinuteValueFormatter,
			Style: chart.Style{
				Show: true,
			},
		},
	}

	graph.TitleStyle = chart.StyleShow()
	graph.Title = fmt.Sprintf("%s [%s]", q.inner[0].Ticker, q.inner[0].Source)
	return graph
}
