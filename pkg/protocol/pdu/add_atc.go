package pdu

import (
	"fmt"
	"strconv"
)

type AddATC struct {
	Base
	Callsign         string
	Name             string
	Cid              string
	Password         string
	Rating           NetworkRating
	ProtocolRevision int
}

func NewPDUAddATC(from, to, name, cid, password string, rating int, protocolRevision int) *AddATC {
	return &AddATC{
		Base: Base{
			From: from,
			To:   to,
		},
		Callsign:         from,
		Name:             name,
		Cid:              cid,
		Password:         password,
		Rating:           NetworkRating(rating),
		ProtocolRevision: protocolRevision,
	}
}

func AddATCFromTokens(tokens []string) (*AddATC, error) {
	if len(tokens) < 7 {
		return nil, fmt.Errorf("add atc: invalid token count, expected >= 7, got %d", len(tokens))
	}

	rating, err := strconv.Atoi(tokens[5])
	if err != nil {
		return nil, fmt.Errorf("add atc: invalid rating value: %w", err)
	}

	protocolRevision, err := strconv.Atoi(tokens[6])
	if err != nil {
		return nil, fmt.Errorf("add atc: invalid protocol revision: %w", err)
	}

	return NewPDUAddATC(tokens[0], tokens[1], tokens[2], tokens[3], tokens[4], rating, protocolRevision), nil
}

func (p *AddATC) ToTokens() []string {
	return []string{p.Callsign, p.To, p.Name, p.Cid, p.Password, strconv.Itoa(int(p.Rating)), strconv.Itoa(p.ProtocolRevision)}
}

func (p *AddATC) GetHeader() string {
	return "#AA"
}
