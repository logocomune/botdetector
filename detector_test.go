package botdetector

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	browserUaList = "./test-ua/browsers.txt"
	spidersUaList = "./test-ua/spiders.txt"
)

var browsers = []string{
	"Flock/14.15 (Android 2.9; fr_BE;)",
	"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-GB; rv:1.9.0.16) Gecko/2010021003 Firefox/3.0.16 Flock/2.5.6",
	"Mozilla/5.0 (Windows; U; Windows NT 6.0; en-US; rv:1.8.1.14) Gecko/20080414 Firefox/2.0.0.14 Flock/1.1.2",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.6) Gecko/20070801 Firefox/2.0.0.6 Flock/0.9.0.2",
	"NetNewsWire/2.x (Mac OS X; http://ranchero.com/netnewswire/)",

	"iTunes/4.0 (Macintosh; U; PPC Mac OS X 10.2)",
	"iTunes/4.2 (Macintosh; U; PPC Mac OS X 10.2)",
	"iTunes/4.7 (Macintosh; N; PPC)",
	"iTunes/4.7 (Macintosh; U; PPC Mac OS X 10.2)",
	"iTunes/4.8 (Macintosh; U; PPC Mac OS X 10.4.1)",
	"iTunes/7.0 (Macintosh; U; PPC Mac OS X 10.4.7)",
	"iTunes/7.0.1 (Windows; N)",
	"iTunes/7.1.1 (Macintosh; N; PPC)",
	"iTunes/7.4.1",
	"iTunes/7.5 (Macintosh; N; PPC)",
	"iTunes/7.6.2.9",
	"iTunes/8.0",
	"iTunes/8.1",
	"iTunes/8.1.1 (Windows; N)",
	"iTunes/8.1.1 (Windows; U)",
	"iTunes/8.2 (Macintosh; U; PPC Mac OS X 10_5_6)",
	"iTunes/9.0",
	"iTunes/9.0 (Macintosh; Intel Mac OS X 10.5.8)",
	"iTunes/9.0 (Macintosh; Intel Mac OS X 10.5.8) AppleWebKit/531.9",
	"iTunes/9.0.2 (Windows; N)",
	"iTunes/9.0.3",
	"iTunes/9.0.3 (Macintosh; U; Intel Mac OS X 10_6_2; en-ca)",
	"iTunes/9.1.1",
}

