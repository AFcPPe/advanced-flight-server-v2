package pdu

import (
	"fmt"
	"strconv"
	"strings"
)

// ATCPosition ATC位置PDU (% 开头)
type ATCPosition struct {
	Callsign        string
	Frequencies     []string
	Facility        NetworkFacility
	VisibilityRange int
	Rating          NetworkRating
	Latitude        float64
	Longitude       float64
}

// NewATCPosition 创建ATCPosition
func NewATCPosition(callsign string, frequencies []string, facility NetworkFacility, visRange int, rating NetworkRating, lat, lon float64) *ATCPosition {
	return &ATCPosition{
		Callsign:        callsign,
		Frequencies:     frequencies,
		Facility:        facility,
		VisibilityRange: visRange,
		Rating:          rating,
		Latitude:        lat,
		Longitude:       lon,
	}
}

// ATCPositionFromTokens 从tokens解析ATCPosition
// 格式: Callsign:Frequencies:Facility:VisRange:Rating:Lat:Lon
// Frequencies 多个频率用 & 分隔
func ATCPositionFromTokens(tokens []string) (*ATCPosition, error) {
	if len(tokens) < 7 {
		return nil, fmt.Errorf("ATCPosition: invalid token count, got %d, need 7", len(tokens))
	}

	frequencies := strings.Split(tokens[1], "&")

	facility, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil, fmt.Errorf("ATCPosition: invalid facility: %w", err)
	}

	visRange, err := strconv.Atoi(tokens[3])
	if err != nil {
		return nil, fmt.Errorf("ATCPosition: invalid visibility range: %w", err)
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

	return NewATCPosition(
		tokens[0],
		frequencies,
		NetworkFacility(facility),
		visRange,
		NetworkRating(rating),
		lat,
		lon,
	), nil
}

// ATCPositionFromRaw 从原始数据解析 (去掉%前缀后的数据)
func ATCPositionFromRaw(data []byte) (*ATCPosition, error) {
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
	freq := strings.Join(p.Frequencies, "&")
	return []string{
		p.Callsign,
		freq,
		strconv.Itoa(int(p.Facility)),
		strconv.Itoa(p.VisibilityRange),
		strconv.Itoa(int(p.Rating)),
		strconv.FormatFloat(p.Latitude, 'f', -1, 64),
		strconv.FormatFloat(p.Longitude, 'f', -1, 64),
		"0",
	}
}
