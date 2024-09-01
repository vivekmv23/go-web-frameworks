# go-web-frameworks
A simple web service written in popular go web frameworks.

# Reference
1. https://www.jetbrains.com/guide/go/tutorials/rest_api_series

## Application

- Exposes a CRUD Restful API
- Uses MongoDB for persistence
- Should have one auth check middleware

## Model: Item

```json
{
  "id": "a79c2798-dc26-40ff-a2ab-3cbca3af5413",
  "name": "Item Name",
  "value": 1000,
  "description": "Item Description",
  "isActive": true,
  "createdOn": "2024-09-01T10:16:35.602Z",
  "updatedOn": "2024-09-01T10:16:35.602Z"
}
```


