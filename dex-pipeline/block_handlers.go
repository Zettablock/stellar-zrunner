package main

import (
	"time"
	"fmt"
	"strings"
	"regexp"
	"strconv"

	"github.com/Zettablock/zsource/utils"
    "github.com/Zettablock/stellar-zrunner/dao"
	"github.com/bitly/go-simplejson"
	"github.com/Zettablock/stellar-zrunner/util"
)

const (
	SoroswapPoolsContract = "CA4HEQTL2WPEUYKYKCDOHCDNIV4QHNJ7EL4J4NQ6VADP7SYHVRYZ7AW2"
	PhoenixPoolsContract = "CB4SVAWJA6TSRNOJZ7W2AWFW46D5VR4ZMFZKDIKXEINZCZEGZCJZCKMI"
	//PhoenixActionContract = "CBISULYO5ZGS32WTNCBMEFCNKNSLFXCQ4Z3XHVDP4X4FLPSEALGSY3PS"
	BlendLendingPoolContract = "CCZD6ESMOGMPWH2KRO4O7RGTAPGTUPFWFQBELQSS7ZUK63V3TZWETGAG"
	 
)
/*var (
	SoroswapLiquidityActionContracts = []string{
		"CACTIOUW5FHYD3Q6ENKAU2IBLO2YFRWST4OGPDB4H32OGFMMJQF6SAJ5",
		"CCH3CJZWG6UMW522ESP3UHL4DCZLNXZLUHKYG5GCGNG5HXRL4A6O4A23",
		"CATUJXDUO7SSSTAKSUV5YU6RSTB4B5AVIHQDV26QTCXOB46T6SLMWNMY",
		"CABIXKWFCRM6VYUPNKF5C24O5LRCINE4XGB7SRUK67T6EVYMQSRUROKH",
		"CDQ4UKVWHJKR465B3NN2YP3IMWBEZ77YYJYTHYZA3BWNTPRHFOJ4OY57",
		"CCXOKQBBNRJ7YKY4Y6HTXAU5ZLV7PKENA7ZT74UAAXZ6XUEE22YQGBCS",
		"CB63RYTOXPVHXJAM7BGN7AUKFGHUCT5KVCK4H2BV4CHPEWIV4J3WDY3W",
	}
	BlendLendingActionContracts = []string{
		"CBP7NO6F7FRDHSOFQBT2L2UWYIZ2PU76JKVRYAQTG3KZSQLYAOKIF2WB",
		"CAO3AGAMZVRMHITL36EJ2VZQWKYRPWMQAPDQD5YEOF3GIF7T44U4JAL3",
	}
)*/


