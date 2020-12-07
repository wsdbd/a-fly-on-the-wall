package main

import (
	"github.com/gordonklaus/portaudio"
  "github.com/mjibson/go-dsp/fft"
  "github.com/joho/godotenv"

	"time"
  "fmt"
  "math"
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
	i      int
}

func newEcho(delay time.Duration) *echo {
	h, err := portaudio.DefaultHostApi()
	chk(err)
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
	p.Input.Channels = 1
	p.Output.Channels = 1
	e := &echo{buffer: make([]float32, int(p.SampleRate*delay.Seconds()))}
	e.Stream, err = portaudio.OpenStream(p, e.processAudio)
	chk(err)
	return e
}

func (e * echo) max(cmlx []complex128) (float64, int) {
  n := float64(len(cmlx))
  pos := 0
  var max float64 = 0.0
  for i, c := range cmlx {
    r := real(c)
    im := imag(c)

    v := math.Sqrt(r * r + im * im)
    if i == 0 {
      v = v/n
    } else {
      v = v * 2 / n
    }

    if v > max {
      max = v
      pos = i
    }

  }

  return max, pos
}

func (e * echo) most(cmlx []complex128) ([]int) {
  n := float64(len(cmlx))
  arr := make([]int, 0)
  var max float64 = 0.0
  for i, c := range cmlx {
    v := e.mod(c)
    if i == 0 {
      v = v/n
    } else {
      v = v * 2 / n
    }

    if v > max {
      max = v
    }
  }


  for i, c := range cmlx {
    v := e.mod(c)

    if i == 0 {
      v = v/n
    } else {
      v = v * 2 / n
    }

    if v > max * mostFactor  {
      arr = append(arr, i)
    }
  }

  return arr
}



func (e * echo) freq(cmlx []complex128) []float64 {
  n := float64(len(cmlx))
  arr := make([]float64, len(cmlx))
  for i, c := range cmlx {
    r := real(c)
    im := imag(c)

    v := math.Sqrt(r * r + im * im)
    if i == 0 {
      v = v/n
    } else {
      v = v * 2 / n
    }
    arr[i] = v
  }

  return arr
}

func (e * echo) mod(c complex128) float64 {
  return math.Sqrt(real(c) * real(c) + imag(c) * imag(c))
}

func (e * echo) isNoise(cmlx []complex128) bool {
  total := 0.0
  for _, c := range cmlx {
    m := e.mod(c)

    total += m
  }

  avg := total/float64(len(cmlx))

  fx := 0.0
  for _, c := range cmlx {
    m := e.mod(c)

    fx += (m - avg) * (m - avg)
  }


  if fx < 0.0001 {
    return true
  }

  return false
}

func inArr(arr []int, m int) bool {
  for _, v := range arr {
    if m == v {
      return true
    }
  }

  return false
}

func isNoise(arr []int, n int) bool {
  newArr := make([]int, 0)

  for _, v := range arr {
    if v != 0 && v < n - 5 && v > 5 {
      newArr = append(newArr, v)
    }
  }

  if len(newArr) > 0 {
    return false
  }

  return false
}

func (e * echo) copy(cmlx []complex128) []complex128 {
  tmp := make([]complex128, 0)

  for _, v := range cmlx {
    tmp = append(tmp, v)
  }

  return tmp
}

var noise []complex128 = make([]complex128, 0)

func (e * echo) filter(cmlx []complex128) []complex128 {
  pos := e.most(cmlx)
  flag := isNoise(pos, len(cmlx))
  // max, _ := e.max(cmlx)
  if (e.isNoise(cmlx)) {
    noise = e.copy(cmlx)
  }
  if !flag {
    // fmt.Println(pos)
  }

  result := make([]complex128, len(cmlx))
  for i, c := range cmlx {
    if i <= 10 && i > len(cmlx) - 10 - 1 {
      result[i] = 0
    } else {
      result[i] = c
    }
  }

  return result
}

func (e * echo) inverse(cmlx []complex128) []float64 {
  t := make([]float64, 0)
  for _, c := range cmlx {
    t = append(t, real(c))
  }

  return t
}

func (e *echo) processAudio(in, out []float32) {
  tmpOut := make([]float64, len(out))
	for i := range out {
    v := float64(e.buffer[e.i])
    tmpOut[i] = v
		e.buffer[e.i] = in[i]
		e.i = (e.i + 1) % len(e.buffer)
	}

  cmplxArr := fft.FFTReal(tmpOut)
  newArr := e.filter(cmplxArr)
  // fmt.Println(newArr)
  ifftResult := fft.IFFT(newArr)
  newXX := e.inverse(ifftResult)

  for i := range out {
    out[i] = float32(float64(scale) * newXX[i])
  }

}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
