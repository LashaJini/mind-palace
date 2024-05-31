from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.prompts.abstract import (
    Prompts,
    JoinableTemplateDict,
    JoinableTemplate,
)


class KeywordsPrompts(Prompts):
    format_start = "KEYWORDS:"
    format_end = "KEYWORDS_END"
    DEFAULT_KEYWORD_EXTRACT_TMPL = (
        "Some text is provided below. Given the text, extract up to {max_keywords} "
        "keywords from the text. Avoid stopwords.\n"
        "---------------------\n"
        "{context_str}\n"
        "---------------------\n"
        f"Provide keywords in the following comma-separated format: '{format_start} <keywords> {format_end}'\n"
    )
    DEFAULT_KEYWORD_EXTRACT_JOINABLE_TMPL: JoinableTemplateDict = {
        "instructions": "extract up to {max_keywords} keywords from the text.",
        "format": f"Provide keywords in the following comma-separated format: '{format_start} <keywords> {format_end}'",
    }

    tmpl = DEFAULT_KEYWORD_EXTRACT_TMPL
    joinable_tmpl = DEFAULT_KEYWORD_EXTRACT_JOINABLE_TMPL
    default_max_keywords = 5

    def standalone_template(self, verbose=False, **kwargs):
        if "max_keywords" not in kwargs:
            print("Missing 'max_keywords'")

        return super()._standalone_template(self.tmpl, verbose, **kwargs)

    def standalone_template_token_count(self, llm: CustomLlamaCPP):
        return llm.token_size(text=self.tmpl)

    def prompt(self, text: str, verbose=False, **kwargs):
        if "max_keywords" not in kwargs:
            kwargs["max_keywords"] = self.default_max_keywords

        result = self.standalone_template().format(context_str=text, **kwargs)
        if verbose:
            print(result)

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
