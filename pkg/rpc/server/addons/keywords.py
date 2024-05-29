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
        program = LLMTextCompletionProgram.from_defaults(
            llm=llm,
            output_parser=KeywordsParser(verbose=verbose),
            output_cls=Keywords,  # type:ignore
            prompt_template_str=KeywordsPrompts().standalone_template(
                verbose=verbose,
            ),
            verbose=verbose,
        )

        if "max_keywords" not in kwargs:
            kwargs["max_keywords"] = KeywordsPrompts.default_max_keywords

        result = (
            program(context_str=input, verbose=verbose, **kwargs).dict().get("value")
        )

        return pbPalace.AddonResult(
            id="", data={"output": pbPalace.Strings(value=result)}
        )
