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
from pkg.rpc.server.addon_names import default, summary, keywords


class AddonScheme(TypedDict):
    addon: Type[Addon]
    prompts: Type[Prompts]
    parser: Type[OutputParser]


addons: Dict[str, AddonScheme] = {
    default: {
        "addon": DefaultAddon,
        "prompts": DefaultPrompts,
        "parser": DefaultParser,
    },
    summary: {
        "addon": SummaryAddon,
        "prompts": SummaryPrompts,
        "parser": SummaryParser,
    },
    keywords: {
        "addon": KeywordsAddon,
        "prompts": KeywordsPrompts,
        "parser": KeywordsParser,
    },
}
