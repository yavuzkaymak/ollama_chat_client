package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	ollama "github.com/ollama/ollama/api"
)

var (
	defaultModel string = "qwen2.5:3b"
	modelName    string
)

type Agent struct {
	client         *ollama.Client
	getUserMessage func() (string, bool)
}

func main() {
	modelName = os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = defaultModel
		fmt.Printf("OLLAMA_MODEL environment variable not set, using default: %s\n", defaultModel)
	} else {
		fmt.Printf("Using model from OLLAMA_MODEL environment variable: %s\n", modelName)
	}

	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost != "" {
		fmt.Printf("OLLAMA_HOST environment variable set: %s\n", ollamaHost)
	} else {
		fmt.Println("OLLAMA_HOST environment variable not set, using default: http://localhost:11434")
	}



	ollamaClient, err := ollama.ClientFromEnvironment()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	getMessage := func() (string, bool) {
		if scanner.Scan() {
			return scanner.Text(), true
		}
		return "", false
	}

	agent := NewAgent(ollamaClient, getMessage)

	ctx, cancel := context.WithCancel(context.Background())
	setupSignalHandler(cancel)
	
	berr := agent.Run(ctx)
	if berr != nil {
		fmt.Printf("Error: %s\n", berr.Error())
	}
}

func setupSignalHandler(cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-signalChan
		fmt.Println("\nShutting down...")
		cancel()
		os.Exit(0)
	}()
}

func NewAgent(client *ollama.Client, getUserMessage func() (string, bool)) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
	}
}

func (a *Agent) Run(ctx context.Context) error {
	var modelContext []int // Store the context from previous exchanges

	fmt.Println("Chat with Ollama (use 'ctrl-c' or type ('exit') to quit)")

	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		userInput, ok := a.getUserMessage()
		if !ok {
			break
		}

		if strings.ToLower(userInput) == "exit" {
			fmt.Println("Exiting...")
			return nil
		}

		if strings.TrimSpace(userInput) == "" {
			continue
		}
		

		newContext, err := a.runInference(ctx, userInput, modelContext)
		if err != nil {
			return err
		}
		

		modelContext = newContext
		
		fmt.Println() 
	}

	return nil
}

func (a *Agent) runInference(ctx context.Context, prompt string, modelContext []int) ([]int, error) {
	// Create generate request with context from previous exchange
	req := &ollama.GenerateRequest{
		Model:    modelName,
		Prompt:   prompt,
		Context:  modelContext, // Pass the context from previous exchange
	}

	var newContext []int
	fmt.Printf("\u001b[93mOllama\u001b[0m: ")
	err := a.client.Generate(ctx, req, func(response ollama.GenerateResponse) error {
		fmt.Print(response.Response)		

		if response.Context != nil {
			newContext = response.Context
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return newContext, nil
}