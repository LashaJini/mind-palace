from typing import Dict, Type, TypedDict
from pkg.rpc.server.addons.abstract import Addon

from pkg.rpc.server.output_parsers.abstract import OutputParser
from pkg.rpc.server.output_parsers.default import DefaultParser
from pkg.rpc.server.output_parsers.keywords import KeywordsParser
from pkg.rpc.server.output_parsers.summary import SummaryParser
from pkg.rpc.server.prompts.abstract import Prompts
from pkg.rpc.server.prompts.default import DefaultPrompts
from pkg.rpc.server.prompts.keywords import KeywordsPrompts
from pkg.rpc.server.prompts.summary import SummaryPrompts
from pkg.rpc.server.addons.default import DefaultAddon
from pkg.rpc.server.addons.summary import SummaryAddon
from pkg.rpc.server.addons.keywords import KeywordsAddon


class AddonScheme(TypedDict):
    addon: Type[Addon]
    prompts: Type[Prompts]
    parser: Type[OutputParser]
    skip: bool


addons: Dict[str, AddonScheme] = {
    "mind-palace-default": {
        "addon": DefaultAddon,
        "prompts": DefaultPrompts,
        "parser": DefaultParser,
        "skip": True,
    },
    "mind-palace-resource-summary": {
        "addon": SummaryAddon,
        "prompts": SummaryPrompts,
        "parser": SummaryParser,
        "skip": False,
    },
    "mind-palace-resource-keywords": {
        "addon": KeywordsAddon,
        "prompts": KeywordsPrompts,
        "parser": KeywordsParser,
        "skip": False,
    },
}
