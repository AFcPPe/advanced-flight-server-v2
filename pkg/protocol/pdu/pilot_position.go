package pdu

import (
	"fmt"
	"strconv"
	"strings"
)

// PilotPosition 飞行员位置PDU (@ 开头)
type PilotPosition struct {
	Callsign         string
	SquawkCode       int
	SquawkingModeC   bool
	Identing         bool
	Rating           NetworkRating
	Latitude         float64
	Longitude        float64
	TrueAltitude     int
	PressureAltitude int
	GroundSpeed      int
	Pitch            float64
	Heading          float64
	Bank             float64
}

// NewPilotPosition 创建PilotPosition
func NewPilotPosition(callsign string, squawkCode int, squawkingModeC, identing bool, rating NetworkRating, lat, lon float64, trueAlt, pressureAlt, gs int, pitch, heading, bank float64) *PilotPosition {
	return &PilotPosition{
		Callsign:         callsign,
		SquawkCode:       squawkCode,
		SquawkingModeC:   squawkingModeC,
		Identing:         identing,
		Rating:           rating,
		Latitude:         lat,
		Longitude:        lon,
		TrueAltitude:     trueAlt,
		PressureAltitude: pressureAlt,
		GroundSpeed:      gs,
		Pitch:            pitch,
		Heading:          heading,
		Bank:             bank,
	}
}

// PilotPositionFromTokens 从tokens解析PilotPosition
// 格式: IdentFlag:Callsign:SquawkCode:Rating:Lat:Lon:TrueAlt:GroundSpeed:PBH:DeltaAlt
// IdentFlag: S=普通, N=ModeC, I=Ident
// PBH: 打包的Pitch/Bank/Heading
// DeltaAlt: PressureAltitude - TrueAltitude
func PilotPositionFromTokens(tokens []string) (*PilotPosition, error) {
	if len(tokens) < 10 {
		return nil, fmt.Errorf("PilotPosition: invalid token count, got %d, need 10", len(tokens))
	}

	identFlag := tokens[0]
	squawkingModeC := identFlag == "N"
	identing := identFlag == "I"

	squawkCode, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid squawk code: %w", err)
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

	gs, err := strconv.Atoi(tokens[7])
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid ground speed: %w", err)
	}

	pbh, err := strconv.ParseInt(tokens[8], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid PBH: %w", err)
	}
	pitch, bank, heading := unpackPBH(pbh)

	deltaAlt, err := strconv.Atoi(tokens[9])
	if err != nil {
		return nil, fmt.Errorf("PilotPosition: invalid delta altitude: %w", err)
	}
	pressureAlt := trueAlt + deltaAlt

	return NewPilotPosition(
		tokens[1],
		squawkCode,
		squawkingModeC,
		identing,
		NetworkRating(rating),
		lat,
		lon,
		trueAlt,
		pressureAlt,
		gs,
		pitch,
		heading,
		bank,
	), nil
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
	identFlag := "S"
	if p.Identing {
		identFlag = "I"
	} else if p.SquawkingModeC {
		identFlag = "N"
	}

	return []string{
		identFlag,
		p.Callsign,
		strconv.Itoa(p.SquawkCode),
		strconv.Itoa(int(p.Rating)),
		strconv.FormatFloat(p.Latitude, 'f', -1, 64),
		strconv.FormatFloat(p.Longitude, 'f', -1, 64),
		strconv.Itoa(p.TrueAltitude),
		strconv.Itoa(p.GroundSpeed),
		strconv.FormatInt(packPBH(p.Pitch, p.Bank, p.Heading), 10),
		strconv.Itoa(p.PressureAltitude - p.TrueAltitude),
	}
}

// unpackPBH 从打包的int64中解包pitch, bank, heading
func unpackPBH(pbh int64) (pitch, bank, heading float64) {
	pitchInt := pbh >> 22
	pitch = float64(pitchInt) / 1024.0 * -360.0
	if pitch > 180.0 {
		pitch -= 360.0
	} else if pitch <= -180.0 {
		pitch += 360.0
	}

	bankInt := (pbh >> 12) & 0x3FF
	bank = float64(bankInt) / 1024.0 * -360.0
	if bank > 180.0 {
		bank -= 360.0
	} else if bank <= -180.0 {
		bank += 360.0
	}

	hdgInt := (pbh >> 2) & 0x3FF
	heading = float64(hdgInt) / 1024.0 * 360.0
	if heading < 0.0 {
		heading += 360.0
	} else if heading >= 360.0 {
		heading -= 360.0
	}
	return
}

// packPBH 将pitch, bank, heading打包为int64
func packPBH(pitch, bank, heading float64) int64 {
	p := pitch / -360.0
	if p < 0 {
		p += 1.0
	}
	p *= 1024.0

	b := bank / -360.0
	if b < 0 {
		b += 1.0
	}
	b *= 1024.0

	h := heading / 360.0 * 1024.0

	return int64(p)<<22 | int64(b)<<12 | int64(h)<<2
}
