package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// FileInfo represents a file in the repository
type FileInfo struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Language string `json:"language"`
}

// RepoInfo represents a GitHub repository submission
type RepoInfo struct {
	GithubUsername string     `json:"githubUsername"`
	RepoName       string     `json:"repoName"`
	Files          []FileInfo `json:"files"`
}

// GitHubTreeResponse represents the GitHub API response for the tree endpoint
type GitHubTreeResponse struct {
	Tree []struct {
		Path string `json:"path"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"tree"`
}

// GitHubContentResponse represents the GitHub API response for the contents endpoint
type GitHubContentResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response structure from OpenAI API
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// TodoResponse represents the response structure for todos
type TodoResponse struct {
	Todos []string `json:"todos"`
}

func main() {
	// Load the .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Error loading .env file: %v\n", err)
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/repo/:username/:reponame", getRepoFiles)
	router.GET("/todos/:username/:reponame", generateTodos)

	router.Run("0.0.0.0:8080")
}

// makeGitHubRequest makes an authenticated request to the GitHub API
func makeGitHubRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add GitHub token if available
	if token := os.Getenv("GITHUB_ACCESS_TOKEN"); token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	return http.DefaultClient.Do(req)
}

// getRepoFiles fetches repository files from GitHub API
func getRepoFiles(c *gin.Context) {
	username := c.Param("username")
	reponame := c.Param("reponame")

	if username == "" || reponame == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and repository name are required"})
		return
	}

	// 1. Get the repository tree
	treeURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/main?recursive=1", username, reponame)
	treeResp, err := makeGitHubRequest(treeURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repository tree"})
		return
	}
	defer treeResp.Body.Close()

	if treeResp.StatusCode == http.StatusNotFound {
		// Try 'master' branch if 'main' doesn't exist
		treeURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/master?recursive=1", username, reponame)
		treeResp, err = makeGitHubRequest(treeURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch repository tree"})
			return
		}
		defer treeResp.Body.Close()
	}

	if treeResp.StatusCode != http.StatusOK {
		c.JSON(treeResp.StatusCode, gin.H{"error": "Repository or branch not found"})
		return
	}

	var treeData GitHubTreeResponse
	if err := json.NewDecoder(treeResp.Body).Decode(&treeData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse repository tree"})
		return
	}

	// 2. Get contents for each file
	var files []FileInfo
	for _, item := range treeData.Tree {
		// Skip directories and large files
		if item.Type != "blob" {
			continue
		}

		// Get file contents
		contentURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", username, reponame, item.Path)
		contentResp, err := makeGitHubRequest(contentURL)
		if err != nil {
			continue
		}
		defer contentResp.Body.Close()

		if contentResp.StatusCode != http.StatusOK {
			continue
		}

		var contentData GitHubContentResponse
		if err := json.NewDecoder(contentResp.Body).Decode(&contentData); err != nil {
			continue
		}

		// Decode base64 content
		content, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(contentData.Content, "\n", ""))
		if err != nil {
			continue
		}

		// Determine language from file extension
		language := getLanguageFromExtension(item.Path)

		files = append(files, FileInfo{
			Path:     item.Path,
			Content:  string(content),
			Language: language,
		})
	}

	repoInfo := RepoInfo{
		GithubUsername: username,
		RepoName:      reponame,
		Files:         files,
	}

	c.JSON(http.StatusOK, repoInfo)
}

// getLanguageFromExtension returns the programming language based on file extension
func getLanguageFromExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return "text"
	}
	ext := parts[len(parts)-1]

	switch ext {
	case "js", "jsx":
		return "javascript"
	case "ts", "tsx":
		return "typescript"
	case "py":
		return "python"
	case "java":
		return "java"
	case "go":
		return "go"
	case "html":
		return "html"
	case "css":
		return "css"
	case "json":
		return "json"
	case "md":
		return "markdown"
	case "php":
		return "php"
	case "rb":
		return "ruby"
	case "rs":
		return "rust"
	case "sh":
		return "shell"
	case "sql":
		return "sql"
	case "xml":
		return "xml"
	case "yaml", "yml":
		return "yaml"
	default:
		return "text"
	}
}

