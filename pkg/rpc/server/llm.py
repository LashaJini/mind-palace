from pathlib import Path
from typing import (
    Callable,
    ClassVar,
    List,
    Optional,
)
from llama_index.core.base.embeddings.base import Embedding
from llama_index.core.base.llms.types import LLMMetadata
from llama_index.core.tools import ToolMetadata
from llama_index.core.types import BaseOutputParser
from llama_index.llms.openai.base import FunctionCallingLLM
import numpy as np

from llama_index.embeddings.huggingface import HuggingFaceEmbedding
from llama_index.llms.llama_cpp import LlamaCPP
from pydantic import BaseModel, PrivateAttr

from pkg.rpc.server import logger
from pkg.rpc.server.output_parsers.abstract import CustomBaseModel

# TODO: read from config
user_home = str(Path.home())
model_path = (
    user_home
    + "/.cache/lm-studio/models/lmstudio-community/Meta-Llama-3-8B-Instruct-GGUF/Meta-Llama-3-8B-Instruct-IQ3_M.gguf"
)


class CustomLlamaCPP(LlamaCPP, FunctionCallingLLM, extra="allow"):
    context_size: ClassVar[int] = 8000
    _context_window: ClassVar[int] = 3900
    _max_new_tokens: ClassVar[int] = 1024
    program: Callable = PrivateAttr()
    tool_choices: List[ToolMetadata] = PrivateAttr()
    _output_parser: Optional[BaseOutputParser] = PrivateAttr()
    output_cls: BaseModel = PrivateAttr()

    def __init__(self, verbose: bool = False, **kwargs):
        super().__init__(
            model_path=model_path,
            temperature=0.1,
            max_new_tokens=CustomLlamaCPP._max_new_tokens,
            context_window=CustomLlamaCPP._context_window,
            verbose=verbose,
            **kwargs,
        )
        self._output_parser = None

    @property
    def metadata(self) -> LLMMetadata:
        return LLMMetadata(
            context_window=CustomLlamaCPP._context_window,
            num_output=CustomLlamaCPP._max_new_tokens,
            is_chat_model=False,
            is_function_calling_model=True,
            model_name="CustomLlamaCPP",
        )

    def set_output_cls(self, output_cls: CustomBaseModel):
        self.output_cls = output_cls

    def set_output_parser(self, output_parser: BaseOutputParser):
        self._output_parser = output_parser

    def set_text_completion_program(self, program: Callable):
        self.program = program

    def set_tool_choices(self, tool_choices: List[ToolMetadata]):
        self.tool_choices = tool_choices

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
        available_tokens = 1

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

    def embeddings(self, text) -> Embedding:
        sentences = text.split(".")

        sentence_embeddings = [
            self.get_text_embedding(sentence) for sentence in sentences if sentence
        ]
        aggregated_embedding = np.mean(sentence_embeddings, axis=0)
        return aggregated_embedding
