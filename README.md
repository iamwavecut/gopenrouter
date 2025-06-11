# gopenrouter

[![Go Reference](https://pkg.go.dev/badge/github.com/iamwavecut/gopenrouter.svg)](https://pkg.go.dev/github.com/iamwavecut/gopenrouter)
[![Go Report Card](https://goreportcard.com/badge/github.com/iamwavecut/gopenrouter)](https://goreportcard.com/report/github.com/iamwavecut/gopenrouter)
[![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-6e45e5.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

An unofficial Go client for the [OpenRouter API](https://openrouter.ai).

This library provides a comprehensive, `go-openai`-inspired interface for interacting with the OpenRouter API, giving you access to a multitude of LLMs through a single, unified client.

## Credits

This library's design and structure are heavily inspired by the excellent [go-openai](https://github.com/sashabaranov/go-openai) library.

## Key Differences from `go-openai`

While this library maintains a familiar, `go-openai`-style interface, it includes several key features and fields specifically tailored for the OpenRouter API:

- **Multi-Provider Models**: Seamlessly switch between models from different providers (e.g., Anthropic, Google, Mistral) by changing the `Model` string.
- **Cost Tracking**: The `Usage` object in responses includes a `Cost` field, providing direct access to the dollar cost of a generation.
- **Native Token Counts**: The `GetGeneration` endpoint provides access to `NativePromptTokens` and `NativeCompletionTokens`, giving you precise, provider-native tokenization data.
- **Advanced Routing**: Use `Models` for fallback chains and `Route` for custom routing logic.
- **Reasoning Parameters**: Control and request "thinking" tokens from supported models using the `Reasoning` parameters.
- **Provider-Specific `ExtraBody`**: Pass custom, provider-specific parameters through the `ExtraBody` field for fine-grained control.
- **Client Utilities**: Includes built-in methods to `ListModels`, `CheckCredits`, and `GetGeneration` stats directly from the client.

## Installation

```bash
go get github.com/iamwavecut/gopenrouter
```

## Usage

For complete, runnable examples, please see the [`examples/`](./examples) directory. A summary of available examples is below:

| Feature                                                      | Description                                                                                     |
| ------------------------------------------------------------ | ----------------------------------------------------------------------------------------------- |
| [Basic Chat](./examples/chat)                                | Demonstrates the standard chat completion flow.                                                 |
| [Streaming Chat](./examples/chat_stream)                     | Shows how to stream responses for real-time output.                                             |
| [Vision (Images)](./examples/chat_vision)                    | Illustrates how to send image data using the `MultiContent` field for vision-enabled models.    |
| [File Attachments](./examples/chat_file_attachment)          | Shows how to attach files (e.g., PDFs) for models that support file-based inputs.               |
| [Prompt Caching](./examples/chat_caching)                    | Reduces cost and latency by using OpenRouter's explicit `CacheControl` for supported providers. |
| [Automatic Caching (OpenAI)](./examples/chat_caching_openai) | Demonstrates OpenAI's automatic caching for large prompts, a cost-saving feature on OpenRouter. |
| [Structured Outputs](./examples/structured_output)           | Enforces a specific JSON schema for model outputs, a powerful OpenRouter feature.               |
| [Reasoning Tokens](./examples/chat_reasoning)                | Shows how to request and inspect the model's "thinking" process, unique to OpenRouter.          |
| [Provider Extras](./examples/chat_extra_body)                | Uses the `ExtraBody` field to pass provider-specific parameters for fine-grained control.       |
| [List Models](./examples/list_models)                        | A client utility to fetch the list of all models available on OpenRouter.                       |
| [Check Credits](./examples/check_credits)                    | A client utility to check your API key's usage, limit, and free tier status on OpenRouter.      |
| [Get Generation](./examples/get_generation)                  | Fetches detailed post-generation statistics, including cost and native token counts.            |

Details on specific features and client utility methods are available in the examples linked above. 