func HandleBlock(blockNumber int64, deps *utils.Deps) (bool, error) {
	//deps.Logger.Info("Sui Dex BlockHandler", "block number", blockNumber)

	var eventDao = dao.NewEventDao(deps.Config.Source.Schema, deps.SourceDB)
	var evts []dao.Event = eventDao.List(blockNumber)

	pools := make(map[string] dao.StellarDexPool)
	liquidityActions := make(map[string] dao.StellarDexLiquidityAction)
	swaps := make(map[string] dao.StellarDexSwap)
	lendingPools := make(map[string] dao.StellarLendingPool)
	lendingActions := make(map[string] dao.StellarLendingAction)
	for _, evt := range evts {
		eventType := parseEventType(&evt, deps)
		if eventType.IsSoroswapPoolEvent {
			data, err := parseSoroswapPool(&evt)
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("parse SoroswapPool error: %v", err))
				continue
			}
			data.FactoryContractID = evt.ContractId
			pools[evt.Id] = data
		} else if eventType.IsPhoenixPoolEvent {
			data, err := parsePhoenixPool(&evt)
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("parse PhoenixPool error: %v", err))
				continue
			}
			if val, ok := pools[data.TransactionHash]; ok {
				val.PoolContractID = data.PoolContractID
				val.FactoryContractID = evt.ContractId
				pools[data.TransactionHash] = val
			} else {
				data.FactoryContractID = evt.ContractId
				pools[data.TransactionHash] = data
			}
		} else if eventType.IsPhoenixPoolDetailsEvent {
			data, err := parsePhoenixPoolDetails(&evt)
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("parse PhoenixPool details error: %v", err))
				continue
			}
			if val, ok := pools[data.TransactionHash]; ok {
				data.PoolContractID = val.PoolContractID
				pools[data.TransactionHash] = data
			} else {
				pools[data.TransactionHash] = data
			}
			//fmt.Println(pools[data.TransactionHash])
		} else if eventType.IsSoroswapLiquidityActionsEvent {
			data, err := parseSoroswapLiquidityAction(&evt)
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("parse SoroswapLiquidityActions error: %v", err))
				continue
			}
			liquidityActions[evt.Id] = data
		} else if eventType.IsPhoenixLiquidityActionsEvent {
			var data dao.StellarDexLiquidityAction
			if _, ok := liquidityActions[evt.TransactionHash]; ok {
				data = liquidityActions[evt.TransactionHash]
			} else {
				data = parseLiquidityAction(&evt)
				data.ActionType = evt.Topic[0]
			}
			err := parsesPhoenixLiquidityAction(&evt, &data)
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("parse SoroswapLiquidityActions error: %v", err))
				continue
			}
			liquidityActions[data.TransactionHash] = data
		} else if eventType.IsSoroswapSwapEvent {
			data, err := parseSoroswapSwap(&evt)
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("parse SoroswapSwap error: %v", err))
				continue
			}
			swaps[evt.Id] = data
		} else if eventType.IsPhoenixSwapEvent {
			var data dao.StellarDexSwap
			if _, ok := swaps[evt.TransactionHash]; ok {
				data = swaps[evt.TransactionHash]
			} else {
				data = parseSwap(&evt)
			}
			err := parsesPhoenixSwap(&evt, &data)
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("parse PhoenixSwap error: %v", err))
				continue
			}
			swaps[data.TransactionHash] = data
		} else if eventType.IsBlendLendingPoolEvent {
			data, err := parseBlendLendingPool(&evt)
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("parse BlendLendingPool error: %v", err))
				continue
			}
			lendingPools[evt.Id] = data
		} else if eventType.IsBlendLendingActionEvent {
			var data dao.StellarLendingAction
			if _, ok := lendingActions[evt.Id]; ok {
				data = lendingActions[evt.Id]
			} else {
				data = parseLendingAction(&evt)
			}
			err := parseBlendLendingAction(&evt, &data)
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("parse BlendLendingAction error: %v", err))
				continue
			}
			lendingActions[evt.Id] = data
		}
	}
	for _, pool := range pools {
		if pool.PoolContractID != "" {
			err := deps.DestinationDB.Table(fmt.Sprintf("%s.%s", deps.DestinationDBSchema, pool.TableName())).Save(&pool).Error
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("Stellar Dex failed to save pool: %v", err))
				continue
			} else {
				if pool.FactoryContractID == PhoenixPoolsContract {
					util.PhoenixPoolsCache.Set(pool.PoolContractID, 1)
				}
			}
		}
	}
	for _, action := range liquidityActions {
		if action.PoolContractID != "" {
			updateLiquidityActionTokenInfo(&action, deps)
			err := deps.DestinationDB.Table(fmt.Sprintf("%s.%s", deps.DestinationDBSchema, action.TableName())).Save(&action).Error
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("Stellar Dex failed to save liquidityAction: %v", err))
				continue
			}
		}
	}
	for _, swap := range swaps {
		if swap.PoolContractID != "" {
			updateSwapTokenInfo(&swap, deps)
			//fmt.Println(swap)
			err := deps.DestinationDB.Table(fmt.Sprintf("%s.%s", deps.DestinationDBSchema, swap.TableName())).Save(&swap).Error
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("Stellar Dex failed to save swap: %v", err))
				continue
			}
		}
	}
	for _, lendingPool := range lendingPools {
		if lendingPool.PoolContractID != "" {
			err := deps.DestinationDB.Table(fmt.Sprintf("%s.%s", deps.DestinationDBSchema, lendingPool.TableName())).Save(&lendingPool).Error
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("Stellar Dex failed to save lendingPool: %v", err))
				continue
			} else {
				if lendingPool.FactoryContractID == BlendLendingPoolContract {
					util.BlendLendingPoolsCache.Set(lendingPool.PoolContractID, 1)
				}
			}
		}
	}
	for _, lendingAction := range lendingActions {
		if lendingAction.PoolContractID != "" {
			err := deps.DestinationDB.Table(fmt.Sprintf("%s.%s", deps.DestinationDBSchema, lendingAction.TableName())).Save(&lendingAction).Error
			if err != nil {
				deps.Logger.Error(fmt.Sprintf("Stellar Dex failed to save lendingAction: %v", err))
				continue
			}
		}
	}
	return false, nil
}

