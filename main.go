package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func readDataFromFile(filename string) ([]float64, []float64) {
	file, _ := os.Open(filename)
	defer file.Close()
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	var years, areas []float64
	for i, record := range records {
		if i == 0 {
			continue
		}
		year, _ := strconv.ParseFloat(record[0], 64)
		area, _ := strconv.ParseFloat(record[1], 64)
		years = append(years, year)
		areas = append(areas, area)
	}
	return years, areas
}

func linearRegression(x, y []float64) (float64, float64) {
	var sumX, sumY, sumXY, sumXX float64
	for i := range x {
		sumX += x[i]
		sumY += y[i]
		//sumXY = sumXY + x[i]*y[i]
		sumXY += x[i] * y[i]
		sumXX += x[i] * x[i]
	}
	n := float64(len(x))
	//formula to calc slope: m = n.sumXY -sumX.sumY / n.SumXX-sumX^2
	slope := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)

	//formula to calc intercept: b = sumY -m.sumX / n
	intercept := (sumY - slope*sumX) / n
	return slope, intercept
}

func predict(slope, intercept, year float64) float64 {
	// to predict y
	//linear equation y = mx + b
	return slope*year + intercept
}

func startWebServer() {
	r := gin.Default()

	r.Static("/static", "./static")

	r.GET("/predict", func(c *gin.Context) {
		yearStr := c.Query("year")
		year, err := strconv.ParseFloat(yearStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
			return
		}

		years, areas := readDataFromFile("city_data.csv")
		slope, intercept := linearRegression(years, areas)
		predictedArea := predict(slope, intercept, year)

		c.JSON(http.StatusOK, gin.H{"year": year, "predicted_area": predictedArea})
	})

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	fmt.Println("Web server started at http://localhost:8080")
	r.Run(":8080")
}

func main() {
	/*
		years, areas := readDataFromFile("city_data.csv")
		slope, intercept := linearRegression(years, areas)

		var futureYear float64
		fmt.Print("Enter a future year for prediction: ")
		fmt.Scan(&futureYear)

		predictedArea := predict(slope, intercept, futureYear)
		fmt.Printf("Predicted urban area for %.0f: %.2f sq km\n", futureYear, predictedArea)
	*/

	startWebServer()
}
