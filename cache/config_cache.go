package cache

import (
	"ehelp/o/config_system"
	"strconv"
)

var SystemCache = make(map[string]string)

func SetCacheSystem() error {
	var res, err = config_system.GetAllConfig()
	if err != nil {
		return err
	}
	if err == nil && len(res) > 0 {
		for _, item := range res {
			SystemCache[item.Name] = item.Value
		}
	}
	return nil
}

func convertToInt(cf string) int {
	var res, _ = strconv.Atoi(SystemCache[config_system.BONUS_TOOL_SERVICE_PERCENT])
	return res
}

func GetCacheBonusToolMoney() int {
	return convertToInt(config_system.BONUS_TOOL_SERVICE_MONEY)
}

func GetCacheHourMoney() int {
	var res = convertToInt(config_system.HOUR_MONEY)
	if res == 0 {
		return 50000
	}
	return res
}

func GetCacheBonusToolPercent() int {
	var res = convertToInt(config_system.BONUS_TOOL_SERVICE_PERCENT)
	if res == 0 {
		return 10
	}
	return res
}
