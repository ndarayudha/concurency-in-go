package chapter4

import (
	"fmt"
	"testing"
	"time"
)

func TestOrChannels(t *testing.T) {
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})

		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()

		return orDone
	}

	// Simulate an exit route with security checks
	monitorExit := func(exitName string, timeToReach, securityCheckTime time.Duration) <-chan interface{} {
		ch := make(chan interface{})
		go func() {
			defer close(ch)

			fmt.Printf("[%s] VIP starts moving towards exit\n", exitName)
			time.Sleep(timeToReach) // Time to reach the exit

			fmt.Printf("[%s] VIP reached exit, starting security check\n", exitName)
			time.Sleep(securityCheckTime) // Security check time

			fmt.Printf("[%s] Security check complete, VIP can exit\n", exitName)
		}()
		return ch
	}

	fmt.Println("=== VIP Exit Protocol Initiated ===")
	fmt.Println("Security Director: All guards at positions")

	start := time.Now()

	// Monitor different exits with realistic timings
	mainEntrance := monitorExit("Main Entrance", 10*time.Second, 5*time.Second)  // Longest but most secure
	backDoor := monitorExit("Back Door", 5*time.Second, 3*time.Second)           // Medium route
	emergencyExit := monitorExit("Emergency Exit", 2*time.Second, 1*time.Second) // Fastest but least secure
	helipad := monitorExit("Helipad", 8*time.Second, 2*time.Second)              // Alternative route

	// Security director waits for first successful exit
	<-or(
		mainEntrance,
		backDoor,
		emergencyExit,
		helipad,
	)

	duration := time.Since(start)
	fmt.Println("\n=== Exit Protocol Complete ===")
	fmt.Printf("Security Director: VIP has successfully exited after %v\n", duration)
	fmt.Println("Security Director: All units stand down")
}
