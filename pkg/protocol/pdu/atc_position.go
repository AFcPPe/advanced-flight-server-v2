package pdu

import (
	"fmt"
	"strconv"
	"strings"
)

type ATCPosition struct {
	Base
	Frequencies     []string
	Facility        NetworkFacility
	VisibilityRange int
	Rating          NetworkRating
	Lat             float64
	Lon             float64
}

func NewPDUATCPosition(from string, freq []string, facility, visRange, rating int, lat, lon float64) *ATCPosition {
	return &ATCPosition{
		Base: Base{
			From: from,
			To:   "",
		},
		Frequencies:     freq,
		Facility:        NetworkFacility(facility),
		VisibilityRange: visRange,
		Rating:          NetworkRating(rating),
		Lat:             lat,
		Lon:             lon,
	}
}

func ATCPositionFromTokens(tokens []string) (*ATCPosition, error) {
	if len(tokens) < 7 {
		return nil, fmt.Errorf("atc position: invalid token count, expected >= 7, got %d", len(tokens))
	}
	freq := strings.Split(tokens[1], "&")
	facility, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil, err
	}
	visRange, err := strconv.Atoi(tokens[3])
	if err != nil {
		return nil, err
	}
	rating, err := strconv.Atoi(tokens[4])
	if err != nil {
		return nil, err
	}
	lat, err := strconv.ParseFloat(tokens[5], 64)
	if err != nil {
		return nil, err
	}
	lon, err := strconv.ParseFloat(tokens[6], 64)
	if err != nil {
		return nil, err
	}
	return NewPDUATCPosition(tokens[0], freq, facility, visRange, rating, lat, lon), nil
}

func (pdu *ATCPosition) ToTokens() []string {
	freq := strings.Join(pdu.Frequencies, "&")
	return []string{
		pdu.From,
		freq,
		strconv.Itoa(int(pdu.Facility)),
		strconv.Itoa(pdu.VisibilityRange),
		strconv.Itoa(int(pdu.Rating)),
		strconv.FormatFloat(pdu.Lat, 'f', -1, 64),
		strconv.FormatFloat(pdu.Lon, 'f', -1, 64),
		"0",
	}
}

func (pdu *ATCPosition) GetHeader() string {
	return "%"
}