func parseEventType(evt *dao.Event, deps *utils.Deps) *stellarEventType {
	var et stellarEventType
	if len(evt.Topic) >= 1 {
		if evt.Topic[0] == "fn_call" {
			if len(evt.Topic) == 3 && evt.Topic[1] == PhoenixPoolsContract  && evt.Topic[2] == "create_liquidity_pool" {
				et.IsPhoenixPoolDetailsEvent = true
			} else if len(evt.Topic) == 3 && evt.ContractId == BlendLendingPoolContract && evt.Topic[2] == "initialize" {
				et.IsBlendLendingPoolEvent = true
			}
		} else if evt.Topic[0] == "SoroswapFactory" {
			if len(evt.Topic) == 2 && evt.Topic[1] == "new_pair" && evt.ContractId == SoroswapPoolsContract {
				et.IsSoroswapPoolEvent = true
			}
		} else if evt.Topic[0] == "create" {
			if len(evt.Topic) == 2 && evt.Topic[1] == "liquidity_pool" && evt.ContractId == PhoenixPoolsContract {
				et.IsPhoenixPoolEvent = true
			}
		} else if evt.Topic[0] == "SoroswapPair" {
			if len(evt.Topic) == 2 {
				if evt.Topic[1] == "swap" {
					et.IsSoroswapSwapEvent = true
				} else if arrContains([]string{"deposit", "withdraw"}, evt.Topic[1]) {
					et.IsSoroswapLiquidityActionsEvent = true
				}
			}
		} else if evt.Topic[0] == "deploy" {
			if len(evt.Topic) == 2 {
				if evt.Topic[1] == "swap" {
					et.IsSoroswapSwapEvent = true
				} else if arrContains([]string{"deposit", "withdraw"}, evt.Topic[1]) {
					et.IsSoroswapLiquidityActionsEvent = true
				}
			}
		} else {
			pt := checkPoolContract(evt.ContractId, deps)
			if pt.IsPhoenixPool {
				if arrContains([]string{"provide_liquidity", "withdraw_liquidity"}, evt.Topic[0]) {
					et.IsPhoenixLiquidityActionsEvent = true
				} else if evt.Topic[0] == "swap" {
					et.IsPhoenixSwapEvent = true
				}
			} else if pt.IsBlendLendingPool {
				if arrContains([]string{"supply", "withdraw", "borrow", "repay", "supply_collateral", "withdraw_collateral", "fill_auction"}, evt.Topic[0]) {
					et.IsBlendLendingActionEvent = true
				}
			}
		}
	}
	return &et
}

func extractAddressValues(s string) string {
	re := regexp.MustCompile(`<Address \[type=([A-Z]+), address=([A-Z0-9]+)\]>`)
	matches := re.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		if len(match) == 3 {
			a := match[1]
			b := match[2]
			newStr := fmt.Sprintf(`{"Address": { "type": "%s", "address": "%s"}}`, a, b)
			s = strings.Replace(s, match[0], newStr, -1)
		}
	}
	s = strings.Replace(s, `'`, `"`, -1)
	return s
}
func extractTokensValue(s string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)
	re := regexp.MustCompile(`'token_(\w+)': <Address \[type=([A-Z]+), address=([A-Z0-9]+)\]>`)

	matches := re.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		if len(match) == 4 {
			token := match[1]
			a := match[2]
			b := match[3]

			if _, ok := result[token]; !ok {
				result[token] = make(map[string]string)
			}
			result[token]["type"] = a
			result[token]["address"] = b
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("No matching tokens found in the event value")
	}
	return result, nil
}

