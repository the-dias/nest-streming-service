package cache

import (
	"log"
	"nats-service/app/cache"
	"nats-service/app/config"
	"nats-service/app/model"
	utils "nats-service/app/utils/sql"
)

func GetDataToCache(conf *config.Config, cacheHolder *cache.Cache) {

	if cacheHolder.Len() == 0 {
		var orderCache map[int]model.Order

		orderCache, err := utils.RetrieveOrdersFromDB(
			conf.DatabaseConfig.User,
			conf.DatabaseConfig.Password,
			conf.DatabaseConfig.DatabaseName,
			conf.DatabaseConfig.Host,
			conf.DatabaseConfig.Port,
		)
		// utils.
		for i, item := range orderCache {
			cacheHolder.Set(i, item)
		}

		if err != nil {
			log.Printf("Erro caching data from postgresql: %v", err)
		}
	}
}
