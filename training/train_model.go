package training

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sajari/regression"
)

type EnergyData struct {
	Population  float64
	Temperature float64
	EnergyKWh   float64
}

const (
	bucketName = "sage-bilguun"
	objectKey  = "train/energy_train.csv"
)

func TrainModel() {
	data, err := downloadFromS3(bucketName, objectKey)
	if err != nil {
		log.Fatalf("Failed to download file from S3: %v", err)
	}

	energyData, err := parseCSVData(data)
	if err != nil {
		log.Fatalf("Failed to parse CSV: %v", err)
	}

	model := trainRegressionModel(energyData)
	saveModel(model)

	fmt.Println("Model trained and saved successfully.")
}

func downloadFromS3(bucket, key string) ([]byte, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	}))

	svc := s3.New(sess)
	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return io.ReadAll(obj.Body)
}

func parseCSVData(data []byte) ([]EnergyData, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var energyData []EnergyData
	for i, record := range records {
		if i == 0 {
			continue
		}
		energyKWh, _ := strconv.ParseFloat(record[2], 64)
		population, _ := strconv.ParseFloat(record[3], 64)
		temperature, _ := strconv.ParseFloat(record[4], 64)
		energyData = append(energyData, EnergyData{population, temperature, energyKWh})
	}
	return energyData, nil
}

func trainRegressionModel(data []EnergyData) *regression.Regression {
	var r regression.Regression
	r.SetObserved("EnergyKWh")
	r.SetVar(0, "Population")
	r.SetVar(1, "Temperature")

	for _, d := range data {
		r.Train(regression.DataPoint(d.EnergyKWh, []float64{d.Population, d.Temperature}))
	}

	r.Run()
	return &r
}

func saveModel(model *regression.Regression) {
	coefficients := model.GetCoeffs()
	modelData, err := json.MarshalIndent(coefficients, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode model data: %v", err)
	}

	if err := os.WriteFile("model.json", modelData, 0644); err != nil {
		log.Fatalf("Failed to save model: %v", err)
	}
}
