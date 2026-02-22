package services

import (
	"image"
	"image/color"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/jaygaha/roadmap-go-projects/intermediate/image-processing-service/internal/models"
)

func newSolidImage(width, height int, c color.Color) image.Image {
	img := imaging.New(width, height, c)
	return img
}

func TestApplyOperationResize(t *testing.T) {
	svc := &ImageService{}
	src := newSolidImage(100, 50, color.RGBA{255, 0, 0, 255})

	dst, contentType, format, err := svc.applyOperation(src, "image/jpeg", string(models.OperationResize), map[string]string{
		"width":  "50",
		"height": "0",
	})
	if err != nil {
		t.Fatalf("applyOperation resize returned error: %v", err)
	}
	if format != imaging.JPEG || contentType != "image/jpeg" {
		t.Fatalf("unexpected format or contentType: got %v, %s", format, contentType)
	}
	b := dst.Bounds()
	if b.Dx() != 50 || b.Dy() != 25 {
		t.Fatalf("unexpected resized dimensions: got %dx%d, want 50x25", b.Dx(), b.Dy())
	}
}

func TestApplyOperationCrop(t *testing.T) {
	svc := &ImageService{}
	src := newSolidImage(100, 100, color.RGBA{0, 255, 0, 255})

	dst, _, _, err := svc.applyOperation(src, "image/png", string(models.OperationCrop), map[string]string{
		"width":  "40",
		"height": "40",
	})
	if err != nil {
		t.Fatalf("applyOperation crop returned error: %v", err)
	}
	b := dst.Bounds()
	if b.Dx() != 40 || b.Dy() != 40 {
		t.Fatalf("unexpected cropped dimensions: got %dx%d, want 40x40", b.Dx(), b.Dy())
	}
}

func TestApplyOperationRotate(t *testing.T) {
	svc := &ImageService{}
	src := newSolidImage(10, 20, color.RGBA{0, 0, 255, 255})

	dst, _, _, err := svc.applyOperation(src, "image/png", string(models.OperationRotate), map[string]string{
		"angle": "90",
	})
	if err != nil {
		t.Fatalf("applyOperation rotate returned error: %v", err)
	}
	b := dst.Bounds()
	if b.Dx() != 20 || b.Dy() != 10 {
		t.Fatalf("unexpected rotated dimensions: got %dx%d, want 20x10", b.Dx(), b.Dy())
	}
}

func TestApplyOperationFlipHorizontal(t *testing.T) {
	svc := &ImageService{}
	img := imaging.New(2, 1, color.RGBA{0, 0, 0, 0})
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	img.Set(1, 0, color.RGBA{0, 0, 255, 255})

	dst, _, _, err := svc.applyOperation(img, "image/png", string(models.OperationFlip), map[string]string{
		"mode": "horizontal",
	})
	if err != nil {
		t.Fatalf("applyOperation flip returned error: %v", err)
	}
	left := color.RGBAModel.Convert(dst.At(0, 0)).(color.RGBA)
	right := color.RGBAModel.Convert(dst.At(1, 0)).(color.RGBA)
	if left.R != 0 || left.B != 255 || right.R != 255 || right.B != 0 {
		t.Fatalf("unexpected flipped colors: left=%v right=%v", left, right)
	}
}

func TestApplyOperationGrayscale(t *testing.T) {
	svc := &ImageService{}
	src := newSolidImage(5, 5, color.RGBA{200, 100, 50, 255})

	dst, _, _, err := svc.applyOperation(src, "image/jpeg", string(models.OperationGrayscale), map[string]string{})
	if err != nil {
		t.Fatalf("applyOperation grayscale returned error: %v", err)
	}
	c := color.RGBAModel.Convert(dst.At(2, 2)).(color.RGBA)
	if !(c.R == c.G && c.G == c.B) {
		t.Fatalf("expected grayscale pixel with equal RGB, got %v", c)
	}
}

func TestApplyOperationCompressKeepsDimensions(t *testing.T) {
	svc := &ImageService{}
	src := newSolidImage(30, 40, color.RGBA{10, 20, 30, 255})

	dst, contentType, format, err := svc.applyOperation(src, "image/jpeg", string(models.OperationCompress), map[string]string{
		"format":  "jpeg",
		"quality": "70",
	})
	if err != nil {
		t.Fatalf("applyOperation compress returned error: %v", err)
	}
	if format != imaging.JPEG || contentType != "image/jpeg" {
		t.Fatalf("unexpected format or contentType for compress: got %v, %s", format, contentType)
	}
	if dst.Bounds() != src.Bounds() {
		t.Fatalf("compress should preserve dimensions: got %v, want %v", dst.Bounds(), src.Bounds())
	}
}
