package main

import (
	// allows interaction with html template
	"fmt"
	"log"
	"net/http" // access to core go http functionality
	"os"
	"text/template"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// Welcome type made for testing
type Welcome struct {
	Name string
}

func main() {

	// connect the css to the html
	http.Handle("/static/", // /static/ is the url that html can refer to when looking for css, can be whatever
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")))) // go looks in the relative "static" directory first using http.FileServer(), then
	// matches it to a url of our choice ("/static/")

	// data to be passed to the template when the page starts
	welcome := Welcome{"Anonymous"}

	// give Go path to the html file and parse
	templates := template.Must(template.ParseFiles("./templates/home-template.html"))

	// what to do on the home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// if there is an error while executing the home page template, print it
		if err := templates.ExecuteTemplate(w, "home-template.html", welcome); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	//---------------- LEARNING AUDIO STUFF --------------------------------

	// open the audio file
	f, err := os.Open("static/audio/Alesis-Fusion-Tubular-Bells-C6.wav")
	if err != nil {
		log.Fatal(err)
	}

	// assign a streamer and format to the wav file that can decode when necessary
	// streamer - can decode and play the audio file, stateful, meaning that once it has been streamed, it cannot be streamed again until reset
	// format - stores info about the audio file, namely, the sample rate
	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	// initialize the speaker with the sample rate and buffer size
	// only need to call this once, at the beginning of the program, otherwise, if called multiple times, cannot play multiple sounds at once
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// speaker.Play(streamer)

	// start the server, open the port to 4200, without a path it assumes localhost
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":4200", nil))
}
