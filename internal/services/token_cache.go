package services

//func InitTokenCache(authList map[string] config.AuthUserConfig, keyConfig config.KeyConfig, integrationConfig *config.IntegrationConfig, cache *types.TokenCache) (*types.TokenCache, error) {
//	if cache == nil {
//		cache = &types.TokenCache{}
//		cache.TokenMap = make(map[string]*types.AuthToken)
//	}
//
//	cache.Lock()
//	for key,_ := range authList {
//
//		token := cache.TokenMap[key]
//		if token == nil {
//			token = &types.AuthToken{}
//		}
//
//		err := token.Login(authList[key], keyConfig, integrationConfig)
//		if err != nil {
//			return nil, fmt.Errorf("failed to create token cache: %s", err)
//		}
//
//		//TODO: uncomment this to get signature key
//		//err = token.GetSignatureKey(keyConfig, integrationConfig)
//		//if err != nil {
//		//	return nil, fmt.Errorf("failed to get signature key: %s", err)
//		//}
//
//		cache.TokenMap[key] = token
//	}
//
//	cache.Unlock()
//	return cache, nil
//}