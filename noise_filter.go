package main


import (
  "github.com/mjibson/go-dsp/fft"
  "github.com/mjibson/go-dsp/window"

  // "github.com/youpy/go-wav"
	// "io"

  // "os"
  // "fmt"
  "math"
  "math/cmplx"
)

func floorDiv(f1 int, f2 int) int {
  return int(math.Floor(float64(f1)/float64(f2)))
}

func sumFloat64(arr []float64) float64 {
  sum := 0.0
  for _, v := range arr {
    sum += v
  }

  return sum
}

func nextpow2(x int) int {
  i := 1
  v := 2
  for {
    v = v * 2
    i += 1

    if v > x {
      break
    }
  }

  return i
}

func pow(f1 int, f2 int) int {
  return int(math.Pow(float64(f1), float64(f2)))
}

func powf(f1 float64, f2 int) float64 {
  return math.Pow(float64(f1), float64(f2))
}

func matMulC(arr1 []float64, fac float64) ([]float64) {
  re := make([]float64, 0)

  for _, v := range arr1 {
    re = append(re, v * fac)
  }

  return re
}

func matMul(arr1 []float64, arr2 []float64) ([]float64) {
  re := make([]float64, 0)

  l2 := len(arr2)

  for i, v := range arr1 {
    v1 := 0.0
    if i < l2 {
      v1 = arr2[i]
    }

    re = append(re, v * v1)
  }

  return re
}

func matMul2(arr1 []float64, arr2 []int) ([]float64) {
  re := make([]float64, 0)

  l2 := len(arr2)

  for i, v := range arr1 {
    v1 := 0.0
    if i < l2 {
      v1 = float64(arr2[i])
    }

    re = append(re, v * v1)
  }

  return re
}

func matPadding(arr1 []float64, len int) []float64 {
  re := make([]float64, len)

  for i, v := range arr1 {
    re[i] = v
  }

  return re
}

func absC(c complex128) float64 {
  return math.Sqrt(real(c) * real(c) + imag(c) * imag(c))
}

func matAbs(arr1 []complex128) []float64 {
  re := make([]float64, len(arr1))

  for i, c := range arr1 {
    re[i] = absC(c)
  }

  return re
}

func matAdd(arr1 []float64, arr2 []float64) []float64 {
  re := make([]float64, len(arr1))

  l2 := len(arr2)
  for i, v := range arr1 {
    v1 := 0.0
    if i < l2 {
      v1 = arr2[i]
    }

    re[i] = v + v1
  }

  return re
}

func matSub(arr1 []float64, arr2 []float64) []float64 {
  re := make([]float64, len(arr1))

  l2 := len(arr2)
  for i, v := range arr1 {
    v1 := 0.0
    if i < l2 {
      v1 = arr2[i]
    }

    re[i] = v - v1
  }

  return re
}

func matDiv(arr1 []float64, fac float64) []float64 {
  re := make([]float64, len(arr1))

  for i, v := range arr1 {
    re[i] = v/fac
  }

  return re
}

func matPow(arr1 []float64, fac float64) []float64 {
  re := make([]float64, len(arr1))

  for i, v := range arr1 {
    re[i] = math.Pow(v, fac)
  }

  return re
}

func matAngle(arr1 []complex128) []float64 {
  re := make([]float64, len(arr1))

  for i, c := range arr1 {
    re[i] = cmplx.Phase(c)
  }

  return re
}

func matNorm2(arr1 []float64) float64 {
  re := 0.0

  for _, v := range arr1 {
    re = re + v * v
  }

  return math.Sqrt(re)
}

func matReplace(arr1 []float64, indexList []int, arr2 []float64) []float64 {
  for _, v := range indexList {
    arr1[v] = arr2[v]
  }

  return arr1
}

func berouti(snr float64) float64 {
  a := 0.0
  if snr >= -5.0 && snr <= 20.0 {
    a = 4.0 - snr * 3 / 20
  } else {
    if snr < -5.0 {
      a = 5
    }

    if snr > 20 {
      a = 1
    }
  }

  return a
}

