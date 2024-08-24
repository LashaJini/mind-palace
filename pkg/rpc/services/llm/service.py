import pkg.rpc.gen.Palace_pb2 as pbPalace
import pkg.rpc.gen.SharedTypes_pb2 as pbShared
from pkg.rpc.services.llm.llm import CustomLlamaCPP, EmbeddingModel


class LLMService:
    def __init__(self, llm: CustomLlamaCPP):
        self.llm = llm

    def TokenSize(self, request: pbPalace.Text, context) -> pbPalace.Size:
        return pbPalace.Size(size=self.llm.token_size(request.text))

    def CalculateAvailableTokens(
        self, request: pbPalace.DecrementList, context
    ) -> pbPalace.Size:
        decrements = [size.size for size in request.sizes]
        return pbPalace.Size(
            size=self.llm.calculate_available_tokens(decrements=decrements)
        )

    def GetConfig(self, request: pbShared.Empty, context) -> pbPalace.LLMConfig:
        return pbPalace.LLMConfig(
            context_size=self.llm.config.context_size,
            context_window=self.llm.config.context_window,
            max_new_tokens=self.llm.config.max_new_tokens,
            map=self.llm.config.kwargs,
        )

    def SetConfig(self, request: pbPalace.LLMConfig, context):
        if request.map is not None:
            self.llm.config.update(**request.map)

        return pbShared.Empty()

    def Ping(self, request, context):
        return pbShared.Empty()


class EmbeddingModelService:
    def __init__(self, embedding_model: EmbeddingModel):
        self.embedding_model = embedding_model

    def CalculateEmbeddings(
        self, request: pbPalace.Text, context
    ) -> pbPalace.Embeddings:
        return pbPalace.Embeddings(
            embedding=self.embedding_model.embeddings(request.text)
        )

    def GetConfig(
        self, request: pbShared.Empty, context
    ) -> pbPalace.EmbeddingModelConfig:
        return pbPalace.EmbeddingModelConfig(
            model_name=self.embedding_model._model_name,
            max_length=self.embedding_model._max_length,
            dimension=self.embedding_model.dimension,
            metric_type=self.embedding_model.metric_type,
        )
