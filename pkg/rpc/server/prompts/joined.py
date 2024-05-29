from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.prompts.abstract import (
    Prompts,
    JoinableTemplate,
)


class JoinedPrompts(Prompts):
    DEFAULT_JOINED_TMPL = (
        "Some text is provided below. Given the text: {instructions}\n"
        "Avoid stopwords. "
        "Try to use only the information provided. "
        "Try to include as many key details as possible.\n"
        "\n"
        "\n"
        "{context_str}\n"
        "\n"
        "\n"
        "{format}"
    )

    tmpl = DEFAULT_JOINED_TMPL

    def standalone_template(self, verbose=False, **kwargs):
        return super()._standalone_template(self.tmpl, verbose, **kwargs)

    def standalone_template_token_count(self, llm: CustomLlamaCPP):
        return llm.token_size(text=self.tmpl)

    def prompt(self, text: str, **kwargs):
        if "instructions" not in kwargs:
            kwargs["instructions"] = ""

        if "format" not in kwargs:
            kwargs["format"] = ""
        return self.standalone_template().format(context_str=text, **kwargs)

    def joinable_template(self, **kwargs) -> JoinableTemplate:
        return JoinableTemplate({"instructions": "", "format": ""})
