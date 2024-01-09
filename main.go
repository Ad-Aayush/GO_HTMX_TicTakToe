package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type game struct {
	Grid   []string
	Turn   string
	Result string
}

func flip_turn(s string) string {
	if s == "X" {
		return "O"
	}
	return "X"
}

func check_win(grid []string) string {
	for i := 0; i < 9; i += 3 {
		if grid[i] != "" && grid[i] == grid[i+1] && grid[i+1] == grid[i+2] {
			return "Win"
		}
	}

	for i := 0; i < 3; i++ {
		if grid[i] != "" && grid[i] == grid[i+3] && grid[i+3] == grid[i+6] {
			return "Win"
		}
	}

	if grid[0] != "" && grid[0] == grid[4] && grid[4] == grid[8] {
		return "Win"
	}

	if grid[2] != "" && grid[2] == grid[4] && grid[4] == grid[6] {
		return "Win"
	}

	for _, x := range grid {
		if x == "" {
			return ""
		}
	}
	return "Draw"
}

func main() {
	state := make(map[string]game)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Game := game{[]string{"", "", "", "", "", "", "", "", ""}, "X", ""}

		state["Grid"] = Game
		temp := template.Must(template.ParseFiles("index.html"))
		temp.Execute(w, state)
	})

	http.HandleFunc("/click", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Click")
		idStr := r.FormValue("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatal("ERROR: ", err)
		}
		if state["Grid"].Result != "" {
			temp := template.Must(template.ParseFiles("index.html"))
			temp.ExecuteTemplate(w, "Play", state)
			return
		}
		cur_turn := state["Grid"].Turn
		// state["Game"].Grid[id] = state["Game"].Turn

		if entry, ok := state["Grid"]; ok {

			// Then we modify the copy
			entry.Grid[id] = state["Grid"].Turn

			// Then we reassign map entry
			state["Grid"] = entry
		}

		res := check_win(state["Grid"].Grid)
		if res == "" {
			if entry, ok := state["Grid"]; ok {

				// Then we modify the copy
				entry.Turn = flip_turn(cur_turn)

				// Then we reassign map entry
				state["Grid"] = entry
			}
		} else {
			if entry, ok := state["Grid"]; ok {

				// Then we modify the copy
				entry.Result = res

				// Then we reassign map entry
				state["Grid"] = entry
			}
		}
		temp := template.Must(template.ParseFiles("index.html"))
		temp.ExecuteTemplate(w, "Play", state)
	})
	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		Game := game{[]string{"", "", "", "", "", "", "", "", ""}, "X", ""}

		state["Grid"] = Game
		temp := template.Must(template.ParseFiles("index.html"))
		temp.ExecuteTemplate(w, "Play", state)
	})
	http.ListenAndServe(":8000", nil)

}
