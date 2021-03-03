package priceTracking

import (
	chart "github.com/wcharczuk/go-chart"
)

func QueueToGraph(q CandleQueue) chart.Chart {
	ser := chart.TimeSeries{Style: chart.StyleShow()}
	for _, c := range q.inner {
		ser.XValues = append(ser.XValues, c.Begin)
		ser.YValues = append(ser.YValues, c.Close)
	}
	graph := chart.Chart{Series: []chart.Series{ser}}
	graph.TitleStyle = chart.StyleShow()
	graph.Title = q.inner[0].Ticker + q.inner[0].Source
	return graph
}
