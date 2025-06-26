package imageutil

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
)

// CompressOptions 压缩选项
type CompressOptions struct {
	Quality      int    // 压缩质量 (1-100)，仅对JPEG和WebP有效
	MaxWidth     uint   // 最大宽度，如果为0则保持原宽度
	MaxHeight    uint   // 最大高度，如果为0则保持原高度
	Format       string // 输出格式 ("jpeg", "png", "gif", "bmp", "webp")，默认使用原格式
	PreserveName bool   // 是否保留原始文件名
}

// ImageInfo 保存图片相关信息
type ImageInfo struct {
	OriginalSize   int64   // 原始大小(字节)
	CompressedSize int64   // 压缩后大小(字节)
	Ratio          float64 // 压缩比率
	Width          int     // 宽度
	Height         int     // 高度
	Format         string  // 格式
	QualityScore   int     // 质量评分 (1-100)
	IsAnimated     bool    // 是否是动画图片
	FrameCount     int     // 动画帧数
}

// CompressImage 压缩图片
// 参数:
//   - file: 上传的图片文件
//   - options: 压缩选项
//
// 返回值:
//   - 压缩后的图片数据
//   - 图片信息
//   - 错误
func CompressImage(file *multipart.FileHeader, options CompressOptions) ([]byte, *ImageInfo, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 读取文件数据
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, nil, fmt.Errorf("读取文件数据失败: %w", err)
	}

	// 检查是否是GIF并预处理
	fileType := detectFileType(fileBytes)

	// 创建图片信息
	imgInfo := &ImageInfo{
		OriginalSize: file.Size,
		Format:       fileType,
		QualityScore: estimateQualityScore(fileBytes, fileType),
	}

	// 设置输出格式
	// 如果用户明确指定了格式，则使用指定的格式，否则保持原格式
	outputFormat := strings.ToLower(options.Format)
	if outputFormat == "" {
		outputFormat = fileType // 保持原格式
	}

	// 对GIF特殊处理以保留动画特性
	if fileType == "gif" && (outputFormat == "gif" || outputFormat == "") {
		return compressGif(fileBytes, imgInfo, options)
	}

	// 处理非GIF格式或需要转换的情况
	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("解码图片失败，可能不是有效的图片格式: %w", err)
	}

	// 补充图片信息
	if imgInfo.Width == 0 {
		imgInfo.Width = img.Bounds().Dx()
	}
	if imgInfo.Height == 0 {
		imgInfo.Height = img.Bounds().Dy()
	}

	// 设置自动调整默认压缩质量
	compressionQuality := options.Quality
	if compressionQuality <= 0 || compressionQuality > 100 {
		// 默认根据文件大小设置压缩质量
		if file.Size > 5*1024*1024 {
			compressionQuality = 60 // 大文件使用更强压缩
		} else if file.Size > 1024*1024 {
			compressionQuality = 70 // 中等文件
		} else {
			compressionQuality = 75 // 小文件
		}
	}

	// 应用尺寸调整以减小文件大小
	isResized := false
	if options.MaxWidth > 0 || options.MaxHeight > 0 {
		img, isResized = resizeImage(img, options.MaxWidth, options.MaxHeight)
		imgInfo.Width = img.Bounds().Dx()
		imgInfo.Height = img.Bounds().Dy()
	} else {
		// 即使未指定尺寸限制，对大图片也进行适当缩小
		if imgInfo.Width > 3000 || imgInfo.Height > 3000 {
			// 对超大图片进行适度缩小（保持比例）
			scaleFactor := 0.7 // 缩小30%
			newWidth := uint(float64(imgInfo.Width) * scaleFactor)
			newHeight := uint(float64(imgInfo.Height) * scaleFactor)
			img, isResized = resizeImage(img, newWidth, newHeight)
			imgInfo.Width = img.Bounds().Dx()
			imgInfo.Height = img.Bounds().Dy()
		}
	}

	// 对于大尺寸图片提供额外的压缩
	if imgInfo.Width*imgInfo.Height > 4000000 { // 400万像素以上
		// 减少压缩质量
		compressionQuality = int(float64(compressionQuality) * 0.8) // 降低20%质量
		if compressionQuality < 40 {
			compressionQuality = 40 // 保证最低质量
		}
	}

	// 应用图像处理优化
	if outputFormat == "jpeg" || outputFormat == "jpg" {
		// 对于JPEG格式，我们不使用锐化，而是直接进行高质量压缩
		// 移除了对imaging库的依赖
	}

	// 多次尝试不同质量等级压缩
	compressedData, err := compressWithMultipleAttempts(img, outputFormat, fileType, compressionQuality, imgInfo.OriginalSize)
	if err != nil {
		// 压缩失败，返回原始数据
		imgInfo.CompressedSize = imgInfo.OriginalSize
		imgInfo.Ratio = 1.0
		return fileBytes, imgInfo, nil
	}

	// 关键修复：检查压缩效果，确保不会让图片变大
	if int64(len(compressedData)) >= imgInfo.OriginalSize && !isResized {
		// 如果没有缩小尺寸，且压缩后文件更大，返回原始文件
		imgInfo.CompressedSize = imgInfo.OriginalSize
		imgInfo.Ratio = 1.0
		imgInfo.Format = fileType // 保持原始格式
		return fileBytes, imgInfo, nil
	}

	// 更新压缩信息
	imgInfo.CompressedSize = int64(len(compressedData))
	imgInfo.Ratio = float64(imgInfo.CompressedSize) / float64(imgInfo.OriginalSize)
	imgInfo.Format = outputFormat

	return compressedData, imgInfo, nil
}