// generateTodos creates a todo list using OpenAI based on repository files
func generateTodos(c *gin.Context) {
	username := c.Param("username")
	reponame := c.Param("reponame")

	fmt.Printf("Generating todos for %s/%s\n", username, reponame)

	// Get OpenAI API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OpenAI API key not configured"})
		return
	}

	// First, get all repository files
	files, err := fetchRepoFiles(username, reponame)
	if err != nil {
		fmt.Printf("Error fetching repo files: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch repository files: %v", err)})
		return
	}

	fmt.Printf("Found %d files in repository\n", len(files))

	// Prepare file content summary for OpenAI
	var filesSummary string
	for _, file := range files {
		// Limit content length to avoid token limits
		content := file.Content
		if len(content) > 1000 {
			content = content[:1000] + "... (truncated)"
		}
		filesSummary += fmt.Sprintf("File: %s\nContent:\n%s\n\n", file.Path, content)
	}

	// Create OpenAI API request
	prompt := fmt.Sprintf("Based on the repository %s/%s, give me a numbered list 1-5 of specific, actionable tasks I can do today to improve this project. Format your response as a clean numbered list with exactly 5 items, starting each line with a number followed by a period. Focus on concrete improvements, bug fixes, and feature additions. Be specific and actionable.", username, reponame)

	openAIReq := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful programming assistant that provides specific, actionable todo items for improving code repositories. Always provide exactly 5 numbered items, no more and no less.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
			{
				Role:    "assistant",
				Content: "I'll analyze the repository and provide exactly 5 specific, actionable tasks to improve the project.",
			},
			{
				Role:    "user",
				Content: filesSummary,
			},
		},
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(openAIReq)
	if err != nil {
		fmt.Printf("Error marshaling OpenAI request: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OpenAI request"})
		return
	}

	fmt.Println("Making request to OpenAI API...")

	// Make request to OpenAI API
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating OpenAI request: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OpenAI request"})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making OpenAI request: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get response from OpenAI"})
		return
	}
	defer resp.Body.Close()

	// Read and log the response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading OpenAI response body: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read OpenAI response"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("OpenAI API error (status %d): %s\n", resp.StatusCode, string(body))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("OpenAI API error: %s", string(body))})
		return
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		fmt.Printf("Error parsing OpenAI response: %v\nResponse body: %s\n", err, string(body))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse OpenAI response"})
		return
	}

	if len(openAIResp.Choices) == 0 {
		fmt.Println("No choices in OpenAI response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No response from OpenAI"})
		return
	}

	// Parse the numbered list into separate todos
	content := openAIResp.Choices[0].Message.Content
	lines := strings.Split(content, "\n")
	var todos []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && (strings.HasPrefix(line, "- ") || (len(line) > 2 && line[1] == '.' && line[0] >= '0' && line[0] <= '9')) {
			// Remove the bullet point or number
			todo := strings.TrimPrefix(strings.TrimPrefix(line, "- "), "0123456789. ")
			todos = append(todos, todo)
		}
	}

	fmt.Printf("Generated %d todos\n", len(todos))
	c.JSON(http.StatusOK, TodoResponse{Todos: todos})
}

// fetchRepoFiles gets all files from a repository
func fetchRepoFiles(username, reponame string) ([]FileInfo, error) {
	// Get the repository tree
	treeURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/main?recursive=1", username, reponame)
	treeResp, err := makeGitHubRequest(treeURL)
	if err != nil {
		return nil, err
	}
	defer treeResp.Body.Close()

	if treeResp.StatusCode == http.StatusNotFound {
		// Try 'master' branch if 'main' doesn't exist
		treeURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/master?recursive=1", username, reponame)
		treeResp, err = makeGitHubRequest(treeURL)
		if err != nil {
			return nil, err
		}
		defer treeResp.Body.Close()
	}

	if treeResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get repository tree")
	}

	var treeData GitHubTreeResponse
	if err := json.NewDecoder(treeResp.Body).Decode(&treeData); err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, item := range treeData.Tree {
		if item.Type != "blob" {
			continue
		}

		contentURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", username, reponame, item.Path)
		contentResp, err := makeGitHubRequest(contentURL)
		if err != nil {
			continue
		}
		defer contentResp.Body.Close()

		if contentResp.StatusCode != http.StatusOK {
			continue
		}

		var contentData GitHubContentResponse
		if err := json.NewDecoder(contentResp.Body).Decode(&contentData); err != nil {
			continue
		}

		content, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(contentData.Content, "\n", ""))
		if err != nil {
			continue
		}

		language := getLanguageFromExtension(item.Path)

		files = append(files, FileInfo{
			Path:     item.Path,
			Content:  string(content),
			Language: language,
		})
	}

	return files, nil
}
