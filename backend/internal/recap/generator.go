package recap

import "context"

type Generator interface {
	Generate(ctx context.Context, input GenerateInput) (GenerateOutput, error)
}
