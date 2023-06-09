package dcr

import (
	"fmt"
	"math"

	sharedW "github.com/crypto-power/cryptopower/libwallet/assets/wallet"
	"github.com/decred/dcrd/chaincfg/v3"
	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/decred/dcrd/wire"
)

const (
	// fetchPercentage is used to increase the initial estimate gotten during cfilters stage
	fetchPercentage = 0.38

	// Use 10% of estimated total headers fetch time to estimate rescan time
	rescanPercentage = 0.1

	// Use 80% of estimated total headers fetch time to estimate address discovery time
	discoveryPercentage = 0.8

	TestnetHDPath       = "m / 44' / 1' / "
	LegacyTestnetHDPath = "m / 44’ / 11’ / "
	MainnetHDPath       = "m / 44' / 42' / "
	LegacyMainnetHDPath = "m / 44’ / 20’ / "
)

// Returns a DCR amount that implements the asset amount interface.
func (asset *DCRAsset) ToAmount(v int64) sharedW.AssetAmount {
	return DCRAmount(dcrutil.Amount(v))
}

func AmountAtom(f float64) int64 {
	amount, err := dcrutil.NewAmount(f)
	if err != nil {
		log.Error(err)
		return -1
	}
	return int64(amount)
}

func calculateTotalTimeRemaining(timeRemainingInSeconds int64) string {
	minutes := timeRemainingInSeconds / 60
	if minutes > 0 {
		return fmt.Sprintf("%d min", minutes)
	}
	return fmt.Sprintf("%d sec", timeRemainingInSeconds)
}

func roundUp(n float64) int32 {
	return int32(math.Round(n))
}

// voteVersion was borrowed from upstream, and needs to always be in
// sync with the upstream method. This is the LOC to the upstream version:
// https://github.com/decred/dcrwallet/blob/master/wallet/wallet.go#L266
func voteVersion(params *chaincfg.Params) uint32 {
	switch params.Net {
	case wire.MainNet:
		return 9
	case 0x48e7a065: // TestNet2
		return 6
	case wire.TestNet3:
		return 10
	case wire.SimNet:
		return 10
	default:
		return 1
	}
}
