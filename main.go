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

type node struct {
	Grid  [9]string
	Turn  string
	score int
	next  []*node
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

func build_tree(root *node) *node {

	res := check_win(root.Grid[:])
	if res == "Draw" {
		root.score = 0
		return root
	} else if res == "Win" {
		if root.Turn == "O" {
			root.score = 1
		} else {
			root.score = -1
		}
		return root
	}
	for i := range root.Grid {
		if root.Grid[i] == "" {
			toSend := new(node)
			*toSend = node{
				Grid:  [9]string{},
				Turn:  flip_turn(root.Turn),
				score: 0,
				next:  []*node{},
			}

			// Copy the grid values and update the current cell
			for j := 0; j < 9; j++ {
				toSend.Grid[j] = root.Grid[j]
			}
			toSend.Grid[i] = root.Turn

			// Append the new node after building its tree
			root.next = append(root.next, build_tree(toSend))
		}
	}
	if root.Turn == "X" {
		for _, x := range root.next {
			root.score = max(root.score, x.score)
		}
	} else {
		for _, x := range root.next {
			root.score = min(root.score, x.score)
		}
	}

	return root
}

func is_equal(a, b [9]string) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	state := make(map[string]game)

	root := new(node)
	root = &node{[9]string{"", "", "", "", "", "", "", "", ""}, "X", 0, []*node{}}
	root = build_tree(root)

	cur_state := root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Game := game{[]string{"", "", "", "", "", "", "", "", ""}, "X", ""}
		cur_state = root
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

		// state["Game"].Grid[id] = state["Game"].Turn

		if entry, ok := state["Grid"]; ok {

			// Then we modify the copy
			entry.Grid[id] = state["Grid"].Turn

			// Then we reassign map entry
			state["Grid"] = entry
		}

		res := check_win(state["Grid"].Grid)
		if res != "" {
			if entry, ok := state["Grid"]; ok {

				// Then we modify the copy
				entry.Result = res

				// Then we reassign map entry
				state["Grid"] = entry
			}
			temp := template.Must(template.ParseFiles("index.html"))
			temp.ExecuteTemplate(w, "Play", state)
			return
		}
		for i := range cur_state.next {
			if is_equal([9]string(state["Grid"].Grid), cur_state.next[i].Grid) {
				cur_state = cur_state.next[i]
				break
			}
		}
		min_ind := 0

		for i := range cur_state.next {
			if cur_state.next[i].score < cur_state.next[min_ind].score {
				min_ind = i
			}
		}
		cur_state = cur_state.next[min_ind]
		new_grid := cur_state.Grid
		if entry, ok := state["Grid"]; ok {

			// Then we modify the copy
			entry.Grid = new_grid[:]

			// Then we reassign map entry
			state["Grid"] = entry
		}
		res = check_win(state["Grid"].Grid)
		if res != "" {
			if entry, ok := state["Grid"]; ok {

				// Then we modify the copy
				entry.Result = res
				entry.Turn = flip_turn(entry.Turn)
				// Then we reassign map entry
				state["Grid"] = entry
			}
			temp := template.Must(template.ParseFiles("index.html"))
			temp.ExecuteTemplate(w, "Play", state)
			return
		}
		temp := template.Must(template.ParseFiles("index.html"))
		temp.ExecuteTemplate(w, "Play", state)
	})
	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		Game := game{[]string{"", "", "", "", "", "", "", "", ""}, "X", ""}
		cur_state = root
		state["Grid"] = Game
		temp := template.Must(template.ParseFiles("index.html"))
		temp.ExecuteTemplate(w, "Play", state)
	})
	http.ListenAndServe(":8000", nil)
}
