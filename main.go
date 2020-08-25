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
	streamers []beep.StreamSeekCloser
}

// function to add something to the queue of streamers
func (q *Queue) Add(streamers []beep.StreamSeekCloser) {
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
	// play the sound
	speaker.Play(streamer)
	streamer.Seek(0)
}

// function to return an array of streamers based on their filepaths
func getStreamers(filepaths ...string) []beep.StreamSeekCloser {
	streamerArray := make([]beep.StreamSeekCloser, len(filepaths))
	for i := 0; i < len(streamerArray); i++ {
		checkIfExists(filepaths[i])
		streamer, _ := getStreamer(filepaths[i])
		streamerArray[i] = streamer
	}
	return streamerArray
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

	// what to do on the callme page
	http.HandleFunc("/callme", func(w http.ResponseWriter, r *http.Request) {
		// initialize the speaker with the sample rate and buffer size with one of the samples in the library
		_, format := getStreamer("static/audio/Alesis-Fusion-Tubular-Bells-C6.wav")
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		// play the audio corresponding to the filepath
		playAudio("static/audio/Alesis-Fusion-Tubular-Bells-C6.wav")
		// log the action
		fmt.Println("Play Audio")
	})

	// multiple noises played after each other
	http.HandleFunc("/multiNoise", func(w http.ResponseWriter, r *http.Request) {
		// initialize the speaker with the sample rate and buffer size with one of the samples in the library
		_, format := getStreamer("static/audio/Alesis-Fusion-Tubular-Bells-C6.wav")
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		// get all of the streamers necessary
		streamerArray := getStreamers("static/audio/drums/kick/kick.wav",
			"static/audio/drums/snare/snare909.wav", "static/audio/drums/kick/kick.wav", "static/audio/drums/snare/snare909.wav")
		// pattern of a drum and a snare for a measure
		// (kick snare kick snare)

		// add them to a queue and play the queue
		var queue Queue
		queue.Add(streamerArray)
		speaker.Play(&queue)
	})

	drumsSelected := []string{}

	http.HandleFunc("/fillForm", func(w http.ResponseWriter, r *http.Request) {
		// if there is an error while executing the home page template, print it
		if err := templates.ExecuteTemplate(w, "home-template.html", welcome); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		r.ParseForm()
		fmt.Printf("%+v\n", r.Form)
		drumsSelected = append(drumsSelected, r.Form["drums"]...)
		// initialize the speaker with the sample rate and buffer size with one of the samples in the library
		_, format := getStreamer("static/audio/Alesis-Fusion-Tubular-Bells-C6.wav")
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		// get the corresponding streamers
		streamerArray := getStreamers(drumsSelected...)

		// add them to a queue depending on which checkboxes are ticked and play the queue
		var queue Queue
		queue.Add(streamerArray)
		speaker.Play(&queue)
	})

	// playing a drum pattern based on the checkboxes ticked
	http.HandleFunc("/playDrums", func(w http.ResponseWriter, r *http.Request) {
		// initialize the speaker with the sample rate and buffer size with one of the samples in the library
		_, format := getStreamer("static/audio/Alesis-Fusion-Tubular-Bells-C6.wav")
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		// get the corresponding streamers
		streamerArray := getStreamers(drumsSelected...)

		// add them to a queue depending on which checkboxes are ticked and play the queue
		var queue Queue
		queue.Add(streamerArray)
		speaker.Play(&queue)
	})

	// start the server, open the port to 4200, without a path it assumes localhost
	fmt.Println("Listening...")
	fmt.Println(http.ListenAndServe(":4200", nil))
}
