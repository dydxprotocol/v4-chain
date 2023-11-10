package daemons

import "time"

// SingleTickTickerAndStop creates a ticker that ticks once before the stop channel is signaled.
func SingleTickTickerAndStop() (*time.Ticker, chan bool) {
	// Create a ticker with a duration long enough that we do not expect to see a tick within the timeframe
	// of a normal unit test.
	ticker := time.NewTicker(10 * time.Minute)
	// Override the ticker's channel with a new channel we can insert into directly, and add a single tick.
	newChan := make(chan time.Time, 1)
	newChan <- time.Now()
	ticker.C = newChan

	stop := make(chan bool, 1)

	// Start a go-routine that will signal the stop channel once the single tick is consumed.
	go func() {
		for {
			// Once the single tick is consumed, stop the ticker and signal the stop channel.
			if len(ticker.C) == 0 {
				stop <- true
				close(stop)
				ticker.Stop()
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	return ticker, stop
}
