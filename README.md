# aws-golang-rest-api-with-dynamodb

Build & Deploy

```
make deploy
```

## Test Apis

you can use the postman collection or bellow curl

## Create

```
curl --request POST \
  --url https://gdcnt3n1ic.execute-api.us-east-1.amazonaws.com/dev/devices \
  --header 'Content-Type: application/json' \
  --data '{
    "id":"2",
    "deviceModel": "model Name",
    "name": "Sensor",
    "note": "Testing a sensor.",
    "serial": "A020000102"
}'
```

## Read

```
curl --request GET \
  --url https://gdcnt3n1ic.execute-api.us-east-1.amazonaws.com/dev/devices/{id}
```
