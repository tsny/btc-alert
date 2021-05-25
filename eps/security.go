package eps

type Security struct {
	Name            string
	Ticker          string
	Type            SecurityType
	Source          string
	AdditionalNames []string
}

type SecurityType int

const (
	Stock SecurityType = iota
	Crypto
)

func (s *SecurityType) String() string {
	switch *s {
	case 0:
		return "Stock"
	case 1:
		return "Crypto"
	}
	return "Unknown"
}

func NewCrypto(name, ticker, source string, addlNames ...string) *Security {
	return &Security{name, ticker, Crypto, source, addlNames}
}

func NewStock(name, ticker, source string, addlNames ...string) *Security {
	return &Security{name, ticker, Stock, source, addlNames}
}
