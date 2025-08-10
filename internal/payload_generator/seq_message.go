package payloadgenerator

import (
	"math"
	"strconv"

	sequencegenerator "github.com/mascanio/uroboros/internal/sequence_generator"
)

type SeqMessageGenerator struct {
	base []byte
	buf  []byte
	seq  sequencegenerator.SequenceGenerator
}

func NewIDMessageGenerator(
	baseMessage string,
	seq sequencegenerator.SequenceGenerator,
) *SeqMessageGenerator {
	rv := &SeqMessageGenerator{seq: seq}
	rv.base = []byte(baseMessage)
	rv.buf = make([]byte, 0, len(baseMessage)+len(strconv.Itoa(math.MaxInt64)))
	rv.buf = append(rv.buf, "id: "...)
	return rv
}

func (g *SeqMessageGenerator) GenerateMessage() []byte {
	i, done := g.seq.Next()
	if done == sequencegenerator.Done {
		return nil
	}
	g.buf = strconv.AppendInt(g.buf[:4], int64(i), 10)
	g.buf = append(g.buf, ' ')
	g.buf = append(g.buf, g.base...)
	return g.buf
}
