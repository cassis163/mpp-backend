package models

import "database/sql"

type Movie struct {
	Id      string   		`json:"id"`
	Name    string   		`json:"name"`
	Year    string  		`json:"year"`
	Score   string   		`json:"score"`
	Summary sql.NullString	`json:"summary,omitempty"`
}