// 调整图片尺寸
func resizeImage(img image.Image, maxWidth, maxHeight uint) (image.Image, bool) {
	origWidth := uint(img.Bounds().Dx())
	origHeight := uint(img.Bounds().Dy())

	// 如果没有指定限制或图片已经小于限制，不需调整
	if (maxWidth == 0 || origWidth <= maxWidth) &&
		(maxHeight == 0 || origHeight <= maxHeight) {
		return img, false
	}

	// 确定新尺寸
	var newWidth, newHeight uint
	if maxWidth == 0 {
		newWidth = origWidth
	} else {
		newWidth = maxWidth
	}

	if maxHeight == 0 {
		newHeight = origHeight
	} else {
		newHeight = maxHeight
	}

	// 保持宽高比
	widthRatio := float64(newWidth) / float64(origWidth)
	heightRatio := float64(newHeight) / float64(origHeight)

	ratio := math.Min(widthRatio, heightRatio)
	newWidth = uint(float64(origWidth) * ratio)
	newHeight = uint(float64(origHeight) * ratio)

	// 调整尺寸
	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3), true
}

// 以多种质量级别尝试压缩
func compressWithMultipleAttempts(img image.Image, outputFormat, origFormat string, startQuality int, origSize int64) ([]byte, error) {
	// 处理不同输出格式
	if outputFormat == "jpeg" || outputFormat == "jpg" {
		// 对于JPEG格式，尝试一系列不同的质量级别
		qualityLevels := []int{startQuality, startQuality - 10, startQuality - 20, 50, 40, 30}
		if startQuality <= 50 {
			qualityLevels = []int{startQuality, startQuality - 5, 40, 35, 30, 25}
		}

		for _, quality := range qualityLevels {
			if quality < 20 {
				quality = 20 // 避免质量太低
			}

			var buf bytes.Buffer
			err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
			if err != nil {
				continue
			}

			data := buf.Bytes()
			// 确保压缩后的数据比原始数据小，或者已经用了最低质量
			if int64(len(data)) < origSize || quality == qualityLevels[len(qualityLevels)-1] {
				return data, nil
			}
		}
		// 如果所有质量级别都无法使文件变小，返回错误
		return nil, fmt.Errorf("无法有效压缩JPEG")
	} else if outputFormat == "png" {
		// PNG压缩
		var buf bytes.Buffer
		err := png.Encode(&buf, img)
		if err != nil {
			return nil, err
		}

		data := buf.Bytes()

		// 如果PNG太大并且原格式不是PNG，可以尝试高质量JPEG
		// 但前提是图像没有透明度且用户没有明确要求PNG格式
		if int64(len(data)) >= origSize && origFormat != "png" && !hasTransparency(img) {
			// 尝试JPEG压缩
			jpegBuf := bytes.Buffer{}
			err = jpeg.Encode(&jpegBuf, img, &jpeg.Options{Quality: 90})
			if err == nil {
				jpegData := jpegBuf.Bytes()
				if int64(len(jpegData)) < origSize {
					return jpegData, nil
				}
			}
		}

		return data, nil
	} else if outputFormat == "gif" {
		var buf bytes.Buffer
		err := gif.Encode(&buf, img, nil)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	} else if outputFormat == "bmp" {
		var buf bytes.Buffer
		err := bmp.Encode(&buf, img)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	return nil, fmt.Errorf("无法压缩为格式: %s", outputFormat)
}

// GIF压缩专用函数
func compressGif(fileBytes []byte, imgInfo *ImageInfo, options CompressOptions) ([]byte, *ImageInfo, error) {
	// 解码GIF
	gifImg, err := gif.DecodeAll(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("解码GIF失败: %w", err)
	}

	// 更新图片信息
	imgInfo.Width = gifImg.Config.Width
	imgInfo.Height = gifImg.Config.Height
	imgInfo.IsAnimated = len(gifImg.Image) > 1
	imgInfo.FrameCount = len(gifImg.Image)
	imgInfo.Format = "gif"

	// 优化：强制进行尺寸缩减以保证压缩效果
	needsResize := true
	needsOptimize := true

	// 如果用户指定了尺寸，则使用用户的设置
	maxWidth := options.MaxWidth
	maxHeight := options.MaxHeight

	// 如果用户没有指定尺寸，则根据原图自动设定合理的尺寸限制
	// 这确保即使用户没有要求缩小尺寸，我们也能实现一定的压缩
	if maxWidth == 0 && maxHeight == 0 {
		// 对于大尺寸GIF，自动适当缩小
		origWidth := uint(imgInfo.Width)
		origHeight := uint(imgInfo.Height)

		// 根据图片大小决定缩放比例
		var scale float64 = 1.0
		if imgInfo.OriginalSize > 2*1024*1024 {
			// 大于2MB的图片，缩小到70%
			scale = 0.7
		} else if imgInfo.OriginalSize > 1024*1024 {
			// 大于1MB的图片，缩小到80%
			scale = 0.8
		} else if imgInfo.OriginalSize > 512*1024 {
			// 大于512KB的图片，缩小到90%
			scale = 0.9
		} else {
			// 小图片不缩小尺寸
			needsResize = false
		}

		// 应用缩放
		if needsResize {
			maxWidth = uint(float64(origWidth) * scale)
			maxHeight = uint(float64(origHeight) * scale)
		}
	}

	// 优化GIF
	compressedGif := gifImg

	// 调整大小（如果需要）
	if needsResize && (maxWidth > 0 || maxHeight > 0) {
		// 调整大小但保留动画
		resizedGif, err := resizeAnimatedGif(gifImg, maxWidth, maxHeight)
		if err != nil {
			// 调整失败，继续使用原始GIF
			compressedGif = gifImg
			needsResize = false
		} else {
			compressedGif = resizedGif
			// 更新尺寸信息
			imgInfo.Width = resizedGif.Config.Width
			imgInfo.Height = resizedGif.Config.Height
		}
	} else {
		needsResize = false
	}

	// 优化GIF帧数据
	if needsOptimize {
		// 1. 减少颜色数量 - 强制使用更小的调色板
		maxColors := 32 // 最大使用32色，这对大多数GIF效果还可以，但大幅减少文件大小

		for i, img := range compressedGif.Image {
			if len(img.Palette) > maxColors {
				// 创建减少颜色的调色板
				newPalette := make(color.Palette, maxColors)
				copy(newPalette, img.Palette[:maxColors])

				// 创建新的减色图像
				bounds := img.Bounds()
				newImg := image.NewPaletted(bounds, newPalette)

				// 复制并映射到新调色板
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					for x := bounds.Min.X; x < bounds.Max.X; x++ {
						c := img.At(x, y)
						newImg.Set(x, y, newPalette.Convert(c))
					}
				}

				compressedGif.Image[i] = newImg
			}
		}

		// 2. 帧优化 - 对多帧GIF进行抽帧处理
		if imgInfo.FrameCount > 5 {
			// 根据帧数决定保留比例
			keepRatio := 0.5 // 默认保留一半帧
			if imgInfo.FrameCount > 30 {
				keepRatio = 0.3 // 帧数很多时，保留更少的帧
			} else if imgInfo.FrameCount > 15 {
				keepRatio = 0.4 // 适中数量的帧，保留40%
			}

			frameCount := len(compressedGif.Image)
			framesToKeep := int(float64(frameCount) * keepRatio)

			// 确保至少保留2帧
			if framesToKeep < 2 {
				framesToKeep = 2
			}

			// 需要减少帧数
			if framesToKeep < frameCount {
				newImages := make([]*image.Paletted, 0, framesToKeep)
				newDelays := make([]int, 0, framesToKeep)
				newDisposal := make([]byte, 0, framesToKeep)

				// 总是保留第一帧
				newImages = append(newImages, compressedGif.Image[0])
				newDelays = append(newDelays, compressedGif.Delay[0])
				if len(compressedGif.Disposal) > 0 {
					newDisposal = append(newDisposal, compressedGif.Disposal[0])
				}

				// 如果需要保留多于2帧，则在中间选择帧
				if framesToKeep > 2 {
					step := frameCount / (framesToKeep - 1)
					if step < 1 {
						step = 1
					}

					for i := step; i < frameCount-1; i += step {
						if len(newImages) < framesToKeep-1 {
							newImages = append(newImages, compressedGif.Image[i])
							newDelays = append(newDelays, compressedGif.Delay[i])
							if i < len(compressedGif.Disposal) {
								newDisposal = append(newDisposal, compressedGif.Disposal[i])
							}
						}
					}
				}

				// 总是保留最后一帧
				lastIdx := frameCount - 1
				if lastIdx > 0 {
					newImages = append(newImages, compressedGif.Image[lastIdx])
					newDelays = append(newDelays, compressedGif.Delay[lastIdx])
					if lastIdx < len(compressedGif.Disposal) {
						newDisposal = append(newDisposal, compressedGif.Disposal[lastIdx])
					}
				}

				// 调整帧延迟以保持总动画时长
				totalOrigDelay := 0
				for _, d := range compressedGif.Delay {
					totalOrigDelay += d
				}

				totalNewDelay := 0
				for _, d := range newDelays {
					totalNewDelay += d
				}

				// 按比例调整延迟
				if totalNewDelay > 0 {
					delayFactor := float64(totalOrigDelay) / float64(totalNewDelay)
					for i := range newDelays {
						newDelays[i] = int(float64(newDelays[i]) * delayFactor)

						// 限制延迟范围
						if newDelays[i] < 2 {
							newDelays[i] = 2
						} else if newDelays[i] > 100 {
							newDelays[i] = 100
						}
					}
				}

				// 更新GIF信息
				compressedGif.Image = newImages
				compressedGif.Delay = newDelays
				compressedGif.Disposal = newDisposal
				imgInfo.FrameCount = len(newImages)
			}
		}
	}

	// 编码GIF
	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, compressedGif)
	if err != nil {
		// 编码失败，返回原始数据
		imgInfo.CompressedSize = imgInfo.OriginalSize
		imgInfo.Ratio = 1.0
		return fileBytes, imgInfo, nil
	}

	compressedData := buf.Bytes()

	// 如果压缩后实际更大，强制使用较小的版本
	// 如果尺寸没有减小，但颜色减少了，仍然认为是合理的压缩
	forceUseCompressed := needsResize || imgInfo.FrameCount < len(gifImg.Image)

	if int64(len(compressedData)) >= imgInfo.OriginalSize && !forceUseCompressed {
		imgInfo.CompressedSize = imgInfo.OriginalSize
		imgInfo.Ratio = 1.0
		return fileBytes, imgInfo, nil
	}

	// 更新压缩信息
	imgInfo.CompressedSize = int64(len(compressedData))
	imgInfo.Ratio = float64(imgInfo.CompressedSize) / float64(imgInfo.OriginalSize)

	return compressedData, imgInfo, nil
}

