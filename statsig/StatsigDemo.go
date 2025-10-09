package main

import (
	"fmt"
	statsig "github.com/statsig-io/go-sdk"
)

func main() {
	// 初始化
	statsig.Initialize("server-secret-key")
	//初始化，可按配置初始化
	//type Options struct {
	//	API                string      `json:"api"`
	//	Environment        Environment `json:"environment"` production、staging
	//	LocalMode          bool        `json:"localMode"`    default false 数据是否外发配置， 不外发情况将存储到缓存中（需定制）
	//	ConfigSyncInterval time.Duration  默认10s 配置同步间隔
	//	IDListSyncInterval time.Duration	默认1分钟ID列表同步间隔
	//  BootstrapValues      string	  默认nil
	//	RulesUpdatedCallback func(rules string, time int64)   规则更新回调
	//  UserStorage	 UserPersistentStorage  // 用户存储
	//	DisableIdListDisableIdList	DisableIdListDisableIdList // 是否禁用ID列表在网络和数据适配器的初始化和后台轮询期间禁用获取id列表的标志。
	//}
	//
	// Or, if you want to initialize with certain options
	statsig.InitializeWithOptions("server-secret-key", &statsig.Options{})

	//Checking a Feature Flag/Gate
	user := statsig.User{UserID: "some_user_id"}
	featureFlg := statsig.CheckGate(user, "use_new_feature")
	if featureFlg {
		// Gate is on, enable new feature
	} else {
		// Gate is off
	}

	//Reading a Dynamic Config
	config := statsig.GetConfig(user, "awesome_product_details")
	// The 2nd parameter is the default value to be used in case the given parameter name does not exist on
	// the Dynamic Config object. This can happen when there is a typo, or when the user is offline and the
	// value has not been cached on the client.
	var itemName = config.GetString("product_name", "Awesome Product v1")
	var price = config.GetNumber("price", 10.0)
	var shouldDiscount = config.GetBool("discount", false) // Or just get the whole json object backing this config if you prefer
	var json = config.Value                                // Values via getLayer
	layer := statsig.GetLayer(user, "user_promo_experiments")
	var promoTitle = layer.GetString("title", "Welcome to Statsig!")
	discount := layer.GetNumber("discount", 0.1)

	// or, via getExperiment
	//Getting a Layer/Experiment
	titleExperiment := statsig.GetExperiment(user, "new_user_promo_title")
	priceExperiment := statsig.GetExperiment(user, "new_user_promo_price")

	var promoTitle = titleExperiment.GetString("title", "Welcome to Statsig!")
	var discount = priceExperiment.GetNumber("discount", 0.1)
	var price = 1 * (1 - discount) // getting the layer name that an experiment belongs to

	var userPromoLayer = statsig.GetExperimentLayer("new_user_promo_title") //Logging an Event
	statsig.LogEvent(statsig.Event{
		User:      user,
		EventName: "add_to_cart",
		Value:     "SKU_12345",
		Metadata:  map[string]string{"price": "9.99", "item_name": "diet_coke_48_pack"},
	})
	//Retrieving Feature Gate Metadata
	user = statsig.User{UserID: "some_user_id"}
	feature := statsig.GetGate(user, "use_new_feature")
	if feature.Value {
		// Gate is on, enable new feature
		fmt.Print(feature.EvaluationDetails.Reason)
	}

}
