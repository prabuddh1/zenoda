package types

import (
	"fmt"

	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyInflationRate     = []byte("InflationRate")
	KeyPredefinedWallets = []byte("PredefinedWallets")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable creates the key table for rewards module parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(inflationRate math.LegacyDec, predefinedWallets []string) Params {
	return Params{
		InflationRate:     inflationRate.String(), // Keep InflationRate as a string
		PredefinedWallets: predefinedWallets,      // List of governance wallets
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		math.LegacyMustNewDecFromStr("0.05"), // Default 5% inflation rate
		[]string{
			"cosmos1lhahcqzx45mssr9wfknx48hy4truyz9p2wj3ht",
			"cosmos1g6k8qf0zksqruq8exv0duw3p9fn33aeffdprl6",
			"cosmos1fmatdzaqv75sxte07nt67arx398v2cdudv25fl",
			"cosmos1hfc7lhageqsln6kj2e225q6hj2muv27ck7ajzx",
			"cosmos1eag08vq68vpmrycaxqs8dlsveu97a6ncxxx9ll",
			"cosmos1v5w82p3x9a4r82lfuztwvzz3hjxtvrx8xz66pu",
			"cosmos13730mjpw988h7wx3dwv0vg2sqm2udjvptjzanh",
			"cosmos13kp6sagqp27f02eetgpp7g35mjxlwsngsk53dm",
			"cosmos1nvwepluydj7xnga6qud3cl46juft7rrnktx5as",
			"cosmos1gj2yqrzkdd9q7yedcasvvke5ls5ahkfq0gm6x4",
		},
	)
}

// ParamSetPairs defines the parameter set pairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyInflationRate, &p.InflationRate, validateInflationRate),
		paramtypes.NewParamSetPair(KeyPredefinedWallets, &p.PredefinedWallets, validatePredefinedWallets),
	}
}

// Validate validates the parameters
func (p Params) Validate() error {
	if err := validateInflationRate(p.InflationRate); err != nil {
		return err
	}
	if err := validatePredefinedWallets(p.PredefinedWallets); err != nil {
		return err
	}
	return nil
}

// validateInflationRate ensures the inflation rate is between 0 and 1
func validateInflationRate(i interface{}) error {
	inflationRateStr, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	inflationRate, err := math.LegacyNewDecFromStr(inflationRateStr)
	if err != nil {
		return fmt.Errorf("invalid inflation rate format: %v", err)
	}

	if inflationRate.IsNegative() || inflationRate.GT(math.LegacyNewDec(1)) {
		return fmt.Errorf("inflation rate must be between 0 and 1 (inclusive)")
	}
	return nil
}

// validatePredefinedWallets ensures that all provided addresses are valid Bech32 addresses
func validatePredefinedWallets(i interface{}) error {
	wallets, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, wallet := range wallets {
		if _, err := sdk.AccAddressFromBech32(wallet); err != nil {
			return fmt.Errorf("invalid Bech32 address: %s", wallet)
		}
	}
	return nil
}

// Helper to get inflation rate as LegacyDec
func (p Params) GetInflationRateAsDec() (math.LegacyDec, error) {
	return math.LegacyNewDecFromStr(p.InflationRate)
}
