from pkg.rpc.server.prompts.abstract import Prompts
from pkg.rpc.server.addons.factory import addons as registry


class PromptsFactory:
    @classmethod
    def construct(cls, name: str) -> Prompts:
        if name not in registry:
            raise ValueError(f"Unknown addon name {name}")

        prompts = registry[name].prompts
        return prompts()
