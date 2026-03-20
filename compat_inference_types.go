package gopenrouter

import (
	anthropicpkg "github.com/iamwavecut/gopenrouter/anthropic"
	embeddingspkg "github.com/iamwavecut/gopenrouter/embeddings"
	responsespkg "github.com/iamwavecut/gopenrouter/responses"
)

// Deprecated: use embeddings.Request from package github.com/iamwavecut/gopenrouter/embeddings.
type EmbeddingRequest = embeddingspkg.Request

// Deprecated: use embeddings.InputPart from package github.com/iamwavecut/gopenrouter/embeddings.
type EmbeddingInputPart = embeddingspkg.InputPart

// Deprecated: use embeddings.MultimodalInput from package github.com/iamwavecut/gopenrouter/embeddings.
type EmbeddingMultimodalInput = embeddingspkg.MultimodalInput

// Deprecated: use embeddings.Response from package github.com/iamwavecut/gopenrouter/embeddings.
type EmbeddingResponse = embeddingspkg.Response

// Deprecated: use embeddings.Datum from package github.com/iamwavecut/gopenrouter/embeddings.
type EmbeddingDatum = embeddingspkg.Datum

// Deprecated: use embeddings.Value from package github.com/iamwavecut/gopenrouter/embeddings.
type EmbeddingValue = embeddingspkg.Value

// Deprecated: use embeddings.Usage from package github.com/iamwavecut/gopenrouter/embeddings.
type EmbeddingUsage = embeddingspkg.Usage

// Deprecated: use anthropic.Request from package github.com/iamwavecut/gopenrouter/anthropic.
type AnthropicMessageRequest = anthropicpkg.Request

// Deprecated: use anthropic.Message from package github.com/iamwavecut/gopenrouter/anthropic.
type AnthropicMessage = anthropicpkg.Message

// Deprecated: use anthropic.Tool from package github.com/iamwavecut/gopenrouter/anthropic.
type AnthropicTool = anthropicpkg.Tool

// Deprecated: use anthropic.Response from package github.com/iamwavecut/gopenrouter/anthropic.
type AnthropicMessageResponse = anthropicpkg.Response

// Deprecated: use anthropic.ContentBlock from package github.com/iamwavecut/gopenrouter/anthropic.
type AnthropicContentBlock = anthropicpkg.ContentBlock

// Deprecated: use anthropic.Stream from package github.com/iamwavecut/gopenrouter/anthropic.
type AnthropicMessageStream = anthropicpkg.Stream

// Deprecated: use anthropic.StreamEvent from package github.com/iamwavecut/gopenrouter/anthropic.
type AnthropicMessageStreamEvent = anthropicpkg.StreamEvent

// Deprecated: use responses.Request from package github.com/iamwavecut/gopenrouter/responses.
type ResponseRequest = responsespkg.Request

// Deprecated: use responses.Tool from package github.com/iamwavecut/gopenrouter/responses.
type ResponseTool = responsespkg.Tool

// Deprecated: use responses.TextConfig from package github.com/iamwavecut/gopenrouter/responses.
type ResponseTextConfig = responsespkg.TextConfig

// Deprecated: use responses.ReasoningConfig from package github.com/iamwavecut/gopenrouter/responses.
type ResponseReasoningConfig = responsespkg.ReasoningConfig

// Deprecated: use responses.Response from package github.com/iamwavecut/gopenrouter/responses.
type Response = responsespkg.Response

// Deprecated: use responses.Usage from package github.com/iamwavecut/gopenrouter/responses.
type ResponseUsage = responsespkg.Usage

// Deprecated: use responses.OutputItem from package github.com/iamwavecut/gopenrouter/responses.
type ResponseOutputItem = responsespkg.OutputItem

// Deprecated: use responses.ContentPart from package github.com/iamwavecut/gopenrouter/responses.
type ResponseContentPart = responsespkg.ContentPart

// Deprecated: use responses.Stream from package github.com/iamwavecut/gopenrouter/responses.
type ResponseStream = responsespkg.Stream

// Deprecated: use responses.StreamEvent from package github.com/iamwavecut/gopenrouter/responses.
type ResponseStreamEvent = responsespkg.StreamEvent
