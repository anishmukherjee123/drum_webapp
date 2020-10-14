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

// type to hold a queue of drum samples
type Queue struct {
	streamers []beep.Streamer
}

// function to add something to the queue of streamers
func (q *Queue) Add(streamers []beep.Streamer) {
	q.streamers = append(q.streamers, streamers...)
}

// function to stream all of the samples in a queue
func (q *Queue) Stream(samples [][2]float64) (n int, ok bool) {
	// We use the filled variable to track how many samples we've
	// successfully filled already. We loop until all samples are filled.
	filled := 0
	for filled < len(samples) {
		// There are no streamers in the queue, so we stream silence.
		if len(q.streamers) == 0 {
			for i := range samples[filled:] {
				samples[i][0] = 0
				samples[i][1] = 0
			}
			break
		}

		// We stream from the first streamer in the queue.
		n, ok := q.streamers[0].Stream(samples[filled:])
		// If it's drained, we pop it from the queue, thus continuing with
		// the next streamer.
		if !ok {
			q.streamers = q.streamers[1:]
		}
		// We update the number of filled samples.
		filled += n
	}
	return len(samples), true
}

// Err: trivial error function for queues
func (q *Queue) Err() error {
	return nil
}

// function to log errors
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// function to check if a file exists
func checkIfExists(filepath string) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		log.Fatal(err)
	}
}

// function to open an audio file and return a streamer
// needs to be provided a filepath
func getStreamer(filepath string) (beep.Streamer, beep.Format) {
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
	// play the sound
	speaker.Play(streamer)
}

// function to return an array of streamers based on their filepaths
func getStreamers(filepaths ...string) []beep.Streamer {
	streamerArray := make([]beep.Streamer, len(filepaths))
	for i := 0; i < len(streamerArray); i++ {
		checkIfExists(filepaths[i])
		streamer, _ := getStreamer(filepaths[i])
		streamerArray[i] = streamer
	}
	return streamerArray
}

// function to remove a certain index from a slice
func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// // function that removes silences from the form input if there is a non-silent element directly after it
// func removeInvalidHidden(form url.Values, key string) []string {
// 	// construct an array of string values from the url.Values
// 	formValue := []string{form.Get("1")

// 	var resultArray []string
// 	copy(resultArray, formValue)
// 	keyCounter := 0
// 	for i := 0; i < len(formValue); i++ {
// 		if keyCounter != 0 && formValue[i] != key {
// 			keyCounter = 0
// 			removeIndex(resultArray, i-1)
// 		} else if formValue[i] == key {
// 			keyCounter++
// 		}
// 	}
// 	return resultArray
// }

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

	// playing drums based on the checkboxes ticked, creating a mixer for each row of checkboxes
	http.HandleFunc("/fillForm", func(w http.ResponseWriter, r *http.Request) {
		// if there is an error while executing the home page template, print it
		if err := templates.ExecuteTemplate(w, "home-template.html", welcome); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		r.ParseForm()
		fmt.Printf("%+v\n", r.Form)
		// // initialize the speaker with the sample rate and buffer size with one of the samples in the library
		_, format := getStreamer("static/audio/Alesis-Fusion-Tubular-Bells-C6.wav")
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		// check if any of the beats are empty, and if they are, include silence to fill the beat

		// get the corresponding streamers and create a mixed streamer on each beat
		// this needs to change dynamically
		streamer1Mix := beep.Mix(getStreamers(r.Form["1"]...)...)
		streamer2Mix := beep.Mix(getStreamers(r.Form["2"]...)...)
		streamer3Mix := beep.Mix(getStreamers(r.Form["3"]...)...)
		streamer4Mix := beep.Mix(getStreamers(r.Form["4"]...)...)
		mixedStreamer := []beep.Streamer{streamer1Mix, streamer2Mix, streamer3Mix, streamer4Mix}

		// add them to a queue depending on which checkboxes are ticked and play the queue
		var queue Queue
		queue.Add(mixedStreamer)
		speaker.Play(&queue)
	})

	http.HandleFunc("/playDrums", func(w http.ResponseWriter, r *http.Request) {
		// // initialize the speaker with the sample rate and buffer size with one of the samples in the library
		// _, format := getStreamer("static/audio/Alesis-Fusion-Tubular-Bells-C6.wav")
		// speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		// // get the corresponding streamers
		// streamerArray := getStreamers(drumsSelected...)

		// // add them to a queue depending on which checkboxes are ticked and play the queue
		// var queue Queue
		// queue.Add(streamerArray)
		// speaker.Play(&queue)
		// adding a test line here for git
	})

	// playing multiple sounds at once using the mixer type in the beep library
	// currently creates a new page that handles the form request of playing audio
	http.HandleFunc("/testMixer", func(w http.ResponseWriter, r *http.Request) {
		// initialize the speaker with the sample rate and buffer size with one of the samples in the library

		mixedStreamer := beep.Mix(getStreamers("static/audio/drums/snare/snare909.wav", "static/audio/drums/kick/kick.wav")...)
		speaker.Play(mixedStreamer)
	})

	// start the server, open the port to 4200, without a path it assumes localhost
	fmt.Println("Listening...")
	fmt.Println(http.ListenAndServe("https://anishmukherjee123.github.io/drum_webapp/", nil))
}
