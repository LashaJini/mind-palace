import inspect
import logging
import os
from typing import Mapping

from pkg.rpc import config

white = "\x1b[1;38;20m"
dark_grey = "\x1b[2;38;20m"
magenta = "\x1b[35;20m"
cyan = "\x1b[36;20m"
green = "\x1b[32;20m"
yellow = "\x1b[33;20m"
red = "\x1b[31;20m"
bold_red = "\x1b[31;1m"
color_reset = "\x1b[0m"


class STDFormatter(logging.Formatter):
    asctime = "{datecolor}{asctime}{color_end}"
    levelname = "{levelcolor}{levelname}{color_end}"
    servicename = "service={servicecolor}{service_name}{color_end}"
    caller = "[{filenamecolor}{caller_filename}:{line}{color_end}]"
    msg = "{message}"
    _format = f"{asctime} {levelname} {caller} {servicename} {cyan}>{color_reset} {msg}"

    FORMATS = {
        logging.DEBUG: {
            "levelcolor": magenta,
            "levelname": "DBG",
        },
        logging.INFO: {"levelcolor": green, "levelname": "INF"},
        logging.WARNING: {"levelcolor": yellow, "levelname": "WRN"},
        logging.ERROR: {"levelcolor": red, "levelname": "ERR"},
        logging.CRITICAL: {"levelcolor": bold_red, "levelname": "CRT"},
    }

    def format(self, record):
        formatter = logging.Formatter(self._format, style="{", datefmt="%H:%M:%S")

        formats = self.FORMATS.get(record.levelno, {})

        record.datecolor = dark_grey
        record.levelname = formats.get("levelname", "")
        record.filenamecolor = white
        record.servicecolor = cyan
        record.levelcolor = formats.get("levelcolor")

        record.color_end = color_reset

        return formatter.format(record)


class FileFormatter(logging.Formatter):
    def format(self, record):
        formatter = logging.Formatter(
            "{asctime} {levelname} [{caller_filename}:{lineno}] {message}",
            style="{",
            datefmt="%H:%M:%S",
        )

        return formatter.format(record)


class Logger(logging.Logger):
    def __init__(self, name=__name__, level=logging.NOTSET):
        super().__init__(name, level)

    def makeRecord(
        self,
        name: str,
        level: int,
        fn: str,
        lno: int,
        msg: object,
        args,
        exc_info,
        func: str | None = None,
        extra: Mapping[str, object] | None = None,
        sinfo: str | None = None,
    ) -> logging.LogRecord:
        if extra is None:
            s = inspect.stack()
            extra = {
                "caller_filename": s[0].filename,
                "line": s[0].lineno,
                "service_name": "Log",
            }

        record = super().makeRecord(
            name,
            level,
            fn,
            lno,
            msg,
            args,
            exc_info,
            func,
            extra,
            sinfo,
        )

        record.caller_filename = os.path.join(
            os.path.basename(os.path.dirname(str(extra["caller_filename"]))),
            os.path.basename(str(extra["caller_filename"])),
        )

        return record

    def db_info(self, msg, *args, **kwargs):
        extra = kwargs.get("extra")
        id = ""
        if extra is not None:
            id = extra.get("id")

        msg = f"{yellow}{msg} (tx={id}){color_reset}"
        super().info(msg, *args, **kwargs)

    def tx_info(self, msg, *args, **kwargs):
        extra = kwargs.get("extra")
        id = ""
        if extra is not None:
            id = extra.get("id")

        msg = f"{yellow}{msg} --- (tx={id}){color_reset}"
        super().info(msg, *args, **kwargs)


logging.setLoggerClass(Logger)
Log = logging.getLogger(__name__)
Log.setLevel(config.LOG_LEVEL)

console_handler = logging.StreamHandler()
console_handler.setFormatter(STDFormatter())
Log.addHandler(console_handler)

if config.MP_ENV != "test":
    file_handler = logging.FileHandler(config.LOG_FILEPATH, mode="a", encoding="utf-8")
    file_handler.setFormatter(FileFormatter())
    Log.addHandler(file_handler)