var google = []string{
	//https://support.google.com/webmasters/answer/1061943?hl=it
	"APIs-Google (+https://developers.google.com/webmasters/APIs-Google.html)",
	"Mediapartners-Google",
	"Mozilla/5.0 (Linux; Android 5.0; SM-G920A) AppleWebKit (KHTML, like Gecko) Chrome Mobile Safari (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
	"AdsBot-Google (+http://www.google.com/adsbot.html)",
	"Googlebot-Image/1.0",
	"Googlebot-News",
	"Googlebot-Video/1.0",
	"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Safari/537.36",
	"Googlebot/2.1 (+http://www.google.com/bot.html)",
	"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.96 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"(compatible; Mediapartners-Google/2.1; +http://www.google.com/bot.html)",
	"AdsBot-Google-Mobile-Apps",
	"FeedFetcher-Google; (+http://www.google.com/feedfetcher.html)",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36 (compatible; Google-Read-Aloud; +https://support.google.com/webmasters/answer/1061943)",
	"google-speakr",
	"Mozilla/5.0 (Linux; Android 8.0; Pixel 2 Build/OPD3.170816.012; DuplexWeb-Google/1.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Mobile Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)  Chrome/49.0.2623.75 Safari/537.36 Google Favicon",
	"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3694.0 Mobile Safari/537.36 Chrome-Lighthouse",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3694.0 Safari/537.36 Chrome-Lighthouse",
}

var baidu = []string{
	"Mozilla/5.0 (compatible; Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (Linux; U; Android 4.1; en-us; GT-N7100 Build/JRO03C;Baiduspider-ads)AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
	"Mozilla/5.0 (Linux;u;Android 4.2.2;zh-cn;) AppleWebKit/534.46 (KHTML,like Gecko) Version/5.1 Mobile Safari/10600.6.3 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
	"Baiduspider-image+(+http://www.baidu.com/search/spider.htm)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1 (compatible; Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:17.0; Baiduspider-ads) Gecko/17.0 Firefox/17.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36 Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html",
	"Baiduspider+(+http://www.baidu.com/search/spider.htm)",
	"Mozilla/5.0 (compatible; Baiduspider-cpro; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (Linux;u;Android 2.3.7;zh-cn;) AppleWebKit/533.1 (KHTML,like Gecko) Version/4.0 Mobile Safari/533.1 (compatible; +http://www.baidu.com/search/spi_der.html)",
	"Mozilla/5.0 (Linux;u;Android 2.3.7;zh-cn;) AppleWebKit/533.1 (KHTML,like Gecko) Version/4.0 Mobile Safari/533.1 (compatible; +http://www.baidu.com/search/spider.html)",
	"Baiduspider-image+(+http://www.baidu.com/search/spider.htm)",
	"Baiduspider+(+http://www.baidu.jp/spider/)",
	"Baiduspider+(+http://help.baidu.jp/system/05.html)",
	"Baiduspider+(+http://www.baidu.com/search/spider_jp.html)",
}

var bing = []string{
	//https://www.bing.com/webmaster/help/which-crawlers-does-bing-use-8c184ec0
	"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11A465 Safari/9537.53 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
	"Mozilla/5.0 (Windows Phone 8.1; ARM; Trident/7.0; Touch; rv:11.0; IEMobile/11.0; NOKIA; Lumia 530) like Gecko (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
	"Mozilla/5.0 (compatible; adidxbot/2.0; +http://www.bing.com/bingbot.htm)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11A465 Safari/9537.53 (compatible; adidxbot/2.0; +http://www.bing.com/bingbot.htm)",
	"Mozilla/5.0 (Windows Phone 8.1; ARM; Trident/7.0; Touch; rv:11.0; IEMobile/11.0; NOKIA; Lumia 530) like Gecko (compatible; adidxbot/2.0; +http://www.bing.com/bingbot.htm)",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534+ (KHTML, like Gecko) BingPreview/1.0b",
	"Mozilla/5.0 (Windows Phone 8.1; ARM; Trident/7.0; Touch; rv:11.0; IEMobile/11.0; NOKIA; Lumia 530) like Gecko BingPreview/1.0b",
}

var yandex = []string{
	//https://yandex.com/support/webmaster/robot-workings/check-yandex-robots.html
	"Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 8_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B411 Safari/600.1.4 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexAccessibilityBot/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 8_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B411 Safari/600.1.4 (compatible; YandexMobileBot/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexDirectDyn/1.0; +http://yandex.com/bots",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36 (compatible; YandexScreenshotBot/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexImages/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexVideo/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexVideoParser/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexMedia/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexFavicons/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexWebmaster/2.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexPagechecker/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexImageResizer/2.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YaDirectFetcher/1.0; Dyatel; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexCalendar/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexSitelinks; Dyatel; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexMetrika/2.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexNews/4.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexVertis/3.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexBot/3.0; MirrorDetector; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexSearchShop/1.0; +http://yandex.com/bots)",
	"Mozilla/5.0 (compatible; YandexVerticals/1.0; +http://yandex.com/bots)",
}

var yahoo = []string{
	"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)",
	"Mozilla/5.0 (compatible; Yahoo! Slurp China; http://misc.yahoo.com.cn/help.html)",
	"Y!J-BRW/1.0 (https://www.yahoo-help.jp/app/answers/detail/p/595/a_id/42716)",
}

var crawler = []string{
	"Blackfire PHP Player/1.0",
}

func TestUABotDetector_IsBotStrict(t *testing.T) {
	rules := []string{
		"^b0t$",
		"^Amazon Simple Notification Service Agent$",
	}

	d := newDetector(rules)
	assert.True(t, d.IsBot("b0t"))
	assert.True(t, d.IsBot("Amazon Simple Notification Service Agent"))
	assert.False(t, d.IsBot("It's b0t"))
}

func TestUABotDetector_IsBotStartWith(t *testing.T) {
	rules := []string{
		"^Java/1.6.0_03",
	}

	d := newDetector(rules)

	assert.True(t, d.IsBot("Java/1.6.0_03"))
	assert.False(t, d.IsBot("It's Java/1.6.0_03"))
}

func TestUABotDetector_IsBotContains(t *testing.T) {
	rules := []string{
		"AHC/",
	}

	u := newDetector(rules)
	assert.True(t, u.IsBot("AHC/1.0"))
	assert.True(t, u.IsBot("It's an AHC/1.0"))
}

func TestUA(t *testing.T) {
	u := New()

	useragents := append(google, bing...)
	useragents = append(useragents, baidu...)
	useragents = append(useragents, yandex...)
	useragents = append(useragents, yahoo...)
	useragents = append(useragents, crawler...)

	for _, c := range useragents {
		isC := u.IsBot(c)
		if !isC {
			t.Log(c)
		}

		assert.True(t, isC)
	}
}
func TestBrowsers(t *testing.T) {
	u := New()

	for _, c := range browsers {
		isC := u.IsBot(c)
		if isC {
			t.Log(c)
		}

		assert.False(t, isC)
	}
}

func TestBrowsersUA(t *testing.T) {
	u := New()
	file, err := os.Open(browserUaList)

	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		userAgent := scanner.Text()
		isBot := u.IsBot(userAgent)

		if isBot {
			t.Log(userAgent)
		}

		assert.False(t, isBot)
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}

func TestSpidersUA(t *testing.T) {
	u := New()

	file, err := os.Open(spidersUaList)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		userAgent := scanner.Text()

		isBot := u.IsBot(userAgent)
		if !isBot {
			t.Log(userAgent)
		}

		assert.True(t, isBot)
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}
