from pkg.rpc.server.addons.abstract import Addon

from pkg.rpc.server.addons.default import DefaultAddon
from pkg.rpc.server.addons.summary import SummaryAddon
from pkg.rpc.server.addons.keywords import KeywordsAddon


class AddonFactory:
    @classmethod
    def construct(cls, name: str) -> Addon:
        if name not in addon_registry:
            print(f"Unknown addon {name}")

        addon = addon_registry[name]
        return addon()


addon_registry = {
    "mind-palace-default": DefaultAddon,
    "mind-palace-resource-summary": SummaryAddon,
    "mind-palace-resource-keywords": KeywordsAddon,
}
