from typing import List, Optional
from llama_index.core import PromptTemplate
from llama_index.core.program import LLMTextCompletionProgram


from pkg.rpc.services.addon.addons.default import DefaultAddon
from pkg.rpc.services.addon.addons.factory import AddonFactory
from pkg.rpc.services.addon.addons.keywords import KeywordsAddon
import pkg.rpc.gen.Palace_pb2 as pbPalace
from pkg.rpc.services.addon.addons.abstract import Addon
from pkg.rpc.services.llm.llm import CustomLlamaCPP
from pkg.rpc.services.addon.prompts.joined import JoinedPrompts
from pkg.rpc.services.addon.output_parsers.joined import Joined, JoinedParser


class JoinedAddons(Addon):
    _parser: JoinedParser
    _output_model: Joined
    _addons: dict[str, Addon]
    _prompt_variables: dict

    def __init__(self, names: List[str], verbose=False, **kwargs):
        super().__init__(**kwargs)

        self._addons = {}
        for name in names:
            self._addons[name] = AddonFactory.construct(name)

        self._parser = JoinedParser(
            parsers=[addon.parser for addon in self._addons.values()], verbose=verbose
        )
        self._output_model = Joined()
        self._prompt_variables = {}

    def prepare_input(self, user_input: str):
        for _, addon in self._addons.items():
            addon.prepare_input(user_input=user_input)
            self._prompt_variables = {
                **self._prompt_variables,
                **addon._prompt_variables,
            }

        return self

    # TODO: dis is ugly
    def input(self, verbose=False) -> str:
        _result = []
        for _, addon in self._addons.items():
            if isinstance(addon, KeywordsAddon):
                addon_input = addon.input(verbose)
                _result.append(addon_input)

        if len(_result) == 0:
            for _, addon in self._addons.items():
                if isinstance(addon, DefaultAddon):
                    _result.append(addon.output_model.default)

        return "\n\n".join(_result)

    def apply(
        self,
        llm: CustomLlamaCPP,
        verbose=False,
        **kwargs,
    ):
        prompt = JoinedPrompts().prompt(
            context_str=self.input(verbose),
            verbose=verbose,
            **self._prompt_variables,
            **kwargs,
        )
        program = LLMTextCompletionProgram(
            llm=llm,
            output_parser=self._parser,
            output_cls=Joined,  # type:ignore
            prompt=PromptTemplate(prompt),
            verbose=verbose,
        )

        llm_output = program(verbose=verbose, **kwargs)
        result = Joined.model_validate(llm_output)

        self._result = result.to_addon_result()

        return self

    def finalize(self, result: Optional[pbPalace.AddonResult] = None, verbose=False):
        for _, addon in self._addons.items():
            addon.finalize(self._result, verbose)

        self._addons = {}
        return self

    @property
    def output_model(self) -> Joined:
        return self._output_model

    @property
    def parser(self) -> JoinedParser:
        return self._parser

    @property
    def prompt_variables(self) -> dict:
        return self._prompt_variables
