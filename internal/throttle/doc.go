// Package throttle implements a fixed-window token-bucket throttle used
// to cap the number of drift-check cycles that driftwatch dispatches
// within a configurable rolling time window.
//
// Typical usage:
//
//	th, err := throttle.New(time.Minute, 10)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Inside the scheduler job:
//	if err := th.Allow(); err != nil {
//		// skip this cycle — budget exhausted
//		return
//	}
//	// … run drift detection …
//
// The window resets automatically once the current period expires;
// no background goroutine is required.
package throttle
