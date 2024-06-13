from typing import ClassVar, List
from pydantic import ConfigDict

from pkg.rpc.server import addon_names
from pkg.rpc.server.output_parsers.abstract import OutputParser, CustomBaseModel


class Default(CustomBaseModel):
    model_config = ConfigDict(from_attributes=True)

    value: List[str]
    name: ClassVar[str] = addon_names.default

    def get_value(self) -> List[str]:
        return self.value


class DefaultParser(OutputParser):
    def __init__(self, verbose: bool = False):
        format_start = ""
        format_end = ""
        group_name = ""
        pattern = ""
        skip = True

        super().__init__(
            format_start=format_start,
            format_end=format_end,
            group_name=group_name,
            pattern=pattern,
            skip=skip,
            verbose=verbose,
        )

    def parse(self, output: str) -> Default:
        self.success = True
        return Default(value=[output])

    @classmethod
    def construct_output(cls, **kwargs) -> str:
        """For testing purposes."""
        return ""
