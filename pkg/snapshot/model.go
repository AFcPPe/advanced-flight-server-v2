package snapshot

// Snapshot 快照顶层结构（写入Redis）
type Snapshot struct {
	Timestamp int64    `json:"timestamp"`
	Pilots    []*Pilot `json:"pilots"`
	ATCs      []*ATC   `json:"controllers"`
	ATISs     []*ATIS  `json:"atis"`
}

// Pilot 飞行员快照
type Pilot struct {
	CID              string      `json:"cid"`
	Callsign         string      `json:"callsign"`
	Name             string      `json:"name"`
	Lat              float64     `json:"lat"`
	Lon              float64     `json:"lon"`
	TrueAltitude     int         `json:"true_altitude"`
	PressureAltitude int         `json:"pressure_altitude"`
	GroundSpeed      int         `json:"ground_speed"`
	Pitch            float64     `json:"pitch"`
	Heading          float64     `json:"heading"`
	Bank             float64     `json:"bank"`
	Squawk           int         `json:"squawk"`
	SquawkingModeC   bool        `json:"squawking_mode_c"`
	Identing         bool        `json:"identing"`
	VisibilityRange  int         `json:"visibility_range"`
	FlightPlan       *FlightPlan `json:"flight_plan,omitempty"`
	FlightPlanLocked bool        `json:"flight_plan_locked"`
	Rating           int         `json:"rating"`
	LogonTime        string      `json:"logon_time"`
}

// ATC 管制员快照
type ATC struct {
	CID             string   `json:"cid"`
	Callsign        string   `json:"callsign"`
	Name            string   `json:"name"`
	Lat             float64  `json:"lat"`
	Lon             float64  `json:"lon"`
	Frequencies     []string `json:"frequencies"`
	Facility        int      `json:"facility"`
	Rating          int      `json:"rating"`
	VisibilityRange int      `json:"visibility_range"`
	TextAtis        []string `json:"text_atis,omitempty"`
	LogonTime       string   `json:"logon_time"`
}

// ATIS 自动终端情报服务快照
type ATIS struct {
	CID             string   `json:"cid"`
	Callsign        string   `json:"callsign"`
	Name            string   `json:"name"`
	Lat             float64  `json:"lat"`
	Lon             float64  `json:"lon"`
	Frequencies     []string `json:"frequencies"`
	Facility        int      `json:"facility"`
	Rating          int      `json:"rating"`
	VisibilityRange int      `json:"visibility_range"`
	TextAtis        []string `json:"text_atis,omitempty"`
	LogonTime       string   `json:"logon_time"`
}

// FlightPlan 飞行计划
type FlightPlan struct {
	FlightRules   string `json:"flight_rules"`
	AircraftType  string `json:"aircraft_type"`
	TAS           string `json:"tas"`
	Departure     string `json:"departure"`
	DepTime       string `json:"dep_time"`
	ActualDepTime string `json:"actual_dep_time"`
	CruiseAlt     string `json:"cruise_alt"`
	Arrival       string `json:"arrival"`
	EnrouteHour   string `json:"enroute_hour"`
	EnrouteMin    string `json:"enroute_min"`
	FobHour       string `json:"fob_hour"`
	FobMin        string `json:"fob_min"`
	Alternate     string `json:"alternate"`
	Route         string `json:"route"`
	Remark        string `json:"remark"`
}
