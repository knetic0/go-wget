package reader

import (
	"fmt"
	"io"
	"sync"
)

const barWidth = 40

type ProgressReader struct {
	mu        sync.RWMutex
	Reader    io.Reader
	TotalSize int
}

func NewProgressReader(reader io.Reader, totalSize int) *ProgressReader {
	return &ProgressReader{
		Reader:    reader,
		TotalSize: totalSize,
	}
}

func (progressReader *ProgressReader) Read(p []byte) (int, error) {
	n, err := progressReader.Reader.Read(p)
	return n, err
}

func (progressReader *ProgressReader) CopyWithProgress(dst io.Writer) error {
	buf := make([]byte, 32*1024)

	var totalRead int
	lastUpdate := -1

	wg := &sync.WaitGroup{}

	for {
		er := make(chan error, 1)

		wg.Add(1)

		go func() {
			progressReader.mu.RLock()
			defer progressReader.mu.RUnlock()
			defer close(er)
			defer wg.Done()

			n, err := progressReader.Read(buf)
			if n > 0 {
				_, err := dst.Write(buf[:n])
				if err != nil {
					er <- err
					return
				}

				totalRead += n
				newProgress := float64(totalRead) / float64(progressReader.TotalSize) * 100

				if int(newProgress) != lastUpdate {
					progressReader.updateProgress(int(newProgress))
					lastUpdate = int(newProgress)
				}

				if err == io.EOF {
					er <- err
					return
				}
			}

			if err != nil && err != io.EOF {
				er <- err
				return
			}
		}()

		wg.Wait()

		if err := <-er; err != nil {
			return fmt.Errorf(err.Error())
		}

		if totalRead == progressReader.TotalSize {
			break
		}
	}

	return nil
}

func (progressReader *ProgressReader) updateProgress(percent int) {
	fmt.Printf("\r[")
	barLength := int(float64(barWidth) * float64(percent) / 100.0)

	for i := 0; i < barWidth; i++ {
		if i < barLength {
			fmt.Print("=")
		} else {
			fmt.Print(" ")
		}
	}

	fmt.Printf("] %d%%", percent)
}