func parseSoroswapPool(evt *dao.Event) (dao.StellarDexPool, error) {
	data := parsePool(evt)
	var evtValue = extractAddressValues(evt.Value)
	rawData, err := simplejson.NewJson([]byte(evtValue))
	if err != nil {
		return data, err
	}else {
		data.PoolContractID = rawData.Get("pair").Get("Address").Get("address").MustString()
		data.TokenAType = rawData.Get("token_0").Get("Address").Get("type").MustString()
		data.TokenAAccount = rawData.Get("token_0").Get("Address").Get("address").MustString()
		data.TokenBType = rawData.Get("token_1").Get("Address").Get("type").MustString()
		data.TokenBAccount = rawData.Get("token_1").Get("Address").Get("address").MustString()
		data.ParsedJSON = evtValue
	}
	return data, nil
}
func parsePhoenixPool(evt *dao.Event) (dao.StellarDexPool, error) {
	var data dao.StellarDexPool
	var evtValue = extractAddressValues(evt.Value)
	rawData, err := simplejson.NewJson([]byte(evtValue))
	if err != nil {
		return data, err
	}else {
		data.PoolContractID = rawData.Get("Address").Get("address").MustString()
		data.ParsedJSON = evtValue
		data.EventID = evt.Id
		data.TransactionHash = evt.TransactionHash
	}
	return data, nil
}
func parsePhoenixPoolDetails(evt *dao.Event) (dao.StellarDexPool, error) {
	data := parsePool(evt)
	var evtValue = extractAddressValues(evt.Value)

	result, err := extractTokensValue(evt.Value)
	if err != nil {
		return data, err
	}

	for token, values := range result {
		if token == "a" {
			data.TokenAType = values["type"]
			data.TokenAAccount = values["address"]
		} else if token == "b" {
			data.TokenBType = values["type"]
			data.TokenBAccount = values["address"]
		}
	}
	data.ParsedJSON = evtValue

	return data, nil
}
func parsePool(evt *dao.Event) dao.StellarDexPool {
	var data dao.StellarDexPool
	data = dao.StellarDexPool{
		PoolContractID:		"",
		TokenAType:			"",
		TokenAAccount:		"",
		TokenBType:			"",
		TokenBAccount:		"",
		FactoryContractID:	"",
		ParsedJSON:			"",
		EventID:			evt.Id,
		Ledger:				evt.Ledger,
		LedgerClosedAt:		evt.LedgerClosedAt,
		Topic:				evt.Topic,
		Value:				evt.Value,
		TransactionHash:	evt.TransactionHash,
		ProcessTime:		time.Now(),
		BlockDate:			evt.BlockDate,
	}
	return data
}

