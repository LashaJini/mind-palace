from typing import Optional
import pkg.rpc.gen.Palace_pb2 as pbPalace
from pkg.rpc.services.addon.addons.abstract import Addon
from pkg.rpc.services.llm.llm import CustomLlamaCPP
from pkg.rpc.services.addon.output_parsers.default import Default, DefaultParser


class DefaultAddon(Addon):
    _parser: DefaultParser
    _output_model: Default
    _prompt_variables: dict

    def __init__(self, verbose=False, **kwargs):
        super().__init__(**kwargs)

        self._parser = DefaultParser(verbose=verbose)
        self._output_model = Default()
        self._prompt_variables = {}

    def prepare_input(self, user_input: str):
        self._output_model.default = user_input
        self._parser.default = user_input
        self._result = self._output_model.to_addon_result()

        return self

    def input(self, verbose=False) -> str:
        return ""

    def apply(
        self,
        llm: CustomLlamaCPP,
        verbose=False,
        **kwargs,
    ):
        """default, identity addon"""
        return self

    def finalize(self, result: Optional[pbPalace.AddonResult] = None, verbose=False):
        return self

    @property
    def output_model(self) -> Default:
        return self._output_model

    @property
    def parser(self) -> DefaultParser:
        return self._parser