func berouti1(snr float64) float64 {
  a := 0.0
  if snr >= -5.0 && snr <= 20.0 {
    a = 3.0 - snr * 2.0 / 20.0
  } else {
    if snr < -5.0 {
      a = 4
    }

    if snr > 20 {
      a = 1
    }
  }

  return a
}

func matFlipud(arr []float64) []float64 {
  l := len(arr)
  re := make([]float64, l)
  for i, _ := range arr {
    re[i] = arr[l-i-1]
  }

  return re
}

// func main() {
//   f, err := os.Open("input_file.wav")
//   if err != nil {
//       fmt.Println(err)
//   }
//   defer f.Close()
//   reader := wav.NewReader(f)
//   pms, err := reader.Format()
//
//   x := make([]int, 0)
//
//   for {
//     samples, err := reader.ReadSamples()
//     if err == io.EOF {
//       break
//     }
//
//     for _, sample := range samples {
//       l := reader.IntValue(sample, 0)
//       // r := reader.IntValue(sample, 1)
//
//       x = append(x, l)
//     }
//   }
//
//   newData := ReduceNoise(x, pms)
//
//
//   of, err := os.Create("output_file.wav")
//   if err != nil {
//       fmt.Println(err)
//   }
//   defer of.Close()
//
//   // func NewWriter(w io.Writer, numSamples uint32, numChannels uint16, sampleRate uint32, bitsPerSample uint16) (writer *Writer)
//
//   writer := wav.NewWriter(of, uint32(len(newData)), 1, pms.SampleRate, pms.BitsPerSample)
//   samples := make([]wav.Sample, 0)
//   for _, v := range newData {
//     samples = append(samples, wav.Sample{ Values: [2]int{ int(v), 0 }})
//   }
//   writer.WriteSamples(samples)
// }

func CalNoiseMean(x []float64, fs int) []float64 {
  len_ := floorDiv(int(fs * 20), 1000)
  // fmt.Println(len_)
  win := window.Hamming(len_)
  nFFT := 2 * pow(2, nextpow2(len_))
  noiseMean := make([]float64, nFFT)

  j := 0
  for i := 1; i < 6; i++ {
    noiseMean = matAdd(noiseMean, matAbs( fft.FFTReal( matPadding( matMul(win, x[j:j+len_]), nFFT ) ) ) )
    j = j + len_
  }

  return matDiv(noiseMean, 5.0)
}

