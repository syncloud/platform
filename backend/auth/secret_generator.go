package auth

import (
	"fmt"
	"github.com/syncloud/platform/cli"
	"strings"
)

type SecretGenerator struct {
	executor cli.Executor
}

type Secret struct {
	Password string
	Hash     string
}

func NewSecretGenerator(executor cli.Executor) *SecretGenerator {
	return &SecretGenerator{
		executor: executor,
	}
}

func (s *SecretGenerator) Generate() (Secret, error) {
	output, err := s.executor.CombinedOutput("snap", "run", "platform.authelia", "crypto", "hash", "generate", "--random")
	if err != nil {
		return Secret{}, err
	}

	parts := strings.Split(string(output), "\n")
	if len(parts) < 2 {
		return Secret{}, fmt.Errorf("not valid authelia crypto response: %s", string(output))
	}
	password, found := strings.CutPrefix(parts[0], "Random Password: ")
	if !found {
		return Secret{}, fmt.Errorf("not valid authelia crypto password: %s", parts[0])
	}
	hash, found := strings.CutPrefix(parts[1], "Digest: ")
	if !found {
		return Secret{}, fmt.Errorf("not valid authelia crypto hash: %s", parts[1])
	}
	return Secret{
		Password: password,
		Hash:     hash,
	}, err
}
