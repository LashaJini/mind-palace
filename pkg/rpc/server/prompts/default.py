from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.prompts.abstract import (
    Prompts,
    JoinableTemplate,
)


class DefaultPrompts(Prompts):
    def standalone_template(self, verbose=False, **kwargs):
        return ""

    def standalone_template_token_count(self, llm: CustomLlamaCPP):
        return 0

    def prompt(self, context_str: str, verbose=False, **kwargs):
        return ""

    def joinable_template(self, **kwargs) -> JoinableTemplate:
        return JoinableTemplate({"instructions": "", "format": ""})

    def joinable_template_token_count(self, llm: CustomLlamaCPP):
        return 0
