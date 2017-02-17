package fcm

type option interface {
	apply(*HTTPClient)
}

type optionHTTPDoer struct {
	v HTTPDoer
}

func OptionHTTPDoer(v HTTPDoer) *optionHTTPDoer {
	return &optionHTTPDoer{v}
}

func (o *optionHTTPDoer) apply(c *HTTPClient) {
	c.http = o.v
}
