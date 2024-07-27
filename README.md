# go-web-frameworks
A simple web service written in popular go web frameworks

## Application

- Exposes a CRUD Restful API
- Exposes additional health checks endpoints
- Uses MongoDB for persistence
- Should have one auth check middleware

## Model: Item

```json
{
  "id": "a79c2798-dc26-40ff-a2ab-3cbca3af5413"
  "name": "ItemName"
  "value": 1000
  "description": "Item description"
  "active": true
  "created-on": "2024-07-26T00.00.000Z"
  "updated-on": "2024-07-27T00.00.000Z"
}
```


