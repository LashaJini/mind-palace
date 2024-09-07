from pkg.rpc.loggers.palace import log
from pkg.rpc.services.llm.llm import CustomLlamaCPP
from pkg.rpc.services.addon.prompts.abstract import (
    Prompts,
    JoinableTemplate,
)

DEFAULT_JOINED_TMPL = (
    "Some text is provided below. Given the text: {instructions}.\n"
    "Avoid stopwords. "
    "Try to use only the information provided. "
    "Try to include as many key details as possible.\n"
    "\n"
    "\n"
    "{context_str}\n"
    "\n"
    "\n"
    "{format}\n\n"
    "Output formats should not overlap with each other.\n"
)


class JoinedPrompts(Prompts):
    tmpl = DEFAULT_JOINED_TMPL

    def standalone_template(self, verbose=False, **kwargs):
        return super()._standalone_template(self.tmpl, verbose, **kwargs)

    def standalone_template_token_count(self, llm: CustomLlamaCPP):
        return llm.token_size(text=self.tmpl)

    def prompt(self, context_str: str, verbose=False, **kwargs):
        kwargs["instructions"] = kwargs.get("instructions", "")
        kwargs["format"] = kwargs.get("format", "")

        result = self.standalone_template().format(
            context_str=context_str, verbose=verbose, **kwargs
        )

        log.debug(f"Prompt: {result}")

        return result

    def joinable_template(self, **kwargs) -> JoinableTemplate:
        return JoinableTemplate({"instructions": "", "format": ""})
