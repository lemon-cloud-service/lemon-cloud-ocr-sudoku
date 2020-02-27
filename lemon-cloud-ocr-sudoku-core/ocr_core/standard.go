package ocr_core

import (
	"fmt"
	"gocv.io/x/gocv"
	"image/color"
)

func Process() {
	//filename := "/Users/liuri/Desktop/sudoku01.jpeg"
	filename := "test-data/hand11.jpg"
	window := gocv.NewWindow("Hello")
	img := gocv.IMRead(filename, gocv.IMReadColor)
	if img.Empty() {
		fmt.Printf("Error reading image from: %v\n", filename)
		return
	}
	cImg := img.Clone()

	// 灰度图
	var grayImg gocv.Mat = gocv.NewMat()
	gocv.CvtColor(img, &grayImg, gocv.ColorRGBToGray)

	// 二值图
	var tvImg gocv.Mat = gocv.NewMat()
	gocv.Threshold(grayImg, &tvImg, 120, 255, gocv.ThresholdBinary)

	// 形态学
	//xtxMat := gocv.GetStructuringElement(gocv.MorphRect, image.Point{20, 20})
	//gocv.Erode(tvImg, &tvImg, xtxMat)
	//gocv.Dilate(tvImg, &tvImg, xtxMat)

	// 寻找轮廓
	countorsPoints := gocv.FindContours(tvImg, gocv.RetrievalList, gocv.ChainApproxSimple)
	countorIndex := 0
	var maxArea float64
	var maxIndex int

	for countorIndex = 0; countorIndex < len(countorsPoints); countorIndex++ {
		points := countorsPoints[countorIndex]
		// 计算面积
		cArea := gocv.ContourArea(points)
		if cArea > maxArea {
			maxArea = cArea
			maxIndex = countorIndex
		}
		if cArea < 10 {
			continue
		}
		// 绘制轮廓，用颜色圈上
		gocv.DrawContours(&cImg, countorsPoints, countorIndex, color.RGBA{
			R: 0,
			G: 255,
			B: 0,
			A: 255,
		}, 2)
		//gocv.ApproxPolyDP(points, 10, true)
	}
	// 绘制最大的轮廓(整个宫的轮廓)
	gocv.DrawContours(&cImg, countorsPoints, maxIndex, color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}, 2)
	// 检测最大轮廓的突起部分,实际上就是最外面宫线的四个顶点
	hull := gocv.NewMat()
	gocv.ConvexHull(countorsPoints[maxIndex], &hull, false, true)
	for i := 0; i < hull.Size()[0]; i++ {
		gocv.Circle(&cImg, countorsPoints[maxIndex][i], 10, color.RGBA{
			R: 0,
			G: 0,
			B: 255,
			A: 255,
		}, 2)
	}

	for {
		window.IMShow(tvImg)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
