package pdu

import (
	"fmt"
	"strconv"
	"strings"
)

// ATCPosition ATC位置PDU (% 开头)
type ATCPosition struct {
	Callsign       string
	Frequency      int
	Facility       int
	VisualRange    int
	Rating         NetworkRating
	Latitude       float64
	Longitude      float64
	AltitudeAGL    int
}

// ATCPositionFromTokens 从tokens解析ATCPosition
// 格式: Callsign:Frequency:Facility:VisualRange:Rating:Lat:Lon:AltitudeAGL
func ATCPositionFromTokens(tokens []string) (*ATCPosition, error) {
	if len(tokens) < 8 {
		return nil, fmt.Errorf("ATCPosition: invalid token count, got %d, need 8", len(tokens))
	}

	frequency, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, fmt.Errorf("ATCPosition: invalid frequency: %w", err)
	}

	facility, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil, fmt.Errorf("ATCPosition: invalid facility: %w", err)
	}

	visualRange, err := strconv.Atoi(tokens[3])
	if err != nil {
		return nil, fmt.Errorf("ATCPosition: invalid visual range: %w", err)
	}

	rating, err := strconv.Atoi(tokens[4])
	if err != nil {
		return nil, fmt.Errorf("ATCPosition: invalid rating: %w", err)
	}

	lat, err := strconv.ParseFloat(tokens[5], 64)
	if err != nil {
		return nil, fmt.Errorf("ATCPosition: invalid latitude: %w", err)
	}

	lon, err := strconv.ParseFloat(tokens[6], 64)
	if err != nil {
		return nil, fmt.Errorf("ATCPosition: invalid longitude: %w", err)
	}

	alt, err := strconv.Atoi(tokens[7])
	if err != nil {
		return nil, fmt.Errorf("ATCPosition: invalid altitude: %w", err)
	}

	return &ATCPosition{
		Callsign:    tokens[0],
		Frequency:   frequency,
		Facility:    facility,
		VisualRange: visualRange,
		Rating:      NetworkRating(rating),
		Latitude:    lat,
		Longitude:   lon,
		AltitudeAGL: alt,
	}, nil
}

// ATCPositionFromRaw 从原始数据解析 (去掉%前缀后的数据)
func ATCPositionFromRaw(data []byte) (*ATCPosition, error) {
	// 去掉%前缀
	s := string(data)
	if len(s) > 0 && s[0] == '%' {
		s = s[1:]
	}
	tokens := strings.Split(s, ":")
	return ATCPositionFromTokens(tokens)
}

func (p *ATCPosition) GetHeader() string {
	return "%"
}

func (p *ATCPosition) ToTokens() []string {
	return []string{
		p.Callsign,
		strconv.Itoa(p.Frequency),
		strconv.Itoa(p.Facility),
		strconv.Itoa(p.VisualRange),
		strconv.Itoa(int(p.Rating)),
		strconv.FormatFloat(p.Latitude, 'f', 6, 64),
		strconv.FormatFloat(p.Longitude, 'f', 6, 64),
		strconv.Itoa(p.AltitudeAGL),
	}
}
