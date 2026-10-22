//go:build !no_face_detection && faces_exif

package face_detection

import (
	"errors"
	"log"
	"strings"

	"github.com/barasher/go-exiftool"
	"github.com/photoview/photoview/api/dataloader"
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/utils"
	"gorm.io/gorm"
)

// FromMetadata reports that faces are read from image metadata rather than detected.
// The scanner then reads them when a file is found or changed, see scanner_tasks.RatingTask.
const FromMetadata = true

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

// DetectFaces replaces the faces of the given image with the named face regions in its metadata
func (fd *faceDetector) DetectFaces(db *gorm.DB, media *models.Media) error {
	if err := db.Where("media_id = ?", media.ID).Delete(&models.ImageFace{}).Error; err != nil {
		return err
	}

	// Only new or changed files get here, so never use a cached result
	fd.dataLoader.Clear(media.Path)
	fileInfo, err := fd.dataLoader.Load(media.Path)
	if err != nil {
		return err
	}

	orientation := 1
	if media.ExifID != nil {
		var exif models.MediaEXIF
		if err := db.Select("orientation").First(&exif, *media.ExifID).Error; err != nil {
			return err
		}
		if exif.Orientation != nil {
			orientation = int(*exif.Orientation)
		}
	}

	for _, face := range ParseFaceRegions(fileInfo.Fields, orientation) {
		imageFace := models.ImageFace{
			MediaID: media.ID,
			Rectangle: models.FaceRectangle{
				MinX: face.X - face.Width/2,
				MaxX: face.X + face.Width/2,
				MinY: face.Y - face.Height/2,
				MaxY: face.Y + face.Height/2,
			},
		}

		var faceGroup models.FaceGroup
		if err := db.First(&faceGroup, "label = ?", face.Name).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			faceGroup = models.FaceGroup{
				Label:        &face.Name,
				DisplayLabel: &strings.Fields(face.Name)[0],
				ImageFaces:   []models.ImageFace{imageFace},
			}
			if err := db.Create(&faceGroup).Error; err != nil {
				return err
			}
			continue
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
