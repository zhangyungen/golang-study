package main

import (
	"fmt"
	"github.com/google/uuid"
	"time"
	"zyj.com/golang-study/util/str"
	timeUtil "zyj.com/golang-study/util/time"
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

	fmt.Println("=== Go 语言位运算 ===\n")
	fmt.Println("TokenTypeBitRefresh:", TokenTypeBitRefresh)
	fmt.Println("TokenTypeBitWeb:", TokenTypeBitWeb)
	fmt.Println("TypeBitPixCake:", TypeBitPixCake)
	fmt.Println("TypeBitPixSoda:", TypeBitPixSoda)
	fmt.Println("TokenTypeBitPixIpad:", TokenTypeBitPixIpad)
	fmt.Println("TokenTypeSodaLogin:", TokenTypeSodaLogin)
	println(timeUtil.AddDays(time.Now(), 1).Unix())
	println(str.CountMatches("test234r4", "test"))
	println(str.StringsStartWith([]string{"test234r4", "estfdsafas"}, "est"))
	println(uuid.New().String())

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
