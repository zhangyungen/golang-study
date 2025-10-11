package main

import (
	"log"
	"zyj.com/golang-study/util/obj"
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/result"
)

/*
TokenType 结构
<- 权限位 -> <- 开关位 ->
<- left  -> <-  8bit ->
*/
type TokenType uint64

const (
	TokenTypeBitRefresh     TokenType = 0x01
	TokenTypeBitWeb         TokenType = 0x02
	TypeBitPixCake          TokenType = 0x0200
	TypeBitPixSoda          TokenType = 0x0400
	TokenTypeBitPixIpad     TokenType = 0x0800
	TokenTypeBitAccount     TokenType = 0x1000
	TokenTypeBitOrgBackend  TokenType = 0x2000
	TokenTypeBitSalesforce  TokenType = 0x4000
	TokenTypeBitCheese      TokenType = 0x8000
	TokenTypeBitCakePhone   TokenType = 0x20000
	TokenTypeBitCakeAndroid TokenType = 0x40000
	TokenTypeBitToast       TokenType = 0x80000
	TokenTypeBitECommerce   TokenType = 0x100000

	TokenTypeSodaLogin    = TypeBitPixSoda | TokenTypeBitRefresh
	TokenTypeSodaWebLogin = TypeBitPixSoda | TokenTypeBitWeb | TokenTypeBitRefresh

	TokenTypeCheeseLogin       TokenType = TokenTypeBitCheese | TokenTypeBitRefresh
	TokenTypeCheeseWebLoginOld TokenType = TokenTypeBitCheese | TokenTypeBitWeb
	TokenTypeCheeseWebLogin    TokenType = TokenTypeBitCheese | TokenTypeBitWeb | TokenTypeBitRefresh

	TokenTypeToastLogin    TokenType = TokenTypeBitToast | TokenTypeBitRefresh
	TokenTypeToastWebLogin TokenType = TokenTypeBitToast | TokenTypeBitWeb | TokenTypeBitRefresh

	TokenTypeCakePhoneLogin   TokenType = TokenTypeBitCakePhone | TokenTypeBitRefresh
	TokenTypeCakeAndroidLogin TokenType = TokenTypeBitCakeAndroid | TokenTypeBitRefresh

	TokenTypeECommerceLogin    TokenType = TokenTypeBitECommerce | TokenTypeBitRefresh
	TokenTypeECommerceWebLogin TokenType = TokenTypeBitECommerce | TokenTypeBitWeb | TokenTypeBitRefresh

	TokenTypeBitAll = ^uint64(0)
)

func main() {

	var param1 = param.PageParam{Page: 1, PageSize: 10}
	var param2 = obj.ObjToObj[result.PageVO[model.User]](&param1)
	log.Println("param2 ObjToObj", obj.ObjToJsonStr(param2))

	param1.PageSize = 1000
	param1.Page = 11

	//var param4 = result.PageVO[model.User]{}
	//obj.CopierObj(&param1, &param4)
	//log.Println("CopierObj param4", obj.ObjToJsonStr(param4))

	obj2Map := obj.ObjToMap(&param1)
	log.Println("obj2Map", obj2Map)

	param1Ptr := obj.MapToObj[param.PageParam](obj2Map)
	log.Println("MapToObj", param1Ptr)

	json := obj.ObjToJsonStr(param1)
	var param3 = obj.JsonStrToObj[result.PageVO[model.User]](json)
	log.Println("ObjToJsonStr param3", obj.ObjToJsonStr(&param3))

	strs := make([]string, 0)
	strs = add(strs)
	strs = add(strs)
	strs = add(strs)
	log.Println("add strs", strs)
	//
	//fmt.Println("=== Go 语言位运算 ===\n")
	//fmt.Println("TokenTypeBitRefresh:", TokenTypeBitRefresh)
	//fmt.Println("TokenTypeBitWeb:", TokenTypeBitWeb)
	//fmt.Println("TypeBitPixCake:", TypeBitPixCake)
	//fmt.Println("TypeBitPixSoda:", TypeBitPixSoda)
	//fmt.Println("TokenTypeBitPixIpad:", TokenTypeBitPixIpad)
	//fmt.Println("TokenTypeSodaLogin:", TokenTypeSodaLogin)
	//println(timeUtil.AddDays(time.Now(), 1).Unix())
	//println(str.CountMatches("test234r4", "test"))
	//println(str.StringsStartWith([]string{"test234r4", "estfdsafas"}, "est"))
	//println(uuid.New().String())

	//logger, _ := zap.NewProduction()
	//defer logger.Sync() // flushes buffer, if any
	//logger.Info("failed to fetch URL",
	//	zap.String("url", "http://example.com"),
	//	zap.Int("attempt", 3),
	//	zap.Duration("backoff", time.Second),
	//)
	//lo
	//log := logger.NewProduction()
	//log.Infof("hello %s", "world")

}

func add(strs []string) []string {
	strs = append(strs, "hello")
	return strs
}