// 根据文件内容检测文件类型
func detectFileType(data []byte) string {
	if len(data) < 8 {
		return ""
	}

	// 检查PNG签名
	if bytes.Equal(data[0:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return "png"
	}

	// 检查GIF签名
	if bytes.Equal(data[0:6], []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}) ||
		bytes.Equal(data[0:6], []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}) {
		return "gif"
	}

	// 检查JPEG签名
	if bytes.Equal(data[0:2], []byte{0xFF, 0xD8}) {
		return "jpeg"
	}

	// 检查BMP签名
	if bytes.Equal(data[0:2], []byte{0x42, 0x4D}) {
		return "bmp"
	}

	// 如果无法识别，尝试用标准库检测
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err == nil {
		return format
	}

	return ""
}

// 调整GIF动画大小，保留所有帧和动画属性
func resizeAnimatedGif(gifImg *gif.GIF, maxWidth, maxHeight uint) (*gif.GIF, error) {
	if gifImg == nil || len(gifImg.Image) == 0 {
		return nil, fmt.Errorf("无效的GIF图像")
	}

	// 获取原始尺寸
	originalWidth := gifImg.Config.Width
	originalHeight := gifImg.Config.Height

	// 如果不需要调整大小，直接返回原始GIF
	if (maxWidth == 0 || uint(originalWidth) <= maxWidth) &&
		(maxHeight == 0 || uint(originalHeight) <= maxHeight) {
		return gifImg, nil
	}

	// 计算新尺寸，保持宽高比
	var newWidth, newHeight uint
	if maxWidth == 0 {
		newWidth = uint(originalWidth)
	} else {
		newWidth = maxWidth
	}

	if maxHeight == 0 {
		newHeight = uint(originalHeight)
	} else {
		newHeight = maxHeight
	}

	// 调整宽高比
	if float64(newWidth)/float64(originalWidth) < float64(newHeight)/float64(originalHeight) {
		newHeight = uint(float64(originalHeight) * float64(newWidth) / float64(originalWidth))
	} else {
		newWidth = uint(float64(originalWidth) * float64(newHeight) / float64(originalHeight))
	}

	// 创建新GIF
	newGif := &gif.GIF{
		Image:           make([]*image.Paletted, len(gifImg.Image)),
		Delay:           make([]int, len(gifImg.Delay)),
		LoopCount:       gifImg.LoopCount,
		Disposal:        make([]byte, len(gifImg.Disposal)),
		Config:          image.Config{ColorModel: gifImg.Config.ColorModel, Width: int(newWidth), Height: int(newHeight)},
		BackgroundIndex: gifImg.BackgroundIndex,
	}

	// 复制延迟和处理方式
	copy(newGif.Delay, gifImg.Delay)
	copy(newGif.Disposal, gifImg.Disposal)

	// 处理每一帧
	for i, srcImg := range gifImg.Image {
		// 调整帧大小
		resized := resize.Resize(newWidth, newHeight, srcImg, resize.Lanczos3)

		// 创建新的调色板图像
		paletted := image.NewPaletted(image.Rect(0, 0, int(newWidth), int(newHeight)), srcImg.Palette)

		// 将调整大小后的图像复制到调色板图像
		draw.Draw(paletted, paletted.Bounds(), resized, resized.Bounds().Min, draw.Over)

		// 添加到新GIF
		newGif.Image[i] = paletted
	}

	return newGif, nil
}

