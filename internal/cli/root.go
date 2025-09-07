package cli

import (
	"context"
	"fmt"
	"strings"

	"lilith/internal/cli/markdown"
	"lilith/internal/handler"
	"lilith/internal/infrastructure/adapters"
	"lilith/internal/infrastructure/anthropic"
	"lilith/internal/infrastructure/deepseek"
	"lilith/internal/infrastructure/session"

	"github.com/spf13/cobra"
)

type Providers struct {
	deepseek adapters.ICompletionAdapter
	claude   adapters.ICompletionAdapter
}

func NewRoot() *cobra.Command {
	app := &Providers{
		deepseek: deepseek.NewDeepSeekClient(),
		claude:   anthropic.NewAnthropicClient(),
	}

	root := &cobra.Command{
		Use:   "lilith",
		Short: "Lilith LLM CLI agent",
		Long: `
   _      _      _      _      _      _   
 _( )_  _( )_  _( )_  _( )_  _( )_  _( )_ 
(_ o _)(_ o _)(_ o _)(_ o _)(_ o _)(_ o _)
 (_,_)  (_,_)  (_,_)  (_,_)  (_,_)  (_,_) 
   _         _      _                 _   
 _( )_      //  .  //  . -/- /_     _( )_ 
(_ o _)   _(/__/__(/__/__/__/ (_   (_ o _)
 (_,_)                              (_,_) 
   _      _      _      _      _      _   
 _( )_  _( )_  _( )_  _( )_  _( )_  _( )_ 
(_ o _)(_ o _)(_ o _)(_ o _)(_ o _)(_ o _)
 (_,_)  (_,_)  (_,_)  (_,_)  (_,_)  (_,_) 

A LLM CLI agent.
Switch between models, perform analysis, discuss methodology, whatever.
Currently supports Anthropic and DeepSeek.
`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no prompt provided")
			}
			client, err := app.pickClient(cmd)
			if err != nil {
				return err
			}
			out, err := handler.Chat(context.Background(), client, strings.Join(args, " "))
			if err != nil {
				return err
			}
			markdown.CLIFormatter(out)
			return nil
		},
	}

	// global provider flags
	root.PersistentFlags().BoolP("cc", "c", false, "use Anthropic provider")
	root.PersistentFlags().BoolP("ds", "s", false, "use DeepSeek provider")

	// analyse
	analyse := &cobra.Command{
		Use:   "analyse [prompt]",
		Short: "Analyse code/issues in context. Includes write and cleanup options",
		RunE: func(cmd *cobra.Command, args []string) error {
			cleanup, _ := cmd.Flags().GetBool("cleanup")
			write, _ := cmd.Flags().GetBool("write")

			switch {
			case cleanup:
				return session.ResetMode("analyse")
			case write:
				return handler.WriteLatestFrom("analyse", "ANALYSIS.md")
			default:
				if len(args) == 0 {
					return fmt.Errorf("missing prompt")
				}
				client, err := app.pickClient(cmd)
				if err != nil {
					return err
				}
				out, err := handler.AnalyseChat(cmd.Context(), client, strings.Join(args, " "), false, false)
				if err != nil {
					return err
				}
				markdown.CLIFormatter(out)
				return nil
			}
		},
	}
	analyse.Flags().Bool("cleanup", false, "clear analysis session")
	analyse.Flags().Bool("write", false, "write last analysis to ANALYSIS.md")
	root.AddCommand(analyse)

	// discuss
	discuss := &cobra.Command{
		Use:   "discuss [prompt]",
		Short: "Discuss architecture/trade-offs. Includes write and cleanup options",
		RunE: func(cmd *cobra.Command, args []string) error {
			cleanup, _ := cmd.Flags().GetBool("cleanup")
			write, _ := cmd.Flags().GetBool("write")

			switch {
			case cleanup:
				return session.ResetMode("discuss")
			case write:
				return handler.WriteLatestFrom("discuss", "DISCUSSION.md")
			default:
				if len(args) == 0 {
					return fmt.Errorf("missing prompt")
				}
				client, err := app.pickClient(cmd)
				if err != nil {
					return err
				}
				out, err := handler.DiscussChat(cmd.Context(), client, strings.Join(args, " "), false, false)
				if err != nil {
					return err
				}
				markdown.CLIFormatter(out)
				return nil
			}
		},
	}
	discuss.Flags().Bool("cleanup", false, "clear discussion session")
	discuss.Flags().Bool("write", false, "write last discussion to DISCUSSION.md")
	root.AddCommand(discuss)

	return root
}

// ---- switch providers
func (provider *Providers) pickClient(cmd *cobra.Command) (adapters.ICompletionAdapter, error) {
	useClaude, _ := cmd.Flags().GetBool("cc")
	useDeepSeek, _ := cmd.Flags().GetBool("ds")

	if useClaude && useDeepSeek {
		return nil, fmt.Errorf("choose one provider: --cc/-c or --ds/-s (not both)")
	}
	if useClaude {
		if provider.claude == nil {
			return nil, fmt.Errorf("anthropic client not configured (API_KEY missing)")
		}
		return provider.claude, nil
	}
	if useDeepSeek {
		if provider.deepseek == nil {
			return nil, fmt.Errorf("deepseek client not configured (API_KEY missing)")
		}
		return provider.deepseek, nil
	}
	// default is DS
	if provider.deepseek != nil {
		return provider.deepseek, nil
	}
	return nil, fmt.Errorf("no provider available. Please check you have setup your providers")
}
