package eps

const capLimit = 10000

type CandleQueue struct {
	Newest *Candlestick
	Oldest *Candlestick
	Array  []*Candlestick
}

func NewQueue(p *Publisher) *CandleQueue {
	q := &CandleQueue{}
	q.Subscribe(p)
	return q
}

// Subscribe makes the queue listen for the publisher's events/updates
func (q *CandleQueue) Subscribe(p *Publisher) {
	p.RegisterPriceUpdateListener(q.handlePriceUpdate)
}

func (q *CandleQueue) handlePriceUpdate(p *Publisher, c *Candlestick, completed bool) {
	// Ignore incomplete candles
	if !completed {
		return
	}
	q.Enqueue(c)
}

func (q *CandleQueue) GetAllPrices() []float64 {
	var arr []float64
	for _, v := range q.Array {
		arr = append(arr, v.Close)
	}
	return arr
}

func (q *CandleQueue) Enqueue(c *Candlestick) {
	// todo: make this better than just truncating
	// this is actually backwards right now which looks fine in graphs
	if len(q.Array) > capLimit {
		q.Array = q.Array[:capLimit-1]
	}
	q.Array = append(q.Array, c)
	q.Newest = c
}

func (q *CandleQueue) GetQueue() []*Candlestick {
	return q.Array
}
