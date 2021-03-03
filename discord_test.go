package main

import (
	"bytes"
	"io"
	"testing"

	chart "github.com/wcharczuk/go-chart"
)

func TestCryptoBot_SendGraph(t *testing.T) {
	cb := initBot("")
	graph := chart.Chart{
		Title:      "Test Title",
		TitleStyle: chart.StyleShow(),
		XAxis:      chart.XAxis{Style: chart.StyleShow()},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 1.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0},
				Name:    "Graph BTC",
			},
		},
	}
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		content string
		reader  io.Reader
	}
	tests := []struct {
		name string
		cb   *CryptoBot
		args args
	}{
		{
			name: "test1",
			cb:   cb,
			args: args{
				content: "test.png",
				reader:  buffer,
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cb.SendGraph(tt.args.content, tt.args.reader)
		})
	}
}
