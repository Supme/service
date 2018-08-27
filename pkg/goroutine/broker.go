package goroutine

import (
	"sync"
)

// Example code
//	func myFunc() {
//		fmt.Println("Starting...")
//		rand.Seed(time.Now().UnixNano())
//		wg := &sync.WaitGroup{}
//
//		broker := NewBroker(20)
//		for i := 1; i <= 300; i++ {
//			if i == 100 {
//				broker.SetMax(10)
//			}
//			if i == 200 {
//				broker.SetMax(15)
//			}
//			broker.Next()
//			wg.Add(1)
//			go func(i int) {
//				fmt.Println("Start worker number", i, "current workers", broker.GetCurrent(), "from maximum", broker.GetMax())
//				time.Sleep(time.Duration(rand.Int63n(500000)+500000) * time.Microsecond)
//				broker.Ready()
//				fmt.Println("Finish worker", i)
//				wg.Done()
//			}(i)
//		}
//		wg.Wait()
//	}

type Broker struct {
	maxCount     int64
	currentCount int64
	next         chan struct{}
	ready        chan struct{}
	mu           sync.RWMutex
}

func NewBroker(count int64) *Broker {
	b := Broker{maxCount: count, currentCount: 0}
	b.next = make(chan struct{}, 1)
	b.ready = make(chan struct{}, 1)
	b.start()
	return &b
}

func (b *Broker) start() {
	go func() {
		for range b.ready {
			b.mu.Lock()
			b.currentCount = b.currentCount - 1
			b.mu.Unlock()
		}
	}()
	go func() {
		for {
			b.mu.RLock()
			current := b.currentCount
			max := b.maxCount
			b.mu.RUnlock()
			if current < max {
				b.mu.Lock()
				b.currentCount = b.currentCount + 1
				b.mu.Unlock()
				b.next <- struct{}{}
			}
		}
	}()
}

func (b *Broker) SetMax(count int64) {
	b.mu.Lock()
	b.maxCount = count
	b.mu.Unlock()
}

func (b *Broker) GetCurrent() int64 {
	b.mu.RLock()
	current := b.currentCount
	b.mu.RUnlock()
	return current
}

func (b *Broker) GetMax() int64 {
	b.mu.RLock()
	max := b.maxCount
	b.mu.RUnlock()
	return max
}

func (b *Broker) Next() {
	<-b.next
}

func (b *Broker) Ready() {
	b.ready <- struct{}{}
}
