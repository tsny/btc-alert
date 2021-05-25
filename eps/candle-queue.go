package eps

const capLimit = 10000

type CandleQueue struct {
	Oldest Candlestick
	Newest Candlestick
	inner  []Candlestick
	cap    int
}

func NewQueue() *CandleQueue {
	return &CandleQueue{}
}

func (q *CandleQueue) Subscribe(p *Publisher) *CandleQueue {
	p.RegisterSubscriber(q.handlePriceUpdate)
	return q
}

func (q *CandleQueue) handlePriceUpdate(p *Publisher, c Candlestick) {
	// Ignore incomplete candles
	if !c.Complete {
		return
	}
	q.Add(c)
}

func (q *CandleQueue) GetAllPrices() []float64 {
	var arr []float64
	for _, v := range q.inner {
		arr = append(arr, v.Close)
	}
	return arr
}

func (q *CandleQueue) Add(c Candlestick) {
	// todo: make this better than just truncating
	// this is actually backwards right now which looks fine in graphs
	if len(q.inner) > q.cap {
		q.inner = q.inner[:capLimit-1]
	}
	q.inner = append(q.inner, c)
	q.Newest = c
}

func (q *CandleQueue) GetCap() int {
	return q.cap
}

func (q *CandleQueue) SetCap(cap int) {
	if cap > capLimit {
		cap = capLimit
	}
	q.cap = cap
}

func (q *CandleQueue) GetQueue() []Candlestick {
	return q.inner
}
