from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.prompts.abstract import (
    Prompts,
    JoinableTemplateDict,
    JoinableTemplate,
)


class SummaryPrompts(Prompts):
    format_start = "SUMMARY:"
    format_end = "SUMMARY_END"
    DEFAULT_SUMMARY_TMPL = (
        "Write a summary of the following. Try to use only the "
        "information provided. "
        "Try to include as many key details as possible.\n"
        "\n"
        "\n"
        "{context_str}\n"
        "\n"
        "\n"
        f"Provide summary in the following format: '{format_start} <summary> {format_end}'\n"
    )
    DEFAULT_SUMMARY_JOINABLE_TMPL: JoinableTemplateDict = {
        "instructions": "write summary",
        "format": f"Provide summary in the following format: '{format_start} <summary> {format_end}'",
    }

    tmpl = DEFAULT_SUMMARY_TMPL
    joinable_tmpl = DEFAULT_SUMMARY_JOINABLE_TMPL

    def standalone_template(self, verbose=False, **kwargs):
        return super()._standalone_template(self.tmpl, verbose, **kwargs)

    def standalone_template_token_count(self, llm: CustomLlamaCPP):
        return llm.token_size(text=self.tmpl)

    def prompt(self, text: str, verbose=False, **kwargs):
        result = self.standalone_template().format(context_str=text, **kwargs)
        if verbose:
            print(result)

        return result

    def joinable_template(self, **kwargs) -> JoinableTemplate:
        return JoinableTemplate(
            {
                "instructions": self.joinable_tmpl["instructions"].format(**kwargs),
                "format": self.joinable_tmpl["format"].format(**kwargs),
            }
        )
