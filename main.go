package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"
)

type Event struct {
	Image     string
	Title     string
	Organizer string
	Location  string
	Time      time.Time
}

type PageData struct {
	Events         []Event
	SearchDate     string
	SearchLocation string
	IsSearching    bool
}

func main() {
	// Sample events data
	events := []Event{
		{Image: "download.jpg", Title: "Concert", Organizer: "Music Inc.", Location: "New York", Time: time.Now().Add(24 * time.Hour)},
		{Image: "download.jpg", Title: "Conference", Organizer: "Tech Co.", Location: "San Francisco", Time: time.Now().Add(48 * time.Hour)},
		{Image: "download.jpg", Title: "Workshop", Organizer: "Learn Co.", Location: "Chicago", Time: time.Now().Add(72 * time.Hour)},
		{Image: "download.jpg", Title: "Seminar", Organizer: "Edu Inc.", Location: "Los Angeles", Time: time.Now().Add(96 * time.Hour)},
		{Image: "download.jpg", Title: "Exhibition", Organizer: "Art Co.", Location: "Boston", Time: time.Now().Add(120 * time.Hour)},
		{Image: "download.jpg", Title: "Concert", Organizer: "Music Inc.", Location: "New York", Time: time.Now().Add(24 * time.Hour)},
		{Image: "download.jpg", Title: "Concert", Organizer: "Music Inc.", Location: "New York", Time: time.Now().Add(24 * time.Hour)},
		{Image: "download.jpg", Title: "Concert", Organizer: "Music Inc.", Location: "New York", Time: time.Now().Add(24 * time.Hour)},
		{Image: "download.jpg", Title: "Concert", Organizer: "Music Inc.", Location: "New York", Time: time.Now().Add(24 * time.Hour)},
		{Image: "download.jpg", Title: "Concert", Organizer: "Music Inc.", Location: "New York", Time: time.Now().Add(24 * time.Hour)},
		{Image: "download.jpg", Title: "Concert", Organizer: "Music Inc.", Location: "New York", Time: time.Now().Add(24 * time.Hour)},
		{Image: "download.jpg", Title: "Concert", Organizer: "Music Inc.", Location: "New York", Time: time.Now().Add(24 * time.Hour)},
		{Image: "download.jpg", Title: "Concert", Organizer: "Music Inc.", Location: "New York", Time: time.Now().Add(24 * time.Hour)},
		{Image: "download.jpg", Title: "Concert", Organizer: "Music Inc.", Location: "New York", Time: time.Now().Add(24 * time.Hour)},
		{Image: "download.jpg", Title: "Concert", Organizer: "Music Inc.", Location: "New York", Time: time.Now().Add(24 * time.Hour)},
	}

	// Parse templates
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/index.html", "templates/card.html"))

	// Handle root route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pageData := PageData{
			Events:      events,
			IsSearching: false,
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}

			searchDate := strings.TrimSpace(r.FormValue("search_date"))
			searchLocation := strings.TrimSpace(r.FormValue("search_location"))

			if searchDate != "" || searchLocation != "" {
				pageData.SearchDate = searchDate
				pageData.SearchLocation = searchLocation
				pageData.IsSearching = true
				pageData.Events = filterEvents(events, searchDate, searchLocation)
			}
		}

		tmpl.ExecuteTemplate(w, "layout", pageData)
	})

	// Handle location suggestions
	http.HandleFunc("/suggest", func(w http.ResponseWriter, r *http.Request) {
		query := normalizeString(r.URL.Query().Get("q"))
		suggestions := []string{}

		for _, event := range events {
			location := normalizeString(event.Location)
			if strings.Contains(location, query) {
				suggestions = append(suggestions, event.Location)
			}
		}

		// Remove duplicates
		uniqueSuggestions := removeDuplicates(suggestions)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(uniqueSuggestions)
	})

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func filterEvents(events []Event, searchDate, searchLocation string) []Event {
	var filteredEvents []Event
	searchLocation = normalizeString(searchLocation)

	for _, event := range events {
		matchesDate := searchDate == "" || event.Time.Format("2006-01-02") == searchDate
		matchesLocation := searchLocation == "" || strings.Contains(normalizeString(event.Location), searchLocation)

		if matchesDate && matchesLocation {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents
}

func normalizeString(s string) string {
	// Convert to lowercase and remove spaces
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1 // drop the space
		}
		return unicode.ToLower(r)
	}, s)
}

func removeDuplicates(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
