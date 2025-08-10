package payloadgenerator

import sequencegenerator "github.com/mascanio/uroboros/internal/sequence_generator"

type StaticMessageGenerator struct {
	message []byte
	seq     sequencegenerator.SequenceGenerator
}

func NewStaticMessageGenerator(
	message string,
	seq sequencegenerator.SequenceGenerator,
) *StaticMessageGenerator {
	rv := &StaticMessageGenerator{seq: seq}
	rv.message = []byte(message)
	return rv
}

func (g *StaticMessageGenerator) GenerateMessage() []byte {
	_, done := g.seq.Next()
	if done == sequencegenerator.Done {
		return nil
	}
	return g.message
}
