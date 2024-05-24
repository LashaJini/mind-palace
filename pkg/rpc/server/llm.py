from typing import List
from pathlib import Path
from llama_index.core import (
    Settings,
)

from llama_index.core.program import LLMTextCompletionProgram
from pkg.rpc.server.embeddings import EMBED_MODEL
from llama_index.llms.llama_cpp import LlamaCPP
from llama_index.llms.llama_cpp.llama_utils import DEFAULT_SYSTEM_PROMPT
from llama_index.core.prompts.default_prompts import (
    DEFAULT_SUMMARY_PROMPT_TMPL,
    DEFAULT_KEYWORD_EXTRACT_TEMPLATE_TMPL,
)

from pkg.rpc.server.output_parsers import (
    Keywords,
    Summary,
    SummaryParser,
    KeywordsParser,
)

# TODO: read from config
user_home = str(Path.home())
model_path = (
    user_home
    + "/.cache/lm-studio/models/lmstudio-community/Meta-Llama-3-8B-Instruct-GGUF/Meta-Llama-3-8B-Instruct-IQ3_M.gguf"
)


class CustomLLamaCPP(LlamaCPP):
    def __init__(self, **kwargs):
        super().__init__(
            model_path=model_path,
            temperature=0.1,
            max_new_tokens=256,
            context_window=3900,
            verbose=False,
            **kwargs,
        )

    @classmethod
    def summary_template(cls):
        prompt = DEFAULT_SUMMARY_PROMPT_TMPL
        return f"<|system|>\n{DEFAULT_SYSTEM_PROMPT}</s>\n<|user|>\n{prompt}</s>\n<|assistant|>\n"

    def summary_completion_to_prompt(self, completion):
        return CustomLLamaCPP.summary_template().format(context_str=completion)

    @classmethod
    def keywords_template(cls):
        prompt = DEFAULT_KEYWORD_EXTRACT_TEMPLATE_TMPL
        return f"<|system|>\n{DEFAULT_SYSTEM_PROMPT}</s>\n<|user|>\n{prompt}</s>\n<|assistant|>\n"

    def keywords_completion_to_prompt(self, completion):
        return CustomLLamaCPP.keywords_template().format(
            text=completion, max_keywords=15
        )

    @classmethod
    def gen_structured_summary(cls, llm, prompt, verbose=False):
        program = LLMTextCompletionProgram.from_defaults(
            llm=llm,
            output_parser=SummaryParser(verbose=verbose),
            output_cls=Summary,  # type:ignore
            prompt_template_str=CustomLLamaCPP.summary_template(),
            verbose=verbose,
        )

        return program(context_str=prompt)

    def gen_summary(self, prompt):
        return self.complete(
            prompt=self.summary_completion_to_prompt(prompt), formatted=True
        )

    @classmethod
    def gen_structured_keywords(cls, llm, prompt, verbose=False):
        """Should return list of keywords."""
        program = LLMTextCompletionProgram.from_defaults(
            llm=llm,
            output_parser=KeywordsParser(verbose=verbose),
            output_cls=Keywords,  # type:ignore
            prompt_template_str=CustomLLamaCPP.keywords_template(),
            verbose=verbose,
        )
        return program(text=prompt, max_keywords=15)

    def gen_keywords(self, prompt):
        return self.complete(
            prompt=self.keywords_completion_to_prompt(prompt), formatted=True
        )


llm = CustomLLamaCPP(
    generate_kwargs={
        "top_k": 1,  # TODO: config
        "stop": ["<|endoftext|>", "</s>"],  # TODO: wtf
    },
    # kwargs to pass to __init__()
    model_kwargs={
        "n_gpu_layers": -1,  # TODO: config
    },
)

Settings.llm = llm
Settings.embed_model = EMBED_MODEL
