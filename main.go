package main

import "log"

func main() {
	newDb, err := newPostgress()
	if err != nil {
		log.Fatal(err)
	}
	serv := newAPIServer(":3000", newDb)
	serv.run()
}
