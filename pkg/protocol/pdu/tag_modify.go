package pdu

import "fmt"

type TagModify struct {
	From          string
	To            string
	Callsign      string
	FlightRules   FlightRule
	Type          string
	TAS           string
	Dep           string
	DepTime       string
	ActualDepTime string
	CruiseAlt     string
	Dest          string
	EnrouteHour   string
	EnrouteMin    string
	FobHour       string
	FobMin        string
	AlterDest     string
	Remark        string
	Route         string
}

func NewTagModify(from, to, callsign string, flightRule FlightRule, pType, tas, dep, depTime, actualDepTime, cruiseAlt, dest, enrouteHour, enrouteMin, fobHour, fobMin, alterDest, remark, route string) *TagModify {
	return &TagModify{
		From:          from,
		To:            to,
		Callsign:      callsign,
		FlightRules:   flightRule,
		Type:          pType,
		TAS:           tas,
		Dep:           dep,
		DepTime:       depTime,
		ActualDepTime: actualDepTime,
		CruiseAlt:     cruiseAlt,
		Dest:          dest,
		EnrouteHour:   enrouteHour,
		EnrouteMin:    enrouteMin,
		FobHour:       fobHour,
		FobMin:        fobMin,
		AlterDest:     alterDest,
		Remark:        remark,
		Route:         route,
	}
}

func TagModifyFromTokens(tokens []string) (*TagModify, error) {
	if len(tokens) < 18 {
		return nil, fmt.Errorf("tag modify: invalid token count, expected >= 18, got %d", len(tokens))
	}

	flightRule := VFR
	if tokens[3] == "I" {
		flightRule = IFR
	}

	return NewTagModify(
		tokens[0],
		tokens[1],
		tokens[2],
		flightRule,
		tokens[4],
		tokens[5],
		tokens[6],
		tokens[7],
		tokens[8],
		tokens[9],
		tokens[10],
		tokens[11],
		tokens[12],
		tokens[13],
		tokens[14],
		tokens[15],
		tokens[16],
		tokens[17],
	), nil
}

func (pdu *TagModify) ToTokens() []string {
	flightRule := "V"
	if pdu.FlightRules == IFR {
		flightRule = "I"
	}

	return []string{
		pdu.From,
		pdu.To,
		pdu.Callsign,
		flightRule,
		pdu.Type,
		pdu.TAS,
		pdu.Dep,
		pdu.DepTime,
		pdu.ActualDepTime,
		pdu.CruiseAlt,
		pdu.Dest,
		pdu.EnrouteHour,
		pdu.EnrouteMin,
		pdu.FobHour,
		pdu.FobMin,
		pdu.AlterDest,
		pdu.Remark,
		pdu.Route,
	}
}

func (pdu *TagModify) GetHeader() string {
	return "$AM"
}
