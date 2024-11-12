package service

import (
	"bytes"
	"fmt"
	"github.com/SaloEater/WhatNot-Webhook-Holder/entity"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/fogleman/gg"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/pkg/errors"
	"io"
	"time"
)

func (s *Service) generateLabels(bundle *entity.Bundle) error {
	bundleBoxes, err := s.BundleBoxesRepositorier.GetAllByBundle(bundle.ID)
	typeIDs := []int64{}
	for _, box := range bundleBoxes {
		unique := true
		for _, id := range typeIDs {
			if id == box.BoxTypeID {
				unique = false
				break
			}
		}
		if unique {
			typeIDs = append(typeIDs, box.BoxTypeID)
		}
	}
	boxTypes, err := s.BoxTypeRepositorier.GetByIDs(typeIDs)
	boxTypeMap := make(map[int64]string, len(boxTypes))
	for _, boxType := range boxTypes {
		boxTypeMap[boxType.ID] = boxType.Name
	}
	if err != nil {
		fmt.Println(errors.WithMessage(err, fmt.Sprintf("failed to get bundle boxes for %d", bundle.ID)))
		return err
	}

	count := calculateCount(bundleBoxes)
	images := make([]io.Reader, count)
	var codeCS barcode.BarcodeIntCS
	var code barcode.Barcode

	boxes := make([]*entity.Box, count)
	counter := 0
	for _, bundleBox := range bundleBoxes {
		for nextIndex := 1; nextIndex <= bundleBox.Count; nextIndex++ {
			locationID := bundle.LocationID.Int64
			bundleID := bundle.ID
			boxTypeID := bundleBox.BoxTypeID
			nextID, err := createID(locationID, bundleID, boxTypeID, nextIndex)
			if err != nil {
				fmt.Println(errors.WithMessage(err, fmt.Sprintf("failed to create ID for %d, %d, %d, %d", locationID, bundleID, boxTypeID, nextIndex)))
				return err
			}
			code, err = createBarcode(codeCS, nextID, code)
			if err != nil {
				fmt.Println(errors.WithMessage(err, fmt.Sprintf("failed to create barcode for ID %s", nextID)))
				return err
			}

			boxTypeName, ok := boxTypeMap[boxTypeID]
			if !ok {
				fmt.Println(errors.WithMessage(err, fmt.Sprintf("failed to get box type name for ID %d", boxTypeID)))
				return err
			}
			dc := createImage(code, locationID, bundleID, boxTypeID, nextIndex, boxTypeName)

			var reader bytes.Buffer
			err = dc.EncodePNG(&reader)
			if err != nil {
				fmt.Println(errors.WithMessage(err, fmt.Sprintf("failed to encode PNG for ID %s", nextID)))
				return err
			}
			images[counter] = &reader
			box := &entity.Box{
				BundleID: bundle.ID,
				Status:   entity.BoxStatusPlanned,
				BoxesID:  bundleBox.ID,
				Index:    nextIndex,
			}
			boxes[counter] = box
			counter++
		}
	}

	imp := &pdfcpu.Import{
		PageDim: &types.Dim{
			Width:  164,
			Height: 85,
		},
		PageSize: "A4",
		Pos:      types.Full,
		Scale:    1,
		InpUnit:  types.POINTS,
	}
	pdfBuffer := new(bytes.Buffer)
	err = api.ImportImages(nil, pdfBuffer, images, imp, nil)
	if err != nil {
		panic(err)
	}

	timestamp := time.Now().Format("20060102_150405")
	url, err := s.DigitalOceaner.SaveLabel(*pdfBuffer, fmt.Sprintf("%s_bundle_%d_labels.pdf", timestamp, bundle.ID))
	labels := &entity.BundleLabels{
		BundleID: bundle.ID,
		URL:      url,
	}
	err = s.BundleLabelRepositorier.Create(labels)
	if err != nil {
		fmt.Println(errors.WithMessage(err, fmt.Sprintf("failed to create labels for bundle %d", bundle.ID)))
		return err
	}
	err = s.BoxRepositorier.BatchCreate(boxes, labels.ID)
	if err != nil {
		fmt.Println(errors.WithMessage(err, fmt.Sprintf("failed to create boxes for bundle %d", bundle.ID)))
		return err
	}
	return nil
}

func createImage(code barcode.Barcode, locationID int64, bundleID int64, boxTypeID int64, nextIndex int, boxName string) *gg.Context {
	width := code.Bounds().Size().X
	height := code.Bounds().Size().Y + 50
	dc := gg.NewContext(width, height)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.DrawImage(code, 0, 15)

	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(boxName, float64(width/2), 0, 0.5, 1)

	trackingNumber := fmt.Sprintf("%03d %07d %03d %04d", locationID, bundleID, boxTypeID, nextIndex)
	dc.DrawString(trackingNumber, 15, 60)

	trackingDescription := fmt.Sprintf(" lc    bndl  tp   bx")
	dc.DrawString(trackingDescription, 15, 70)
	return dc
}

func createBarcode(codeCS barcode.BarcodeIntCS, nextID string, code barcode.Barcode) (barcode.Barcode, error) {
	codeCS, err := code128.Encode(nextID)
	if err != nil {
		fmt.Println(errors.WithMessage(err, fmt.Sprintf("failed to encode ID %s", nextID)))
		return nil, err
	}

	code, err = barcode.Scale(codeCS, 164, 35)
	if err != nil {
		fmt.Println(errors.WithMessage(err, fmt.Sprintf("failed to scale barcode for ID %s", nextID)))
		return nil, err
	}
	return code, nil
}

func validateParameter(name string, value, max int) error {
	if value < 0 || value > max {
		return fmt.Errorf("%s must be between 0 and %d", name, max)
	}
	return nil
}

func createID(locationID int64, bundleID int64, boxTypeID int64, boxIndex int) (string, error) {
	if err := validateParameter("locationID", int(locationID), 999); err != nil {
		return "", err
	}
	if err := validateParameter("bundleID", int(bundleID), 9999999); err != nil {
		return "", err
	}
	if err := validateParameter("boxTypeID", int(boxTypeID), 999); err != nil {
		return "", err
	}
	if err := validateParameter("boxIndex", boxIndex, 9999); err != nil {
		return "", err
	}

	id := fmt.Sprintf("%03d%07d%03d%04d", locationID, bundleID, boxTypeID, boxIndex)
	return id, nil
}

func calculateCount(boxes []*entity.BundleBoxes) int {
	count := 0
	for _, box := range boxes {
		count += box.Count
	}
	return count
}