func ReduceNoise(x []float64, fs int, noiseMu[]float64) []float64 {
  len_ := floorDiv(int(fs * 20), 1000)
  // fmt.Println(len_)
  PERC := 50
  len1 := floorDiv(int(len_ * PERC), 100)  // 100  # 重叠窗口
  len2 := len_ - len1
  // fmt.Println(len1, len2)

  win := window.Hamming(len_)
  // fmt.Println(win)
  winGain := float64(len2) / sumFloat64(win)
  // # 设置默认参数
  Thres := -4.0 // 3
  Expnt := 2.0  //2.0
  beta := 0.002 //0.002
  G := 0.9 // 0.9

  nFFT := 2 * pow(2, nextpow2(len_))
  // fmt.Println(nFFT)
  // fmt.Println(win)
  // fmt.Println(winGain)

  // noiseMean := make([]float64, nFFT)
  //
  // j := 0
  // for i := 1; i < 6; i++ {
  //   // fmt.Println(len(matPadding( matMul2(win, x[j:j+len_]), nFFT )))
  //   // fmt.Println( fft.FFTReal( matPadding( matMul2(win, x[j:j+len_]), nFFT ) ) )
  //
  //   noiseMean = matAdd(noiseMean, matAbs( fft.FFTReal( matPadding( matMul(win, x[j:j+len_]), nFFT ) ) ) )
  //   // fmt.Println(fft.FFTReal( matPadding( matMul2(win, x[j:j+len_]), nFFT ) ))
  //   // noise_mean = noise_mean + abs(np.fft.fft(win * x[j:j + len_], nFFT))
  //   j = j + len_
  // }
  //
  // noiseMu := matDiv(noiseMean, 5.0)
  // fmt.Println(noiseMu)
  // noise_mu = noise_mean / 5
  k := 1
  // img := 1j
  xOld := make([]float64, len1)
  // x_old = np.zeros(len1)
  Nframes := floorDiv(len(x), len2) - 1
  // xfinal = np.zeros(Nframes * len2)
  xfinal := make([]float64, Nframes * len2)

  for i := 0; i < Nframes; i++ {
    insign := matMul(win, x[k-1:k + len_ - 1])
    // fmt.Println(insign)
    spec := fft.FFTReal( matPadding(insign, nFFT) )
    // fmt.Println(spec)
    sig := matAbs(spec)

    theta := matAngle(spec)
    // fmt.Println(theta)

    // fmt.Println(matNorm2(sig))

    SNRseg := math.Log10( powf( matNorm2(sig), 2) / powf ( matNorm2(noiseMu), 2) ) * 10
    alpha := 0.0
    if Expnt == 1.0 {
      alpha = berouti1(SNRseg)
    } else {
      alpha = berouti(SNRseg)
    }

    subSpeech := matSub(matPow(sig, Expnt), matMulC(matPow(noiseMu, Expnt), alpha))
    // fmt.Println(subSpeech)
    diffw := matSub( subSpeech, matMulC( matPow(noiseMu, Expnt), beta) )


    z := findIndex(diffw)

    if len(z) > 0 {
      subSpeech = matReplace(subSpeech, z, matMulC( matPow(noiseMu, Expnt), beta ))
    }

    noiseTemp := make([]float64, 0)
    // fmt.Println(SNRseg, Thres)
    if SNRseg < float64(Thres) {
      // fmt.Println(SNRseg)
      noiseTemp = matAdd(matMulC( matPow(noiseMu, Expnt) , G), matMulC( matPow(sig, Expnt), (1 - G)))  //平滑处理噪声功率谱
      noiseMu = matPow(noiseTemp, (1.0 / Expnt)) //  # 新的噪声幅度谱
    }

    subSpeech = matReplaceRange(subSpeech, floorDiv(nFFT, 2)+1, nFFT, matFlipud(subSpeech[1:floorDiv(nFFT, 2)]))
    // subSpeech[floorDiv(nFFT, 2)+1:nFFT] = matFlipud(subSpeech[1:floorDiv(nFFT, 2)])
    // sub_speech[nFFT // 2 + 1:nFFT] = np.flipud(sub_speech[1:nFFT // 2])

    xPhase := make([]complex128, 0)

    cosTheta := make([]float64, len(theta))
    sinTheta := make([]float64, len(theta))

    for j, v := range theta {
      cosTheta[j] = math.Cos(v)
      sinTheta[j] = math.Sin(v)
    }

    rPart := matMul( matPow(subSpeech, 1.0/Expnt), cosTheta )
    iPart := matMul( matPow(subSpeech, 1.0/Expnt), sinTheta )

    for j, v := range rPart {
      xPhase = append(xPhase, complex(v, iPart[j]))
    }

    xi := matReal(fft.IFFT(xPhase))

    xfinal = matReplaceRange(xfinal, k-1, k + len2 - 1, matAdd(xOld, xi[0:len1]))
    // xfinal[k-1:k + len2 - 1] = x_old + xi[0:len1]

    // if i == Nframes - 1 {
    //   fmt.Println(Nframes, i, k, len2, k-1, k + len2 - 1, xfinal[k-1:k + len2 - 1])
    // }
    xOld = xi[0 + len1:len_]
    k = k + len2
  }

  newData := matMulC(xfinal, winGain)

  return newData
  // fmt.Println(newData)
  // fmt.Println("total:", x)

}

func matReplaceRange(arr []float64, s int, e int, arr2 []float64) []float64 {
  j := 0
  for i := 0; i < len(arr); i++ {
    if i >= s && i < e {
      arr[i] = arr2[j]

      j += 1
    }
  }

  return arr
}

func matReal(arr []complex128) []float64 {
  re := make([]float64, len(arr))

  for i, c := range arr {
    re[i] = real(c)
  }

  return re
}

func findIndex(arr []float64) []int {
  indexList := make([]int, 0)
  for i, v := range arr {
    if v < 0 {
      indexList = append(indexList, i)
    }
  }

  return indexList
}
