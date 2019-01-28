package sse

import (
	"context"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stellar/go/support/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StreamTestSuite struct {
	suite.Suite
	ctx    context.Context
	w      *httptest.ResponseRecorder
	stream Stream
}

// Helper method to check that the preamble has been sent and all HTTP response headers are correctly set.
func (suite *StreamTestSuite) checkHeadersAndPreamble() {
	assert.Equal(suite.T(), "text/event-stream; charset=utf-8", suite.w.Header().Get("Content-Type"))
	assert.Equal(suite.T(), "no-cache", suite.w.Header().Get("Cache-Control"))
	assert.Equal(suite.T(), "keep-alive", suite.w.Header().Get("Connection"))
	assert.Equal(suite.T(), "*", suite.w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), 200, suite.w.Code)
	assert.Contains(suite.T(), suite.w.Body.String(), "retry: 1000\nevent: open\ndata: \"hello\"\n\n")
}

func (suite *StreamTestSuite) SetupTest() {
	suite.ctx, _ = test.ContextWithLogBuffer()
	suite.w = httptest.NewRecorder()
	suite.stream = NewStream(suite.ctx, suite.w)
}

// Tests that the stream sends the preamble before any events and that events are correctly sent.
func (suite *StreamTestSuite) TestStream_Send() {
	e := Event{Data: "test message"}
	suite.stream.Send(e)
	// Before sending, it should have sent the preamble first and set the headers.
	suite.checkHeadersAndPreamble()
	// Now check that the data got written
	assert.Contains(suite.T(), suite.w.Body.String(), "data: \"test message\"\n\n")
	suite.stream.Done()
	assert.Equal(suite.T(), 1, suite.stream.SentCount())
}

// Tests that Stream can send error events.
func (suite *StreamTestSuite) TestStream_Err() {
	err := errors.New("example error")
	suite.stream.Err(err)
	// Even if no events have been sent, Err should still send the preamble before the error event.
	suite.checkHeadersAndPreamble()
	assert.Contains(suite.T(), suite.w.Body.String(), "event: err\ndata: example error\n\n")
	assert.True(suite.T(), suite.stream.IsDone())
}

// Tests that SetLimit sets stream.done to true after the limit has been reached.
func (suite *StreamTestSuite) TestStream_SetLimit() {
	suite.stream.SetLimit(3)
	// Send more than the limit
	for i := 0; i < 5; i++ {
		message := "test message " + strconv.Itoa(i)
		suite.stream.Send(Event{Data: message})
	}
	assert.True(suite.T(), suite.stream.IsDone())
}

// Tests that SentCount reports the correct number.
func (suite *StreamTestSuite) TestStream_SentCount() {
	for i := 0; i < 5; i++ {
		message := "test message " + strconv.Itoa(i)
		suite.stream.Send(Event{Data: message})
	}
	assert.Equal(suite.T(), 5, suite.stream.SentCount())
	suite.stream.Err(errors.New("example error"))
	// Make sure that errors don't contribute to the send count
	assert.Equal(suite.T(), 5, suite.stream.SentCount())
}

// Runs the test suite.
func TestStreamTestSuite(t *testing.T) {
	suite.Run(t, new(StreamTestSuite))
}