func updateSwapTokenInfo(data *dao.StellarDexSwap, deps *utils.Deps) {
	if data.TokenAAccount == "" {
		var pools []dao.StellarDexPool
		var p dao.StellarDexPool
		deps.DestinationDB.Table(p.TableName()).Where("pool_contract_id = ?", data.PoolContractID).Limit(1).Scan(&pools)
		if len(pools) > 0 {
			data.TokenAType = pools[0].TokenAType
			data.TokenAAccount = pools[0].TokenAAccount
			data.TokenBType = pools[0].TokenBType
			data.TokenBAccount = pools[0].TokenBAccount
		}
	}
}
func updateLiquidityActionTokenInfo(data *dao.StellarDexLiquidityAction, deps *utils.Deps) {
	if data.TokenAAccount == "" {
		var pools []dao.StellarDexPool
		var p dao.StellarDexPool
		deps.DestinationDB.Table(p.TableName()).Where("pool_contract_id = ?", data.PoolContractID).Limit(1).Scan(&pools)
		if len(pools) > 0 {
			data.TokenAType = pools[0].TokenAType
			data.TokenAAccount = pools[0].TokenAAccount
			data.TokenBType = pools[0].TokenBType
			data.TokenBAccount = pools[0].TokenBAccount
		}
	}
}
func parseSoroswapLiquidityAction(evt *dao.Event) (dao.StellarDexLiquidityAction, error) {
	data := parseLiquidityAction(evt)
	var evtValue = extractAddressValues(evt.Value)
	rawData, err := simplejson.NewJson([]byte(evtValue))
	if err != nil {
		return data, err
	}else {
		data.ActionType = evt.Topic[1]
		data.UserType = rawData.Get("to").Get("Address").Get("type").MustString()
		data.UserAccount = rawData.Get("to").Get("Address").Get("address").MustString()
		data.Amount0 = strconv.FormatInt(rawData.Get("amount_0").MustInt64(),10)
		data.Amount1 = strconv.FormatInt(rawData.Get("amount_1").MustInt64(),10)
		data.Liquidity = strconv.FormatInt(rawData.Get("liquidity").MustInt64(),10)
		data.NewReserve0 = strconv.FormatInt(rawData.Get("new_reserve_0").MustInt64(),10)
		data.NewReserve1 = strconv.FormatInt(rawData.Get("new_reserve_1").MustInt64(),10)
		data.ParsedJSON = evtValue
	}
	return data, nil
}
func parsesPhoenixLiquidityAction(evt *dao.Event, data *dao.StellarDexLiquidityAction) error {
	var addressArr = []string{"sender", "token_a", "token_b"}
	var jsonStr = ""
	if arrContains(addressArr, evt.Topic[1]) {
		var evtValue = extractAddressValues(evt.Value)
		rawData, err := simplejson.NewJson([]byte(evtValue))
		if err != nil {
			return err
		} else {
			if evt.Topic[1] == "token_a" {
				data.TokenAType = rawData.Get("Address").Get("type").MustString()
				data.TokenAAccount = rawData.Get("Address").Get("address").MustString()
			} else if evt.Topic[1] == "token_b" {
				data.TokenBType = rawData.Get("Address").Get("type").MustString()
				data.TokenBAccount = rawData.Get("Address").Get("address").MustString()
			} else if evt.Topic[1] == "sender" {
				data.UserType = rawData.Get("Address").Get("type").MustString()
				data.UserAccount = rawData.Get("Address").Get("address").MustString()
			}
			jsonStr = evtValue
		}
	} else {
		if evt.Topic[1] == "token_a-amount" {
			data.Amount0 = evt.Value
		} else if evt.Topic[1] == "token_b-amount" {
			data.Amount1 = evt.Value
		} else if evt.Topic[1] == "return_amount_a" {
			data.NewReserve0 = evt.Value
		} else if evt.Topic[1] == "return_amount_b" {
			data.NewReserve1 = evt.Value
		} else if evt.Topic[1] == "shares_amount" {
			data.Liquidity = evt.Value
		}
		jsonStr = evt.Value
	}
	jsonStr = `"` + evt.Topic[1] + `":` + jsonStr
	if data.ParsedJSON == "" {
		data.ParsedJSON = "{" + jsonStr + "}"
	} else {
		data.ParsedJSON = data.ParsedJSON[:len(data.ParsedJSON)-1] + "," + jsonStr + "}"
	}
	return nil
}
func parseLiquidityAction(evt *dao.Event) dao.StellarDexLiquidityAction {
	var data dao.StellarDexLiquidityAction
	data = dao.StellarDexLiquidityAction{
		EventID:			evt.Id,
		PoolContractID:		evt.ContractId,
		UserType:			"",
		UserAccount:		"",
		ActionType:			"",
		TokenAType:			"",
		TokenAAccount:		"",
		TokenBType:			"",
		TokenBAccount:		"",
		Amount0:			"",
		Amount1:			"",
		Liquidity:			"",
		NewReserve0:		"",
		NewReserve1:		"",
		ParsedJSON:			"",
		Ledger:				evt.Ledger,
		LedgerClosedAt:		evt.LedgerClosedAt,
		Topic:				evt.Topic,
		Value:				evt.Value,
		TransactionHash:	evt.TransactionHash,
		ProcessTime:		time.Now(),
		BlockDate:			evt.BlockDate,
	}
	return data
}
func parseSoroswapSwap(evt *dao.Event) (dao.StellarDexSwap, error) {
	data := parseSwap(evt)
	var evtValue = extractAddressValues(evt.Value)
	rawData, err := simplejson.NewJson([]byte(evtValue))
	if err != nil {
		return data, err
	}else {
		data.UserType = rawData.Get("to").Get("Address").Get("type").MustString()
		data.UserAccount = rawData.Get("to").Get("Address").Get("address").MustString()
		data.Amount0In = strconv.FormatInt(rawData.Get("amount_0_in").MustInt64(),10)
		data.Amount0Out = strconv.FormatInt(rawData.Get("amount_0_out").MustInt64(),10)
		data.Amount1In = strconv.FormatInt(rawData.Get("amount_1_in").MustInt64(),10)
		data.Amount1Out = strconv.FormatInt(rawData.Get("amount_1_out").MustInt64(),10)
		data.ParsedJSON = evtValue
	}
	return data, nil
}
func parsesPhoenixSwap(evt *dao.Event, data *dao.StellarDexSwap) error {
	var addressArr = []string{"sender", "sell_token", "buy_token"}
	var jsonStr = ""
	if arrContains(addressArr, evt.Topic[1]) {
		var evtValue = extractAddressValues(evt.Value)
		rawData, err := simplejson.NewJson([]byte(evtValue))
		if err != nil {
			return err
		} else {
			if evt.Topic[1] == "sell_token" {
				data.TokenAType = rawData.Get("Address").Get("type").MustString()
				data.TokenAAccount = rawData.Get("Address").Get("address").MustString()
			} else if evt.Topic[1] == "buy_token" {
				data.TokenBType = rawData.Get("Address").Get("type").MustString()
				data.TokenBAccount = rawData.Get("Address").Get("address").MustString()
			} else if evt.Topic[1] == "sender" {
				data.UserType = rawData.Get("Address").Get("type").MustString()
				data.UserAccount = rawData.Get("Address").Get("address").MustString()
			}
			jsonStr = evtValue
		}
	} else {
		if evt.Topic[1] == "offer_amount" {
			data.Amount0Out = evt.Value
		} else if evt.Topic[1] == "return_amount" {
			data.Amount1In = evt.Value
		} else if evt.Topic[1] == "spread_amount" {
			num, _ := strconv.ParseInt(evt.Value, 10, 32)
			data.SpreadAmount = int32(num)
		} else if evt.Topic[1] == "referral_fee_amount" {
			num, _ := strconv.ParseInt(evt.Value, 10, 32)
			data.ReferralFeeAmount = int32(num)
		}
		jsonStr = evt.Value
	}
	jsonStr = `"` + evt.Topic[1] + `":` + jsonStr
	if data.ParsedJSON == "" {
		data.ParsedJSON = "{" + jsonStr + "}"
	} else {
		data.ParsedJSON = data.ParsedJSON[:len(data.ParsedJSON)-1] + "," + jsonStr + "}"
	}
	return nil
}
func parseSwap(evt *dao.Event) dao.StellarDexSwap {
	var data dao.StellarDexSwap
	data = dao.StellarDexSwap{
		EventID:			evt.Id,
		PoolContractID:		evt.ContractId,
		UserType:			"",
		UserAccount:		"",
		TokenAType:			"",
		TokenAAccount:		"",
		TokenBType:			"",
		TokenBAccount:		"",
		Amount0In:			"0",
		Amount1In:			"0",
		Amount0Out:			"0",
		Amount1Out:			"0",
		SpreadAmount:		0,
		ReferralFeeAmount:	0,
		ParsedJSON:			"",
		Ledger:				evt.Ledger,
		LedgerClosedAt:		evt.LedgerClosedAt,
		Topic:				evt.Topic,
		Value:				evt.Value,
		TransactionHash:	evt.TransactionHash,
		ProcessTime:		time.Now(),
		BlockDate:			evt.BlockDate,
	}
	return data
}

