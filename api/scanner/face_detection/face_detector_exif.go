//go:build !no_face_detection && faces_exif

package face_detection

import (
	"errors"
	"log"

	"github.com/barasher/go-exiftool"
	"github.com/photoview/photoview/api/dataloader"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/utils"
	"gorm.io/gorm"
)

type faceDetector struct {
	et         *exiftool.Exiftool
	dataLoader *dataloader.ExiftoolLoader
}

func InitializeFaceDetector(db *gorm.DB) error {
	if utils.EnvDisableFaceRecognition.GetBool() {
		log.Printf("Face detection disabled (%s=1)\n", utils.EnvDisableFaceRecognition.GetName())
		return nil
	}

	log.Println("Initializing face extractions from EXIF")

	buf := make([]byte, 256*1024)

	et, err := exiftool.NewExiftool(exiftool.NoPrintConversion(), exiftool.Buffer(buf, 64*1024))

	if err != nil {
		log.Printf("Error initializing ExifTool: %s\n", err)
		return err
	}

	GlobalFaceDetector = &faceDetector{
		et:         et,
		dataLoader: dataloader.NewExiftoolLoader(et),
	}

	return nil
}

// ReloadFacesFromDatabase replaces the in-memory face descriptors with the ones in the database
func (fd *faceDetector) ReloadFacesFromDatabase(db *gorm.DB) error {
	return nil
}

type FaceRegion struct {
	Name   string
	X, Y   float64 // top-left
	Width  float64
	Height float64
}

func toStringSlice(raw interface{}) []string {
	var result []string
	switch v := raw.(type) {
	case string:
		result = []string{v}
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
	}
	return result
}

func toFloatSlice(raw interface{}) []float64 {
	var result []float64
	switch v := raw.(type) {
	case float64:
		result = []float64{v}
	case []interface{}:
		for _, item := range v {
			if f, ok := item.(float64); ok {
				result = append(result, f)
			}
		}
	}
	return result
}

// DetectFaces finds the faces in the given image and saves them to the database
func (fd *faceDetector) DetectFaces(db *gorm.DB, media *models.Media) error {
	if err := db.Model(media).Preload("MediaURL").First(&media).Error; err != nil {
		return err
	}

	log.Printf("Looking for EXIF faces in %s", media.Path)

	fileInfo, err := fd.dataLoader.Load(media.Path)
	if err != nil {
		return err
	}

	// Collect faces
	var faces []FaceRegion

	// exiftool often reports regions as indexed fields:
	// RegionName, RegionType, RegionAreaX, RegionAreaY, RegionAreaW, RegionAreaH
	names := toStringSlice(fileInfo.Fields["RegionName"])
	types := toStringSlice(fileInfo.Fields["RegionType"])
	xs := toFloatSlice(fileInfo.Fields["RegionAreaX"])
	ys := toFloatSlice(fileInfo.Fields["RegionAreaY"])
	ws := toFloatSlice(fileInfo.Fields["RegionAreaW"])
	hs := toFloatSlice(fileInfo.Fields["RegionAreaH"])

	// Sanity check lengths
	for i := 0; i < len(types); i++ {
		if types[i] != "Face" {
			continue
		}
		// Convert coordinates from strings to float64

		name := ""
		if i < len(names) {
			name = names[i]
		}

		faces = append(faces, FaceRegion{
			Name:   name,
			X:      xs[i],
			Y:      ys[i],
			Width:  ws[i],
			Height: hs[i],
		})
	}

	// Print results
	for _, face := range faces {
		log.Printf("Face: %s at (%.3f, %.3f, %.3f, %.3f)\n", face.Name, face.X, face.Y, face.Width, face.Height)
		imageFace := models.ImageFace{
			MediaID: media.ID,
			// Descriptor: models.FaceDescriptor(face.Descriptor), --- IGNORE ---
			Rectangle: models.FaceRectangle{
				MinX: face.X,
				MaxX: face.X + face.Width,
				MinY: face.Y,
				MaxY: face.Y + face.Height,
			},
		}
		var faceGroup models.FaceGroup
		if err := db.First(&faceGroup, "label = ?", face.Name).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				faceGroup = models.FaceGroup{
					Label:      &face.Name,
					ImageFaces: []models.ImageFace{imageFace},
				}
				if err := db.Create(&faceGroup).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		if err := db.Model(&faceGroup).Association("ImageFaces").Append(&imageFace); err != nil {
			return err
		}
	}

	return nil
}

func (fd *faceDetector) MergeCategories(sourceID int32, destID int32) {
}

func (fd *faceDetector) MergeImageFaces(imageFaceIDs []int, destFaceGroupID int32) {
}

func (fd *faceDetector) RecognizeUnlabeledFaces(tx *gorm.DB, user *models.User) ([]*models.ImageFace, error) {
	updatedImageFaces := make([]*models.ImageFace, 0)
	return updatedImageFaces, nil
}
