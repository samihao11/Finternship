import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Light as SyntaxHighlighter } from 'react-syntax-highlighter';
import { github } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import './App.css';

// Import language support
import javascript from 'react-syntax-highlighter/dist/esm/languages/hljs/javascript';
import python from 'react-syntax-highlighter/dist/esm/languages/hljs/python';
import java from 'react-syntax-highlighter/dist/esm/languages/hljs/java';
import go from 'react-syntax-highlighter/dist/esm/languages/hljs/go';
import css from 'react-syntax-highlighter/dist/esm/languages/hljs/css';
import xml from 'react-syntax-highlighter/dist/esm/languages/hljs/xml';
import markdown from 'react-syntax-highlighter/dist/esm/languages/hljs/markdown';

// Register languages
SyntaxHighlighter.registerLanguage('javascript', javascript);
SyntaxHighlighter.registerLanguage('python', python);
SyntaxHighlighter.registerLanguage('java', java);
SyntaxHighlighter.registerLanguage('go', go);
SyntaxHighlighter.registerLanguage('css', css);
SyntaxHighlighter.registerLanguage('xml', xml);
SyntaxHighlighter.registerLanguage('markdown', markdown);

function LandingPage({ onSubmit }) {
  const [username, setUsername] = useState('');
  const [reponame, setReponame] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!username.trim() || !reponame.trim()) {
      setError('Please enter both username and repository name');
      return;
    }
    onSubmit(username.trim(), reponame.trim());
  };

  return (
    <div className="landing-page">
      <div className="landing-content">
        <h1>Finternship</h1>
        <p>Enter a GitHub username and repository name to get started</p>
        {error && <div className="error-message">{error}</div>}
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            placeholder="GitHub Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <input
            type="text"
            placeholder="Repository Name"
            value={reponame}
            onChange={(e) => setReponame(e.target.value)}
          />
          <button type="submit">View Repository</button>
        </form>
      </div>
    </div>
  );
}

function TodoPanel({ username, reponame }) {
  const [todos, setTodos] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const generateTodos = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await axios.get(`http://localhost:8080/todos/${username}/${reponame}`);
      setTodos(response.data.todos);
    } catch (err) {
      setError('Failed to generate todos. Please try again.');
      console.error('Error:', err);
    }
    setLoading(false);
  };

  const handleResolve = (index) => {
    setTodos(prevTodos => prevTodos.filter((_, i) => i !== index));
  };

  return (
    <div className="todo-panel">
      <div className="todo-header">
        <h2>Today's To-Dos</h2>
        <button onClick={generateTodos} disabled={loading} className="generate-button">
          {loading ? 'Generating...' : 'Generate To-Dos'}
        </button>
      </div>
      {error && <div className="error-container">{error}</div>}
      <div className="todo-list">
        {todos.length === 0 ? (
          <p>Click "Generate To-Dos" to get started</p>
        ) : (
          todos.map((todo, index) => (
            <div key={index} className="todo-item">
              <span>{todo}</span>
              <button onClick={() => handleResolve(index)} className="resolve-button">
                Resolved
              </button>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

function RepoViewer({ username, reponame, onBack }) {
  const [repoData, setRepoData] = useState(null);
  const [selectedFile, setSelectedFile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchRepoData();
  }, [username, reponame]);

  const fetchRepoData = async () => {
    try {
      setLoading(true);
      const response = await axios.get(`http://localhost:8080/repo/${username}/${reponame}`);
      setRepoData(response.data);
      setLoading(false);
    } catch (err) {
      setError('Failed to fetch repository data');
      setLoading(false);
    }
  };

  const handleFileClick = (file) => {
    setSelectedFile(file);
  };

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  if (error) {
    return (
      <div className="error-container">
        <div className="error">{error}</div>
        <button onClick={onBack} className="back-button">Go Back</button>
      </div>
    );
  }

  return (
    <div className="app">
      <header className="header">
        <div className="header-content">
          <h1>Finternship</h1>
          <button onClick={onBack} className="back-button">Choose Different Repo</button>
        </div>
        {repoData && (
          <div className="repo-info">
            <span>{repoData.githubUsername}</span>
            <span>/</span>
            <span>{repoData.repoName}</span>
          </div>
        )}
      </header>

      <div className="content">
        <div className="file-list">
          <h2>Files</h2>
          {repoData?.files.map((file, index) => (
            <div
              key={index}
              className={`file-item ${selectedFile === file ? 'selected' : ''}`}
              onClick={() => handleFileClick(file)}
            >
              {file.path}
            </div>
          ))}
        </div>

        <div className="main-content">
          <div className="file-viewer">
            {selectedFile ? (
              <>
                <div className="file-header">
                  <h3>{selectedFile.path}</h3>
                </div>
                <div className="code-container">
                  <SyntaxHighlighter
                    language={selectedFile.language}
                    style={github}
                    showLineNumbers
                    customStyle={{
                      margin: 0,
                      borderRadius: '0 0 4px 4px',
                    }}
                  >
                    {selectedFile.content}
                  </SyntaxHighlighter>
                </div>
              </>
            ) : (
              <div className="no-file">
                Select a file from the list to view its contents
              </div>
            )}
          </div>
          
          <TodoPanel username={username} reponame={reponame} />
        </div>
      </div>
    </div>
  );
}

function App() {
  const [selectedRepo, setSelectedRepo] = useState(null);

  const handleRepoSubmit = (username, reponame) => {
    setSelectedRepo({ username, reponame });
  };

  const handleBack = () => {
    setSelectedRepo(null);
  };

  return (
    selectedRepo ? (
      <RepoViewer
        username={selectedRepo.username}
        reponame={selectedRepo.reponame}
        onBack={handleBack}
      />
    ) : (
      <LandingPage onSubmit={handleRepoSubmit} />
    )
  );
}

export default App; 