func parseBlendLendingPool(evt *dao.Event) (dao.StellarLendingPool, error) {
	data := parseLendingPool(evt)
	var evtValue = extractAddressValues(evt.Value)
	rawData, err := simplejson.NewJson([]byte(evtValue))
	if err != nil {
		return data, err
	}else {
		l := len(rawData.MustArray())
		if l > 1 {
			data.PoolName = rawData.GetIndex(1).MustString()
		}
		data.ParsedJSON = evtValue
	}
	return data, nil
}
func parseLendingPool(evt *dao.Event) dao.StellarLendingPool {
	var data dao.StellarLendingPool
	data = dao.StellarLendingPool{
		EventID:			evt.Id,
		PoolContractID:		evt.Topic[1],
		PoolName:			"",
		FactoryContractID:	evt.ContractId,
		ParsedJSON:			"",
		Ledger:				evt.Ledger,
		LedgerClosedAt:		evt.LedgerClosedAt,
		Topic:				evt.Topic,
		Value:				evt.Value,
		TransactionHash:	evt.TransactionHash,
		ProcessTime:		time.Now(),
		BlockDate:			evt.BlockDate,
	}
	return data
}
func parseBlendLendingAction(evt *dao.Event, data *dao.StellarLendingAction) error {
	var evtValue = extractAddressValues(evt.Value)
	rawData, err := simplejson.NewJson([]byte(evtValue))
	if err != nil {
		return err
	}else {
		var contract = extractAddressValues(evt.Topic[1])
		rawDataContract, errContract := simplejson.NewJson([]byte(contract))
		if errContract != nil {
			return errContract
		} else {
			data.TokenType = rawDataContract.Get("Address").Get("type").MustString()
			data.TokenAccount = rawDataContract.Get("Address").Get("address").MustString()
		}
		if evt.Topic[0] == "fill_auction" {
			l := len(rawData.MustArray())
			if l > 1 {
				data.RequestAmount = strconv.FormatInt(rawData.GetIndex(1).MustInt64(),10)
				data.UserType = rawData.GetIndex(0).Get("Address").Get("type").MustString()
				data.UserAccount = rawData.GetIndex(0).Get("Address").Get("address").MustString()
			}
		} else {
			var account = extractAddressValues(evt.Topic[2])
			rawDataAccount, errAccount := simplejson.NewJson([]byte(account))
			if errAccount != nil {
				return errAccount
			} else {
				data.UserType = rawDataAccount.Get("Address").Get("type").MustString()
				data.UserAccount = rawDataAccount.Get("Address").Get("address").MustString()
			}
			l := len(rawData.MustArray())
			if l > 1 {
				data.RequestAmount = strconv.FormatInt(rawData.GetIndex(0).MustInt64(),10)
				if arrContains([]string{"supply", "supply_collateral"}, evt.Topic[0]) {
					data.BtokenAmount = strconv.FormatInt(rawData.GetIndex(1).MustInt64(),10)
				} else if arrContains([]string{"borrow", "repay"}, evt.Topic[0]) {
					data.DtokenAmount = strconv.FormatInt(rawData.GetIndex(1).MustInt64(),10)
				}
			}
		}
		jsonStr := `"` + evt.Topic[0] + `":` + evtValue
		if data.ParsedJSON == "" {
			data.ParsedJSON = "{" + jsonStr + "}"
		} else {
			data.ParsedJSON = data.ParsedJSON[:len(data.ParsedJSON)-1] + "," + jsonStr + "}"
		}
	}
	return nil
}
func parseLendingAction(evt *dao.Event) dao.StellarLendingAction {
	var data dao.StellarLendingAction
	data = dao.StellarLendingAction{
		EventID:			evt.Id,
		PoolContractID:		evt.ContractId,
		UserType:			"",
		UserAccount:		"",
		ActionType:			evt.Topic[0],
		TokenType:			"",
		TokenAccount:		"",
		RequestAmount:		"",
		BtokenAmount:		"",
		BtokenType:			"",
		BtokenAccount:		"",
		DtokenAmount:		"",
		DtokenType:			"",
		DtokenAccount:		"",
		ParsedJSON:			"",
		Ledger:				evt.Ledger,
		LedgerClosedAt:		evt.LedgerClosedAt,
		Topic:				evt.Topic,
		Value:				evt.Value,
		TransactionHash:	evt.TransactionHash,
		ProcessTime:		time.Now(),
		BlockDate:			evt.BlockDate,
	}
	return data
}

