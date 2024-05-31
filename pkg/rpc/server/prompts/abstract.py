from abc import ABC, abstractmethod
from typing import TypedDict
from llama_index.llms.llama_cpp.llama_utils import DEFAULT_SYSTEM_PROMPT

from pkg.rpc.server.llm import CustomLlamaCPP


class JoinableTemplateDict(TypedDict):
    instructions: str
    format: str


class JoinableTemplate:
    def __init__(self, data: JoinableTemplateDict):
        self._data = data

    @property
    def instructions(self) -> str:
        return self._data["instructions"]

    @property
    def format(self) -> str:
        return self._data["format"]


class Prompts(ABC):
    sys_prompt = DEFAULT_SYSTEM_PROMPT

    @abstractmethod
    def standalone_template(self, verbose=False, **kwargs) -> str:
        pass

    @abstractmethod
    def standalone_template_token_count(self, llm: CustomLlamaCPP) -> int:
        pass

    @abstractmethod
    def prompt(self, text: str, verbose: bool, **kwargs) -> str:
        pass

    @abstractmethod
    def joinable_template(self, **kwargs) -> JoinableTemplate:
        pass

    def joinable_template_token_count(self, llm: CustomLlamaCPP):
        tmpl = self.joinable_template()
        text = tmpl.instructions + tmpl.format
        return llm.token_size(text=text)

    def _standalone_template(self, tmpl: str, verbose=False, **kwargs):
        result = f"<|system|>\n{Prompts.sys_prompt}</s>\n<|user|>\n{tmpl}</s>\n<|assistant|>\n"
        if verbose:
            print(result)

        return result

    @classmethod
    def system_prompt_token_count(cls, llm: CustomLlamaCPP):
        return llm.token_size(text=Prompts.sys_prompt)
