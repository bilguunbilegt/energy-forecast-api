# energy-forecast-api

ssh -i ~/.ssh/jan_09_2025_virginia.pem ubuntu@44.201.248.65

sudo apt update -y
sudo apt install -y docker.io

sudo systemctl enable docker
sudo systemctl start docker

docker --version

sudo apt install -y git
git clone https://github.com/bilguunbilegt/energy-forecast-api.git
cd energy-forecast-api

docker build -t energy-forecast-api .

docker run -d -p 5000:5000 energy-forecast-api

curl -X POST "http://ec2-44-201-248-65.compute-1.amazonaws.com:5000/predict" -H "Content-Type: application/json" -d '{
  "population": 13000000,
  "temperature": 75
}'
