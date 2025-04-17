# Ollama CLI Chat

A simple command-line interface (CLI) application written in Go to interact with a locally running Ollama instance. It supports streaming responses and maintains conversation context.

## Features

*   Connects to a local Ollama instance.
*   Streams responses from the language model.
*   Maintains conversation context between turns.
*   Allows specifying the Ollama model via an environment variable.
*   Graceful shutdown using Ctrl+C or typing `exit`.

## Prerequisites

*   **Go:** Version 1.18 or later installed (Installation Guide).
*   **Ollama:** A running Ollama instance on your local machine (Ollama Website). Make sure you have pulled the models you intend to use (e.g., `ollama pull qwen2.5:3b`).

## Building

To build the application, navigate to the project's root directory in your terminal and run:

```bash
go build -o ollama-chat .
