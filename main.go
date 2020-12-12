package main

import (
	"github.com/gordonklaus/portaudio"
  // "github.com/mjibson/go-dsp/fft"
  "github.com/joho/godotenv"

	"time"
  "fmt"
  // "math"
  "os"
	"os/signal"
  "strconv"
)

var (
  scale int = 100
  mostFactor float64 = 0.5
  step int = 0
)

func main() {
  sig := make(chan os.Signal, 1)
  signal.Notify(sig, os.Interrupt, os.Kill)

  err := godotenv.Load()
  if err != nil {
    fmt.Println("Error loading .env file")
  }

  scale, _ = strconv.Atoi(os.Getenv("SCALE"))

  if scale <= 0 {
    scale = 100
  }

  if scale > 400 {
    scale = 400
  }

	portaudio.Initialize()
	defer portaudio.Terminate()
	e := newEcho(time.Second / 3)
	defer e.Close()
	chk(e.Start())
  for {

    select {
		case <-sig:
			return
		default:
		}
  }

	chk(e.Stop())
}

type echo struct {
	*portaudio.Stream
	buffer []float32
	noiseMean []float64
	SampleRate  int
	i      int
	step   int
	hm   int
}

func newEcho(delay time.Duration) *echo {
	h, err := portaudio.DefaultHostApi()
	chk(err)
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
	p.Input.Channels = 1
	p.Output.Channels = 1
	// e := &echo{buffer: make([]float32, int(p.SampleRate))}
	e := &echo{buffer: make([]float32, int(p.SampleRate * delay.Seconds()) )}
	e.Stream, err = portaudio.OpenStream(p, e.processAudio)
	e.SampleRate = int(p.SampleRate)
	e.noiseMean = make([]float64, 0)
	e.hm = floorDiv(int(e.SampleRate * 20), 1000)

	chk(err)
	return e
}

func sum(arr []float32) float32 {
	var s float32 = 0.0

	for _, v := range arr {
		s = s + v
	}

	return s
}

func sum64(arr []float64) float64 {
	var s float64 = 0.0

	for _, v := range arr {
		s = s + v
	}

	return s
}


func (e *echo) processAudio(in, out []float32) {
	// fmt.Println(len(out))
	// oldBuffer := make([]float64, len(e.buffer))
	//
	// for i, v := range e.buffer {
	// 	oldBuffer[i] = float64(v)
	// }

	// newBuffer := ReduceNoise(oldBuffer, e.SampleRate)
	// oldBuffer := make([]float64, len(e.buffer))
	//
	// for i, v := range e.buffer {
	// 	oldBuffer[i] = float64(v)
	// }
	//
	// newBuffer := ReduceNoise(oldBuffer, e.SampleRate, noiseMean)

  tmpOut := make([]float64, len(out))
	indexList := make([]int, 0)
	for i := range out {
		indexList = append(indexList, e.i)
		v := float64(e.buffer[e.i])
    tmpOut[i] = v
		e.buffer[e.i] = in[i]
		e.i = (e.i + 1) % len(e.buffer)
		// fmt.Println(e.i)
	}

	// fmt.Print	ln(len(in))

	// fmt.Println(sum(in))

	if sum(in) !=  0 {
		e.step = e.step + len(in)
	}

	if (e.step >= e.hm) && len(e.noiseMean) == 0 {
		oldBuffer := make([]float64, len(e.buffer))

		for i, v := range e.buffer {
			oldBuffer[i] = float64(v)
		}

		// fmt.Println(oldBuffer[e.i - e.step + 10:])

		// fmt.Println(oldBuffer[e.i:])
		e.noiseMean = CalNoiseMean(oldBuffer[e.i - e.step + 10:], e.SampleRate)
		// fmt.Println(e.noiseMean)
	}

	for i := range out {
		out[i] = float32(float64(scale) * tmpOut[i])
	}

	if len(e.noiseMean) > 0 {
		oldBuffer := make([]float64, len(e.buffer))

		for i, v := range e.buffer {
			oldBuffer[i] = float64(v)
		}

		// ReduceNoise(oldBuffer, e.SampleRate, e.noiseMean)
	// 	// newBuffer := ReduceNoise(oldBuffer, e.SampleRate, e.noiseMean)
	// 	// fmt.Println(indexList)
	// 	// fmt.Println(len(oldBuffer), len(newBuffer))
	//
	// 	// fmt.Println(len(indexList))
	//
	//
	// 	for i, v := range indexList {
	// 		if v < len(newBuffer) {
	// 			fmt.Println("out:", out[i], "org:", e.buffer[v] , "old:", oldBuffer[v], "new:", newBuffer[v])
	// 			// out[i] = float32(float64(scale) * oldBuffer[v])
	// 		} else {
	// 			// out[i] = float32(tmpOut[i])
	// 		}
	// 	}
	// 	// fmt.Println(sum64(oldBuffer), sum64(newBuffer))
	//
	// 	return
	}

}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
