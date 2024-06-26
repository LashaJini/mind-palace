from pkg.rpc.server.addons.abstract import Addon
from pkg.rpc.server.addons.scheme import addons


class AddonFactory:
    @classmethod
    def construct(cls, name: str) -> Addon:
        if name not in addons:
            print(f"Unknown addon {name}")

        addon = addons[name].get("addon")
        return addon()
