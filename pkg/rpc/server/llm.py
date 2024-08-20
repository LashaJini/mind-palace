from pathlib import Path
from typing import (
    ClassVar,
    List,
)
from llama_index.core.program import LLMTextCompletionProgram
from llama_index.core.base.embeddings.base import Embedding
from llama_index.core.base.llms.types import LLMMetadata
import numpy as np

from llama_index.embeddings.huggingface import HuggingFaceEmbedding
from llama_index.llms.llama_cpp import LlamaCPP

from pkg.rpc.server import logger

# TODO: read from config
user_home = str(Path.home())
model_path = (
    user_home
    + "/.cache/lm-studio/models/lmstudio-community/Meta-Llama-3-8B-Instruct-GGUF/Meta-Llama-3-8B-Instruct-IQ3_M.gguf"
)


class LLMConfig:
    context_size: ClassVar[int] = 8000
    context_window: ClassVar[int] = 3900
    max_new_tokens: ClassVar[int] = 1024

    def __init__(self, verbose=False, **kwargs):
        self.verbose = verbose
        self.kwargs = kwargs

    def update(self, **kwargs):
        if len(kwargs.keys()):
            self.kwargs.update(kwargs)
        else:
            self.kwargs = {}

    def __getattr__(self, name):
        if name in self.kwargs:
            return self.kwargs[name]
        return None

    def __setattr__(self, name, value):
        if name in ["verbose", "kwargs"]:
            super().__setattr__(name, value)
        else:
            self.kwargs[name] = value

    def __repr__(self):
        return f"ServerConfig(verbose={self.verbose}, kwargs={self.kwargs})"


class CustomLlamaCPP(LlamaCPP, extra="allow"):
    def __init__(self, llm_config: LLMConfig, **kwargs):
        super().__init__(
            model_path=model_path,
            temperature=0.1,
            max_new_tokens=llm_config.max_new_tokens,
            context_window=llm_config.context_window,
            verbose=llm_config.verbose,
            **kwargs,
        )
        self.llm_config = llm_config

    @property
    def metadata(self) -> LLMMetadata:
        return LLMMetadata(
            context_window=self.llm_config.context_window,
            num_output=self.llm_config.max_new_tokens,
            is_chat_model=False,
            is_function_calling_model=True,
            model_name="CustomLlamaCPP",
        )

    def text_completion_program(self, parser, output_cls, prompt):
        program = LLMTextCompletionProgram(
            llm=self,
            output_parser=parser,
            output_cls=output_cls,
            prompt=prompt,
            verbose=self.llm_config.verbose,
        )

        return program

    def encode(self, text: str):
        """get tokens"""
        return self._model.tokenizer().encode(text=text)

    def token_size(self, text: str):
        return len(self.encode(text))

    def calculate_available_tokens(
        self,
        decrements: List[int],
    ):
        available_tokens = (
            int(self.llm_config.available_tokens)
            if self.llm_config.available_tokens is not None
            else (self.llm_config.context_size - sum(decrements))
        )

        logger.log.debug(f"> Available tokens: {available_tokens}")

        return available_tokens


class EmbeddingModel(HuggingFaceEmbedding):
    # https://www.sbert.net/docs/sentence_transformer/pretrained_models.html
    _model_name = "sentence-transformers/all-MiniLM-L6-v2"
    _max_length: ClassVar[int] = 512
    dimension: ClassVar[int] = 384
    metric_type: ClassVar[str] = "COSINE"

    def __init__(self):
        super().__init__(
            model_name=EmbeddingModel._model_name,
            max_length=EmbeddingModel._max_length,
            cache_folder=user_home + "/.mind-palace/.cache/",
        )

    def _get_text_embeddings(self, texts: List[str]) -> List[List[float]]:
        return self._embed(texts, prompt_name="text")

    def embeddings(self, text: str) -> Embedding:
        sentences = text.split(".")

        sentence_embeddings = [
            self.get_text_embedding(sentence) for sentence in sentences if sentence
        ]
        aggregated_embedding = np.mean(sentence_embeddings, axis=0)
        return aggregated_embedding
