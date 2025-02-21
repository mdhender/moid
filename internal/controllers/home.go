// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package controllers

import (
	"database/sql"
	"fmt"
	"github.com/mdhender/moid/internal/views"
	"log"
	"math/rand/v2"
	"net/http"
)

type Home struct {
	db   *sql.DB
	view *views.View
}

// NewHomeController creates a new instance of the Home controller
func NewHomeController(db *sql.DB, view *views.View) (*Home, error) {
	c := &Home{
		db:   db,
		view: view,
	}
	// add any initialization logic here if needed
	return c, nil
}

type PageData struct {
	CounterMessage string
	ViewCount      string
}

var viewCount int

func (c Home) Show(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: %s\n", r.Method, r.URL.Path, r.RemoteAddr)

	messages := []string{
		`citizen interactions logged.`,
		`visitors. Noted. Next.`,
		`inquiries processed. Yawn.`,
		`reports filed. Thrilling.`,
		`forms reviewed. Barely.`,
		`subjects observed. Again.`,
		`data points. Nobody cares.`,
		`log entries. Bureaucracy endures.`,
		`visits documented. As required.`,
		`records updated. Policy met.`,
		`accesses. This matters? Sure.`,
		`individuals noted. No excitement.`,
		`transactions archived. Sigh.`,
		`engagements. It’s a living.`,
		`incidents catalogued. Proceed.`,
		`units of interest. Barely.`,
		`entries. No anomalies. Yet.`,
		`contacts recorded. Moving on.`,
		`digital echoes. Pointless.`,
		`glimpses of bureaucracy. Next.`,
		`observations logged. No urgency.`,
		`moments wasted. Keep going.`,
		`interactions. Wow. Or not.`,
		`accesses. No need to celebrate.`,
		`records filed. All routine.`,
		`disturbances detected. Unlikely.`,
		`data points. Thrilling work.`,
		`forms checked. Feigned interest.`,
		`units noted. Paperwork pending.`,
		`events catalogued. Awaiting coffee.`,
		`log entries. Carry on.`,
		`visitors. System remains unimpressed.`,
		`transactions noted. Fascinating.`,
		`glimpses into the void. Huh.`,
		`instances of participation. Barely.`,
		`interactions. Policy dictates logging.`,
	}

	// unsafe increment the view count, but who cares
	viewCount++

	data := PageData{
		CounterMessage: messages[rand.IntN(len(messages))],
		ViewCount:      fmt.Sprintf("%012d", viewCount),
	}

	// TODO: Implement home page logic
	// - Fetch required data from database
	// - Process any markdown content
	// - Handle any necessary encryption/decryption

	// - Render the template
	c.view.Render(w, r, "home.gohtml", data)
}
