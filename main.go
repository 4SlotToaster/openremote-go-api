package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	db := make(database, 0)
	http.HandleFunc("/houses/populate", db.populate)
	http.HandleFunc("/houses/get", db.find)
	// An available port will be selected automatically
	log.Fatal(http.ListenAndServe(":0", nil))
}

type database map[int]string

func (db database) populate(w http.ResponseWriter, req *http.Request) {
	f, err := os.Open("houses.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	var i int
	for sc.Scan() {
		i++
		db[i] = sc.Text()
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	fmt.Println(db)
}

func (db database) find(w http.ResponseWriter, req *http.Request) {

	lineNumStr := req.URL.Query().Get("id")
	lineNum, _ := strconv.Atoi(lineNumStr)

	json, ok := db[lineNum]
	if !ok {
		http.Error(w, fmt.Sprintf("invalid query parameter: %s\n", lineNumStr), http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s\n", json)
	fmt.Fprintf(os.Stdout, "get: %v\n", req.URL.String())
}
