package input_test

import (
	"testing"

	"github.com/halicoming/rod/lib/input"
	"github.com/halicoming/rod/lib/proto"
	"github.com/ysmood/got"
)

func TestMouseEncode(t *testing.T) {
	g := got.T(t)

	b, flag := input.EncodeMouseButton([]proto.InputMouseButton{proto.InputMouseButtonLeft})

	g.Eq(b, proto.InputMouseButtonLeft)
	g.Eq(flag, 1)
}
