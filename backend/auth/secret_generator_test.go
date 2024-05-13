package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type GeneratorExecutorStub struct {
	output string
}

func (e *GeneratorExecutorStub) CombinedOutput(name string, arg ...string) ([]byte, error) {
	return []byte(e.output), nil
}

func TestSecretGenerator_Generate(t *testing.T) {

	generator := &SecretGenerator{executor: &GeneratorExecutorStub{
		`Random Password: Dtf0qf8eoVJBaSCPU7hbYFUAOKbBahP5Pgf9VTDssA17dfCG1ilFph6PrBljr1aQFhCfy3TW
Digest: $argon2id$v=19$m=65536,t=3,p=4$7TGF3l00V1y6pJQPqmalnQ$qziJ0fKC23V/tpGhze97RJ9TbVbKae3CyZCcBiqmI5I
`}}
	secret, err := generator.Generate()
	assert.NoError(t, err)

	assert.Equal(t, "$argon2id$v=19$m=65536,t=3,p=4$7TGF3l00V1y6pJQPqmalnQ$qziJ0fKC23V/tpGhze97RJ9TbVbKae3CyZCcBiqmI5I", secret.Hash)
	assert.Equal(t, "Dtf0qf8eoVJBaSCPU7hbYFUAOKbBahP5Pgf9VTDssA17dfCG1ilFph6PrBljr1aQFhCfy3TW", secret.Password)
}
