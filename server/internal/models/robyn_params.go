package models

type RobynParams struct {
	InputDataPath    string
	HolidaysDataPath string
	OutputDirectory  string
	CreateFiles      bool
	DepVar           string
	DepVarType       string
	ProphetVars      string // comma-separated list
	ProphetCountry   string
	PaidMediaSpends  string // comma-separated list
	PaidMediaVars    string // comma-separated list
	WindowStart      string
	WindowEnd        string
	AdstockType      string
	Hyperparameters  string // formatted hyperparameters list
	TrainSizeMin     float64
	TrainSizeMax     float64
	Cores            int
	Iterations       int
	Trials           int
	TSValidation     bool
	AddPenaltyFactor bool
	ParetoFronts     string
	CSVOut           string
	Clusters         bool
	ChannelConstrLow float64
	ChannelConstrUp  string // comma-separated list
	Scenario         string
}
