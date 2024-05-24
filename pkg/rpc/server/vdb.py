from typing import List, TypedDict
from pymilvus import MilvusClient, connections, DataType
from pkg.rpc.server.llm import Settings


class InputDataDict(TypedDict):
    id: str
    input: str


class Milvus:
    dim: int = 384

    def __init__(self, host: str, port: int, db_name: str, collection_name: str):
        self.host = host
        self.port = port
        self.db_name = db_name
        self.collection_name = collection_name

        connections.connect(host=self.host, port=self.port)

        self.client = MilvusClient(
            uri=f"http://{self.host}:{self.port}",
            db_name=self.db_name,
        )

        if not self.client.has_collection(collection_name):
            schema = self.client.create_schema(
                auto_id=False,
                enable_dynamic_field=True,
            )

            schema.add_field(
                field_name="id",
                datatype=DataType.VARCHAR,
                is_primary=True,
                max_length=64,
            )
            schema.add_field(
                field_name="vector", datatype=DataType.FLOAT_VECTOR, dim=Milvus.dim
            )

            index_params = self.client.prepare_index_params()
            index_params.add_index(field_name="id", index_type="INVERTED")
            index_params.add_index(
                field_name="vector",
                index_type="IVF_FLAT",
                metric_type="L2",
                params={"nlist": 1536},
            )

            print(f"Collection {collection_name} not found. Creating...")
            self.client.create_collection(
                collection_name=collection_name,
                schema=schema,
                index_params=index_params,
            )

    def insert(self, data: List[InputDataDict]):
        embedded_data = self._get_text_embedding(data)
        self.client.insert(collection_name=self.collection_name, data=embedded_data)

    def _get_text_embedding(self, data: List[InputDataDict]):
        embedded_data = []
        for item in data:
            new_item = {
                "id": item["id"],
                "vector": Settings.embed_model.get_text_embedding(item["input"]),
            }
            embedded_data.append(new_item)
        return embedded_data