type stellarEventType struct {
	IsSoroswapPoolEvent				bool
	IsPhoenixPoolEvent				bool
	IsPhoenixPoolDetailsEvent		bool
	IsSoroswapLiquidityActionsEvent	bool
	IsPhoenixLiquidityActionsEvent	bool
	IsSoroswapSwapEvent				bool
	IsPhoenixSwapEvent				bool
	IsBlendLendingPoolEvent			bool
	IsBlendLendingActionEvent		bool
}
type poolType struct {
	IsPhoenixPool				bool
	IsBlendLendingPool			bool
}


func initPoolContractCache(deps *utils.Deps) {
	if util.LoadedDb == false {
		util.LoadedDb = true
		var phoenixPools []dao.StellarDexPool
		var phoenixPool dao.StellarDexPool
		deps.DestinationDB.Table(phoenixPool.TableName()).Where("factory_contract_id = ?", PhoenixPoolsContract).Limit(util.CacheCapacity).Scan(&phoenixPools)
		if len(phoenixPools) > 0 {
			for _, pp := range phoenixPools {
				util.PhoenixPoolsCache.Set(pp.PoolContractID, 1)
			}
		}

		var lendingPools []dao.StellarLendingPool
		var lendingPool dao.StellarLendingPool
		deps.DestinationDB.Table(lendingPool.TableName()).Where("factory_contract_id = ?", BlendLendingPoolContract).Limit(util.CacheCapacity).Scan(&lendingPools)
		if len(lendingPools) > 0 {
			for _, pp2 := range lendingPools {
				util.BlendLendingPoolsCache.Set(pp2.PoolContractID, 1)
			}
		}
	}
}
func checkPoolContract(contract string, deps *utils.Deps) *poolType {
	initPoolContractCache(deps)
	var pt poolType
	_, ok := util.PhoenixPoolsCache.Get(contract)
	if ok {
		pt.IsPhoenixPool = true
		return &pt
	}
	_, ok = util.BlendLendingPoolsCache.Get(contract)
	if ok {
		pt.IsBlendLendingPool = true
		return &pt
	}
	var pools []dao.StellarDexPool
	var p dao.StellarDexPool
	deps.DestinationDB.Table(p.TableName()).Where("pool_contract_id = ?", contract).Limit(1).Scan(&pools)
	if len(pools) > 0 {
		if pools[0].FactoryContractID == PhoenixPoolsContract{
			pt.IsPhoenixPool = true
			util.PhoenixPoolsCache.Set(contract, 1)
			return &pt
		}
	}

	var poolsLending []dao.StellarLendingPool
	var pLending dao.StellarLendingPool
	deps.DestinationDB.Table(pLending.TableName()).Where("pool_contract_id = ?", contract).Limit(1).Scan(&poolsLending)
	if len(poolsLending) > 0 {
		if poolsLending[0].FactoryContractID == BlendLendingPoolContract{
			pt.IsBlendLendingPool = true
			util.BlendLendingPoolsCache.Set(contract, 1)
			return &pt
		}
	}
	return &pt
}

func arrContains(arr []string, s string) bool {
	for _, element := range arr {
		if element == s {
			return true
		}
	}
	return false
}