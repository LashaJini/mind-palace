import re
from typing import ClassVar, List

from pydantic import PrivateAttr

import pkg.rpc.gen.Palace_pb2 as pbPalace
from pkg.rpc.services.addon import addon_names
from pkg.rpc.loggers.palace import log
from pkg.rpc.services.addon.output_parsers.abstract import (
    OutputParser,
    CustomBaseModel,
)
from pkg.rpc.services.addon.prompts.keywords import KeywordsPrompts


class Keywords(CustomBaseModel):
    name: ClassVar[str] = addon_names.keywords
    _keywords: List[List[str]] = PrivateAttr([])
    _chunks: List[str] = PrivateAttr([])

    class Config:
        arbitrary_types_allowed = True

    def __init__(
        self, keywords: List[List[str]] = [], chunks: List[str] = [], **kwargs
    ):
        super().__init__(**kwargs)

        self._keywords = keywords
        self._chunks = chunks

    @property
    def keywords(self) -> List[List[str]]:
        return self._keywords

    @keywords.setter
    def keywords(self, keywords):
        self._keywords = keywords

    @property
    def chunks(self) -> List[str]:
        return self._chunks

    @chunks.setter
    def chunks(self, chunks: List[str]):
        self._chunks = chunks

    def to_addon_result(self) -> pbPalace.AddonResult:
        result: List[pbPalace.KeywordsResponse.KeywordChunk] = []

        if len(self.keywords) != len(self.chunks):
            raise ValueError("Keywords and chunks must be the same length")

        for keywords, chunk in zip(self.keywords, self.chunks):
            result.append(
                pbPalace.KeywordsResponse.KeywordChunk(
                    keywords=keywords,
                    chunk=chunk,
                )
            )

        response = pbPalace.AddonResponse(
            keywords_response=pbPalace.KeywordsResponse(list=result),
            success=self.success,
        )

        return pbPalace.AddonResult(map={Keywords.name: response})


class KeywordsParser(OutputParser):
    _chunks: List[str] = []

    def __init__(self, verbose: bool = False, **kwargs):
        format_start = KeywordsPrompts.format_start
        format_end = KeywordsPrompts.format_end
        group_name = "keywords"
        pattern = rf"{format_start}\s*(?P<{group_name}>.*?)\s*{format_end}"

        super().__init__(
            format_start=format_start,
            format_end=format_end,
            group_name=group_name,
            pattern=pattern,
            verbose=verbose,
            **kwargs,
        )

        self._chunks: List[str] = kwargs.get("chunks", [])

    def parse(self, output: str) -> Keywords:
        log.debug(f"Raw output: {output}")

        self.success = False
        keywords: List[List[str]] = []
        matches = re.finditer(self.pattern, output, re.DOTALL)

        for match in matches:
            chunk_keywords = match.group(self.group_name)
            chunk_keywords = [
                chunk_keyword.strip().lower()
                for chunk_keyword in chunk_keywords.split(",")
                if chunk_keyword
            ]
            keywords.append(chunk_keywords)

            # TODO: this is not a correct way to validate
            self.success = len(chunk_keywords) > 0

        if not self.success:
            e = (
                f"Keywords output format should be: {self.format_start} <{self.group_name}> {self.format_end}\n",
                f"Got: {output}",
            )
            raise ValueError(e)

        log.debug(f"Extracted keywords: {keywords}")

        return Keywords(keywords=keywords, chunks=self._chunks, success=self.success)

    @classmethod
    def construct_output(cls, **kwargs) -> str:
        """
        Args:
            input (List[str]): List of input text.
        """
        input = kwargs.get("input", [])
        return f"{KeywordsPrompts.format_start} {','.join(input)} {KeywordsPrompts.format_end}"

    @property
    def chunks(self) -> List[str]:
        return self._chunks

    @chunks.setter
    def chunks(self, chunks: List[str]):
        self._chunks = chunks
