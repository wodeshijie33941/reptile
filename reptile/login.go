package reptile

import (
	"time"

	"github.com/chromedp/chromedp"
)

func Onlogin(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://yiqixie.com/e/login`),
		chromedp.WaitVisible(`.copyright`),

		// second param format: username \ pwd
		chromedp.SendKeys(`//input[@type="text"]`, "wodeshijie3394@126.com\t1qaz1QAZ"),

		chromedp.Sleep(2 * time.Second),
		chromedp.Click(`//button[@type="button"]`, chromedp.NodeVisible),
		chromedp.WaitVisible(`.i-o-sr-ho`),
		chromedp.OuterHTML(`.office-infinite-list-items-wrapper-outter`, res),
		chromedp.Sleep(2 * time.Second),
	}
}
