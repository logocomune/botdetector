package botdetector

import (
	"fmt"
	"testing"
)

const (
	botUA     = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
	browserUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	lynxUA    = "Lynx/2.8.9dev.16 libwww-fm/2.14 SSL-MM/1.4.1 GNUTLS/3.5.17"
	cubotUA   = "Mozilla/5.0 (Linux; Android 9; CUBOT NOTE 20) AppleWebKit/537.36 yandexsearch/7.0"
)

// BenchmarkNew measures the cost of constructing a BotDetector with the full default rule set (~1446 rules).
func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = New()
	}
}

// BenchmarkIsBot_BotMatch measures detection of a well-known bot UA (Googlebot).
// With early-return this should exit as soon as the first matching rule fires.
func BenchmarkIsBot_BotMatch(b *testing.B) {
	d, _ := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = d.IsBot(botUA)
	}
}

// BenchmarkIsBot_BrowserNoMatch measures the worst-case path: a pure browser UA that
// matches no rule, forcing the detector to scan every rule before returning false.
func BenchmarkIsBot_BrowserNoMatch(b *testing.B) {
	d, _ := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = d.IsBot(browserUA)
	}
}

// BenchmarkIsBot_WithCache_Hit measures the hot path when the LRU cache already contains the result.
func BenchmarkIsBot_WithCache_Hit(b *testing.B) {
	d, _ := New(WithCache(1024))
	// warm up: populate the cache entry we will be hitting.
	_ = d.IsBot(botUA)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = d.IsBot(botUA)
	}
}

// BenchmarkIsBot_WithCache_Miss measures the cold-cache path by cycling through a pool
// of unique browser UAs so that each call is a cache miss.
func BenchmarkIsBot_WithCache_Miss(b *testing.B) {
	d, _ := New(WithCache(1))
	const poolSize = 1000
	uas := make([]string, poolSize)
	for i := range uas {
		uas[i] = fmt.Sprintf("Mozilla/5.0 Chrome/%d.0.0.0 Safari/537.36", i)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = d.IsBot(uas[i%poolSize])
	}
}

// BenchmarkNormalize_Normal measures normalize() on a typical browser UA with no special tokens.
func BenchmarkNormalize_Normal(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = normalize(browserUA)
	}
}

// BenchmarkNormalize_Lynx measures the Lynx fast-path in normalize() (early return after one Replace).
func BenchmarkNormalize_Lynx(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = normalize(lynxUA)
	}
}

// BenchmarkNormalize_WithReplacements measures normalize() when the UA contains tokens
// that trigger multiple replacements ("cubot", "yandexsearch/", …).
func BenchmarkNormalize_WithReplacements(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = normalize(cubotUA)
	}
}

// BenchmarkIsBot_ParallelBot measures throughput under concurrent access for a bot UA.
func BenchmarkIsBot_ParallelBot(b *testing.B) {
	d, _ := New()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = d.IsBot(botUA)
		}
	})
}

// BenchmarkIsBot_ParallelBrowser measures throughput under concurrent access for a browser UA (worst-case).
func BenchmarkIsBot_ParallelBrowser(b *testing.B) {
	d, _ := New()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = d.IsBot(browserUA)
		}
	})
}
