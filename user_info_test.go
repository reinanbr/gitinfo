package gitinfo
import (
	"testing"
	"github.com/joho/godotenv"
)

func TestGetUserInfo(t *testing.T) {
	tokens, err := godotenv.Read(".env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v\n", err)
	}

	githubToken, ok := tokens["GITHUB_TOKEN"]
	if !ok {
		t.Fatal("GITHUB_TOKEN not found in .env file")
	}

	userInfo, err := GetUserInfo("reinanbr", githubToken)
	if err != nil {
		t.Fatalf("Error fetching user info: %v\n", err)
	}

	if userInfo.Login == "" {
		t.Fatal("Login should not be empty")
	}

	if userInfo.Name == "" {
		t.Fatal("Name should not be empty")
	}
	if userInfo.Bio == "" {
		t.Fatal("Bio should not be empty")
	}
	if userInfo.AvatarUrl == "" {
		t.Fatal("AvatarUrl should not be empty")
	}
	if userInfo.CreatedAt == "" {
		t.Fatal("CreatedAt should not be empty")
	}
	if userInfo.ID == "" {
		t.Fatal("ID should not be empty")
	}
	if userInfo.URL == "" {
		t.Fatal("URL should not be empty")
	}
	if userInfo.Repositories.TotalCount < 0 {
		t.Fatal("Repositories total count should be non-negative")
	}
	if userInfo.Followers.TotalCount < 0 {
		t.Fatal("Followers total count should be non-negative")
	}
	if userInfo.Following.TotalCount < 0 {
		t.Fatal("Following total count should be non-negative")
	}


	t.Logf("Login: %s\n", userInfo.Login)
	t.Logf("Name: %s\n", userInfo.Name)
	t.Logf("Created At: %s\n", userInfo.CreatedAt)
	t.Logf("Avatar URL: %s\n", userInfo.AvatarUrl)
	t.Logf("Bio: %s\n", userInfo.Bio)
	t.Logf("ID: %s\n", userInfo.ID)
	t.Logf("URL: %s\n", userInfo.URL)
	t.Logf("Total Repositories: %d\n", userInfo.Repositories.TotalCount)
	t.Logf("Total Followers: %d\n", userInfo.Followers.TotalCount)
	t.Logf("Total Following: %d\n", userInfo.Following.TotalCount)
	
}