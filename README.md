# OpenAI Chat CLI Application in Go

A command-line interface application that interacts with OpenAI's GPT-3.5-turbo model while managing conversation context and tracking token usage. Project was a modified example of [this project from OpenAI](https://github.com/openai/openai-go/blob/main/examples/chat-completion-streaming/main.go)
## Features

- Real-time conversation with GPT-3.5-turbo
- Token usage tracking and context window management
- Streaming responses for immediate feedback
- Session statistics display
- Context window limit enforcement (4096 tokens)

## Prerequisites

- Go 1.23.2 or higher
- OpenAI API key set as an environment variable

## Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```
3. Enable Environment Variable
```bash
export OPENAI_API_KEY='your-api-key-here'
```

## Features

- Real-time conversation with GPT-3.5-turbo
- Token usage tracking and context window management
- Streaming responses for immediate feedback
- Session statistics display
- Context window limit enforcement (4096 tokens)

## Prerequisites

- Go 1.23.2 or higher
- OpenAI API key set as an environment variable

## Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```
3. Enable Environment Variable
```bash
export OPENAI_API_KEY='your-api-key-here'
```