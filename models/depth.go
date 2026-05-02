package projectx

type GatewayDepth struct {
	Timestamp     string  `json:"timestamp"`
	Type          int     `json:"type"`
	Price         float64 `json:"price"`
	Volume        int64   `json:"volume"`
	CurrentVolume int64   `json:"currentVolume"`
}

func GatewayDepthCSVHeader() string {
	return "timestamp,type,price,volume,current_volume"
}

func (d GatewayDepth) String() string {
	return fmt.Sprintf(
		"GatewayDepth{Timestamp:%s Type:%d Price:%f Volume:%d CurrentVolume:%d}",
		d.Timestamp,
		d.Type,
		d.Price,
		d.Volume,
		d.CurrentVolume,
	)
}

func (d GatewayDepth) ToCSVRow() string {
	return fmt.Sprintf(
		"%s,%d,%.2f,%d,%d",
		d.Timestamp,
		d.Type,
		d.Price,
		d.Volume,
		d.CurrentVolume,
	)
}
