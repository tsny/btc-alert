package eps

const capLimit = 10000

type CandleStack struct {
	Newest *Candlestick
	Oldest *Candlestick
	Array  []*Candlestick
}

func NewStack(p *Publisher) *CandleStack {
	q := &CandleStack{}
	q.Subscribe(p)
	return q
}

// Subscribe makes the stack listen for the publisher's events/updates
func (q *CandleStack) Subscribe(p *Publisher) {
	p.RegisterPriceUpdateListener(q.handlePriceUpdate)
}

func (q *CandleStack) handlePriceUpdate(p *Publisher, c *Candlestick, completed bool) {
	// Ignore incomplete candles
	if !completed {
		return
	}
	q.Push(c)
}

func (q *CandleStack) GetAllPrices() []float64 {
	var arr []float64
	for _, v := range q.Array {
		arr = append(arr, v.Close)
	}
	return arr
}

func (q *CandleStack) Push(c *Candlestick) {
	// todo: make this better than just truncating
	// this is actually backwards right now which looks fine in graphs
	if len(q.Array) > capLimit {
		q.Array = q.Array[:capLimit-1]
	}
	q.Array = append(q.Array, c)
	q.Newest = c
}
