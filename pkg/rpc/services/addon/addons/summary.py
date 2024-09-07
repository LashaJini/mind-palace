from typing import Optional
from llama_index.core import PromptTemplate
from llama_index.core.program import LLMTextCompletionProgram

import pkg.rpc.gen.Palace_pb2 as pbPalace
from pkg.rpc.services.addon.addons.abstract import Addon
from pkg.rpc.services.llm.llm import CustomLlamaCPP
from pkg.rpc.services.addon.prompts.summary import SummaryPrompts
from pkg.rpc.services.addon.output_parsers.summary import Summary, SummaryParser


class SummaryAddon(Addon):
    _parser: SummaryParser
    _output_model: Summary
    _input: str
    _prompt_variables: dict

    def __init__(self, verbose=False, **kwargs):
        super().__init__(**kwargs)

        self._parser = SummaryParser(verbose=verbose)
        self._output_model = Summary()
        self._input = ""
        self._prompt_variables = {}

    def prepare_input(self, user_input: str):
        self._input = user_input
        return self

    def input(self, verbose=False) -> str:
        return self._input

    def apply(
        self,
        llm: CustomLlamaCPP,
        verbose=False,
        **kwargs,
    ):
        """input -> generate summary -> insert embeddings -> return summary"""
        prompt = SummaryPrompts().prompt(
            context_str=self.input(verbose),
            verbose=verbose,
            **kwargs,
        )
        program = LLMTextCompletionProgram(
            llm=llm,
            output_parser=self._parser,
            output_cls=Summary,  # type:ignore
            prompt=PromptTemplate(prompt),
            verbose=verbose,
        )

        llm_output = program(verbose=verbose, **kwargs)
        result = Summary.model_validate(llm_output)

        self._result = result.to_addon_result()

        return self

    def finalize(self, result: Optional[pbPalace.AddonResult] = None, verbose=False):
        return self

    @property
    def output_model(self) -> Summary:
        return self._output_model

    @property
    def parser(self) -> SummaryParser:
        return self._parser
