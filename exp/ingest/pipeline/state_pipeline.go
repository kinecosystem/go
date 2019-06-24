package pipeline

import (
	"github.com/stellar/go/exp/ingest/io"
	supportPipeline "github.com/stellar/go/exp/support/pipeline"
)

func (p *StatePipeline) Node(processor StateProcessor) *supportPipeline.PipelineNode {
	return &supportPipeline.PipelineNode{
		Processor: &stateProcessorWrapper{processor},
	}
}

func (p *StatePipeline) Process(readCloser io.StateReadCloser) <-chan error {
	return p.Pipeline.Process(&stateReadCloserWrapper{readCloser})
}
