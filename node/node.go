package node

import (
	"context"
	"database/sql"

	"github.com/Factom-Asset-Tokens/factom"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pegnet/pegnet/modules/grader"
	"github.com/pegnet/pegnetd/config"
	"github.com/pegnet/pegnetd/node/pegnet"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var OPRChain = *factom.NewBytes32FromString("a642a8674f46696cc47fdb6b65f9c87b2a19c5ea8123b3d2f0c13b6f33a9d5ef")
var TransactionChain = *factom.NewBytes32FromString("cffce0f409ebba4ed236d49d89c70e4bd1f1367d86402a3363366683265a242d")
var PegnetActivation uint32 = 206421
var GradingV2Activation uint32 = 210330

// TransactionConversionActivation indicates when tx/conversions go live on mainnet.
// Target Activation Height is Oct 7, 2019 15 UTC
var TransactionConversionActivation uint32 = 213237

// Estimated to be Oct 14 2019, 15:00:00 UTC
var PEGPricingActivation uint32 = 214287

type Pegnetd struct {
	FactomClient *factom.Client
	Config       *viper.Viper

	Sync   *pegnet.BlockSync
	Pegnet *pegnet.Pegnet
}

func NewPegnetd(ctx context.Context, conf *viper.Viper) (*Pegnetd, error) {
	// TODO : Update emyrk's factom library
	n := new(Pegnetd)
	n.FactomClient = FactomClientFromConfig(conf)
	n.Config = conf

	n.Pegnet = pegnet.New(conf)
	if err := n.Pegnet.Init(); err != nil {
		return nil, err
	}

	if sync, err := n.Pegnet.SelectSynced(ctx); err != nil {
		if err == sql.ErrNoRows {
			n.Sync = new(pegnet.BlockSync)
			n.Sync.Synced = PegnetActivation
			log.Debug("connected to a fresh database")
		} else {
			return nil, err
		}
	} else {
		n.Sync = sync
	}

	grader.InitLX()
	return n, nil
}

func FactomClientFromConfig(conf *viper.Viper) *factom.Client {
	cl := factom.NewClient()
	cl.FactomdServer = conf.GetString(config.Server)
	cl.WalletdServer = conf.GetString(config.Wallet)
	if config.WalletUser != "" {
		cl.Walletd.BasicAuth = true
		cl.Walletd.User = conf.GetString(config.WalletUser)
		cl.Walletd.Password = conf.GetString(config.WalletPass)
	}

	return cl
}
