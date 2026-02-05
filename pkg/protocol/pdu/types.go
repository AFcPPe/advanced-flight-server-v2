package pdu

// ClientType 客户端类型
type ClientType int

const (
	ClientTypePilot ClientType = iota // 飞行员
	ClientTypeATC                     // 管制员
)

func (c ClientType) String() string {
	switch c {
	case ClientTypePilot:
		return "Pilot"
	case ClientTypeATC:
		return "ATC"
	default:
		return "Unknown"
	}
}

// NetworkRating 网络等级
type NetworkRating int

const (
	RatingSUS NetworkRating = iota // Suspended
	RatingOBS                      // Observer
	RatingS1                       // Student 1
	RatingS2                       // Student 2
	RatingS3                       // Student 3
	RatingC1                       // Controller 1
	RatingC2                       // Controller 2
	RatingC3                       // Controller 3
	RatingI1                       // Instructor 1
	RatingI2                       // Instructor 2
	RatingI3                       // Instructor 3
	RatingSUP                      // Supervisor
	RatingADM                      // Administrator
)

func (r NetworkRating) String() string {
	switch r {
	case RatingSUS:
		return "SUS"
	case RatingOBS:
		return "OBS"
	case RatingS1:
		return "S1"
	case RatingS2:
		return "S2"
	case RatingS3:
		return "S3"
	case RatingC1:
		return "C1"
	case RatingC2:
		return "C2"
	case RatingC3:
		return "C3"
	case RatingI1:
		return "I1"
	case RatingI2:
		return "I2"
	case RatingI3:
		return "I3"
	case RatingSUP:
		return "SUP"
	case RatingADM:
		return "ADM"
	default:
		return "Unknown"
	}
}

// NetworkError 网络错误码
type NetworkError int

const (
	NetworkErrorOK                       NetworkError = iota // 无错误
	NetworkErrorCallsignInUse                                // 呼号已被使用
	NetworkErrorCallsignInvalid                              // 呼号无效
	NetworkErrorAlreadyRegistered                            // 已注册
	NetworkErrorSyntaxError                                  // 语法错误
	NetworkErrorPDUSourceInvalid                             // PDU来源无效
	NetworkErrorInvalidLogon                                 // 无效登录
	NetworkErrorNoSuchCallsign                               // 呼号不存在
	NetworkErrorNoFlightPlan                                 // 无飞行计划
	NetworkErrorNoWeatherProfile                             // 无天气配置
	NetworkErrorInvalidProtocolRevision                      // 无效协议版本
	NetworkErrorRequestedLevelTooHigh                        // 请求的等级过高
	NetworkErrorServerFull                                   // 服务器已满
	NetworkErrorCertificateSuspended                         // 证书已暂停
	NetworkErrorInvalidControl                               // 无效管制
	NetworkErrorInvalidPositionForRating                     // 等级不匹配的位置
	NetworkErrorUnauthorizedSoftware                         // 未授权软件
)

func (ne NetworkError) String() string {
	switch ne {
	case NetworkErrorOK:
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

// NetworkFacility ATC设施类型
type NetworkFacility int

const (
	FacilityOBS NetworkFacility = iota // Observer
	FacilityFSS                        // Flight Service Station
	FacilityDEL                        // Delivery
	FacilityGND                        // Ground
	FacilityTWR                        // Tower
	FacilityAPP                        // Approach
	FacilityCTR                        // Center
)

func (nf NetworkFacility) String() string {
	switch nf {
	case FacilityOBS:
		return "OBS"
	case FacilityFSS:
		return "FSS"
	case FacilityDEL:
		return "DEL"
	case FacilityGND:
		return "GND"
	case FacilityTWR:
		return "TWR"
	case FacilityAPP:
		return "APP"
	case FacilityCTR:
		return "CTR"
	default:
		return ""
	}
}
