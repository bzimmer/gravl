package ngrok

type Fault struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"msg"`
	Details    struct {
		Path string `json:"path"`
	} `json:"details"`
}

func (f *Fault) Error() string {
	return f.Message
}

type Config struct {
	Address string `json:"addr"`
	Inspect bool   `json:"inspect"`
}

type Conns struct {
	Count  int     `json:"count"`
	Gauge  int     `json:"gauge"`
	Rate1  float64 `json:"rate1"`
	Rate5  float64 `json:"rate5"`
	Rate15 float64 `json:"rate15"`
	P50    float64 `json:"p50"`
	P90    float64 `json:"p90"`
	P95    float64 `json:"p95"`
	P99    float64 `json:"p99"`
}

type HTTP struct {
	Count  int     `json:"count"`
	Rate1  float64 `json:"rate1"`
	Rate5  float64 `json:"rate5"`
	Rate15 float64 `json:"rate15"`
	P50    float64 `json:"p50"`
	P90    float64 `json:"p90"`
	P95    float64 `json:"p95"`
	P99    float64 `json:"p99"`
}

type Metrics struct {
	Connections *Conns `json:"conns"`
	HTTP        *HTTP  `json:"http"`
}

type Tunnel struct {
	Name      string   `json:"name"`
	URI       string   `json:"uri"`
	PublicURL string   `json:"public_url"`
	Proto     string   `json:"proto"`
	Config    *Config  `json:"config"`
	Metrics   *Metrics `json:"metrics"`
}
