package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/faiface/beep"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// Welcome type made for testing
type Welcome struct {
	Name string
}

// function to log errors
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// function to open an audio file and return a streamer
// needs to be provided a filepath
func getStreamer(filepath string) (beep.StreamSeekCloser, beep.Format) {
	f, err := os.Open(filepath)
	checkError(err)
	// assign a streamer and format to the wav file that can decode when necessary
	// streamer - can decode and play the audio file, stateful, meaning that once it has been streamed, it cannot be streamed again until reset
	// format - stores info about the audio file, namely, the sample rate
	streamer, format, err2 := wav.Decode(f)
	checkError(err2)
	return streamer, format
}

// function to play an audio file given the filepath and a pre-initialized speaker
func playAudio(filepath string) {
	streamer, _ := getStreamer(filepath)
	speaker.Play(streamer)
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

	// initialize the speaker with the sample rate and buffer size with one of the samples in the library
	_, format := getStreamer("static/audio/Alesis-Fusion-Tubular-Bells-C6.wav")
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// what to do on the callme page
	http.HandleFunc("/callme", func(w http.ResponseWriter, r *http.Request) {
		// play the audio corresponding to the filepath
		playAudio("static/audio/Alesis-Fusion-Tubular-Bells-C6.wav")
		// log the action
		fmt.Println("Play Audio")
	})

	// start the server, open the port to 4200, without a path it assumes localhost
	fmt.Println("Listening...")
	fmt.Println(http.ListenAndServe(":4200", nil))
}
