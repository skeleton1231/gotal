package posix

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/skeleton1231/gotal/pkg/shutdown"
)

// Name is the identifier for the PosixSignalManager.
const Name = "PosixSignalManager"

// PosixSignalManager listens for POSIX signals and triggers a graceful shutdown.
type PosixSignalManager struct {
	signals []os.Signal
}

// NewPosixSignalManager creates a new PosixSignalManager that listens for the provided signals.
// If no signals are provided, it defaults to listening for SIGINT and SIGTERM.
func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = []os.Signal{os.Interrupt, syscall.SIGINT, syscall.SIGTERM} // Simplified signal slice initialization
	}

	return &PosixSignalManager{
		signals: sig,
	}
}

// GetName returns the name of this ShutdownManager.
func (pm *PosixSignalManager) GetName() string {
	return Name
}

// Start begins the signal listening process.
// It runs in its own goroutine to avoid blocking the caller.
func (pm *PosixSignalManager) Start(gs shutdown.GSInterface) error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, pm.signals...)
	go func() {
		for {
			sig := <-signalChan
			log.Printf("Received signal: %v\n", sig)
			gs.StartShutdown(pm)
			// Depending on your shutdown logic, you might want to break after initiating shutdown
		}
	}()

	return nil
}

// ShutdownStart is called when a shutdown sequence begins.
// This implementation currently does nothing but can be extended if needed.
func (pm *PosixSignalManager) ShutdownStart() error {
	// Implement any pre-shutdown initiation logic here if necessary.
	return nil
}

// ShutdownFinish is called after all shutdown callbacks have been completed.
// It terminates the application by calling os.Exit(0).
func (pm *PosixSignalManager) ShutdownFinish() error {
	// Implement any last-minute cleanup here if necessary.
	os.Exit(0) // Terminates the program. The line below will not be executed.
	return nil // This is unreachable code. Document that os.Exit is called above.
}
