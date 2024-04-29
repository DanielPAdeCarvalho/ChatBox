package chat

import (
	"context"
	"encoding/json"
	"log"

	"chat-bot/pkg/session"

	"github.com/jdkato/prose/v2"
)

type Response struct {
	Message   string `json:"message"`
	SessionID string `json:"sessionID,omitempty"`
}

func ProcessMessage(ctx context.Context, sessionID, msg string) string {
	var response Response
	doc, err := prose.NewDocument(msg)
	if err != nil {
		log.Println("Error creating Prose document:", err)
		response.Message = "I'm having trouble understanding that."
		return encodeJSON(response)
	}

	// Check if we are expecting a product name
	expectingProduct := session.Get(ctx, sessionID, "expecting_product") == "true"
	if expectingProduct {
		productName := extractProductName(doc)
		session.Set(ctx, sessionID, "expecting_product", "false") // Reset expectation
		if productName != "" {
			session.Set(ctx, sessionID, "last_ordered_product", productName)
			return "I've added " + productName + " to your cart. Anything else?"
		} else {
			return "I'm sorry, I still didn't catch that. What product were you interested in?"
		}
	}

	name := session.Get(ctx, sessionID, "name")
	if name == "" {
		response.Message = handleInitialGreeting(ctx, sessionID, doc)
	} else {
		response.Message = handleUserCommand(ctx, sessionID, doc, name)
	}
	return encodeJSON(response)
}

func handleInitialGreeting(ctx context.Context, sessionID string, doc *prose.Document) string {
	found := false
	for _, tok := range doc.Tokens() {
		log.Printf("Token: %s, POS Tag: %s", tok.Text, tok.Tag) // Log every token and its tag
		if tok.Tag == "NNP" {
			session.Set(ctx, sessionID, "name", tok.Text)
			found = true
			break
		}
	}
	if found {
		return "Nice to meet you, " + doc.Text + "! How can I help you today?"
	}
	return "Hello! I don't believe we've met. What's your name?"
}

func handleUserCommand(ctx context.Context, sessionID string, doc *prose.Document, name string) string {
	tokens := doc.Tokens()
	ordered := false
	productName := ""
	for i, tok := range tokens {
		if tok.Text == "order" && (tok.Tag == "VB" || tok.Tag == "VBP") { // Verb 'order'
			// Try to find the product name next to the verb 'order'
			if product := getNextNoun(tokens, i); product != "" {
				productName = product
				ordered = true
				break
			}
		}
	}

	if ordered {
		if productName == "" {
			session.Set(ctx, sessionID, "expecting_product", "true")
			return "What product would you like to order?"
		} else {
			session.Set(ctx, sessionID, "expecting_product", "false")
			return "Excellent choice, " + name + "! I've added " + productName + " to your cart..."
		}
	}
	return "Sorry, " + name + ", I didn't understand that. Could you specify what you would like to do?"
}

func extractProductName(doc *prose.Document) string {
	// Simplistic implementation, can be improved
	for _, ent := range doc.Entities() {
		if ent.Label == "PRODUCT" { // Label for entities as products
			return ent.Text
		}
	}
	return ""
}

func getNextNoun(tokens []prose.Token, index int) string {
	// Helper function to find the next noun after a given index
	for i := index + 1; i < len(tokens); i++ {
		if tokens[i].Tag == "NN" || tokens[i].Tag == "NNP" { // Checking for common noun (NN) or proper noun (NNP)
			return tokens[i].Text
		}
	}
	return ""
}

func encodeJSON(response Response) string {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("JSON encoding error: %v", err)
		return ""
	}
	return string(jsonBytes)
}
