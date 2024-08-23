from typing import Dict, Type
from pkg.rpc.services.addon.addons.abstract import Addon

from pkg.rpc.services.addon.output_parsers.abstract import OutputParser
from pkg.rpc.services.addon.output_parsers.default import DefaultParser
from pkg.rpc.services.addon.output_parsers.keywords import KeywordsParser
from pkg.rpc.services.addon.output_parsers.summary import SummaryParser
from pkg.rpc.services.addon.prompts.abstract import Prompts
from pkg.rpc.services.addon.prompts.default import DefaultPrompts
from pkg.rpc.services.addon.prompts.keywords import KeywordsPrompts
from pkg.rpc.services.addon.prompts.summary import SummaryPrompts
from pkg.rpc.services.addon.addons.default import DefaultAddon
from pkg.rpc.services.addon.addons.summary import SummaryAddon
from pkg.rpc.services.addon.addons.keywords import KeywordsAddon
from pkg.rpc.services.addon.addon_names import default, summary, keywords


class AddonScheme:
    addon: Type[Addon]
    prompts: Type[Prompts]
    parser: Type[OutputParser]
    name: str

    def __init__(self, **kwargs):
        self.addon = kwargs["addon"]
        self.prompts = kwargs["prompts"]
        self.parser = kwargs["parser"]
        self.name = kwargs["name"]


addons: Dict[str, AddonScheme] = {
    default: AddonScheme(
        addon=DefaultAddon,
        prompts=DefaultPrompts,
        parser=DefaultParser,
        name=default,
    ),
    summary: AddonScheme(
        addon=SummaryAddon,
        prompts=SummaryPrompts,
        parser=SummaryParser,
        name=summary,
    ),
    keywords: AddonScheme(
        addon=KeywordsAddon,
        prompts=KeywordsPrompts,
        parser=KeywordsParser,
        name=keywords,
    ),
}
