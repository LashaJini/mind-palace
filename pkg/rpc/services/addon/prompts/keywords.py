from pkg.rpc.loggers.palace import log
from pkg.rpc.services.llm.llm import CustomLlamaCPP
from pkg.rpc.services.addon.prompts.abstract import (
    Prompts,
    JoinableTemplateDict,
    JoinableTemplate,
)


class KeywordsPrompts(Prompts):
    format_start = "KEYWORDS:"
    format_end = "KEYWORDS_END"
    DEFAULT_KEYWORD_EXTRACT_TMPL = (
        "Some text is provided below. "
        "Extract the most relevant keywords only from the text within <CHUNK></CHUNK> sections. "
        "Each keyword should represent significant concepts, names, or terms from the chunk. "
        "Avoid stopwords.\n"
        "\n"
        "\n"
        "{context_str}\n"
        "\n"
        "\n"
        f"For each chunk, provide up to {{max_keywords}} keywords in the following comma-separated format: '{format_start} <keywords> {format_end}'\n"
    )
    DEFAULT_KEYWORD_EXTRACT_JOINABLE_TMPL: JoinableTemplateDict = {
        "instructions": "extract the most relevant keywords only from the text within <CHUNK></CHUNK> sections. Each keyword should represent significant concepts, names, or terms from the chunk",
        "format": f"For each chunk, provide up to {{max_keywords}} keywords in the following comma-separated format: '{format_start} <keywords> {format_end}'",
    }

    tmpl = DEFAULT_KEYWORD_EXTRACT_TMPL
    joinable_tmpl = DEFAULT_KEYWORD_EXTRACT_JOINABLE_TMPL
    default_max_keywords = 5

    def standalone_template(self, verbose=False, **kwargs):
        if "max_keywords" not in kwargs:
            log.warning("Missing 'max_keywords'")

        return super()._standalone_template(self.tmpl, verbose, **kwargs)

    def standalone_template_token_count(self, llm: CustomLlamaCPP):
        return llm.token_size(text=self.tmpl)

    def prompt(self, context_str: str, verbose=False, **kwargs):
        if "max_keywords" not in kwargs:
            kwargs["max_keywords"] = self.default_max_keywords

        result = self.standalone_template().format(
            context_str=context_str, verbose=verbose, **kwargs
        )

        log.debug(f"Prompt: {result}")

        return result

    def joinable_template(self, **kwargs) -> JoinableTemplate:
        if "max_keywords" not in kwargs:
            kwargs["max_keywords"] = self.default_max_keywords

        return JoinableTemplate(
            {
                "instructions": self.joinable_tmpl["instructions"].format(**kwargs),
                "format": self.joinable_tmpl["format"].format(**kwargs),
            }
        )
