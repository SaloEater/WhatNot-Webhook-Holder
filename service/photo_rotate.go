package service

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"net/http"
)

type PhotoRotateRequest struct {
	Id      int64 `json:"id"`
	Degrees int   `json:"degrees"`
}

type PhotoRotateResponse struct {
	Url string `json:"url"`
}

func (s *Service) PhotoRotate(req *PhotoRotateRequest) (*PhotoRotateResponse, error) {
	photo, err := s.PhotoRepositorier.GetById(req.Id)
	if err != nil {
		return nil, fmt.Errorf("get photo: %w", err)
	}

	resp, err := http.Get(photo.Url)
	if err != nil {
		return nil, fmt.Errorf("download image: %w", err)
	}
	defer resp.Body.Close()

	src, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	rotated := rotateImage(src, req.Degrees)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, rotated, &jpeg.Options{Quality: 95}); err != nil {
		return nil, fmt.Errorf("encode image: %w", err)
	}

	filename := fmt.Sprintf("photo_%d_r%d.jpg", req.Id, req.Degrees)
	newUrl, err := s.DigitalOceaner.SaveCardPhoto(buf.Bytes(), photo.SeriesId, filename)
	if err != nil {
		return nil, fmt.Errorf("upload image: %w", err)
	}

	if err := s.PhotoRepositorier.UpdateUrl(req.Id, newUrl); err != nil {
		return nil, fmt.Errorf("update url: %w", err)
	}

	return &PhotoRotateResponse{Url: newUrl}, nil
}

func rotateImage(src image.Image, degrees int) image.Image {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	switch degrees {
	case 90:
		dst := image.NewRGBA(image.Rect(0, 0, h, w))
		draw.Draw(dst, dst.Bounds(), image.Transparent, image.Point{}, draw.Src)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.Set(h-1-y, x, src.At(b.Min.X+x, b.Min.Y+y))
			}
		}
		return dst
	case 180:
		dst := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.Set(w-1-x, h-1-y, src.At(b.Min.X+x, b.Min.Y+y))
			}
		}
		return dst
	case 270:
		dst := image.NewRGBA(image.Rect(0, 0, h, w))
		draw.Draw(dst, dst.Bounds(), image.Transparent, image.Point{}, draw.Src)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				dst.Set(y, w-1-x, src.At(b.Min.X+x, b.Min.Y+y))
			}
		}
		return dst
	}
	return src
}
