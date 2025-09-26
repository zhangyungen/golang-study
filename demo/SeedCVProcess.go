package main

import (
	"git.inner.truesightai.com/zyj1/ai-gc/httpclient"
	"git.inner.truesightai.com/zyj1/ai-gc/seedream"
	"log"
	"strconv"
	"time"
)

func main() {
	// 从环境变量或直接设置 API Key
	//apiKey := "e4aba69a-ce15-4709-b73e-35bb44230eb5" // 测试使用 或者从环境变量获取
	apiKey := "cafbfee7-f78a-4494-b721-5fce238f2869" // 生产使用 或者从环境变量获取
	//if apiKey == "" {
	//	apiKey = os.Getenv("ARK_API_KEY")
	//	if apiKey == "" {
	//		log.Fatal("请设置 ARK_API_KEY 环境变量或直接提供 API Key")
	//	}
	//}

	// 创建客户端
	client := seedream.NewARKClient(apiKey)

	// 提示词
	//prompt := "星际穿越，黑洞，黑洞里面映射出21维空间的图书馆，抢视觉冲击力，电影大片，末日既视感，动感，对比色，oc渲染，光线追踪，动态模糊，景深，超现实主义，深蓝，画面通过细腻的丰富的色彩层次塑造主体与场景，质感真实，暗黑风背景的光影效果营造出氛围，整体兼具艺术幻想感，夸张的广角透视效果，耀光，反射，极致的光影，强引力，吞噬"
	//prompt2 := "青春男士打篮球投篮的动作的照片，全景、侧面。照片宽高为 1500*1000像素"

	//生成图像
	log.Println("正在生成图像...")
	//response, err := client.GenerateImageWithPrompt(prompt2, seedream.WithSeed(1))

	//response, err := client.GenerateImageWithPrompt(prompt2)
	//if err != nil {
	//	log.Fatalf("生成图像失败: %v", err)
	//}
	//
	//// 输出结果
	//log.Printf("\n=== 图像生成成功 ===\n")
	//log.Printf("模型: %s\n", response.Model)
	//log.Printf("创建时间: %d\n", response.Created)
	//log.Printf("生成图片数量: %d\n", response.Usage.GeneratedImages)
	//log.Printf("输出Token: %d\n", response.Usage.OutputTokens)
	//log.Printf("总Token: %d\n", response.Usage.TotalTokens)
	//
	//for i, data := range response.Data {
	//	log.Printf("\n图片 %d:\n", i+1)
	//	log.Printf("URL: %s\n", data.URL)
	//	if data.Size != "" {
	//		log.Printf("尺寸: %s\n", data.Size)
	//	}
	//}

	//// 使用完整请求对象的方式
	//log.Println("\n=== 使用完整请求对象 ===")
	//fullRequest := &seedream.ImageGenerationRequest{
	//	Prompt:                    prompt2,
	//	Model:                     seedream.DefaultModel,
	//	Size:                      seedream.Size2K,
	//	ResponseFormat:            seedream.FormatURL,
	//	SequentialImageGeneration: seedream.SequentialDisabled,
	//	Watermark:                 false,
	//	Stream:                    false,
	//}
	//
	//response2, err := client.GenerateImage(fullRequest)
	//if err != nil {
	//	log.Printf("生成多张图像失败: %v", err)
	//}
	//
	//log.Printf("成功生成 %d 张图片\n", response2.Usage.GeneratedImages)

	//for i, data := range response2.Data {
	//	fmt.Println("图片 URL:", i, data.URL)
	//}
	//err = httpclient.NewHTTPClient("").DownloadFile(response.Data[0].URL, "seedream_generated_image"+strconv.FormatInt(time.Now().Unix(), 10)+".png", nil)
	//if err != nil {
	//	log.Printf("下载图片失败: %v", err)
	//} else {
	//	log.Println("图片下载成功")
	//}

	url := "https://ark-content-generation-v2-cn-beijing.tos-cn-beijing.volces.com/doubao-seedream-4-0/021758617304361860cde38290461f4b74ab5abe39cd41d213235_0.jpeg?X-Tos-Algorithm=TOS4-HMAC-SHA256&X-Tos-Credential=AKLTYWJkZTExNjA1ZDUyNDc3YzhjNTM5OGIyNjBhNDcyOTQ%2F20250923%2Fcn-beijing%2Ftos%2Frequest&X-Tos-Date=20250923T084829Z&X-Tos-Expires=86400&X-Tos-Signature=89760a1d43689dfcae5ec21fb596e4edced912f5d78b731a258ce150a7b58b27&X-Tos-SignedHeaders=host"

	response, err := client.GenerateImageWithPromptAndImage("增加万圣节元素和万圣节背景", url)
	err = httpclient.NewHTTPClient("").DownloadFile(response.Data[0].URL, "seedream_generated_image"+strconv.FormatInt(time.Now().Unix(), 10)+".png", nil)
	if err != nil {
		log.Printf("下载图片失败: %v", err)
	} else {
		log.Println("图片下载成功")
	}

}
