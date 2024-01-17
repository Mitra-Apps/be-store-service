package storage_test

// func TestUploadImage(t *testing.T) {
// 	os.Setenv("MINIO_ENDPOINT", "172.19.14.120:5000")
// 	os.Setenv("MINIO_ACCESS_KEY", "ae1LiTEzMxnFvg3UHnW8")
// 	os.Setenv("MINIO_SECRET_KEY", "mLBLwRXc6DVUovyeIZ53mL7gMy4LpWmX0OWl4vKv")
// 	os.Setenv("MINIO_BUCKET_NAME", "test-bucket")

// 	store := storage.New()

// 	fakeImageData := generateFakeBase64Image()

// 	url, err := store.UploadImage(context.Background(), fakeImageData)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, url)
// }

// func generateFakeBase64Image() string {
// 	return "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAWElEQVR42mJ8UgR2GgAxEZgBNqgAVIMKQRYWUoikgphzE4AAI0G8IZtJYAAAAAElFTkSuQmCC=="
// }
