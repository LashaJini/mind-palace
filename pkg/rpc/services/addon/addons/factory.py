from pkg.rpc.services.addon.addons.abstract import Addon
from pkg.rpc.services.addon.addons.scheme import addons


class AddonFactory:
    @classmethod
    def construct(cls, name: str) -> Addon:
        if name not in addons:
            raise ValueError(f"Unknown addon {name}")

        addon = addons[name].addon
        return addon()
