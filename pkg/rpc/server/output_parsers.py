from typing import List
from llama_index.core.output_parsers.base import ChainableOutputParser
from pydantic import BaseModel


class Summary(BaseModel):
    """Data model for a summary."""

    summary: str


class SummaryParser(ChainableOutputParser):
    def __init__(self, verbose: bool = False):
        self.verbose = verbose

    def parse(self, output: str) -> Summary:
        if self.verbose:
            print(f"> Raw output: {output}")

        lines = output.split("\n")
        summary = ""
        for i in range(2, len(lines)):
            summary += lines[i] + "\n"

        return Summary(summary=summary)


class Keywords(BaseModel):
    """Data model for a keywords."""

    keywords: List[str]


class KeywordsParser(ChainableOutputParser):
    def __init__(self, verbose: bool = False):
        self.verbose = verbose

    def parse(self, output: str) -> Keywords:
        if self.verbose:
            print(f"> Raw output: {output}")

        keywords_part = output.split(":", 1)[1].strip()
        keywords = [keyword.strip().lower() for keyword in keywords_part.split(",")]

        return Keywords(keywords=keywords)
