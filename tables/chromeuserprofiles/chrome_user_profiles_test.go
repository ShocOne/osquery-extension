package chromeuserprofiles

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBtoi(t *testing.T) {
	assert.Equal(t, 1, btoi(true))
	assert.Equal(t, 0, btoi(false))
}

func TestGoogleChromeProfilesColumns(t *testing.T) {
	columns := GoogleChromeProfilesColumns()
	assert.Len(t, columns, 4)

	expectedColumnNames := []string{"username", "email", "name", "ephemeral"}
	for i, column := range columns {
		assert.Equal(t, expectedColumnNames[i], column.Name)
	}
}

func TestFindFileInUserDirs(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a user directory inside the temporary directory
	userDir := filepath.Join(tempDir, "testuser")
	os.Mkdir(userDir, os.ModePerm)

	// Create a test file inside the user directory
	testFile := filepath.Join(userDir, "testfile.txt")
	os.WriteFile(testFile, []byte("test data"), os.ModePerm)

	// Set the home directory location for the current platform
	homeDirLocations[runtime.GOOS] = []string{tempDir}

	// Test with a username
	foundFiles, err := findFileInUserDirs("testfile.txt", WithUsername("testuser"))
	assert.NoError(t, err)
	assert.Len(t, foundFiles, 1)
	assert.Equal(t, "testuser", foundFiles[0].user)
	assert.Equal(t, testFile, foundFiles[0].path)

	// Test without a username
	foundFiles, err = findFileInUserDirs("testfile.txt")
	assert.NoError(t, err)
	assert.Len(t, foundFiles, 1)
	assert.Equal(t, "testuser", foundFiles[0].user)
	assert.Equal(t, testFile, foundFiles[0].path)
}

func TestGenerateForPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a test Chrome local state file
	localStateFile := filepath.Join(tempDir, "Local State")
	localStateData := `{
		"profile": {
			"info_cache": {
				"profile1": {
					"name": "Profile 1",
					"is_ephemeral": false,
					"user_name": "profile1@example.com"
				},
				"profile2": {
					"name": "Profile 2",
					"is_ephemeral": true,
					"user_name": "profile2@example.com"
				}
			}
		}
	}`

	os.WriteFile(localStateFile, []byte(localStateData), os.ModePerm)

	// Test generateForPath
	fileInfo := userFileInfo{
		user: "testuser",
		path: localStateFile,
	}

	results, err := generateForPath(context.Background(), fileInfo)
	assert.NoError(t, err)
	assert.Len(t, results, 2)

	expectedProfiles := []map[string]string{
		{
			"username":  "testuser",
			"email":     "profile1@example.com",
			"name":      "Profile 1",
			"ephemeral": "0",
		},
		{
			"username":  "testuser",
			"email":     "profile2@example.com",
			"name":      "Profile 2",
			"ephemeral": "1",
		},
	}

	assert.ElementsMatch(t, expectedProfiles, results)
}
