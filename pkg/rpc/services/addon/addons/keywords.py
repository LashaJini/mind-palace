from typing import Any, List, Optional
from llama_index.core import PromptTemplate
from llama_index.core.program import LLMTextCompletionProgram
from llama_index.core import Document
from llama_index.core.node_parser import SemanticSplitterNodeParser
from pydantic import Field
from transformers import AutoTokenizer, AutoModelForTokenClassification, Pipeline
from transformers import pipeline

import pkg.rpc.gen.Palace_pb2 as pbPalace
from pkg.rpc.services.addon.prompts.keywords import KeywordsPrompts
from pkg.rpc.services.addon.addons.abstract import Addon
from pkg.rpc.services.llm.llm import CustomLlamaCPP, EmbeddingModel
from pkg.rpc.services.addon.output_parsers.keywords import Keywords, KeywordsParser
from pkg.rpc.loggers.palace import log


ner_model_name = "dslim/bert-base-NER"


class KeywordsAddon(Addon):
    _parser: KeywordsParser
    _output_model: Keywords
    _prompt_variables: dict

    def __init__(self, verbose=False, **kwargs):
        super().__init__(**kwargs)

        self._parser = KeywordsParser(verbose=verbose)
        self._output_model = Keywords()
        self._prompt_variables = {}

    def prepare_input(self, user_input: str):
        chunks = self.semantic_chunks(user_input)

        self._output_model.chunks = chunks
        self._parser.chunks = chunks
        self._prompt_variables["total_chunks"] = len(chunks)

        return self

    def input(self, verbose=False) -> str:
        log.debug(f"Total Chunks {len(self._output_model.chunks)}")

        result = "\n\n".join(
            [
                f"<CHUNK {i+1}>" + chunk.strip() + f"\n</CHUNK {i+1}>"
                for i, chunk in enumerate(self._output_model.chunks)
            ]
        )
        return result

    def apply(
        self,
        llm: CustomLlamaCPP,
        verbose=False,
        **kwargs,
    ):
        if "max_keywords" not in kwargs:
            kwargs["max_keywords"] = KeywordsPrompts.default_max_keywords

        prompt = KeywordsPrompts().prompt(
            context_str=self.input(verbose),
            verbose=verbose,
            **kwargs,
        )
        program = LLMTextCompletionProgram(
            llm=llm,
            output_parser=self._parser,
            output_cls=Keywords,  # type:ignore
            prompt=PromptTemplate(prompt),
            verbose=verbose,
        )

        llm_output = program(verbose=verbose, **kwargs)
        result = Keywords.model_validate(llm_output)

        self._result = result.to_addon_result()

        return self

    def finalize(self, result: Optional[pbPalace.AddonResult] = None, verbose=False):
        _result = self._result if self._result is not None else result
        if _result is None:
            return self

        keywords_response: pbPalace.AddonResponse | None = _result.map.get(
            self._output_model.name
        )

        if keywords_response is not None and keywords_response.success:
            pipe = self.ner()

            for keywords_chunk in keywords_response.keywords_response.list:
                text = keywords_chunk.chunk
                llm_keywords = keywords_chunk.keywords
                ner_keywords: List[str] = []

                if text is None:
                    continue

                possible_entities = pipe(text)
                self.ner_join_words(possible_entities, ner_keywords)

                unique_keywords = list(set(ner_keywords + list(llm_keywords)))

                # llm sometimes generates way too many (duplicated) keyword sets.
                # If length of keyword sets is greater than total chunks, we take
                # first len(chunks) keyword sets. In this case, keywords may or may
                # not exist in corresponding chunk. This forces us to redo the
                # uniqueness check.
                #
                # issue may be related to token limit and/or context window size
                unique_existing_keywords = []
                for keyword in unique_keywords:
                    if keyword in text.lower():
                        unique_existing_keywords.append(keyword)

                keywords_chunk.CopyFrom(
                    pbPalace.KeywordsResponse.KeywordChunk(
                        keywords=unique_existing_keywords, chunk=text
                    )
                )

                log.debug(
                    f"Updated keywords: {unique_keywords if len(unique_keywords) > 0 else None}"
                )

        return self

    def ner(self) -> Pipeline:
        """Named entity recognition

        model output schema:
            List[{
                    'entity_group': 'LOC'|'ORG'|'PER'|'MISC'|'O',
                    'score': np.float32,
                    'word': str,
                    'start': int,
                    'end': int
                }]

        example output:
            [{
                'entity_group': 'LOC',
                'score': np.float32(0.9997018),
                'word': 'Paris',
                'start': 118,
                'end': 123
            }]

        model:
            https://huggingface.co/dslim/bert-base-NER

        When model returns "##sometext", this means that it is a continuation of the previous word
        """
        tokenizer = AutoTokenizer.from_pretrained(ner_model_name)
        model = AutoModelForTokenClassification.from_pretrained(ner_model_name)

        return pipeline(
            "ner",
            model=model,
            tokenizer=tokenizer,
            aggregation_strategy="simple",  # tries to join results into single entity
            device=-1,  # on cpu
            ignore_labels=["O"],  # useless in our case
        )

    def ner_join_words(self, features: Any, list: List[str]):
        if features is not None:
            for feature in features:
                f: str = feature.get("word", "")  # type:ignore
                lastIndexOfHashtag = f.rfind("#")
                if lastIndexOfHashtag != -1:
                    list[-1] += f[lastIndexOfHashtag + 1 :].lower()
                else:
                    list.append(f.lower())

    def semantic_chunks(
        self, text: str = Field(description="Text to split")
    ) -> List[str]:
        """Split text into semantic chunks"""
        embed_model = EmbeddingModel()
        documents = [Document(text=text)]
        splitter = SemanticSplitterNodeParser(
            buffer_size=2,
            breakpoint_percentile_threshold=65,
            embed_model=embed_model,
        )
        nodes = splitter.get_nodes_from_documents(documents)
        return [node.get_content() for node in nodes]

    @property
    def output_model(self) -> Keywords:
        return self._output_model

    @property
    def parser(self) -> KeywordsParser:
        return self._parser
