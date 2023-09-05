package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	targetURL := "https://www.amazon.com/-/zh/dp/B00006RVTS/ref=sr_1_2?_encoding=UTF8&content-id=amzn1.sym.eb39b83d-c690-496d-9f16-0a9bd66ca6c8&pd_rd_r=069fff53-a1e8-42e7-bbbc-953f2a73c5bf&pd_rd_w=1cbpz&pd_rd_wg=cZwfd&pf_rd_p=eb39b83d-c690-496d-9f16-0a9bd66ca6c8&pf_rd_r=BDNT5B4G9XDCHRYVG3Y9&qid=1693907429&refinements=p_36%3A-3000&rnid=386491011&s=toys-and-games&sr=1-2&th=1"
	//禁用headless-chrome 无头浏览器  推荐debug模式下使用
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // 禁用 Headless 模式
		//chromedp.Flag("window-size", "1920,1080"), // 设置窗口大小
		chromedp.Flag("incognito", true), // 设置无痕浏览模式
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// 创建锁，以确保只有一个操作能够打印内容
	var wg sync.WaitGroup
	wg.Add(1)

	// 创建自定义动作函数
	printAfterLoad := func(ctx context.Context) error {
		defer wg.Done()

		// 获取HTML内容
		var htmlContent string
		err := chromedp.Run(ctx,
			chromedp.InnerHTML("html", &htmlContent),
		)
		if err != nil {
			return err
		}

		WriteFile(time.Now().Format("2006-01-02-15-04-05")+"test.html", htmlContent)
		if err != nil {
			return err
		}
		fmt.Println(htmlContent)

		return nil
	}

	// 提供目标URL给chromedp库，并在页面加载完成后执行自定义动作函数
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(printAfterLoad),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 等待自定义动作函数完成打印操作
	wg.Wait()
}

// 写入文件信息  调用a := models.WriteFile("testaaa.txt", "xxxxxxxxxxxxx")
func WriteFile(filename string, str string) bool {
	//绝对路径
	pf := "D:\\go-test\\test_curl\\" + filename

	file, err := os.Create(pf)

	if err != nil {
		fmt.Println(err)
	}

	//测试打印
	fmt.Println(" Write to file : " + pf)

	n, err := io.WriteString(file, str)

	if err != nil {
		fmt.Println(n, err)
	}

	file.Close()

	return true
}