// 检查图像是否包含透明像素
func hasTransparency(img image.Image) bool {
	if rgba, ok := img.(*image.RGBA); ok {
		bounds := rgba.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				_, _, _, a := rgba.At(x, y).RGBA()
				if a < 65535 {
					return true
				}
			}
		}
	} else if nrgba, ok := img.(*image.NRGBA); ok {
		bounds := nrgba.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				_, _, _, a := nrgba.At(x, y).RGBA()
				if a < 65535 {
					return true
				}
			}
		}
	}
	return false
}

// BatchCompressImages 批量压缩图片
func BatchCompressImages(files []*multipart.FileHeader, options CompressOptions) (map[string][]byte, map[string]*ImageInfo, []error) {
	results := make(map[string][]byte)
	infos := make(map[string]*ImageInfo)
	errors := make([]error, 0)

	for _, file := range files {
		data, info, err := CompressImage(file, options)
		if err != nil {
			errors = append(errors, fmt.Errorf("压缩 %s 失败: %w", file.Filename, err))
			continue
		}

		// 确定输出文件名（原文件名或添加_compressed后缀）
		filename := file.Filename
		if !options.PreserveName {
			ext := filepath.Ext(filename)
			baseFilename := strings.TrimSuffix(filename, ext)
			newExt := "." + info.Format
			filename = baseFilename + "_compressed" + newExt
		}

		results[filename] = data
		infos[filename] = info
	}

	return results, infos, errors
}

// 评估图像质量（基于文件大小与尺寸的比率以及格式）
func estimateQualityScore(imageData []byte, format string) int {
	// 此处仅为示例实现，实际应用中应使用更复杂的图像质量评估算法
	// 例如：PSNR, SSIM, MS-SSIM 等

	dataSize := len(imageData)
	if dataSize < 1024 {
		return 30 // 文件太小，可能质量较低
	} else if dataSize > 5*1024*1024 {
		return 95 // 大文件，可能质量较高
	}

	// 根据格式给予初始分数
	baseScore := 70
	switch format {
	case "jpeg", "jpg":
		baseScore = 75
	case "png":
		baseScore = 85
	case "webp":
		baseScore = 80
	case "gif":
		baseScore = 60
	}

	// 实际应用中，这里应该结合图像数据进行分析
	// 例如分析噪声、细节保留、色彩度量等因素

	return baseScore
}
