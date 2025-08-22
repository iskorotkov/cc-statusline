package shell

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func String(ctx context.Context, s ...string) (string, error) {
	cmd := exec.CommandContext(ctx, s[0], s[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: \n\n %s", err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

func JSON[T any](ctx context.Context, s ...string) (T, error) {
	var result T
	cmd := exec.CommandContext(ctx, s[0], s[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return result, fmt.Errorf("%w: \n\n %s", err, strings.TrimSpace(string(output)))
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return result, err
	}
	return result, nil
}
