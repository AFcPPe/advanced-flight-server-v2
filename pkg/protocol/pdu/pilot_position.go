package pdu

import (
	"fmt"
	"strconv"
	"strings"
)

// PilotPosition 飞行员位置PDU (@ 开头)
type PilotPosition struct {
	Callsign       string
	Transponder    string
	Rating         NetworkRating
	Latitude       float64
	Longitude      float64
	TrueAltitude   int
	PressureAlt    int
	GroundSpeed    int
	Pitch          float64
	Heading        float64
	Bank           float64
}

// PilotPositionFromTokens 从tokens解析PilotPosition
// 格式: IdentFlag:Callsign:Transponder:Rating:Lat:Lon:TrueAlt:PressureAlt:GroundSpeed:Pitch:Heading:Bank
func PilotPositionFromTokens(tokens []string) (*PilotPosition, error) {
	if len(tokens) < 12 {
		return nil, fmt.Errorf("PilotPosition: invalid token count, got %d, need 12", len(tokens))
	}

	rating, err := strconv.Atoi(tokens[3])
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid rating: %w", err)
	}

	lat, err := strconv.ParseFloat(tokens[4], 64)
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid latitude: %w", err)
	}

	lon, err := strconv.ParseFloat(tokens[5], 64)
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid longitude: %w", err)
	}

	trueAlt, err := strconv.Atoi(tokens[6])
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid true altitude: %w", err)
	}

	pressureAlt, err := strconv.Atoi(tokens[7])
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid pressure altitude: %w", err)
	}

	groundSpeed, err := strconv.Atoi(tokens[8])
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid ground speed: %w", err)
	}

	pitch, err := strconv.ParseFloat(tokens[9], 64)
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid pitch: %w", err)
	}

	heading, err := strconv.ParseFloat(tokens[10], 64)
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid heading: %w", err)
	}

	bank, err := strconv.ParseFloat(tokens[11], 64)
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid bank: %w", err)
	}

	return &PilotPosition{
		Callsign:     tokens[1],
		Transponder:  tokens[2],
		Rating:       NetworkRating(rating),
		Latitude:     lat,
		Longitude:    lon,
		TrueAltitude: trueAlt,
		PressureAlt:  pressureAlt,
		GroundSpeed:  groundSpeed,
		Pitch:        pitch,
		Heading:      heading,
		Bank:         bank,
	}, nil
}

// PilotPositionFromRaw 从原始数据解析 (去掉@前缀后的数据)
func PilotPositionFromRaw(data []byte) (*PilotPosition, error) {
	s := string(data)
	if len(s) > 0 && s[0] == '@' {
		s = s[1:]
	}
	tokens := strings.Split(s, ":")
	return PilotPositionFromTokens(tokens)
}

func (p *PilotPosition) GetHeader() string {
	return "@"
}

func (p *PilotPosition) ToTokens() []string {
	return []string{
		"", // IdentFlag placeholder
		p.Callsign,
		p.Transponder,
		strconv.Itoa(int(p.Rating)),
		strconv.FormatFloat(p.Latitude, 'f', 6, 64),
		strconv.FormatFloat(p.Longitude, 'f', 6, 64),
		strconv.Itoa(p.TrueAltitude),
		strconv.Itoa(p.PressureAlt),
		strconv.Itoa(p.GroundSpeed),
		strconv.FormatFloat(p.Pitch, 'f', 2, 64),
		strconv.FormatFloat(p.Heading, 'f', 2, 64),
		strconv.FormatFloat(p.Bank, 'f', 2, 64),
	}
}
