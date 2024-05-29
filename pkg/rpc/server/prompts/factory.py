from pkg.rpc.server.prompts.abstract import Prompts
from pkg.rpc.server.prompts.keywords import KeywordsPrompts
from pkg.rpc.server.prompts.default import DefaultPrompts
from pkg.rpc.server.prompts.summary import SummaryPrompts


class PromptsFactory:
    @classmethod
    def construct(cls, name: str) -> Prompts:
        if name not in prompts_registry:
            print(f"Unknown prompt {name}")

        prompts = prompts_registry[name]
        return prompts()


prompts_registry = {
    "mind-palace-default": DefaultPrompts,
    "mind-palace-resource-summary": SummaryPrompts,
    "mind-palace-resource-keywords": KeywordsPrompts,
}
