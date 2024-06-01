from llama_index.core import PromptTemplate
from llama_index.core.program import LLMTextCompletionProgram

import gen.Palace_pb2 as pbPalace

from pkg.rpc.server.prompts.keywords import KeywordsPrompts
from pkg.rpc.server.addons.abstract import Addon
from pkg.rpc.server.llm import CustomLlamaCPP
from pkg.rpc.server.vdb import Milvus
from pkg.rpc.server.output_parsers.keywords import Keywords, KeywordsParser


class KeywordsAddon(Addon):
    def apply(
        self,
        id: str,
        input: str,
        llm: CustomLlamaCPP,
        client: Milvus,
        verbose=False,
        **kwargs,
    ):
        """input -> generate keywords -> return keywords"""
        if "max_keywords" not in kwargs:
            kwargs["max_keywords"] = KeywordsPrompts.default_max_keywords

        prompt = KeywordsPrompts().prompt(text=input, verbose=verbose, **kwargs)
        program = LLMTextCompletionProgram.from_defaults(
            llm=llm,
            output_parser=KeywordsParser(verbose=verbose),
            output_cls=Keywords,  # type:ignore
            prompt=PromptTemplate(prompt),
            verbose=verbose,
        )

        result = (
            program(context_str=input, verbose=verbose, **kwargs).dict().get("value")
        )

        return pbPalace.AddonResult(
            id="", data={"output": pbPalace.Strings(value=result)}
        )
