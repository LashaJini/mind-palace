from typing import List
from pathlib import Path
from llama_index.core import (
    Settings,
)
from llama_index.core.bridge.pydantic import BaseModel
from llama_index.core.program import LLMTextCompletionProgram
from pkg.rpc.server.embeddings import EMBED_MODEL
from llama_index.llms.llama_cpp import LlamaCPP
from llama_index.llms.llama_cpp.llama_utils import DEFAULT_SYSTEM_PROMPT
from llama_index.core.prompts.default_prompts import (
    DEFAULT_SUMMARY_PROMPT_TMPL,
    DEFAULT_KEYWORD_EXTRACT_TEMPLATE_TMPL,
)

# TODO: read from config
user_home = str(Path.home())
model_path = (
    user_home
    + "/.cache/lm-studio/models/lmstudio-community/Meta-Llama-3-8B-Instruct-GGUF/Meta-Llama-3-8B-Instruct-IQ3_M.gguf"
)


class Summary(BaseModel):
    """Data model for a summary."""

    text: str


class Song(BaseModel):
    """Data model for a song."""

    title: str
    length_seconds: int


class Album(BaseModel):
    """Data model for an album."""

    name: str
    artist: str
    songs: List[Song]


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

    def summary_completion_to_prompt(self, completion):
        prompt = DEFAULT_SUMMARY_PROMPT_TMPL.format(context_str=completion)
        return f"<|system|>\n{DEFAULT_SYSTEM_PROMPT}</s>\n<|user|>\n{prompt}</s>\n<|assistant|>\n"

    def keywords_completion_to_prompt(self, completion):
        prompt = DEFAULT_KEYWORD_EXTRACT_TEMPLATE_TMPL.format(
            max_keywords=15, text=completion
        )
        return f"<|system|>\n{DEFAULT_SYSTEM_PROMPT}</s>\n<|user|>\n{prompt}</s>\n<|assistant|>\n"

    def gen_summary(self, prompt):
        return self.complete(
            prompt=self.summary_completion_to_prompt(prompt), formatted=True
        )

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
