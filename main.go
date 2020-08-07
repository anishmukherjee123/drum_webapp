package main

import (
	"fmt"      // allows interaction with html template
	"net/http" // access to core go http functionality
)

func main() {

	// connect the css to the html
	http.Handle("/static/", // /static/ is the url that html can refer to when looking for css, can be whatever
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")))) // go looks in the relative "static" directory first using http.FileServer(), then
	// matches it to a url of our choice ("/static/")

	// start the server, open the port to 4200, without a path it assumes localhost
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":4200", nil))
}
