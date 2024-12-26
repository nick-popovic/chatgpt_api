// Package main provides a CLI interface to interact with OpenAI's GPT models
// while managing conversation context and tracking token usage.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/pkoukk/tiktoken-go"
)

// MaxTokens defines the maximum context window size for GPT-3.5-turbo
// This limit includes both input and output tokens.
var (
	MaxTokens = 4096 // GPT-3.5-turbo default context window
)

// countTokens calculates the number of tokens in the given text using the
// cl100k_base encoding used by GPT-3.5-turbo and GPT-4.
// Returns the token count and any error encountered during encoding.
func countTokens(text string) (int, error) {
	encoding := "cl100k_base" // encoding for GPT-3.5-turbo and GPT-4
	tkm, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return 0, err
	}
	tokens := tkm.Encode(text, nil, nil)
	return len(tokens), nil
}

// SessionStats tracks the conversation metrics including total tokens used,
// number of messages exchanged, and remaining tokens in the context window.
type SessionStats struct {
	totalTokens     int // Total tokens used in the conversation
	messageCount    int // Number of messages exchanged
	remainingTokens int // Remaining tokens in the context window
}

// update recalculates the session statistics after processing new tokens.
// It updates the total tokens used, increments message count, and
// adjusts the remaining tokens available in the context window.
func (s *SessionStats) update(newTokens int) {
	s.totalTokens += newTokens
	s.messageCount++
	s.remainingTokens = MaxTokens - s.totalTokens
}

// display prints the current session statistics to standard output,
// showing the number of messages exchanged, total tokens used,
// and remaining tokens in the context window.
func (s *SessionStats) display() {
	fmt.Printf("\n=== Session Stats ===\n")
	fmt.Printf("Messages: %d\n", s.messageCount)
	fmt.Printf("Tokens Used: %d\n", s.totalTokens)
	fmt.Printf("Tokens Remaining: %d\n", s.remainingTokens)
	fmt.Printf("Context Usage: %.1f%%\n", float64(s.totalTokens)/float64(MaxTokens)*100)
	fmt.Println("===================\n")
}

// main initializes the ChatGPT client and runs the main conversation loop,
// handling user input, tracking token usage, and maintaining conversation context.
func main() {
	client := openai.NewClient()
	ctx := context.Background()
	scanner := bufio.NewScanner(os.Stdin)

	stats := &SessionStats{
		remainingTokens: MaxTokens,
	}

	// Initialize conversation history
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful assistant."),
	}

	for {
		fmt.Print("\nEnter your question (or 'q' to exit): ") // Updated prompt
		scanner.Scan()
		input := scanner.Text()

		if input == "q" {
			fmt.Println("Final Stats:")
			stats.display()
			break
		}

		// Count tokens for new input
		inputTokens, err := countTokens(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error counting tokens: %v\n", err)
			continue
		}

		if inputTokens+stats.totalTokens > MaxTokens {
			fmt.Printf("Warning: Adding this input would exceed token limit (%d + %d > %d)\n",
				stats.totalTokens, inputTokens, MaxTokens)
			continue
		}

		// Add user message to history
		messages = append(messages, openai.UserMessage(input))
		stats.update(inputTokens)

		// Echo the question
		fmt.Printf("\n> %s\n\n", input)

		stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Messages: openai.F(messages),
			Seed:     openai.Int(1),
			Model:    openai.F(openai.ChatModelGPT3_5Turbo),
		})

		// Collect assistant's response
		var assistantResponse strings.Builder
		for stream.Next() {
			evt := stream.Current()
			if len(evt.Choices) > 0 {
				content := evt.Choices[0].Delta.Content
				fmt.Print(content)
				assistantResponse.WriteString(content)
			}
		}
		fmt.Println()

		// Add assistant's response to history
		messages = append(messages, openai.AssistantMessage(assistantResponse.String()))

		responseTokens, _ := countTokens(assistantResponse.String())
		stats.update(responseTokens)
		stats.display()

		if err := stream.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error from OpenAI:", err)
			continue
		}
	}
}
