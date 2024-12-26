# Go multipart/form-data image upload handling

This toy API demonstrates how to handle multipart/form-data requests with images in a Go HTTP server.

Currently, the uploaded images are stored on the local file system.

### Endpoints
- `POST localhost:4321/api/users`
    - All photos in this toy app are associated with a user
    - in the request body, provide a JSON object with a `username` key with a `string` value
- `POST localhost:4321/api/users/<userId>/profile_pic`
    - Create a new photo for the given `userId`
    - provide the photo file under the `profilePic` key in a `multipart/form-data` encoded request body
- `GET localhost:4321/api/users/<userId>/profile_pic/<imageId>`
    - Retrieve an uploaded image
- `DELETE localhost:4321/api/users/<userId>/profile_pic/<imageId>`
    - Delete an uploaded image
