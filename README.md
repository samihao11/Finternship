# Finternship

Finternship is your AI-powered coding companion designed for students who want to make the most of their time between internships or while job searching. It helps you turn your personal projects into professional-grade portfolios by providing intelligent suggestions and improvements. Think of it as having a senior developer reviewing your code and guiding you on what to improve next!

## Why Finternship?

- 🎯 **Stay Productive**: Keep improving your projects even without an internship
- 🚀 **Level Up Your Code**: Get professional-grade suggestions for your projects
- 💼 **Portfolio Building**: Transform your projects into impressive portfolio pieces
- 📈 **Continuous Learning**: Receive actionable tasks that teach you best practices
- 🤖 **AI Guidance**: Get suggestions from an AI trained on millions of repositories

## Prerequisites

- Go 1.16 or higher
- Node.js 14.0 or higher
- npm or yarn
- GitHub Personal Access Token
- OpenAI API Key

## Quick Start

### Backend Setup

1. Clone the repository:

```bash
git clone https://github.com/yourusername/Finternship.git
cd Finternship
```

2. Create a `.env` file in the root directory:

```bash
# Create .env file
touch .env

# Add your API keys to the .env file
echo "GITHUB_ACCESS_TOKEN=your_github_token_here
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_API_BASE_URL=https://api.openai.com/v1
GITHUB_API_BASE_URL=https://api.github.com" > .env
```

3. Install Go dependencies:

```bash
go mod tidy
```

4. Run the Go backend:

```bash
go run main.go
```

The backend server will start on `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:

```bash
cd frontend
```

2. Install dependencies:

```bash
npm install
```

3. Start the development server:

```bash
npm start
```

The frontend will be available at `http://localhost:3000`

## How to Use

1. **Upload Your Project**

   - Enter your GitHub username and repository name
   - Get instant access to your codebase

2. **Generate Tasks**

   - Click "Generate To-Dos" to get your daily improvement tasks
   - Each task is specifically tailored to your project

3. **Track Progress**
   - View your code with syntax highlighting
   - Mark tasks as resolved as you complete them
   - Generate new tasks when you're ready for more challenges

## Technologies Used

### Backend

- Go
- Gin Web Framework
- GitHub API
- OpenAI API
- godotenv (for environment variables)

### Frontend

- React
- Axios
- react-syntax-highlighter
- CSS with glass-morphism design

## Getting Started with API Keys

### GitHub Token

1. Go to GitHub.com → Settings → Developer Settings → Personal Access Tokens → Tokens (classic)
2. Generate a new token with `repo` scope
3. Add as `GITHUB_ACCESS_TOKEN` in your `.env` file

### OpenAI API Key

1. Visit https://platform.openai.com/api-keys
2. Create a new API key
3. Add as `OPENAI_API_KEY` in your `.env` file

## Environment Variables

The project uses the following environment variables in the `.env` file:

```bash
GITHUB_ACCESS_TOKEN=your_github_token_here
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_API_BASE_URL=https://api.openai.com/v1
GITHUB_API_BASE_URL=https://api.github.com
```

## Security Notes

- Never commit your `.env` file to version control (it's already in `.gitignore`)
- Keep your API keys secure and never share them
- Keep your dependencies updated
- Use environment variables for all sensitive data

## Troubleshooting

Common issues and solutions:

1. **"Error loading .env file"**: Make sure you've created the `.env` file in the root directory with the correct variables
2. **"OpenAI API key not configured"**: Check that your OpenAI API key is correctly set in the `.env` file
3. **"Failed to fetch repository"**: Verify your GitHub access token is valid and has the correct permissions
4. **"node_modules not found"**: Run `npm install` in the frontend directory

## Support

Need help? Create an issue in the GitHub repository. We're here to help you make the most of your learning journey!

## Contributing

Have ideas to make Finternship better for students? We'd love your contributions:

1. Fork the repository
2. Create a feature branch
3. Make your improvements
4. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
