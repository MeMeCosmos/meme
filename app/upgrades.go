package app

import (
	"fmt"

	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	
	intertxtypes "github.com/cosmos/interchain-accounts/x/inter-tx/types"
)

const upgradeName = "v2.0.7"

func equalTraces(dtA, dtB ibctransfertypes.DenomTrace) bool {
	return dtA.BaseDenom == dtB.BaseDenom && dtA.Path == dtB.Path
}

func (app *MEMEApp) RegisterUpgradeHandlers(cfg module.Configurator) {
	app.upgradeKeeper.SetUpgradeHandler(upgradeName, func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
////
	vm[icatypes.ModuleName] = app.mm.Modules[icatypes.ModuleName].ConsensusVersion()
////
		
//
	var newTraces []ibctransfertypes.DenomTrace
		app.transferKeeper.IterateDenomTraces(ctx,
			func(dt ibctransfertypes.DenomTrace) bool {
				newTrace := ibctransfertypes.ParseDenomTrace(dt.GetFullDenomPath())
				if err := newTrace.Validate(); err == nil && !equalTraces(newTrace, dt) {
					newTraces = append(newTraces, newTrace)
				}
				return false
			})
		for _, nt := range newTraces {
			app.transferKeeper.SetDenomTrace(ctx, nt)
		}
//
		
////
	// initialize ICS27 module
	icamodule, correctTypecast := app.mm.Modules[icatypes.ModuleName].(ica.AppModule)
	if !correctTypecast {
		panic("mm.Modules[icatypes.ModuleName] is not of type ica.AppModule")
	}
	icamodule.InitModule(ctx, controllerParams, hostParams)
////
		
		return app.mm.RunMigrations(ctx, cfg, vm)

	})

	upgradeInfo, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == upgradeName && !app.upgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := store.StoreUpgrades{
			Added: []string{icacontrollertypes.StoreKey, icahosttypes.StoreKey, intertxtypes.StoreKey},
		}
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}
