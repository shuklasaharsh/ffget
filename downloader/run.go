package downloader

import (
	"ffget/constants"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func (c *Client) Run() error {
	var wg sync.WaitGroup
	if c.Auto {
		// Get headers with size and make chunks
		// Create a new http client
		client := &http.Client{}
		// Create a head request to the link
		headers, err := client.Head(c.Link)
		if err != nil {
			return err
		}
		// Get the content length
		c.contentLength = headers.ContentLength
		log.Printf("Content length: %d\n", c.contentLength)
		// Create batches
		// Depending on the total length of the file, we will create a batch of 100MB
		c.makePairs()
		// Download Range pairs are done
		// Now we will download the file
		wg.Add(1)
		for index, eachPair := range c.downloadRangePairs {
			go func(byteRangeToDownload [2]int64, index int) {
				defer wg.Done()
				// For this pair
				req, err := http.NewRequest("GET", c.Link, nil)
				if err != nil {
					log.Println(err.Error())
					return
				}
				// Set the range header
				headerValue := fmt.Sprintf("bytes=%d-%d", byteRangeToDownload[0], byteRangeToDownload[1])
				req.Header.Add("Range", headerValue)
				resp, err := client.Do(req)
				defer resp.Body.Close()
				if err != nil {
					log.Printf("ERROR | %s\n", err.Error())
					return
				}
				if resp.StatusCode != 200 {
					log.Printf("ERROR | %s\n", resp.Status)
					return
				}
				// Save the packet
				err = c.savePacket(&resp.Body, index)
				if err != nil {
					log.Println(err.Error())
				}
			}(eachPair, index)
		}
		wg.Wait()
	}
	return nil
}

func (c *Client) savePacket(body *io.ReadCloser, index int) error {
	data, err := io.ReadAll(*body)
	if err != nil {
		return err
	}
	log.Printf("Saving packet %d\n", index)
	if fileInfo, err := os.Stat("/data"); err != nil {
		if os.IsNotExist(err) {
			// Create the directory
			err := os.Mkdir("/data", 0777)
			if err != nil {
				return err
			}
		}
	} else {
		log.Println("saving to : ", fileInfo.Name())
	}
	err = os.WriteFile(fmt.Sprintf("/data/%d", index), data, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) makePairs() {
	checkContentLength := int64(0)
	for checkContentLength < c.contentLength {
		// We divide batches by 100MB
		c.downloadRangePairs = append(c.downloadRangePairs, [2]int64{checkContentLength, checkContentLength + 100*constants.MegaByte})
		checkContentLength += 100 * constants.MegaByte
	}
	c.Parts = len(c.downloadRangePairs)
	return
}

func (c *Client) RunAuto() error {

	return nil
}
