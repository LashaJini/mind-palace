import logging
import os

from pkg.rpc.server.config import LOG_FILEPATH


class STDFormatter(logging.Formatter):
    white = "\x1b[1;38;20m"
    dark_grey = "\x1b[2;38;20m"
    grey = "\x1b[38;20m"
    green = "\x1b[32;20m"
    yellow = "\x1b[33;20m"
    red = "\x1b[31;20m"
    bold_red = "\x1b[31;1m"
    reset = "\x1b[0m"
    _format = "{datecolor}{asctime}{color_end} {levelcolor}{levelname}{color_end} [{filenamecolor}{relative_filename}:{lineno}{color_end}] {message}"

    FORMATS = {
        logging.DEBUG: {
            "levelcolor": grey,
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

        record.datecolor = self.dark_grey
        record.levelname = formats.get("levelname", "")
        record.filenamecolor = self.white
        record.relative_filename = os.path.join(
            os.path.basename(os.path.dirname(record.pathname)),
            os.path.basename(record.pathname),
        )
        record.levelcolor = formats.get("levelcolor")

        record.color_end = self.reset

        return formatter.format(record)


log = logging.getLogger(__name__)
log.setLevel(logging.DEBUG)

console_handler = logging.StreamHandler()
console_handler.setFormatter(STDFormatter())
log.addHandler(console_handler)


class FileFormatter(logging.Formatter):
    def format(self, record):
        formatter = logging.Formatter(
            "{asctime} {levelname} [{relative_filename}:{lineno}] {message}",
            style="{",
            datefmt="%H:%M:%S",
        )

        return formatter.format(record)


file_handler = logging.FileHandler(LOG_FILEPATH, mode="a", encoding="utf-8")
file_handler.setFormatter(FileFormatter())
log.addHandler(file_handler)
