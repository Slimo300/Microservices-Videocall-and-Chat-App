package handlers_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/textproto"

	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/database/mock"
	"github.com/Slimo300/MicroservicesChatApp/backend/group-service/handlers"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/auth"
	"github.com/Slimo300/MicroservicesChatApp/backend/lib/storage"
)

func setupTestServerWithHub() handlers.Server {
	mockDB := mock.NewMockDB()
	s := handlers.NewServerWithMockHub(mockDB, storage.MockStorage{})
	go s.RunHub()
	return *s
}

func setupTestServer() *handlers.Server {

	mockDB := mock.NewMockDB()
	mockAuthClient := auth.NewMockAuthClient()
	s := handlers.NewServer(mockDB, storage.MockStorage{}, mockAuthClient)
	return s
}

func setupTestServerWithAuthClient(auth auth.TokenClient) *handlers.Server {
	s := handlers.NewServer(mock.NewMockDB(), storage.MockStorage{}, auth)
	return s
}

func createImage() *image.RGBA {
	width := 200
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	return img
}

func createTestFormFile(fileName, cType string) (*bytes.Buffer, *multipart.Writer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			fileName, "img.png"))
	h.Set("Content-Type", cType)
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, nil, err
	}

	if err = png.Encode(part, createImage()); err != nil {
		return nil, nil, err
	}
	writer.Close()
	return body, writer, nil
}
