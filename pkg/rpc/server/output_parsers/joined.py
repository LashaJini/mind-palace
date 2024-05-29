import re
from typing import List, Dict

from pkg.rpc.server.output_parsers.abstract import OutputParser, CustomBaseModel
from pkg.rpc.server.output_parsers.factory import OutputParserFactory


class Joined(CustomBaseModel):
    """Data model for joined data models"""

    value: Dict[str, List[str]]


class JoinedParser(OutputParser):
    def __init__(self, addons: List[str], verbose=False):
        self.verbose = verbose
        self.addons = addons

    def parse(self, output: str) -> Joined:
        if self.verbose:
            print(f"> Raw output: {output}")

        patterns = []
        parsers: List[OutputParser] = []
        for addon_name in self.addons:
            parser, skip = OutputParserFactory.construct(addon_name)
            if parser is None or skip:
                continue

            parser_instance: OutputParser = parser(verbose=self.verbose)  # type: ignore
            parsers.append(parser_instance)
            patterns.append(parser_instance.pattern)

        pattern = ".*".join(patterns)

        if self.verbose:
            print(f"> Output parser pattern: {pattern}")

        match = re.search(pattern, output, re.DOTALL)

        results: Dict[str, List[str]] = {}

        if match:
            outputs = []
            for parser in parsers:
                output = match.group(parser.group_name)
                outputs.append(output)

            for i, output in enumerate(outputs):
                parser = parsers[i]
                result = parser.parse(
                    f"{parser.format_start} {output} {parser.format_end}"
                )
                results[result.name] = result.value

            if self.verbose:
                print(f"> Output parser results {results}")
        else:
            print("Failed to parse output.")

        return Joined(value=results)
