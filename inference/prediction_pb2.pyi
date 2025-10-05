from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class PredictionRequest(_message.Message):
    __slots__ = ("message", "return_embedding")
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    RETURN_EMBEDDING_FIELD_NUMBER: _ClassVar[int]
    message: str
    return_embedding: bool
    def __init__(self, message: _Optional[str] = ..., return_embedding: bool = ...) -> None: ...

class PredictionResponse(_message.Message):
    __slots__ = ("label", "embedding")
    LABEL_FIELD_NUMBER: _ClassVar[int]
    EMBEDDING_FIELD_NUMBER: _ClassVar[int]
    label: int
    embedding: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, label: _Optional[int] = ..., embedding: _Optional[_Iterable[float]] = ...) -> None: ...

class EmbedRequest(_message.Message):
    __slots__ = ("message",)
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    message: str
    def __init__(self, message: _Optional[str] = ...) -> None: ...

class EmbedResponse(_message.Message):
    __slots__ = ("embedding",)
    EMBEDDING_FIELD_NUMBER: _ClassVar[int]
    embedding: _containers.RepeatedScalarFieldContainer[float]
    def __init__(self, embedding: _Optional[_Iterable[float]] = ...) -> None: ...
