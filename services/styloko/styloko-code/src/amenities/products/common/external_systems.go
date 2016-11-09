package common

import (
	"fmt"
	"github.com/jabong/floRest/src/common/utils/http"
	"github.com/jabong/floRest/src/common/utils/logger"
	"time"
)

type BoutiqueSystem struct {
	Url string
}

func (bs *BoutiqueSystem) SetUrl() error {
	conf := GetConfig()
	boutique := conf.Boutique.External
	bs.Url = boutique.Host + boutique.Path
	return nil
}

func (bs *BoutiqueSystem) UpdateProducts(productIds ...int) error {
	if len(productIds) == 0 {
		return nil
	}
	var url string = bs.Url + "/products/keys/"
	logger.Debug(fmt.Sprintf("Sending request to: [%s]", url))
	//prepare data
	configIds := make([]map[string]int, 0)
	for _, id := range productIds {
		configIds = append(configIds, map[string]int{"id": id})
	}
	data := M{
		"id_list": configIds,
	}

	response, err := http.HttpDelete(url, nil, data.ToJson(), 30*time.Second)
	if err != nil {
		return fmt.Errorf("(bs BoutiqueSystem)#UpdateMemcache Req failed: %s", err.Error())
	}
	logger.Debug(response.Body)
	if response.HttpStatus != 200 {
		return fmt.Errorf("(bs BoutiqueSystem)#UpdateMemcache Status: %s", response.HttpStatus)
	}
	return nil
}
