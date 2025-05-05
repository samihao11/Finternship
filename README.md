# Finternship

Finternship is your AI-powered coding companion designed for students who want to make the most of their time between internships or while job searching. It helps you turn your personal projects into professional-grade portfolios by providing intelligent suggestions and improvements. Think of it as having a senior developer reviewing your code and guiding you on what to improve next!

## Why Finternship?

- 🎯 **Stay Productive**: Keep improving your projects even without an internship
- 🚀 **Level Up Your Code**: Get professional-grade suggestions for your projects
- 💼 **Portfolio Building**: Transform your projects into impressive portfolio pieces
- 📈 **Continuous Learning**: Receive actionable tasks that teach you best practices
- 🤖 **AI Guidance**: Get suggestions from an AI trained on millions of repositories

## Features

- 🔍 **Smart Project Analysis**: Upload any GitHub repository and get instant insights
- ✅ **Daily Todo Generator**: Get 5 actionable tasks each day to improve your code
- 📚 **Code Viewer**: Browse your project files with beautiful syntax highlighting
- 💡 **Best Practices**: Learn industry standards and coding patterns
- 📝 **Progress Tracking**: Mark completed tasks and see your project evolve

## How It Helps Students

1. **Project Enhancement**: Turn basic projects into professional-quality work
2. **Learning Opportunities**: Each suggestion teaches you something new
3. **Interview Preparation**: Build stronger projects to showcase in interviews
4. **Skill Development**: Learn best practices and industry standards
5. **Time Management**: Focus on high-impact improvements with daily todo lists

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
git clone <your-repo-url>
cd Finternship
```

2. Set up environment variables:

```bash
export GITHUB_TOKEN=your_github_token_here
export OPENAI_API_KEY=your_openai_api_key_here
```

3. Run the Go backend:

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

### Frontend

- React
- Axios
- react-syntax-highlighter
- CSS with glass-morphism design

## Getting Started with API Keys

### GitHub Token

1. Go to GitHub.com → Settings → Developer Settings → Personal Access Tokens → Tokens (classic)
2. Generate a new token with `repo` scope
3. Set as `GITHUB_TOKEN` environment variable

### OpenAI API Key

1. Visit https://platform.openai.com/api-keys
2. Create a new API key
3. Set as `OPENAI_API_KEY` environment variable

## Security Notes

- Never commit API keys to version control
- Use environment variables for sensitive data
- Keep your dependencies updated

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
