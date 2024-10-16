package tools

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/horoshi10v/grpc-experiment/proto"
)

// LoadGenerator генерує навантаження на сервер, створюючи кілька одночасних запитів
func LoadGenerator(client pb.ExperimentServiceClient, concurrentRequests int, duration time.Duration) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					_, err := client.RequestResponse(ctx, &pb.Request{Message: "Load Test Request"})
					if err != nil {
						log.Printf("Request failed: %v", err)
						return
					}
				}
			}
		}()
	}

	wg.Wait()
}
