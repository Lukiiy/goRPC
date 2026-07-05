package main

import (
    "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/hugolgst/rich-go/client"
)

func main() {
	mainGui()
}

func update(clientId string, activity client.Activity, once bool, retry int) error {
	if clientId == "" {
        fmt.Fprintln(os.Stderr, "Client ID is required")
        flag.Usage()

        os.Exit(1)
    }

	for {
		err := client.Login(clientId)
		if err == nil {
			break
		}

		fmt.Fprintf(os.Stderr, "Connection error: %v. Retrying in %ds...\n", err, retry)
		time.Sleep(time.Duration(retry) * time.Second)
	}

	if err := client.SetActivity(activity); err != nil {
		client.Logout()

		return fmt.Errorf("Couldn't update presence due to %w", err)
	}

	fmt.Println("Updated!")

	if once {
		client.Logout()
		return nil
	}

	// run & signal handling
	sig := make(chan os.Signal, 1)

    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

    tick := time.NewTicker(15 * time.Second)
    defer tick.Stop()

	for {
		select {
			case <- sig:
				fmt.Println("\nClearing...")

				client.Logout()
				return nil

			case <- tick.C: // Periodically refresh the activity state
				_ = client.SetActivity(activity)
		}
	}
}