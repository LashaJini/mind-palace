import uuid
from pkg.rpc.server.llm import llm
from pkg.rpc.server.vdb import Milvus
import gen.Palace_pb2 as pbPalace


class Addon:
    @classmethod
    def default(cls, original_id: str, input: str, client: Milvus):
        """input -> insert embeddings"""
        client.insert({"id": original_id, "input": input})

        return pbPalace.AddonResult(id="", data=None)

    @classmethod
    def summary(cls, original_id: str, input: str, client: Milvus):
        """input -> generate summary -> insert embeddings -> return summary"""
        result = [str(llm.gen_structured_summary(llm, input).dict().get("summary"))]

        id = str(uuid.uuid4())
        client.insert({"id": id, "input": input})

        return pbPalace.AddonResult(
            id=id, data={"output": pbPalace.Strings(value=result)}
        )

    @classmethod
    def keywords(cls, original_id: str, input: str, client: Milvus):
        """input -> generate keywords -> return keywords"""
        result = llm.gen_structured_keywords(llm, input).dict().get("keywords")

        return pbPalace.AddonResult(
            id="", data={"output": pbPalace.Strings(value=result)}
        )


AddonsDict = {
    "mind-palace-default": Addon.default,
    "mind-palace-resource-summary": Addon.summary,
    "mind-palace-resource-keywords": Addon.keywords,
}
