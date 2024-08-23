from pkg.rpc.services.addon.prompts.abstract import Prompts
from pkg.rpc.services.addon.addons.factory import addons as registry


class PromptsFactory:
    @classmethod
    def construct(cls, name: str) -> Prompts:
        if name not in registry:
            raise ValueError(f"Unknown addon name {name}")

        prompts = registry[name].prompts
        return prompts()
