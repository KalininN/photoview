package scanner_tasks

import (
	"github.com/photoview/photoview/api/graphql/models"
	"github.com/photoview/photoview/api/scanner/face_detection"
	"github.com/photoview/photoview/api/scanner/scanner_task"
	"github.com/photoview/photoview/api/scanner/scanner_utils"
)

type FaceDetectionTask struct {
	scanner_task.ScannerTaskBase
}

func (t FaceDetectionTask) AfterMediaFound(ctx scanner_task.TaskContext, media *models.Media, newMedia bool) error {
	if !newMedia {
		return nil
	}

	if media.Type != models.MediaTypePhoto {
		return nil
	}

	if face_detection.GlobalFaceDetector == nil {
		return nil
	}

	if err := face_detection.GlobalFaceDetector.DetectFaces(ctx.GetDB(), media); err != nil {
		scanner_utils.ScannerError(ctx, "Error detecting faces in image (%s): %s", media.Path, err)
	}
	return nil
}
