package golangcodes

import (
	"fmt"
	"net/http"
	"os"

	reader "github.com/knetic0/go-wget/reader"
)

func DownloadFile(url, fileName string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error handled on getting url: %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%d", response.StatusCode)
	}

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	err = reader.NewProgressReader(response.Body, int(response.ContentLength)).CopyWithProgress(out)
	if err != nil {
		return err
	}

	fmt.Printf("\n%s downloaded successfully!\n", fileName)
	return nil
}
