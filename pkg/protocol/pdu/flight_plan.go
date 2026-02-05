package pdu

import (
	"fmt"
	"strconv"
)

type FlightPlan struct {
	Base
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

func NewPDUFlightPlan(from, to string, flightRule int, flightType string, tas int, dep string, depTime int, actualDepTime int,
	cruiseAlt string, dest string, enrouteHour int, enrouteMin int, fobHour int, fobMin int, alterDest string, remark string, route string) *FlightPlan {
	return &FlightPlan{
		Base: Base{
			From: from,
			To:   to,
		},
		Callsign:      from,
		FlightRules:   FlightRule(flightRule),
		Type:          flightType,
		TAS:           strconv.Itoa(tas),
		Dep:           dep,
		DepTime:       strconv.Itoa(depTime),
		ActualDepTime: strconv.Itoa(actualDepTime),
		CruiseAlt:     cruiseAlt,
		Dest:          dest,
		EnrouteHour:   strconv.Itoa(enrouteHour),
		EnrouteMin:    strconv.Itoa(enrouteMin),
		FobHour:       strconv.Itoa(fobHour),
		FobMin:        strconv.Itoa(fobMin),
		AlterDest:     alterDest,
		Remark:        remark,
		Route:         route,
	}
}

func FlightPlanFromTokens(tokens []string) (*FlightPlan, error) {
	if len(tokens) < 17 {
		return nil, fmt.Errorf("flight plan: invalid token count, expected >= 17, got %d", len(tokens))
	}

	flightRule := VFR
	if tokens[2] == "I" {
		flightRule = IFR
	}

	tas, err := strconv.Atoi(tokens[4])
	if err != nil {
		tas = 0
	}

	depTime, err := strconv.Atoi(tokens[6])
	if err != nil {
		depTime = 0
	}

	actualDepTime, err := strconv.Atoi(tokens[7])
	if err != nil {
		actualDepTime = 0
	}

	enrouteHour, err := strconv.Atoi(tokens[10])
	if err != nil {
		enrouteHour = 0
	}

	enrouteMin, err := strconv.Atoi(tokens[11])
	if err != nil {
		enrouteMin = 0
	}

	fobHour, err := strconv.Atoi(tokens[12])
	if err != nil {
		fobHour = 0
	}

	fobMin, err := strconv.Atoi(tokens[13])
	if err != nil {
		fobMin = 0
	}

	return NewPDUFlightPlan(tokens[0], tokens[1], int(flightRule), tokens[3], tas, tokens[5], depTime, actualDepTime,
		tokens[8], tokens[9], enrouteHour, enrouteMin, fobHour, fobMin, tokens[14], tokens[15], tokens[16]), nil
}

func (pdu *FlightPlan) ToTokens() []string {
	flightRule := "V"
	if pdu.FlightRules == IFR {
		flightRule = "I"
	}

	return []string{
		pdu.Callsign,
		pdu.To,
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

func (pdu *FlightPlan) GetHeader() string {
	return "$FP"
}
