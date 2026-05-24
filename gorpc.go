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
    id := flag.String("client-id", "", "Client ID (required)")
    once := flag.Bool("once", false, "Execute and exit")
    retry := flag.Int("retry", 5, "Seconds to wait before reconnect attempts")
    activity := client.Activity{}

    flag.StringVar(&activity.Details, "details", "", "Main details text")
    flag.StringVar(&activity.State, "state", "", "Secondary text")
    flag.StringVar(&activity.LargeImage, "large-image", "", "Large image asset key")
    flag.StringVar(&activity.LargeText, "large-text", "", "Large image hover text")
    flag.StringVar(&activity.SmallImage, "small-image", "", "Small image asset key")
    flag.StringVar(&activity.SmallText, "small-text", "", "Small image hover text")
    flag.Parse()

    if *id == "" {
        fmt.Fprintln(os.Stderr, "Error: Client ID is required")
        flag.Usage()

        os.Exit(1)
    }

    for {
        err := client.Login(*id)
        if err == nil {
            break
        }

        fmt.Fprintf(os.Stderr, "Connection error: %v. Retrying in %ds...\n", err, *retry)
        time.Sleep(time.Duration(*retry) * time.Second)
    }

    if err := client.SetActivity(activity); err != nil {
        fmt.Fprintln(os.Stderr, "Error while updating presence:", err)
        client.Logout()

        os.Exit(1)
    }

    fmt.Println("Updated!")

    if *once {
        client.Logout()
        return
    }

    sig := make(chan os.Signal, 1)

    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

    tick := time.NewTicker(15 * time.Second)
    defer tick.Stop()

    for {
        select {
            case <-sig:
                fmt.Println("\nClearing...")
                client.Logout()

                return

            case <-tick.C:
                _ = client.SetActivity(activity)

        }
    }
}