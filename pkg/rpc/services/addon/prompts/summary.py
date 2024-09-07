from pkg.rpc.loggers.palace import log
from pkg.rpc.services.llm.llm import CustomLlamaCPP
from pkg.rpc.services.addon.prompts.abstract import (
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
        f"**Summary:** Provide the summary in the following format: `{format_start} <summary> {format_end}`. Ensure that you end the summary with `{format_end}`"
    )
    DEFAULT_SUMMARY_JOINABLE_TMPL: JoinableTemplateDict = {
        "instructions": "write summary",
        "format": f"**Summary:** Provide the summary in the following format: `{format_start} <summary> {format_end}`. Ensure that you end the summary with `{format_end}`",
    }

    tmpl = DEFAULT_SUMMARY_TMPL
    joinable_tmpl = DEFAULT_SUMMARY_JOINABLE_TMPL

    def standalone_template(self, verbose=False, **kwargs):
        return super()._standalone_template(self.tmpl, verbose, **kwargs)

    def standalone_template_token_count(self, llm: CustomLlamaCPP):
        return llm.token_size(text=self.tmpl)

    def prompt(self, context_str: str, verbose=False, **kwargs):
        result = self.standalone_template().format(
            context_str=context_str, verbose=verbose, **kwargs
        )

        log.debug(f"Prompt: {result}")

        return result

    def joinable_template(self, **kwargs) -> JoinableTemplate:
        return JoinableTemplate(
            {
                "instructions": self.joinable_tmpl["instructions"].format(**kwargs),
                "format": self.joinable_tmpl["format"].format(**kwargs),
            }
        )
