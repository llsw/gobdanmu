package main

import (
	"github.com/electricbubble/go-toast"
)

func main() {
	// _ = toast.Push("test message")
	// _ = toast.Push("test message", toast.WithTitle("app title"))
	_ = toast.Push("test message",
		toast.WithTitle("app title"),
		toast.WithSubtitle("app sub title"),
		toast.WithAudio(toast.Ping),
		// toast.WithObjectiveC(true),
	)
}
