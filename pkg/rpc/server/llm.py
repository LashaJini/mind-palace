from pathlib import Path
from typing import ClassVar, List
from llama_index.core.base.embeddings.base import BaseEmbedding, Embedding
import numpy as np

from llama_index.embeddings.huggingface import HuggingFaceEmbedding
from llama_index.llms.llama_cpp import LlamaCPP

# TODO: read from config
user_home = str(Path.home())
model_path = (
    user_home
    + "/.cache/lm-studio/models/lmstudio-community/Meta-Llama-3-8B-Instruct-GGUF/Meta-Llama-3-8B-Instruct-IQ3_M.gguf"
)


class CustomLlamaCPP(LlamaCPP):
    context_size: ClassVar[int] = 8000

    def __init__(self, verbose: bool = False, **kwargs):
        super().__init__(
            model_path=model_path,
            temperature=0.1,
            max_new_tokens=1024,
            context_window=3900,
            verbose=verbose,
            **kwargs,
        )

    def encode(self, text: str):
        """get tokens"""
        return self._model.tokenizer().encode(text=text)

    def token_size(self, text: str):
        return len(self.encode(text))

    def calculate_available_tokens(
        self,
        input_text_token_count,
        sys_prompt_token_count,
        joined_prompt_token_count,
        verbose=False,
    ):
        available_tokens = (
            CustomLlamaCPP.context_size
            - input_text_token_count
            - sys_prompt_token_count
            - joined_prompt_token_count
        )

        if verbose:
            print(f"> Available tokens: {available_tokens}")

        return available_tokens


class EmbeddingModel(HuggingFaceEmbedding):
    # https://www.sbert.net/docs/sentence_transformer/pretrained_models.html
    _model_name = "sentence-transformers/all-MiniLM-L6-v2"
    _max_length = 512
    dimension = 384
    metric_type = "COSINE"

    def __init__(self):
        super().__init__(
            model_name=EmbeddingModel._model_name,
            max_length=EmbeddingModel._max_length,
            cache_folder=user_home + "/.mind-palace/.cache/",
        )

    def _get_text_embeddings(self, texts: List[str]) -> List[List[float]]:
        return self._embed(texts, prompt_name="text")

    def embeddings(self, text) -> Embedding:
        sentences = text.split(".")

        sentence_embeddings = [
            self.get_text_embedding(sentence) for sentence in sentences if sentence
        ]
        aggregated_embedding = np.mean(sentence_embeddings, axis=0)
        return aggregated_embedding
