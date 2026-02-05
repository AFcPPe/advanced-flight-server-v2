package pdu

// NetworkError 网络错误类型
type NetworkError int

const (
	NetworkErrorOk NetworkError = iota
	NetworkErrorCallsignInUse
	NetworkErrorCallsignInvalid
	NetworkErrorAlreadyRegistered
	NetworkErrorSyntaxError
	NetworkErrorPDUSourceInvalid
	NetworkErrorInvalidLogon
	NetworkErrorNoSuchCallsign
	NetworkErrorNoFlightPlan
	NetworkErrorNoWeatherProfile
	NetworkErrorInvalidProtocolRevision
	NetworkErrorRequestedLevelTooHigh
	NetworkErrorServerFull
	NetworkErrorCertificateSuspended
	NetworkErrorInvalidControl
	NetworkErrorInvalidPositionForRating
	NetworkErrorUnauthorizedSoftware
)

func (ne NetworkError) String() string {
	switch ne {
	case NetworkErrorOk:
		return "Ok"
	case NetworkErrorCallsignInUse:
		return "CallsignInUse"
	case NetworkErrorCallsignInvalid:
		return "CallsignInvalid"
	case NetworkErrorAlreadyRegistered:
		return "AlreadyRegistered"
	case NetworkErrorSyntaxError:
		return "SyntaxError"
	case NetworkErrorPDUSourceInvalid:
		return "PDUSourceInvalid"
	case NetworkErrorInvalidLogon:
		return "InvalidLogon"
	case NetworkErrorNoSuchCallsign:
		return "NoSuchCallsign"
	case NetworkErrorNoFlightPlan:
		return "NoFlightPlan"
	case NetworkErrorNoWeatherProfile:
		return "NoWeatherProfile"
	case NetworkErrorInvalidProtocolRevision:
		return "InvalidProtocolRevision"
	case NetworkErrorRequestedLevelTooHigh:
		return "RequestedLevelTooHigh"
	case NetworkErrorServerFull:
		return "ServerFull"
	case NetworkErrorCertificateSuspended:
		return "CertificateSuspended"
	case NetworkErrorInvalidControl:
		return "InvalidControl"
	case NetworkErrorInvalidPositionForRating:
		return "InvalidPositionForRating"
	case NetworkErrorUnauthorizedSoftware:
		return "UnauthorizedSoftware"
	default:
		return ""
	}
}

// NetworkFacility 网络设施类型
type NetworkFacility int

const (
	Observer NetworkFacility = iota
	FSS
	DEL
	GND
	TWR
	APP
	CTR
)

func (nf NetworkFacility) String() string {
	switch nf {
	case Observer:
		return "OBS"
	case FSS:
		return "FSS"
	case DEL:
		return "DEL"
	case GND:
		return "GND"
	case TWR:
		return "TWR"
	case APP:
		return "APP"
	case CTR:
		return "CTR"
	default:
		return ""
	}
}

// NetworkRating 网络等级
type NetworkRating int

const (
	RatingSUS NetworkRating = iota
	RatingOBS
	RatingS1
	RatingS2
	RatingS3
	RatingC1
	RatingC2
	RatingC3
	RatingI1
	RatingI2
	RatingI3
	RatingSUP
	RatingADM
)

// FlightRule 飞行规则
type FlightRule int

const (
	IFR FlightRule = iota
	VFR
)

func (fr FlightRule) String() string {
	if fr == IFR {
		return "I"
	}
	return "V"
}

// NetworkFlightPlan 网络飞行计划
type NetworkFlightPlan struct {
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
	Locked        bool
}
