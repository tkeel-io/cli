package downloadutil

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Download(filepath string, url string) error {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return fmt.Errorf("http get err:%w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("create file err:%w", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		err = fmt.Errorf("io err:%w", err)
	}

	return err
}
