package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/iskorotkov/cc-statusline/parts"
	"github.com/iskorotkov/cc-statusline/style"
)

var r = parts.Rows(
	parts.Row(
		style.Dim(style.Blue("CC")),
		parts.CCVersion(),
		parts.CCModel(),
		parts.CCOutputStyle(),
		parts.CCDir(),
		parts.CCStats(),
		parts.CC200KContextBadge(),
	),
	parts.Row(
		style.Dim(style.Blue("API")),
		parts.CCSessionUsage(),
		parts.CCHourUsage(),
		parts.CCDayUsage(),
		parts.CCWeekUsage(),
	),
	parts.Row(
		style.Dim(style.Blue("GIT")),
		parts.GitRemoteOrigin(),
		parts.GitBranch(),
		parts.GitStatus(),
		parts.GitDiffStats(),
	),
	parts.Row(
		style.Dim(style.Blue("PR")),
		parts.GHPRNumber(),
		parts.GHPRTitle(),
		parts.GHPRStats(),
	),
	parts.Row(
		style.Dim(style.Blue("PR")),
		parts.GHPRURL(),
	),
	parts.Row(
		style.Dim(style.Blue("TASK")),
		parts.GHIssueURL(),
		parts.JiraURL(),
	),
)

func main() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic: %#v", p)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	var hook parts.CCHook
	if err := json.NewDecoder(os.Stdin).Decode(&hook); err != nil {
		return err
	}
	s, err := r(ctx, hook)
	if err != nil {
		return err
	}
	fmt.Print(s)
	return nil
}
