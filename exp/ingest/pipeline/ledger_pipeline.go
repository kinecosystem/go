package pipeline

import (
	"github.com/stellar/go/exp/ingest/io"
	supportPipeline "github.com/stellar/go/exp/support/pipeline"
)

func (p *LedgerPipeline) Node(processor LedgerProcessor) *supportPipeline.PipelineNode {
	return &supportPipeline.PipelineNode{
		Processor: &ledgerProcessorWrapper{processor},
	}
}

func (p *LedgerPipeline) Process(readCloser io.LedgerReadCloser) <-chan error {
	return p.Pipeline.Process(&ledgerReadCloserWrapper{readCloser})
}
